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
	data	[]hashVector // each should be short
	count	int // to compute load factor
}

// HashPair is a key and a value.
// Iter() yields HashPairs.
type HashPair struct {
	Key Hashable
	Value interface{}
}

func (self *HashMap) loadFactor() float {
//	fmt.Printf("loadFactor %d/%d\n", self.count, len(self.data))
	return float(self.count) / float(len(self.data))
}

func (self *HashMap) rehashInto(data []hashVector) {
//	fmt.Printf("rehashInto %d\n", len(data))
	l := uint(len(data))
	for b := range self.data {
		if self.data[b].count > 0 && self.data[b].data != nil {
			for i := 0; i < self.data[b].count; i++ {
				e := self.data[b].data[i]
				h := e.Key.Hash() % l
				data[h].push(e)
			}
		}
	}
}

func (self *HashMap) grow() {
//	fmt.Printf("grow\n")
	d := make([]hashVector, len(self.data)*2)
	self.rehashInto(d)
	self.data = d
}

func (self *HashMap) shrink() {
//	fmt.Printf("shrink\n")
	d := make([]hashVector, len(self.data)/2)
	self.rehashInto(d)
	self.data = d
}

func (self *HashMap) find(key Hashable) (bucket int, position int) {
//	fmt.Printf("find %s\n", key)
	h := key.Hash() % uint(len(self.data))
	p := self.data[h].find(key)
	return int(h), p
}

// Init initializes or clears a HashMap.
func (self *HashMap) Init() *HashMap {
//	fmt.Printf("Init %s\n", self)
	self.data = make([]hashVector, 8)
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

	bucket, position := self.find(key)
	if position != -1 {
		panic("HashMap.Insert: duplicate key")
	}

	self.data[bucket].push(HashPair{key, value})
	self.count++
}

func (self *HashMap) Remove(key Hashable) {
//	fmt.Printf("Remove %s\n", key)
	bucket, position := self.find(key)
	if position == -1 {
		panic("HashMap.Remove: key not found")
	}

	self.data[bucket].pop(position)
	self.count--

	if self.loadFactor() <= loadShrink {
		self.shrink()
	}
}

func (self *HashMap) At(key Hashable) interface{} {
//	fmt.Printf("At %s\n", key)
	bucket, position := self.find(key)
	if position == -1 {
		panic("HashMap.At: key not found")
	}
	e := self.data[bucket].data[position]
	return e.Value
}

func (self *HashMap) Set(key Hashable, value interface{}) {
//	fmt.Printf("Set %s->%s\n", key, value)
	bucket, position := self.find(key)
	if position == -1 {
		panic("HashMap.Set: key not found")
	}
	self.data[bucket].data[position].Value = value
}

func (self *HashMap) Has(key Hashable) bool {
//	fmt.Printf("Has %s\n", key)
	_, position := self.find(key)
	return position != -1
}

func (self *HashMap) Len() int {
//	fmt.Printf("Len %d\n", self.count)
	return self.count
}

func (self *HashMap) Do(f func(key Hashable, value interface{})) {
//	fmt.Printf("Do %s\n", f)
	for b := range self.data {
		if self.data[b].count > 0 {
			for i := 0; i < self.data[b].count; i++ {
				e := self.data[b].data[i]
				f(e.Key, e.Value)
			}
		}
	}
}

func (self *HashMap) iterate(c chan<- interface{}) {
//	fmt.Printf("Iterate %s\n", c)
	for b := range self.data {
		if self.data[b].count > 0 {
			for i := 0; i < self.data[b].count; i++ {
				e := self.data[b].data[i]
				c <- e
			}
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
