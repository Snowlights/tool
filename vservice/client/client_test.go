package client

import (
	"fmt"
	"testing"
	"vtool/vservice/common"
)

func TestHttpClient_Do(t *testing.T) {

	var cli common.Caller
	c, _ := NewHttpClient(&common.ClientConfig{
		RegistrationType: 0,
		Cluster:          nil,
		ServGroup:        "",
		ServName:         "",
	})

	cli = c

	fmt.Println(cli)
}
