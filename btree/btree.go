package btree

// TODO:
// 1. don't use Keys.find instead define one on Node and use that method.
// 2. define util methods to change state of node such as shift_keys_by_n, shift_pointers_by_n, truncate_at_n
//    and use them in methods like InsertAt and SplitNode
// 3. Put a constraint on Key to make it fix sized maybe? It will solve many problems when we try to persist nodes on a disk
// 4. Use interface methods in delete.go as well.

import (
	"fmt"
	"reflect"
	"sort"
)

var pager Pager = &NoopPager{}

type Pointer int64

type Key interface {
	Less(than Key) bool
}

type Keys []Key

func (k Keys) find(item Key) (index int, found bool) {
	i := sort.Search(len(k), func(i int) bool {
		return item.Less(k[i])
	})
	if i > 0 && !k[i-1].Less(item) {
		return i - 1, true
	}
	return i, false
}

type NodeIndexPair struct {
	Node  Node
	Index int // pointer index for internal nodes and value index for leaf nodes
}

type Node interface {
	// findAndGetStack is used to recursively find the given key and it also passes a stack object recursively to
	// keep the path it followed down to leaf node. value is nil when key does not exist.
	findAndGetStack(key Key, stackIn []NodeIndexPair) (value interface{}, stackOut []NodeIndexPair)
	findKey(key Key) (index int, found bool)
	shiftKeyValueAt(n int)
	setKeyAt(idx int, key Key)
	setValueAt(idx int, val interface{})
	SplitNode(index int) (left Pointer, keyAtLeft Key, keyAtRight Key)
	PrintNode()
	IsOverFlow(degree int) bool
	InsertAt(index int, key Key, val interface{})
	GetPageId() Pointer
}

type InternalNode struct {
	PersistentPage
	Keys     Keys
	Pointers []Pointer
}

type LeafNode struct {
	PersistentPage
	Keys   Keys
	Values []interface{}
	Right  *LeafNode
	Left   *LeafNode
}

type BTree struct {
	degree   int
	length   int
	Root     *InternalNode
	LeafRoot *LeafNode
}

func NewBtree(degree int) *BTree {
	l := pager.NewLeafNode()

	root := pager.NewInternalNode()
	// TODO: use a method here to be persistent
	root.Pointers = []Pointer{l.GetPageId()}

	return &BTree{
		degree: degree,
		length: 0,
		Root:   root,
	}
}

func newNode() (n *InternalNode) {
	return pager.NewInternalNode()
}

func newLeafNode() (n *LeafNode) {
	return pager.NewLeafNode()
}

func (n *LeafNode) findKey(key Key) (index int, found bool) {
	return n.Keys.find(key)
}

func (n *InternalNode) findKey(key Key) (index int, found bool) {
	return n.Keys.find(key)
}

func (n *LeafNode) shiftKeyValueAt(at int) {
	n.Keys = append(n.Keys, nil)
	n.Values = append(n.Values, nil)
	copy(n.Keys[at+1:], n.Keys[at:])
	copy(n.Values[at+1:], n.Values[at:])
}

func (n *InternalNode) shiftKeyValueAt(at int) {
	n.Keys = append(n.Keys, nil)
	var zeroPointer Pointer
	n.Pointers = append(n.Pointers, zeroPointer)
	copy(n.Keys[at+1:], n.Keys[at:])
	copy(n.Pointers[at+2:], n.Pointers[at+1:])
}

func (n *LeafNode) setKeyAt(idx int, key Key) {
	n.Keys[idx] = key
}

func (n *InternalNode) setKeyAt(idx int, key Key) {
	n.Keys[idx] = key
}

func (n *LeafNode) setValueAt(idx int, val interface{}) {
	n.Values[idx] = val
}

func (n *InternalNode) setValueAt(idx int, val interface{}) {
	n.Pointers[idx] = val.(Pointer)
}

func (n *LeafNode) IsOverFlow(degree int) bool {
	return len(n.Values) == degree
}

func (n *InternalNode) IsOverFlow(degree int) bool {
	return len(n.Pointers) == degree+1
}

func (n *LeafNode) InsertAt(index int, key Key, value interface{}) {
	n.shiftKeyValueAt(index)
	n.setKeyAt(index, key)
	n.setValueAt(index, value)
	pager.SyncLeafNode(n)
}

func (n *InternalNode) InsertAt(index int, key Key, pointer interface{}) {
	n.shiftKeyValueAt(index)
	n.setKeyAt(index, key)
	n.setValueAt(index+1, pointer)
	pager.SyncInternalNode(n)
}

func (n *LeafNode) SplitNode(index int) (rightNode Pointer, keyAtLeft Key, keyAtRight Key) {
	right := newLeafNode()
	keyAtLeft = n.Keys[index-1]
	keyAtRight = n.Keys[index]

	right.Keys = append(right.Keys, n.Keys[index:]...)
	right.Values = append(right.Values, n.Values[index:]...)
	n.Keys = n.Keys[:index]
	n.Values = n.Values[:index]
	right.Right = n.Right
	n.Right = right
	right.Left = n

	pager.SyncLeafNode(n)
	pager.SyncLeafNode(right)
	return right.GetPageId(), keyAtLeft, keyAtRight
}

func (n *InternalNode) SplitNode(index int) (rightNode Pointer, keyAtLeft Key, keyAtRight Key) {
	right := newNode()
	keyAtLeft = n.Keys[index-1]
	keyAtRight = n.Keys[index]

	right.Keys = append(right.Keys, n.Keys[index+1:]...)
	right.Pointers = append(right.Pointers, n.Pointers[index+1:]...)
	n.Keys = n.Keys[:index]
	n.Pointers = n.Pointers[:index]

	n.truncate(index)

	pager.SyncInternalNode(n)
	pager.SyncInternalNode(right)
	return right.GetPageId(), keyAtLeft, keyAtRight
}

func (n *LeafNode) findAndGetStack(key Key, stackIn []NodeIndexPair) (value interface{}, stackOut []NodeIndexPair) {
	i, found := n.findKey(key)
	stackOut = append(stackIn, NodeIndexPair{n, i})
	if !found {
		return nil, stackOut
	}
	return n.Values[i], stackOut
}

func (n *InternalNode) findAndGetStack(key Key, stackIn []NodeIndexPair) (value interface{}, stackOut []NodeIndexPair) {
	i, found := n.findKey(key)
	if found {
		i++
	}
	stackOut = append(stackIn, NodeIndexPair{n, i})
	node := pager.GetNode(n.Pointers[i])
	res, stackOut := node.findAndGetStack(key, stackOut)
	return res, stackOut
}

func (n *LeafNode) PrintNode() {
	fmt.Printf("Node( ")
	for i := 0; i < len(n.Keys); i++ {
		fmt.Printf("%d | ", n.Keys[i])
	}
	fmt.Printf(")    ")
}

func (n *InternalNode) PrintNode() {
	fmt.Printf("Node( ")
	for i := 0; i < len(n.Keys); i++ {
		fmt.Printf("%d | ", n.Keys[i])
	}
	fmt.Printf(")    ")
}

func (tree *BTree) Insert(key Key, value interface{}) {
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.Root.findAndGetStack(key, stack)
	if i != nil {
		panic("key already exists")
	}

	var rightNod = value
	var rightKey = key

	for len(stack) > 0 {
		topOfStack := stack[len(stack)-1].Node
		i, _ := topOfStack.findKey(key)
		topOfStack.InsertAt(i, rightKey, rightNod)

		popped := stack[len(stack)-1].Node
		stack = stack[:len(stack)-1]
		if popped.IsOverFlow(tree.degree) {
			rightNod, _, rightKey = popped.SplitNode((tree.degree) / 2)
			if popped == tree.Root {
				leftNode := popped

				// TODO: a method should be called here
				tree.Root = pager.NewInternalNode()
				tree.Root.Keys = Keys{rightKey}
				tree.Root.Pointers = []Pointer{leftNode.GetPageId(), rightNod.(Pointer)}
				pager.SyncInternalNode(tree.Root)
			}
		} else {
			break
		}
	}
}

func (tree *BTree) InsertOrReplace(key Key, value interface{}) (isInserted bool) {
	var stack = make([]NodeIndexPair, 0, 0)
	var i interface{}
	i, stack = tree.Root.findAndGetStack(key, stack)
	if i != nil {
		// top of stack is the leaf Node
		topOfStack := stack[len(stack)-1]
		leafNode := topOfStack.Node.(*LeafNode)
		leafNode.Values[topOfStack.Index] = value
		return false
	}

	var rightNod = value
	var rightKey = key

	for len(stack) > 0 {
		topOfStack := stack[len(stack)-1].Node
		i, _ := topOfStack.findKey(key)
		topOfStack.InsertAt(i, rightKey, rightNod)

		popped := stack[len(stack)-1].Node
		stack = stack[:len(stack)-1]
		if popped.IsOverFlow(tree.degree) {
			rightNod, _, rightKey = popped.SplitNode((tree.degree) / 2)
			if popped == tree.Root {
				leftNode := popped

				// TODO: a method should be called here
				tree.Root = pager.NewInternalNode()
				tree.Root.Keys = Keys{rightKey}
				tree.Root.Pointers = []Pointer{leftNode.GetPageId(), rightNod.(Pointer)}
				pager.SyncInternalNode(tree.Root)
			}
		} else {
			break
		}
	}

	return true
}

func (tree *BTree) Find(key Key) interface{} {
	res, _ := tree.Root.findAndGetStack(key, []NodeIndexPair{})
	return res
}

func (tree *BTree) Height() int {
	var currentNode Node = tree.Root
	acc := 0
	for {
		switch currentNode.(type) {
		case *InternalNode:
			currentNode = pager.GetNode(currentNode.(*InternalNode).Pointers[0])
		case *LeafNode:
			return acc + 1
		}
		acc++
	}
}

func (tree BTree) Print() {
	queue := make([]Pointer, 0, 2)
	queue = append(queue, tree.Root.GetPageId())
	queue = append(queue, 0)
	for i := 0; i < len(queue); i++ {
		if queue[i] != 0 && reflect.TypeOf(queue[i]) == reflect.TypeOf(newLeafNode()) {
			break
		}
		if queue[i] == 0 {
			queue = append(queue, 0)
			continue
		}

		node := pager.GetNode(queue[i]).(*InternalNode)

		queue = append(queue, node.Pointers...)
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
