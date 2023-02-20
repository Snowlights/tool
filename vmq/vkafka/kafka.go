package vkafka

import (
	"context"
	"fmt"
	"github.com/Snowlights/tool/vlog"
	"log"
	"strings"
	"time"
)

const (
	// REBALANCE_IN_PROGRESS
	ErrorMsgReBalanceInProgress = "Rebalance In Progress"
)

type KafkaReaderConf struct {
	Brokers        []string
	Topic          string
	Group          string
	Partition      int
	CommitInterval time.Duration
	MinByte        int
	MaxByte        int
	StartOffset    int64
}

type KafkaWriterConf struct {
	Brokers []string
	Topic   string
}

type infoLogger struct{}

func getInfoLogger() *infoLogger {
	return &infoLogger{}
}

func (i *infoLogger) Printf(format string, v ...interface{}) {
	errMsg := fmt.Sprintf(format, v...)
	if strings.Contains(errMsg, ErrorMsgReBalanceInProgress) {
		return
	}
	vlog.ErrorF(context.Background(), errMsg)
}

type errorLogger struct {
}

func getErrorLogger() *errorLogger {
	return &errorLogger{}
}

func (m *errorLogger) Printf(format string, v ...interface{}) {
	errMsg := fmt.Sprintf(format, v...)
	if strings.Contains(errMsg, ErrorMsgReBalanceInProgress) {
		log.Println(context.Background(), errMsg)
		return
	}
	vlog.ErrorF(context.Background(), errMsg)
}
