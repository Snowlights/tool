package vmq

import (
	"context"
	"github.com/segmentio/kafka-go"
)

func ReadMsgWithTopicGroup(ctx context.Context, cluster, topic, group string, ov interface{}) error {
	return defaultManager.ReadMsg(ctx, cluster, topic, group, 0, ov)
}

func ReadMsgWithTopicAndPartition(cluster, topic string, partition int, ov interface{}) error {
	return defaultManager.ReadMsg(context.Background(), cluster, topic, "", partition, ov)
}

func FetchMsgWithTopicGroup(cluster, topic, group string, ov interface{}) (Handler, error) {
	return defaultManager.FetchMsg(context.Background(), cluster, topic, group, 0, ov)
}

func FetchMsgWithTopicAndPartition(cluster, topic string, partition int, ov interface{}) (Handler, error) {
	return defaultManager.FetchMsg(context.Background(), cluster, topic, "", partition, ov)
}

func WriteMsgWithTopic(cluster, topic string, k string, v interface{}) error {

	return nil
}

func WriteMsgsWithTopic(cluster, topic string, msgs ...kafka.Message) error {
	return nil
}
