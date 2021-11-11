package btree

import (
	"fmt"
	"reflect"
	"sort"
)

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

type node interface {
	findAndGetStack(key Key, stackIn []NodeIndexPair) (value interface{}, stackOut []NodeIndexPair)
	findKey(key Key) (index int, found bool)
	SplitNode(index int) (left node, keyAtLeft Key, keyAtRight Key)
	PrintNode()
	IsOverFlow(degree int) bool
	InsertAt(index int, key Key, val interface{})
}

type InternalNode struct {
	Keys     Keys
	Pointers []node
}

type LeafNode struct {
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
	p := make([]node, 1, 2)
	l := LeafNode{
		Keys:   make(Keys, 0, 2),
		Values: make([]interface{}, 0, 2),
		Right:  nil,
		Left:   nil,
	}
	p[0] = &l

	return &BTree{
		degree: degree,
		length: 0,
		Root: &InternalNode{
			Keys:     make(Keys, 0, 2),
			Pointers: p,
		},
	}
}

func newNode() (n *InternalNode) {
	return &InternalNode{
		Keys:     make([]Key, 0, 2),
		Pointers: make([]node, 0, 2),
	}
}

func newLeafNode() (n *LeafNode) {
	return &LeafNode{
		Keys:   make([]Key, 0, 2),
		Values: make([]interface{}, 0, 2),
	}
}

func (n *LeafNode) findKey(key Key) (index int, found bool) {
	return n.Keys.find(key)
}

func (n *InternalNode) findKey(key Key) (index int, found bool) {
	return n.Keys.find(key)
}

func (n *LeafNode) IsOverFlow(degree int) bool {
	return len(n.Values) == degree
}

func (n *InternalNode) IsOverFlow(degree int) bool {
	return len(n.Pointers) == degree+1
}

func (n *InternalNode) truncate(index int) {
	n.Keys = n.Keys[:index]
	n.Pointers = n.Pointers[:index+1]
}

func (n *InternalNode) InsertAt(index int, key Key, pointer interface{}) {
	n.Keys = append(n.Keys, key)
	copy(n.Keys[index+1:], n.Keys[index:])
	n.Keys[index] = key

	n.Pointers = append(n.Pointers, pointer.(node))
	copy(n.Pointers[index+2:], n.Pointers[index+1:])
	n.Pointers[index+1] = pointer.(node)
}

func (n *LeafNode) InsertAt(index int, key Key, value interface{}) {
	n.Keys = append(n.Keys, key)
	copy(n.Keys[index+1:], n.Keys[index:])
	n.Keys[index] = key

	n.Values = append(n.Values, value)
	copy(n.Values[index+1:], n.Values[index:])
	n.Values[index] = value
}

func (n *InternalNode) SplitNode(index int) (rightNode node, keyAtLeft Key, keyAtRight Key) {
	right := newNode()
	keyAtLeft = n.Keys[index-1]
	keyAtRight = n.Keys[index]

	right.Keys = append(right.Keys, n.Keys[index+1:]...)
	right.Pointers = append(right.Pointers, n.Pointers[index+1:]...)
	n.Keys = n.Keys[:index]
	n.Pointers = n.Pointers[:index]

	n.truncate(index)
	return right, keyAtLeft, keyAtRight
}

func (n *LeafNode) SplitNode(index int) (rightNode node, keyAtLeft Key, keyAtRight Key) {
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
	return right, keyAtLeft, keyAtRight
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
				tree.Root = &InternalNode{
					Keys:     Keys{rightKey},
					Pointers: []node{leftNode, rightNod.(*InternalNode)},
				} // if it is root it should be the last item in the stack so loop will be breaked
			}
		} else {
			break
		}
	}
}

func (tree *BTree) Find(key Key) interface{} {
	res, _ := tree.Root.findAndGetStack(key, []NodeIndexPair{})
	return res
}

func (tree *BTree) Height() int {
	var currentNode node = tree.Root
	acc := 0
	for {
		switch currentNode.(type) {
		case *InternalNode:
			currentNode = currentNode.(*InternalNode).Pointers[0]
		case *LeafNode:
			return acc + 1
		}
		acc++
	}
}

type NodeIndexPair struct {
	Node  node
	Index int // pointer index for internal nodes and value index for leaf nodes
}

func (n *InternalNode) findAndGetStack(key Key, stackIn []NodeIndexPair) (value interface{}, stackOut []NodeIndexPair) {
	i, found := n.Keys.find(key)
	if found {
		i++
	}
	stackOut = append(stackIn, NodeIndexPair{n, i})
	res, stackOut := n.Pointers[i].findAndGetStack(key, stackOut)
	return res, stackOut
}

func (n *LeafNode) findAndGetStack(key Key, stackIn []NodeIndexPair) (value interface{}, stackOut []NodeIndexPair) {
	i, found := n.Keys.find(key)
	stackOut = append(stackIn, NodeIndexPair{n, i})
	if !found {
		return nil, stackOut
	}
	return n.Values[i], stackOut
}

func (tree BTree) Print() {
	queue := make([]node, 0, 2)
	queue = append(queue, tree.Root)
	queue = append(queue, nil)
	for i := 0; i < len(queue); i++ {
		if queue[i] != nil && reflect.TypeOf(queue[i]) == reflect.TypeOf(newLeafNode()) {
			break
		}
		if queue[i] == nil {
			queue = append(queue, nil)
			continue
		}
		node := queue[i].(*InternalNode)

		queue = append(queue, node.Pointers...)
	}
	for _, n := range queue {
		if n != nil {
			n.PrintNode()
		} else {
			fmt.Print("\n ### \n")
		}
	}
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
