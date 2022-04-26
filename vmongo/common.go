package vmongo

import "strings"

const (
	comma = ","
)

type MongoConfig struct {
	Conf map[string]InstanceConfig `json:"conf" properties:"conf"`
}

type InstanceConfig struct {
	Host     string
	Username string
	Password string
	Document string

	Timeout      int64
	ReadTimeout  int64
	WriteTimeout int64
	PoolSize     int64
}

func (ic InstanceConfig) buildInsCfgKey() string {
	return strings.Join([]string{ic.Host, ic.Username, ic.Password}, comma)
}
