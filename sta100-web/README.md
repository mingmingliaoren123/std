# STA-100 Go Web 应用

本目录包含可运行的 STA-100 页面和 Go 服务。页面业务演示数据仍为原型数据；OpenClaw 网关状态、模型配置和 24 个 Agent 管理通过独立的 `/home/User/gsx/gitdir/std/openclaw-orchestrator` 接入真实本机部署。

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

写接口要求请求头 `X-STA100-Request: 1`。健康检查为 `GET /api/health`，当前返回 `prototype: false`，表示 Go 服务和 OpenClaw 编排接口已真实运行；客户、订单、单据等业务数据仍需后续接入 SQLite。

## 本机登录账户

- 首次运行默认用户名和密码均为 `admin`。
- 登录成功后使用 12 小时 HttpOnly、SameSite=Strict 会话 Cookie；未登录访问 `/api/v1/*` 业务接口返回 `401`。
- 设置 → 数据安全 → 本机登录账户可修改用户名和密码，必须验证当前密码；下次登录页只显示修改后的用户名并要求输入密码。
- 忘记密码时从登录页进入账户恢复，使用设备维护万能密码后重置当前用户名和密码。
- 账户保存在 `~/.config/sta100/auth.json`，文件权限为 `0600`；密码使用随机盐和 150,000 轮迭代 SHA-256 哈希，不保存明文。可通过 `STA100_AUTH_FILE` 指定其它持久化位置。

## 本次页面能力

- 客户、报价单、订单、单据和产品已补齐前端新增、查看、编辑、删除流程；数据保存在当前页面内存中，刷新服务或页面后恢复示例数据。
- 报价关联客户、订单关联客户/报价、单据关联订单支持输入名称或编号进行包含式模糊匹配，并要求从有效关系数据中选择后保存。
- 客户列表在订单数量、累计金额、最近更新字段旁提供升序/降序箭头；报价单和订单在金额字段旁提供同样排序。三个表格的表头全选、行复选框和当前筛选结果联动，排序后按记录 ID 保留选择。
- 报价和订单表单直接选择产品库中的启用产品，支持多条明细、数量、单价、折扣/小计和总额联动；报价转订单复制明细，订单生成单据再复制订单快照。旧示例记录通过摘要兼容恢复明细。
- 抽屉内发起编辑、新建或生成单据时会关闭抽屉再显示弹窗；弹窗内打开辅助抽屉时保留原弹窗，避免同层遮挡。
- 产品库同级新增供应商页面，支持公司、电话、联系人、邮件、产品、规格、报价、备注、来源等字段的前端维护、搜索、排序和导出；客户档案同步增加来源字段。
- OEM 和客户搜索统一展示“本地知识库”和“联网检索”；RAG 是联网检索的实现方式，不再作为第三个并列来源。
- 单据页支持关键字、单据类型、状态组合筛选、结果计数、重置和空状态。
- 模型设置展示 OpenClaw 真实默认模型、模型列表和凭据提供商状态；列表直接来自 `openclaw models list`。已保存 Key 只显示配置状态且默认隐藏，本次新输入的 Key 可显示/隐藏。
- 智能体管理展示 `main + 24` 的真实配置，并可按清单同步。
- 顶级模块和 24 个 Agent 使用原应用 Emoji，工具按钮继续使用本地 Lucide；概览三个工具标题不再重复显示旧快捷图片。
- STA-100 左上使用旧应用的骑自行车图标和 STA-100 名称。当前主题采用截图对应的白色内容面、浅灰工作区和深灰正文；品牌主色为粉色 `#ed3f93`，蓝色用于链接与选择，绿色用于在线和成功状态，琥珀色与红色分别用于待处理和错误风险。
- 新闻设置支持每次 1-100 条、1/2/3/6/8/12/24 小时频率和输入校验；超过 20 条在新闻列表内部滚动，概览仍固定展示 3 条。
- 本地客户发现国家下拉覆盖中国和欧洲国家的主要城市，每次仍限制单选国家、城市和客户类型。
- 智能体聊天来源改为复选项；联网来源可按 Agent 配置域名白名单。消息通过真实 OpenClaw Gateway 调用，超时或模型错误不会伪造答案。
- 桌面 1440px 和移动 390px 均支持，宽表在局部容器内滚动。
- 模板中心已实现图片与文件模板的本地选择、格式/20 MB 大小校验、待处理状态和接口缺口说明。模板上传、OCR、字段映射和发布后端尚未实现，页面不会把本地选择伪装成上传成功。
- 应用启动先显示本机账户登录页，侧栏账户按钮用于退出登录；设置页支持验证当前密码后修改账户，忘记密码支持维护密码验证后重置。

已保存 API Key 不支持明文读取。OpenClaw 的结构化接口只提供凭据是否配置、提供商和凭据类型；STA-100 不读取 OpenClaw 配置文件中的密钥材料，也不把 Key 放入接口响应、日志或 CLI 参数。

在线热升级仍不实现；设置页只保留管理员导入离线升级包、校验、安装和自动重启流程。

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
