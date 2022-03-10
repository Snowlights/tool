package zk

type Node struct {
	key string
	val string
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) Val() string {
	return n.val
}
