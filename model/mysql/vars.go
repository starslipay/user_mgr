package mysql

import (
	"database/sql"
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

var ErrRowsAffectedNotOne = errors.New("affected rows is not 1")

func checkOneRowAffected(ret sql.Result) error {
	rows, err := ret.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrRowsAffectedNotOne
	}
	return nil
}
