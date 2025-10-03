package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogMiddleware struct {
}

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 读取请求体（仅对POST/PUT请求）
		var requestBody []byte
		if r.Method == "POST" || r.Method == "PUT" {
			if r.Body != nil {
				requestBody, _ = io.ReadAll(r.Body)
				// 重新设置请求体，以便后续处理器可以读取
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}
		}

		// 创建响应写入器包装器，用于捕获响应数据
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			statusCode:     200, // 默认状态码
		}

		// 记录请求开始信息
		logData := map[string]interface{}{
			"timestamp":    start.Format(time.RFC3339),
			"method":       r.Method,
			"path":         r.URL.Path,
			"query":        r.URL.RawQuery,
			"remote_addr":  r.RemoteAddr,
			"user_agent":   r.Header.Get("User-Agent"),
			"content_type": r.Header.Get("Content-Type"),
		}

		// 如果有请求体，记录请求体内容
		if len(requestBody) > 0 {
			logData["request_body"] = string(requestBody)
		}

		// 从 context 中获取用户ID（如果有的话）
		if userID := r.Context().Value("user-id"); userID != nil {
			logData["user_id"] = userID
		}

		logx.Infof("[日志中间件] 请求开始: %s", m.formatLogData(logData))

		// 调用下一个处理器
		next(wrappedWriter, r)

		// 计算耗时
		duration := time.Since(start)

		// 记录响应信息
		responseData := map[string]interface{}{
			"timestamp":     time.Now().Format(time.RFC3339),
			"method":        r.Method,
			"path":          r.URL.Path,
			"status_code":   wrappedWriter.statusCode,
			"duration_ms":   duration.Milliseconds(),
			"response_size": wrappedWriter.body.Len(),
		}

		// 如果响应体不太大，记录响应内容
		if wrappedWriter.body.Len() > 0 && wrappedWriter.body.Len() < 1000 {
			responseData["response_body"] = wrappedWriter.body.String()
		}

		logx.Infof("[日志中间件] 请求完成: %s", m.formatLogData(responseData))

		// 将响应写入原始响应写入器
		w.WriteHeader(wrappedWriter.statusCode)
		w.Write(wrappedWriter.body.Bytes())
	}
}

// responseWriter 响应写入器包装器
type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

// formatLogData 格式化日志数据
func (m *LogMiddleware) formatLogData(data map[string]interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "failed to marshal log data"
	}
	return string(jsonData)
}
