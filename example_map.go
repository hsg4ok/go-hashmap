package main

import "fmt"
import "io/ioutil"
import "strings"

func main() {
	raw, error := ioutil.ReadFile("/usr/share/dict/cracklib-words")
	if error == nil {
		data := string(raw)
		words := strings.Split(data, "\n", 0)
		dict := make(map[string]bool)
		for _, w := range words {
			dict[w] = true
		}
		fmt.Printf("%d words\n", len(dict))
	}
}
