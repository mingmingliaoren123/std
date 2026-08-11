# STA-100 OpenClaw 部署说明

本目录归档 STA-100 使用的 OpenClaw 官方稳定版、独立 Node.js 运行时、部署脚本和 systemd 单元；业务 Agent 清单由使用方项目维护。

## 固定版本与来源

- OpenClaw：`2026.7.1-2`，2026-08-10 查询时 npm `latest` 与 GitHub 最新稳定 Release 均为该版本。
- 官方 Release：`https://github.com/openclaw/openclaw/releases/tag/v2026.7.1-2`，发布时间 2026-08-04。
- Node.js：`22.22.3`，Linux ARM64，满足 OpenClaw 的 `>=22.22.3 <23` 要求。
- npm 包 SHA-256：`5bb525f36f471a41239615d321c441778c7e1c007018ed6d84b795be77803276`。
- npm integrity：`sha512-ycF3yPcbjN6bUPeaUx6Mh6vze1hQWoD3CT/wWcmD7a8xaHHHRUaAlaq+lFxMHf1ssEgODVAwjlzYqp2twkYZ7g==`。
- Node.js 归档 SHA-256：`1c4a9933a5e45bc88f54f70b5f91232c127ec49f1a5989d23fb85824c7adf9b7`。
- GitHub tag 源码归档 SHA-256：`176513d36ff61b809f768a7043667a19dc20e0e8fb864fd77524774bec0e852b`。

上游 tag 名为 `v2026.7.1-2`，其仓库 `package.json` 基础版本仍写作 `2026.7.1`；官方 npm 稳定发布包版本为 `2026.7.1-2`。部署版本以 npm Release 修订号为准。

完整版本元数据见 `VERSION.json`。`source/` 是已校验 npm Release 中的完整发布包，包含可执行分发代码、文档、skills、模板、License 和第三方声明；`upstream/` 是同一 GitHub tag 的完整源码快照，包含 `src`、extensions、测试和构建配置。设备部署使用 npm Release，不依赖 Git 历史。

## 目录说明

```text
openclaw/
├── app/                 # OpenClaw 及生产依赖
├── runtime/             # Node.js 22.22.3 ARM64 独立运行时
├── source/              # 官方 npm Release 解包内容
├── upstream/            # GitHub v2026.7.1-2 完整源码快照
├── downloads/           # 已校验的官方归档
├── backups/             # 部署前的 ~/.openclaw 状态快照及 SHA-256
├── bin/openclaw         # 固定使用本目录运行时的 CLI 包装器
├── scripts/
│   ├── deploy.sh
│   ├── backup-state.sh
│   └── sync-agents.sh
└── systemd/openclaw-gateway.service
```

客户运行状态仍保存在 `~/.openclaw`。部署脚本先生成状态快照，再安装服务，不删除模型凭据、会话和客户配置。

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

`deploy.sh` 会校验两个官方归档的 SHA-256、备份状态、安装固定版本、部署并启动 user systemd 服务。`sync-agents.sh` 按清单创建或校正 24 个 Agent；重复执行时，一致的 Agent 直接跳过。

## Go 编排组件

通用编排能力已独立归档到 `/home/User/gsx/gitdir/std/openclaw-orchestrator`，包括 Go package、`openclaw-operator` CLI、manifest 同步和本机 HTTP API。STA-100 页面通过该组件接入本目录的 CLI；官方 OpenClaw 源码和发布包不包含 STA-100 业务代码。

STA-100 的 Agent 清单位于 `../sta100-web/config/sta100-agents.json`。本目录的 `scripts/sync-agents.sh` 只是兼容入口，会优先调用已构建的 Go 编排器，未构建时回退到 `go run`。

模型密钥和 Gateway 令牌不写入本目录、页面、接口响应或命令参数；清单同步不会删除清单之外的 Agent 或工作区。
