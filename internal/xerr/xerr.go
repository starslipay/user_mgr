package xerr

import (
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"
)

// 错误码  10000 0000 ~~99999 9999
// 模块id  20000
// 错误码 = 模块id + 业务错误码
var (
	ModuleId = int64(455904)
)

var (
	// 系统错误 000-099
	ErrCodeDBError        = int64(455904000)
	ErrCodeServerInternal = int64(455904001)
	ErrCodeCallRpc        = int64(455904002)

	// 业务错误码 100-999
	ErrCodeParam                 = int64(455904100)
	ErrCodeUserNotExist          = int64(455904101)
	ErrCodePasswordWrong         = int64(455904102)
	ErrCodeUserAlreadyRegistered = int64(455904103)
	ErrCodeRelationStateInvalid  = int64(455904104)
	ErrCodeTokenInvalid          = int64(455904105)
)

func ParseRPCError(err error) error {
	// 解析下游业务错误
	bizError, isSuccessParse := xerror.ParseBizError(err)
	if isSuccessParse {
		return xerror.NewBizError(codes.Internal, bizError.Code, bizError.Message)
	}

	// 如果没有解析到业务错误，返回rpc错误码
	return xerror.NewBizError(codes.Internal, ErrCodeCallRpc, err.Error())
}
