package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"test/config"
	"test/logger/writer"
)

var logger *log.Logger

// Init 初始化日志系统，设置日志输出目录和文件名。
func Init(logConfig config.LogConfig) error {
	// 初始化日志目录和文件
	rw, err := writer.InitWriter(logConfig)
	if err != nil {
		return err
	}

	// 创建一个新的日志记录器，输出到标准输出和日志文件
	logger = log.New(io.MultiWriter(os.Stdout, rw.CurrentFile), "", log.Ldate|log.Ltime|log.Lshortfile)
	cfgJSON, _ := json.Marshal(logConfig)
	Debug("初始化日志模块成功: %s", string(cfgJSON))
	return nil
}

// Debug 记录调试信息日志，通常用于开发和调试阶段。
func Debug(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("DEBUG "+format, v...)
}

// Info 记录普通信息日志。
func Info(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("INFO  "+format, v...)
}

// Warn 记录警告日志，但不会退出程序。
func Warn(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("WARN  "+format, v...)
}

// Error 记录错误日志，但不会退出程序。
func Error(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("ERROR "+format, v...)
}

// Fatalf 记录致命错误日志，并关闭日志文件后退出程序。
func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if logger != nil {
		if rw, ok := logger.Writer().(*writer.RotatingWriter); ok {
			_ = rw.Close()
		}
	}
	log.Fatalf(msg)
}

// ensureInitialized 确保日志系统已经初始化，如果未初始化则抛出 panic。
func ensureInitialized() {
	if logger == nil {
		panic("请先调用初始化日志系统")
	}
}
