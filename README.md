# Go 学习文档阅读路径

这个目录包含 4 篇 Go 学习文档。建议按顺序阅读，这样更容易把基础、并发、网络和项目结构串起来。

## 推荐阅读顺序

### 第 1 步：基础语法

[documents/01.Go语言基础入门教程.md](documents/01.Go语言基础入门教程.md)

先建立基础语法和核心概念，内容包括：

- 变量、常量、基本类型
- 条件判断和循环
- 函数
- 数组、切片、map
- struct、方法、指针、错误处理
- 包和模块管理、Go 常用命令
- 日志输出的最基础写法

### 第 2 步：并发专题

[documents/02.Go并发专题入门.md](documents/02.Go并发专题入门.md)

专门展开 Go 并发模型，重点包括：

- goroutine
- channel
- select
- 超时控制
- WaitGroup
- worker pool
- 扇出 / 扇入
- 并发中的常见坑

### 第 3 步：网络编程

[documents/03.Go网络编程入门.md](documents/03.Go网络编程入门.md)

把 Go 在网络请求和 HTTP 服务上的常用能力串起来，内容包括：

- `net/http` 基础
- GET / POST 请求
- 读取响应体
- JSON 编码和解码
- HTTP 客户端超时
- 并发请求
- HTTP 服务端基础
- 网络日志记录

### 第 4 步：项目实战模板

[documents/04.Go日志与项目实战模板.md](documents/04.Go日志与项目实战模板.md)

把前面学过的内容组织成更像真实项目的结构，重点包括：

- 目录结构怎么拆
- 配置怎么管理
- 日志怎么统一初始化
- HTTP 客户端怎么封装
- 并发任务怎么组织
- `main.go` 怎么保持清晰

## 一条最顺的学习路线

1. 先看 [documents/01.Go语言基础入门教程.md](documents/01.Go语言基础入门教程.md)，建立语法和基础概念。
2. 再看 [documents/02.Go并发专题入门.md](documents/02.Go并发专题入门.md)，把 goroutine、channel、select 彻底看明白。
3. 接着看 [documents/03.Go网络编程入门.md](documents/03.Go网络编程入门.md)，把 HTTP、JSON、超时和并发请求串起来。
4. 最后看 [documents/04.Go日志与项目实战模板.md](documents/04.Go日志与项目实战模板.md)，学习怎么把这些知识组织成小项目。

## 不同目标的读法

### 如果你主要想学 Go 语法

1. [documents/01.Go语言基础入门教程.md](documents/01.Go语言基础入门教程.md)
2. [documents/02.Go并发专题入门.md](documents/02.Go并发专题入门.md)
3. [documents/03.Go网络编程入门.md](documents/03.Go网络编程入门.md)

### 如果你主要想学并发

1. [documents/01.Go语言基础入门教程.md](documents/01.Go语言基础入门教程.md)
2. [documents/02.Go并发专题入门.md](documents/02.Go并发专题入门.md)

### 如果你主要想写网络请求或 HTTP 服务

1. [documents/01.Go语言基础入门教程.md](documents/01.Go语言基础入门教程.md)
2. [documents/02.Go并发专题入门.md](documents/02.Go并发专题入门.md)
3. [documents/03.Go网络编程入门.md](documents/03.Go网络编程入门.md)

### 如果你主要想看项目怎么组织

1. [documents/01.Go语言基础入门教程.md](documents/01.Go语言基础入门教程.md)
2. [documents/03.Go网络编程入门.md](documents/03.Go网络编程入门.md)
3. [documents/04.Go日志与项目实战模板.md](documents/04.Go日志与项目实战模板.md)

## 目录说明

- [documents](documents)：放学习文档
- [samples](samples)：放可运行示例

如果你是从零开始学 Go，建议先按“推荐阅读顺序”完整走一遍，再看 `samples` 目录里的代码。