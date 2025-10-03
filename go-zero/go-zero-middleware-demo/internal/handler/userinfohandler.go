package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-git-demo/go-zero-middleware-demo/internal/logic"
	"go-git-demo/go-zero-middleware-demo/internal/svc"
	"go-git-demo/go-zero-middleware-demo/internal/types"
)

func userinfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserinfoLogic(r.Context(), svcCtx)
		resp, err := l.Userinfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
