# STA-100 私有智能体设备系统用户需求规格说明书（Go 重构版）

## 1. 文档范围

本文基于以下资料整理：

- 客户与 OpenClaw 沟通记录：`/home/User/gsx/test/std/webchat-ALL-history-2026-08-04.md`
- 现有项目代码：`/home/User/gsx/test/std/projects/cycling-agent`
- 澄清答复：`/home/User/gsx/gitdir/std/需求澄清问题.md`
- 真假数据对照：`/home/User/gsx/gitdir/std/现有功能真假对照表.md`

本文目标：

- 明确 STA-100 系统的业务目标、页面范围、接口范围、数据范围、非功能要求和验收标准。
- 梳理现有每个页面、每组接口的功能、作用、实现方式。
- 将现有 Python/Flask 系统完整替换为 Go 单体可执行程序。
- 保留现有前端视觉结构、交互方式、路由路径、接口路径和主要响应字段；如旧接口存在冲突、Mock 或错误，应在 Go 版中修正并保持前端兼容。

## 2. 核心业务目标

STA-100 的业务目标不是进入欧洲自行车市场本身，而是进入欧洲 AI 市场。

产品定位：

- `STA-100` 是私有智能体设备，面向欧洲公司客户。
- 系统是生产力工具，不是单一骑行行业工具。
- `cycling-agent` 是 STA-100 的第一个行业样板，用于证明“行业知识库 + 智能体 + 私有部署”的能力。
- 后续应扩展到更多行业和公司客户，例如法律、医疗、制造、金融、咨询、教育、政府、跨境贸易等。

客户确认要求：

- 全部页面保留。
- 全部接口保留。
- 全部现有功能都要实现。
- 之前未实现、Mock、演示、占位的功能，在 Go 版中需要补成真实实现。
- 使用 Go 完整替换现有 Python Web 服务和业务逻辑。
- 沿用现有 `SQLite / JSON / Markdown / 静态资源 / 模板 / 数据目录`。
- 支持 ARM Linux，产物为单个可执行文件。
- 端口和访问方式保持当前系统可运行：`0.0.0.0:7777`。

## 3. 总体架构要求

### 3.1 架构形态

Go 版本采用单体服务架构：

- 一个 Go 可执行文件负责 HTTP 服务、页面路由、API、静态资源、PWA、调度任务、数据读写、日志、备份、权限和健康检查。
- 不允许长期依赖 Python 业务逻辑。
- 可以为了数据兼容读取原 Python 项目生成的数据，但业务执行逻辑应由 Go 重写。
- 原有脚本和历史数据可以作为测试数据、迁移来源和兼容检查依据。

### 3.2 推荐模块划分

```text
cmd/sta100/
  main.go                         # 程序入口
internal/http/
  router.go                       # 路由注册
  middleware.go                   # 日志、鉴权、错误处理、CORS
internal/pages/
  templates.go                    # 页面渲染
  static.go                       # 静态资源/PWA
internal/console/
  scan.go                         # 文件系统扫描
  search.go                       # 全局搜索
  detail.go                       # KPI 明细
internal/customer/
  search.go                       # SQL 客户搜索
  rag.go                          # 客户 RAG 搜索
  lead.go                         # 线索评分/开发
  local.go                        # 本地客户发现
  crm.go                          # CRM 状态
internal/agent/
  chat.go                         # 智能体对话
  route.go                        # Master 路由
  registry.go                     # 智能体列表
internal/rag/
  index.go                        # 文档索引
  search.go                       # RAG 检索
internal/tools/
  quotation.go                    # 报价
  shipping.go                     # 物流
  compliance.go                   # 合规
  oem.go                          # OEM 匹配
internal/docs/
  invoice.go                      # 单据/发票
  pdf.go                          # PDF 生成
  document.go                     # 文档处理
internal/automation/
  scheduler.go                    # Go 调度器
  daily.go                        # 日报/周报
  push.go                         # 推送订阅/出箱
  ingest.go                       # RSS/B站/价格快照任务
internal/store/
  sqlite.go                       # SQLite 连接
  jsonstore.go                    # JSON 文件读写
  markdown.go                     # Markdown 文件读取
internal/auth/
  user.go                         # 临时登录/用户体系
  session.go                      # session_id/user_id 兼容
internal/system/
  health.go                       # 健康检查
  backup.go                       # 备份/恢复
  hardware.go                     # 本机硬件状态
```

### 3.3 智能体底层模型确认项

当前 Python 项目中，Web 服务初始化为 `STA100LLM(mock_mode=True)`，因此现有 `/api/chat` 和 `/api/agent/call` 的智能体对话链路默认是演示模式，不是真实模型推理。

现有代码和前端文案中体现的目标模型为：

| 项 | 当前设计 | Go 版要求 |
|---|---|---|
| 模型名称 | `DeepSeek V4-Flash 7B INT4` | 作为首期默认模型名称写入配置 |
| 模型格式 | `GGUF` | 保持本地量化模型路线 |
| 默认路径 | `models/deepseek-v4-flash-7b-int4/model.gguf` | Go 版继续兼容该路径 |
| Python 推理方式 | `llama_cpp.Llama` | Go 版需提供真实 LLM Adapter |
| 当前运行状态 | `mock_mode=True` | Go 版生产默认必须为真实推理，不允许默认 mock |

Go 版需要确认并实现：

- 提供统一 `LLMAdapter`，支持本地 `llama.cpp/GGUF` 或本机 OpenAI-compatible 服务。
- 生产环境 `STA100_MOCK_ENABLED=false`，模型不可用时返回明确错误，不静默伪造回答。
- `/api/health` 返回 `model_name`、`model_loaded`、`model_path`、`rag_available`、`mock_mode`。
- `/api/chat`、`/api/agent/call` 返回 `response`、`sources`、`confidence`、`latency_ms`、`agent`、`agent_name`，并保留模型状态字段。
- RAG 无命中时可以回答通用建议，但必须标注“未命中本地知识库”。
- 模型升级包放在 `models/`，不写入 SQLite；支持配置切换和 OTA 校验。

详细确认项见：`/home/User/gsx/gitdir/std/需求确认项补充-智能体模型.md`。

### 3.4 数据目录沿用

Go 版必须沿用以下目录：

| 目录 | 用途 | Go 实现要求 |
|---|---|---|
| `data/` | SQLite、JSON、推荐、订阅、日志、brief、CRM | 保持原结构，可新增 migration 版本表 |
| `_customers/` | 客户 JSON 源 | 继续可预览、可导入 |
| `_donald_chunks/` | Donald RAG 文档块 | 继续可搜索 |
| `_raw_*` | 原始采集数据源 | 继续扫描统计 |
| `web_static/` | CSS/JS/图标/manifest | 原样提供静态服务 |
| `web_templates/` | HTML 页面 | 保持页面结构，允许小幅样式微调 |
| `logs/` | Web/Cron/任务日志 | Go 写入结构化日志和普通日志 |
| `output/` | PDF、备份、邮件等生成文件 | Go 写入 |

### 3.5 当前数据缺失风险

现有项目代码依赖但当前拷贝中未看到以下关键文件：

- `data/customers.db`
- `data/cycling_knowledge.db`
- `data/oem_factory.db`
- 部分 `_raw_*`、`_customers`、`_donald_chunks` 实际数据

Go 版必须处理：

- 文件存在：正常读写。
- 文件不存在：自动初始化空库或测试库，不允许页面崩溃。
- 历史数据导入：提供 `sta100 migrate` 或启动时自动迁移能力。
- 测试数据清理：支持后台脚本/接口清理测试数据。

## 4. 页面需求

### 4.1 `/` 首页

| 项 | 说明 |
|---|---|
| 当前行为 | 跳转 `/console` |
| Go 版要求 | 保持 302 跳转到 `/console` |
| 作用 | 用户零操作进入控制台 |
| 实现 | Go HTTP handler 返回 redirect |

### 4.2 `/console` 控制台

对应文件：`web_templates/console.html`，静态资源 `console.css`、`console_api.js`、`console.js`。

页面模块：

| 模块 | 页面功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|---|
| 仪表盘 | 展示智能体、知识库、OEM、客户、新闻、对话、线索、RAG KPI | `/api/console/stats` | 系统总览 | 扫描本地目录、SQLite、JSON、日志并聚合 |
| 智能体 | 展示智能体目录，点击打开聊天面板 | `/api/console/sections/agents`、`/api/chat` | 让用户直接调用业务智能体 | 从 Go 智能体注册表读取 |
| 知识库 | 展示 12 类知识库、RAG 检索 | `/api/console/sections/knowledge`、`/api/console/kb-detail`、`/api/console/rag-search` | 管理和检索行业知识 | 扫描 Markdown、RAG chunks、SQLite |
| 数据源 | 展示 `_raw_*` 数据源、Cron 任务 | `/api/console/sections/datasources`、`/api/console/cron-status` | 监控数据采集状态 | 扫描目录和日志 |
| 工具 | 展示 API 工具目录 | `/api/console/sections/tools`、`/api/console/api-docs` | 让用户知道系统有哪些能力 | 从 Go 路由注册表生成 |
| 设置 | 语言、安装、刷新、快捷键 | `/api/console/scan/refresh`、`/api/console/shortcut/<platform>` | 系统维护入口 | Go 实现缓存刷新与快捷方式生成 |
| Ctrl+K 搜索 | 搜索智能体、工具、文档、客户、Cron | `/api/console/search` | 一个搜索框找到全部入口 | Go 搜索聚合器 |

保留要求：

- 左侧栏、顶部栏、KPI 卡片、搜索触发、快捷操作、聊天抽屉保持现有结构。
- 可以小幅优化样式，但不能破坏页面层级和交互。
- 原有 Mock KPI 必须换成真实扫描/真实数据库结果。

### 4.3 `/search` 全局搜索页

对应文件：`web_templates/search.html`。

| 功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|
| 输入关键词实时搜索 | `/api/console/search?q=` | 搜索智能体、工具、文件、客户、RAG、Cron、配置 | 建立统一搜索服务，按类型返回结果 |
| 快捷 chip | 同上 | 快速搜索 Trek、Shimano、suppliers、RAG、cron 等 | 前端保持，后端按关键词匹配 |
| Enter 跳转第一结果 | 搜索结果里的 `action` | 提高操作效率 | Go 返回兼容 action 或更安全的 `url/action_type` |

Go 实现建议：

- 兼容旧字段：`kind`、`icon`、`title`、`sub`、`action`。
- 同时新增安全字段：`url`、`action_type`、`params`，逐步替代前端 `eval(action)` 风险。

### 4.4 `/dashboard` STA-100 行业智能体看板

对应文件：`web_templates/dashboard.html`，加载：

- `dashboard.js`
- `dashboard_i18n.js`
- `dashboard_recs.js`
- `dashboard_oem.js`
- `dashboard_leads.js`
- `dashboard_local.js`
- `dashboard_search.js`
- `dashboard_tools.js`
- `dashboard_subscribe.js`

页面 Section：

| Section | 功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|---|
| Overview | KPI、推荐、快捷智能体、日程 | `/api/recommendations/<user_id>`、`/api/push/subscribe` | 用户进入后看到今日重点 | 推荐、订阅、KPI 从真实数据生成 |
| Agents | 24 个智能体列表 | `/api/agent/call`、`/api/chat` | 快速进入智能体能力 | Go 智能体注册和调度 |
| Documents | PI/CI/PL/CO 单据表和生成按钮 | `/api/docs`、`/api/docs/<filename>` | 外贸单据生成与下载 | Go PDF/文档生成 |
| Orders | 订单表 | `/api/orders` | 查看订单 | 旧版为 Mock，Go 版必须接真实 SQLite/JSON 订单表 |
| Customers | 客户表 | `/api/customers`、`/api/customer/search` | 查客户和风险 | 旧版冲突，Go 版统一真实客户接口 |
| Private DB | 私有数据库/文件上传/备份 | `/api/document/process`、`/api/backup/create` | 本地私有知识库 | Go 文件处理和备份 |
| Daily Brief | 日报/周报/推荐 | `/api/brief/daily`、`/api/brief/weekly`、`/api/brief/latest/<user_id>` | 自动行业情报推送 | Go scheduler + brief 生成 |
| Settings | 语言、模型、RAG、Cron | `/api/health`、`/api/scheduler/stats` | 系统设置 | Go 状态和配置 |

Dashboard 额外挂载 Widget：

| Widget | JS | 调用接口 | 功能 | Go 版要求 |
|---|---|---|---|---|
| 个性化推荐 | `dashboard_recs.js` | `/api/recommendations/<user_id>` | 根据用户行为推内容 | 使用真实日志和推荐缓存 |
| OEM 工厂匹配 | `dashboard_oem.js` | `/api/oem/match`、`/api/oem/stats` | 查询 Top OEM 工厂，支持 Markdown 导出 | 用真实工厂索引和评分逻辑 |
| 客户开发 | `dashboard_leads.js` | `/api/lead/top`、`/api/lead/outreach`、`/api/lead/crm` | Top 线索、邮件模板、CRM 状态 | 用真实客户库和 CRM JSON |
| 本地客户发现 | `dashboard_local.js` | `/api/local/countries`、`/api/local/find`、`/api/local/cluster`、`/api/local/hot` | 按国家/类型找本地客户 | 用真实客户库聚类 |
| 客户统一搜索 | `dashboard_search.js` | `/api/customer/search`、`/api/customer/rag-search` | 多维客户搜索、CSV 导出 | SQL + RAG 双模式 |
| 海外客户工具集 | `dashboard_tools.js` | `/api/tools/oem-match-v2`、`/api/tools/shipping`、`/api/tools/compliance` | OEM、物流、合规三工具 | Go 规则和计算引擎 |
| 自动订阅 | `dashboard_subscribe.js` | `/api/push/subscribe`、`/api/push/unsubscribe` | 首次访问自动订阅本地 06:00 日报 | Go 订阅表和 scheduler |

### 4.5 `/all` 全量看板

对应文件：`dashboard_all.html`。

| 功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|
| 一页展示全量 KPI | `/api/console/stats` | 老板查看完整系统状态 | 调用同一 stats 聚合器 |
| Cron 状态 | `/api/console/cron-status` | 查看定时任务是否运行 | 读取 Go scheduler 状态和旧日志 |
| KPI 二层详情 | `/api/console/detail/<kind>` | 从 KPI 进入真实数据 | 返回真实 items |
| MD 查看 | `/api/console/md-view` | 预览文档 | 安全白名单读取 |

### 4.6 `/sta100` 产品页

对应文件：`sta100.html`。

| 功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|
| 展示 STA-100 产品定位、定价、Agent 能力 | `/api/sta100/agents` | 市场展示和客户理解产品 | 从 Go agent registry 返回 |
| 询价/联系表单 | `/api/sta100/quote` | 收集销售线索 | 写入 `data/sta100_leads/` 或 CRM |

要求：

- 保持现有视觉和文案体系，中文优先。
- 产品页内容可以产品化，但表单必须真实落库。

### 4.7 `/sta100/app` 客户使用界面

对应文件：`sta100_app.html`。

| 功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|
| 推送卡片展示 | `/api/openclaw/card/<kind>` | 客户少操作，看关键业务卡片 | 从真实订单/物流/报价/日报生成 |
| 聊天入口 | `/api/chat` | 客户问 STA-100 | Go 智能体对话 |
| 一次性看板体验 | 多个 `/api/openclaw/*` | 用于 WhatsApp/微信/邮件链接 | Go 卡片渲染和鉴权 |

旧版卡片数值多为示例，Go 版必须改为真实业务数据或明确返回空状态。

### 4.8 `/console/install`

对应文件：`console_install.html`。

| 功能 | 调用接口 | 作用 | Go 实现 |
|---|---|---|---|
| 桌面/手机/平板安装引导 | `/api/console/shortcut/<platform>` | 支持 PWA 和快捷方式 | 生成 Windows `.url`、macOS `.webloc`、Linux `.desktop` |

## 5. 接口需求清单

### 5.1 页面和静态资源接口

| 路由 | 方法 | 功能 | 作用 | Go 实现要求 |
|---|---|---|---|---|
| `/` | GET | 跳转控制台 | 默认入口 | 302 到 `/console` |
| `/old` | GET | 老版内联页面 | 兼容旧入口 | 可保留，建议标注 deprecated |
| `/console`、`/console/` | GET | 控制台 | 主系统后台 | 渲染 `console.html` |
| `/console/<section>`、`/console/<section>/` | GET | hash 跳转兼容 | 支持 `/console/dashboard` | 合法 section 转 `/console#section` |
| `/search` | GET | 全局搜索页 | 单搜索框入口 | 渲染 `search.html` |
| `/dashboard`、`/dashboard/<path:subpath>` | GET | STA-100 看板 | 客户/用户操作页 | 返回 `dashboard.html` |
| `/all`、`/all/` | GET | 全量看板 | 老板视图 | 渲染 `dashboard_all.html` |
| `/console/install`、`/console/install/` | GET | 安装页 | PWA/快捷方式 | 渲染 `console_install.html` |
| `/sta100`、`/sta100/` | GET | 产品页 | 销售展示 | 渲染 `sta100.html` |
| `/sta100/app`、`/sta100/app/` | GET | 客户使用端 | 推送卡片+聊天 | 渲染 `sta100_app.html` |
| `/_customers/<path:filename>` | GET | 客户 JSON 预览 | 控制台查看客户源 | 安全读取 `_customers/*.json` 前若干 KB |
| `/static/*` | GET | 静态资源 | CSS/JS/icon/manifest | Go 嵌入或文件系统提供 |

### 5.2 Console 管理接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/console/stats` | GET | 聚合控制台 KPI | 仪表盘主数据 | 扫描文件夹、SQLite、JSON、日志，带缓存 |
| `/api/console/scan/refresh` | GET | 清缓存重扫 | 手动刷新系统状态 | 清空 stats cache |
| `/api/openclaw/snapshot` | GET | 一次返回完整快照 | 给 OpenClaw/外部系统读取 | 返回 health、stats、KPI 摘要、链接 |
| `/api/console/detail/<kind>` | GET | KPI 二层详情 | 点击 KPI 看真实列表 | kind 支持 agents/docs/factories/customers/news/convos/leads/rag/llm/cron/dbs/raw |
| `/api/console/sections/<section>` | GET | 模块详情 | console 子页面加载 | section 支持 agents/knowledge/datasources/tools |
| `/api/console/search` | GET | 全局搜索 | 搜索智能体、工具、文档、客户、Cron | 多源聚合搜索，返回兼容字段 |
| `/api/console/source` | GET | 源码查看 | 管理员调试 | 白名单文件读取；Go 版可改为模块说明和源码预览 |
| `/api/console/cron-log` | GET | Cron 日志 | 查看任务运行结果 | 读取 `logs/cron_*.log` 和 Go scheduler 日志 |
| `/api/console/dir-files` | GET | 目录文件列表 | 查看 `_raw_*`、数据目录 | 防路径穿越，最多返回指定数量 |
| `/api/console/md-view` | GET | Markdown 预览 | 查看知识库文档 | 白名单读取 Markdown/JSONL 文本 |
| `/api/console/kb-detail` | GET | 知识库分类详情 | 查看某类文件和 RAG 命中 | 扫描分类目录并可做 RAG 搜索 |
| `/api/console/customer-dist` | GET | 客户分布 | 图表数据 | 统计 `_customers` 或 `customers.db` |
| `/api/console/api-docs` | GET | API 文档 | 前端工具目录 | 从 Go 路由元数据生成 |
| `/api/console/news` | GET | 最新新闻 | 控制台新闻区 | 读取采集结果，缺失时返回空数组 |
| `/api/console/cron-status` | GET | 任务状态 | 数据源页 | 汇总 Go scheduler 和日志 |
| `/api/console/rag-search` | GET | RAG 检索测试 | 知识库搜索 | 使用 Go RAG 索引，失败不得返回 demo，返回明确错误 |
| `/api/bundle/<filename>` | GET | 下载打包文件 | 代码/资料下载 | 读取允许包名；Go 版应支持生成新包 |
| `/api/console/shortcut/<platform>` | GET | 快捷方式文件 | 安装到桌面 | 生成 windows/macos/linux 文件 |
| `/api/openclaw/preview/<kind>` | GET | 富卡片 HTML 预览 | 给管理员/销售预览推送效果 | 复用卡片 JSON，再由 Go 模板渲染 HTML |
| `/api/openclaw/push/demo` | GET | 今日推送演示清单 | 展示今日会推哪些卡片 | 调用卡片生成器并汇总渠道、时间和数量 |

### 5.3 智能体和对话接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/chat` | POST | 智能对话 | 用户与 STA-100/子智能体对话 | Go 实现真实 agent router，不再 `mock_mode=True` |
| `/api/agent/call` | POST | 指定智能体调用 | Dashboard 点击智能体 | 归一化 agent ID，调用对应 Go agent |
| `/api/master/route` | POST | Master 意图分发 | 自动路由任务 | 基于规则+RAG/LLM 判断意图 |
| `/api/rag/query` | POST | RAG 客服问答 | 产品/维修/兼容问答 | 使用知识库索引返回答案和来源 |
| `/api/session/new` | POST | 新建会话 | 多轮对话 | 生成 session_id，保存在内存/SQLite |
| `/api/session/<sid>/reset` | POST | 重置会话 | 清上下文 | 清理指定 session 上下文 |
| `/api/health` | GET | 健康检查 | 前端和部署检测 | 返回版本、模型、RAG、DB、scheduler 状态 |

### 5.4 客户、线索、CRM 接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/customer/search` | GET | 多维客户搜索 | 按国家/名称/邮箱/电话/业务/type/source 搜客户 | 查询 `data/customers.db`，支持分页、facets、has_contact |
| `/api/customer/rag-search` | GET/POST | 客户语义搜索 | 语义找客户 | Go TF-IDF/BM25/向量检索，支持 country/type 过滤 |
| `/api/customer/rag-stats` | GET | RAG 统计 | 查看客户索引状态 | 返回客户数、词表/索引大小、构建时间 |
| `/api/lead/<int:lead_id>` | GET | 线索详情 | 查看评分明细 | 返回客户、score、reasons、activity |
| `/api/local/countries` | GET | 国家列表 | 本地客户发现下拉 | 统计每国客户数和联系比例 |
| `/api/local/find` | GET | 本地客户查找 | 按国家/type 找 Top 客户 | SQL + lead scoring |
| `/api/local/cluster` | GET | 客户类型聚类 | 看某国客户结构 | group by customer_type |
| `/api/local/hot` | GET | 活跃客户 | 最近被搜索/打开客户 | 基于 user logs |
| `/api/lead/top` | GET | Top 线索 | 客户开发列表 | 客户评分：国家、类型、兴趣、联系完整度、活跃度 |
| `/api/lead/<id>` | GET | 线索详情 | 查看评分明细 | 返回客户、score、reasons、activity |
| `/api/lead/outreach` | POST | 生成开发模板 | Email/WhatsApp/微信文案 | 基于客户画像、产品、语言模板生成 |
| `/api/lead/crm` | POST | 更新 CRM 状态 | 跟进新/已联系/谈判/成交/流失 | 写 `data/sta100_crm.json` 或 SQLite |
| `/api/lead/stats` | GET | 线索统计 | 看 CRM 和客户开发状态 | 统计客户数、活跃数、状态分布 |
| `/api/customers` | GET | 客户列表 | Dashboard 客户表 | 旧版重复/Mock，Go 版统一为真实客户分页接口 |
| `/api/orders` | GET | 订单列表 | Dashboard Orders 表 | 读取真实订单表/JSON，支持分页与筛选 |
| `/api/customer` | POST | 客户风险评估 | 外贸风险判断 | 写/读客户库，输出风险等级、建议账期 |

必须修复：

- 现有 Python 中 `/api/customers` 出现重复定义且一个是 Mock。Go 版只能保留一个真实实现。
- 如果前端仍需要旧 mock 字段，应从真实客户库映射字段，不能写死。

### 5.5 订阅、推荐、日报和调度接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/push/subscribe` | POST | 自动订阅日报 | 用户首次打开 dashboard 自动订阅 | 按 `X-User-Id`、tz、lang、country 写订阅表 |
| `/api/push/unsubscribe` | POST | 取消订阅 | 用户关闭日报 | 更新订阅状态 |
| `/api/push/status` | GET | 订阅状态 | 查看当前用户和总订阅数 | 读取订阅 JSON/SQLite |
| `/api/recommendations/<user_id>` | GET | 个性化推荐 | Dashboard 顶部推荐卡 | 基于用户日志、RAG、brief 缓存生成 |
| `/api/brief/daily` | POST | 生成日报 | 行业日报推送 | 从新闻、价格、合规、用户兴趣生成 |
| `/api/brief/weekly` | POST | 生成周报 | 周度汇总 | 聚合一周新闻和用户行为 |
| `/api/brief/latest/<user_id>` | GET | 最新 brief | Dashboard 显示最近日报 | 读取用户 brief outbox/cache |
| `/api/daily/today` | GET | 今日日报 | 旧 dashboard 日报 | Go 版接真实日报，不再模板 mock |
| `/api/daily/date/<date>` | GET | 指定日期日报 | 历史日报查询 | 读取历史 brief |
| `/api/scheduler/tick` | POST | 手动触发调度 | 测试 daily/weekly/recs | Go scheduler 执行一次 |
| `/api/scheduler/stats` | GET | 调度统计 | 查看任务状态 | 返回任务列表、上次运行、下次运行、错误 |
| `/api/skills` | GET | 永久技能列表 | 展示 STA-100 技能包 | 从 Go registry 返回 |

Go scheduler 必须实现：

- 每日当地 06:00 推日报。
- 周报任务。
- 推荐任务。
- RSS 采集。
- B站采集。
- 翻译处理。
- 价格快照。
- 支持与旧 cron/log 兼容。

### 5.6 工具、报价、物流、合规接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/tools/oem-match-v2` | GET/POST | 自然语言 OEM 匹配 | 海外客户工具集 | 解析 query、国家、类型，返回 Top 工厂 |
| `/api/tools/shipping` | GET/POST | 跨境物流计算 | 计算海运/空运/快递/落地价 | Go 规则表和国家港口配置 |
| `/api/tools/compliance` | GET/POST | 合规检查 | CE/FDA/CBAM/UN38.3 等清单 | Go 规则引擎 |
| `/api/oem/match` | POST | OEM 工厂匹配 | Dashboard OEM Widget | 真实工厂索引、评分、日志记录 |
| `/api/oem/stats` | GET | OEM 工厂统计 | 显示工厂库规模 | 返回分类和构建时间 |
| `/api/quotation` | POST | 报价计算 | 快速计算 FOB/保费等 | Go 报价模块 |
| `/api/factory-price` | POST | FOB 转出厂价 | 退税测算 | Go Incoterm/退税计算 |
| `/api/section301` | POST | 美国 301 关税 | 美国出口风险 | Go 关税规则 |
| `/api/cbam-info` | GET | CBAM 信息 | 欧盟碳关税提示 | Go 返回规则表，可配置 |
| `/api/incoterms` | GET | Incoterms 11 术语 | 外贸术语查询 | Go 常量/配置 |
| `/api/inquiry` | POST | 直接询价 | 快捷入口 | 真实报价流程 |
| `/api/shipping` | POST | 直接查物流 | 快捷入口 | 真实物流计算/ETA |
| `/api/payment` | POST | 收款方式推荐 | LC/TT/OA 决策 | Go 支付规则 |

### 5.7 单据、发票、文档接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/docs` | POST | 生成 PI/CI/PL/CO + PDF | 外贸单据 | Go PDF 生成，写 `output/` |
| `/api/docs/<path:filename>` | GET | 下载/预览 PDF | 用户下载单据 | 安全读取 output |
| `/api/invoice/create` | POST | 创建发票 | 销售/采购/PI/CN/DN/服务等 | SQLite 保存发票和明细 |
| `/api/invoice/ocr/upload` | POST | OCR 上传 | 识别发票图片/文件 | Go OCR 插件接口；无 OCR 时返回明确错误 |
| `/api/invoice/ocr/create` | POST | 从 OCR 生成发票 | 自动建单 | 根据 upload_id 生成 invoice |
| `/api/invoice/sync/crm` | POST | 同步 CRM | 发票进入客户系统 | 先做本地 CRM，外部系统预留 adapter |
| `/api/invoice/sync/erp` | POST | 同步 ERP | 采购/财务系统 | 本地记录同步状态 |
| `/api/invoice/<invoice_id>` | GET | 查询发票 | 查看详情 | SQLite 查询 |
| `/api/invoice/list` | GET | 发票列表 | 财务列表 | 支持 status/type/limit |
| `/api/document/process` | POST | 文档处理 | 上传 PDF/Word/Excel 解析摘要 | Go 文件读取、OCR/解析、关键词 |
| `/api/document/upload` | POST | 同上 | 上传兼容路径 | 与 process 共用 |
| `/api/email/scenarios` | GET | 邮件场景 | 邮件生成器配置 | 返回模板场景 |
| `/api/email/generate` | POST | 生成邮件 | 开发信/跟进信 | Go 模板引擎，写 output/emails |
| `/api/email/<filename>` | GET | 下载邮件 | 下载 txt/html | 安全读取 output/emails |

### 5.8 行业、设计、维修、内容接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/design/color-scheme` | POST | 配色推荐 | 设计顾问 | 规则+RAG |
| `/api/design/geometry` | POST | 几何推荐 | 车型设计 | 本地规则 |
| `/api/design/ergonomics` | POST | 人体工程学 | 车架/尺码 | 公式+规则 |
| `/api/design/components` | POST | 配件推荐 | BOM 搭配 | 价格层级规则 |
| `/api/design/inspiration` | GET | 设计灵感 | 品牌案例 | 本地知识库 |
| `/api/design/spec-sheet` | POST | 规格表生成 | 完整产品定义 | 组合上述模块 |
| `/api/design/scenarios` | GET | 设计案例 | 演示/训练 | Markdown/JSON 配置化 |
| `/api/design/compare` | GET | 品牌对标 | 市场定位 | 本地比较表 |
| `/api/design/rag-search` | POST | 设计 RAG | 搜设计知识 | Go RAG |
| `/api/market/industries` | GET | 子行业列表 | 市场分析入口 | 本地行业配置 |
| `/api/market/overview` | POST | 行业报告 | 市场决策 | RAG/Markdown 生成 |
| `/api/market/competitors` | POST | 竞品搜索 | 品牌分析 | 本地品牌库 |
| `/api/market/company` | POST | 公司画像/对比 | 客户/竞品分析 | 本地数据+RAG |
| `/api/market/timeline` | GET | 行业事件时间线 | 趋势分析 | 本地事件库 |
| `/api/market/full-report` | GET | 完整报告 | 一次导出 | 聚合报告 |
| `/api/cn-brands` | GET | 中国品牌列表 | 出海推荐 | 本地品牌库 |
| `/api/cn-brands/navihood` | GET | NAVIHOOD 详情 | 重点品牌推荐 | 配置化品牌档案 |
| `/api/cn-brands/overseas-recommend` | POST | 推荐中国品牌 | 对海外客户销售 | 规则+画像 |
| `/api/cn-brands/compare` | GET | 中外品牌对比 | 销售材料 | 本地对比 |
| `/api/repair/search` | POST | 维修 SOP 搜索 | 售后/维修知识 | RAG |
| `/api/repair/quick-diagnosis` | POST | 快速诊断 | 维修建议 | 规则树 |
| `/api/repair/pricing` | POST | 维修定价 | 服务报价 | 价格规则 |
| `/api/bom/list` | GET | BOM 模板列表 | 产品设计 | 读取 Markdown |
| `/api/bom/<bike_type>` | GET | BOM 详情 | 某车型配置 | 读取 Markdown |
| `/api/qc/checklist/<stage>` | GET | QC 检查清单 | 质检 | 读取 Markdown |
| `/api/sourcing/country/<country>` | GET | 国别选品 | 市场进入 | 读取 Markdown |
| `/api/content/topics` | GET | 内容选题 | 营销内容 | 读取 Markdown |
| `/api/content/platform/<platform>` | GET | 平台内容模板 | 小红书/抖音/YouTube/Instagram | 读取 Markdown |

### 5.9 库存、B2B、二手交易、尺码和兼容接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/clearance/plan` | POST | 清仓计划 | 库存清理 | Go 清仓策略 |
| `/api/b2b/order` | POST | B2B 下单 | 厂家-车店撮合 | 订单规划和价格计算 |
| `/api/used-bike/trade` | POST | 二手交易 | 二手估值/撮合 | 规则+估值 |
| `/api/measurement/match` | POST | 量体匹配 | 0 库存预售/尺码 | 尺码公式 |
| `/api/compatibility/recommend` | POST | 配件兼容推荐 | 搭配升级 | 兼容规则库 |
| `/api/inventory/add` | POST | 进货登记 | 入库 | SQLite 库存表 |
| `/api/inventory/scan-alerts` | POST | 库存预警 | 滞销/库存风险 | SQL 扫描规则 |
| `/api/inventory/damage` | POST | 损坏记录 | 售后/库存损耗 | 写库存损耗表和图片路径 |

### 5.10 系统、会议、备份接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/schedule/list` | GET | 日程列表 | 日程管理 | SQLite 持久化，不再仅内存 |
| `/api/schedule/create` | POST | 新建日程 | 创建任务 | 写 SQLite |
| `/api/schedule/update`、`/api/schedule/update/<sid>` | PUT | 更新日程 | 修改任务 | 更新 SQLite |
| `/api/schedule/delete`、`/api/schedule/delete/<sid>` | DELETE | 删除日程 | 删除任务 | 软删或硬删 |
| `/api/meeting/transcribe` | POST | 音频转文字 | 会议记录 | Go 接 ASR 插件；无插件时返回未配置 |
| `/api/meeting/summarize` | POST | 会议总结 | 生成纪要/待办 | 真实 LLM 或本地规则摘要 |
| `/api/rag/search` | POST | 通用 RAG 搜索 | 不经过聊天直接检索 | Go RAG 索引 |
| `/api/hardware/status` | GET | 硬件状态 | STA-100 设备监控 | 读取 CPU/内存/磁盘/温度 |
| `/api/backup/create` | POST | 创建备份 | 备份 data/配置 | Go zip/tar |
| `/api/backup/restore` | POST | 恢复备份 | 回滚数据 | 支持 dry_run，真实恢复需权限 |

### 5.11 OpenClaw 富卡片与外部预览接口

| 接口 | 方法 | 功能 | 作用 | Go 实现方式 |
|---|---|---|---|---|
| `/api/openclaw/card/<kind>` | GET | 富卡片数据 | 提供订单、物流、报价、业绩、日报卡片 JSON | 组合订单/物流/报价/日报/统计模块，缺数据时返回空态 |
| `/api/openclaw/preview/<kind>` | GET | 卡片 HTML 预览 | 给销售/管理员直接看推送效果 | 复用同一份卡片 JSON，再用 Go 模板渲染 HTML |
| `/api/openclaw/push/demo` | GET | 今日推送演示 | 展示用户今天会收到哪些卡片和渠道 | 调用卡片生成器，汇总渠道、时间、触达策略 |
| `/api/openclaw/snapshot` | GET | 系统快照 | 外部系统一次性获取健康、统计和摘要 | 聚合 health、stats、KPI、卡片和链接 |

### 5.12 路由兼容补充

以下路由在 Go 版中必须同时兼容带尾斜杠与不带尾斜杠的访问方式：

- `/console/<section>` 和 `/console/<section>/`
- `/dashboard` 和 `/dashboard/<path:subpath>`
- `/all` 和 `/all/`
- `/console/install` 和 `/console/install/`
- `/sta100` 和 `/sta100/`
- `/sta100/app` 和 `/sta100/app/`
- `/_customers/<path:filename>`
- `/api/docs/<path:filename>`
- `/api/lead/<int:lead_id>`

Go 路由实现要求：

- 旧前端或外链访问不能因为尾斜杠差异而 404。
- 需要统一 canonical URL，并在必要时做 301/302 归一化。
- 路径变量中只允许白名单字符和范围，不允许路径穿越。

## 6. 数据实现要求

### 6.1 SQLite

必须支持：

- `customers.db`
- `cycling_knowledge.db`
- `oem_factory.db`
- 新增 Go 系统库，例如 `sta100_system.db`

建议表：

| 表 | 用途 |
|---|---|
| `customers` | 客户数据 |
| `customer_fts` | 客户全文索引 |
| `leads_crm` | 线索状态 |
| `sessions` | 对话会话 |
| `messages` | 对话消息 |
| `orders` | 订单 |
| `invoices` | 发票 |
| `invoice_items` | 发票明细 |
| `schedules` | 日程 |
| `inventory` | 库存 |
| `jobs` | 调度任务状态 |
| `migrations` | 数据迁移记录 |

### 6.2 JSON

继续兼容：

- `data/sta100_subscribers.json`
- `data/sta100_crm.json`
- `data/sta100_leads.json`
- `data/sta100_recommendations/*.json`
- `_customers/*.json`
- 原有索引 JSON

Go 写 JSON 必须：

- 原子写入：临时文件 + rename。
- 保留 UTF-8。
- 保留现有字段。

### 6.3 Markdown / JSONL / 原始数据

用于：

- 知识库预览。
- RAG 索引。
- 行业报告。
- BOM、QC、内容模板。
- 用户行为日志。

Go 实现：

- 支持扫描、读取、搜索、摘要。
- 防止路径穿越。
- 大文件限制返回片段，避免接口卡死。

## 7. Go 调度任务需求

现有 `scripts/daily_update.sh` 包含：

- RSS 英文新闻。
- B站中文视频。
- 翻译队列状态。
- 价格快照。

Go 版需要实现自己的 scheduler，同时兼容旧脚本概念。

| 任务 | 频率 | 作用 | Go 实现 |
|---|---|---|---|
| RSS 采集 | 每日/可手动 | 更新行业新闻 | Go RSS parser，写 `cycling_knowledge.db` |
| B站采集 | 每日/可手动 | 更新中文视频 | Go HTTP client，关键词 URL encode |
| 翻译处理 | 定时 | 多语言内容覆盖 | Go 翻译队列，支持状态和批处理 |
| 价格快照 | 每日 | 价格历史 | Go price snapshot |
| 日报推送 | 用户当地 06:00 | 用户主动接收情报 | 根据订阅时区执行 |
| 周报推送 | 每周 | 周度汇总 | 聚合 brief |
| 推荐生成 | 每日/触发 | 个性化推荐 | 根据日志和 RAG |
| 备份 | 每日/手动 | 数据保护 | 打包 data/output/config |

## 8. 权限与用户需求

客户答复“暂时需要”，因此 Go 版需实现轻量权限体系，并兼容旧 `user_id/session_id`。

要求：

- 支持登录。
- 支持角色：`admin`、`operator`、`viewer`。
- 支持用户隔离，但首期可单租户。
- 兼容 `X-User-Id`、URL `?user_id=`、cookie/localStorage 中的旧 user_id。
- 未登录访问公开页：`/sta100`、可选 `/search`。
- 控制台、数据、备份、恢复、源码、日志等接口需鉴权。

## 9. 国际化要求

- 保留现有 17/18 语种能力。
- 默认中文优先。
- 保持现有词汇体系和产品命名。
- 前端继续支持 `dashboard_i18n.js` 当前机制，或 Go 后端直接提供 i18n JSON。
- 多语言字段不得破坏现有 `data-i18n` key。

## 10. Mock 清理要求

以下旧功能必须改成真实实现：

| 旧功能 | 当前问题 | Go 版要求 |
|---|---|---|
| `/api/customers` | Mock 固定 4 条，且路由重复 | 改为真实客户列表 |
| `/api/orders` | Mock 固定 3 条 | 改为真实订单表 |
| `/api/openclaw/card/<kind>` | 卡片数值示例 | 接真实订单/物流/报价/日报 |
| `/api/chat` | `mock_mode=True` | 接真实 Go agent / LLM adapter / RAG |
| `/api/meeting/transcribe` | 无音频时 mock | 接 ASR，未配置返回明确错误 |
| `/api/meeting/summarize` | 抽句 mock | 接 LLM/摘要模块 |
| `/api/daily/today` | 模板日报 | 接真实日报生成 |
| `/api/console/rag-search` 兜底 | 返回 demo | 返回错误或空结果，不假装命中 |
| Dashboard 静态订单/客户 DOM | 前端写死示例 | 改为接口加载真实数据 |

## 11. 部署与打包要求

目标：

- ARM Linux。
- 单个 Go 可执行文件。
- 默认监听 `0.0.0.0:7777`。
- 支持启动脚本打包。
- 支持静态资源和模板随二进制嵌入，或可选择外部目录覆盖。

必须提供：

| 文件/命令 | 作用 |
|---|---|
| `sta100` | Go 可执行文件 |
| `start.sh` | 启动服务 |
| `stop.sh` | 停止服务 |
| `build-arm.sh` | ARM 构建 |
| `pack.sh` | 打包单文件和资源 |
| `.env.example` | 环境变量模板 |
| `DEPLOY_GO.md` | Go 部署文档 |

建议环境变量：

```text
STA100_HOST=0.0.0.0
STA100_PORT=7777
STA100_DATA_DIR=./data
STA100_TEMPLATE_DIR=./web_templates
STA100_STATIC_DIR=./web_static
STA100_LOG_DIR=./logs
STA100_AUTH_ENABLED=true
STA100_DEFAULT_LANG=zh
STA100_SCHEDULER_ENABLED=true
```

## 12. 验收标准

### 12.1 页面验收

- `/console` 正常打开，6 大模块可切换。
- `/search` 正常搜索，结果可点击。
- `/dashboard` 正常打开，所有 widget 能加载。
- `/all` 正常显示系统总览和详情。
- `/sta100` 正常展示产品页，询价能落库。
- `/sta100/app` 正常显示卡片和聊天。
- `/console/install` 可生成快捷方式。
- PWA manifest、icon、service worker 不报错。

### 12.2 接口验收

- 本文列出的全部接口均有 Go 实现。
- 不允许存在未处理的 404/500。
- Mock 接口必须替换为真实实现或真实空状态。
- 旧前端调用参数和主要响应字段保持兼容。
- 有问题的旧接口允许修正，但前端必须同步兼容。

### 12.3 数据验收

- 可读取现有 SQLite / JSON / Markdown。
- 缺失数据文件时可初始化空库。
- 历史数据可查询。
- 支持测试数据清理。
- 支持备份和 dry-run 恢复。

### 12.4 调度验收

- Go scheduler 可启动、停止、查看状态。
- RSS/B站/翻译/价格快照任务保留。
- 日报/周报/推荐任务可手动触发。
- `/api/scheduler/tick` 和 `/api/scheduler/stats` 可用。
- 日志可在控制台查看。

### 12.5 部署验收

- ARM Linux 上可运行。
- 单个可执行文件可启动服务。
- 默认端口 `7777`。
- 页面和 API 均可访问。
- 提供打包脚本和启动脚本。

## 13. 开发优先顺序

虽然客户要求所有功能都要实现，但工程上建议按依赖顺序推进：

1. Go HTTP 服务、静态资源、页面路由。
2. 数据层：SQLite/JSON/Markdown/日志。
3. Console stats/search/detail。
4. 客户搜索、客户 RAG、线索开发、本地客户发现。
5. Dashboard widgets。
6. 智能体对话和 RAG。
7. 工具集：OEM、物流、合规、报价。
8. 单据、发票、文档处理。
9. 日报、订阅、推荐、scheduler。
10. 库存、B2B、二手、兼容、维修、市场、设计等行业模块。
11. 权限、备份、部署脚本、验收测试。

## 14. 最终结论

Go 版 STA-100 不是简单“换语言”，而是要把当前 Python/Flask 形成的页面、接口、数据、工具、智能体、调度和演示功能全部产品化：

- 页面保持一致。
- 接口保持兼容。
- 数据目录沿用。
- Mock 全部清理。
- 智能体和工具变为真实可用。
- scheduler 内置到 Go。
- ARM 单文件部署。

最终交付应是一套可在当前系统中直接运行的 STA-100 私有智能体设备后端和前端控制台。
