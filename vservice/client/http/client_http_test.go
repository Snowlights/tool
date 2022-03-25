package http

import (
	"fmt"
	"testing"
	"vtool/vservice/common"
)

func TestHttpClient_Do(t *testing.T) {

	c, _ := NewHttpClient(&common.ClientConfig{
		RegistrationType: 0,
		Cluster:          nil,
		ServGroup:        "",
		ServName:         "",
	})
	fmt.Println(c)
}
