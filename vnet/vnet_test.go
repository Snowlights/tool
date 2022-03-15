package vnet

import (
	"context"
	"fmt"
	"testing"
)

func TestGetLocalIp(t *testing.T) {

	fmt.Println(ListenServAddr(context.Background(), ":"))

}
