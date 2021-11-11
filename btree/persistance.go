package btree

/* InternalNode and LeafNode structures should extend a PersistentPage implementation to be able to be disk persistent */

type PersistentPage interface {
	GetData() []byte

	// GetPageId returns the page_id of the physical page.
	GetPageId() Pointer
}

type Pager interface {
	// NewInternalNode first should create a PersistentPage which points to a byte array.
	// Then initialize an InternalNode structure.
	// Finally, it should serialize the structure on to pointed byte array.
	// NOTE: the node should have a reference(by extending it for example) to the created PersistentPage
	// so that it can be serialized in the future when its state changes.
	NewInternalNode() *InternalNode

	// NewLeafNode first should create an PersistentPage which points to a byte array.
	// Then initialize an LeafNode structure.
	// Finally, it should serialize the structure on to pointed byte array
	NewLeafNode() *LeafNode

	// SyncInternalNode should be called after every change in an internal node and it synces the serialized data with
	// the structure.
	SyncInternalNode(serializer *InternalNode)

	// SyncLeafNode should be called after every change in a leaf node.
	SyncLeafNode(serializer *LeafNode)

	// GetNode returns a Node given a Pointer. Should be able to deserialize a node from byte arr and should be able to
	// recognize if it is an InternalNode or LeafNode and return the correct type.
	GetNode(p Pointer) Node
}

type InternalNodeSerializer interface {
	Serialize(node *InternalNode, dest []byte)
}

type LeafNodeSerializer interface {
	Serialize(node *LeafNode, dest []byte)
}

/* NOOP IMPLEMENTATION*/

type NoopPersistentPage struct {
	pageId Pointer
}

func (n NoopPersistentPage) GetData() []byte {
	panic("implement me")
}

func (n NoopPersistentPage) GetPageId() Pointer {
	return n.pageId
}

type NoopPager struct {
	internalNodeSerializer InternalNodeSerializer
	leafNodeSerializer     LeafNodeSerializer
}

var lastPageId Pointer = 0

func (b *NoopPager) NewInternalNode() *InternalNode {
	// TODO: create persistent page from buffer pool
	lastPageId++
	i := InternalNode{
		PersistentPage: NoopPersistentPage{lastPageId},
		Keys:           make([]Key, 0, 2),
		Pointers:       make([]Pointer, 0, 2),
	}
	mapping[lastPageId] = &i

	//b.internalNodeSerializer.Serialize(&i, i.GetData())

	return &i
}

func (b *NoopPager) NewLeafNode() *LeafNode {
	// TODO: create persistent page from buffer pool
	lastPageId++
	l := LeafNode{
		PersistentPage: NoopPersistentPage{lastPageId},
		Keys:           make([]Key, 0, 2),
		Values:         make([]interface{}, 0, 2),
	}
	mapping[lastPageId] = &l

	//b.leafNodeSerializer.Serialize(&l, l.GetData())

	return &l
}

func (b *NoopPager) SyncInternalNode(n *InternalNode) {
	//b.internalNodeSerializer.Serialize(n, n.GetData())
}

func (b *NoopPager) SyncLeafNode(n *LeafNode) {
	//b.leafNodeSerializer.Serialize(n, n.GetData())
}

var mapping = make(map[Pointer]Node)

func (b *NoopPager) GetNode(p Pointer) Node {
	return mapping[p]
}

/* REAL IMPLEMENTATION */
