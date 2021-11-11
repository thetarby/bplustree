package btree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDemo(t *testing.T) {
	got := 4 + 6
	want := 10

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

type myInt int

func (key myInt) Less(than Key) bool {
	return key < than.(myInt)
}

func TestInsert_Should_Split_Root_When_It_Has_M_Keys(t *testing.T) {
	tree := NewBtree(3)
	tree.Insert(myInt(1), "1")
	tree.Insert(myInt(5), "5")
	tree.Insert(myInt(3), "3")

	var stack []NodeIndexPair

	res, stack := tree.Root.findAndGetStack(myInt(5), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "5", res)
	assert.Equal(t, myInt(3), tree.Root.Keys[0])
}
