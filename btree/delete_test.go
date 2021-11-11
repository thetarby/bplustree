package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete_Should_Decrease_Height_Size_When_Root_Is_Empty(t *testing.T) {
	tree := NewBtree(4)
	for _, val := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		tree.Insert(myInt(val), "selam")
	}
	var stack []NodeIndexPair
	res, stack := tree.Root.findAndGetStack(myInt(1), stack)
	assert.Len(t, stack, 3)
	assert.Equal(t, "selam", res.(string))

	tree.Delete(myInt(1))
	stack = []NodeIndexPair{}
	_, stack = tree.Root.findAndGetStack(myInt(1), stack)

	assert.Len(t, stack, 2)
}

func TestDelete_Should_Decrease_Height_Size_When_Root_Is_Empty_2(t *testing.T) {
	tree := NewBtree(3)

	for _, val := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		tree.Insert(myInt(val), "selam")
	}

	var stack []NodeIndexPair
	res, stack := tree.Root.findAndGetStack(myInt(1), stack)
	assert.Len(t, stack, 4)
	assert.Equal(t, "selam", res.(string))

	tree.Print()
	for _, i := range []int{1, 2, 3, 4, 5} {
		var stack []NodeIndexPair
		tree.Delete(myInt(i))
		tree.Print()
		res, stack := tree.Root.findAndGetStack(myInt(10), stack)
		assert.Len(t, stack, 3)
		assert.Equal(t, "selam", res.(string))
	}

	tree.Delete(myInt(6))
	stack = []NodeIndexPair{}
	_, stack = tree.Root.findAndGetStack(myInt(10), stack)
	assert.Len(t, stack, 2)
}
