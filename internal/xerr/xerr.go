package xerr

import "fmt"

type CodeMsg struct {
	Code int
	Msg  string
}

func (c *CodeMsg) Error() string {
	return fmt.Sprintf("[%d]%s", c.Code, c.Msg)
}

// 错误码  10000 0000 ~~99999 9999
// 模块id  10000
// 错误码 = 模块id + 业务错误码
var (
	ModuleId        = int64(20000)
	ModuleErrorBase = ModuleId * 10000
)

var (
	// 系统错误 0000-0999
	ErrCodeServerInternal = ModuleErrorBase + 0
	ErrCodeDBError        = ModuleErrorBase + 1

	// 业务错误码 1000-1999
	ErrCodeParam                                   = ModuleErrorBase + 1000
	ErrCodeUserNotExist                            = ModuleErrorBase + 1001
	ErrCodePasswordWrong                           = ModuleErrorBase + 1002
	ErrCodeUserAlreadyRegistered                   = ModuleErrorBase + 1003
	ErrCodeRelationStateNotRegisteringOrRegistered = ModuleErrorBase + 1004
	ErrCodeTokenInvalid                            = ModuleErrorBase + 1005
)
