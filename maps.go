package maps

// A brief study in interface design.
// Comments welcome!

// Interfaces modeling key properties.

type Any interface {
	Equal(other Any) bool
}

type Ordered interface {
	Any
	Less(other Ordered) bool
	Greater(other Ordered) bool
}

type Hashable interface {
	Any
	Hash() uint
}

// Alternative for comparisons, Equal(), Less(), and
// Greater() can be done in terms of this. (Not used
// in the map interfaces.)

type Comparable interface {
	Compare(other Comparable) int
}

// Operations that can be made to apply to any kind
// of map. Len() is obvious, Iter() will return some
// kind of "key/value pair" depending on the concrete
// map implementation, Do() requires clients to cast
// back to their concrete types anyway so separating
// it out by specific map doesn't help much.

type CommonMap interface {
	Len() int
        Iter() <-chan interface{}
	Do(f func(key interface{}, value interface{}))
}

// Simplistic maps.

type SimpleMap interface {
	CommonMap
	Insert(key Any, value interface{})
	Remove(key Any)
	At(key Any) interface{}
	Set(key Any, value interface{})
	Has(key Any) bool
}

// All maps implemented as some kind of hash table.

type HashMap interface {
	CommonMap
	Insert(key Hashable, value interface{})
	Remove(key Hashable)
	At(key Hashable) interface{}
	Set(key Hashable, value interface{})
	Has(key Hashable) bool
}

// All maps implemented as some kind of ordered
// structure (e.g. some kind of tree).

type OrderedMap interface {
	CommonMap
	Insert(key Ordered, value interface{})
	Remove(key Ordered)
	At(key Ordered) interface{}
	Set(key Ordered, value interface{})
	Has(key Ordered) bool
}
