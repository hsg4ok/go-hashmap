// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A specialized container/vector clone for hash buckets.
//
// These never shrink: If they empty out, the hash table
// will hopefully shrink/rehash soon anyway.

package hashmap

const initialLength = 2

type hashVector struct {
	data []HashPair
	count int
}

func (self *hashVector) find(key Hashable) int {
	d := self.data
	if d != nil {
		l := self.count
		for i := 0; i < l; i++ {
			if key.Equal(d[i].key) {
				return i
			}
		}
	}
	return -1
}

func (self *hashVector) grow() {
	d := make([]HashPair, len(self.data)*2)
	copy(d, self.data)
	self.data = d
}

func (self *hashVector) push(pair HashPair) {
	d := self.data
	if d == nil {
		// lazy: avoid allocation for empty buckets
		// small: assuming good hash function
		self.data = make([]HashPair, initialLength)
		d = self.data
	}

	c := self.count
	if c == len(d) {
		self.grow();
		d = self.data
	}

	d[c] = pair
	self.count++
}

func (self *hashVector) pop(i int) {
	d := self.data
	copy(d[i:], d[i+1:]) // explicit loop does worth despite slice allocation
	self.count--
}
