package main

import (
	"fmt"
	"log"
	"net/http"
)

// 处理根路径请求的函数
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为GET
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET请求", http.StatusMethodNotAllowed)
		return
	}

	// 返回简单的问候信息
	fmt.Fprintf(w, "你好，这是一个使用Go标准库net/http的简单服务器！\n")
}

// 处理/info路径请求的函数
func infoHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	// 返回JSON格式的信息
	fmt.Fprintf(w, `{"message": "这是一个使用Go标准库net/http的示例", "version": "1.0"}`)
}

func main() {
	// 注册处理函数
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/info", infoHandler)

	// 打印启动信息
	fmt.Println("服务器启动在 http://localhost:8080")
	fmt.Println("访问 http://localhost:8080 查看问候信息")
	fmt.Println("访问 http://localhost:8080/info 查看JSON信息")
	fmt.Println("按Ctrl+C停止服务器")

	// 启动HTTP服务器，监听8080端口
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}