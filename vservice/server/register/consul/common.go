package consul

import (
	"context"
	"vtool/vlog"
)

type RegisterConfig struct {
	Host  string
	Port  string
	Token string
}

const (
	ConsulNamespace = "consul"

	defaultHost  = "127.0.0.1"
	defaultPort  = "8500"
	defaultToken = "1f8afae5-32e7-c38f-eaec-497dd0532b88"
)

var DefaultConsulInstance *Register

// must init, todo change to apollo config
func init() {
	ins, err := NewRegistry(&RegisterConfig{
		Host:  defaultHost,
		Port:  defaultPort,
		Token: defaultToken,
	})
	if err != nil {
		vlog.Error(context.Background(), err)
	}
	DefaultConsulInstance = ins
}
