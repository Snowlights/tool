package vmq

import "sync"

type Manager struct {
	instances map[string]interface{}
	mutex     sync.Mutex
}
