package server

import (
	"context"
	"time"
)

type RegistrationType int64

const (
	ETCD      RegistrationType = 1
	ZOOKEEPER RegistrationType = 2
)

type RegisterConfig struct {
	RegistrationType RegistrationType
}

type ETCDConfig struct{}

type Register interface {
	// 获取指定key的值
	Get(ctx context.Context, path string) (string, error)
	// 获取指定key对应的节点，会将节点及子节点返回
	GetNode(ctx context.Context, path string) (Node, error)
	// 设置key值，并指定过期时间
	SetTtl(ctx context.Context, path, val string, ttl time.Duration) error
	// 刷新节点的过期时间，不更新值
	RefreshTtl(ctx context.Context, path string, ttl time.Duration) error
	// 执行注册，注册后会一直维持心跳。调用注册时会将当前值设置到节点上
	Register(ctx context.Context, path, val string, ttl time.Duration) error
	// 监听一个节点
	Watch(ctx context.Context, path string, handler func(Handler))
}

// 节点信息
type Node interface {
	Key() string
	Val() string
	Ttl() int64
	Children() []Node
}

// 监听回调
type Handler interface {
	Node() Node
}