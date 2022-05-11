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

type instanceConf struct {
	group     string
	role      RoleType
	topic     string
	groupId   string
	partition int
}
