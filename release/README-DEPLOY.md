# STA-100 部署说明模板

## 运行方式

```bash
./scripts/sta100.sh start
```

## 控制命令

```bash
./scripts/sta100.sh start
./scripts/sta100.sh stop
./scripts/sta100.sh restart
./scripts/sta100.sh status
./scripts/sta100.sh install-autostart
./scripts/sta100.sh uninstall-autostart
```

## 默认内容

- 应用：STA-100 Go Web
- 端口：18080
- 数据目录：`./data`
- OpenClaw：`./openclaw/bin/openclaw`
- Agent 同步工具：`./bin/openclaw-operator`
- Agent 初始化：`./config/sta100-agents.json`
- OpenClaw 独立状态：`./data/openclaw-home` + `./data/openclaw-state`

## 首次启动

首次启动会创建空数据库、空账户配置和包内 OpenClaw 状态，不包含本机测试数据、历史聊天记录或旧环境配置。

本包不会使用系统已有的 `~/.openclaw` 状态目录；如果机器上曾安装过旧的系统 OpenClaw，可保留或卸载，当前应用只认本包内的 OpenClaw。

启动时默认会执行 Agent 同步，把 `config/sta100-agents.json` 中的 24 个业务 Agent 和系统 Agent 初始化到 OpenClaw。

页面 Emoji 图标字体已随应用内置，目标机器无需额外安装 Emoji 字体。

安装开机自启动：

```bash
./scripts/sta100.sh install-autostart
```

取消开机自启动：

```bash
./scripts/sta100.sh uninstall-autostart
```

如需临时改端口：

```bash
STA100_ADDR=0.0.0.0:18081 ./scripts/sta100.sh start
```

如需跳过 Agent 同步：

```bash
STA100_SYNC_AGENTS_ON_START=0 ./scripts/sta100.sh start
```
