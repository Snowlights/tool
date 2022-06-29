package etcd

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vservice/common"
)

type Register struct {
	client *clientv3.Client
}

func (c *Register) Register(ctx context.Context, path, val string, ttl time.Duration) (string, error) {

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

	kv := clientv3.NewKV(c.client)
	lease := clientv3.NewLease(c.client)
	leaseRes, err := lease.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return err
	}

	tx := kv.Txn(ctx)
	// Transaction gun lock
	tx.If(clientv3.Compare(clientv3.CreateRevision(path), common.Equals, 0)).
		Then(clientv3.OpPut(path, "", clientv3.WithLease(leaseRes.ID))).
		Else(clientv3.OpGet(path))

	txRes, err := tx.Commit()
	if err != nil {
		return err
	}
	if !txRes.Succeeded {
		return lockFailed
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
				vlog.Info(ctx, strings.Join([]string{path, val, leaseSuccess}, common.Colon))
			}
		}
	}
}

func (c *Register) calculateCurrentServID(ctx context.Context, path string) (string, error) {
	fun := "etcd.Register.calculateCurrentServID --> "

	nodes, err := c.GetNode(ctx, path)
	if err != nil {
		return "", err
	}
	idList := make([]int, 0)

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

func (c *Register) UnRegister(ctx context.Context, path string) error {
	c.client.Delete(ctx, path)
	return nil
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

func (c *Register) GetNode(ctx context.Context, path string) ([]*common.RegisterServiceInfo, error) {

	res, err := c.client.Get(ctx, path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	nodeList := make([]*common.RegisterServiceInfo, 0, len(res.Kvs))
	for _, v := range res.Kvs {
		val := new(common.RegisterServiceInfo)
		err = json.Unmarshal(v.Value, &val)
		if err != nil {
			continue
		}
		val.ServPath = string(v.Key)
		nodeList = append(nodeList, val)
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
