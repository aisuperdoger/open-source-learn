package logic

import (
	"context"
	"fmt"

	"go-git-demo/go-zero-middleware-demo/internal/svc"
	"go-git-demo/go-zero-middleware-demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 简单的用户名密码验证（实际项目中应该从数据库查询并验证密码哈希）
	var token string
	var userID int64

	switch req.Username {
	case "user123":
		if req.Password != "password123" {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		token = "valid-token-123"
		userID = 123
	case "admin456":
		if req.Password != "admin123" {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		token = "admin-token-456"
		userID = 456
	default:
		return nil, fmt.Errorf("用户不存在")
	}

	l.Infof("用户登录成功 - Username: %s, UserID: %d", req.Username, userID)

	// 返回登录响应
	resp = &types.LoginResponse{
		Token: token,
		User: types.UserResponse{
			Id:   userID,
			Name: req.Username,
			Age:  25,
		},
	}

	return resp, nil
}
