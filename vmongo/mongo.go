package vmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"vtool/vconfig"
)

type Manager struct {
	insMu sync.RWMutex

	// cluster to db instance
	insMap map[string]*Instance

	center vconfig.Center
	cMu    sync.RWMutex
	cfg    *MongoConfig
}

// Collection is a handle to a MongoDB collection. It is safe for concurrent use by multiple goroutines.

func (m *Manager) Exec(ctx context.Context) error {
	return nil
}

func (m *Manager) getInstance(ctx context.Context) (*mongo.Client, error) {

	return nil, nil
}
