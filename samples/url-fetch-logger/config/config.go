package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// Config 保存这个示例程序的全部可调参数。
type Config struct {
	// LogFileName 是日志文件名。
	LogFileName string
	// RequestTimeout 是单个 HTTP 请求允许等待的最长时间。
	RequestTimeout time.Duration
	// URLs 是示例程序要并发请求的网址列表。
	URLs []string
	// WorkerCount 是 worker pool 版本使用的固定 worker 数量。
	WorkerCount int
	// LogLevel 是日志等级，支持 debug、info、warn、error。
	LogLevel string
}

// FileConfig 对应 config.json 里的字段。
// 这里把请求超时时间写成秒，方便在 JSON 里直接修改。
type FileConfig struct {
	// LogFileName 是日志文件名。
	LogFileName string `json:"logFileName"`
	// RequestTimeoutSeconds 是单个 HTTP 请求允许等待的最长秒数。
	RequestTimeoutSeconds int `json:"requestTimeoutSeconds"`
	// URLs 是示例程序要并发请求的网址列表。
	URLs []string `json:"urls"`
	// WorkerCount 是 worker pool 版本使用的固定 worker 数量。
	WorkerCount int `json:"workerCount"`
	// LogLevel 是日志等级，支持 debug、info、warn、error。
	LogLevel string `json:"logLevel"`
}

// DefaultLogFileName 是默认日志文件名。
// 日志建议放到独立目录里，方便做轮转和清理，不会碰到源码文件。
const DefaultLogFileName = "logs/app.log"

// DefaultRequestTimeout 是默认的请求超时时间。
const DefaultRequestTimeout = 5 * time.Second

// DefaultWorkerCount 是 worker pool 版本的默认 worker 数量。
const DefaultWorkerCount = 3

// DefaultLogLevel 是默认日志等级。
const DefaultLogLevel = "info"

// DefaultURLs 保存示例程序默认要请求的地址。
var DefaultURLs = []string{
	"https://example.com",
	"https://golang.org",
	"https://httpbin.org/get",
}

// Default 返回一份可以直接运行的默认配置。
// 如果 JSON 配置文件缺少某些字段，程序会自动回退到这里的默认值。
func Default() Config {
	return Config{
		LogFileName:    DefaultLogFileName,
		RequestTimeout: DefaultRequestTimeout,
		URLs:           DefaultURLs,
		WorkerCount:    DefaultWorkerCount,
		LogLevel:       DefaultLogLevel,
	}
}

// LoadFromFile 从 JSON 文件读取配置，并在缺省字段上补默认值。
// 这样即使 JSON 里只改一两个字段，程序也能继续运行。
func LoadFromFile(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	var fileCfg FileConfig
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return cfg, err
	}

	if fileCfg.LogFileName != "" {
		cfg.LogFileName = fileCfg.LogFileName
	}
	if fileCfg.RequestTimeoutSeconds > 0 {
		cfg.RequestTimeout = time.Duration(fileCfg.RequestTimeoutSeconds) * time.Second
	}
	if len(fileCfg.URLs) > 0 {
		cfg.URLs = fileCfg.URLs
	}
	if fileCfg.WorkerCount > 0 {
		cfg.WorkerCount = fileCfg.WorkerCount
	}
	if fileCfg.LogLevel != "" {
		cfg.LogLevel = fileCfg.LogLevel
	}

	return cfg, nil
}

// ParseLogLevel 把配置里的字符串日志等级转换成 slog.Level。
func (c Config) ParseLogLevel() (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(c.LogLevel)); err != nil {
		return 0, fmt.Errorf("invalid log level %q: %w", c.LogLevel, err)
	}

	return level, nil
}
