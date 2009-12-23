package main

import "fmt"
import "container/hashmap"
import "io/ioutil"
import "strings"

type String string

func (self String) Hash() uint {
	var h uint = 5381
	// explicit for loop is slower than range
	for _, r := range self {
		h = (h << 5) + h + uint(r)
		// h = (h << 5) + h ^ uint(r)
	}
	return h
}

func (self String) sillyHash() uint {
	var h uint;
	for _, r := range self {
		h = h + uint(r)
	}
	return h;
}

func (self String) Equal(other hashmap.Hashable) bool {
	s := other.(String)
	return self == s;
}

func main() {
	raw, error := ioutil.ReadFile("/usr/share/dict/cracklib-words")
	if error == nil {
		data := string(raw)
		words := strings.Split(data, "\n", 0)
		dict := hashmap.New()
		for _, w := range words {
			dict.Insert(String(w), true)
		}
		fmt.Printf("%d words\n", dict.Len())
	}
}
