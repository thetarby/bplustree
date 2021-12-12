package btree

type TreeIterator struct {
	tree    *BTree
	curr    Pointer
	currIdx int
	pager   Pager
}

func (it *TreeIterator) Next() interface{} {
	currNode := it.pager.GetNode(it.curr)
	h := currNode.GetHeader()

	// if there is no element left in node proceed to next node
	if h.KeyLen == int16(it.currIdx) {
		it.pager.Unpin(currNode, false)
		if h.Right == 0 {
			return nil
		}
		it.curr = h.Right
		currNode = it.pager.GetNode(it.curr)
		it.currIdx = 0
	}

	val := currNode.GetValueAt(it.currIdx)
	it.pager.Unpin(currNode, false)
	it.currIdx++
	return val
}

// NewTreeIterator creates an iterator which starts from the smallest key in the tree and iterates through up until
// the largest key.
func NewTreeIterator(tree *BTree, pager Pager) *TreeIterator {
	curr := tree.GetRoot()
	for !curr.IsLeaf() {
		old := curr
		curr = tree.pager.GetNode(curr.GetValueAt(0).(Pointer))
		tree.pager.Unpin(old, false)
	}

	defer tree.pager.Unpin(curr, false)

	return &TreeIterator{
		tree:    tree,
		curr:    curr.GetPageId(),
		currIdx: 0,
		pager:   pager,
	}
}

// NewTreeIteratorWithKey starts the iterator from the given key if it exists otherwise starts from the smallest key
// which is larger than given key.
func NewTreeIteratorWithKey(key Key, tree *BTree, pager Pager) *TreeIterator {
	_, stack := tree.FindAndGetStack(key, Read)
	leaf, idx := stack[len(stack)-1].Node, stack[len(stack)-1].Index

	return &TreeIterator{
		tree:    tree,
		curr:    leaf,
		currIdx: idx,
		pager:   pager,
	}
}
