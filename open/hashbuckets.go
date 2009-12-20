// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashmap

// Buckets start out fresh; once they contain a
// pair they become used; once their pair is
// removed they become deleted.
const (
	fresh = iota;
	used;
	deleted;
)

// Buckets are just HashPairs and a state.
type bucket struct {
	pair HashPair
	state int
}

// Special bucket written to array for deleted buckets.
var deletedBucket bucket = bucket{HashPair{nil, nil}, deleted}

// All we do is wrap a bucket slice. Yes there's a pointer
// indirection, sue me.
type bucketArray struct {
	data []bucket;
}

// Find the bucket index for the given key. If the bucket
// used, we found an exact match; if the bucket is fresh
// or deleted, we can insert there. Note that we relocate
// buckets on successful searches (if possible).
func (self bucketArray) find(key Hashable) (index int) {
	d := self.data
	l := uint(len(d))
	h := key.Hash() % l
	i := h // current probe
	j := uint(1) // next probe offset
	r := -1 // relocation index
	for {
		b := d[i]

		switch b.state {
		case fresh:
			// key not found
			if r != -1 {
				return r
			}
			return int(i)
		case deleted:
			// remember if it's the first
			if r == -1 {
				r = int(i)
			}
		case used:
			// winner if keys match
			if key.Equal(b.pair.Key) {
				// relocate from i to r if possible
				if r != -1 {
					d[r] = b
					d[i] = deletedBucket
					return r
				}
				// otherwise just return it
				return int(i)
			}
		}

		// next index, wrapping around
		i = (i + j) % l
		// XXX tried +2 but that's worse, probably because
		// we trash into other bucket-ranges then
		// XXX tried j = 2*j but that's at least not better,
		// not sure why; "exponential probing"?
		// XXX tried i = (i + uint(j*j)) % l; j++ which is
		// quadratic probing, and that's better, by about
		// half a second in fact for example_hashmap; but
		// I am not sure about the termination properties
		// yet...

		// back to where we started?
		if i == h {
			if r != -1 {
				// if we have a deleted one, return that
				return r
			} else {
				// table full, should never happen
				panic("bucketArray.find: table full")
			}
		}
	}

	// unreachable
	panic("bucketArray.find: unreachable")
	return -1
}

func (self bucketArray) push(key Hashable, value interface{}) bool {
	p := self.find(key)
	if self.data[p].state == used {
		return false
	}
	self.data[p] = bucket{HashPair{key, value}, used}
	return true
}

func (self bucketArray) pop(key Hashable) bool {
	p := self.find(key)
	if self.data[p].state != used {
		return false
	}
	self.data[p] = deletedBucket
	return true
}
