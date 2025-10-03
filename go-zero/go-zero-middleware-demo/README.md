# go-zero 自定义中间件完全指南

## 概述

本项目演示了go-zero框架中自定义中间件的各种实现方式和使用场景。go-zero提供了两种主要的中间件注册方式：

1. **API文件声明的中间件** - 通过在.api文件中使用`middleware`关键字声明
2. **全局中间件** - 在server启动时通过`server.Use()`方法注册

## 项目结构

```
go-zero-middleware-demo/
├── user.api                          # API定义文件
├── user.go                           # 主程序入口
├── etc/
│   └── user.yaml                     # 配置文件
└── internal/
    ├── config/                       # 配置定义
    ├── types/                        # 类型定义
    ├── svc/                         # 服务上下文
    ├── handler/                     # HTTP处理器
    ├── logic/                       # 业务逻辑
    └── middleware/                  # 中间件实现
        ├── useragentmiddleware.go   # 用户代理中间件
        ├── authmiddleware.go        # 认证中间件
        ├── logmiddleware.go         # 日志中间件
        ├── ratelimitmiddleware.go   # 限流中间件
        └── globalmiddleware.go      # 全局中间件
```

## 中间件类型

### 1. API文件声明的中间件

在API文件中通过`@server(middleware: MiddlewareName)`声明：

```go
// 单个中间件
@server(
    middleware: UserAgentMiddleware
)
service user {
    @handler userinfo
    get /user/info/:id (UserRequest) returns (UserResponse)
}

// 多个中间件（按顺序执行）
@server(
    middleware: AuthMiddleware,LogMiddleware,RateLimitMiddleware
)
service user {
    @handler updateUser
    put /user/:id (UserRequest) returns (UserResponse)
}
```

### 2. 全局中间件

在main.go中通过`server.Use()`注册：

```go
func main() {
    server := rest.MustNewServer(c.RestConf)
    
    // 注册全局中间件（所有路由都会应用）
    server.Use(middleware.GlobalMiddleware)
    server.Use(middleware.SecurityMiddleware)
    
    // 注册路由处理器
    handler.RegisterHandlers(server, ctx)
    
    server.Start()
}
```

## 中间件实现示例

### 1. UserAgent中间件 - 提取和存储用户代理信息

```go
func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 获取User-Agent信息
        userAgent := r.Header.Get("User-Agent")
        if userAgent == "" {
            userAgent = "Unknown"
        }

        // 存储到context中
        ctx := context.WithValue(r.Context(), "User-Agent", userAgent)
        newReq := r.WithContext(ctx)

        // 记录日志
        logx.Infof("[UserAgent中间件] 请求开始 - UserAgent: %s", userAgent)

        // 调用下一个处理器
        next(w, newReq)
    }
}
```

### 2. 认证中间件 - JWT Token验证

```go
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 获取Authorization请求头
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            m.writeErrorResponse(w, "Invalid Authorization header", http.StatusUnauthorized)
            return
        }

        // 验证token
        token := strings.TrimPrefix(auth, "Bearer ")
        userID, err := m.validateToken(token)
        if err != nil {
            m.writeErrorResponse(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // 存储用户信息到context
        ctx := context.WithValue(r.Context(), "user-id", userID)
        ctx = context.WithValue(ctx, "token", token)
        newReq := r.WithContext(ctx)

        next(w, newReq)
    }
}
```

### 3. 日志中间件 - 请求响应记录

```go
func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // 读取请求体
        var requestBody []byte
        if r.Method == "POST" || r.Method == "PUT" {
            requestBody, _ = io.ReadAll(r.Body)
            r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
        }

        // 包装响应写入器以捕获响应数据
        wrappedWriter := &responseWriter{
            ResponseWriter: w,
            body:          &bytes.Buffer{},
            statusCode:    200,
        }

        // 记录请求信息
        logData := map[string]interface{}{
            "method":      r.Method,
            "path":        r.URL.Path,
            "remote_addr": r.RemoteAddr,
            "user_agent":  r.Header.Get("User-Agent"),
        }

        next(wrappedWriter, r)

        // 记录响应信息
        duration := time.Since(start)
        logx.Infof("请求完成 - 耗时: %v, 状态码: %d", duration, wrappedWriter.statusCode)
    }
}
```

### 4. 限流中间件 - 基于IP的请求频率控制

```go
func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        clientIP := m.getClientIP(r)
        
        if !m.isAllowed(clientIP) {
            m.writeRateLimitResponse(w, clientIP)
            return
        }

        next(w, r)
    }
}

func (m *RateLimitMiddleware) isAllowed(clientIP string) bool {
    // 实现滑动窗口限流算法
    // 每分钟最多100个请求
    visitor := m.getOrCreateVisitor(clientIP)
    
    now := time.Now()
    cutoff := now.Add(-m.timeWindow)
    
    // 清理过期请求记录
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
```

## 中间件执行顺序

中间件的执行遵循以下顺序：

1. **全局中间件**（按注册顺序）
2. **路由中间件**（按API文件中声明的顺序）
3. **业务处理器**

示例执行流程：
```
请求 -> GlobalMiddleware -> SecurityMiddleware -> AuthMiddleware -> LogMiddleware -> RateLimitMiddleware -> Handler -> 响应
```

## 在业务逻辑中使用中间件数据

在logic层可以通过context获取中间件传递的数据：

```go
func (l *UserInfoLogic) UserInfo(req *types.UserRequest) (*types.UserResponse, error) {
    // 获取中间件传递的数据
    userAgent, _ := l.ctx.Value("User-Agent").(string)
    userID, _ := l.ctx.Value("user-id").(string)
    
    l.Infof("处理用户信息请求 - UserAgent: %s, CurrentUser: %s", userAgent, userID)
    
    // 业务逻辑处理...
    return resp, nil
}
```

## API测试示例

### 1. 用户登录（无中间件）
```bash
curl -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user123","password":"password123"}'
```

### 2. 获取用户信息（UserAgent中间件）
```bash
curl -X GET http://localhost:8888/user/info/123 \
  -H "User-Agent: Mozilla/5.0 (Mobile; Android)"
```

### 3. 更新用户信息（需要认证+日志+限流中间件）
```bash
curl -X PUT http://localhost:8888/user/123 \
  -H "Authorization: Bearer valid-token-123" \
  -H "Content-Type: application/json"
```

## 最佳实践

### 1. 中间件设计原则
- **单一职责**：每个中间件只负责一个特定功能
- **无状态**：避免在中间件中存储状态，使用context传递数据
- **错误处理**：妥善处理错误，提供清晰的错误信息
- **性能考虑**：避免在中间件中进行耗时操作

### 2. Context使用规范
```go
// 推荐：使用有意义的key
ctx := context.WithValue(r.Context(), "user-id", userID)

// 不推荐：使用字符串常量
ctx := context.WithValue(r.Context(), "uid", userID)
```

### 3. 错误响应格式
```go
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    errorResp := map[string]interface{}{
        "code":    statusCode,
        "message": message,
        "data":    nil,
    }
    
    json.NewEncoder(w).Encode(errorResp)
}
```

### 4. 日志记录规范
```go
// 使用结构化日志
logx.WithFields(logx.Field("user_id", userID), logx.Field("path", r.URL.Path)).
    Info("用户认证成功")
```

### 5. 性能优化建议
- 使用连接池管理数据库连接
- 实现中间件结果缓存（如用户权限）
- 合理设置限流参数
- 定期清理过期数据

## 配置中间件

通过配置文件可以控制内置中间件的启用状态：

```yaml
# etc/user.yaml
Name: user.api
Host: 0.0.0.0
Port: 8888

# 中间件配置
Middlewares:
  Trace: true      # 链路追踪
  Log: true        # 日志记录
  Prometheus: true # 指标监控
  MaxConns: true   # 连接限制
  Breaker: true    # 熔断器
  Shedding: true   # 负载保护
  Timeout: true    # 超时处理
  Recover: true    # 异常恢复
  Metrics: true    # 指标统计
  MaxBytes: true   # 请求大小限制
  Gunzip: true     # 压缩处理
```

## 总结

go-zero的中间件系统提供了灵活而强大的请求处理能力：

1. **API声明式中间件**：适用于特定路由组的功能需求
2. **全局中间件**：适用于所有请求的通用处理
3. **内置中间件**：提供了丰富的开箱即用功能
4. **Context传递**：实现中间件间的数据共享
5. **性能考量**：支持高并发和低延迟的处理

通过合理使用这些中间件功能，可以构建出安全、高效、可维护的微服务应用。