package logic

import (
	"context"

	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/user_mgr_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTokenLogic {
	return &GetUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserTokenLogic) GetUserToken(in *user_mgr_pb.GetUserTokenReq) (*user_mgr_pb.GetUserTokenRsp, error) {
	if err := CheckUserId(in.UserId); err != nil {
		return nil, err
	}

	relation, err := CheckUserRegistered(l.ctx, l.svcCtx.TRelationModelMaster, in.UserId)
	if err != nil {
		return nil, err
	}

	userInfo, err := l.svcCtx.TUserInfoModelMaster.FindOne(l.ctx, relation.Uid)
	if err != nil {
		return nil, xerr.NewDBError("find user info failed: " + err.Error())
	}

	inPasswordMD5 := GenMD5(in.Password)
	if userInfo.Password != inPasswordMD5 {
		return nil, xerr.ErrPasswordWrong
	}

	userToken := GenUserToken(in.UserId, in.BusinessInfo)

	return &user_mgr_pb.GetUserTokenRsp{
		UserToken: userToken,
	}, nil
}
