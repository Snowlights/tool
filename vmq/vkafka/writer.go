package vkafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"strings"
)

type Writer struct {
	*kafka.Writer
	conf *KafkaWriterConf
}

func NewKafkaWriter(insConf *KafkaWriterConf) *Writer {
	writer := &kafka.Writer{
		Addr:        kafka.TCP(insConf.Brokers...),
		Topic:       insConf.Topic,
		Balancer:    &kafka.Hash{},
		Logger:      getInfoLogger(),
		ErrorLogger: getErrorLogger(),
	}

	kafkaWriter := &Writer{
		Writer: writer,
		conf:   insConf,
	}

	return kafkaWriter
}

func (w *Writer) Cluster() string {
	return strings.Join(w.conf.Brokers, ",")
}

func (w *Writer) Topic() string {
	return w.conf.Topic
}

func (w *Writer) Close() error {
	return w.Writer.Close()
}

func (w *Writer) WriteMsg(ctx context.Context, k string, v interface{}) error {
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(k),
		Value: msg,
	})

	return err
}

func (w *Writer) WriteMsgs(ctx context.Context, msgs ...kafka.Message) error {
	return w.WriteMessages(ctx, msgs...)
}
