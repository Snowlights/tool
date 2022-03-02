package vnet

import (
	"fmt"
	"testing"
)

func TestGetLocalIp(t *testing.T) {

	fmt.Println(GetServAddr(":4445"))

}
