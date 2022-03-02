package server

import "errors"

var (
	UnSupportedRegistrationType = errors.New("不支持的注册类型")
	RegisterFailed              = errors.New("注册服务失败")
)
