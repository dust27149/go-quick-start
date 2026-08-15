package main

import (
	"io"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	for i := 0; i < 10; i++ {
		logger.Printf("DEBUG 初始化日志: %d", i)
	}
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("ERROR 打开日志文件失败: %v", err)
	}
	defer file.Close()
	logger.SetOutput(file)
	for i := 0; i < 10; i++ {
		logger.Printf("DEBUG 修改Output后的日志: %d", i)
	}
	io.MultiWriter(os.Stdout, file)
	logger.SetOutput(io.MultiWriter(os.Stdout, file))
	for i := 0; i < 10; i++ {
		logger.Printf("DEBUG 同时输出到控制台和文件: %d", i)
	}
}
