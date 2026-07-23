package logic

import (
	"math"

	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/user_mgr/internal/xerr"
	"google.golang.org/grpc/codes"
)

func CheckName(name string) error {
	if len(name) < 1 || len(name) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "name length is not in range [1,64]")
	}
	return nil
}

func CheckAge(age int32) error {
	if age < 1 || age > 256 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "age is not in range [1,256]")
	}
	return nil
}

func CheckGender(gender int32) error {
	if gender < 1 || gender > 2 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "gender is not in range [1,2]")
	}
	return nil
}

func CheckAddress(address string) error {
	if len(address) < 1 || len(address) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "address length is not in range [1,64]")
	}
	return nil
}

func CheckPhone(phone string) error {
	if len(phone) < 1 || len(phone) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "phone length is not in range [1,64]")
	}
	return nil
}

func CheckEmail(email string) error {
	if len(email) < 1 || len(email) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "email length is not in range [1,64]")
	}
	return nil
}

func CheckIdType(idType int32) error {
	if idType < 1 || idType > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "idType is not in range [1,64]")
	}
	return nil
}

func CheckIdCard(idCard string) error {
	if len(idCard) < 1 || len(idCard) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "idCard length is not in range [1,64]")
	}
	return nil
}

func CheckUserId(userId string) error {
	if len(userId) < 1 || len(userId) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "userId length is not in range [1,64]")
	}
	return nil
}

func CheckPassword(password string) error {
	if len(password) < 1 || len(password) > 64 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "password length is not in range [1,64]")
	}
	return nil
}

func CheckAmount(amount int64) error {
	if amount < 1 || amount > math.MaxInt64-1 {
		return xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "amount is not in range [1,int64.Max-1]")
	}
	return nil
}
