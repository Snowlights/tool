package vkafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"strings"
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

func (r *Reader) ReadMsg(ctx context.Context) (kafka.Message, error) {
	msg, err := r.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}

	return msg, nil
}

func (r *Reader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return r.Reader.FetchMessage(ctx)
}

func (r *Reader) Commit(ctx context.Context, msg ...kafka.Message) error {
	return r.Reader.CommitMessages(ctx, msg...)
}

func (r *Reader) Cluster() string {
	return strings.Join(r.conf.Brokers, ",")
}

func (r *Reader) Close() error {
	return r.Reader.Close()
}
