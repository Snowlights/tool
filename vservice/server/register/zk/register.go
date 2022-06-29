package zk

import (
	"context"
	"encoding/json"
	"github.com/samuel/go-zookeeper/zk"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vservice/common"
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

func (c *Register) Register(ctx context.Context, path, val string, ttl time.Duration) (string, error) {
	if err := c.ensureAllPathExit(path); err != nil {
		return "", err
	}

	rand.Seed(time.Now().Unix())
	retry := 0
	for retry = 0; ; retry++ {
		id, err := c.calculateCurrentServID(ctx, path)
		if err != nil {
			retry++
			time.Sleep(common.DefaultTTl)
			continue
		}
		servPath := path + common.Slash + id

		err = c.register(ctx, servPath, val, ttl)
		if err == nil {
			vlog.Info(ctx, servPath, val, "register success")
			return id, nil
		} else {
			vlog.ErrorF(ctx, servPath, val, "register failed error is %s", err.Error())
		}
		retry++
		time.Sleep(common.DefaultTTl)
	}
}

func (c *Register) register(ctx context.Context, path, val string, ttl time.Duration) error {
	_, err := c.conn.Create(path, []byte(val), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}

func (c *Register) calculateCurrentServID(ctx context.Context, path string) (string, error) {
	fun := "etcd.Register.calculateCurrentServID --> "

	idList := make([]int, 0)
	nodes, err := c.GetNode(ctx, path)
	if err != nil {
		return "", err
	}
	for _, n := range nodes {
		id := n.ServPath[strings.LastIndex(n.ServPath, common.Slash)+1:]
		idInt, err := strconv.Atoi(id)
		if err != nil || idInt < 0 {
			vlog.ErrorF(ctx, "%s id error key:%s", fun, n.ServPath)
		} else {
			idList = append(idList, idInt)
		}
	}

	sort.Ints(idList)
	idRes := 0
	for _, id := range idList {
		if idRes == id {
			idRes++
		} else {
			break
		}
	}
	return strconv.FormatInt(int64(idRes), 10), nil
}

// It is not recommended to use the new method here. It is
// recommended to use a special group for maintenance.
//It is maintained uniformly within the company through work orders or other methods
func (c *Register) ensureAllPathExit(path string) error {
	parts := strings.Split(path, common.Slash)

	if len(parts) == 1 {
		return c.ensurePath(path)
	}

	i := 2
	for i < len(parts) {
		err := c.ensurePath(strings.Join(parts[:i], common.Slash))
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

func (c *Register) UnRegister(ctx context.Context, path string) error {
	c.conn.Delete(path, 0)
	return nil
}

func (c *Register) Get(ctx context.Context, path string) (string, error) {

	res, _, err := c.conn.Get(path)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (c *Register) GetNode(ctx context.Context, path string) ([]*common.RegisterServiceInfo, error) {
	c.ensurePath(path)
	res, _, err := c.conn.Children(path)
	if err != nil {
		return nil, err
	}

	nodeList := make([]*common.RegisterServiceInfo, 0, len(res))
	for _, child := range res {
		fullPath := path + common.Slash + child
		data, _, err := c.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return nil, err
		}
		val := new(common.RegisterServiceInfo)
		err = json.Unmarshal(data, &val)
		if err != nil {
			continue
		}
		val.ServPath = string(fullPath)
		nodeList = append(nodeList, val)
	}

	return nodeList, nil
}

// TTL node is added in version 3.5.5
// In version 3.5.8, zookeeper does not support nodes with expiration time by default.
// However, in 3.6.3, the expiration time is directly supported by default.
// The referenced SDK has no - t option and will not be demonstrated here.
// It is generally used as a distributed lock. In this scenario,
// the - e parameter is used to meet the requirements of
// the server registration and discovery scenario
func (c *Register) RefreshTtl(ctx context.Context, path, val string, ttl time.Duration) error {

	return nil
}
