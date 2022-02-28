package zk

import (
	"context"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type Register struct {
	conn *zk.Conn
}

func (c *Register) Register(ctx context.Context, path, val string, ttl time.Duration) error {

	if err := c.ensureName(path); err != nil {
		return err
	}

	_, err := c.conn.CreateProtectedEphemeralSequential(path, []byte(val), zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

func (c *Register) ensureName(path string) error {
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

	return "", nil
}

func (c *Register) GetNode(ctx context.Context, path string) []*Node {

	return nil
}

func (c *Register) RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error {

	return nil
}
