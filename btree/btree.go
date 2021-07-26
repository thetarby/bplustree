package btree

import (
	"fmt"
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
}

type InternalNode struct {
	keys     Keys
	Pointers []node
}

type LeafNode struct {
	keys   Keys
	values []interface{}
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
	}
	p[0] = &l

	return BTree{
		degree: degree,
		length: 0,
		Root: &InternalNode{
			keys:     make(Keys, 0, 2),
			Pointers: p,
		},
	}
}

func newNode() (n *InternalNode) {
	return &InternalNode{
		keys:     make([]Key, 0, 2),
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
	fmt.Println("node:")
	for i := 0; i < len(n.keys); i++ {
		fmt.Printf("Key: %d, Value: %s \n", n.keys[i], n.values[i])
	}
}

func (n *InternalNode) PrintNode() {
	fmt.Println("node:")
	for i := 0; i < len(n.keys); i++ {
		fmt.Printf("Key: %d, Value: %s \n", n.keys[i], n.Pointers[i])
	}
}

func (n *InternalNode) truncate(index int) {
	n.keys = n.keys[:index]
	n.Pointers = n.Pointers[:index+1]
}

func (n *InternalNode) internalInsertAt(index int, key Key, pointer node) {
	n.keys = append(n.keys, key)
	copy(n.keys[index+1:], n.keys[index:])
	n.keys[index] = key

	n.Pointers = append(n.Pointers, pointer)
	copy(n.Pointers[index+1:], n.Pointers[index:])
	n.Pointers[index] = pointer
}

func (n *LeafNode) leafInsertAt(index int, key Key, value interface{}) {
	n.keys = append(n.keys, key)
	copy(n.keys[index+1:], n.keys[index:])
	n.keys[index] = key

	n.values = append(n.values, value)
	copy(n.values[index+1:], n.values[index:])
	n.values[index] = value
}

func (n *InternalNode) SplitNode(index int) (leftNode node, keyAtLeft Key, keyAtRight Key) {
	left := newNode()
	keyAtLeft = n.keys[index]
	keyAtRight = n.keys[index+1]
	left.keys = n.keys
	left.Pointers = n.Pointers
	n.Pointers = make([]node, 0, 2)
	n.keys = make(Keys, 0, 2)

	n.keys = append(n.keys, left.keys[index+1:]...)
	// TODO: yeni node bir eksik pointer ile kaldıyor ya la
	n.Pointers = append(n.Pointers, left.Pointers[index+2:]...)

	left.truncate(index + 1)
	return left, keyAtLeft, keyAtRight
}

func (n *LeafNode) SplitNode(index int) (leftNode node, keyAtLeft Key, keyAtRight Key) {
	left := newLeafNode()
	keyAtLeft = n.keys[index]
	keyAtRight = n.keys[index+1]
	left.keys = n.keys
	left.values = n.values
	n.values = make([]interface{}, 0, 2)
	n.keys = make(Keys, 0, 2)

	n.keys = append(n.keys, left.keys[index+1:]...)
	n.values = append(n.values, left.values[index+1:]...)

	left.keys = left.keys[:index+1]
	left.values = left.values[:index+1]
	return left, keyAtLeft, keyAtRight
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

			i, _ := top.keys.find(key)
			if i == len(top.keys) {
				top.Pointers[len(top.Pointers)-1] = leftNode
				top.Pointers = append(top.Pointers, rightNode)
				top.keys = append(top.keys, rightKey)
			} else {
				top.internalInsertAt(i, rightKey, leftNode)
			}
			if len(top.Pointers) == tree.degree+1 {
				leftNode, _, rightKey = top.SplitNode((len(top.Pointers) - 1) / 2)
				rightNode = top
				// if top is root node special case, create new root
				if top == tree.Root {
					tree.Root = &InternalNode{
						keys:     Keys{rightKey},
						Pointers: []node{leftNode, top},
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
	i, found := n.keys.find(key)
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
