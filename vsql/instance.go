package vsql

import (
	"database/sql"
	"fmt"
	"time"
)

type Instance struct {
	dbName   string
	host     string
	username string
	password string

	maxIdleConn     int
	maxOpenConn     int
	maxConnLifetime time.Duration

	timeout, readTimeout, writeTimeout time.Duration

	db *sql.DB
}

func NewInstance(instanceConfig *InstanceConfig) (*Instance, error) {

	ins := &Instance{
		dbName:          instanceConfig.DBName,
		host:            instanceConfig.Host,
		username:        instanceConfig.Username,
		password:        instanceConfig.Password,
		maxIdleConn:     defaultMaxIdleConn,
		maxOpenConn:     defaultMaxOpenConn,
		maxConnLifetime: defaultMaxLifetime,
		timeout:         defaultTimeout,
		readTimeout:     defaultReadTimeout,
		writeTimeout:    defaultWriteTimeout,
	}

	if instanceConfig.MaxLifeTime > 0 {
		ins.maxConnLifetime = time.Duration(instanceConfig.MaxLifeTime) * time.Millisecond
	}
	if instanceConfig.MaxIdleConn > 0 {
		ins.maxIdleConn = instanceConfig.MaxIdleConn
	}
	if instanceConfig.MaxOpenConn > 0 {
		ins.maxOpenConn = instanceConfig.MaxOpenConn
	}
	if instanceConfig.Timeout > 0 {
		ins.timeout = time.Duration(instanceConfig.Timeout) * time.Millisecond
	}
	if instanceConfig.ReadTimeout > 0 {
		ins.readTimeout = time.Duration(instanceConfig.ReadTimeout) * time.Millisecond
	}
	if instanceConfig.WriteTimeout > 0 {
		ins.writeTimeout = time.Duration(instanceConfig.WriteTimeout) * time.Millisecond
	}

	err := ins.initDB()
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *Instance) initDB() error {

	db, err := NewConf(i.dbName, i.username, i.password, i.host).WithDSN(i.dsnList()...).Open()
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(i.maxIdleConn)
	db.SetMaxOpenConns(i.maxOpenConn)
	db.SetConnMaxLifetime(i.maxConnLifetime)

	i.db = db
	return nil
}

func (i *Instance) Reset(cfg *InstanceConfig) {
	i.maxIdleConn = cfg.MaxIdleConn
	i.maxOpenConn = cfg.MaxOpenConn
	i.maxConnLifetime = time.Duration(cfg.MaxLifeTime) * time.Millisecond

	i.db.SetMaxIdleConns(cfg.MaxIdleConn)
	i.db.SetMaxOpenConns(cfg.MaxOpenConn)
	i.db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Millisecond)
}

func (i *Instance) Close() error {
	return i.db.Close()
}

func (i *Instance) dsnList() []DSN {
	timeoutDSN := DSN(func() string {
		return fmt.Sprintf(dsnFormat, "timeout", i.timeout)
	})
	readTimeoutDSN := DSN(func() string {
		return fmt.Sprintf(dsnFormat, "readTimeout", i.readTimeout)
	})
	writeTimeoutDSN := DSN(func() string {
		return fmt.Sprintf(dsnFormat, "writeTimeout", i.writeTimeout)
	})

	return []DSN{
		utf8mb4DSN, collationDSN, locDSN, parseTimeDSN,
		timeoutDSN, writeTimeoutDSN, readTimeoutDSN,
	}
}
