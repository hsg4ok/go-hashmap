# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ../../../Make.$(GOARCH)

TARG=container/hashmap
GOFILES=hashmap.go hashvec.go
CLEANFILES+=example_map example_hashmap primer test_random

include ../../../Make.pkg

benchmark:
	gotest -benchmarks=".*"

example_map: install example_map.go
	$(GC) example_map.go
	$(LD) -o $@ example_map.$O

example_hashmap: install example_hashmap.go
	$(GC) example_hashmap.go
	$(LD) -o $@ example_hashmap.$O

test_random: install test_random.go
	$(GC) test_random.go
	$(LD) -o $@ test_random.$O

primer: primer.go
	$(GC) $^
	$(LD) -o $@ primer.$O
