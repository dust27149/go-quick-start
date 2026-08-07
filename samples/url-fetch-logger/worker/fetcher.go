package worker

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Result 记录一次网址请求的结果。
// 主协程通过这个结构拿到每个任务的结果，避免只传一个状态码而丢失错误和耗时信息。
type Result struct {
	// URL 记录请求的目标地址。
	URL string
	// StatusCode 记录请求成功时的 HTTP 状态码。
	StatusCode int
	// Err 记录请求或读取响应体时发生的错误。
	Err error
	// Latency 记录一次请求总共花了多长时间。
	Latency time.Duration
}

// FetchURL 请求一个网址，并把最终结果通过 resultCh 发送回去。
// 这里使用 channel 回传结果，主协程就可以统一汇总所有 goroutine 的执行结果。
func FetchURL(httpClient *http.Client, logger *slog.Logger, url string, resultCh chan<- Result) {
	start := time.Now()
	logger.Debug("start request", "url", url)

	resp, err := httpClient.Get(url)
	if err != nil {
		logger.Error("request failed", "url", url, "err", err)
		resultCh <- Result{URL: url, Err: err, Latency: time.Since(start)}
		return
	}
	defer resp.Body.Close()
	logger.Debug("response received", "url", url, "status", resp.StatusCode)

	// 读完响应体可以确保请求完整结束，也避免连接资源悬挂在未消费状态。
	_, readErr := io.Copy(io.Discard, resp.Body)
	if readErr != nil {
		logger.Error("read response failed", "url", url, "err", readErr)
		resultCh <- Result{URL: url, Err: readErr, Latency: time.Since(start)}
		return
	}
	logger.Debug("response body drained", "url", url, "status", resp.StatusCode)

	resultCh <- Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Latency:    time.Since(start),
	}
}
