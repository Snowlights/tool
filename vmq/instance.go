package vmq

type Message struct {
	key string
	val interface{}
}

type RoleType int

const (
	RoleTypeKafkaReader RoleType = iota
	RoleTypeKafkaWriter

	instanceConfSep = "@"
)

type InstanceConf struct {
	Cluster        string   `json:"cluster" properties:"cluster"`
	Brokers        []string `json:"brokers" properties:"brokers"`
	Topic          string   `json:"topic" properties:"topic"`
	Group          string   `json:"group" properties:"group"`
	Role           RoleType `json:"-" properties:"-"`
	Partition      int      `json:"partition" properties:"partition"`
	CommitInterval int      `json:"commitInterval" properties:"commit_interval"`
	StartOffset    int64    `json:"startOffset" properties:"start_offset"`
	MinBytes       int      `json:"minBytes" properties:"min_bytes"`
	MaxBytes       int      `json:"maxBytes" properties:"max_bytes"`
}
