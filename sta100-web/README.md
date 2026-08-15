# STA-100 Go Web 应用

本目录包含可运行的 STA-100 页面和 Go 服务。客户、报价、订单、单据、产品、供应商、私有文件元数据、设置、插件和任务已接入本机 SQLite；OpenClaw 网关状态、模型配置及 `2 个系统 Agent + 24 个业务 Agent` 通过独立的 `/home/User/gsx/gitdir/std/openclaw-orchestrator` 接入真实本机部署。

空数据库首次启动会写入一组带明确 `seeded_demo_data` 标记的业务演示记录，用于页面和接口验收。它们是持久化测试数据，不是客户正式业务事实。客户原始知识数据、OEM 工厂库、联网/RAG 来源、OCR、正式单据模板和离线升级包规范尚未提供，对应接口返回明确 `501 TODO_*`，不会伪造成功。

STA-100 只保留页面和业务适配层，通用 OpenClaw CLI 调用、模型校验、凭据写入、Agent 清单同步均在独立编排组件中。STA-100 专属清单位于 `config/sta100-agents.json`，不再放在 OpenClaw 官方运行目录。

## 运行

```bash
cd /home/User/gsx/gitdir/std/sta100-web
go run . -addr 0.0.0.0:8080
```

访问 `http://127.0.0.1:8080`。

ARM64 静态构建：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -trimpath -ldflags='-s -w' \
  -o dist/sta100-web-linux-arm64 .
```

## 已接 OpenClaw 接口

| 方法 | 接口 | 作用 |
|---|---|---|
| GET | `/api/v1/openclaw/status` | 返回版本、systemd、Gateway、RPC 和配置审计状态 |
| GET | `/api/v1/openclaw/models` | 返回默认模型、可用模型和不含密钥的凭据状态 |
| PUT | `/api/v1/openclaw/models/default` | 调用 `openclaw models set` 更新默认模型 |
| POST | `/api/v1/openclaw/models/auth` | 通过 stdin 写入 API Key，不回显密钥 |
| GET | `/api/v1/openclaw/agents` | 返回当前 Agent 的 ID、身份、Emoji 和继承模型 |
| POST | `/api/v1/openclaw/agents/sync` | 通过独立 Go 编排组件按 STA-100 清单同步 Agent |
| POST | `/api/v1/agents/chat` | 校验 Agent、来源和会话后，通过临时 `0600` 文件调用 `openclaw agent --message-file`，返回真实回复或明确错误 |
| POST | `/api/v1/assistant/query` | 统一执行本地检索、Knowledge Agent、页面领域 Agent 并行调用和 Coordinator 汇总；冲突证据并列保留，模型失败时返回可审计的部分结果 |
| GET/DELETE | `/api/v1/agent-token-usage` | 前期测试用 Token 可视化统计；读取或清空 OpenClaw 实际返回的输入、输出、缓存和总 Token，不保存提示词正文 |

写接口要求请求头 `X-STA100-Request: 1`。健康检查为 `GET /api/health`，当前返回 `prototype: false`。完整接口状态、数据性质、校验规则和联调结果见 [STA100-后端接口实现与验证说明.md](../STA100-后端接口实现与验证说明.md)。

## 本地数据

- SQLite 默认路径：`~/.local/share/sta100/sta100.db`，可用 `STA100_DB_PATH` 覆盖。
- 私有上传文件默认目录：`~/.local/share/sta100/private-files`，可用 `STA100_PRIVATE_DATA_DIR` 覆盖。
- 登录账户默认路径：`~/.config/sta100/auth.json`，可用 `STA100_AUTH_FILE` 覆盖。
- 前端业务数组登录后只接受 `/api/v1/bootstrap` 返回结果；接口失败时保持空状态，不回退到 JavaScript 假业务数据。

## 本机登录账户

- 首次运行默认用户名和密码均为 `admin`。
- 登录成功后使用 12 小时 HttpOnly、SameSite=Strict 会话 Cookie；未登录访问 `/api/v1/*` 业务接口返回 `401`。
- 设置 → 数据安全 → 本机登录账户可修改用户名和密码，必须验证当前密码；下次登录页只显示修改后的用户名并要求输入密码。
- 忘记密码时从登录页进入账户恢复，使用设备维护万能密码后重置当前用户名和密码。
- 账户保存在 `~/.config/sta100/auth.json`，文件权限为 `0600`；密码使用随机盐和 150,000 轮迭代 SHA-256 哈希，不保存明文。可通过 `STA100_AUTH_FILE` 指定其它持久化位置。

## 本次页面能力

- 客户、报价单、订单、单据、产品和供应商的新增、查看、编辑、归档/停用已接 Go REST API 和 SQLite，刷新页面或重启服务后数据保留。
- 客户详情和编辑流程支持追加沟通记录；记录包含方式、主题、内容、联系人、发生时间、创建人和创建时间。历史记录没有修改和删除接口，客户归档后仍可读但不能新增。
- 报价关联客户、订单关联客户/报价、单据关联订单支持输入名称或编号进行包含式模糊匹配，并要求从有效关系数据中选择后保存。
- 客户列表在订单数量、累计金额、最近更新字段旁提供升序/降序箭头；报价单和订单在金额字段旁提供同样排序。三个表格的表头全选、行复选框和当前筛选结果联动，排序后按记录 ID 保留选择。
- 报价和订单表单直接选择产品库中的启用产品，支持多条明细、数量、单价、折扣/小计和总额联动；报价转订单复制明细，订单生成单据再复制订单快照。旧示例记录通过摘要兼容恢复明细。
- 抽屉内发起编辑、新建或生成单据时会关闭抽屉再显示弹窗；弹窗内打开辅助抽屉时保留原弹窗，避免同层遮挡。
- 产品库同级新增供应商页面，支持公司、电话、联系人、邮件、产品、规格、报价、备注、来源等字段的持久化维护、搜索、排序和真实 XLSX 导出；客户档案同步增加来源字段。
- OEM、客户统一搜索和本地客户发现不再要求用户选择“本地/联网/RAG”来源。页面把输入交给统一查询接口，由 Go 先检索本地结构化数据，再由 Knowledge Agent 整理证据、业务 Agent 提供专业结果，最后由 Coordinator 汇总。
- 单据页支持关键字、单据类型、状态组合筛选、结果计数、重置和空状态。
- 模型设置展示 OpenClaw 真实默认模型和凭据状态；可选目录固定对应 OpenClaw `2026.7.1-2`，包含 20 个 API Key 提供商和 106 个文本模型，不在页面加载时调用远端目录。已保存 Key 只显示配置状态且默认隐藏，本次新输入的 Key 可显示/隐藏；升级 OpenClaw 时必须同步更新编排层目录版本。
- 智能体管理只展示 24 个业务 Agent；2 个系统 Agent 不在业务页面显示，但会随清单同步。OpenClaw 自带 `main` 也不属于页面业务 Agent。
- 顶级模块和 24 个 Agent 使用原应用 Emoji，工具按钮继续使用本地 Lucide；概览三个工具标题不再重复显示旧快捷图片。
- STA-100 左上使用旧应用的骑自行车图标和 STA-100 名称。当前主题采用截图对应的白色内容面、浅灰工作区和深灰正文；品牌主色为粉色 `#ed3f93`，蓝色用于链接与选择，绿色用于在线和成功状态，琥珀色与红色分别用于待处理和错误风险。
- 新闻设置支持每次 1-100 条、1/2/3/6/8/12/24 小时频率和输入校验；超过 20 条在新闻列表内部滚动，概览仍固定展示 3 条。
- 本地客户发现国家下拉覆盖中国和欧洲国家的主要城市，每次仍限制单选国家、城市和客户类型。
- 智能体聊天直接调用用户选择的业务 Agent，不再显示来源复选项。业务 Agent 的联网域名白名单保留在设置中由管理员配置；消息通过真实 OpenClaw Gateway 调用，超时、额度或模型错误不会伪造答案。
- 桌面 1440px 和移动 390px 均支持，宽表在局部容器内滚动。
- 模板中心已实现图片与文件模板的本地选择、格式/20 MB 大小校验、待处理状态和接口缺口说明。模板上传、OCR、字段映射和发布后端返回明确 `501`，页面不会把本地选择伪装成上传成功。
- 应用启动先显示本机账户登录页，侧栏账户按钮用于退出登录；设置页支持验证当前密码后修改账户，忘记密码支持维护密码验证后重置。
- 顶栏右侧显示登录用户；侧栏不再显示欧洲节点和在线服务。全局固定显示“内容由AI生成，请仔细甄别”，弹窗和抽屉遮罩不使用背景模糊，打开详情时原列表仍可辨认。
- 顶栏加入 `red_logo.png` 横向 STRATRONIX Logo；桌面端位于页面标题和搜索框之间，窄屏自动隐藏以保证操作区空间。
- 顶栏 Token 测试入口展示累计用量，详情提供最近一次请求、输入、输出、缓存、调用次数和按 Agent 汇总。数据来自 OpenClaw `result.meta.agentMeta.usage` 并写入 SQLite；未返回 usage 的失败调用单独计数，不按文本长度估算。该能力标记为临时测试功能，可独立清空，后续删除不会影响 Agent 会话和业务数据。

已保存 API Key 不支持明文读取。OpenClaw 的结构化接口只提供凭据是否配置、提供商和凭据类型；STA-100 不读取 OpenClaw 配置文件中的密钥材料，也不把 Key 放入接口响应、日志或 CLI 参数。

在线热升级仍不实现；设置页只保留管理员导入离线升级包、校验、安装和自动重启流程。

## 统一智能查询链路

```text
用户输入与页面上下文
  -> Go 检索 SQLite 与私有文件元数据
  -> sta100-knowledge 整理本地证据
  -> Go 根据 page/feature 并行调用 1-2 个业务 Agent
  -> sta100-coordinator 统一汇总
  -> 页面展示统一结果、证据、冲突、已用 Agent 和各阶段状态
```

页面不按“互联网/本地知识库”拆分结果区。每条证据仍保留来源、记录编号和更新时间；发现冲突时全部并列展示，不自动覆盖。Knowledge Agent 或 Coordinator 不可用时，接口以 `partial=true` 返回本地证据、可展示记录、成功的业务 Agent 输出及失败阶段，不把降级结果标记为完整模型结论。

客户私有文件正文解析、分块和向量索引等待正式原始数据格式，目前统一链路只读取 SQLite 结构化数据和私有文件元数据。联网检索也须在客户确认站点白名单、授权和留存规则后接入工具层。

## 部署顺序

```bash
cd /home/User/gsx/gitdir/std/openclaw
./scripts/deploy.sh

cd ../openclaw-orchestrator
./scripts/build.sh
./bin/openclaw-operator --bin ../openclaw/bin/openclaw \
  --manifest ../sta100-web/config/sta100-agents.json agents sync

cd ../sta100-web
go run . -addr 0.0.0.0:8080
```

如果运行目录与上述目录不同，请显式设置 `OPENCLAW_BIN`、`OPENCLAW_CONFIG_PATH` 和 `STA100_AGENT_MANIFEST`。STA-100 Go module 使用本地 `replace` 引用编排组件，交付时应将两个目录一起归档，或改成内部 Go module 版本。
