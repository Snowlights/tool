package vcpu

import (
	"fmt"
	"testing"
)

func TestUsage(t *testing.T) {
	fmt.Println(Usage())
}

func TestInfo(t *testing.T) {
	fmt.Println(Info())
}
