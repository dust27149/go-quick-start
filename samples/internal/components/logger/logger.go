package logger

import (
	"log"
	"os"
)

var Logger *log.Logger // 全局日志记录器

// Init 初始化日志系统
func init() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}
