package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TRelationModel = (*customTRelationModel)(nil)

type (
	// TRelationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTRelationModel.
	TRelationModel interface {
		tRelationModel
		withSession(session sqlx.Session) TRelationModel
	}

	customTRelationModel struct {
		*defaultTRelationModel
	}
)

// NewTRelationModel returns a model for the database table.
func NewTRelationModel(conn sqlx.SqlConn) TRelationModel {
	return &customTRelationModel{
		defaultTRelationModel: newTRelationModel(conn),
	}
}

func (m *customTRelationModel) withSession(session sqlx.Session) TRelationModel {
	return NewTRelationModel(sqlx.NewSqlConnFromSession(session))
}
