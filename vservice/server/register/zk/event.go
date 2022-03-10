package zk

import "vtool/vservice/common"

type Event struct {
	EventType common.EventType
}

func (e Event) Event() common.EventType {
	return e.EventType
}
