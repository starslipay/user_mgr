package svc

import (
	"github.com/starslipay/user_mgr/internal/config"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                 config.Config
	TRelationModelMaster   mysql.TRelationModel
	TUserInfoModelMaster   mysql.TUserInfoModel
	TUidSegmentModelMaster mysql.TUidSegmentModel

	TRelationModelSlave mysql.TRelationModel
	TUserInfoModelSlave mysql.TUserInfoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlMasterConn := sqlx.NewMysql(c.MasterDBConfig.DataSource)
	SqlSlaveConn := sqlx.NewMysql(c.SlaveDBConfig.DataSource)
	return &ServiceContext{
		Config:                 c,
		TRelationModelMaster:   mysql.NewTRelationModel(SqlMasterConn),
		TUserInfoModelMaster:   mysql.NewTUserInfoModel(SqlMasterConn),
		TUidSegmentModelMaster: mysql.NewTUidSegmentModel(SqlMasterConn),

		TRelationModelSlave: mysql.NewTRelationModel(SqlSlaveConn),
		TUserInfoModelSlave: mysql.NewTUserInfoModel(SqlSlaveConn),
	}
}
