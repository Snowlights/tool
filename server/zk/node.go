package zk

type Node struct {
	ip   string
	port string

	key string
	val string
}

func (n *Node) Key() string {
	return n.key
}

func (n *Node) Val() string {
	return n.val
}

func (n *Node) Ip() string {
	return n.ip
}

func (n *Node) Port() string {
	return n.port
}