package consul

import (
	"fmt"
	"testing"
)

func TestNewConsulServiceRegistry(t *testing.T) {

	c, err := NewConsulServiceRegistry("127.0.0.1", 8500, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(c)

}
