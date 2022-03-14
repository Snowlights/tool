package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"vtool/vlog"
	"vtool/vservice/common"
)

type Client struct {
	client *clientv3.Client

	servGroup string
	servName  string
	baseLoc   string

	servMu   sync.RWMutex
	servList []*common.RegisterServiceInfo
}

func NewEtcdClient(config *ClientConfig) (*Client, error) {
	timeOut := common.DefaultTTl
	if config.TimeOut > 0 {
		timeOut = config.TimeOut
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Cluster,
		DialTimeout: timeOut,
	})

	if err != nil {
		return nil, err
	}

	cli := &Client{
		client:    client,
		servGroup: config.ServGroup,
		servName:  config.ServName,
		baseLoc:   common.DefaultRegisterPath,
	}
	return cli, nil
}

func (c *Client) GetAllServAddr() []*common.RegisterServiceInfo {
	c.servMu.RLock()
	defer c.servMu.RUnlock()

	nodeList := make([]*common.RegisterServiceInfo, 0, len(c.servList))
	for _, v := range c.servList {
		val := make(map[common.ServiceType]*common.ServiceInfo, len(v.ServList))
		for path, serv := range v.ServList {
			val[path] = &common.ServiceInfo{
				Type: serv.Type,
				Addr: serv.Addr,
			}
		}

		node := &common.RegisterServiceInfo{
			ServPath: v.ServPath,
			ServList: val,
		}
		nodeList = append(nodeList, node)
	}

	return nodeList
}

func (c *Client) watch(ctx context.Context) {
	watchChan := c.client.Watch(ctx, c.servPath(), clientv3.WithPrefix())

	go func() {
		for {
			msg := <-watchChan
			for _, event := range msg.Events {
				err := c.reloadAllServ(ctx)
				if err != nil {
					vlog.ErrorF(ctx, "Client.watch.reloadAllServ failed, error is %s, event is %+v", err.Error(), event)
				}
			}
		}
	}()

}

func (c *Client) servPath() string {
	return c.baseLoc + common.Slash + c.servGroup + common.Slash + c.servName
}

func (c *Client) reloadAllServ(ctx context.Context) error {
	res, err := c.client.Get(ctx, c.servPath(), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	servList := make([]*common.RegisterServiceInfo, 0, len(res.Kvs))
	for _, v := range res.Kvs {
		val := make(map[common.ServiceType]*common.ServiceInfo, 1)
		err = json.Unmarshal(v.Value, &val)
		if err != nil {
			continue
		}
		node := &common.RegisterServiceInfo{
			ServPath: string(v.Key),
			ServList: val,
		}
		servList = append(servList, node)
	}

	c.updateServ(servList)
	return nil
}

func (c *Client) updateServ(servList []*common.RegisterServiceInfo) {
	c.servMu.Lock()
	defer c.servMu.Unlock()

	c.servList = servList
}
