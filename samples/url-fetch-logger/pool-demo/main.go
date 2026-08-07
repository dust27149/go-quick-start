package main

import (
	"fmt"
	"log/slog"
	"os"

	"url-fetch-logger/client"
	"url-fetch-logger/config"
	"url-fetch-logger/logger"
	"url-fetch-logger/worker"
)

// configFileName 是 worker pool 示例读取的 JSON 配置文件名。
const configFileName = "config.json"

// printResult 按日志等级输出一条请求结果。
func printResult(appLogger *slog.Logger, item worker.Result) {
	if item.Err != nil {
		appLogger.Error("request failed", "url", item.URL, "latency", item.Latency, "err", item.Err)
		return
	}

	if item.StatusCode >= 400 {
		appLogger.Warn("request finished with non-success status", "url", item.URL, "status", item.StatusCode, "latency", item.Latency)
		return
	}

	appLogger.Info("request finished", "url", item.URL, "status", item.StatusCode, "latency", item.Latency)
}

// main 展示固定 worker 数量的 worker pool 写法。
func main() {
	// 读取 JSON 配置，复用基础示例同一套配置文件。
	cfg, err := config.LoadFromFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			// 让新手一眼看出缺少配置文件时发生了什么。
			panic(fmt.Sprintf("cannot find %s", configFileName))
		}
		panic(fmt.Sprintf("failed to load %s: %v", configFileName, err))
	}

	level, err := cfg.ParseLogLevel()
	if err != nil {
		panic(err)
	}

	// 先用一个临时 logger 打印配置摘要，便于确认当前运行参数。
	previewLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	previewLogger.Debug("config loaded", "logFileName", cfg.LogFileName, "requestTimeout", cfg.RequestTimeout, "workerCount", cfg.WorkerCount, "logLevel", cfg.LogLevel, "urlCount", len(cfg.URLs))

	// 统一创建日志对象，保证基础示例和 worker pool 示例的日志格式一致。
	appLogger, logFile, err := logger.NewLogger(cfg.LogFileName, level)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// 使用配置里的请求超时时间创建 HTTP 客户端。
	httpClient := client.NewHTTPClient(cfg.RequestTimeout)
	appLogger.Debug("http client created", "timeout", cfg.RequestTimeout)

	appLogger.Info("worker pool demo started")

	// 使用固定数量的 worker 处理任务，适合观察“任务队列 + worker”模型。
	appLogger.Debug("start worker pool dispatch", "workerCount", cfg.WorkerCount, "urlCount", len(cfg.URLs))
	results := worker.FetchURLsWithPool(httpClient, appLogger, cfg.URLs, cfg.WorkerCount)
	for _, item := range results {
		appLogger.Debug("result received", "url", item.URL, "latency", item.Latency, "hasError", item.Err != nil, "status", item.StatusCode)
		printResult(appLogger, item)
	}

	appLogger.Info("worker pool demo finished")
	fmt.Println("Worker pool demo completed. Check " + cfg.LogFileName + " for details.")
}
