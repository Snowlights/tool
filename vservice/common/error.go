package common

import "errors"

var (
	UnSupportedRegistrationType = errors.New("不支持的注册类型")
)

var (
	ServiceNameIsNil  = errors.New("service name can not be nil")
	ServiceGroupIsNil = errors.New("service group can not be nil")
	RegisterTypeIsNil = errors.New("register type can not be nil")
	LogDirIsNil       = errors.New("logDir can not be nil")
)
