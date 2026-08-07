# worker 包说明

这个包负责实际的并发请求逻辑。

这里有两种写法：

- `FetchURL`：一网址一个 goroutine 的基础版本
- `FetchURLsWithPool`：固定 worker 数量的 worker pool 版本

如果你是新手，建议先看基础版本，再对照 worker pool 版本理解任务如何被分配和回收。