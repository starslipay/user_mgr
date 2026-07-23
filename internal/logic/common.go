package logic

import (
	"context"

	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/model/mysql"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	RelationStateRegistering = 1 // 注册中
	RelationStateRegistered  = 2 // 注册成功
)

const (
	FlowTypeIn  = 1 // 入
	FlowTypeOut = 2 // 出
)

const (
	BizTypeDeposit  = 1 // 充值
	BizTypeWithdraw = 2 // 提现
)

func CheckUserRegistered(ctx context.Context, model mysql.TRelationModel, userId string) (*mysql.TRelation, error) {
	relation, err := model.FindOne(ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeUserNotExist, "user not exist")
		}
		return nil, err
	}
	if relation.State != RelationStateRegistered {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeUserNotExist, "user not exist")
	}
	return relation, nil
}
