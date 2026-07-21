package logic

import (
	"context"

	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/user_mgr_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ValidStatusValid   = 1 // 有效
	ValidStatusInvalid = 2 // 无效
)

type CheckUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserTokenLogic {
	return &CheckUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckUserTokenLogic) CheckUserToken(in *user_mgr_pb.CheckUserTokenReq) (*user_mgr_pb.CheckUserTokenRsp, error) {
	if err := CheckUserId(in.UserId); err != nil {
		return nil, err
	}

	relation, err := CheckUserRegistered(l.ctx, l.svcCtx.TRelationModelMaster, in.UserId)
	if err != nil {
		return nil, err
	}

	isValid := CheckUserToken(in.UserToken, in.UserId, in.BusinessInfo)
	validStatus := ValidStatusInvalid
	if isValid {
		validStatus = ValidStatusValid
	}

	return &user_mgr_pb.CheckUserTokenRsp{
		UserId:      relation.UserId,
		ValidStatus: int32(validStatus),
	}, nil
}
