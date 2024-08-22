#!/bin/sh

r=3
for s in 10 100 1000 10000 100000 1000000 10000000 ; do
	go run pingpong.go -s $s -r $r | awk '{ v += $1; n++ } END {
		printf("%-15d: %-16d byte/sec\n", '$s', v/n)
	}'
done
