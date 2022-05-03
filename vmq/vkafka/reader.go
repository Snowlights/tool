package vkafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Reader struct {
	*kafka.Reader
	topic string
}

func NewKafkaReader(brokers []string, topic, groupId string, partition, minBytes, maxBytes int) *Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupId,
		Partition:      partition,
		MinBytes:       minBytes,
		MaxBytes:       maxBytes,
		CommitInterval: 0,
		StartOffset:    kafka.LastOffset,
		ErrorLogger:    getErrorLogger(),
	})

	kafkaReader := &Reader{
		Reader: reader,
		topic:  "",
	}

	return kafkaReader
}

func (m *Reader) ReadMsg(ctx context.Context, ov interface{}) error {
	msg, err := m.ReadMessage(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal(msg.Value, ov)
	if err != nil {
		return err
	}

	return nil
}

func (m *Reader) Commit(ctx context.Context, msg ...kafka.Message) error {
	return m.Reader.CommitMessages(ctx, msg...)
}

func (m *Reader) Close() error {
	return m.Reader.Close()
}
