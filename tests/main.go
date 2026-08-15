package main

import (
	"context"
	"encoding/json"
	"os/signal"
	"syscall"
	"tests/internal/components/https"
	"tests/internal/components/logger"
	"tests/internal/components/writer"
	"tests/internal/utils/config"
	"time"
)

func main() {
	// 等待退出信号
	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 初始化配置
	cfg := config.Cfg
	data, _ := json.Marshal(cfg)
	writer.Init(cfg.Log)
	logger.Logger.Printf("DEBUG 当前配置: %s", data)
	// 初始化HTTP客户端和服务
	https.Init(cfg.Http)

	go func() {
		i := 0
		for {
			logger.Logger.Printf("DEBUG 服务正在运行...,%d\n", i)
			i++
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	// 阻塞直到收到退出信号
	<-runCtx.Done()
	logger.Logger.Println("DEBUG 收到退出信号，开始关闭服务...")
	https.DeInit()
	writer.DeInit()
	logger.Logger.Println("DEBUG 服务已关闭，退出程序。")
}
