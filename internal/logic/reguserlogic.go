package logic

import (
	"context"

	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/user_mgr_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegUserLogic {
	return &RegUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegUserLogic) RegUser(in *user_mgr_pb.RegUserReq) (*user_mgr_pb.RegUserRsp, error) {
	// todo: add your logic here and delete this line

	return &user_mgr_pb.RegUserRsp{}, nil
}
