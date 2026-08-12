# STA-100 Agent 分类与页面调用链路说明

版本：V1.0

日期：2026-08-12

实现范围：`sta100-web`、`openclaw-orchestrator`、`openclaw`

## 1. 设计结论

采用“Go 负责流程控制，OpenClaw Agent 负责语义处理”的分层方案。页面搜索不直接绑定 24 个业务 Agent 中的任意一个，也不让一个 Agent 自主递归调度全部 Agent。Go 服务掌握本地数据检索、路由、并发、超时、审计和降级；两个页面不可见的系统 Agent 分别负责本地证据整理和最终结果汇总；24 个业务 Agent 保持独立专业职责。

这一方案满足以下要求：

1. 页面搜索能力独立于 24 个业务 Agent，但能够识别并分发给已有业务 Agent。
2. 同一套编排能力可通过 `openclaw-orchestrator` 和清单移植到其它 Go 项目。
3. 页面不要求用户理解本地知识库、RAG 或联网检索的技术差异。
4. 数据冲突不被模型静默覆盖，所有值连同来源、记录编号和时间并列保留。
5. 任一模型阶段失败时，已有本地证据和成功结果仍可返回，并明确标记为部分结果。

## 2. Agent 分类

| 分类 | 数量 | 页面可见 | 职责 | 禁止事项 |
|---|---:|---|---|---|
| OpenClaw 默认 Agent `main` | 1 | 否 | OpenClaw 自身默认入口 | 不作为 STA-100 业务路由目标 |
| Knowledge Agent `sta100-knowledge` | 1 | 否 | 整理 Go 已检索出的本地证据，保留来源、记录 ID、更新时间和冲突 | 不联网、不自行检索、不调度其它 Agent |
| Coordinator Agent `sta100-coordinator` | 1 | 否 | 整合用户问题、本地证据和业务 Agent 返回内容，形成页面统一答复 | 不递归调用 Agent、不声称使用未返回的数据 |
| 业务 Agent | 24 | 是 | 按出口、支付、物流、单据、产品、供应商、客户、市场等领域给出专业结果 | 不替代 Go 的权限、路由、事务和审计 |

OpenClaw 实际注册数量为 27：`main + 2 个系统 Agent + 24 个业务 Agent`。`GET /api/v1/agents` 和页面 Agent 管理只返回 24 个业务 Agent，避免客户误操作系统编排角色。

## 3. 统一调用链路

```text
页面输入 page、feature、message、context
  -> POST /api/v1/assistant/query
  -> Go 校验登录、写请求头、消息长度、会话键和 targetAgent
  -> Go 从 SQLite 检索客户、供应商、产品和私有文件元数据
  -> Knowledge Agent 整理本地证据
  -> Go 根据 feature 选择业务 Agent，并行调用
  -> Coordinator Agent 接收原始证据、证据摘要和业务 Agent 输出
  -> Go 返回统一文本、结构化记录、冲突、证据、调用 Agent 和阶段状态
  -> 页面统一展示，不按“本地/联网”拆成两个结果区
```

### 3.1 为什么由 Go 路由

- 可对 Agent ID 使用固定白名单，阻止页面调用系统 Agent 或未知 Agent。
- 可限制消息长度、会话键格式、总超时和并发数量。
- 可记录页面、功能、证据数、实际成功 Agent 和是否降级。
- 可避免 Agent 自主循环调用、重复计费和不可预测的递归链路。
- 本地数据库查询无需消耗模型 Token，且字段过滤和权限更可靠。

### 3.2 数据处理顺序

1. **本地检索**：当前真实检索客户、供应商、产品和私有文件元数据，最多返回 30 条证据；无关键词命中时最多补充 5 条可用客户记录。
2. **本地整理**：Knowledge Agent 只接收 Go 找到的证据，不能将推断写成本地事实。
3. **专业分析**：Go 按页面功能选择 1 到 2 个业务 Agent，并行执行以控制总耗时。
4. **统一汇总**：Coordinator 接收所有实际返回内容，输出面向当前页面的答复。
5. **页面展示**：响应同时包含文字、结构化 `items`、`evidence`、`conflicts`、`usedAgents`、`agentOutputs` 和 `pipeline`，便于验收和问题追踪。

## 4. 页面与 Agent 路由

| 页面/功能 | feature | 业务 Agent | 作用 | 当前数据边界 |
|---|---|---|---|---|
| 概览/OEM 工厂匹配 | `oem-match` | `supplier-aggregator` + `brand-value-crawler` | 聚合供应商能力、品牌与公开情报，形成候选和匹配说明 | SQLite 供应商可检索；正式工厂分类、评分、联网工具待提供 |
| 概览/客户统一搜索 | `customer-search` | `customer-measurement-agent` + `market-analyzer` | 整理客户字段、联系方式和市场判断 | SQLite 客户真实可查；私有正文与联网补充待接 |
| 概览/本地客户发现 | `customer-discovery` | `customer-measurement-agent` + `market-analyzer` | 按单一国家、城市和类型筛选后形成客户洞察 | 当前可筛 SQLite；外部候选源待客户确认 |
| 客户/详情分析 | `customers` | `customer-measurement-agent` + `market-analyzer` | 结合客户档案、交易关系和输入生成跟进建议 | 档案和沟通记录为本地事实；模型建议须复核 |
| 报价单 | `quote` | `export-agent` + `payment-advisor` | 分析出口条款、报价结构和支付风险 | 报价数据已持久化；正式邮件与 PDF 模板待提供 |
| 订单/物流 | `order` / `logistics` | `shipping-eta` + `inventory-agent` | 分析交付、库存和物流事项 | 订单和库存为本地数据；实时物流平台未接入 |
| 单据 | `document` | `invoice-agent` + `export-agent` | 检查单据字段和出口业务要求 | 单据草稿已保存；正式模板、签章和导出待提供 |
| 产品库 | `product` | `compatibility-agent` + `design-advisor` | 分析规格兼容性和产品设计建议 | 产品表可查；技术规则库和私有正文待接 |
| 行业新闻 | `news` | `market-analyzer` + `country-advisor` | 归纳市场和国家影响 | 当前新闻为缓存/种子；真实来源白名单待确认 |
| 智能体/直接聊天 | `context.targetAgent` | 用户选中的 1 个业务 Agent | 保持原 24 个专业智能体直接对话能力 | 目标必须属于业务 Agent 白名单 |
| 未识别功能 | 其它值 | `market-analyzer` | 提供通用分析降级路由 | 后续新增功能应补充显式路由 |

## 5. 数据与冲突规则

页面不再给用户提供“本地知识库/联网检索/RAG”的来源模式选择，也不按来源拆分结果卡片，但后端不能丢失来源信息。

每条本地证据至少包含：

| 字段 | 作用 |
|---|---|
| `id` | 本地记录或文件编号，可回查原始记录 |
| `entity` | 客户、供应商、产品或私有文件元数据类型 |
| `title` | 记录名称 |
| `content` | 传给 Agent 的必要字段摘要 |
| `updatedAt` | 数据更新时间 |
| `source` | 当前实际来源，例如 `STA-100 SQLite` |

冲突判定和展示遵循以下规则：

1. 同一实体出现不同内容时写入 `conflicts`。
2. Coordinator 提示词强制要求并列展示，不平均、不覆盖、不静默选取。
3. 页面可不按来源分区，但详情和证据区域必须允许查看来源、记录编号和时间。
4. 模型推断必须与数据库事实区分；没有证据时不得描述为“本地检索结果”。
5. 所有 AI 汇总适用全局声明“内容有AI生成，请仔细甄别”。

## 6. 接口契约

### 6.1 请求

`POST /api/v1/assistant/query`

```json
{
  "page": "overview",
  "feature": "customer-search",
  "message": "查找德国有邮箱的经销商",
  "sessionKey": "overview-customer-001",
  "context": {
    "hasContact": true,
    "country": "德国"
  }
}
```

约束：必须登录；请求头需包含 `X-STA100-Request: 1`；`message` 必填且不超过 32 KiB；`sessionKey` 最长 80 字符；`targetAgent` 只能是 24 个业务 Agent 之一。

### 6.2 响应

| 字段 | 含义 |
|---|---|
| `text` | Coordinator 汇总文本；失败时为 Go 生成的本地降级说明 |
| `items` | 页面可直接展示的结构化客户等记录 |
| `evidence` | 实际检索证据及来源 |
| `conflicts` | 检出的冲突组及“并列展示”规则 |
| `usedAgents` | 实际成功返回的 Agent，不包含失败调用 |
| `agentOutputs` | 每个 Agent 的文本、runId 或受控错误 |
| `pipeline` | `local-retrieval`、`knowledge-agent`、`domain-agents`、`coordinator-agent` 阶段状态 |
| `aiGenerated` | 固定为 `true`，用于页面声明和审计 |
| `partial` | 任一模型阶段失败时为 `true` |
| `todo` | 尚缺客户资料的能力边界 |

## 7. 失败与降级

| 失败点 | 接口行为 | 页面行为 |
|---|---|---|
| 本地无命中 | 返回空证据或有限兜底记录，不伪造数据 | 显示无本地匹配，可继续查看 Agent 建议 |
| Knowledge Agent 失败 | 原始 Go 证据直接交给 Coordinator，`partial=true` | 标明证据整理阶段失败 |
| 部分业务 Agent 失败 | 保留成功 Agent 输出，失败项写入阶段状态 | 显示部分结果和实际使用 Agent |
| Coordinator 失败 | Go 根据本地证据和成功 Agent 数生成降级摘要 | 不把摘要标为完整智能结论 |
| 模型额度不足 | OpenClaw 返回失败，接口不暴露密钥或底层敏感错误 | 提示模型服务不可用，可稍后重试 |
| 私有正文未解析 | 只返回文件元数据，并在证据中明确说明 | 不声称已读取文件正文 |

当前 DeepSeek 请求已经真实到达提供商，但账户余额不足，Knowledge Agent 和 Coordinator 的生成内容尚不能完成端到端验收。接口的部分结果降级已验证正常，恢复额度后需补做完整内容质量验收。

## 8. 客户沟通记录规则

客户沟通记录属于不可变历史，不交给模型作为可修改内容：

- `GET /api/v1/accounts/{id}/communications`：按发生时间倒序读取。
- `POST /api/v1/accounts/{id}/communications`：追加记录，`createdBy` 取当前登录用户。
- 不提供 `PATCH` 和 `DELETE`；其它方法返回 `405 COMMUNICATION_IMMUTABLE`。
- 客户归档后仍允许查询历史，新增返回 `409 ACCOUNT_ARCHIVED`。
- 当前字段为沟通方式、主题、内容、联系人、发生时间、创建时间和创建人。

后续如需纠错，应追加一条“更正记录”并引用原记录 ID，不能覆盖原文。该规则保证客户跟进历史、审计和交付验收一致。

## 9. 后续 TODO

1. 客户提供原始文件样例、目录、格式、编码、增量标识和字段映射后，实现正文解析、分块和向量索引。
2. 客户确认联网站点白名单、访问授权、抓取频率、留存与引用规则后，在工具层接入联网检索；不能仅依赖提示词限制访问。
3. 客户提供 OEM 类目、评分权重、Top 排序口径和报告模板后完成正式匹配输出。
4. 模型余额恢复后，验证两个系统 Agent 和各页面路由的完整模型内容、冲突保留和中文输出质量。
5. 后续新增页面智能能力时，必须先补充显式 `feature -> Agent` 路由、输入上下文和验收样例，避免落入通用 Agent 后无法验收。
