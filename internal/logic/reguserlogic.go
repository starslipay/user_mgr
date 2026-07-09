package logic

import (
	"context"
	"errors"

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

	// 先查询relation是否已经存在
	relation, err := l.svcCtx.TRelationModelMaster.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err != sqlx.ErrNotFound {
			return nil, err
		}
	} else {
		if relation.State == RelationStateRegistering {
			// 继续关联中，不执行后续操作
		} else if relation.State == RelationStateRegistered {
			// 重入关联成功的用户, 要校验关键信息一致性
			userInfoTmp, err := l.svcCtx.TUserInfoModelMaster.FindOne(l.ctx, relation.Uid)
			if err != nil {
				return nil, err
			}

			if userInfoTmp.Name != in.Name ||
				userInfoTmp.Password != in.Password ||
				userInfoTmp.Age != int64(in.Age) ||
				userInfoTmp.Gender != int64(in.Gender) ||
				userInfoTmp.Address != in.Address ||
				userInfoTmp.Phone != in.Phone ||
				userInfoTmp.Email != in.Email ||
				userInfoTmp.IdType != int64(in.IdType) ||
				userInfoTmp.IdCard != in.IdCard {
				return nil, errors.New("user info is not consistent")
			}

			return &user_mgr_pb.RegUserRsp{
				UserId:   in.UserId,
				IsRepeat: 1,
			}, nil
		} else {
			return nil, errors.New("relation state is not registering or registered")
		}
	}

	// 生成用户ID
	newUid, err := l.generateUid()

	if err != nil {
		l.Logger.Errorf("generateUid failed: %v", err)
		return nil, err
	}

	_, err = l.svcCtx.TRelationModelMaster.Insert(l.ctx, &mysql.TRelation{
		UserId: in.UserId,
		Uid:    newUid,
		State:  RelationStateRegistering, // 注册中
	})
	if err != nil {
		l.Logger.Errorf("insert relation failed: %v", err)
		return nil, err
	}

	// 插入用户信息
	_, err = l.svcCtx.TUserInfoModelMaster.Insert(l.ctx, &mysql.TUserInfo{
		Uid:      newUid,
		UserId:   in.UserId,
		Password: in.Password,
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

	// _, err = l.svcCtx.TAccountModelMaster.Insert(l.ctx, &mysql.TAccount{
	// 	Uid:     newUid,
	// 	UserId:  in.UserId,
	// 	Balance: 0,
	// })
	// if err != nil {
	// 	l.Logger.Errorf("insert account failed: %v", err)
	// 	return nil, err
	// }

	err = l.svcCtx.TRelationModelMaster.Update(l.ctx, &mysql.TRelation{
		UserId: in.UserId,
		Uid:    newUid,
		State:  RelationStateRegistered, // 注册成功
	})
	if err != nil {
		l.Logger.Errorf("update relation state failed: %v", err)
		return nil, err
	}

	return &user_mgr_pb.RegUserRsp{
		UserId:   in.UserId,
		IsRepeat: 0,
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
