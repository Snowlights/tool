package vnet

import "errors"

const (
	networkTypeTCP = "tcp"
)

var NoInternalIp = errors.New("no internal ip")
