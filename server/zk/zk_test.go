package zk

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewRegister(t *testing.T) {

	err := DefaultZkInstance.Register(context.Background(), "/group/base/censor/service1", "127.0.0.1:9909", time.Second*10)
	if err != nil {
		fmt.Println(err)
	}

	err = DefaultZkInstance.Register(context.Background(), "/group/base/censor/service2", "127.0.0.1:9910", time.Second*10)
	if err != nil {
		fmt.Println(err)
	}

	err = DefaultZkInstance.Register(context.Background(), "/group/base/censor/service3", "127.0.0.1:9911", time.Second*10)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Hour)
}

func TestRegister_Get(t *testing.T) {

	fmt.Println(DefaultZkInstance.Get(context.Background(), "/group/base/censor/service1"))

}

func TestRegister_GetNode(t *testing.T) {
	DefaultZkInstance.GetNode(context.Background(), "/group/base/censor")
}

func TestRegister_Watch(t *testing.T) {
	resChan, eventChan, errChan := DefaultZkInstance.Watch(context.Background(), "/group/base/censor")

	go func() {
		for {
			res := <-resChan
			fmt.Println(res)
		}
	}()

	go func() {
		for {
			res := <-eventChan
			fmt.Println(res)
		}
	}()

	go func() {
		for {
			res := <-errChan
			fmt.Println(res)
		}
	}()

	time.Sleep(time.Hour)
}
