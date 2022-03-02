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

	//for {
	//	msg := <-watchChan
	//	fmt.Println(fmt.Sprintf("%+v", msg))
	//	if len(msg.Events) > 0 {
	//		for _, e := range msg.Events {
	//			fmt.Println(fmt.Sprintf("%+v", e))
	//		}
	//	}
	//
	//	time.Sleep(time.Second * 2)
	//}

}
