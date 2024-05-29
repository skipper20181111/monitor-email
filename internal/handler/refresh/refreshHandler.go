package refresh

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"monitor/internal/logic/refresh"
	"monitor/internal/svc"
)

func RefreshHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := refresh.NewRefreshLogic(r.Context(), svcCtx)
		resp, err := l.Refresh()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
