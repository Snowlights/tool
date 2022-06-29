package http

import (
	"fmt"
	"testing"
	"github.com/Snowlights/tool/vservice/common"
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
