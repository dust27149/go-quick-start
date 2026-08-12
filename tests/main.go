package main

import (
	"test/config"
	"test/logger"
	"time"
)

func main() {
	// 初始化config
	cfg := config.Cfg

	if err := logger.Init(cfg.Log); err != nil {
		panic(err)
	}
	for i := 1; i <= 20000; i++ {
		logger.Debug("程序运行第 %d 次", i)
		time.Sleep(20 * time.Millisecond)
	}
}
