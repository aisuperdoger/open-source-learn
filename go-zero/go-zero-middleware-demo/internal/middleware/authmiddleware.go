package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取Authorization请求头
		auth := r.Header.Get("Authorization")
		if auth == "" {
			m.writeErrorResponse(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 检查Bearer token格式
		if !strings.HasPrefix(auth, "Bearer ") {
			m.writeErrorResponse(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// 提取token
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			m.writeErrorResponse(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// 简单的token验证（实际项目中应该使用JWT解析）
		userID, err := m.validateToken(token)
		if err != nil {
			m.writeErrorResponse(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// 将用户ID存储到context中
		ctx := context.WithValue(r.Context(), "user-id", userID)
		ctx = context.WithValue(ctx, "token", token)
		newReq := r.WithContext(ctx)

		logx.Infof("[认证中间件] 用户认证成功 - UserID: %s, Path: %s", userID, r.URL.Path)

		// 调用下一个处理器
		next(w, newReq)
	}
}

// validateToken 简单的token验证（实际项目中应该使用JWT解析）
func (m *AuthMiddleware) validateToken(token string) (string, error) {
	// 这里仅作为示例，实际项目中应该使用JWT解析
	if token == "valid-token-123" {
		return "user123", nil
	}
	if token == "admin-token-456" {
		return "admin456", nil
	}
	return "", fmt.Errorf("invalid token")
}

// writeErrorResponse 写入错误响应
func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := map[string]interface{}{
		"code":    statusCode,
		"message": message,
		"data":    nil,
	}

	json.NewEncoder(w).Encode(errorResp)
	logx.Errorf("[认证中间件] 认证失败: %s", message)
}
