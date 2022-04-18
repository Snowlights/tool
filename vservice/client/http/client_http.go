package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"vtool/vlog"
	clientCommon "vtool/vservice/client/common"
	"vtool/vservice/common"
	"vtool/vtrace"
)

type HttpClient struct {
	client common.Client
}

func NewHttpClient(cliConfig *common.ClientConfig) (*HttpClient, error) {
	c := &HttpClient{}

	cli, err := clientCommon.NewClientWithClientConfig(cliConfig)
	if err != nil {
		return nil, err
	}
	c.client = cli
	return c, nil
}

func (hc HttpClient) Do(args *common.ClientCallerArgs, option interface{}) (interface{}, error) {

	if option == nil {
		return nil, fmt.Errorf("%s, caller option is nil", common.NeedHttpCallerOption)
	}

	opt, ok := option.(*common.HttpCallerOptions)
	if !ok {
		return nil, fmt.Errorf("%s, caller option is %+v", common.NeedHttpCallerOption, option)
	}

	if len(args.HashKey) == 0 {
		args.HashKey = clientCommon.NewHashKey()
	}

	serv, ok := hc.client.GetServAddr(args.Lane, common.HTTP, args.HashKey)
	if !ok {
		return nil, fmt.Errorf("%s caller args is %+v", common.NotFoundServInfo, args)
	}

	if serv.Type != common.Gin {
		return nil, fmt.Errorf("%s serv info is %+v, caller args is %+v", common.NotFoundServEngine, serv, args)
	}

	return hc.do(serv, opt, args.TimeOut)
}

func (hc HttpClient) do(serv *common.ServiceInfo, option *common.HttpCallerOptions, timeout time.Duration) (interface{}, error) {

	if timeout == 0 {
		timeout = common.DefaultMaxTimeOut
	}

	url := common.HttpPrefix + serv.Addr + option.API

	request, err := http.NewRequest(option.Method, url, bytes.NewReader(option.Body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Connection", "Keep-Alive")
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(timeout, func() {
		cancel()
	})
	err = vtrace.TraceHTTPRequest(ctx, request)
	if err != nil {
		vlog.ErrorF(ctx, "HttpClient.do trace http request error: %v", err)
		err = nil
	}
	request = request.WithContext(ctx)

	response, err := common.DefaultHttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status is %d, serv is %+v, options is %+v", response.StatusCode, serv, option)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
