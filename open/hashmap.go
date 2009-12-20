// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The hashmap package re-implements Go's builtin map type.
package hashmap

//import "fmt"

const loadGrow = 0.5
const loadShrink = 0.1

// Hashable is an interface that keys have to implement.
type Hashable interface {
	Hash() uint
	Equal(other Hashable) bool
}

// HashMap is the container itself.
// You must call Init() before using it.
type HashMap struct {
	buckets bucketArray
	count int // to compute load factor
}

// HashPair is a key and a value.
// Iter() yields HashPairs.
type HashPair struct {
	Key Hashable
	Value interface{}
}

func (self *HashMap) loadFactor() float {
//	fmt.Printf("loadFactor %d/%d\n", self.count, len(self.buckets.data))
	return float(self.count) / float(len(self.buckets.data))
}

func (self *HashMap) rehashInto(dest bucketArray) {
//	fmt.Printf("rehashInto %d\n", len(data))
	for _, b := range self.buckets.data {
		if  b.state == used {
			dest.push(b.pair.Key, b.pair.Value)
		}
	}
}

func (self *HashMap) grow() {
//	fmt.Printf("grow\n")
	var newBuckets bucketArray
	newBuckets.data = make([]bucket, len(self.buckets.data)*2)
	self.rehashInto(newBuckets)
	self.buckets = newBuckets
}

func (self *HashMap) shrink() {
//	fmt.Printf("shrink\n")
	var newBuckets bucketArray
	newBuckets.data = make([]bucket, len(self.buckets.data)/2)
	self.rehashInto(newBuckets)
	self.buckets = newBuckets
}

// Init initializes or clears a HashMap.
func (self *HashMap) Init() *HashMap {
//	fmt.Printf("Init %s\n", self)
	self.buckets.data = make([]bucket, 8)
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

	if !self.buckets.push(key, value) {
		panic("HashMap.Insert: duplicate key")
	}
	self.count++
}

func (self *HashMap) Remove(key Hashable) {
//	fmt.Printf("Remove %s\n", key)
	if !self.buckets.pop(key) {
		panic("HashMap.Remove: key not found")
	}
	self.count--

	if self.loadFactor() <= loadShrink {
		self.shrink()
	}
}

func (self *HashMap) At(key Hashable) interface{} {
//	fmt.Printf("At %s\n", key)
	b := self.buckets
	p := b.find(key)
	if b.data[p].state != used {
		panic("HashMap.At: key not found")
	}
	return b.data[p].pair.Value
}

func (self *HashMap) Set(key Hashable, value interface{}) {
//	fmt.Printf("Set %s->%s\n", key, value)
	b := self.buckets
	p := b.find(key)
	if b.data[p].state != used {
		panic("HashMap.Set: key not found")
	}
	b.data[p].pair.Value = value
}

func (self *HashMap) Has(key Hashable) bool {
//	fmt.Printf("Has %s\n", key)
	b := self.buckets
	p := b.find(key)
	return b.data[p].state == used;
}

func (self *HashMap) Len() int {
//	fmt.Printf("Len %d\n", self.count)
	return self.count
}

func (self *HashMap) Do(f func(key Hashable, value interface{})) {
//	fmt.Printf("Do %s\n", f)
	for _, b := range self.buckets.data {
		if  b.state == used {
			p := b.pair
			f(p.Key, p.Value)
		}
	}
}

func (self *HashMap) iterate(c chan<- interface{}) {
//	fmt.Printf("Iterate %s\n", c)
	for _, b := range self.buckets.data {
		if  b.state == used {
			c <- b.pair
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
