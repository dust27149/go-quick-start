# URL 并发请求示例

这个目录里放了一个适合新手学习的 Go 示例，主题是“并发请求网址并写日志文件”。

## 包含的内容

- `config.json`：示例配置文件
- `config`：默认参数
- `logger`：日志初始化
- `client`：HTTP 客户端封装
- `worker`：并发请求逻辑
- `pool-demo`：worker pool 版本示例

## 运行基础版本

在这个目录下执行：

```powershell
go run .
```

## 运行 worker pool 版本

在这个目录下执行：

```powershell
go run ./pool-demo
```

运行后会生成 `logs/app.log`，里面记录了每个请求的开始、结束、失败信息和耗时。

## 修改配置

你可以直接编辑同目录下的 `config.json`，调整：

- `logFileName`：日志输出路径
- `requestTimeoutSeconds`：单个请求超时时间，单位是秒
- `workerCount`：worker pool 使用的 worker 数量
- `urls`：要请求的网址列表

## 日志等级

这个示例使用的是 Go 标准库的 `slog`，所以可以像 Java 那样区分日志等级：

- `debug`：最详细，适合排查问题
- `info`：普通运行信息
- `warn`：有异常但程序还能继续
- `error`：请求失败或发生错误

如果你想多看一些中间过程，可以把 `config.json` 里的 `logLevel` 改成 `debug`；如果想少一点输出，可以改成 `info`、`warn` 或 `error`。