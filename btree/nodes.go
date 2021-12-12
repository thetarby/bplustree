package btree

import (
	"sort"
)

/*
  Node structure:

  Leaf node and an internal node behaves the same for all operations in Node interface. Nodes consist of key-value
  pairs. Logically we could say each key is associated with the value right after itself. Internal Node is treated
  as only a slightly special case of a Leaf Node which has another value attached to its beginning with no associated key.

  A leaf node:
  ----------------------------------------------------------------------------
  | Key_0 | Val_0 | Key_1 | Val_1 | Key_2 | Val_2 | Key_3 | Val_3 | Key_4 | Val_4 |
  ----------------------------------------------------------------------------
      ^ key at 0th index      ^ value 1st index

  First pair is Key_0-Val_0, second is Key_1-Val_1

  An internal node:
  ----------------------------------------------------------------------------
  | Val_0 | Key_0 | Val_1 | Key_1 | Val_2 | Key_2 | Val_3 | Key_3 | Val_4 | Key_4 | Val_5 |
  ----------------------------------------------------------------------------
     ^ value at 0th index     ^ key at 1st index

  pairs are Key_0-Val_1, Key_1-Val_2 notice that the value indexes are one more than key indexes for internal
  nodes because of the extra value at the beginning.
*/

/*
  Splitting a node:

  ----------------------------------------------------------------------------
  | Key_1 | Val_1 | Key_2 | Val_2 | Key_3 | Val_3 | Key_4 | Val_4 | Key_5 | Val_5 |
  ----------------------------------------------------------------------------
                                      ^ split point
  after split, node is split into two nodes namely Left Node and Right Node

	Left Node:
  ----------------------------------------------------------------------------
  | Key_1 | Val_1 | Key_2 | Val_2 |
  ----------------------------------------------------------------------------

   Right Node:
  ----------------------------------------------------------------------------
  | Key_3 | Val_3 | Key_4 | Val_4 | Key_5 | Val_5 |
  ----------------------------------------------------------------------------
*/

/*
  Before redistribute:

  ..................| Pointer_1 | Key_3 | Pointer_2 |............................
                       /                   \
                      /                     \
           | Key_1 | Val_1 |       | Key_3 | Val_3 | Key_4 | Val_4 | Key_5 | Val_5 |


  After redistribute is called:

  ..................| Pointer_1 | Key_4 | Pointer_2 |............................
                       /                   \
                      /                     \
 | Key_1 | Val_1 | Key_3 | Val_3 |       | Key_4 | Val_4 | Key_5 | Val_5 |

 Keys and values in children are redistributed and the key in parent is that is separating the nodes are updated to
 maintain the integrity of the index.

*/

/*
  Before merging two internal nodes:

  ..................| Pointer_1 | Key_3 | Pointer_2 |............................
                       /                   \
                      /                     \
           | Val_0 | Key_0 | Val_1 |      | Val_5 | Key_5 | Val_6 | Key_6 | Val_7 | Key_7 | Val_8 |


  After merging:

  ..................| Pointer_1 |............................
                       /
                      /
 | Val_0 | Key_0 | Val_1 | Key_3 | Val_5 | Key_5 | Val_6 | Key_6 | Val_7 | Key_7 | Val_8 |

 Notice that the key-value pair which points to right node is deleted from parent. DeleteAt method with the
 deleted key's index could be called to delete that pair.

 Also, Key_3 is pushed down to merged node.

 NOTE: In this case in a b+ tree, separating key cannot be equal to a key in the right node(it cannot be Key_5 for above example)
 if right node is internal node as well since then right node would have to point to a smaller key than separating key
 but would violate b+ tree since when you go right you should always find equal or bigger keys.
*/

type Pointer int64

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

type NodeIndexPair struct {
	Node  Pointer
	Index int // pointer index for internal nodes and value index for leaf nodes
}

// Node defines the simple api that a b+ tree node exposes to outer world. Both leaf and internal nodes should
// implement this interface. Methods in this interface are small and is designed to implement a b+ tree easier
// by another upper layer.
type Node interface {
	// findKey returns the index of the key in the node if it is found. If it is not found it returns the index
	// that key would reside if found.
	findKey(key Key) (index int, found bool)

	// shiftKeyValueToRightAt shifts keys and values that comes after nth index to right by amount that a new key-value
	// pair should fit
	shiftKeyValueToRightAt(n int)

	// shiftKeyValueToLeftAt is same as shiftKeyValueToRightAt but to left
	shiftKeyValueToLeftAt(n int)

	// setKeyAt sets the key at given index to key
	setKeyAt(idx int, key Key)

	// setValueAt sets the value at given index to value
	setValueAt(idx int, val interface{})

	// GetKeyAt returns the key in the given index
	GetKeyAt(idx int) Key

	// GetValueAt returns the value in the given index
	GetValueAt(idx int) interface{}
	GetValues() []interface{}

	// SplitNode splits the node it is called into two nodes at the given index. Split is done so that the key at the given
	// index is moved to newly created node along with every key and value after itself. All keys and values that
	// comes before that key stays in the current node and current node is truncated after split.
	// right is a pointer for the newly created node, keyAtLeft is the last key of the current node and keyAtRight is
	// the first key in newly created node
	SplitNode(index int) (right Pointer, keyAtLeft Key, keyAtRight Key)
	PrintNode()
	IsOverFlow(degree int) bool

	// InsertAt inserts the given key value pair after the given index.
	InsertAt(index int, key Key, val interface{})

	// DeleteAt deletes the key and value right after the key from the node. Cannot be used to delete first value
	// in an internal node since first value does not have an associated key, hence a key index to pass to DeleteAt
	DeleteAt(index int)

	GetPageId() Pointer
	IsLeaf() bool
	GetHeader() *PersistentNodeHeader
	SetHeader(*PersistentNodeHeader)

	// IsSafeForSplit returns true if there is at least one empty place in the node meaning it
	// won't split even one key is inserted
	IsSafeForSplit(degree int) bool

	// IsSafeForMerge returns true if it is more than half full meaning it won't underflow and merge even
	// one key is deleted
	IsSafeForMerge(degree int) bool

	/* delete related methods */

	// Keylen returns number of keys in the node
	Keylen() int

	// GetRight is meant only for leaf nodes which has a right pointer
	GetRight() Pointer

	// MergeNodes merges the node it is called on with its parameter rightNode. Merging two leaf nodes is trivial.
	// Directly appending key-value pairs of rightNode to leftNode is enough. And the pointer which points to the rightNode
	// in the parent should be deleted from parent with its value.
	// Merging two internal nodes is a little different. The key in parent separating the right and left nodes should be
	// pushed down and placed in between right and left nodes. And it should be deleted from parent with its value as well.
	MergeNodes(rightNode Node, parent Node)

	// Redistribute is called for a parent node and its two children which are adjacent to each other. It redistributes
	// keys and values between children and updates the key separating them in the parent.
	Redistribute(rightNode_ Node, parent_ Node)

	IsUnderFlow(degree int) bool
}
