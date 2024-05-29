package monitor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"monitor/internal/logic/monitor"
	"monitor/internal/svc"
)

func ShrcbmonitorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := monitor.NewShrcbmonitorLogic(r.Context(), svcCtx)
		resp, err := l.Shrcbmonitor()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
