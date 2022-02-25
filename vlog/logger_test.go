package vlog

import (
	"context"
	"testing"
	"time"
)

func TestInitLogger(t *testing.T) {

	for {
		Debug(context.Background(), "1", "2", "3")
		time.Sleep(time.Second * 5)
	}

}
