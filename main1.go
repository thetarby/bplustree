package main

import (
	"awesomeProject/btree"
	"fmt"
)

type myInt int

func (key myInt) Less(than btree.Key) bool {
	return key < than.(myInt)
}

func main() {
	main2()
}

func main1() {
	fmt.Println("selam")
	tree := btree.NewBtree(3)
	tree.Insert(myInt(1), "selam1")
	tree.Insert(myInt(5), "selam5")
	tree.Insert(myInt(3), "selam3")
	tree.Insert(myInt(2), "selam2")
	for _, val := range []int{81, 87, 47, 59, 82, 88, 89} {
		tree.Insert(myInt(val), "selam")
	}

	//leftNode, _ , _ := tree.LeafRoot.SplitNode(1)
	tree.Root.PrintNode()
	tree.Root.Pointers[0].(*btree.LeafNode).PrintNode()
	tree.Root.Pointers[1].(*btree.LeafNode).PrintNode()

	res := tree.Find(myInt(7))
	fmt.Println(res)

	//leftNode.PrintNode()
}

func main2() {
	fmt.Println("selam")
	tree := btree.NewBtree(3)
	tree.Insert(myInt(1), "selam1")
	tree.Insert(myInt(5), "selam5")
	tree.Insert(myInt(3), "selam3")
	tree.Insert(myInt(2), "selam2")
	tree.Print()
	for _, val := range []int{81, 87, 47, 59, 82, 88, 89} {
		tree.Insert(myInt(val), "selam")
		fmt.Println("new tree: !!!")
		tree.Print()
	}

	tree.Print()

	tree.Delete(myInt(3))
	fmt.Println("After Delete 3 !!!!!!!!!!!!!!!!!!!!!!!!!!!")
	tree.Print()
	tree.Delete(myInt(5))
	fmt.Println("After Delete 5 !!!!!!!!!!!!!!!!!!!!!!!!!!!")
	tree.Print()
	tree.Delete(myInt(1))
	fmt.Println("After Delete 1 !!!!!!!!!!!!!!!!!!!!!!!!!!!")
	tree.Print()

	//leftNode, _ , _ := tree.LeafRoot.SplitNode(1)
	tree.Root.PrintNode()
	tree.Root.Pointers[0].(*btree.LeafNode).PrintNode()
	tree.Root.Pointers[1].(*btree.LeafNode).PrintNode()

	res := tree.Find(myInt(7))
	fmt.Println(res)

	//leftNode.PrintNode()
}

func main3() {
	tree := btree.NewBtree(4)
	for _, val := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		tree.Insert(myInt(val), "selam")
		fmt.Println("new tree: !!!")
		tree.Print()
	}

}
