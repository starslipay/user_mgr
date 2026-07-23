package logic

import (
	"context"

	"github.com/starslipay/account_mgr/account_mgr_pb"
	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/trade_id_mgr/trade_id_mgr_pb"
	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/starslipay/user_mgr/user_mgr_pb"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	if err := CheckUserId(in.UserId); err != nil {
		return nil, err
	}
	if err := CheckName(in.Name); err != nil {
		return nil, err
	}
	if err := CheckAge(in.Age); err != nil {
		return nil, err
	}
	if err := CheckGender(in.Gender); err != nil {
		return nil, err
	}
	if err := CheckAddress(in.Address); err != nil {
		return nil, err
	}
	if err := CheckPhone(in.Phone); err != nil {
		return nil, err
	}
	if err := CheckEmail(in.Email); err != nil {
		return nil, err
	}
	if err := CheckIdType(in.IdType); err != nil {
		return nil, err
	}
	if err := CheckIdCard(in.IdCard); err != nil {
		return nil, err
	}
	if err := CheckPassword(in.Password); err != nil {
		return nil, err
	}
	// TODO password这种参数需要加密传输
	PasswordMD5 := GenMD5(in.Password)

	isExistRelation := true
	// 先查询relation是否已经存在
	relation, err := l.svcCtx.TRelationModelMaster.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			isExistRelation = false
		} else {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "find relation failed: "+err.Error())
		}
	} else {
		if RelationStateRegistering == relation.State {
			// 继续关联中，不执行后续操作
		} else if RelationStateRegistered == relation.State {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeUserAlreadyRegistered, "user already registered")
		} else {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeRelationStateNotRegisteringOrRegistered, "relation state not registering or registered")
		}
	}

	var uid int64
	if isExistRelation {
		uid = relation.Uid
	} else {
		uid, err = l.genUid()
		if err != nil {
			return nil, err
		}

		_, err = l.svcCtx.TRelationModelMaster.Insert(l.ctx, &mysql.TRelation{
			UserId: in.UserId,
			Uid:    uid,
			State:  RelationStateRegistering, // 注册中
		})
		if err != nil {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "insert relation failed: "+err.Error())
		}
	}

	isExistUserInfo := true
	userInfo, err := l.svcCtx.TUserInfoModelMaster.FindOne(l.ctx, uid)
	if err != nil {
		if err == sqlx.ErrNotFound {
			isExistUserInfo = false
		} else {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "find user info failed: "+err.Error())
		}
	}
	if !isExistUserInfo {
		// 插入用户信息
		_, err = l.svcCtx.TUserInfoModelMaster.Insert(l.ctx, &mysql.TUserInfo{
			Uid:      uid,
			UserId:   in.UserId,
			Password: PasswordMD5,
			Name:     in.Name,
			Gender:   int64(in.Gender),
			Age:      int64(in.Age),
			Address:  in.Address,
			Phone:    in.Phone,
			Email:    in.Email,
			IdType:   int64(in.IdType),
			IdCard:   in.IdCard,
		})
		if err != nil {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "insert user info failed: "+err.Error())
		}
	} else {
		userInfo.UserId = in.UserId
		userInfo.Password = PasswordMD5
		userInfo.Name = in.Name
		userInfo.Gender = int64(in.Gender)
		userInfo.Age = int64(in.Age)
		userInfo.Address = in.Address
		userInfo.Phone = in.Phone
		userInfo.Email = in.Email
		userInfo.IdType = int64(in.IdType)
		userInfo.IdCard = in.IdCard
		// 更新用户信息
		err = l.svcCtx.TUserInfoModelMaster.Update(l.ctx, userInfo)
		if err != nil {
			return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "update user info failed: "+err.Error())
		}
	}

	// 创建资金账户
	createAccountRsp, err := l.svcCtx.AccountMgrRpcClient.CreateAccount(l.ctx, &account_mgr_pb.CreateAccountReq{
		Uid:     uid,
		UserId:  in.UserId,
		CurType: 1, // 1-人民币
	})
	if err != nil {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "create account failed: "+err.Error())
	}
	if createAccountRsp.IsRepeat {
		l.Logger.Info("create account already exist, create repeat")
	}

	err = l.svcCtx.TRelationModelMaster.Update(l.ctx, &mysql.TRelation{
		UserId: in.UserId,
		Uid:    uid,
		State:  RelationStateRegistered, // 注册成功
	})
	if err != nil {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "update relation state failed: "+err.Error())
	}

	return &user_mgr_pb.RegUserRsp{
		UserId: in.UserId,
	}, nil
}

func (l *RegUserLogic) genUid() (int64, error) {
	genUidRsp, err := l.svcCtx.TradeIdMgrRpcClient.GenUid(l.ctx, &trade_id_mgr_pb.GenUidReq{})
	if err != nil {
		return 0, err
	}
	return genUidRsp.Uid, nil
}
