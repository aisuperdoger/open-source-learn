package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	// 每个IP的访问记录
	visitors map[string]*Visitor
	mu       sync.RWMutex
	// 每分钟允许的最大请求数
	maxRequests int
	// 时间窗口
	timeWindow time.Duration
}

// Visitor 访问者信息
type Visitor struct {
	lastSeen time.Time
	requests []time.Time
	mu       sync.Mutex
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	m := &RateLimitMiddleware{
		visitors:    make(map[string]*Visitor),
		maxRequests: 100, // 每分钟100个请求
		timeWindow:  time.Minute,
	}

	// 启动清理协程，定期清理过期的访问者记录
	go m.cleanupVisitors()

	return m
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取客户端IP
		clientIP := m.getClientIP(r)

		// 检查是否超过限流
		if !m.isAllowed(clientIP) {
			m.writeRateLimitResponse(w, clientIP)
			return
		}

		logx.Infof("[限流中间件] 请求通过 - IP: %s, Path: %s", clientIP, r.URL.Path)

		// 调用下一个处理器
		next(w, r)
	}
}

// getClientIP 获取客户端IP
func (m *RateLimitMiddleware) getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// 检查X-Real-IP头
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// 使用RemoteAddr
	return r.RemoteAddr
}

// isAllowed 检查是否允许请求
func (m *RateLimitMiddleware) isAllowed(clientIP string) bool {
	m.mu.Lock()
	visitor, exists := m.visitors[clientIP]
	if !exists {
		visitor = &Visitor{
			lastSeen: time.Now(),
			requests: make([]time.Time, 0),
		}
		m.visitors[clientIP] = visitor
	}
	m.mu.Unlock()

	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	now := time.Now()
	visitor.lastSeen = now

	// 清理过期的请求记录
	cutoff := now.Add(-m.timeWindow)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range visitor.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	visitor.requests = validRequests

	// 检查是否超过限制
	if len(visitor.requests) >= m.maxRequests {
		return false
	}

	// 记录当前请求
	visitor.requests = append(visitor.requests, now)
	return true
}

// writeRateLimitResponse 写入限流响应
func (m *RateLimitMiddleware) writeRateLimitResponse(w http.ResponseWriter, clientIP string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)

	errorResp := map[string]interface{}{
		"code":    http.StatusTooManyRequests,
		"message": "请求过于频繁，请稍后再试",
		"data":    nil,
	}

	json.NewEncoder(w).Encode(errorResp)
	logx.Errorf("[限流中间件] 请求被限流 - IP: %s", clientIP)
}

// cleanupVisitors 定期清理过期的访问者记录
func (m *RateLimitMiddleware) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			// 清理超过10分钟未访问的记录
			for ip, visitor := range m.visitors {
				if now.Sub(visitor.lastSeen) > 10*time.Minute {
					delete(m.visitors, ip)
				}
			}
			m.mu.Unlock()
		}
	}
}
