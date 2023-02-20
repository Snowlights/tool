package vsql

import (
	"github.com/Snowlights/tool/parse"
	"github.com/Snowlights/tool/vconfig"
	"sync"
	"time"
)

var defaultManager *Manager

type Manager struct {
	insMu sync.RWMutex

	// cluster to db instance
	insMap map[string]*Instance

	center vconfig.Center
	cMu    sync.RWMutex
	cfg    *DBConfig
}

func GetDB(cluster string) (*DB, error) {
	if defaultManager == nil {
		return nil, NotInitManager
	}

	return defaultManager.getDB(cluster)
}

func InitManager(center vconfig.Center) (*Manager, error) {

	manager := &Manager{
		insMap: make(map[string]*Instance),
		center: center,
	}
	defaultManager = manager

	err := manager.loadConfig()
	if err != nil {
		return nil, err
	}

	center.AddListener(&MysqlListener{manager.changeEvent})

	return defaultManager, nil
}

func (m *Manager) changeEvent() {
	m.loadConfig()
}

func (m *Manager) loadConfig() error {
	cfg := new(DBConfig)
	err := m.center.UnmarshalWithNameSpace(vconfig.ServerDB, parse.PropertiesTagName, cfg)
	if err != nil {
		return err
	}
	m.setConfig(cfg)
	m.resetAndCLose()
	return nil
}

func (m *Manager) setConfig(cfg *DBConfig) {
	m.cMu.Lock()
	defer m.cMu.Unlock()

	m.cfg = cfg
}

func (m *Manager) getConfig() *DBConfig {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	cfg := m.cfg
	return cfg
}

func (m *Manager) resetAndCLose() {
	cfg := m.getConfig()

	m.insMu.Lock()
	defer m.insMu.Unlock()

	closeInsList, resetInsList := make([]*Instance, 0), make(map[*Instance]InstanceConfig)
	for cluster, ins := range m.insMap {
		instanceConfig, ok := cfg.Conf[cluster]
		if !ok {
			closeInsList = append(closeInsList, ins)
			continue
		}

		switch m.compareInstanceCfgChanged(instanceConfig, ins) {
		case Reopen:
			closeInsList = append(closeInsList, ins)
		case Reset:
			resetInsList[ins] = instanceConfig
		}

	}

	for _, ins := range closeInsList {
		ins.Close()
	}
	for ins, dbNameCfg := range resetInsList {
		ins.Reset(&dbNameCfg)
	}
}

func (m *Manager) compareInstanceCfgChanged(insCfg InstanceConfig, ins *Instance) ChangeType {
	if time.Duration(insCfg.ReadTimeout)*time.Millisecond != ins.readTimeout ||
		time.Duration(insCfg.ReadTimeout)*time.Millisecond != ins.writeTimeout {
		return Reopen
	}

	if insCfg.MaxIdleConn != ins.maxIdleConn ||
		insCfg.MaxOpenConn == ins.maxOpenConn ||
		time.Duration(insCfg.MaxLifeTime)*time.Millisecond != ins.maxConnLifetime {
		return Reset
	}
	return 0
}

func (m *Manager) getDB(cluster string) (*DB, error) {
	if m.center == nil {
		return nil, NotInitManager
	}

	ins, err := m.getInstance(cluster)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:      ins.db,
		cluster: cluster,
		schema:  ins.dbName,
	}, nil
}

func (m *Manager) getInstance(cluster string) (*Instance, error) {
	if m.center == nil {
		return nil, NotInitManager
	}

	m.insMu.RLock()
	ins, ok := m.insMap[cluster]
	if ok {
		m.insMu.RUnlock()
		return ins, nil
	}

	m.insMu.RUnlock()
	m.insMu.Lock()
	defer m.insMu.Unlock()

	cfg := m.getConfig()
	if cfg == nil {
		return nil, NotInitManager
	}

	instanceConfig, ok := cfg.Conf[cluster]
	if !ok {
		return nil, NotFoundCluster
	}

	newIns, err := NewInstance(&instanceConfig)
	if err != nil {
		return nil, err
	}

	m.insMap[cluster] = newIns
	return newIns, nil
}
