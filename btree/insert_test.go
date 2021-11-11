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

func TestInsert_Internals(t *testing.T) {
	tree := NewBtree(4)
	tree.Insert(myInt(1), "1")

	stack := make([]NodeIndexPair, 0)
	val, stack := tree.Root.findAndGetStack(myInt(1), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "1", val.(string))

	tree.Insert(myInt(2), "2")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(2), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "2", val.(string))

	tree.Insert(myInt(3), "3")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(3), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "3", val.(string))

	tree.Insert(myInt(4), "4")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(4), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "4", val.(string))

	tree.Insert(myInt(5), "5")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(5), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "5", val.(string))

	tree.Insert(myInt(6), "6")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(6), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "6", val.(string))

	tree.Insert(myInt(7), "7")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(7), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "7", val.(string))

	tree.Insert(myInt(8), "8")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(8), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "8", val.(string))

	tree.Insert(myInt(9), "9")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(9), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, "9", val.(string))

	tree.Insert(myInt(10), "10")
	stack = make([]NodeIndexPair, 0)
	val, stack = tree.Root.findAndGetStack(myInt(10), stack)

	assert.Len(t, stack, 3)
	assert.Equal(t, "10", val.(string))
}
