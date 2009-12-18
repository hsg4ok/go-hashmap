#!/usr/bin/env python
with open("/usr/share/dict/cracklib-words") as f:
	data = f.read()
	words = data.split("\n")
	dict = {}
	for w in words:
		dict[w] = True
	print "%d words" % len(dict)
