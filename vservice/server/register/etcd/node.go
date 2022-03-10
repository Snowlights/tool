package etcd

type Node struct {
	key string
	val string

	Lease int64
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) Val() string {
	return n.val
}
