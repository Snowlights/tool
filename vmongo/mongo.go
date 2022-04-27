package vmongo

import (
	"context"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"sync"
	"vtool/parse"
	"vtool/vconfig"
	"vtool/vtrace"
)

const (
	mongoExec = "mongoExec"
)

var (
	ReadPrefPrimary            = readpref.Primary
	ReadPrefPrimaryPreferred   = readpref.PrimaryPreferred
	ReadPrefSecondary          = readpref.Secondary
	ReadPrefSecondaryPreferred = readpref.SecondaryPreferred
	ReadPrefNearest            = readpref.Nearest
)

var defaultManager *Manager

type Manager struct {
	insMu sync.RWMutex

	// cluster to db instance
	insMap map[string]*Instance

	center vconfig.Center
	cMu    sync.RWMutex
	cfg    *MongoConfig
}

func Exec(ctx context.Context, cluster, database, collection string, query func(c *mongo.Collection) error) error {
	return defaultManager.Exec(ctx, cluster, database, collection, query, nil, nil, nil)
}

func ExecWithOpt(ctx context.Context, cluster, database, collection string, query func(c *mongo.Collection) error,
	readPref *readpref.ReadPref, readConcern *readconcern.ReadConcern,
	writeConcern *writeconcern.WriteConcern) error {
	return defaultManager.Exec(ctx, cluster, database, collection, query, readPref, readConcern, writeConcern)
}

func NewManager(center vconfig.Center) (*Manager, error) {
	manager := &Manager{
		insMap: make(map[string]*Instance),
		center: center,
	}

	err := manager.loadConfig()
	if err != nil {
		return nil, err
	}

	manager.center.AddListener(&MongoListener{manager.changeEvent})
	defaultManager = manager
	return defaultManager, nil
}

func (m *Manager) Exec(ctx context.Context, cluster, database, collection string, query func(c *mongo.Collection) error,
	readPref *readpref.ReadPref, readConcern *readconcern.ReadConcern,
	writeConcern *writeconcern.WriteConcern) error {
	if readPref == nil {
		readPref = ReadPrefPrimaryPreferred()
	}

	return m.exec(ctx, cluster, database, collection, query, readPref, readConcern, writeConcern)
}

func (m *Manager) exec(ctx context.Context, cluster, database, collection string, query func(c *mongo.Collection) error,
	readPref *readpref.ReadPref, readConcern *readconcern.ReadConcern,
	writeConcern *writeconcern.WriteConcern) error {

	ins, err := m.getInstance(ctx, cluster)
	if err != nil {
		return err
	}

	collectionIns := ins.Database(database).Collection(collection, &options.CollectionOptions{
		ReadConcern:    readConcern,
		WriteConcern:   writeConcern,
		ReadPreference: readPref,
	})
	span, ctx := opentracing.StartSpanFromContext(ctx, mongoExec)
	defer span.Finish()
	span.SetTag(vtrace.MongoCluster, cluster)
	span.SetTag(vtrace.DataBase, database)
	span.SetTag(vtrace.Collection, collection)
	span.SetTag(vtrace.Component, vtrace.ComponentMongo)
	span.SetTag(vtrace.SpanKind, vtrace.SpanKindMongo)
	return query(collectionIns)
}

func (m *Manager) loadConfig() error {
	cfg := new(MongoConfig)
	err := m.center.UnmarshalWithNameSpace(vconfig.Mongo, parse.PropertiesTagName, cfg)
	if err != nil {
		return err
	}

	m.setConfig(cfg)
	return nil
}

func (m *Manager) changeEvent() {
	m.loadConfig()
}

func (m *Manager) setConfig(cfg *MongoConfig) {
	m.cMu.Lock()
	defer m.cMu.Unlock()
	m.cfg = cfg
}

// Collection is a handle to a MongoDB collection. It is safe for concurrent use by multiple goroutines.

func (m *Manager) getConfig() *MongoConfig {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	cfg := m.cfg
	return cfg
}

func (m *Manager) getInstance(ctx context.Context, cluster string) (*mongo.Client, error) {

	m.insMu.RLock()
	ins, ok := m.insMap[cluster]
	if ok {
		m.insMu.RUnlock()
		return ins.client, nil
	}

	m.insMu.RUnlock()
	m.insMu.Lock()
	defer m.insMu.Unlock()
	ins, ok = m.insMap[cluster]
	if ok {
		return ins.client, nil
	}

	cfg := m.getConfig()
	if cfg == nil {
		return nil, NotInitManager
	}

	instanceConfig, ok := cfg.Conf[cluster]
	if !ok {
		return nil, NotFoundCluster
	}

	newIns, err := NewInstance(instanceConfig)
	if err != nil {
		return nil, err
	}
	m.insMap[cluster] = newIns
	return newIns.client, nil
}

type MongoListener struct {
	Change func()
}

func (cl *MongoListener) OnChange(event *storage.ChangeEvent) {

}

func (cl *MongoListener) OnNewestChange(event *storage.FullChangeEvent) {
	if event.Namespace != vconfig.Mongo {
		return
	}
	cl.Change()
}
