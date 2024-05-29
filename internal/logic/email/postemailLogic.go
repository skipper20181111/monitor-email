package email

import (
	"context"
	"monitor/internal/svc"
	"monitor/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostemailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	eul    *EmailUtilLogic
}

func NewPostemailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostemailLogic {
	return &PostemailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		eul:    NewEmailUtilLogic(ctx, svcCtx),
	}
}

func (l *PostemailLogic) Postemail(req *types.PostEmailRes) (resp *types.PostEmailResp, err error) {
	resp = &types.PostEmailResp{
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
		if l.eul.SendMailRandom(emailInfo, req.Address, req.Subject, req.Body) != nil {
			resp.Msg = "not all success"
		}
	}

	return resp, nil
}
