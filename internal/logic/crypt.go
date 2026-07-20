package logic

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func GenMD5(input string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(input)))
}

func GenUserToken(user_id, businessInfo string) string {
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	md5Str := GenMD5(businessInfo + user_id + timestampStr)
	return md5Str + timestampStr
}
