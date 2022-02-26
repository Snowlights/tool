package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func NewClient(cluster []string) (*Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cluster,
		DialTimeout: time.Second * 20,
	})

	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}
