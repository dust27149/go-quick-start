package main

import (
	"fmt"
	"net/http"
	"time"
)

type fetchResult struct {
	url string
	msg string
}

// fetch 负责请求一个网址，并把结果通过通道发回主流程。
func fetch(url string, resultCh chan fetchResult) {
	fmt.Println("请求", url)

	// 给每个请求加超时，避免某个网址一直卡住。
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		// 请求失败时，把错误信息发回主流程。
		resultCh <- fetchResult{
			url: url,
			msg: fmt.Sprintf("请求失败: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	// 请求成功时，把状态码发回主流程。
	resultCh <- fetchResult{
		url: url,
		msg: fmt.Sprintf("状态码: %d", resp.StatusCode),
	}
}

func main() {
	// 准备要并发请求的网址列表。
	urls := []string{
		"https://example.com",
		"https://golang.org",
	}

	// 用通道统一接收每个 goroutine 返回的结果。
	resultCh := make(chan fetchResult)

	for _, url := range urls {
		// 每个网址启动一个 goroutine，这样多个请求可以同时进行。
		go fetch(url, resultCh)
	}

	fmt.Println("按完成顺序输出:")
	for range urls {
		result := <-resultCh
		fmt.Printf("%s 请求结果: %s\n", result.url, result.msg)
	}
}
