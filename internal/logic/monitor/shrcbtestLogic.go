package monitor

import (
	"context"

	"monitor/internal/svc"
	"monitor/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShrcbtestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShrcbtestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShrcbtestLogic {
	return &ShrcbtestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShrcbtestLogic) Shrcbtest(req *types.ShrcbMonitorRes) (resp *types.ShrcbMonitorResp, err error) {
	return &types.ShrcbMonitorResp{Msg: "ok"}, nil
}
