package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
	"vtool/vlog"
)

type Register struct {
	client *clientv3.Client
}

func (c *Register) Register(ctx context.Context, path, val string, ttl time.Duration) error {

	kv := clientv3.NewKV(c.client)

	lease := clientv3.NewLease(c.client)
	leaseRes, err := lease.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return err
	}

	_, err = kv.Put(ctx, path, val, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		return err
	}

	keepAliveRes, err := lease.KeepAlive(ctx, leaseRes.ID)
	if err != nil {
		return err
	}

	go c.keepAlive(ctx, keepAliveRes, path, val)
	return nil
}

func (c *Register) keepAlive(ctx context.Context, keepAliveRes <-chan *clientv3.LeaseKeepAliveResponse, path, val string) {
	for {
		select {
		case ret := <-keepAliveRes:
			if ret != nil {
				vlog.Info(ctx, strings.Join([]string{path, val, leaseSuccess}, character))
			}
		}
	}
}

func (c *Register) Get(ctx context.Context, path string) (string, error) {

	res, err := c.client.Get(ctx, path)
	if err != nil {
		return "", err
	}

	if len(res.Kvs) > 0 {
		return string(res.Kvs[0].Value), nil
	}

	return "", nil
}

func (c *Register) GetNode(ctx context.Context, path string) ([]*Node, error) {

	res, err := c.client.Get(ctx, path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	nodeList := make([]*Node, 0, len(res.Kvs))

	for k, v := range res.Kvs {
		valStr := string(v.Value)
		node := &Node{
			key:   string(rune(k)),
			val:   valStr,
			lease: v.Lease,
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

func (c *Register) RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error {
	kv := clientv3.NewKV(c.client)

	lease := clientv3.NewLease(c.client)
	leaseRes, err := lease.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return err
	}

	_, err = kv.Put(ctx, path, val, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		return err
	}

	_, err = lease.KeepAliveOnce(ctx, leaseRes.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Register) Watch(ctx context.Context, path string) clientv3.WatchChan {
	return c.client.Watch(ctx, path, clientv3.WithPrefix())
}
