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

// configFileName 是示例程序读取的 JSON 配置文件名。
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

// main 负责把配置、日志、HTTP 客户端和并发任务组装起来。
func main() {
	// 从 JSON 文件读取配置，让日志路径、超时时间和网址列表都可以直接调整。
	cfg, err := config.LoadFromFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在时直接退出，避免程序在默认值下悄悄运行。
			panic(fmt.Sprintf("cannot find %s", configFileName))
		}
		panic(fmt.Sprintf("failed to load %s: %v", configFileName, err))
	}

	level, err := cfg.ParseLogLevel()
	if err != nil {
		panic(err)
	}

	// 这里先输出一次调试信息，便于确认实际加载到的配置。
	// debug 模式下会显示；info 及以上时会被过滤掉。
	previewLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	previewLogger.Debug("config loaded", "logFileName", cfg.LogFileName, "requestTimeout", cfg.RequestTimeout, "workerCount", cfg.WorkerCount, "logLevel", cfg.LogLevel, "urlCount", len(cfg.URLs))

	// NewLogger 返回一个同时写终端和写日志文件的 logger。
	appLogger, logFile, err := logger.NewLogger(cfg.LogFileName, level)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// HTTP 客户端统一使用配置中的超时时间，避免请求长时间挂起。
	httpClient := client.NewHTTPClient(cfg.RequestTimeout)
	appLogger.Debug("http client created", "timeout", cfg.RequestTimeout)

	// resultCh 用来接收每个 goroutine 回传的结果。
	resultCh := make(chan worker.Result)

	appLogger.Info("program started")

	// 基础版本使用“一网址一个 goroutine”的方式，便于理解。
	for _, url := range cfg.URLs {
		appLogger.Debug("dispatch request goroutine", "url", url)
		go worker.FetchURL(httpClient, appLogger, url, resultCh)
	}

	// 主协程逐个接收结果，直到收满所有请求的返回值。
	for range cfg.URLs {
		item := <-resultCh
		appLogger.Debug("result received", "url", item.URL, "latency", item.Latency, "hasError", item.Err != nil, "status", item.StatusCode)
		printResult(appLogger, item)
	}

	appLogger.Info("program finished")
	fmt.Println("All requests completed. Check " + cfg.LogFileName + " for details.")
}
