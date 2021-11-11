package btree

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

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

func TestInsert_Or_Replace_Should_Return_False_When_Key_Exists(t *testing.T) {
	tree := NewBtree(3)
	for i := 0; i < 1000; i++ {
		tree.Insert(myInt(i), strconv.Itoa(i))
	}

	isInserted := tree.InsertOrReplace(myInt(500), "new_500")

	assert.False(t, isInserted)
}

func TestInsert_Or_Replace_Should_Replace_Value_When_Key_Exists(t *testing.T) {
	tree := NewBtree(3)
	for i := 0; i < 1000; i++ {
		tree.Insert(myInt(i), strconv.Itoa(i))
	}

	tree.InsertOrReplace(myInt(500), "new_500")
	val := tree.Find(myInt(500))

	assert.Equal(t, "new_500", val.(string))
}

func TestAll_Inserts_Should_Be_Found_By_Find_Method(t *testing.T) {
	tree := NewBtree(3)
	arr := make([]int, 0)
	for i := 0; i < 1000; i++ {
		arr = append(arr, i)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

	for _, item := range arr {
		tree.Insert(myInt(item), strconv.Itoa(item))
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	for _, item := range arr {
		val := tree.Find(myInt(item))
		assert.NotNil(t, val)
		assert.Equal(t, strconv.Itoa(item), val.(string))
	}
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
