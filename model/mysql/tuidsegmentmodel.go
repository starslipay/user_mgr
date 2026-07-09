package mysql

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TUidSegmentModel = (*customTUidSegmentModel)(nil)

type (
	// TUidSegmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTUidSegmentModel.
	TUidSegmentModel interface {
		tUidSegmentModel
		WithSession(session sqlx.Session) TUidSegmentModel
		TransactCtx(ctx context.Context, fn func(ctx context.Context, tx TUidSegmentModel) error) error
		FindOneForUpdate(ctx context.Context, id int64) (*TUidSegment, error)
	}

	customTUidSegmentModel struct {
		*defaultTUidSegmentModel
	}
)

// NewTUidSegmentModel returns a model for the database table.
func NewTUidSegmentModel(conn sqlx.SqlConn) TUidSegmentModel {
	return &customTUidSegmentModel{
		defaultTUidSegmentModel: newTUidSegmentModel(conn),
	}
}

func (m *customTUidSegmentModel) WithSession(session sqlx.Session) TUidSegmentModel {
	return NewTUidSegmentModel(sqlx.NewSqlConnFromSession(session))
}

// TransactCtx 事务封装，回调中直接使用 TUidSegmentModel 方法
func (m *customTUidSegmentModel) TransactCtx(ctx context.Context, fn func(ctx context.Context, tx TUidSegmentModel) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		return fn(ctx, m.WithSession(tx))
	})
}

// FindOneForUpdate 查询并锁定记录
func (m *customTUidSegmentModel) FindOneForUpdate(ctx context.Context, id int64) (*TUidSegment, error) {
	query := "select " + tUidSegmentRows + " from " + m.table + " where `id` = ? limit 1 for update"
	var resp TUidSegment
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
