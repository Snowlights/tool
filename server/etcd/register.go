package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
	"vtool/vlog"
)

type Client struct {
	client *clientv3.Client
}

func (c *Client) Register(ctx context.Context, path, val string, ttl time.Duration) error {

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

func (c *Client) keepAlive(ctx context.Context, keepAliveRes <-chan *clientv3.LeaseKeepAliveResponse, path, val string) {
	for {
		select {
		case ret := <-keepAliveRes:
			if ret != nil {
				vlog.Info(ctx, strings.Join([]string{path, val, leaseSuccess}, character))
			}
		}
	}
}

func (c *Client) Get(ctx context.Context, path string) (string, error) {
	// todo 实验

	return "", nil
}

func (c *Client) GetNode() {

}

func (c *Client) SetTtl(ctx context.Context, path, val string, ttl time.Duration) error {
	return nil
}

func (c *Client) RefreshTtl(ctx context.Context, path string, ttl time.Duration) error {
	return nil
}
