package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

// DBConfig 数据库连接配置
type DBConfig struct {
	DataSource         string // 数据库连接字符串
	MaxOpenConns       int    // 最大打开连接数
	MaxIdleConns       int    // 最大空闲连接数
	ConnMaxLifetimeSec int    // 连接最大生命周期秒数
}

type Config struct {
	zrpc.RpcServerConf
	MasterDBConfig      DBConfig
	SlaveDBConfig       DBConfig
	AccountMgrRpcConfig zrpc.RpcClientConf
	TradeIdMgrRpcConfig zrpc.RpcClientConf

	// Kafka 生产者: relation 缓存更新消息, 由 cache_consumer 消费后回调 cache_mgr 写 redis
	KafkaProducerConf   []string
	TopicUpdateRelation string

	// AccessLog 请求/响应日志脱敏配置
	AccessLog AccessLogConf
}

// AccessLogConf 访问日志配置
type AccessLogConf struct {
	Enable          bool     `json:",default=true"`
	SensitiveFields []string `json:",optional"`
}
