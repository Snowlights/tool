package vkafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"vtool/vlog"
)

const (
	// REBALANCE_IN_PROGRESS
	ErrorMsgRebalanceInProgress = "Rebalance In Progress"
)

type Kafka struct {
	brokers []string
	topic   string
}

type infoLogger struct {
}

func getInfoLogger() *infoLogger {
	return &infoLogger{}
}

func (i *infoLogger) Printf(format string, v ...interface{}) {
	errMsg := fmt.Sprintf(format, v...)
	if strings.Contains(errMsg, ErrorMsgRebalanceInProgress) {
		log.Println(context.Background(), errMsg)
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
	if strings.Contains(errMsg, ErrorMsgRebalanceInProgress) {
		log.Println(context.Background(), errMsg)
		return
	}
	vlog.ErrorF(context.Background(), errMsg)
}
