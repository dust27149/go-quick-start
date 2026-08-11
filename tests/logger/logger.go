package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"test/config"
)

var logger *log.Logger
var file *os.File

func Init(logConfig config.LogConfig) error {
	if err := os.MkdirAll(logConfig.DirName, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(logConfig.DirName, logConfig.FileName)
	var err error
	file, err = os.OpenFile(
		logPath,                             // 文件名：日志输出到当前目录下 app.log
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, // 打开方式：不存在则创建、只写、追加写入
		0666,                                // 文件权限：rw-rw-rw-（实际权限还会受 umask 影响）
	)
	if err != nil {
		return err
	}

	writer := io.MultiWriter(os.Stdout, file) // 同时写控制台和文件
	logger = log.New(
		writer,                             // 输出目标：写入日志文件
		"",                                 // 日志前缀
		log.Ldate|log.Ltime|log.Lshortfile, // 日志格式：日期 + 时间 + 短文件名:行号
	)

	cfgJSON, _ := json.Marshal(logConfig)
	Debug("配置初始化成功: %s", string(cfgJSON))
	return nil
}

func Close() error {
	if file == nil {
		return nil
	}
	err := file.Close()
	file = nil
	return err
}

func ensureInitialized() {
	if logger == nil {
		panic("请先调用初始化日志系统")
	}
}

func Debug(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("DEBUG "+format, v...)
}

func Info(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("INFO  "+format, v...)
}

func Warn(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("WARN  "+format, v...)
}

func Error(format string, v ...interface{}) {
	ensureInitialized()
	logger.Printf("ERROR "+format, v...)
}

func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	if file != nil {
		_ = file.Close()
	}
	log.Fatalf(msg)
}
