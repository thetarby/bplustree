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

	for _, val := range []int{10, 9, 8, 7} {
		tree.Delete(myInt(val))
		tree.Print()
		var stack []NodeIndexPair
		res, stack := tree.Root.findAndGetStack(myInt(1), stack)
		assert.Len(t, stack, 3)
		assert.Equal(t, "selam", res.(string))
	}
	tree.Delete(myInt(6))
	var stack []NodeIndexPair
	_, stack = tree.Root.findAndGetStack(myInt(1), stack)

	assert.Len(t, stack, 2)
}
