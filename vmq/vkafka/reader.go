package vkafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Reader struct {
	*kafka.Reader
	conf *KafkaReaderConf
}

func NewKafkaReader(insConf *KafkaReaderConf) *Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        insConf.Brokers,
		Topic:          insConf.Topic,
		GroupID:        insConf.Group,
		Partition:      insConf.Partition,
		MinBytes:       insConf.MinByte,
		MaxBytes:       insConf.MaxByte,
		CommitInterval: insConf.CommitInterval,
		StartOffset:    insConf.StartOffset,
		ErrorLogger:    getErrorLogger(),
	})

	kafkaReader := &Reader{
		Reader: reader,
		conf:   insConf,
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

func (m *Reader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return m.Reader.FetchMessage(ctx)
}

func (m *Reader) Commit(ctx context.Context, msg ...kafka.Message) error {
	return m.Reader.CommitMessages(ctx, msg...)
}

func (m *Reader) Close() error {
	return m.Reader.Close()
}
