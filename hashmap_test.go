// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashmap

import "testing"

type Integer int;

func (self Integer) Hash() uint { return uint(self*self) }
func (self Integer) Equal(other Hashable) bool { return self == other.(Integer) }

func TestZeroLen(t *testing.T) {
	a := New()
	if a.Len() != 0 {
		t.Errorf("expected 0, got %d", a.Len())
	}
}

func TestZeroLookup(t *testing.T) {
	const Len = 10000
	a := New()
	for i := 0; i < Len; i++ {
		if a.Has(Integer(i)) {
			t.Errorf("found %d in empty hashmap", i)
		}
	}
}

func TestInsert(t *testing.T) {
	const Len = 10000
	a := New()
	for i := 0; i < Len; i++ {
		a.Insert(Integer(i), i)
	}
	for i := 0; i < Len; i++ {
		if !a.Has(Integer(i)) {
			t.Errorf("inserted %d not found", i)
		}
	}
}

func TestRemove(t *testing.T) {
	const Len = 10000
	a := New()
	for i := 0; i < Len; i++ {
		a.Insert(Integer(i), i)
	}
	for i := 0; i < Len; i++ {
		a.Remove(Integer(i))
	}
	for i := 0; i < Len; i++ {
		if a.Has(Integer(i)) {
			t.Errorf("removed %d was found", i)
		}
	}
}

func TestIter(t *testing.T) {
	const Len = 100
	x := New()
	for i := 0; i < Len; i++ {
		x.Insert(Integer(i), i*i)
	}
	i := 0
	for v := range x.Iter() {
		p := v.(HashPair)
		key := p.key.(Integer)
		val := p.value.(int)
		if key*key != Integer(val) {
			t.Error("Iter expected", key*key, "got", val)
		}
		i++
	}
	if i != Len {
		t.Error("Iter stopped at", i, "not", Len)
	}
}

func BenchmarkLen(b *testing.B) {
	b.StopTimer()
	m := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Len()
	}
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	m := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Insert(Integer(i), true);
	}
}

func BenchmarkRemove(b *testing.B) {
	b.StopTimer()
	m := New()
	for i := 0; i < b.N; i++ {
		m.Insert(Integer(i), true);
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(Integer(i));
	}
}

func BenchmarkSuccessfulLookup(b *testing.B) {
	b.StopTimer()
	m := New()
	for i := 0; i < b.N; i++ {
		m.Insert(Integer(i), true);
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Has(Integer(i));
	}
}

func BenchmarkFailedLookup(b *testing.B) {
	b.StopTimer()
	m := New()
	for i := 0; i < b.N; i++ {
		m.Insert(Integer(i), true);
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Has(Integer(-1));
	}
}
