package merkletree

type Node struct {
	k string
	v string
	l *Node
	r *Node
}

type Tree struct {
	r *Node
}
