package logic

import (
	"errors"
	"math"
)

func CheckUserId(userId string) error {
	if len(userId) < 1 || len(userId) > 64 {
		return errors.New("userId length is not in range [1,64]")
	}
	return nil
}

func CheckAmount(amount int64) error {
	if amount < 1 || amount > math.MaxInt64-1 {
		return errors.New("amount is not in range [1,int64.Max-1]")
	}
	return nil
}

func CheckLimit(limit int32) (int32, error) {
	if limit < 0 || limit > 100 {
		return 0, errors.New("limit is not in range [0,100]")
	}
	if limit <= 0 {
		return 10, nil
	}
	return limit, nil
}

func CheckOffset(offset int32) int32 {
	if offset < 0 {
		return 0
	}
	return offset
}

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
