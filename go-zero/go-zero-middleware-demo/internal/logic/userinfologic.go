package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-git-demo/go-zero-middleware-demo/internal/svc"
	"go-git-demo/go-zero-middleware-demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserinfoLogic {
	return &UserinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserinfoLogic) Userinfo(req *types.UserRequest) (resp *types.UserResponse, err error) {
	// 从context中获取中间件传递的数据
	userAgent, _ := l.ctx.Value("User-Agent").(string)
	requestStartTime, _ := l.ctx.Value("request-start-time").(time.Time)

	// 记录从中间件获取的信息
	l.Infof("处理用户信息请求 - UserID: %d, UserAgent: %s, 请求开始时间: %v",
		req.Id, userAgent, requestStartTime)

	// 模拟业务逻辑
	resp = &types.UserResponse{
		Id:   req.Id,
		Name: fmt.Sprintf("用户_%d", req.Id),
		Age:  25,
	}

	// 可以根据User-Agent做不同的处理
	if strings.Contains(userAgent, "Mobile") {
		l.Info("移动端用户访问")
		// 移动端特殊处理逻辑
	} else if strings.Contains(userAgent, "Bot") {
		l.Info("爬虫访问")
		// 爬虫特殊处理逻辑
	}

	return resp, nil
}
