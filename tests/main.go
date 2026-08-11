package main

import (
	"test/config"
	"test/logger"
)

func main() {
	// 初始化config
	cfg := config.Cfg

	if err := logger.Init(cfg.Log); err != nil {
		panic(err)
	}
	defer logger.Close()
}
