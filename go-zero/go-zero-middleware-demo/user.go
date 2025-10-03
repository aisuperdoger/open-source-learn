package main

import (
	"flag"
	"fmt"

	"go-git-demo/go-zero-middleware-demo/internal/config"
	"go-git-demo/go-zero-middleware-demo/internal/handler"
	"go-git-demo/go-zero-middleware-demo/internal/middleware"
	"go-git-demo/go-zero-middleware-demo/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建server
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册全局中间件（所有路由都会应用）
	// 注意：中间件的注册顺序很重要，先注册的中间件会先执行
	server.Use(middleware.GlobalMiddleware)   // 全局处理：CORS、通用头、日志
	server.Use(middleware.SecurityMiddleware) // 安全处理：安全头、请求大小限制

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 注册路由处理器（包含API文件中声明的中间件）
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Println("\n=== go-zero 中间件示例服务 ===")
	fmt.Println("API端点:")
	fmt.Println("  POST   /user/login           - 用户登录（无中间件）")
	fmt.Println("  GET    /user/info/:id       - 获取用户信息（UserAgent中间件）")
	fmt.Println("  PUT    /user/:id            - 更新用户信息（Auth+Log+RateLimit中间件）")
	fmt.Println("  DELETE /user/:id            - 删除用户（Auth+Log+RateLimit中间件）")
	fmt.Println("\n测试用户:")
	fmt.Println("  用户名: user123, 密码: password123, Token: valid-token-123")
	fmt.Println("  用户名: admin456, 密码: admin123, Token: admin-token-456")
	fmt.Println("\n中间件执行顺序:")
	fmt.Println("  全局中间件 -> 安全中间件 -> 路由中间件 -> 业务逻辑")
	fmt.Println("==========================================\n")

	server.Start()
}
