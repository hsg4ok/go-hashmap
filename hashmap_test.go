// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashmap

import "testing"

func TestLen(t *testing.T) {
	var m HashMap
	m.Len();
}

func BenchmarkLen(b *testing.B) {
	b.StopTimer()
	var m HashMap
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Len();
	}
}
