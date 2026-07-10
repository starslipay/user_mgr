package logic

import (
	"context"
	"errors"

	"github.com/starslipay/user_mgr/internal/svc"
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
		return nil, err
	}
	if relation.State != RelationStateRegistered {
		return nil, errors.New("user not registered")
	}
	return &user_mgr_pb.GetRelationRsp{
		UserId: in.UserId,
		Uid:    relation.Uid,
	}, nil
}
