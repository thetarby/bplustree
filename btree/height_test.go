package btree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeight_Should_Return_Correct_Height(t *testing.T) {
	tests := []struct {
		tree     *BTree
		toInsert []myInt
		expected int
	}{
		{
			tree:     NewBtree(3),
			toInsert: []myInt{myInt(1), myInt(2), myInt(3), myInt(4), myInt(5), myInt(6), myInt(7), myInt(8), myInt(9)},
			expected: 4,
		},
		{
			tree:     NewBtree(4),
			toInsert: []myInt{myInt(1), myInt(2), myInt(3), myInt(4)},
			expected: 2,
		},
		{
			tree:     NewBtree(5),
			toInsert: []myInt{myInt(1), myInt(2), myInt(3), myInt(4), myInt(5)},
			expected: 2,
		},
	}
	for _, test := range tests {
		for _, m := range test.toInsert {
			test.tree.Insert(m, "value")
		}
		h := test.tree.Height()
		assert.Equal(t, test.expected, h)

	}
}
