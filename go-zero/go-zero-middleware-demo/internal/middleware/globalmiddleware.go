package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// GlobalMiddleware 全局中间件，所有请求都会经过
// 这种中间件不是通过API文件声明的，而是直接在server启动时注册的
func GlobalMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 记录请求开始时间
		start := time.Now()

		// 设置通用响应头
		w.Header().Set("X-API-Version", "1.0")
		w.Header().Set("X-Powered-By", "go-zero")

		// 跨域处理（CORS）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		logx.Infof("[全局中间件] 请求开始 - Method: %s, Path: %s, RemoteAddr: %s",
			r.Method, r.URL.Path, r.RemoteAddr)

		// 调用下一个处理器
		next(w, r)

		// 计算请求耗时
		duration := time.Since(start)
		logx.Infof("[全局中间件] 请求完成 - Path: %s, 耗时: %v", r.URL.Path, duration)
	}
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置安全相关的响应头
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 检查请求大小限制
		if r.ContentLength > 10*1024*1024 { // 10MB限制
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		logx.Infof("[安全中间件] 安全检查通过 - Path: %s", r.URL.Path)

		next(w, r)
	}
}
