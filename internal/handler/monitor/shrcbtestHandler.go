package monitor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"monitor/internal/logic/monitor"
	"monitor/internal/svc"
	"monitor/internal/types"
)

func ShrcbtestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShrcbMonitorRes
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := monitor.NewShrcbtestLogic(r.Context(), svcCtx)
		resp, err := l.Shrcbtest(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
