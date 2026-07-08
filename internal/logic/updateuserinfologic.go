package logic

import (
	"context"

	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/user_mgr_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user_mgr_pb.UpdateUserInfoReq) (*user_mgr_pb.UpdateUserInfoRsp, error) {
	// todo: add your logic here and delete this line

	return &user_mgr_pb.UpdateUserInfoRsp{}, nil
}
