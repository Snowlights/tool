package vkafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
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

func (m *Writer) Close() error {
	return m.Writer.Close()
}

func (m *Writer) WriteMsg(ctx context.Context, k string, v interface{}) error {
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = m.WriteMessages(ctx, kafka.Message{
		Key:   []byte(k),
		Value: msg,
	})

	return err
}
