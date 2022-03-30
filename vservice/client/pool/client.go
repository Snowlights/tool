package pool

import (
	"context"
	"sync"
	"time"
	"vtool/vconfig"
	"vtool/vservice/common"
)

type ClientPoolConfig struct {
	ServiceName string

	Idle, Active int64
	IdleTimeout  time.Duration

	Wait        bool
	WaitTimeOut time.Duration

	StatTime       time.Duration
	GetConnTimeout time.Duration
}

type ClientPool struct {
	cMu     sync.RWMutex
	conf    *ClientPoolConfig
	newConn func(string) (common.RpcConn, error)

	mu         sync.Mutex
	clientPool sync.Map
}

func NewClientPool(conf *ClientPoolConfig, newConn func(string) (common.RpcConn, error)) *ClientPool {
	return &ClientPool{conf: conf, newConn: newConn}
}

// todo reset pool config
func (c *ClientPool) ResetConfig(cfg *vconfig.ClientConfig) {
	c.cMu.Lock()
	defer c.cMu.Unlock()

	c.conf.Idle = cfg.Idle
	c.conf.Active = cfg.MaxActive
	c.conf.IdleTimeout = time.Duration(cfg.IdleTimeout)
	c.conf.Wait = cfg.Wait
	c.conf.WaitTimeOut = time.Duration(cfg.WaitTimeout)
	c.conf.StatTime = time.Duration(cfg.StatTime)
	c.conf.GetConnTimeout = cfg.GetConnTimeout
}

func (c *ClientPool) ResetConnConfig(cfg *vconfig.ClientConfig) {
	c.cMu.Lock()
	defer c.cMu.Unlock()

	c.clientPool.Range(func(key, value interface{}) bool {
		conn, ok := value.(*ConnPool)
		if !ok {
			return false
		}
		conn.ResetConfig(cfg)
		return true
	})

}

func (c *ClientPool) getConfig() *ClientPoolConfig {
	c.cMu.RLock()
	defer c.cMu.RUnlock()

	cfg := c.conf
	return cfg
}

func (c *ClientPool) getPool(serv *common.ServiceInfo) *ConnPool {
	var cp *ConnPool
	value, ok := c.clientPool.Load(serv.Addr)
	if ok == true {
		cp = value.(*ConnPool)
	} else {
		c.mu.Lock()
		defer c.mu.Unlock()
		value, ok := c.clientPool.Load(serv.Addr)
		if ok == true {
			cp = value.(*ConnPool)
		} else {
			cp = NewConnPool(&ConnPoolConfig{
				serviceName: c.conf.ServiceName,
				addr:        serv.Addr,
				idle:        c.conf.Idle,
				maxActive:   c.conf.Active,
				idleTimeout: c.conf.IdleTimeout,
				wait:        c.conf.Wait,
				waitTimeOut: c.conf.WaitTimeOut,
				statTime:    c.conf.StatTime,
			}, c.newConn)
			c.clientPool.Store(serv.Addr, cp)
		}
	}
	return cp
}

func (c *ClientPool) Delete(ctx context.Context, addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.clientPool.Load(addr)
	if !ok {
		return
	}
	connPool, ok := value.(*ConnPool)
	if !ok {
		return
	}
	c.clientPool.Delete(addr)
	connPool.Close()
	return
}

func (c *ClientPool) Get(ctx context.Context, serv *common.ServiceInfo) (common.RpcConn, error) {
	cp := c.getPool(serv)

	cfg := c.getConfig()
	timeout := cfg.GetConnTimeout
	if timeout == 0 {
		timeout = DefaultGetConnTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return cp.Get(ctx)
}

func (c *ClientPool) Close() {
	closeConnectionPool := func(key, value interface{}) bool {
		if connectionPool, ok := value.(*ConnPool); ok {
			connectionPool.Close()
		}
		return true
	}
	c.clientPool.Range(closeConnectionPool)
}

func (c *ClientPool) Put(ctx context.Context, serv *common.ServiceInfo, conn common.RpcConn) error {
	return c.getPool(serv).Put(ctx, conn)
}
