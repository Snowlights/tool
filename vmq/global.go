package vmq

func ReadMsgWithTopic(cluster, topic string, ov interface{}) error {

	return nil
}

func ReadMsgWithTopicAndPartition(cluster, topic string, partition int, ov interface{}) error {

	return nil
}

func FetchMsgWithTopic(cluster, topic string, partition int, ov interface{}) error {

	return nil
}

func FetchMsgWithTopicAndPartition(cluster, topic string, partition int, ov interface{}) error {

	return nil
}

func WriteMsgWithTopic(cluster, topic string, msgs ...Message) error {
	return nil
}
