package main

import (
	"awesomeProject/btree"
	"fmt"
	"math/rand"
)

type myInt int

func (key myInt) Less(than btree.Key) bool{
	return key < than.(myInt)
}

func main() {
	fmt.Println("selam")
	tree := btree.NewBtree(3)
	tree.Insert(myInt(1), "selam1")
	tree.Insert(myInt(5), "selam5")
	tree.Insert(myInt(3), "selam3")
	tree.Insert(myInt(2), "selam2")
	for i := 0; i < 5; i++ {
		tree.Insert(myInt(rand.Intn(100)), "selam")
	}
	//leftNode, _ , _ := tree.LeafRoot.SplitNode(1)
	tree.Root.PrintNode()
	tree.Root.Pointers[0].(*btree.LeafNode).PrintNode()
	tree.Root.Pointers[1].(*btree.LeafNode).PrintNode()

	res := tree.Find(myInt(7))
	fmt.Println(res)
	//leftNode.PrintNode()
}
