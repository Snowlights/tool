package server

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"vtool/server/common"
	"vtool/server/consul"
	"vtool/server/etcd"
	"vtool/server/zk"
	"vtool/vlog"
)

func RegisterService(ctx context.Context, config *common.RegisterConfig) error {

	// todo make group to os.env config
	path := config.Group + common.Slash + config.ServName

	var engine common.Register
	switch config.RegistrationType {
	case common.ETCD:
		engine = etcd.DefaultEtcdInstance
	case common.ZOOKEEPER:
		engine = zk.DefaultZkInstance
	case common.Consul:
		// only for metric collection
		engine = consul.DefaultConsulInstance
	default:
		return common.UnSupportedRegistrationType
	}

	return retryRegister(ctx, engine, config.RegistrationType, path, config.ServAddr)
}

func retryRegister(ctx context.Context, engine common.Register, registrationType common.RegistrationType, path, servAddr string) error {

	rand.Seed(time.Now().Unix())
	retry := 0
	for retry = 0; ; retry++ {
		id, err := calculateCurrentServID(ctx, registrationType, path)
		if err != nil {
			retry++
			time.Sleep(common.DefaultTTl)
			continue
		}
		servPath := path + common.Slash + id

		err = engine.Register(ctx, servPath, servAddr, common.DefaultTTl)
		if err == nil {
			vlog.InfoF(ctx, servPath, servAddr, "register success")
			return nil
		} else {
			vlog.ErrorF(ctx, servPath, servAddr, "register failed error is %s", err.Error())
		}
		retry++
		time.Sleep(common.DefaultTTl)
	}

}

func calculateCurrentServID(ctx context.Context, registrationType common.RegistrationType, path string) (string, error) {
	fun := "calculateCurrentServID --> "
	idList := make([]int, 0)
	switch registrationType {
	case common.ETCD:
		nodes, err := etcd.DefaultEtcdInstance.GetNode(ctx, path)
		if err != nil {
			return "", err
		}
		for _, n := range nodes {
			id := n.Key()[strings.LastIndex(n.Key(), common.Slash)+1:]
			idInt, err := strconv.Atoi(id)
			if err != nil || idInt < 0 {
				vlog.ErrorF(ctx, "%s id error key:%s", fun, n.Key())
			} else {
				idList = append(idList, idInt)
			}
		}
	case common.ZOOKEEPER:
		nodes, err := zk.DefaultZkInstance.GetNode(ctx, path)
		if err != nil {
			return "", err
		}
		for _, n := range nodes {
			id := n.Key()[strings.LastIndex(n.Key(), common.Slash)+1:]
			idInt, err := strconv.Atoi(id)
			if err != nil || idInt < 0 {
				vlog.ErrorF(ctx, "%s id error key:%s", fun, n.Key())
			} else {
				idList = append(idList, idInt)
			}
		}
	default:
		return common.DefaultID, common.UnSupportedRegistrationType
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
