package vmq

import (
	"fmt"
	"strconv"
	"strings"
)

type RoleType int

const (
	RoleTypeKafkaReader RoleType = iota
	RoleTypeKafkaWriter

	confStrSep = "@"
)

const (
	WriteMsg = "WriteMsg"
	ReadMsg  = "ReadMsg"
	FetchMsg = "FetchMsg"
)

type KafkaConf struct {
	Conf map[string]*InstanceConf `json:"conf" properties:"conf"`
}

type InstanceConf struct {
	Brokers        []string `json:"brokers" properties:"brokers"`
	CommitInterval int      `json:"commitInterval" properties:"commit_interval"`
	StartOffset    int64    `json:"startOffset" properties:"start_offset"`
	MinBytes       int      `json:"minBytes" properties:"min_bytes"`
	MaxBytes       int      `json:"maxBytes" properties:"max_bytes"`
}

type Conf struct {
	cluster   string
	topic     string
	group     string
	partition int
	role      RoleType
}

func (c *Conf) String() string {
	return strings.Join([]string{c.cluster, c.topic,
		c.group, fmt.Sprintf("%d", c.partition),
		fmt.Sprintf("%d", c.role)}, confStrSep)
}

func confStrToConf(s string) *Conf {
	c := &Conf{}
	parts := strings.Split(s, confStrSep)

	if len(parts) != 5 {
		return nil
	}
	c.cluster = parts[0]
	c.topic = parts[1]
	c.group = parts[2]
	c.partition, _ = strconv.Atoi(parts[3])
	role, _ := strconv.Atoi(parts[4])
	c.role = RoleType(role)
	return c
}
