package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestClient_Register(t *testing.T) {

	err := DefaultEtcdInstance.Register(context.Background(), "/group/base/censor/1", "127.0.0.1:9909", time.Second*10)

	if err != nil {
		fmt.Println(err)
	}
	err = DefaultEtcdInstance.Register(context.Background(), "/group/base/censor/1", "127.0.0.1:9910", time.Second*10)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Minute)
}

func TestClient_Register1(t *testing.T) {

	err := DefaultEtcdInstance.Register(context.Background(), "/group/base/censor/2", "127.0.0.1:9909", time.Second*10)
	if err != nil {
		fmt.Println(err)
	}
	err = DefaultEtcdInstance.Register(context.Background(), "/group/base/censor/3", "127.0.0.1:9910", time.Second*10)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Minute)
}

func TestClient_Get(t *testing.T) {

	base := "/group/base/censor"

	fmt.Println(DefaultEtcdInstance.GetNode(context.Background(), base))
	//watchChan, err := DefaultEtcdInstance.Watch(context.Background(), base)
	//if err != nil {
	//
	//}
	//for {
	//	msg := <-watchChan
	//	fmt.Println(msg)
	//	// todo
	//}

}
