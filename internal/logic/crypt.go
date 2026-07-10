package logic

import (
	"crypto/md5"
	"fmt"
)

func GenMD5(input string) string {
	// 计算input的md5值
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
