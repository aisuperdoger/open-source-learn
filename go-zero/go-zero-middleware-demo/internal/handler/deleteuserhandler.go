package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-git-demo/go-zero-middleware-demo/internal/logic"
	"go-git-demo/go-zero-middleware-demo/internal/svc"
	"go-git-demo/go-zero-middleware-demo/internal/types"
)

func deleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDeleteUserLogic(r.Context(), svcCtx)
		err := l.DeleteUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
