package common

type Client interface {
	Watch() chan Event
}
