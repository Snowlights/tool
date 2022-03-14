package zk

import (
	"context"
	"encoding/json"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"vtool/vlog"
	"vtool/vservice/common"
)

type Client struct {
	conn *zk.Conn

	servGroup string
	servName  string
	baseLoc   string

	servMu   sync.RWMutex
	servList []*common.RegisterServiceInfo
}

func NewZkClient(config *ClientConfig) (*Client, error) {
	timeOut := common.DefaultTTl
	if config.TimeOut > 0 {
		timeOut = config.TimeOut
	}

	conn, _, err := zk.Connect(config.Cluster, timeOut)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		conn:      conn,
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

func (c *Client) ensurePath(path string) error {
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

func (c *Client) watch(ctx context.Context) {
	_, _, watchChan, _ := c.conn.ChildrenW(c.servPath())

	go func() {
		for {
			msg := <-watchChan
			if msg.Type > 0 {
				err := c.reloadAllServ(ctx)
				if err != nil {
					vlog.ErrorF(ctx, "Client.watch.reloadAllServ failed, error is %s, event is %+v", err.Error(), msg.Type)
				}
			}
		}
	}()

}

func (c *Client) servPath() string {
	return c.baseLoc + common.Slash + c.servGroup + common.Slash + c.servName
}

func (c *Client) reloadAllServ(ctx context.Context) error {
	c.ensurePath(c.servPath())
	res, _, err := c.conn.Children(c.servPath())
	if err != nil {
		return err
	}

	servList := make([]*common.RegisterServiceInfo, 0, len(res))
	for _, child := range res {
		fullPath := c.servPath() + common.Slash + child
		data, _, err := c.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return err
		}
		val := make(map[common.ServiceType]*common.ServiceInfo, 1)
		err = json.Unmarshal(data, &val)
		if err != nil {
			continue
		}
		node := &common.RegisterServiceInfo{
			ServPath: fullPath,
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
