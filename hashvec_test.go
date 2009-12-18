// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashmap

import "testing"

func BenchmarkHashVectorPush(b *testing.B) {
	b.StopTimer()
	var m hashVector
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.push(HashPair{Integer(i), true})
	}
}

func BenchmarkHashVectorPop(b *testing.B) {
	b.StopTimer()
	var m hashVector
	for i := 0; i < b.N; i++ {
		m.push(HashPair{Integer(i), true})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.pop(0)
	}
}
