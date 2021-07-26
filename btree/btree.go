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
	findAndGetStack(key Key, stackIn []node) (value interface{}, stackOut []node)
	SplitNode(index int) (left node, keyAtLeft Key, keyAtRight Key)
	PrintNode()
}

type InternalNode struct {
	Keys     Keys
	Pointers []node
}

type LeafNode struct {
	keys   Keys
	values []interface{}
	Right  *LeafNode
	Left   *LeafNode
}

type BTree struct {
	degree   int
	length   int
	Root     *InternalNode
	LeafRoot *LeafNode
}

func NewBtree(degree int) BTree {
	p := make([]node, 1, 2)
	l := LeafNode{
		keys:   make(Keys, 0, 2),
		values: make([]interface{}, 0, 2),
		Right:  nil,
		Left:   nil,
	}
	p[0] = &l

	return BTree{
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
		keys:   make([]Key, 0, 2),
		values: make([]interface{}, 0, 2),
	}
}

func (n *LeafNode) PrintNode() {
	fmt.Printf("Node( ")
	for i := 0; i < len(n.keys); i++ {
		fmt.Printf("%d | ", n.keys[i])
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

func (n *InternalNode) truncate(index int) {
	n.Keys = n.Keys[:index]
	n.Pointers = n.Pointers[:index+1]
}

func (n *InternalNode) internalInsertAt(index int, key Key, pointer node) {
	n.Keys = append(n.Keys, key)
	copy(n.Keys[index+1:], n.Keys[index:])
	n.Keys[index] = key

	n.Pointers = append(n.Pointers, pointer)
	copy(n.Pointers[index+1:], n.Pointers[index:])
	n.Pointers[index] = pointer
}

func (n *InternalNode) internalInsertAt2(index int, key Key, pointer node) {
	n.Keys = append(n.Keys, key)
	copy(n.Keys[index+1:], n.Keys[index:])
	n.Keys[index] = key

	n.Pointers = append(n.Pointers, pointer)
	copy(n.Pointers[index+2:], n.Pointers[index+1:])
	n.Pointers[index+1] = pointer
}

func (n *LeafNode) leafInsertAt(index int, key Key, value interface{}) {
	n.keys = append(n.keys, key)
	copy(n.keys[index+1:], n.keys[index:])
	n.keys[index] = key

	n.values = append(n.values, value)
	copy(n.values[index+1:], n.values[index:])
	n.values[index] = value
}

func (n *LeafNode) leafInsertAt2(index int, key Key, value interface{}) {
	n.keys = append(n.keys, key)
	copy(n.keys[index+1:], n.keys[index:])
	n.keys[index] = key

	n.values = append(n.values, value)
	copy(n.values[index+1:], n.values[index:])
	n.values[index] = value
}

func (n *InternalNode) SplitNode(index int) (leftNode node, keyAtLeft Key, keyAtRight Key) {
	left := newNode()
	keyAtLeft = n.Keys[index-1]
	keyAtRight = n.Keys[index]
	left.Keys = n.Keys
	left.Pointers = n.Pointers
	n.Pointers = make([]node, 0, 2)
	n.Keys = make(Keys, 0, 2)

	n.Keys = append(n.Keys, left.Keys[index+1:]...)
	// TODO: yeni node bir eksik pointer ile kaldıyor ya la
	n.Pointers = append(n.Pointers, left.Pointers[index+1:]...)

	left.truncate(index)
	return left, keyAtLeft, keyAtRight
}

func (n *InternalNode) SplitNode2(index int) (rightNode node, keyAtLeft Key, keyAtRight Key) {
	right := newNode()
	keyAtLeft = n.Keys[index-1]
	keyAtRight = n.Keys[index]

	right.Keys = append(right.Keys, n.Keys[index:]...)
	right.Pointers = append(right.Pointers, n.Pointers[index:]...)
	n.Keys = n.Keys[:index]
	n.Pointers = n.Pointers[:index]

	n.truncate(index)
	return right, keyAtLeft, keyAtRight
}

func (n *LeafNode) SplitNode(index int) (leftNode node, keyAtLeft Key, keyAtRight Key) {
	left := newLeafNode()
	keyAtLeft = n.keys[index-1]
	keyAtRight = n.keys[index]
	left.keys = n.keys
	left.values = n.values
	n.values = make([]interface{}, 0, 2)
	n.keys = make(Keys, 0, 2)

	n.keys = append(n.keys, left.keys[index:]...)
	n.values = append(n.values, left.values[index:]...)

	left.keys = left.keys[:index]
	left.values = left.values[:index]
	left.Left = n.Left
	left.Right = n
	n.Left = left
	return left, keyAtLeft, keyAtRight
}

func (n *LeafNode) SplitNode2(index int) (rightNode node, keyAtLeft Key, keyAtRight Key) {
	right := newLeafNode()
	keyAtLeft = n.keys[index-1]
	keyAtRight = n.keys[index]

	right.keys = append(right.keys, n.keys[index:]...)
	right.values = append(right.values, n.values[index:]...)
	n.keys = n.keys[:index]
	n.values = n.values[:index]
	n.Right = right
	right.Left = n
	right.Right = n.Right
	return right, keyAtLeft, keyAtRight
}

func (tree *BTree) Insert(key Key, value interface{}) {
	var stack = make([]node, 0, 0)
	var i interface{}
	i, stack = tree.Root.findAndGetStack(key, stack)
	if i != nil {
		panic("key already exists")
	}

	leafNode := stack[len(stack)-1].(*LeafNode)
	index, _ := leafNode.keys.find(key)
	leafNode.leafInsertAt(index, key, value)
	stack = stack[:len(stack)-1]

	if len(leafNode.values) == tree.degree {
		var leftNode node
		var rightKey Key

		leftNode, _, rightKey = leafNode.SplitNode((len(leafNode.values) - 1) / 2)
		var rightNode node = leafNode
		for len(stack) > 0 {
			top := stack[len(stack)-1].(*InternalNode)
			stack = stack[:len(stack)-1]

			i, _ := top.Keys.find(key)
			top.internalInsertAt(i, rightKey, leftNode)
			if len(top.Pointers) == tree.degree+1 {
				leftNode, _, rightKey = top.SplitNode((len(top.Pointers) - 1) / 2)
				rightNode = top
				// if top is root node special case, create new root
				if top == tree.Root {
					tree.Root = &InternalNode{
						Keys:     Keys{rightKey},
						Pointers: []node{leftNode, rightNode},
					}
					break
				}

				// TODO: special case if top's right most pointer is the one that is in the stack

			} else {
				break
			}
		}
	}
}

func (tree BTree) Find(key Key) interface{} {
	res, _ := tree.Root.findAndGetStack(key, []node{})
	return res
}

func (n *InternalNode) findAndGetStack(key Key, stackIn []node) (value interface{}, stackOut []node) {
	i, found := n.Keys.find(key)
	if found {
		i++
	}
	stackOut = append(stackIn, n)
	res, stackOut := n.Pointers[i].findAndGetStack(key, stackOut)
	return res, stackOut
}

func (n *LeafNode) findAndGetStack(key Key, stackIn []node) (value interface{}, stackOut []node) {
	i, found := n.keys.find(key)
	stackOut = append(stackIn, n)
	if !found {
		return nil, stackOut
	}
	return n.values[i], stackOut
}

func (tree BTree) Print() {
	arr := make([]node, 0, 2)
	arr = append(arr, tree.Root)
	//arr := tree.Root.Pointers
	arr = append(arr, nil)
	for i := 0; i < len(arr); i++ {
		if arr[i] != nil && reflect.TypeOf(arr[i]) == reflect.TypeOf(newLeafNode()) {
			break
		}
		if arr[i] == nil {
			arr = append(arr, nil)
			continue
		}
		node := arr[i].(*InternalNode)

		arr = append(arr, node.Pointers...)
	}
	for _, n := range arr {
		if n != nil {
			n.PrintNode()
		} else {
			fmt.Print("\n ### \n")
		}
	}
}
