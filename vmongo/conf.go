package vmongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

const (
	mongoDbFormat = "mongodb://%s:%s@%s/%s?%s"

	dsnCharacter = "&"
)

type DSN func() string

var (
	connectTimeout = DSN(func() string {
		return "connectTimeoutMS=30000"
	})
	socketTimeout = DSN(func() string {
		return "socketTimeoutMS=30000"
	})
	readPreference = DSN(func() string {
		return "readPreference=primaryPreferred"
	})
	writeConcern = DSN(func() string {
		return "w=majority"
	})
	readConcern = DSN(func() string {
		return "readConcernLevel=majority"
	})
)

type Conf struct {
	dbName   string
	Host     string
	user     string
	password string
	dsn      []DSN
}

func NewConf(dbName, user, pwd, host string) *Conf {
	return &Conf{
		dbName:   dbName,
		Host:     host,
		user:     user,
		password: pwd,
	}
}

func (c *Conf) WithDSN(dsn ...DSN) *Conf {
	c.dsn = append(c.dsn, dsn...)
	return c
}

func (c *Conf) Open() (*mongo.Client, error) {
	ctx := context.Background()
	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI(c.formatDataSourceName()))
	if err != nil {
		return nil, err
	}

	return mongoCli, mongoCli.Ping(ctx, nil)
}

func (c *Conf) formatDSN() string {
	var dsn []string
	for _, d := range c.dsn {
		dsn = append(dsn, d())
	}
	return strings.Join(dsn, dsnCharacter)
}

func (c *Conf) formatDataSourceName() string {
	// like: mongodb://user1:pwd1@localhost:40000/test
	return fmt.Sprintf(mongoDbFormat, c.user, c.password, c.Host, c.dbName, c.formatDSN())
}
