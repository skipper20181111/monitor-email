// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	email "monitor/internal/handler/email"
	monitor "monitor/internal/handler/monitor"
	refresh "monitor/internal/handler/refresh"
	"monitor/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/shrcbmonitor",
				Handler: monitor.ShrcbmonitorHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/shrcbtest",
				Handler: monitor.ShrcbtestHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/encrypt",
				Handler: monitor.EncryptHandler(serverCtx),
			},
		},
		rest.WithPrefix("/monitor"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/refresh",
				Handler: refresh.RefreshHandler(serverCtx),
			},
		},
		rest.WithPrefix("/refresh"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/postemail",
				Handler: email.PostemailHandler(serverCtx),
			},
		},
		rest.WithPrefix("/email"),
	)
}
