package logic

import (
	"context"

	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/user_mgr_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	CheckResultValid   = 1
	CheckResultInvalid = 2
)

type CheckPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPasswordLogic {
	return &CheckPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckPasswordLogic) CheckPassword(in *user_mgr_pb.CheckPasswordReq) (*user_mgr_pb.CheckPasswordRsp, error) {
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

	return &user_mgr_pb.CheckPasswordRsp{}, nil
}
