package logic

import (
	"context"

	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/starslipay/user_mgr/user_mgr_pb"
	"google.golang.org/grpc/codes"

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
	if err := CheckUserId(in.UserId); err != nil {
		return nil, err
	}

	relation, err := CheckUserRegistered(l.ctx, l.svcCtx.TRelationModelMaster, in.UserId)
	if err != nil {
		return nil, err
	}

	// 先查询用户信息，在事务中更新
	err = l.svcCtx.TUserInfoModelMaster.TransactCtx(l.ctx, func(ctx context.Context, tx mysql.TUserInfoModel) error {
		userInfo, err := tx.FindOne(ctx, relation.Uid)
		if err != nil {
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeUnKnownDBError, "find user info failed: "+err.Error())
		}

		// 传入的字段，才更新
		if in.Name != "" {
			userInfo.Name = in.Name
		}
		if in.Age != 0 {
			userInfo.Age = int64(in.Age)
		}
		if in.Gender != 0 {
			userInfo.Gender = int64(in.Gender)
		}
		if in.Address != "" {
			userInfo.Address = in.Address
		}
		if in.Phone != "" {
			userInfo.Phone = in.Phone
		}
		if in.Email != "" {
			userInfo.Email = in.Email
		}
		if in.IdType != 0 {
			userInfo.IdType = int64(in.IdType)
		}
		if in.IdCard != "" {
			userInfo.IdCard = in.IdCard
		}

		err = tx.Update(ctx, userInfo)
		if err != nil {
			return xerror.NewBizError(codes.Internal, xerr.ErrCodeUnKnownDBError, "update user info failed: "+err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &user_mgr_pb.UpdateUserInfoRsp{
		UserId: relation.UserId,
	}, nil
}
