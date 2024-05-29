package email

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"monitor/internal/logic/email"
	"monitor/internal/svc"
	"monitor/internal/types"
)

func PostemailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostEmailRes
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := email.NewPostemailLogic(r.Context(), svcCtx)
		resp, err := l.Postemail(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
