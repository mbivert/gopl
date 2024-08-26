#!/bin/sh

set -e

for s in 1024 2048; do
	for x in 10 50 100 200 400 1000 2000 3000 5000; do
		printf "%-6d %-6d %-6d " $s $s $x
		go run parallel-mandelbrot.go -width $s -height $s -ngo $x > t3.png
	done
done

for s in 4096; do
	for x in 1000 2000 3000 5000 7000; do
		printf "%-6d %-6d %-6d " $s $s $x
		go run parallel-mandelbrot.go -width $s -height $s -ngo $x > t3.png
	done
done

for s in 8192; do
	for x in 3000 5000 7000 8000 9000; do
		printf "%-6d %-6d %-6d " $s $s $x
		go run parallel-mandelbrot.go -width $s -height $s -ngo $x > t3.png
	done
done

for s in 16384; do
	for x in 5000 10000 15000 20000 30000; do
		printf "%-6d %-6d %-6d " $s $s $x
		go run parallel-mandelbrot.go -width $s -height $s -ngo $x > t3.png
	done
done
