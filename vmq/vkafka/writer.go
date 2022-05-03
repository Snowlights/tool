package vkafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	*kafka.Writer
	topic string
}

func NewKafkaWriter(ctx context.Context, brokers []string, topic string) *Writer {
	config := kafka.WriterConfig{
		Brokers:     brokers,
		Topic:       topic,
		Balancer:    &kafka.Hash{},
		Logger:      getInfoLogger(),
		ErrorLogger: getErrorLogger(),
	}
	writer := kafka.NewWriter(config)
	kafkaWriter := &Writer{
		Writer: writer,
		topic:  topic,
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
