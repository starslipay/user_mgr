package mysql

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TUserInfoModel = (*customTUserInfoModel)(nil)

type (
	// TUserInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTUserInfoModel.
	TUserInfoModel interface {
		tUserInfoModel
		WithSession(session sqlx.Session) TUserInfoModel
		TransactCtx(ctx context.Context, fn func(ctx context.Context, tx TUserInfoModel) error) error
	}

	customTUserInfoModel struct {
		*defaultTUserInfoModel
	}
)

// NewTUserInfoModel returns a model for the database table.
func NewTUserInfoModel(conn sqlx.SqlConn) TUserInfoModel {
	return &customTUserInfoModel{
		defaultTUserInfoModel: newTUserInfoModel(conn),
	}
}

func (m *customTUserInfoModel) WithSession(session sqlx.Session) TUserInfoModel {
	return NewTUserInfoModel(sqlx.NewSqlConnFromSession(session))
}

// TransactCtx 事务封装
func (m *customTUserInfoModel) TransactCtx(ctx context.Context, fn func(ctx context.Context, tx TUserInfoModel) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		return fn(ctx, m.WithSession(tx))
	})
}
