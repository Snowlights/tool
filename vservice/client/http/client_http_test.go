package http

import (
	"fmt"
	"github.com/Snowlights/tool/vservice/common"
	"testing"
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
