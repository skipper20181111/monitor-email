package email

import (
	"context"

	"monitor/internal/svc"
	"monitor/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EasyemailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	eul    *EmailUtilLogic
}

func NewEasyemailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EasyemailLogic {
	return &EasyemailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		eul:    NewEmailUtilLogic(ctx, svcCtx),
	}
}

func (l *EasyemailLogic) Easyemail(req *types.EasyEmailRes) (resp *types.EasyEmailResp, err error) {
	resp = &types.EasyEmailResp{
		Code: "10000",
		Msg:  "",
		Data: &types.PostEmailRp{
			Success: true,
		},
	}
	EmailInfoList := make([]*types.EmailInfo, 0)
	emaillist, ok := l.svcCtx.LocalCache.Get(svc.EmailListKey)
	if ok {
		EmailInfoList = emaillist.([]*types.EmailInfo)
	}
	for _, emailInfo := range EmailInfoList {
		EsayEmailInfo := &types.EmailInfo{
			Host:      emailInfo.Host,
			Port:      emailInfo.Port,
			Send2Who:  req.Send2Who,
			EmailUser: req.EmailUser,
		}
		if l.eul.SendMailRandom(EsayEmailInfo, req.Send2Who, req.Subject, req.Body) != nil {
			resp.Msg = "not all success"
		}
	}
	return resp, nil
}
