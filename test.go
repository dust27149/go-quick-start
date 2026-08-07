package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIResponse struct {
	UserID    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	client := http.Client{Timeout: 5 * time.Second}                         // 设置超时
	resp, err := client.Get("https://jsonplaceholder.typicode.com/todos/1") // 发起请求
	if err != nil {
		fmt.Println("请求失败:", err) // 请求失败直接返回
		return
	}
	defer resp.Body.Close() // 关闭响应体

	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败, 状态码:", resp.StatusCode) // 状态码不为 200，直接返回
		return
	}
	var result APIResponse                           // 用于接收解析后的 JSON
	err = json.NewDecoder(resp.Body).Decode(&result) // 从响应体直接解码
	if err != nil {
		fmt.Println("JSON 解析失败:", err) // 解析失败直接返回
		return
	}

	fmt.Printf("请求成功(%d), UserID: %d, Id: %d, Title: %s, Completed: %t\n", resp.StatusCode, result.UserID, result.Id, result.Title, result.Completed) // 打印字段
}
