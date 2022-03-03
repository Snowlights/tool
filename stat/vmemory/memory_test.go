package vmemory

import (
	"fmt"
	"testing"
	"time"
)

func TestNewMemory(t *testing.T) {

	for {

		fmt.Println(Virtual())
		fmt.Println(GoroutineNums())
		fmt.Println(HeapObjects())
		fmt.Println(HeapAlloc())
		fmt.Println(GCPause())

		time.Sleep(time.Second * 10)
	}

}
