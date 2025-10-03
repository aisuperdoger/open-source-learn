package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAgentMiddleware struct {
}

func NewUserAgentMiddleware() *UserAgentMiddleware {
	return &UserAgentMiddleware{}
}

func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取开始时间，用于计算请求耗时
		start := time.Now()

		// 从请求头获取User-Agent信息
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			userAgent = "Unknown"
		}

		// 将User-Agent信息存储到context中
		ctx := context.WithValue(r.Context(), "User-Agent", userAgent)
		ctx = context.WithValue(ctx, "request-start-time", start)

		// 创建新的请求对象，包含更新后的context
		newReq := r.WithContext(ctx)

		// 记录请求开始日志
		logx.Infof("[UserAgent中间件] 请求开始 - Method: %s, Path: %s, UserAgent: %s",
			r.Method, r.URL.Path, userAgent)

		// 调用下一个处理器
		next(w, newReq)

		// 计算请求耗时
		duration := time.Since(start)
		logx.Infof("[UserAgent中间件] 请求完成 - Path: %s, 耗时: %v", r.URL.Path, duration)
	}
}
