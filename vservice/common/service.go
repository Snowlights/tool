package common

type ServiceInfo struct {
	Type string `json:"type"`
	Addr string `json:"addr"`
}

const (
	ServiceTypeGin = "gin"
)
