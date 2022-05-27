package vmq

import (
	"context"
	"fmt"
	"github.com/apolloconfig/agollo/v4/storage"
	"sync"
	"time"
	"vtool/parse"
	"vtool/vconfig"
	"vtool/vlog"
	"vtool/vmq/vkafka"
)

var defaultManager *Manager

type Manager struct {
	instances map[string]interface{}
	mutex     sync.Mutex

	center    vconfig.Center
	cMu       sync.RWMutex
	kafkaConf *KafkaConf
}

func NewManager(center vconfig.Center) (*Manager, error) {
	manager := &Manager{
		instances: make(map[string]interface{}),
	}

	err := manager.loadConfig()
	if err != nil {
		return nil, err
	}
	manager.center.AddListener(&KafkaListener{manager.changeEvent})
	defaultManager = manager
	return defaultManager, nil
}

func (m *Manager) getKafkaReader(ctx context.Context, conf *Conf) *vkafka.Reader {

	ins := m.getIns(ctx, conf)
	reader, ok := ins.(*vkafka.Reader)
	if ok == false {
		vlog.ErrorF(ctx, "in.(Reader) err, topic: %s", conf.topic)
		return nil
	}

	return reader
}

func (m *Manager) getKafkaWriter(ctx context.Context, conf *Conf) *vkafka.Writer {

	ins := m.getIns(ctx, conf)
	writer, ok := ins.(*vkafka.Writer)
	if ok == false {
		vlog.ErrorF(ctx, "in.(Writer) err, topic: %s", conf.topic)
		return nil
	}

	return writer
}

func (m *Manager) getIns(ctx context.Context, conf *Conf) interface{} {

	var in interface{}
	in, ok := m.instances[conf.String()]
	if !ok {
		m.mutex.Lock()
		in, ok = m.instances[conf.String()]
		if ok {
			m.mutex.Unlock()
			return in
		}

		vlog.InfoF(ctx, "newInstance, role:%v, topic: %s", conf.role, conf.topic)
		newIns, err := m.newInstance(ctx, conf)
		if err != nil {
			vlog.ErrorF(ctx, "newInstance err, topic: %s, err: %s", conf.topic, err.Error())
			m.mutex.Unlock()
			return nil
		}
		m.instances[conf.String()] = newIns
		in, _ = m.instances[conf.String()]
		m.mutex.Unlock()
	}
	return in
}

func (m *Manager) getKafkaConfigWithCluster(cluster string) (*InstanceConf, bool) {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	conf, ok := m.kafkaConf.Conf[cluster]
	if ok {
		return conf, true
	}
	return nil, false
}

func (m *Manager) newInstance(ctx context.Context, conf *Conf) (interface{}, error) {

	insConfig, ok := m.getKafkaConfigWithCluster(conf.cluster)
	if !ok {
		return nil, fmt.Errorf("get kafka conf failed, cluster is %s", conf.cluster)
	}

	switch conf.role {
	case RoleTypeKafkaReader:
		return vkafka.NewKafkaReader(&vkafka.KafkaReaderConf{
			Brokers:        insConfig.Brokers,
			Topic:          conf.topic,
			Group:          conf.group,
			Partition:      conf.partition,
			CommitInterval: time.Millisecond * time.Duration(insConfig.CommitInterval),
			MinByte:        insConfig.MinBytes,
			MaxByte:        insConfig.MaxBytes,
			StartOffset:    insConfig.StartOffset,
		}), nil
	case RoleTypeKafkaWriter:
		return vkafka.NewKafkaWriter(&vkafka.KafkaWriterConf{
			Brokers: insConfig.Brokers,
			Topic:   conf.topic,
		}), nil
	default:
		return nil, fmt.Errorf("role %d error", conf.role)
	}
}

func (m *Manager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, ins := range m.instances {
		switch ins.(type) {
		case *vkafka.Reader:
			ins.(*vkafka.Reader).Close()
		case *vkafka.Writer:
			ins.(*vkafka.Writer).Close()
		}
	}
	m.instances = make(map[string]interface{})
}

func (m *Manager) loadConfig() error {
	cfg := new(KafkaConf)
	err := m.center.UnmarshalWithNameSpace(vconfig.Kafka, parse.PropertiesTagName, cfg)
	if err != nil {
		return err
	}

	m.setConfig(cfg)
	return nil
}

func (m *Manager) changeEvent() {
	m.loadConfig()
}

func (m *Manager) setConfig(cfg *KafkaConf) {
	m.cMu.Lock()
	defer m.cMu.Unlock()
	m.kafkaConf = cfg
}

type KafkaListener struct {
	Change func()
}

func (cl *KafkaListener) OnChange(event *storage.ChangeEvent) {

}

func (cl *KafkaListener) OnNewestChange(event *storage.FullChangeEvent) {
	if event.Namespace != vconfig.Kafka {
		return
	}
	cl.Change()
}
