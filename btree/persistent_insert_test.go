package btree

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestPersistent_Insert_Should_Split_Root_When_It_Has_M_Keys(t *testing.T) {
	tree := NewBtreeWithPager(3, NoopPersistentPager{KeySerializer: &PersistentKeySerializer{}, KeySize: 8})
	p := SlotPointer{
		PageId:  10,
		SlotIdx: 10,
	}
	tree.Insert(PersistentKey(1), p)
	tree.Insert(PersistentKey(5), p)
	tree.Insert(PersistentKey(3), p)

	var stack []NodeIndexPair

	res, stack := tree.GetRoot().findAndGetStack(PersistentKey(5), stack)

	assert.Len(t, stack, 2)
	assert.Equal(t, p, res.(SlotPointer))
	assert.Equal(t, PersistentKey(3), tree.GetRoot().GetKeyAt(0))
}

func TestPersistentEvery_Inserted_Should_Be_Found(t *testing.T) {
	tree := NewBtreeWithPager(80, NoopPersistentPager{KeySerializer: &PersistentKeySerializer{}, KeySize: 8})
	for i := 0; i < 10000; i++ {
		tree.Insert(PersistentKey(i), SlotPointer{
			PageId:  int64(i),
			SlotIdx: int16(i),
		})
	}

	for i := 0; i < 10000; i++ {
		val := tree.Find(PersistentKey(i))
		if val == nil {
			print("")
			tree.Print()
			val = tree.Find(PersistentKey(i))
		}
		assert.Equal(t, SlotPointer{
			PageId:  int64(i),
			SlotIdx: int16(i),
		}, val.(SlotPointer))
	}
}

func TestPersistentInsert_Or_Replace_Should_Return_False_When_Key_Exists(t *testing.T) {
	tree := NewBtreeWithPager(3, NoopPersistentPager{KeySerializer: &PersistentKeySerializer{}, KeySize: 8})
	for i := 0; i < 1000; i++ {
		tree.Insert(PersistentKey(i), SlotPointer{
			PageId:  int64(i),
			SlotIdx: int16(i),
		})
	}

	isInserted := tree.InsertOrReplace(PersistentKey(500), SlotPointer{
		PageId:  int64(1500),
		SlotIdx: int16(1500),
	})

	assert.False(t, isInserted)
}

func TestPersistentInsert_Or_Replace_Should_Replace_Value_When_Key_Exists(t *testing.T) {
	tree := NewBtree(3)
	for i := 0; i < 1000; i++ {
		tree.Insert(MyInt(i), strconv.Itoa(i))
	}

	tree.InsertOrReplace(MyInt(500), "new_500")
	val := tree.Find(MyInt(500))

	assert.Equal(t, "new_500", val.(string))
}
