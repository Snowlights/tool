package consul

import (
	"context"
	"fmt"
	"testing"
	"time"
	"vtool/vservice/common"
)

func TestNewConsulServiceRegistry(t *testing.T) {
	fmt.Println(DefaultConsulInstance.Register(context.Background(), "/group/base/censor/1", "127.0.0.1:1", common.DefaultTTl))
	fmt.Println(DefaultConsulInstance.Register(context.Background(), "/group/base/censor/1", "127.0.0.1:2", common.DefaultTTl))

	time.Sleep(time.Hour)
}
