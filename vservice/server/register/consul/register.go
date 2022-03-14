package consul

import (
	"context"
	"github.com/hashicorp/consul/api"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vtool/vservice/common"
)

// consul agent -vservice -ui -bootstrap-expect=1 -client=0.0.0.0 -bind {ip addr} -data-dir={data dir} >> {log dir}
type Register struct {
	client *api.Client
	check  api.AgentServiceCheck
}

func (c Register) Register(ctx context.Context, path, servAddr string, ttl time.Duration) (string, error) {
	registration := new(api.AgentServiceRegistration)
	parts := strings.Split(servAddr, common.Colon)
	// health check caller
	check := new(api.AgentServiceCheck)
	check.HTTP = common.HttpPrefix + servAddr + DefaultCheckPath
	check.Method = http.MethodGet
	check.Timeout = ttl.String()
	check.Interval = ttl.String()
	check.DeregisterCriticalServiceAfter = ttl.String()

	registration.Check = check
	port, _ := strconv.ParseInt(parts[1], 10, 64)
	registration.Port = int(port)
	registration.Name = common.ConsulNamespace + path
	registration.Address = parts[0]
	registration.ID = c.buildInsID(path)

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (c *Register) buildInsID(path string) string {
	return strings.ReplaceAll(path, common.Slash, common.Bar)
}

func (c *Register) UnRegister(ctx context.Context, path string) error {
	return c.client.Agent().ServiceDeregister(c.buildInsID(path))
}

func (c Register) Get(ctx context.Context, path string) (string, error) {
	return "", nil
}

// Refresh the expiration time of the node without updating the value
func (c Register) RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error {
	return nil
}

func (c Register) GetNode(ctx context.Context, path string) ([]*common.RegisterServiceInfo, error) {
	return nil, nil
}
