package pool

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
	"vtool/vlog"
	"vtool/vservice/common"
)

const (
	DefaultStatTime    = time.Millisecond * 100
	DefaultIdle        = 128
	DefaultMaxActive   = 256
	DefaultIdleTimeout = time.Minute
)

type ConnPool struct {
	newConn func(string) (common.RpcConn, error)

	cond   chan struct{}
	active int

	conf *ConnPoolConfig

	mu       sync.Mutex
	closed   bool
	connList *list.List
}

type ConnPoolConfig struct {
	serviceName     string
	addr            string
	idle, maxActive int
	idleTimeout     time.Duration

	wait        bool
	waitTimeOut time.Duration

	statTime time.Duration
}

type connItem struct {
	addTime time.Time
	conn    common.RpcConn
}

// todo reset pool config

func (ci *connItem) expire(timeout time.Duration) bool {
	if timeout <= 0 {
		return false
	}
	return ci.addTime.Add(timeout).Before(time.Now())
}

func NewConnPool(conf *ConnPoolConfig, newConn func(string) (common.RpcConn, error)) *ConnPool {

	if conf.statTime == 0 {
		conf.statTime = DefaultStatTime
	}
	if conf.idle == 0 {
		conf.idle = DefaultIdle
	}
	if conf.maxActive == 0 {
		conf.idle = DefaultMaxActive
	}
	if conf.maxActive > conf.idle {
		conf.maxActive = DefaultMaxActive
		conf.idle = DefaultIdle
	}

	c := &ConnPool{
		newConn:  newConn,
		cond:     make(chan struct{}),
		conf:     conf,
		closed:   false,
		connList: list.New(),
	}
	go c.stat()
	return c
}

func (cp *ConnPool) stat() {
	if cp.conf.idleTimeout == 0 {
		return
	}

	ticker := time.NewTicker(cp.conf.statTime)
	for {
		vlog.DebugF(context.Background(), fmt.Sprintf("ConnPool.stat, cp.active:%d, cp.idle:%d ", cp.active, cp.connList.Len()))
		select {
		case <-ticker.C:
			cp.mu.Lock()
			if cp.closed || cp.conf.idleTimeout <= 0 {
				cp.mu.Unlock()
				return
			}
			for i, n := 0, cp.connList.Len(); i < n; i++ {
				e := cp.connList.Back()
				ci := e.Value.(connItem)
				if !ci.expire(cp.conf.idleTimeout) {
					continue
				}
				cp.connList.Remove(e)
				cp.active--
				cp.mu.Unlock()
				ci.conn.Close()
				cp.mu.Lock()
			}
			cp.mu.Unlock()
		}
	}
}

func (cp *ConnPool) Close() error {
	cp.mu.Lock()
	connList := cp.connList
	cp.connList.Init()
	cp.closed = true
	cp.active -= connList.Len()
	cp.mu.Unlock()

	for e := connList.Front(); e != nil; e = e.Next() {
		e.Value.(connItem).conn.Close()
	}
	return nil
}

func (cp *ConnPool) Get(ctx context.Context) (common.RpcConn, error) {
	cp.mu.Lock()
	if cp.closed {
		cp.mu.Unlock()
		return nil, nil
	}
	fmt.Println("Get", "cp.active:", cp.active, "cp.idle:", cp.connList.Len())
	for {
		for index, length := 0, cp.connList.Len(); index < length; index++ {
			e := cp.connList.Front()
			if e == nil {
				break
			}
		}

		if cp.closed {
			cp.mu.Unlock()
			return nil, nil
		}

		if cp.active < cp.conf.maxActive {
			newItem := cp.newConn
			cp.active++
			cp.mu.Unlock()

			c, err := newItem(cp.conf.addr)
			if err != nil {
				cp.mu.Lock()
				cp.active--
				cp.mu.Unlock()
				c = nil
			}
			return c, err
		}

		if !cp.conf.wait && cp.conf.waitTimeOut == 0 {
			cp.mu.Unlock()
			return nil, nil
		}

		timeOut := cp.conf.waitTimeOut
		cp.mu.Unlock()

		newCtx := ctx
		cancel := func() {}
		if timeOut > 0 {
			newCtx, cancel = context.WithTimeout(ctx, timeOut)
		}

		select {
		case <-newCtx.Done():
			cancel()
			return nil, newCtx.Err()
		case <-cp.cond:
		}

		cancel()
		cp.mu.Lock()
	}

}

func (cp *ConnPool) Put(ctx context.Context, rpcConn common.RpcConn) error {
	cp.mu.Lock()
	fmt.Println("Put", "cp.active:", cp.active, "cp.idle:", cp.connList.Len())

	if cp.closed {
		cp.mu.Unlock()
		return nil
	}

	if cp.connList.Len() > cp.active {
		cp.active--
		cp.mu.Unlock()
		rpcConn.Close()
		return nil
	}

	cp.connList.PushFront(connItem{
		addTime: time.Now(),
		conn:    rpcConn,
	})
	if cp.connList.Len() > cp.conf.idle {
		rpcConn = cp.connList.Remove(cp.connList.Back()).(connItem).conn
	} else {
		rpcConn = nil
	}
	cp.active--
	cp.mu.Unlock()

	if rpcConn == nil {
		return nil
	}

	return rpcConn.Close()
}
