# OpenClaw Orchestrator

独立的 Go 编排组件，用于把应用需要的 OpenClaw 管理能力与官方运行时解耦。组件不修改 OpenClaw 官方源码，也不保存模型 API Key 或 Gateway token。

## 职责边界

```text
openclaw/                    官方 OpenClaw、Node 运行时、Gateway 和部署脚本
openclaw-orchestrator/       通用 Go package、CLI、HTTP API、manifest 同步
sta100-web/                  STA-100 页面、业务接口和 STA-100 Agent 清单
```

`orchestrator` package 可被其它 Go 项目直接导入；`openclaw-operator` CLI 适用于部署脚本、运维和不使用 HTTP 的应用；`serve` 提供本机管理 API。

## 构建

```bash
cd /home/User/gsx/gitdir/std/openclaw-orchestrator
go test ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -trimpath -ldflags='-s -w' \
  -o bin/openclaw-operator ./cmd/openclaw-operator
```

## CLI

```bash
./bin/openclaw-operator \
  --bin ../openclaw/bin/openclaw \
  --manifest ../sta100-web/config/sta100-agents.json status

./bin/openclaw-operator --bin ../openclaw/bin/openclaw models
./bin/openclaw-operator --bin ../openclaw/bin/openclaw agents list
./bin/openclaw-operator --bin ../openclaw/bin/openclaw \
  --manifest ../sta100-web/config/sta100-agents.json agents sync
```

默认情况下，组件通过 `OPENCLAW_BIN`、`OPENCLAW_CONFIG_PATH`、`OPENCLAW_AGENT_MANIFEST` 和相邻目录自动发现 OpenClaw。生产部署建议显式传入路径。

## HTTP API

```bash
./bin/openclaw-operator \
  --bin ../openclaw/bin/openclaw \
  --manifest ../sta100-web/config/sta100-agents.json \
  --addr 127.0.0.1:18790 serve
```

| 方法 | 路径 | 作用 |
|---|---|---|
| GET | `/health` | 编排器健康检查 |
| GET | `/v1/status` | Gateway、RPC 和配置审计状态 |
| GET | `/v1/models` | 固定 OpenClaw 版本模型目录、默认模型和非敏感凭据状态 |
| PUT | `/v1/models/default` | 校验可用模型后切换默认模型 |
| POST | `/v1/models/auth` | 通过 stdin 写入 API Key，不回显密钥 |
| GET | `/v1/agents` | 读取 Agent 公共字段 |
| POST | `/v1/agents/sync` | 按调用方 manifest 幂等创建/校正 Agent |

写接口要求 `X-OpenClaw-Operator-Request: 1`。HTTP 服务默认只监听 `127.0.0.1`，部署到其它地址前必须由上层网关提供认证和访问控制。

模型目录固定对应 OpenClaw `2026.7.1-2`，包含 20 个 API Key 提供商和 106 个文本模型；读取目录时不调用 `models list`，只调用 `models status` 获取当前默认模型和凭据状态。升级 OpenClaw 时必须同步更新 `orchestrator/model_catalog.go`。

## Manifest

Agent 清单属于使用方项目，不属于 `openclaw/` 官方目录。清单至少包含 `id`、`name`、`workspace`，可通过 `workspace_root` 指定工作区根目录；相对路径以清单目录为基准。

```json
{
  "schema_version": 1,
  "workspace_root": "../../openclaw",
  "agents": [
    {
      "id": "example-agent",
      "name": "示例助手",
      "technical_name": "ExampleAgent",
      "emoji": "🤖",
      "workspace": "workspaces/example-agent",
      "model": "provider/model"
    }
  ]
}
```

同步只创建不存在的工作区和 `IDENTITY.md`，通过官方 CLI 执行 `agents add` 与 `agents set-identity`。它不会删除清单外 Agent，也不会删除工作区；因此适合重复执行和后续迁移。
