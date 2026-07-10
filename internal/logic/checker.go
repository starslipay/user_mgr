package logic

import (
	"errors"
	"math"
)

func CheckName(name string) error {
	if len(name) < 1 || len(name) > 64 {
		return errors.New("name length is not in range [1,64]")
	}
	return nil
}

func CheckAge(age int32) error {
	if age < 1 || age > 256 {
		return errors.New("age is not in range [1,256]")
	}
	return nil
}

func CheckGender(gender int32) error {
	if gender < 1 || gender > 2 {
		return errors.New("gender is not in range [1,2]")
	}
	return nil
}

func CheckAddress(address string) error {
	if len(address) < 1 || len(address) > 64 {
		return errors.New("address length is not in range [1,64]")
	}
	return nil
}

func CheckPhone(phone string) error {
	if len(phone) < 1 || len(phone) > 64 {
		return errors.New("phone length is not in range [1,64]")
	}
	return nil
}

func CheckEmail(email string) error {
	if len(email) < 1 || len(email) > 64 {
		return errors.New("email length is not in range [1,64]")
	}
	return nil
}

func CheckIdType(idType int32) error {
	if idType < 1 || idType > 64 {
		return errors.New("idType is not in range [1,64]")
	}
	return nil
}

func CheckIdCard(idCard string) error {
	if len(idCard) < 1 || len(idCard) > 64 {
		return errors.New("idCard length is not in range [1,64]")
	}
	return nil
}

func CheckUserId(userId string) error {
	if len(userId) < 1 || len(userId) > 64 {
		return errors.New("userId length is not in range [1,64]")
	}
	return nil
}

func CheckPassword(password string) error {
	if len(password) < 1 || len(password) > 64 {
		return errors.New("password length is not in range [1,64]")
	}
	return nil
}

func CheckAmount(amount int64) error {
	if amount < 1 || amount > math.MaxInt64-1 {
		return errors.New("amount is not in range [1,int64.Max-1]")
	}
	return nil
}
