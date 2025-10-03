package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"go-git-demo/go-zero-middleware-demo/internal/config"
	"go-git-demo/go-zero-middleware-demo/internal/middleware"
)

type ServiceContext struct {
	Config              config.Config
	UserAgentMiddleware rest.Middleware
	AuthMiddleware      rest.Middleware
	LogMiddleware       rest.Middleware
	RateLimitMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		UserAgentMiddleware: middleware.NewUserAgentMiddleware().Handle,
		AuthMiddleware:      middleware.NewAuthMiddleware().Handle,
		LogMiddleware:       middleware.NewLogMiddleware().Handle,
		RateLimitMiddleware: middleware.NewRateLimitMiddleware().Handle,
	}
}
