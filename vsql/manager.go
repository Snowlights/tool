package vsql

import (
	"sync"
	"vtool/vconfig"
)

var defaultManager *Manager

type Manager struct {
	insMu  sync.RWMutex
	insMap map[string]*Instance

	center vconfig.Center
}

func GetDB() *DB {
	return defaultManager.GetDB()
}

func InitManager() *Manager {

	// todo init with config
	manager := &Manager{}
	defaultManager = manager

	return defaultManager
}

func (m *Manager) GetDB() *DB {

	// todo: get instance and return
	return nil
}
