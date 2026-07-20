package svc

import (
	"github.com/starslipay/account_mgr/account_mgr_pb"
	"github.com/starslipay/trade_id_mgr/trade_id_mgr_pb"
	"github.com/starslipay/user_mgr/internal/config"
	"github.com/starslipay/user_mgr/model/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	TRelationModelMaster   mysql.TRelationModel
	TUserInfoModelMaster   mysql.TUserInfoModel
	TUidSegmentModelMaster mysql.TUidSegmentModel

	TRelationModelSlave mysql.TRelationModel
	TUserInfoModelSlave mysql.TUserInfoModel

	AccountMgrRpcClient account_mgr_pb.AccountMgrClient
	TradeIdMgrRpcClient trade_id_mgr_pb.TradeIdMgrClient
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

		AccountMgrRpcClient: account_mgr_pb.NewAccountMgrClient(zrpc.MustNewClient(c.AccountMgrRpcConfig).Conn()),
		TradeIdMgrRpcClient: trade_id_mgr_pb.NewTradeIdMgrClient(zrpc.MustNewClient(c.TradeIdMgrRpcConfig).Conn()),
	}
}
