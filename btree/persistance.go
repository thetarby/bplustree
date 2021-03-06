package btree

import (
	"bytes"
	"encoding/binary"
)

// PersistentPage is an interface that InternalNode and SlottedPage structures should implement to be able to be
// disk persistent. It encapsulates methods which would be useful to flush nodes to disk.
type PersistentPage interface {
	GetData() []byte

	// GetPageId returns the page_id of the physical page.
	GetPageId() Pointer
}

// Pager abstracts away the logic about allocating and managing tree nodes. For a truly persistent b+ tree implementation
// a real Pager should be implemented. This package only implements mock or in memory implementations of the Pager.
type Pager interface {
	// NewInternalNode first should create a PersistentPage which points to a byte array.
	// Then initialize an InternalNode structure.
	// Finally, it should serialize the structure on to pointed byte array.
	// NOTE: the node should have a reference(by extending it for example) to the created PersistentPage
	// so that it can be serialized in the future when its state changes.
	NewInternalNode(firstPointer Pointer) Node

	// NewLeafNode first should create an PersistentPage which points to a byte array.
	// Then initialize an LeafNode structure.
	// Finally, it should serialize the structure on to pointed byte array
	NewLeafNode() Node

	// GetNode returns a Node given a Pointer. Should be able to deserialize a node from byte arr and should be able to
	// recognize if it is an InternalNode or LeafNode and return the correct type.
	GetNode(p Pointer) Node

	// Unpin is called when we are done with reading or writing to a tree node. In a truly persistent implementation
	// it could flush nodes to disk or not depending on isDirty parameter. An in memory implementation of this method could
	// be a noop
	Unpin(n Node, isDirty bool)

	// UnpinByPointer does the same thing as Unpin to the node pointed by p Pointer.
	UnpinByPointer(p Pointer, isDirty bool)
}

/* NOOP IMPLEMENTATION*/

type NoopPersistentPage struct {
	pageId Pointer
	data   []byte
}

func NewNoopPersistentPage(pageId Pointer) *NoopPersistentPage {
	return &NoopPersistentPage{
		pageId: pageId,
		data:   make([]byte, 4096, 4096),
	}
}

func (n NoopPersistentPage) GetData() []byte {
	return n.data
}

func (n NoopPersistentPage) GetPageId() Pointer {
	return n.pageId
}

// will be used by noop persistent pager. Making them global is not good but NoopPager is only intented
// for testing purposes
var lastPageId Pointer = 0
var mapping = make(map[Pointer]Node)

type NoopPersistentPager struct {
	KeySerializer   KeySerializer
	ValueSerializer ValueSerializer
}

func (n2 *NoopPersistentPager) UnpinByPointer(p Pointer, isDirty bool) {}

func (n2 *NoopPersistentPager) Unpin(n Node, isDirty bool) {}

func (n *NoopPersistentPager) NewInternalNode(firstPointer Pointer) Node {
	h := PersistentNodeHeader{
		IsLeaf: 0,
		KeyLen: 0,
	}

	// create a new node
	// TODO: should use an adam ak??ll?? pager
	lastPageId++
	node := PersistentInternalNode{PersistentPage: NewNoopPersistentPage(lastPageId), pager: n, keySerializer: n.KeySerializer}

	// write header
	data := node.GetData()
	WritePersistentNodeHeader(&h, data)

	// write first pointer
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, firstPointer)
	CheckErr(err)
	asByte := buf.Bytes()
	copy(data[PersistentNodeHeaderSize:], asByte)

	mapping[lastPageId] = &node
	return &node
}

func (n *NoopPersistentPager) NewLeafNode() Node {
	h := PersistentNodeHeader{
		IsLeaf: 1,
		KeyLen: 0,
	}

	// create a new node
	// TODO: should use an adam ak??ll?? pager
	lastPageId++
	var node PersistentLeafNode
	if n.ValueSerializer == nil {
		node = PersistentLeafNode{PersistentPage: NewNoopPersistentPage(lastPageId), pager: n, keySerializer: n.KeySerializer, valSerializer: &SlotPointerValueSerializer{}}
	} else {
		node = PersistentLeafNode{PersistentPage: NewNoopPersistentPage(lastPageId), pager: n, keySerializer: n.KeySerializer, valSerializer: n.ValueSerializer}

	}

	// write header
	data := node.GetData()
	WritePersistentNodeHeader(&h, data)

	mapping[lastPageId] = &node
	return &node
}

func (n *NoopPersistentPager) GetNode(p Pointer) Node {
	return mapping[p]
}

func NewNoopPager(serializer KeySerializer, valSerializer ValueSerializer) *NoopPersistentPager {
	return &NoopPersistentPager{
		KeySerializer:   serializer,
		ValueSerializer: valSerializer,
	}
}
