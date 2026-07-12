package logic

import (
	"context"

	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/user_mgr_pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationLogic {
	return &GetRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationLogic) GetRelation(in *user_mgr_pb.GetRelationReq) (*user_mgr_pb.GetRelationRsp, error) {
	relation, err := l.svcCtx.TRelationModelSlave.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, xerr.NewDBError("find relation failed: " + err.Error())
	}

	if relation.State != RelationStateRegistered {
		return nil, xerr.ErrUserNotExist
	}
	return &user_mgr_pb.GetRelationRsp{
		UserId: in.UserId,
		Uid:    relation.Uid,
	}, nil
}
