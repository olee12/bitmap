#!/bin/bash

FILENAME=$(basename $(pwd))
go test -run=. -bench=. -benchmem -memprofile=mem.out -timeout 3000s
go tool pprof -pdf --alloc_space $FILENAME.test mem.out > alloc_space.pdf && open alloc_space.pdf
