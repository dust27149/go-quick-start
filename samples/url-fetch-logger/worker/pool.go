package worker

import (
	"log/slog"
	"net/http"
	"sync"
)

// FetchURLsWithPool 使用固定数量的 worker 处理网址列表。
// 这个版本的重点是“固定 worker 数量”，适合你理解任务队列和资源复用。
func FetchURLsWithPool(httpClient *http.Client, logger *slog.Logger, urls []string, workerCount int) []Result {
	jobs := make(chan string)
	resultCh := make(chan Result)
	var wg sync.WaitGroup

	// worker 函数不断从 jobs 中取任务，直到 jobs 被关闭。
	// 每个 worker 处理完一个网址后，会继续等待下一个网址。
	workerFn := func(id int) {
		defer wg.Done()
		logger.Debug("worker started", "worker", id)

		for url := range jobs {
			logger.Debug("worker picked job", "worker", id, "url", url)
			FetchURL(httpClient, logger, url, resultCh)
		}

		logger.Debug("worker finished", "worker", id)
	}

	// 启动固定数量的 worker。
	// workerCount 越大，并发能力越强，但也会占用更多资源。
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go workerFn(i + 1)
	}

	// 把所有 URL 送进 jobs 通道，随后关闭通道告诉 worker 没有新任务了。
	// 这里用一个单独的 goroutine 发送任务，避免主流程被发送操作阻塞。
	go func() {
		for _, url := range urls {
			logger.Debug("enqueue job", "url", url)
			jobs <- url
		}
		close(jobs)
		logger.Debug("job queue closed")
	}()

	// 等所有 worker 结束后关闭 resultCh，方便下面 range 读完所有结果。
	go func() {
		wg.Wait()
		close(resultCh)
		logger.Debug("result channel closed")
	}()

	// 主协程通过 range 读取所有结果，直到 resultCh 被关闭。
	results := make([]Result, 0, len(urls))
	for item := range resultCh {
		results = append(results, item)
	}

	return results
}
