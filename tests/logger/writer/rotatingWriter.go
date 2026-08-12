package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"test/config"
	"test/utils/compress"
	"time"
)

const dateLayout = "2006-01-02" // 日期格式，用于按天归档日志文件

var maxDirSizeBytes int64  // 日志目录大小
var maxFileSizeBytes int64 // 单个日志文件大小

type RotatingWriter struct {
	mu                    sync.Mutex // 互斥锁，保证并发安全
	dirSizeExcludeLogFile int64      // 当前日志目录的大小，不含当前日志文件的大小
	currentFile           *os.File   // 当前正在写入的日志文件
	currentFileSize       int64      // 当前日志文件的大小
	currentDate           string     // 当前日志文件的日期，用于按天归档
}

// InitWriter 初始化日志写入器，创建日志目录和文件，并返回 RotatingWriter 实例。
func InitWriter(logConfig config.LogConfig) (*RotatingWriter, error) {
	// 设置日志目录和文件的大小限制
	maxDirSizeBytes = int64(logConfig.DirMaxSizeMB) * 1024 * 1024
	maxFileSizeBytes = int64(logConfig.FileMaxSizeMB) * 1024 * 1024
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
	dirSizeExcludeLogFile, err := calculateDirSizeExcludeLogFile(logConfig.DirName, logPath)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingWriter{
		dirSizeExcludeLogFile: dirSizeExcludeLogFile,
		currentFile:           file,
		currentFileSize:       info.Size(),
		currentDate:           info.ModTime().Format(dateLayout),
	}, nil
}

// Write 实现了 io.Writer 接口，用于写入日志数据。
func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// 检查日志目录大小是否超过限制，如果超过则删除最旧的文件
	if w.dirSizeExcludeLogFile+w.currentFileSize >= maxDirSizeBytes {
		go func() {
			if err := deleteOldestArchiveFiles(filepath.Dir(w.currentFile.Name()), w.currentFile.Name()); err != nil {
				fmt.Printf("日志目录清理失败: %v\n", err)
			}
		}()
	}
	// 检查是否需要归档日志文件
	if w.currentFileSize >= maxFileSizeBytes || time.Now().Format(dateLayout) != w.currentDate {
		if err := w.archiveLogFile(); err != nil {
			return 0, fmt.Errorf("日志归档失败: %v", err)
		}
	}
	n, err := w.currentFile.Write(p)
	w.currentFileSize += int64(n)
	return n, err
}

// Close 关闭当前日志文件。
func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentFile == nil {
		return nil
	}

	err := w.currentFile.Close()
	w.currentFile = nil
	return err
}

// compressAndCreateNewLogFile 关闭当前日志文件，压缩归档，并创建新的日志文件。
func (w *RotatingWriter) archiveLogFile() error {
	// 关闭当前日志文件
	if err := w.currentFile.Close(); err != nil {
		return err
	}

	// 生成下一个归档日志文件的名称
	archiveFileName, err := nextArchiveFileName(filepath.Dir(w.currentFile.Name()), w.currentDate)
	if err != nil {
		return err
	}
	// 重命名当前日志文件为归档文件
	if err := os.Rename(w.currentFile.Name(), archiveFileName); err != nil {
		return err
	}
	// 异步压缩归档日志文件，避免阻塞主线程。
	go func() {
		err := compress.CompressGzipFile(archiveFileName)
		if err == nil {
			fmt.Printf("日志归档压缩成功: %s\n", archiveFileName)
			// 压缩成功后删除原始归档文件
			os.Remove(archiveFileName)
		} else {
			fmt.Printf("日志归档压缩失败: %v\n", err)
		}
	}()

	// 创建新的日志文件
	newFile, err := os.OpenFile(w.currentFile.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	// 获取新日志文件的大小和日期
	info, err := newFile.Stat()
	if err != nil {
		return err
	}
	// 重新计算日志目录大小，排除当前日志文件的大小
	dirSizeExcludeLogFile, err := calculateDirSizeExcludeLogFile(filepath.Dir(w.currentFile.Name()), w.currentFile.Name())
	if err != nil {
		return err
	}
	// 更新 RotatingWriter 的状态
	w.currentFile = newFile
	w.currentFileSize = info.Size()
	w.currentDate = info.ModTime().Format(dateLayout)
	w.dirSizeExcludeLogFile = dirSizeExcludeLogFile
	return nil
}

// nextArchiveFileName 生成下一个归档日志文件的名称，如:2026-08-01.000.log
func nextArchiveFileName(dir string, datePart string) (string, error) {
	// 同一天内可能轮转多次，所以归档名里要带日期和序号。
	prefix := datePart + "."
	gzipSuffix := ".log.gz"

	// 遍历日志目录，查找最新归档文件，序号按“最新序号+1再取模”生成。
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	latestSeq := -1
	var latestModTime time.Time
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, gzipSuffix) {
			continue
		}

		seqText := strings.TrimSuffix(strings.TrimPrefix(name, prefix), gzipSuffix) // 获取序号部分的字符串
		seq, err := strconv.Atoi(seqText)                                           // 解析序号，将字符串转换为整数
		if err != nil {
			continue
		}
		// 获取文件信息，主要是为了获取文件的修改时间
		info, err := entry.Info()
		if err != nil {
			continue
		}
		// 记录最新修改时间的文件的序号
		if latestModTime.IsZero() || info.ModTime().After(latestModTime) {
			latestModTime = info.ModTime()
			latestSeq = seq
		}
	}
	// 计算下一个序号，确保在 0-999 范围内循环。
	nextSeq := (latestSeq + 1) % 1000
	archiveFileName := fmt.Sprintf("%s%03d.log", prefix, nextSeq)
	return filepath.Join(dir, archiveFileName), nil
}

// calculateDirSizeExcludeLogFile 计算日志目录的总大小，排除当前日志文件的大小。
func calculateDirSizeExcludeLogFile(dir string, currentLogFile string) (int64, error) {
	var totalSize int64
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if currentLogFile != "" && entry.Name() == filepath.Base(currentLogFile) {
			continue // 排除当前日志文件
		}
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		totalSize += info.Size()
	}
	return totalSize, nil
}

// deleteOldestArchiveFiles 删除最旧的文件，直到日志目录的大小低于限制。
func deleteOldestArchiveFiles(dir, excludeFile string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	// 按文件修改时间排序，最旧的文件排在前面
	sort.Slice(entries, func(i, j int) bool {
		infoI, _ := entries[i].Info()
		infoJ, _ := entries[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// 删除最旧的文件，直到日志目录的大小低于限制
	for _, entry := range entries {
		dirSize, err := calculateDirSizeExcludeLogFile(dir, excludeFile)
		if err != nil {
			return err
		}
		if dirSize < maxDirSizeBytes {
			break
		}
		if entry.IsDir() || entry.Name() == excludeFile || !strings.HasSuffix(entry.Name(), ".log.gz") {
			continue
		}
		filePath := filepath.Join(dir, entry.Name())
		if err := os.Remove(filePath); err == nil {
			fmt.Printf("已删除旧日志文件: %s\n", filePath)
		}
	}
	return nil
}
