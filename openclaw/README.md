# STA-100 OpenClaw 部署说明

本目录保留 STA-100 使用的 OpenClaw 运行包、部署脚本和 systemd 单元；Node.js 使用当前盒子的本机环境，业务 Agent 清单由使用方项目维护。

## 固定版本与来源

- OpenClaw：`2026.7.1-2`，2026-08-10 查询时 npm `latest` 与 GitHub 最新稳定 Release 均为该版本。
- 官方 Release：`https://github.com/openclaw/openclaw/releases/tag/v2026.7.1-2`，发布时间 2026-08-04。
- Node.js：使用本机 `/usr/bin/node v22.23.1`，Linux ARM64，满足 OpenClaw 的 `>=22.22.3 <23` 要求。
- npm 包 SHA-256：`5bb525f36f471a41239615d321c441778c7e1c007018ed6d84b795be77803276`。
- npm integrity：`sha512-ycF3yPcbjN6bUPeaUx6Mh6vze1hQWoD3CT/wWcmD7a8xaHHHRUaAlaq+lFxMHf1ssEgODVAwjlzYqp2twkYZ7g==`。
- Node.js 归档 SHA-256：`1c4a9933a5e45bc88f54f70b5f91232c127ec49f1a5989d23fb85824c7adf9b7`（旧独立运行时的历史核验记录，归档已删除）。
- GitHub tag 源码归档 SHA-256：`176513d36ff61b809f768a7043667a19dc20e0e8fb864fd77524774bec0e852b`（历史核验记录，不随运行目录保存源码归档）。

上游 tag 名为 `v2026.7.1-2`，其仓库 `package.json` 基础版本仍写作 `2026.7.1`；官方 npm 稳定发布包版本为 `2026.7.1-2`。部署版本以 npm Release 修订号为准。

完整版本元数据见 `VERSION.json`。运行目录只保留已安装的 npm Release 和依赖；源码快照、下载缓存和独立 Node.js 运行时不属于设备运行必需内容。

## 目录说明

```text
openclaw/
├── app/                 # OpenClaw 及生产依赖
├── backups/             # 部署前的 ~/.openclaw 状态快照及 SHA-256
├── bin/openclaw         # 使用本机 Node、固定进入本目录 OpenClaw 的 CLI 包装器
├── scripts/
│   ├── deploy.sh
│   ├── backup-state.sh
│   └── sync-agents.sh
└── systemd/openclaw-gateway.service
```

STA-100 发布包运行时使用部署目录内的独立 OpenClaw 状态：`data/openclaw-home` 和 `data/openclaw-state`。旧的系统级 `~/.openclaw` 不再作为 STA-100 部署包的默认状态来源，避免与系统原有 OpenClaw 配置、Agent 同步和模型凭据互相污染。

## 部署与复核

```bash
cd /home/User/gsx/gitdir/std/openclaw
chmod +x bin/openclaw scripts/*.sh scripts/*.mjs
./scripts/deploy.sh
./scripts/sync-agents.sh

./bin/openclaw --version
./bin/openclaw gateway status --json --require-rpc
./bin/openclaw health --json
./bin/openclaw models status --json
./bin/openclaw agents list --json
systemctl --user status openclaw-gateway.service
```

`deploy.sh` 会检查本机 Node.js 版本、备份状态、安装或复用固定版本 OpenClaw，并部署启动 user systemd 服务。缺少本地安装时会通过 npm 源安装 `openclaw@2026.7.1-2`。`sync-agents.sh` 按清单创建或校正 26 个应用 Agent：2 个页面不可见的系统 Agent 和 24 个业务 Agent；重复执行时，一致的 Agent 直接跳过。

## Go 编排组件

通用编排能力已独立归档到 `/home/User/gsx/gitdir/std/openclaw-orchestrator`，包括 Go package、`openclaw-operator` CLI、manifest 同步和本机 HTTP API。STA-100 页面通过该组件接入本目录的 CLI；官方 OpenClaw 源码和发布包不包含 STA-100 业务代码。

STA-100 的 Agent 清单位于 `../sta100-web/config/sta100-agents.json`。本目录的 `scripts/sync-agents.sh` 只是兼容入口，会优先调用已构建的 Go 编排器，未构建时回退到 `go run`。

模型密钥和 Gateway 令牌不写入页面、接口响应或命令参数；在发布包中会落到部署目录的 `data/openclaw-state`，不写入系统默认 `~/.openclaw`。清单同步不会删除清单之外的 Agent 或工作区。

## Agent 分层

- `sta100-knowledge`：只整理 Go 服务检索出的本地证据，不自行联网。
- `sta100-coordinator`：根据页面功能分发结果并统一汇总，冲突数据并列保留。
- 24 个业务 Agent：承担出口、支付、物流、单据、产品、供应商、客户、市场等专业任务，可由协调链路按功能并行调用，也可在智能体页面直接对话。
- OpenClaw 自带 `main` 不属于 STA-100 页面业务 Agent。实际 OpenClaw 列表为 `main + 2 个系统 Agent + 24 个业务 Agent`，STA-100 的业务 Agent 列表接口只返回 24 个页面可见 Agent。

完整页面路由、证据处理及失败降级规则见 `../STA100-Agent分类与页面调用链路说明.md`。

## 运行目录边界

保留 `openclaw/app`、部署脚本、systemd 单元、工作区和发布包内 `data/openclaw-state` 客户状态。旧独立 Node.js、下载缓存、源码快照和上游构建目录不参与盒子运行。系统原有 OpenClaw 可卸载；STA-100 运行只应调用本目录的 `openclaw/bin/openclaw`。不得删除仍由 Codex 使用的 `/usr/local/lib/node_modules/@openai`。
