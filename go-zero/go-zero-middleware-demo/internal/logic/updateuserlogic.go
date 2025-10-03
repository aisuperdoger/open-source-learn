package logic

import (
	"context"
	"fmt"
	"strings"

	"go-git-demo/go-zero-middleware-demo/internal/svc"
	"go-git-demo/go-zero-middleware-demo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UserRequest) (resp *types.UserResponse, err error) {
	// 从context中获取认证中间件传递的用户信息
	currentUserID, _ := l.ctx.Value("user-id").(string)
	token, _ := l.ctx.Value("token").(string)
	userAgent, _ := l.ctx.Value("User-Agent").(string)

	// 检查权限：用户只能更新自己的信息，或者管理员可以更新任何用户
	if !l.checkPermission(currentUserID, fmt.Sprintf("%d", req.Id)) {
		l.Errorf("用户 %s 试图更新用户 %d 的信息，权限不足", currentUserID, req.Id)
		return nil, fmt.Errorf("权限不足")
	}

	l.Infof("用户 %s 正在更新用户 %d 的信息, Token: %s, UserAgent: %s",
		currentUserID, req.Id, l.maskToken(token), userAgent)

	// 模拟更新逻辑
	resp = &types.UserResponse{
		Id:   req.Id,
		Name: fmt.Sprintf("更新后的用户_%d", req.Id),
		Age:  26, // 模拟更新年龄
	}

	l.Infof("用户信息更新成功 - UserID: %d, 操作者: %s", req.Id, currentUserID)

	return resp, nil
}

// checkPermission 检查用户权限
func (l *UpdateUserLogic) checkPermission(currentUserID, targetUserID string) bool {
	// 管理员可以更新任何用户
	if strings.HasPrefix(currentUserID, "admin") {
		return true
	}
	// 用户只能更新自己的信息
	return currentUserID == fmt.Sprintf("user%s", targetUserID)
}

// maskToken 隐藏token的部分内容
func (l *UpdateUserLogic) maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}
