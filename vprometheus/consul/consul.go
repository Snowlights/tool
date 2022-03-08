package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
	"strings"
)

// consul agent -server -ui -bootstrap-expect=1 -client=0.0.0.0 -bind 你的ip地址 -data-dir=/状态数据存储文件夹/data >> /日志记录文件夹/logs/consul.log
type consulServiceRegistry struct {
	client *api.Client
}

func NewConsulServiceRegistry(host string, port int, token string) (*consulServiceRegistry, error) {

	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = "1f8afae5-32e7-c38f-eaec-497dd0532b88"
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulServiceRegistry{client: client}, nil
}

func (c consulServiceRegistry) Register(path, servAddr string) bool {
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)

	parts := strings.Split(servAddr, ":")

	// 增加consul健康检查回调函数
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("%s://%s/health", "http", servAddr)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "10s"
	registration.Check = check
	port, _ := strconv.ParseInt(parts[1], 10, 64)
	registration.Port = int(port)
	registration.Name = path
	registration.Address = parts[0]
	// 注册服务到consul
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
