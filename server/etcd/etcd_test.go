package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {

	err := DefaulEtcdInstance.Register(context.Background(), "/group/base/censor", "127.0.0.1:9909", time.Second*10)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Minute)
}
