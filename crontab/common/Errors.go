package common

import "errors"

var (
	ERROR_LOCK_ALREADY_REQUIRED = errors.New("锁已被占用")
	ERROR_NO_LOCAL_IP_FOUND     = errors.New("没有找到网卡IP")
)
