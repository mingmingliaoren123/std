# Go 开发环境说明

更新时间：2026-08-10

## 已安装组件

| 组件 | 版本/位置 |
|---|---|
| Go | `go1.26.5 linux/arm64` |
| GOROOT | `/usr/local/go1.26.5` |
| 固定入口 | `/usr/local/go -> /usr/local/go1.26.5` |
| go | `/usr/local/bin/go` |
| gofmt | `/usr/local/bin/gofmt` |
| GOPATH | `/home/User/go` |
| gopls | `v0.23.0`，位于 `/home/User/go/bin/gopls` |
| VS Code Go 扩展 | `golang.go 0.56.0` |

安装包来自 Go 官方中国镜像域名，安装前已按官方值完成 SHA-256 校验。

## Go 环境配置

```text
GOPROXY=https://goproxy.cn,direct
GOTOOLCHAIN=local
```

`proxy.golang.org` 在本机网络中连接超时，因此使用可正常访问的国内模块代理。`GOTOOLCHAIN=local` 用于确保项目明确使用本机安装的 Go 工具链。

`/home/User/go/bin` 已加入 `/home/User/.profile` 和 `/home/User/.bashrc`。新打开的终端和 VS Code 终端会自动加载；当前已打开的旧终端可重新打开后使用。

## 环境检查

```bash
go version
go env GOROOT GOPATH GOPROXY GOTOOLCHAIN
gopls version
```

## STA-100 页面项目检查与编译

```bash
cd /home/User/gsx/gitdir/std/sta100-web
gofmt -w main.go
go test ./...
go vet ./...
mkdir -p dist
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -trimpath -ldflags='-s -w' \
  -o dist/sta100-web-linux-arm64 .
```

使用 `CGO_ENABLED=0` 生成静态链接 ARM64 单文件，适合后续设备交付和压缩包归档。

## 当前编译结果

```text
文件：sta100-web/dist/sta100-web-linux-arm64
类型：ELF 64-bit ARM aarch64，statically linked，stripped
大小：约 6.1 MB
SHA-256：c54451a6abd4a812cf728c8917ef8bab4674d63a036902e83a8cb469f0cd4cc0
```

该可执行文件已实际启动验证：页面入口返回 HTTP 200，`GET /api/health` 返回 `status: ok`。
