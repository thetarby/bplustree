## What is this ?
This is a b+ tree implementation in Golang which is written keeping persistency in mind. All the logic about persistency is abstracted away under persistence.go

Pager interface type encapsulates methods to achieve persistency. In this package only in memory or mock implementations of the Pager exists. To check a disk-persistent implementation, you can have a look at buffer pool pager implementation [here.](https://github.com/thetarby/helindb)


## Usage
Btree type exposes these five methods which are useful.
```go
func (tree *BTree) Find(key Key) interface{}

func (tree *BTree) InsertOrReplace(key Key, value interface{}) (isInserted bool)

func (tree *BTree) Insert(key Key, value interface{})

func (tree *BTree) Delete(key Key) bool

func (tree *BTree) Print()
```

To create a b+ tree one must first create a pager instance.

```go
// first parameter is degree of the tree. Second parameter is the pager.
tree := NewBtreeWithPager(3, 
		NewNoopPagerWithValueSize(&StringKeySerializer{Len: 10}, 
			&StringValueSerializer{Len: 10},
		),
	)

for i := 0; i < 1000; i++ {
	tree.Insert(StringKey(strconv.Itoa(i)), fmt.Sprintf("value_%v", i))
}

val := tree.Find(StringKey("500")) // val is "value_500"
```

Pager needs a ValueSerializer and KeySerializer. These are interfaces that are used to serialize a key or a value. There are some types (StringKeySerializer, StringValueSerializer etc...) already implementing these interfaces but more could be added to support other types as keys or values such as dates. All keys in a b+ tree instance should have the same byte length when serialized. This is required to be able to do binary search in a node.

Also all values are enforced to have the same length by the ValueSerializer type but implementation could easily be tweaked to support variable length values since they only exists in leaf nodes. But I kept them fixed size here since I do not think for now it would be useful to have variable length values.

More examples are in `*_test.go` files.

## Tests

To run tests
```sh
go test ./... -v
```

## License
[MIT](https://choosealicense.com/licenses/mit/)