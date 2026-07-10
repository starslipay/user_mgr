package logic

import (
	"context"
	"errors"

	"github.com/starslipay/account_mgr/account_mgr_pb"
	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/starslipay/user_mgr/user_mgr_pb"

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
			return nil, err
		}
	} else {
		if RelationStateRegistering == relation.State {
			// 继续关联中，不执行后续操作
		} else if RelationStateRegistered == relation.State {
			return nil, errors.New("user already registered")
		} else {
			return nil, errors.New("relation state is not registering or registered")
		}
	}

	var uid int64
	if isExistRelation {
		uid = relation.Uid
	} else {
		uid, err = l.generateUid()
		if err != nil {
			l.Logger.Errorf("generateUid failed: %v", err)
			return nil, err
		}

		_, err = l.svcCtx.TRelationModelMaster.Insert(l.ctx, &mysql.TRelation{
			UserId: in.UserId,
			Uid:    uid,
			State:  RelationStateRegistering, // 注册中
		})
		if err != nil {
			l.Logger.Errorf("insert relation failed: %v", err)
			return nil, err
		}
	}

	isExistUserInfo := true
	userInfo, err := l.svcCtx.TUserInfoModelMaster.FindOne(l.ctx, uid)
	if err != nil {
		if err == sqlx.ErrNotFound {
			isExistUserInfo = false
		} else {
			return nil, err
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
			l.Logger.Errorf("insert user info failed: %v", err)
			return nil, err
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
			l.Logger.Errorf("update user info failed: %v", err)
			return nil, err
		}
	}

	// 创建资金账户
	createAccountRsp, err := l.svcCtx.AccountMgrRpcClient.CreateAccount(l.ctx, &account_mgr_pb.CreateAccountReq{
		Uid:     uid,
		UserId:  in.UserId,
		CurType: 1, // 1-人民币
	})
	if err != nil {
		l.Logger.Errorf("create account failed: %v", err)
		return nil, err
	}
	if createAccountRsp.IsRepeat {
		l.Logger.Info("create account repeat")
	}

	err = l.svcCtx.TRelationModelMaster.Update(l.ctx, &mysql.TRelation{
		UserId: in.UserId,
		Uid:    uid,
		State:  RelationStateRegistered, // 注册成功
	})
	if err != nil {
		l.Logger.Errorf("update relation state failed: %v", err)
		return nil, err
	}

	return &user_mgr_pb.RegUserRsp{
		UserId: in.UserId,
	}, nil
}

// generateUid 从 t_uid_segment 获取一个UID，需要在事务中完成
func (l *RegUserLogic) generateUid() (int64, error) {
	var newUid int64
	err := l.svcCtx.TUidSegmentModelMaster.TransactCtx(l.ctx, func(ctx context.Context, tx mysql.TUidSegmentModel) error {
		// 查询当前 segment (使用 FOR UPDATE 锁定行)
		segment, err := tx.FindOneForUpdate(ctx, 1)
		if err != nil {
			return err
		}

		// 计算新的 UID
		newUid = segment.UidMax + 1

		// 更新 uid_max
		segment.UidMax = newUid
		return tx.Update(ctx, segment)
	})

	if err != nil {
		return 0, err
	}

	return newUid, nil
}
