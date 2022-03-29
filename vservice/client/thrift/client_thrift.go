package thrift

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"strconv"
	"sync"
	"time"
	"vtool/parse"
	"vtool/vconfig"
	"vtool/vlog"
	clientCommon "vtool/vservice/client/common"
	"vtool/vservice/client/pool"
	"vtool/vservice/common"
)

type ThriftClient struct {
	client common.Client

	serviceClient func(thrift.TTransport, thrift.TProtocolFactory) interface{}

	clientPool *pool.ClientPool

	center vconfig.Center

	mu           sync.RWMutex
	clientConfig *vconfig.ClientConfig
}

func NewThriftClient(client common.Client, servCli func(thrift.TTransport, thrift.TProtocolFactory) interface{}) common.RpcClient {

	tc := &ThriftClient{
		client:        client,
		serviceClient: servCli,
	}

	err := tc.initCenter()
	if err != nil {
		vlog.ErrorF(context.Background(), "initCenter error: %v", err)
	}

	tc.center.AddListener(&common.ClientListener{Change: tc.reload})

	tc.clientPool = pool.NewClientPool(&pool.ClientPoolConfig{
		ServiceName: client.ServName(),
		Idle:        pool.DefaultIdle,
		Active:      pool.DefaultMaxActive,
		IdleTimeout: pool.DefaultIdleTimeout,
		Wait:        true,
		WaitTimeOut: time.Second * 3,
		StatTime:    pool.DefaultStatTime,
	}, tc.newConn)
	tc.client.AddPoolHandler(tc.deleteAddrHandler)
	return tc
}

func (t *ThriftClient) Rpc(args *common.ClientCallerArgs, fnRpc func(interface{}) error) error {

	if len(args.HashKey) == 0 {
		args.HashKey = clientCommon.NewHashKey()
	}

	serv, ok := t.client.GetServAddr(args.Lane, common.Thrift, args.HashKey)
	if !ok {
		return fmt.Errorf("%s caller args is %+v", common.NotFoundServInfo, args)
	}
	if serv.Type != common.Rpc {
		return fmt.Errorf("%s serv info is %+v, caller args is %+v", common.NotFoundServEngine, serv, args)
	}

	return t.do(context.TODO(), serv, fnRpc)
}

func (t *ThriftClient) do(ctx context.Context, serv *common.ServiceInfo, fnRpc func(interface{}) error) error {

	var err error
	retry := 3
	for ; retry >= 0; retry-- {
		err = t.rpc(ctx, serv, fnRpc)
		if err == nil {
			return nil
		}
	}
	return err
}

func (t *ThriftClient) rpc(ctx context.Context, serv *common.ServiceInfo, fnRpc func(interface{}) error) error {
	conn, err := t.clientPool.Get(ctx, serv)
	if err != nil {
		return err
	}

	return fnRpc(conn.GetConn())
}

func (t *ThriftClient) updateConfig(cfg *vconfig.ClientConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.clientConfig = cfg
}

func (t *ThriftClient) getConfig() *vconfig.ClientConfig {
	t.mu.RLock()
	defer t.mu.RUnlock()

	cfg := t.clientConfig
	return cfg
}

func (t *ThriftClient) reload() {

	cfg := new(vconfig.ClientConfig)
	err := t.center.UnmarshalWithNameSpace(vconfig.Client, parse.PropertiesTagName, cfg)
	if err != nil {
		return
	}
	go t.updateConfig(cfg)
	go t.resetPoolConfig(cfg)
}

func (t *ThriftClient) initCenter() error {
	cfg, err := t.parseConfigEnv()
	if err != nil {
		return err
	}

	center, err := vconfig.NewCenter(cfg)
	if err != nil {
		return err
	}

	t.center = center
	return nil
}

func (t *ThriftClient) parseConfigEnv() (*vconfig.CenterConfig, error) {
	centerConfig, err := vconfig.ParseConfigEnv()
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(centerConfig.Port, 10, 64)
	if err != nil {
		return nil, err
	}

	return &vconfig.CenterConfig{
		AppID:            t.client.ServGroup() + common.Slash + t.client.ServName(),
		Cluster:          centerConfig.Cluster,
		Namespace:        []string{vconfig.Client},
		IP:               centerConfig.IP,
		Port:             int(port),
		IsBackupConfig:   false,
		BackupConfigPath: "",
		MustStart:        centerConfig.MustStart,
	}, nil
}

func (t *ThriftClient) resetPoolConfig(cfg *vconfig.ClientConfig) {
	t.clientPool.ResetConnConfig(cfg)
}

func (t *ThriftClient) deleteAddrHandler(addr []string) {
	for _, addr := range addr {
		t.clientPool.Delete(context.Background(), addr)
	}
}

func (t *ThriftClient) newConn(addr string) (common.RpcConn, error) {
	ctx := context.Background()

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		vlog.ErrorF(ctx, "open thriftSocket addr:%s, err:%v", addr, err)
		return nil, err
	}
	useTransport := transportFactory.GetTransport(transport)

	if err := useTransport.Open(); err != nil {
		vlog.ErrorF(ctx, "open addr:%s err:%v", addr, err)
		return nil, err
	}

	return &ClientConn{
		thriftSocket:  transport,
		transport:     useTransport,
		serviceClient: t.serviceClient(useTransport, protocolFactory),
	}, nil
}
