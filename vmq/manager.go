package vmq

import (
	"sync"
)

var defaultManager *Manager

type Manager struct {
	instances map[string]interface{}
	mutex     sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		instances: make(map[string]interface{}),
	}
}

func WriteMsgWithTopic(cluster, topic string, msgs ...Message) error {

	return nil
}
