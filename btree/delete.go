package btree

func (leftNode *LeafNode) Redistribute(rightNode *LeafNode, parent *InternalNode) {
	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	keys := append(leftNode.Keys, rightNode.Keys...)
	vals := append(leftNode.Values, rightNode.Values...)

	leftNode.Keys = keys[:len(keys)/2]
	leftNode.Values = vals[:len(vals)/2]
	rightNode.Keys = keys[len(keys)/2:]
	rightNode.Values = vals[len(vals)/2:]

	parent.Keys[i] = rightNode.Keys[0]
}

func (leftNode *InternalNode) Redistribute(rightNode *InternalNode, parent *InternalNode) {
	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	keys := append(leftNode.Keys, parent.Keys[i])
	keys = append(keys, rightNode.Keys...)

	vals := append(leftNode.Pointers, rightNode.Pointers...)

	numKeysInLeft := len(keys) / 2
	leftNode.Keys = keys[:numKeysInLeft]
	leftNode.Pointers = vals[:1+numKeysInLeft]
	rightNode.Keys = keys[numKeysInLeft+1:]
	rightNode.Pointers = vals[numKeysInLeft+1:]

	parent.Keys[i] = keys[numKeysInLeft]
}

func (leftNode *LeafNode) MergeNodes(rightNode *LeafNode, parent *InternalNode) {
	var i int
	for i := 0; parent.Pointers[i] != leftNode; i++ {
	}

	keys := append(leftNode.Keys, rightNode.Keys...)
	parent.DeleteAt(i)
	vals := append(leftNode.Values, rightNode.Values...)

	leftNode.Keys = keys
	leftNode.Values = vals

	// delete at shifts to left by one
	parent.Pointers[i] = leftNode

	leftNode.Right = rightNode.Right
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
	n.Keys = append(n.Keys[:index], n.Keys[index+1:]...)
	n.Values = append(n.Values[:index], n.Values[index+1:]...)
}

func (tree *BTree) Delete(key Key) bool {
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.Root.findAndGetStack(key, stack)
	if i == nil {
		return false
	}

	leafNode := stack[len(stack)-1].Node.(*LeafNode)
	index, _ := leafNode.Keys.find(key)
	leafNode.DeleteAt(index)
	stack = stack[:len(stack)-1]
	top := stack[len(stack)-1].Node.(*InternalNode)

	if len(leafNode.Values) < (tree.degree)/2 {
		// should merge or redistribute
		if leafNode.Right != nil && len(leafNode.Right.Keys) >= ((tree.degree)/2)+1 {
			leafNode.Redistribute(leafNode.Right, top)
			return true
		} else if leafNode.Left != nil && len(leafNode.Left.Keys) >= ((tree.degree)/2)+1 {
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
				if rightSibling != nil && len(rightSibling.Pointers) > (tree.degree+1)/2 {
					top.Redistribute(rightSibling, parent)
					return true
				} else if leftSibling != nil && len(leftSibling.Pointers) > (tree.degree+1)/2 {
					leftSibling.Redistribute(top, parent)
					return true
				}

				// if redistribution is not valid merge
				if rightSibling != nil {
					merged = top.MergeNodes(rightSibling, parent)
				} else {
					if leftSibling == nil {
						panic("Both siblings are null for an internal Node! This should not be possible.")
					}
					merged = leftSibling.MergeNodes(top, parent)
				}
				if parent == tree.Root && len(parent.Keys) == 0 {
					tree.Root = merged
				}
			}
		}
	}
	return true
}
