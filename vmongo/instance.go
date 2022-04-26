package vmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Instance struct {
	document string
	host     string
	username string
	password string

	timeout      time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	poolSize     int64

	client *mongo.Client
}

func NewInstance(config InstanceConfig) (*Instance, error) {

	ins := &Instance{
		document:     config.Document,
		host:         config.Host,
		username:     config.Username,
		password:     config.Password,
		timeout:      defaultTimeout,
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
		poolSize:     defaultPoolLimit,
	}
	if config.PoolSize > 0 {
		ins.poolSize = config.PoolSize
	}
	if config.Timeout > 0 {
		ins.timeout = time.Duration(config.Timeout) * time.Millisecond
	}
	if config.ReadTimeout > 0 {
		ins.readTimeout = time.Duration(config.ReadTimeout) * time.Millisecond
	}
	if config.WriteTimeout > 0 {
		ins.writeTimeout = time.Duration(config.WriteTimeout) * time.Millisecond
	}

	err := ins.initClient()
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *Instance) initClient() error {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()

	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI(i.host).SetAuth(options.Credential{
		Username: i.username,
		Password: i.password}).
		SetMaxPoolSize(uint64(i.poolSize)))
	if err != nil {
		return err
	}

	i.client = mongoCli
	return nil
}

func (i *Instance) Close() error {
	return i.client.Disconnect(context.TODO())
}
