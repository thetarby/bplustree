package btree

func (n *InternalNode) truncate(index int) {
	n.Keys = n.Keys[:index]
	n.Pointers = n.Pointers[:index+1]
}
