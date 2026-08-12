package writer

import (
	"os"
	"path/filepath"
	"sync"
	"test/config"
)

type RotatingWriter struct {
	mu          sync.Mutex // 互斥锁，保证并发安全
	CurrentFile *os.File   // 当前正在写入的日志文件
	currentSize int64      // 当前日志文件的大小
}

// InitWriter 初始化日志写入器，创建日志目录和文件，并返回 RotatingWriter 实例。
func InitWriter(logConfig config.LogConfig) (*RotatingWriter, error) {
	// 创建日志目录
	if err := os.MkdirAll(logConfig.DirName, 0755); err != nil {
		return nil, err
	}
	// 打开或创建日志文件
	logPath := filepath.Join(logConfig.DirName, logConfig.FileName)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	// 获取当前日志文件的大小
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingWriter{
		CurrentFile: file,
		currentSize: info.Size(),
	}, nil
}

// Write 实现了 io.Writer 接口，用于写入日志数据。
func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err := w.CurrentFile.Write(p)
	w.currentSize += int64(n)
	return n, err
}

// Close 关闭当前日志文件。
func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.CurrentFile == nil {
		return nil
	}

	err := w.CurrentFile.Close()
	w.CurrentFile = nil
	return err
}
