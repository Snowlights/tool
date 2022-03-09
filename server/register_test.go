package server

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
	"vtool/server/common"
)

func TestRegisterService(t *testing.T) {
	for i := 0; i < 3; i++ {
		err := RegisterService(context.Background(), &common.RegisterConfig{
			RegistrationType: common.ZOOKEEPER,
			ServName:         "censor",
			ServAddr:         "127.0.0.1:" + strconv.FormatInt(int64(i), 10),
			Group:            "/group/base",
		})
		if err != nil {
			fmt.Println(err)
		}

		//time.Sleep(time.Second)
	}

	time.Sleep(time.Hour)

}
