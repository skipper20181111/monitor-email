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
}

func NewPostemailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostemailLogic {
	return &PostemailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostemailLogic) Postemail(req *types.PostEmailRes) (resp *types.PostEmailResp, err error) {
	// todo: add your logic here and delete this line

	return
}
