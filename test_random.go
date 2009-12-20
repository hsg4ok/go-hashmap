package main

import hashmap "container/hashmap/open"
import "rand"
import "fmt"
import "time"

type Integer int;

func (self Integer) Hash() uint { return uint(self*self) }
func (self Integer) Equal(other hashmap.Hashable) bool { return self == other.(Integer) }

const N = 300000
const S = 300000

func main() {
	duplicate := 0
	unknown := 0

	rand.Seed(time.Nanoseconds())

	m := hashmap.New()

	for i := 0; i < N; i++ {
		r := rand.Intn(S)
		if m.Has(Integer(r)) {
			duplicate++
		} else {
			m.Insert(Integer(r), true)
		}
	}

	for i := 0; i < N; i++ {
		r := rand.Intn(S)
		if m.Has(Integer(r)) {
			m.Remove(Integer(r))
		} else {
			unknown++
		}
	}

	fmt.Printf("%d insertions, %d duplicates\n", N, duplicate)
	fmt.Printf("%d removals, %d unknown\n", N, unknown)
	fmt.Printf("%d left\n", m.Len())
}
