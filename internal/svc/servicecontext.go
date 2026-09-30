package svc

import (
	"database/sql"
	"time"

	"github.com/starslipay/account_mgr/account_mgr_pb"
	"github.com/starslipay/trade_id_mgr/trade_id_mgr_pb"
	"github.com/starslipay/user_mgr/internal/config"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	TRelationModelMaster mysql.TRelationModel
	TUserInfoModelMaster mysql.TUserInfoModel

	TRelationModelSlave mysql.TRelationModel
	TUserInfoModelSlave mysql.TUserInfoModel

	AccountMgrRpcClient account_mgr_pb.AccountMgrClient
	TradeIdMgrRpcClient trade_id_mgr_pb.TradeIdMgrClient
}

// newDBConn 自建 *sql.DB 以支持自定义连接池参数
func newDBConn(dataSource string, maxOpen, maxIdle, lifetimeSec int) sqlx.SqlConn {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		logx.Must(err)
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(lifetimeSec) * time.Second)
	if err = db.Ping(); err != nil {
		logx.Must(err)
	}
	// 关闭db熔断保护
	return sqlx.NewSqlConnFromDB(db, sqlx.WithAcceptable(func(err error) bool { return true }))
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlMasterConn := newDBConn(c.MasterDBConfig.DataSource,
		c.MasterDBConfig.MaxOpenConns, c.MasterDBConfig.MaxIdleConns, c.MasterDBConfig.ConnMaxLifetimeSec)
	SqlSlaveConn := newDBConn(c.SlaveDBConfig.DataSource,
		c.SlaveDBConfig.MaxOpenConns, c.SlaveDBConfig.MaxIdleConns, c.SlaveDBConfig.ConnMaxLifetimeSec)
	return &ServiceContext{
		Config:               c,
		TRelationModelMaster: mysql.NewTRelationModel(SqlMasterConn),
		TUserInfoModelMaster: mysql.NewTUserInfoModel(SqlMasterConn),

		TRelationModelSlave: mysql.NewTRelationModel(SqlSlaveConn),
		TUserInfoModelSlave: mysql.NewTUserInfoModel(SqlSlaveConn),

		AccountMgrRpcClient: account_mgr_pb.NewAccountMgrClient(zrpc.MustNewClient(c.AccountMgrRpcConfig).Conn()),
		TradeIdMgrRpcClient: trade_id_mgr_pb.NewTradeIdMgrClient(zrpc.MustNewClient(c.TradeIdMgrRpcConfig).Conn()),
	}
}
