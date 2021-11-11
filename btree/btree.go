package btree

// TODO:
// 1. don't use Keys.find instead define one on Node and use that method.
// 2. define util methods to change state of node such as shift_keys_by_n, shift_pointers_by_n, truncate_at_n
//    and use them in methods like InsertAt and SplitNode
// 3. Put a constraint on Key to make it fix sized maybe? It will solve many problems when we try to persist nodes on a disk
// 4. Use interface methods in delete.go as well.

import (
	"fmt"
)

type BTree struct {
	degree int
	length int
	Root   Pointer
	pager  Pager
}

func NewBtree(degree int) *BTree {
	pager := NoopPager{}
	l := pager.NewLeafNode()

	root := pager.NewInternalNode(l.GetPageId())

	return &BTree{
		degree: degree,
		length: 0,
		Root:   root.GetPageId(),
		pager:  &pager,
	}
}

func NewBtreeWithPager(degree int, pager Pager) *BTree {
	l := pager.NewLeafNode()

	root := pager.NewInternalNode(l.GetPageId())

	return &BTree{
		degree: degree,
		length: 0,
		Root:   root.GetPageId(),
		pager:  pager,
	}
}

func (tree *BTree) GetRoot() Node {
	return tree.pager.GetNode(tree.Root)
}

func (tree *BTree) Insert(key Key, value interface{}) {
	pager := tree.pager
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.GetRoot().findAndGetStack(key, stack)
	if i != nil {
		panic("key already exists")
	}

	var rightNod = value
	var rightKey = key

	for len(stack) > 0 {
		popped := tree.pager.GetNode(stack[len(stack)-1].Node)
		stack = stack[:len(stack)-1]
		i, _ := popped.findKey(key)
		popped.InsertAt(i, rightKey, rightNod)
		//topOfStack.PrintNode()

		if popped.IsOverFlow(tree.degree) {
			rightNod, _, rightKey = popped.SplitNode((tree.degree) / 2)

			if popped.GetPageId() == tree.Root {
				leftNode := popped

				newRoot := pager.NewInternalNode(leftNode.GetPageId())
				newRoot.InsertAt(0, rightKey, rightNod.(Pointer))
				tree.Root = newRoot.GetPageId()
			}
		} else {
			break
		}
	}
}

func (tree *BTree) InsertOrReplace(key Key, value interface{}) (isInserted bool) {
	pager := tree.pager
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.GetRoot().findAndGetStack(key, stack)
	if i != nil {
		// top of stack is the leaf Node
		topOfStack := stack[len(stack)-1]
		leafNode := tree.pager.GetNode(topOfStack.Node)
		leafNode.setValueAt(topOfStack.Index, value)
		return false
	}

	var rightNod = value
	var rightKey = key

	for len(stack) > 0 {
		topOfStack := tree.pager.GetNode(stack[len(stack)-1].Node)
		i, _ := topOfStack.findKey(key)
		topOfStack.InsertAt(i, rightKey, rightNod)

		popped := tree.pager.GetNode(stack[len(stack)-1].Node)
		stack = stack[:len(stack)-1]
		if popped.IsOverFlow(tree.degree) {
			rightNod, _, rightKey = popped.SplitNode((tree.degree) / 2)
			if popped.GetPageId() == tree.Root {
				leftNode := popped

				newRoot := pager.NewInternalNode(leftNode.GetPageId())
				newRoot.InsertAt(0, rightKey, rightNod.(Pointer))
				tree.Root = newRoot.GetPageId()
			}
		} else {
			break
		}
	}

	return true
}

func (tree *BTree) Find(key Key) interface{} {
	res, _ := tree.GetRoot().findAndGetStack(key, []NodeIndexPair{})
	return res
}

func (tree *BTree) Height() int {
	pager := tree.pager
	var currentNode Node = tree.pager.GetNode(tree.Root)
	acc := 0
	for {
		if currentNode.IsLeaf() {
			return acc + 1
		} else {
			currentNode = pager.GetNode(currentNode.GetValueAt(0).(Pointer))
		}
		acc++
	}
}

func (tree BTree) Print() {
	pager := tree.pager
	queue := make([]Pointer, 0, 2)
	queue = append(queue, tree.Root)
	queue = append(queue, 0)
	for i := 0; i < len(queue); i++ {
		node := tree.pager.GetNode(queue[i])
		if node != nil && node.IsLeaf() {
			break
		}
		if node == nil {
			queue = append(queue, 0)
			continue
		}

		pointers := make([]Pointer, 0)
		vals := node.GetValues()
		for _, val := range vals {
			pointers = append(pointers, val.(Pointer))
		}
		queue = append(queue, pointers...)
	}
	for _, n := range queue {
		if n != 0 {
			currNode := pager.GetNode(n)
			currNode.PrintNode()
		} else {
			fmt.Print("\n ### \n")
		}
	}
}

func (tree *BTree) DeleteOld(key Key) bool {
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.GetRoot().findAndGetStack(key, stack)
	if i == nil {
		return false
	}

	leafNode := tree.pager.GetNode(stack[len(stack)-1].Node)
	index, _ := leafNode.findKey(key)
	leafNode.DeleteAt(index)
	stack = stack[:len(stack)-1]
	top := tree.pager.GetNode(stack[len(stack)-1].Node)

	if leafNode.IsUnderFlow(tree.degree) { //len(leafNode.Values) < (tree.degree)/2 {
		// should merge or redistribute
		rightOfLeaf := tree.pager.GetNode(leafNode.GetRight())
		leftOfLeaf := tree.pager.GetNode(leafNode.GetLeft())
		if rightOfLeaf != nil && rightOfLeaf.Keylen() >= ((tree.degree)/2)+1 {
			leafNode.Redistribute(rightOfLeaf, top)
			return true
		} else if leftOfLeaf != nil && leftOfLeaf.Keylen() >= ((tree.degree)/2)+1 {
			leftOfLeaf.Redistribute(leafNode, top)
			return true
		} else {
			if rightOfLeaf != nil {
				leafNode.MergeNodes(rightOfLeaf, top)
			} else if leftOfLeaf != nil {
				leftOfLeaf.MergeNodes(leafNode, top)
			} else {
				// TODO: maybe log here
				return true
			}
		}

		for len(stack) > 0 {
			top := tree.pager.GetNode(stack[len(stack)-1].Node)
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				// if no parent left in stack it is done
				return true
			}
			parent := tree.pager.GetNode(stack[len(stack)-1].Node)
			index, _ = top.findKey(key)

			if top.IsUnderFlow(tree.degree) {
				// get siblings
				indexAtParent, _ := parent.findKey(key)
				var rightSibling, leftSibling, merged Node
				if indexAtParent > 0 {
					leftSibling = tree.pager.GetNode(parent.GetValueAt(indexAtParent - 1).(Pointer)) //leftSibling = parent.Pointers[indexAtParent-1].(*InternalNode)
				}
				if indexAtParent+1 < (parent.Keylen() + 1) { // +1 is the length of pointers
					rightSibling = tree.pager.GetNode(parent.GetValueAt(indexAtParent + 1).(Pointer)) //rightSibling = parent.Pointers[indexAtParent+1].(*InternalNode)
				}

				//try redistribute
				if rightSibling != nil && (rightSibling.Keylen()+1) > (tree.degree+1)/2 {
					top.Redistribute(rightSibling, parent)
					return true
				} else if leftSibling != nil && (leftSibling.Keylen()+1) > (tree.degree+1)/2 {
					leftSibling.Redistribute(top, parent)
					return true
				}

				// if redistribution is not valid merge
				if rightSibling != nil {
					top.MergeNodes(rightSibling, parent)
					merged = top
				} else {
					if leftSibling == nil {
						panic("Both siblings are null for an internal Node! This should not be possible.")
					}
					leftSibling.MergeNodes(top, parent)
					merged = leftSibling
				}
				if parent.GetPageId() == tree.Root && parent.Keylen() == 0 {
					tree.Root = merged.GetPageId()
				}
			}
		}
	}
	return true
}

func (tree *BTree) Delete(key Key) bool {
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.GetRoot().findAndGetStack(key, stack)
	if i == nil {
		return false
	}

	for len(stack) > 0 {
		popped := tree.pager.GetNode(stack[len(stack)-1].Node)
		stack = stack[:len(stack)-1]
		if popped.IsLeaf() {
			index, _ := popped.findKey(key)
			popped.DeleteAt(index)
		}

		if len(stack) == 0 {
			// if no parent left in stack it is done
			return true
		}
		indexAtParent := stack[len(stack)-1].Index
		parent := tree.pager.GetNode(stack[len(stack)-1].Node)

		if popped.IsUnderFlow(tree.degree) {
			// get siblings
			var rightSibling, leftSibling, merged Node
			if indexAtParent > 0 {
				leftSibling = tree.pager.GetNode(parent.GetValueAt(indexAtParent - 1).(Pointer)) //leftSibling = parent.Pointers[indexAtParent-1].(*InternalNode)
			}
			if indexAtParent+1 < (parent.Keylen() + 1) { // +1 is the length of pointers
				rightSibling = tree.pager.GetNode(parent.GetValueAt(indexAtParent + 1).(Pointer)) //rightSibling = parent.Pointers[indexAtParent+1].(*InternalNode)
			}

			//try redistribute
			if rightSibling != nil &&
				((popped.IsLeaf() && rightSibling.Keylen() >= (tree.degree/2)+1) ||
					(!popped.IsLeaf() && rightSibling.Keylen()+1 > (tree.degree+1)/2)) { // TODO: second check is actually different for internal and leaf nodes since internal nodes have one more value than they have keys
				popped.Redistribute(rightSibling, parent)
				return true
			} else if leftSibling != nil &&
				((popped.IsLeaf() && leftSibling.Keylen() >= (tree.degree/2)+1) ||
					(!popped.IsLeaf() && leftSibling.Keylen()+1 > (tree.degree+1)/2)) {
				leftSibling.Redistribute(popped, parent)
				return true
			}

			// if redistribution is not valid merge
			if rightSibling != nil {
				popped.MergeNodes(rightSibling, parent)
				merged = popped
			} else {
				if leftSibling == nil {
					if !popped.IsLeaf() {
						panic("Both siblings are null for an internal Node! This should not be possible except for root")
					}

					// TODO: may be log here? if it is a leaf node its both left and right nodes can be nil
					return true
				}
				leftSibling.MergeNodes(popped, parent)
				merged = leftSibling
			}
			if parent.GetPageId() == tree.Root && parent.Keylen() == 0 {
				tree.Root = merged.GetPageId()
			}
		} else {
			break
		}
	}
	return true
}
