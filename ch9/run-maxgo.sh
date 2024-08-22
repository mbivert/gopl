#!/bin/sh

echo Do not run blindly.
exit 1

echo "Time to go through the pipeline:"
for n in 1 10 100 1000 10000 100000 1000000 10000000; do
	go run maxgo.go -m $n | sed 's/^/	/g'
done

# This is reasonable if you have about 10Go RAM available at least;
# 50000000 (IIRC) was eating 10Go of swap and as much RAM (killed).
echo "Creating a huge pipeline:"
for n in 1 10 100 1000 10000 100000 1000000; do
	go run maxgo.go -n -m $n | sed 's/^/	/g'
done
