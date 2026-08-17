package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Logger *log.Logger // 全局日志记录器

// init 初始化日志系统
func init() {
	Logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

// SetOutput 设置日志输出目标，如文件或其他 io.Writer
func SetOutput(w io.Writer) {
	Logger.SetOutput(w)
}

// output 输出日志信息
func output(level string, format string, v ...interface{}) {
	_ = Logger.Output(3, fmt.Sprintf(level+" "+format, v...)) // 使用 Output 方法输出日志，调用栈深度为 3，以便正确显示调用者信息
}

// Debug 输出调试级别日志
func Debug(format string, v ...interface{}) {
	output("DEBUG", format, v...)
}

// Info 输出信息级别日志
func Info(format string, v ...interface{}) {
	output("INFO ", format, v...)
}

// Warn 输出警告级别日志
func Warn(format string, v ...interface{}) {
	output("WARN ", format, v...)
}

// Error 输出错误级别日志
func Error(format string, v ...interface{}) {
	output("ERROR", format, v...)
}
