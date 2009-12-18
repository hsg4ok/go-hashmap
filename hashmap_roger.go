// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The hashmap package re-implements Go's builtin map type.
package hashmap

// Hashable is an interface that keys have to implement.
type Hashable interface {
	Hash() uint
	Equal(other Hashable) bool
}

type bucket struct {
	HashPair
	next *bucket
}

// HashMap is the container itself.
// You must call Init() before using it.
type HashMap struct {
	m     map[uint]*bucket
	count int
}

// HashPair is a key and a value.
// Iter() yields HashPairs.
type HashPair struct {
	Key   Hashable
	Value interface{}
}

// Init initializes or clears a HashMap.
func (h *HashMap) Init() { h.m = make(map[uint]*bucket) }

// New returns an initialized hashmap.
func New() *HashMap {
	var h HashMap
	h.Init()
	return &h
}

func (h *HashMap) find(key Hashable, insert bool) (fb *bucket, found bool) {
	hash := key.Hash()
	b, ok := h.m[hash]
	if ok {
		for p := b; p != nil; p = p.next {
			if p.Key.Equal(key) {
				return p, true
			}
		}
		if insert {
			fb := &bucket{HashPair{key, nil}, b.next}
			b.next = fb
			h.count++
			return fb, false
		}
	} else {
		if insert {
			b = &bucket{HashPair{key, nil}, nil}
			h.m[hash] = b
			h.count++
			return b, false
		}
	}
	return nil, false
}

func (h *HashMap) At(key Hashable) (value interface{}) {
	if b, _ := h.find(key, false); b != nil {
		value = b.Value
	}
	return
}

func (h *HashMap) Insert(key Hashable, value interface{}) {
	b, found := h.find(key, true)
	if found {
		panic("HashMap.Insert: duplicate key")
	}
	b.Value = value
}

func (h *HashMap) Remove(key Hashable) {
	var prev, p *bucket

	hash := key.Hash()
	b := h.m[hash]
	for p = b; p != nil; prev, p = p, p.next {
		if p.Key.Equal(key) {
			break
		}
	}
	switch {
	case p == nil:
		panic("HashMap.Remove: key not found")

	case prev == nil:
		if p.next == nil {
			h.m[hash] = nil, false
		} else {
			*p = *p.next
		}

	default:
		prev.next = p.next
	}
	h.count--
}

func (h *HashMap) Set(key Hashable, value interface{}) {
	b, _ := h.find(key, false)
	if b == nil {
		panic("HashMap.Set: key not found")
	}
	b.Value = value
}

func (h *HashMap) Has(key Hashable) (found bool) {
	_, found = h.find(key, false)
	return
}

func (h *HashMap) Len() int { return h.count }

func (h *HashMap) Do(f func(key Hashable, value interface{})) {
	for _, b := range h.m {
		for ; b != nil; b = b.next {
			f(b.Key, b.Value)
		}
	}
}

func (h *HashMap) iterate(c chan<- interface{}) {
	for _, b := range h.m {
		for ; b != nil; b = b.next {
			c <- b.HashPair
		}
	}
	close(c)
}

func (h *HashMap) Iter() <-chan interface{} {
	c := make(chan interface{})
	go h.iterate(c)
	return c
}
