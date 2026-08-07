package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// rotateSizeBytes 是单个日志文件达到这个大小时的轮转阈值。
	rotateSizeBytes int64 = 100 * 1024 * 1024
	// maxDirSizeBytes 是日志目录允许的最大总大小。
	maxDirSizeBytes int64 = 1 * 1024 * 1024 * 1024
)

// NewLogger 创建一个同时输出到终端和日志文件的 slog logger。
// 返回的 closer 需要在主程序退出前关闭，这样可以确保当前日志文件被正常刷盘和关闭。
func NewLogger(logPath string, level slog.Level) (*slog.Logger, io.Closer, error) {
	writer, err := newRotatingWriter(logPath)
	if err != nil {
		return nil, nil, err
	}

	output := io.MultiWriter(os.Stdout, writer)
	handler := slog.NewTextHandler(output, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
	appLogger := slog.New(handler)
	return appLogger, writer, nil
}

// rotatingWriter 负责把日志写到当前文件里，并在达到阈值时自动轮转压缩。
type rotatingWriter struct {
	mu            sync.Mutex
	path          string
	dir           string
	file          *os.File
	currentSize   int64
	rotateSize    int64
	maxDirSize    int64
	archivePrefix string
}

// newRotatingWriter 打开日志文件，并准备轮转所需的状态。
func newRotatingWriter(logPath string) (*rotatingWriter, error) {
	dir := filepath.Dir(logPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &rotatingWriter{
		path:          logPath,
		dir:           dir,
		file:          file,
		currentSize:   info.Size(),
		rotateSize:    rotateSizeBytes,
		maxDirSize:    maxDirSizeBytes,
		archivePrefix: filepath.Base(logPath) + ".",
	}, nil
}

// Write 将日志同时写入控制台和当前日志文件；如果文件超过阈值就会轮转压缩。
func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.rotateSize {
		if err := w.rotateLocked(); err != nil {
			return 0, err
		}
	}

	n, err := w.file.Write(p)
	w.currentSize += int64(n)
	if err != nil {
		return n, err
	}

	return n, nil
}

// Close 关闭当前活动日志文件。
func (w *rotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		return nil
	}

	err := w.file.Close()
	w.file = nil
	return err
}

// rotateLocked 把当前日志文件压缩成 .gz 归档，并重新打开一个新的空日志文件。
func (w *rotatingWriter) rotateLocked() error {
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return err
		}
	}

	archiveName, err := w.nextArchiveNameLocked()
	if err != nil {
		return err
	}
	if err := compressFile(w.path, archiveName); err != nil {
		return err
	}

	newFile, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	w.file = newFile
	w.currentSize = 0

	return w.cleanupArchivesLocked()
}

// cleanupArchivesLocked 会在日志目录总大小超过限制时，删除最早的压缩日志文件。
func (w *rotatingWriter) cleanupArchivesLocked() error {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return err
	}

	type archiveInfo struct {
		path    string
		size    int64
		modTime time.Time
	}

	var archives []archiveInfo
	var totalSize int64

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		if entry.Name() == filepath.Base(w.path) {
			totalSize += info.Size()
			continue
		}

		totalSize += info.Size()
		if strings.HasPrefix(entry.Name(), w.archivePrefix) && strings.HasSuffix(entry.Name(), ".gz") {
			archives = append(archives, archiveInfo{
				path:    filepath.Join(w.dir, entry.Name()),
				size:    info.Size(),
				modTime: info.ModTime(),
			})
		}
	}

	if totalSize <= w.maxDirSize {
		return nil
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].modTime.Before(archives[j].modTime)
	})

	for _, archive := range archives {
		if totalSize <= w.maxDirSize {
			break
		}

		if err := os.Remove(archive.path); err != nil {
			return err
		}
		totalSize -= archive.size
	}

	return nil
}

// nextArchiveNameLocked 生成当天下一份可用的压缩日志文件名。
// 文件名格式为 app.log.YYYY-MM-DD.00.log.gz，序号会根据当天已有文件自动递增。
func (w *rotatingWriter) nextArchiveNameLocked() (string, error) {
	datePart := time.Now().Format("2006-01-02")
	prefix := w.archivePrefix + datePart + "."
	suffix := ".log.gz"

	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return "", err
	}

	maxSeq := -1
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
			continue
		}

		seqText := strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
		seq, err := strconv.Atoi(seqText)
		if err != nil {
			continue
		}

		if seq > maxSeq {
			maxSeq = seq
		}
	}

	nextSeq := maxSeq + 1
	if nextSeq > 99 {
		return "", fmt.Errorf("daily log archive sequence exceeds 99: %s", datePart)
	}

	archiveName := fmt.Sprintf("%s%02d.log.gz", prefix, nextSeq)
	return filepath.Join(w.dir, archiveName), nil
}

// compressFile 把 sourcePath 压缩到 targetPath 指定的 gzip 文件里。
func compressFile(sourcePath, targetPath string) error {
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	gz := gzip.NewWriter(target)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	return nil
}
