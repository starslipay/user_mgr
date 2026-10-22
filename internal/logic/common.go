package logic

import (
	"context"
	"encoding/json"

	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/model/mysql"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	RelationStateRegistering = 1 // 注册中
	RelationStateRegistered  = 2 // 注册成功
	RelationStateCancelled   = 3 // 已注销(禁止重新注册)
)

// relation 缓存消息操作类型
const (
	RelationCacheOpUpdate = "update" // 写入/更新 relation|{user_id}
	RelationCacheOpDelete = "delete" // 删除 relation|{user_id}
)

// UpdateRelationMessage relation 缓存更新消息, 投递到 Kafka topic t_update_relation
// 消费端 cache_consumer 收到后回调 cache_mgr 执行实际的 redis 写入/删除
type UpdateRelationMessage struct {
	Op     string `json:"op"` // update-写入, delete-删除; 空兼容旧消息, 按 update 处理
	UserId string `json:"user_id"`
	Uid    int64  `json:"uid"`
}

// sendUpdateRelationMessage 发送 relation 缓存更新消息(t_update_relation), 替代直接写 redis, 解耦主流程与缓存写入
// 发送失败仅记录日志, 不影响注册主流程(缓存最终一致性, 由对账任务兜底)
func sendUpdateRelationMessage(ctx context.Context, svcCtx *svc.ServiceContext, userId string, uid int64) {
	// 总开关关闭时不发送消息, relation 缓存不做更新
	if !svcCtx.Config.EnableUpdateRelationCmq {
		logx.WithContext(ctx).Infof("update relation cmq disabled, skip, user_id: %s, uid: %d", userId, uid)
		return
	}

	msg := &UpdateRelationMessage{
		Op:     RelationCacheOpUpdate,
		UserId: userId,
		Uid:    uid,
	}
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		logx.WithContext(ctx).Errorf("marshal update relation message failed, user_id: %s, err: %v", userId, err)
		return
	}
	if err := svcCtx.UpdateRelationProducer.Push(ctx, string(jsonStr)); err != nil {
		logx.WithContext(ctx).Errorf("send update relation message failed, user_id: %s, uid: %d, err: %v", userId, uid, err)
	}
}

// sendDeleteRelationMessage 发送 relation 缓存删除消息(t_update_relation), 注销用户后删除 redis 中的 relation 缓存
// 发送失败仅记录日志, 不影响注销主流程
func sendDeleteRelationMessage(ctx context.Context, svcCtx *svc.ServiceContext, userId string) {
	// 与更新消息共用同一个总开关
	if !svcCtx.Config.EnableUpdateRelationCmq {
		logx.WithContext(ctx).Infof("update relation cmq disabled, skip delete, user_id: %s", userId)
		return
	}

	msg := &UpdateRelationMessage{
		Op:     RelationCacheOpDelete,
		UserId: userId,
	}
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		logx.WithContext(ctx).Errorf("marshal delete relation message failed, user_id: %s, err: %v", userId, err)
		return
	}
	if err := svcCtx.UpdateRelationProducer.Push(ctx, string(jsonStr)); err != nil {
		logx.WithContext(ctx).Errorf("send delete relation message failed, user_id: %s, err: %v", userId, err)
	}
}

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
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeDBError, "find user relation failed: "+err.Error())
	}
	if relation.State != RelationStateRegistered {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeUserNotExist, "user not exist")
	}
	return relation, nil
}
