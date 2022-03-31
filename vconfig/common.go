package vconfig

import (
	"errors"
	"github.com/apolloconfig/agollo/v4/storage"
	"os"
	"strings"
)

const (
	// server app config
	Application = "application"
	// server config
	Server = "server"
	// client config
	Client = "client"

	// middleware config
	MiddlewareAppID          = "middleware"
	MiddlewareNamespaceTrace = "trace"
	MiddlewareNamespaceMQ    = "mq"
	MiddlewareNamespaceDB    = "db"

	BackupPath = "backupPath"
	httpScheme = "http://"

	dot   = "."
	comma = ","
	colon = ":"
	slash = "/"
	under = "_"
	minus = "-"

	portMin = 0
	portMax = 65535
)

var (
	InvalidIp      = errors.New("invalid ip")
	InvalidPort    = errors.New("invalid port")
	InvalidBackup  = errors.New("invalid backup")
	InvalidCluster = errors.New("apollo cluster is empty")
)

type Center interface {
	GetValue(key string) (string, bool)
	GetValueWithNamespace(namespace, key string) (string, bool)
	UnmarshalWithNameSpace(namespace, tag string, v interface{}) error
	AddListener(listeners storage.ChangeListener)
}

type CenterConfig struct {
	AppID             string
	Cluster           string
	Namespace         []string
	IP                string
	Port              int
	IsBackupConfig    bool
	BackupConfigPath  string
	SecretKey         string
	SyncServerTimeout int
	MustStart         bool
}

type ServerConfig struct {
	RegisterType    int64    `json:"register_type" properties:"register_type"`
	RegisterCluster []string `json:"register_cluster" properties:"register_cluster"`

	NeedMetric  bool   `json:"need_metric" properties:"need_metric"`
	ConsulHost  string `json:"consul_host" properties:"consul_host"`
	ConsulPort  string `json:"consul_port" properties:"consul_port"`
	ConsulToken string `json:"consul_token" properties:"consul_token"`

	LogLevel string `json:"log_level" properties:"log_level"`

	// todo db、mq、redis config
}

type ClientConfig struct {
	Idle        int64 `json:"idle" properties:"idle"`
	IdleTimeout int64 `json:"idle_timeout" properties:"idle_timeout"`
	MaxActive   int64 `json:"max_active" properties:"max_active"`
	StatTime    int64 `json:"stat_time" properties:"stat_time"`
	Wait        bool  `json:"wait" properties:"wait"`
	WaitTimeout int64 `json:"wait_timeout" properties:"wait_timeout"`

	GetConnTimeout int64 `json:"get_conn_timeout" properties:"get_conn_timeout"`
}

const (
	APOLLO_CLUSTER    = "APOLLO_CLUSTER"
	APOLLO_IP         = "APOLLO_IP"
	APOLLO_PORT       = "APOLLO_PORT"
	APOLLO_SECRET_KEY = "APOLLO_SECRET_KEY"
)

type CenterConfigEnv struct {
	Cluster        string
	IP             string
	Port           string
	IsBackupConfig bool
	SecretKey      string
	MustStart      bool
}

func ParseConfigEnv() (*CenterConfigEnv, error) {
	apolloCluster, ok := os.LookupEnv(APOLLO_CLUSTER)
	if !ok {
		return nil, InvalidCluster
	}
	apolloIP, ok := os.LookupEnv(APOLLO_IP)
	if !ok {
		return nil, InvalidIp
	}
	apolloPort, ok := os.LookupEnv(APOLLO_PORT)
	if !ok {
		return nil, InvalidPort
	}

	return &CenterConfigEnv{
		Cluster:        apolloCluster,
		IP:             apolloIP,
		Port:           apolloPort,
		IsBackupConfig: false,
		SecretKey:      os.Getenv(APOLLO_SECRET_KEY),
		MustStart:      true,
	}, nil
}

// Invalid AppId format: Only digits, alphabets and symbol - _ . are allowed
func ReplaceServiceName(servName string) string {
	str := strings.ReplaceAll(servName, slash, dot)
	str = strings.ReplaceAll(str, comma, dot)
	str = strings.ReplaceAll(str, minus, dot)
	str = strings.ReplaceAll(str, under, dot)
	return str
}
