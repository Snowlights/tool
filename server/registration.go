package server

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"vtool/server/etcd"
	"vtool/server/zk"
	"vtool/vlog"
	"vtool/vnet"
)

func RegisterService(ctx context.Context, config *RegisterConfig) error {

	// todo make group to os.env config
	path := config.Group + zk.Slash + config.ServName
	servAddr, err := vnet.GetServAddr(config.ServAddr)
	if err != nil {
		return err
	}

	var engine Register
	switch config.RegistrationType {
	case ETCD:
		engine = etcd.DefaultEtcdInstance
	case ZOOKEEPER:
		engine = zk.DefaultZkInstance
	default:
		return UnSupportedRegistrationType
	}

	return retryRegister(ctx, engine, config.RegistrationType, path, servAddr)
}

func retryRegister(ctx context.Context, engine Register, registrationType RegistrationType, path, servAddr string) error {

	rand.Seed(time.Now().Unix())
	retry := 0
	for retry < retryTime {
		s := rand.Intn(100)
		time.Sleep(time.Millisecond * time.Duration(s))

		id, err := calculateCurrentServID(ctx, registrationType, path)
		if err != nil {
			retry++
			continue
		}
		servPath := path + zk.Slash + id

		err = engine.Register(ctx, servPath, servAddr, defaultTTl)
		if err == nil {
			return nil
		}
		retry++
	}

	return RegisterFailed
}

func calculateCurrentServID(ctx context.Context, registrationType RegistrationType, path string) (string, error) {
	fun := "calculateCurrentServID --> "
	idList := make([]int, 0)
	switch registrationType {
	case ETCD:
		nodes, err := etcd.DefaultEtcdInstance.GetNode(ctx, path)
		if err != nil {
			return "", err
		}
		for _, n := range nodes {
			id := n.Key()[strings.LastIndex(n.Key(), zk.Slash)+1:]
			idInt, err := strconv.Atoi(id)
			if err != nil || idInt < 0 {
				vlog.ErrorF(ctx, "%s id error key:%s", fun, n.Key())
			} else {
				idList = append(idList, idInt)
			}
		}
	case ZOOKEEPER:
		nodes, err := zk.DefaultZkInstance.GetNode(ctx, path)
		if err != nil {
			return "", err
		}
		for _, n := range nodes {
			id := n.Key()[strings.LastIndex(n.Key(), zk.Slash)+1:]
			idInt, err := strconv.Atoi(id)
			if err != nil || idInt < 0 {
				vlog.ErrorF(ctx, "%s id error key:%s", fun, n.Key())
			} else {
				idList = append(idList, idInt)
			}
		}
	default:
		return _defaultID, UnSupportedRegistrationType
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
