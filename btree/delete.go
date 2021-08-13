package btree

func (n *InternalNode) MergeNode(rightNode *InternalNode, middlePointer node) (mergedNode *InternalNode) {
	n.Keys = append(n.Keys, rightNode.Keys...)
	n.Pointers[len(n.Pointers)-1] = middlePointer
	n.Pointers = append(n.Pointers, rightNode.Pointers...)

	return n
}

func (leftNode *LeafNode) Redistribute(rightNode *LeafNode, parent *InternalNode) {
	keys := append(leftNode.keys, rightNode.keys...)
	vals := append(leftNode.values, rightNode.values)

	leftNode.keys = keys[:len(keys)/2]
	leftNode.values = vals[:len(vals)/2]
	rightNode.keys = keys[len(keys)/2:]
	rightNode.values = vals[len(vals)/2:]

	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	parent.Keys[i] = rightNode.keys[0]
}

func (leftNode *InternalNode) Redistribute(rightNode *InternalNode, parent *InternalNode) {
	keys := append(leftNode.Keys, rightNode.Keys...)
	vals := append(leftNode.Pointers, rightNode.Pointers...)

	leftNode.Keys = keys[:len(keys)/2]
	leftNode.Pointers = vals[:len(vals)/2]
	rightNode.Keys = keys[len(keys)/2:]
	rightNode.Pointers = vals[len(vals)/2:]

	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	parent.Keys[i] = rightNode.Keys[0]
}

func (leftNode *LeafNode) MergeNodes(rightNode *LeafNode, parent *InternalNode) {
	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	keys := append(leftNode.keys, rightNode.keys...)
	parent.DeleteAt(i)
	vals := append(leftNode.values, rightNode.values...)

	leftNode.keys = keys
	leftNode.values = vals

	// delete at shifts to left by one
	parent.Pointers[i] = leftNode
}

func (leftNode *InternalNode) MergeNodes(rightNode *InternalNode, parent *InternalNode) (merged *InternalNode) {
	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	keys := append(leftNode.Keys, parent.Keys[i])
	keys = append(keys, rightNode.Keys...)
	parent.DeleteAt(i)
	pointers := append(leftNode.Pointers, rightNode.Pointers...)

	leftNode.Keys = keys
	leftNode.Pointers = pointers

	// delete at shifts to left by one
	parent.Pointers[i] = leftNode

	return leftNode
}

func (n *InternalNode) DeleteAt(index int) {
	n.Keys = append(n.Keys[:index], n.Keys[index+1:]...)
	n.Pointers = append(n.Pointers[:index], n.Pointers[index+1:]...)
}

func (n *LeafNode) DeleteAt(index int) {
	n.keys = append(n.keys[:index], n.keys[index+1:]...)
	n.values = append(n.values[:index], n.values[index+1:]...)
}

func (tree *BTree) Delete(key Key) bool {
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.Root.findAndGetStack(key, stack)
	if i == nil {
		return false
	}

	leafNode := stack[len(stack)-1].Node.(*LeafNode)
	index, _ := leafNode.keys.find(key)
	leafNode.DeleteAt(index)
	stack = stack[:len(stack)-1]
	top := stack[len(stack)-1].Node.(*InternalNode)

	if len(leafNode.values) < (tree.degree)/2 {
		// should merge or redistribute
		if leafNode.Right != nil && len(leafNode.Right.keys) >= ((tree.degree)/2)+1 {
			leafNode.Redistribute(leafNode.Right, top)
			return true
		} else if leafNode.Left != nil && len(leafNode.Left.keys) >= ((tree.degree)/2)+1 {
			leafNode.Left.Redistribute(leafNode, top)
			return true
		} else {
			if leafNode.Right != nil {
				leafNode.MergeNodes(leafNode.Right, top)
			} else {
				leafNode.Left.MergeNodes(leafNode, top)
			}
		}

		for len(stack) > 0 {
			top := stack[len(stack)-1].Node.(*InternalNode)
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				// if no parent left in stack it is done
				return true
			}
			parent := stack[len(stack)-1].Node.(*InternalNode)
			index, _ = top.Keys.find(key)

			if len(top.Keys) < (tree.degree)/2 {
				// get siblings
				indexAtParent, _ := parent.Keys.find(key)
				var rightSibling, leftSibling, merged *InternalNode
				if indexAtParent > 0 {
					leftSibling = parent.Pointers[indexAtParent-1].(*InternalNode)
				}
				if indexAtParent+1 < len(parent.Pointers) {
					rightSibling = parent.Pointers[indexAtParent+1].(*InternalNode)
				}

				//try redistribute
				if rightSibling != nil && len(rightSibling.Keys) > (tree.degree+1)/2 {
					top.Redistribute(rightSibling, parent)
					return true
				} else if leftSibling != nil && len(leftSibling.Keys) > (tree.degree+1)/2 {
					leftSibling.Redistribute(top, parent)
					return true
				}

				// if redistribution is not valid merge
				if rightSibling != nil {
					merged = top.MergeNodes(rightSibling, parent)
				} else {
					if leftSibling == nil {
						panic("Both siblings are null for an internal node! This should not be possible.")
					}
					merged = leftSibling.MergeNodes(top, parent)
				}
				if top == tree.Root && len(top.Keys) == 0 {
					tree.Root = merged
				}
			}
		}
	}

	return true
}
