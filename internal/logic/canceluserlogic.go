package logic

import (
	"context"

	"github.com/starslipay/account_mgr/account_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/starslipay/user_mgr/user_mgr_pb"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CancelUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelUserLogic {
	return &CancelUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelUserLogic) CancelUser(in *user_mgr_pb.CancelUserReq) (*user_mgr_pb.CancelUserRsp, error) {
	if err := CheckUserId(in.UserId); err != nil {
		return nil, err
	}

	relation, err := l.svcCtx.TRelationModelMaster.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeUserNotExist, "user not exist")
		}
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "find relation failed: "+err.Error())
	}

	switch relation.State {
	case RelationStateCancelled:
		// 已注销，幂等返回成功
		return &user_mgr_pb.CancelUserRsp{UserId: in.UserId}, nil
	case RelationStateRegistered:
		// 正常注销流程
	default:
		// 注册中/其他状态不允许注销
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeRelationStateInvalid, "relation state invalid, cannot cancel")
	}

	// 余额不为零不允许注销，提示用户先转出。查主库(qry_mode=2)确保数据准确
	balanceRsp, err := l.svcCtx.AccountMgrRpcClient.GetUserBalanceInfo(l.ctx, &account_mgr_pb.GetUserBalanceInfoReq{
		Uid:     relation.Uid,
		QryMode: 2, // 2-查主库
	})
	if err != nil {
		return nil, xerror.HandleRPCError(err, "account_mgr.GetUserBalanceInfo")
	}
	if balanceRsp.Balance != 0 {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeBalanceNotZero, "balance not zero, please transfer out before cancel")
	}

	// 软删除：relation 状态置为已注销，userinfo 数据保留
	err = l.svcCtx.TRelationModelMaster.Update(l.ctx, &mysql.TRelation{
		UserId: in.UserId,
		Uid:    relation.Uid,
		State:  RelationStateCancelled,
	})
	if err != nil {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "update relation state failed: "+err.Error())
	}

	// 注销成功，异步发送 relation 缓存删除消息(t_update_relation)，由 cache_consumer 消费后回调 cache_mgr 删除 redis 缓存
	sendDeleteRelationMessage(l.ctx, l.svcCtx, in.UserId)

	return &user_mgr_pb.CancelUserRsp{UserId: in.UserId}, nil
}
