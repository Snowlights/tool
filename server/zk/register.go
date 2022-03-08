package zk

import (
	"context"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
	"vtool/server/common"
)

type Register struct {
	conn *zk.Conn
}

// ZK new nodes need to be built from top to bottom
// example :
// /public/base/web/service1
// first /public/
// next /public/base/
// then /public/base/web/
// final /public/base/web/service1

func (c *Register) Register(ctx context.Context, path, val string, ttl time.Duration) error {
	if err := c.ensureAllPathExit(path); err != nil {
		return err
	}

	_, err := c.conn.Create(path, []byte(val), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	return nil
}

// It is not recommended to use the new method here. It is
// recommended to use a special group for maintenance.
//It is maintained uniformly within the company through work orders or other methods
func (c *Register) ensureAllPathExit(path string) error {
	parts := strings.Split(path, Slash)

	if len(parts) == 1 {
		return c.ensurePath(path)
	}

	i := 2
	for i < len(parts) {
		err := c.ensurePath(strings.Join(parts[:i], Slash))
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func (c *Register) ensurePath(path string) error {
	exists, _, err := c.conn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		_, err := c.conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (c *Register) Get(ctx context.Context, path string) (string, error) {

	res, _, err := c.conn.Get(path)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (c *Register) GetNode(ctx context.Context, path string) ([]common.Node, error) {
	c.ensurePath(path)
	res, _, err := c.conn.Children(path)
	if err != nil {
		return nil, err
	}

	nodeList := make([]common.Node, 0, len(res))
	for _, child := range res {
		fullPath := path + Slash + child
		data, _, err := c.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return nil, err
		}
		valStr := string(data)
		node := &Node{
			key: fullPath,
			val: valStr,
		}
		parts := strings.Split(valStr, character)
		if len(parts) == 2 {
			node.ip = parts[0]
			node.port = parts[1]
		}
		nodeList = append(nodeList, node)
	}

	return nodeList, nil
}

func (c *Register) Watch(ctx context.Context, path string) (chan common.Event, error) {
	snapshots := make(chan []string)
	eventChan := make(chan common.Event)
	errors := make(chan error)

	go func() {
		for {
			snapshot, _, events, err := c.conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			// todo if you would use snapshot you can do something here
			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
			eventChan <- Event{common.ChildrenChanged}
		}
	}()

	return eventChan, nil
}

// TTL node is added in version 3.5.5
// In version 3.5.8, zookeeper does not support nodes with expiration time by default.
// However, in 3.6.3, the expiration time is directly supported by default.
// The referenced SDK has no - t option and will not be demonstrated here.
// It is generally used as a distributed lock. In this scenario,
// the - e parameter is used to meet the requirements of
// the service registration and discovery scenario
func (c *Register) RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error {

	return nil
}
