package xerr

import (
	"github.com/starslipay/paycomm/xerror"
	"google.golang.org/grpc/codes"
)

// 错误码  10000 0000 ~~99999 9999
// 模块id  20000
// 错误码 = 模块id + 业务错误码
var (
	ModuleId        = int64(20000)
	ModuleErrorBase = ModuleId * 10000
)

var (
	// 系统错误 0000-0999
	ErrCodeDBError        = ModuleErrorBase + 0
	ErrCodeServerInternal = ModuleErrorBase + 1
	ErrCodeCallRpc        = ModuleErrorBase + 2

	// 业务错误码 1000-1999
	ErrCodeParam                                   = ModuleErrorBase + 1000
	ErrCodeUserNotExist                            = ModuleErrorBase + 1001
	ErrCodePasswordWrong                           = ModuleErrorBase + 1002
	ErrCodeUserAlreadyRegistered                   = ModuleErrorBase + 1003
	ErrCodeRelationStateNotRegisteringOrRegistered = ModuleErrorBase + 1004
	ErrCodeTokenInvalid                            = ModuleErrorBase + 1005
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
