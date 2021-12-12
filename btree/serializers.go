package btree

import (
	"bytes"
	"encoding/binary"
)

// KeySerializer is the interface a b+ tree uses to serialize keys in a b+ tree. A type should have a corresponding
// KeySerializer to be used as key in b+tree. All keys of the same type should have the same byte length when they are
// serialized. That behaviour is enforced by Size method.
type KeySerializer interface {
	Serialize(key Key) ([]byte, error)
	Deserialize([]byte) (Key, error)

	// Size return byte length of the serialized key.
	Size() int
}

type PersistentKeySerializer struct{}

func (p *PersistentKeySerializer) Serialize(key Key) ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, key.(PersistentKey))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *PersistentKeySerializer) Deserialize(data []byte) (Key, error) {
	reader := bytes.NewReader(data)
	var key PersistentKey
	err := binary.Read(reader, binary.BigEndian, &key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (p *PersistentKeySerializer) Size() int {
	return 10
}

type StringKeySerializer struct {
	Len int
}

func (s *StringKeySerializer) Serialize(key Key) ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, ([]byte)(key.(StringKey)))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *StringKeySerializer) Deserialize(data []byte) (Key, error) {
	return StringKey(data[:s.Len]), nil
}

func (s *StringKeySerializer) Size() int {
	return s.Len
}

// ValueSerializer is very similar to KeySerializer. A type should have a ValueSerializer implemented to be
// used as a value in b+ tree. All values serialized by a ValueSerializer should also have the same length although
// that behaviour can be changed easily by slightly modifying the implementation since values are only stored in
// leaf nodes.
type ValueSerializer interface {
	Serialize(val interface{}) ([]byte, error)
	Deserialize([]byte) (interface{}, error)
	Size() int
}

type StringValueSerializer struct {
	Len int
}

func (s *StringValueSerializer) Serialize(val interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, ([]byte)(val.(string)))
	if err != nil {
		return nil, err
	}
	res := make([]byte, s.Len)
	copy(res, buf.Bytes())
	return res, nil
}

func (s *StringValueSerializer) Deserialize(data []byte) (interface{}, error) {
	return string(data[:s.Len]), nil
}

func (s *StringValueSerializer) Size() int {
	return s.Len
}

type SlotPointerValueSerializer struct {
}

func (s *SlotPointerValueSerializer) Serialize(val interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.BigEndian, val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *SlotPointerValueSerializer) Deserialize(data []byte) (interface{}, error) {
	reader := bytes.NewReader(data)
	var val SlotPointer
	err := binary.Read(reader, binary.BigEndian, &val)

	return val, err
}

func (s *SlotPointerValueSerializer) Size() int {
	return 10
}
