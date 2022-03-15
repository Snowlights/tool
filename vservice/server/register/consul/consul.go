package consul

import (
	"github.com/hashicorp/consul/api"
	"net"
)

// must init in serv base func
func NewRegistry(regConfig *RegisterConfig) (*Register, error) {
	config := api.DefaultConfig()
	config.Address = net.JoinHostPort(regConfig.Host, regConfig.Port)
	config.Token = regConfig.Token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	register := &Register{client: client}
	return register, nil
}
