package zk

import (
	"context"
	"encoding/json"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
	"sync"
	"vtool/load_balance/consistent"
	"vtool/vlog"
	"vtool/vservice/common"
)

type Client struct {
	conn *zk.Conn

	servGroup string
	servName  string
	baseLoc   string

	// Here, in order to ensure that all services are available,
	// mutex locks are used instead of read-write locks
	servMu       sync.Mutex
	servList     []*common.RegisterServiceInfo
	servLaneHash map[string]*consistent.Consistent

	poolMu      sync.Mutex
	handlerList []func([]string)
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
		conn:         conn,
		servGroup:    config.ServGroup,
		servName:     config.ServName,
		baseLoc:      common.DefaultRegisterPath,
		servLaneHash: make(map[string]*consistent.Consistent),
	}
	cli.reloadAllServ(context.Background())
	go cli.watch(context.Background())
	return cli, nil
}

func (c *Client) AddPoolHandler(handle func([]string)) {
	c.poolMu.Lock()
	defer c.poolMu.Unlock()

	c.handlerList = append(c.handlerList, handle)
}

func (c *Client) getPoolHandler() []func([]string) {
	c.poolMu.Lock()
	defer c.poolMu.Unlock()

	return c.handlerList
}

func (c *Client) resetPool(addr []string) {
	for _, handle := range c.getPoolHandler() {
		handle(addr)
	}
}

func (c *Client) ServName() string {
	return c.servName
}

func (c *Client) ServGroup() string {
	return c.servGroup
}

func (c *Client) GetAllServAddr() []*common.RegisterServiceInfo {
	c.servMu.Lock()
	defer c.servMu.Unlock()

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

func (c *Client) GetServAddr(lane string, serviceType common.ServiceType, hashKey string) (*common.ServiceInfo, bool) {
	ctx := context.Background()
	c.servMu.Lock()
	defer c.servMu.Unlock()

	if c.servLaneHash == nil {
		return nil, false
	}

	hash, ok := c.servLaneHash[lane]
	if !ok {
		hash, ok = c.servLaneHash[""]
		if !ok {
			vlog.ErrorF(ctx, "c.servLaneHash[] == nil, serv path:%s, lane:%s, key:%s", c.servPath(), lane, hashKey)
			return nil, false
		}
	}

	key, err := hash.Get(hashKey)
	if err != nil {
		vlog.ErrorF(ctx, "hash get serv key failed, error is %s, serv path:%s, lane:%s, key:%s", err.Error(), c.servPath(), lane, hashKey)
		return nil, false
	}

	servPathPartIndex := strings.LastIndex(key, common.HashKey)
	servPath := key[:servPathPartIndex]

	for _, serv := range c.servList {
		if servPath == serv.ServPath && lane == serv.Lane {
			servInfo, ok := serv.ServList[serviceType]
			if ok {
				return &common.ServiceInfo{
					Type: servInfo.Type,
					Addr: servInfo.Addr,
				}, true
			}
		}
	}

	return nil, false
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
	servLaneToHashKeyList := make(map[string][]string)
	for _, child := range res {
		fullPath := c.servPath() + common.Slash + child
		data, _, err := c.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return err
		}
		val := new(common.RegisterServiceInfo)
		err = json.Unmarshal(data, &val)
		if err != nil {
			continue
		}
		val.ServPath = fullPath
		var keyList []string
		if _, ok := servLaneToHashKeyList[val.Lane]; ok {
			keyList = servLaneToHashKeyList[val.Lane]
		} else {
			keyList = make([]string, 0, common.ServWeight)
		}
		for i := 0; i < common.ServWeight; i++ {
			keyList = append(keyList, strings.Join([]string{val.ServPath, strconv.FormatInt(int64(i), 10)}, common.HashKey))
		}
		servList = append(servList, val)
		servLaneToHashKeyList[val.Lane] = keyList
	}

	servHash := make(map[string]*consistent.Consistent)
	for lane, keyList := range servLaneToHashKeyList {
		hash := consistent.NewConsistentWithServKeys(keyList)
		if hash == nil {
			continue
		}
		servHash[lane] = hash
	}

	go c.resetPool(c.diffServAndResetClientPool(servList))
	c.updateServ(servList, servHash)
	return nil
}

func (c *Client) diffServAndResetClientPool(servList []*common.RegisterServiceInfo) []string {
	newIpMap, diffIpList := make(map[string]bool, len(servList)), make([]string, 0, len(servList))
	for _, serv := range servList {
		for _, s := range serv.ServList {
			if s.Type == common.Rpc {
				newIpMap[s.Addr] = true
			}
		}
	}

	c.servMu.Lock()
	defer c.servMu.Unlock()

	for _, serv := range servList {
		for _, s := range serv.ServList {
			if s.Type == common.Rpc && !newIpMap[s.Addr] {
				diffIpList = append(diffIpList, s.Addr)
			}
		}
	}
	return diffIpList
}

func (c *Client) updateServ(servList []*common.RegisterServiceInfo, servHash map[string]*consistent.Consistent) {
	c.servMu.Lock()
	defer c.servMu.Unlock()

	c.servList = servList
	c.servLaneHash = servHash
}
