// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The hashmap package re-implements Go's builtin map type.
package hashmap

//import "fmt"

// These seem right, Java's lower 0.75 bound resizes too
// much, a higher 1.15 or 1.25 bound makes chains grow
const loadGrow = 1.0
const loadShrink = 0.25

// Hashable is an interface that keys have to implement.
type Hashable interface {
	Hash() uint
	Equal(other Hashable) bool
}

// HashMap is the container itself.
// You must call Init() before using it.
type HashMap struct {
	data	[]*bucket // each should be short
	count	int // to compute load factor
}

// HashPair is a key and a value.
// Iter() yields HashPairs.
type HashPair struct {
	key Hashable
	value interface{}
}

type bucket struct {
	hp HashPair
	next *bucket
}

func (self *HashMap) loadFactor() float {
//	fmt.Printf("loadFactor %d/%d\n", self.count, len(self.data))
	return float(self.count) / float(len(self.data))
}

func (self *HashMap) rehashInto(data []*bucket) {
//	fmt.Printf("rehashInto %d\n", len(data))
	for _, b := range self.data {
		for n := b; n != nil; n = n.next {
			e := n.hp
			h := e.key.Hash() % uint(len(data))
			x := &bucket{e, data[h]}
			data[h] = x
		}
	}
}

func (self *HashMap) grow() {
//	fmt.Printf("grow\n")
	d := make([]*bucket, len(self.data)*2)
	self.rehashInto(d)
	self.data = d
}

func (self *HashMap) shrink() {
//	fmt.Printf("shrink\n")
	d := make([]*bucket, len(self.data)/2)
	self.rehashInto(d)
	self.data = d
}

func (self *HashMap) find(key Hashable) (b int, position *bucket, prev *bucket) {
//	fmt.Printf("find %s\n", key)
	h := key.Hash() % uint(len(self.data))
	for n := self.data[h]; n != nil; prev, n = n, n.next {
		if key.Equal(n.hp.key) {
			return int(h), n, prev
		}
	}
	return int(h), nil, prev
}

// Init initializes or clears a HashMap.
func (self *HashMap) Init() *HashMap {
//	fmt.Printf("Init %s\n", self)
	self.data = make([]*bucket, 8)
	self.count = 0
	return self
}

// New returns an initialized hashmap.
func New() *HashMap {
//	fmt.Printf("New\n")
	return new(HashMap).Init()
}

func (self *HashMap) Insert(key Hashable, value interface{}) {
//	fmt.Printf("Insert %s->%s\n", key, value)
	if self.loadFactor() >= loadGrow {
		self.grow()
	}

	b, position, _ := self.find(key)
	if position != nil {
		panic("HashMap.Insert: duplicate key")
	}

	head := self.data[b]
	node := &bucket{HashPair{key, value}, head}
	self.data[b] = node
	self.count++
}

func (self *HashMap) Remove(key Hashable) {
//	fmt.Printf("Remove %s\n", key)
//	fmt.Printf("%s\n", self)
	b, position, prev := self.find(key)
	if position == nil {
		panic("HashMap.Remove: key not found")
	}

	if prev == nil && position.next == nil {
		self.data[b] = nil
	} else if prev == nil {
		self.data[b] = position.next
	} else {
		prev.next = position.next
	}
	self.count--

	if self.loadFactor() <= loadShrink {
		self.shrink()
	}
}

func (self *HashMap) At(key Hashable) interface{} {
//	fmt.Printf("At %s\n", key)
	_, position, _ := self.find(key)
	if position == nil {
		panic("HashMap.At: key not found")
	}
	return position.hp.value
}

func (self *HashMap) Set(key Hashable, value interface{}) {
//	fmt.Printf("Set %s->%s\n", key, value)
	_, position, _ := self.find(key)
	if position == nil {
		panic("HashMap.Set: key not found")
	}
	position.hp.value = value
}

func (self *HashMap) Has(key Hashable) bool {
//	fmt.Printf("Has %s\n", key)
	_, position, _ := self.find(key)
	return position != nil
}

func (self *HashMap) Len() int {
//	fmt.Printf("Len %d\n", self.count)
	return self.count
}

func (self *HashMap) Do(f func(key Hashable, value interface{})) {
//	fmt.Printf("Do %s\n", f)
}

func (self *HashMap) iterate(c chan<- interface{}) {
//	fmt.Printf("Iterate %s\n", c)
	for _, b := range self.data {
		for n := b; n != nil; n = n.next {
			e := n.hp
			c <- e
		}
	}
	close(c)
}

func (self *HashMap) Iter() <-chan interface{} {
//	fmt.Printf("Iter\n")
	c := make(chan interface{})
	go self.iterate(c)
	return c
}

func (self *HashMap) String() string {
	s := "{"
	for r := range self.Iter() {
		q := r.(HashPair)
		s = s + fmt.Sprintf("%s: %s, ", q.key, q.value)
	}
	s = s + "}"
	return s
}
