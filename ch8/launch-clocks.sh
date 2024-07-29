#!/bin/sh

set -e

clock="go run ./clock.go"
clockwall="go run ./clockwall.go"

TZ=US/Eastern    $clock -port :8010 &
TZ=Europe/London $clock -port :8020 &
TZ=Asia/Tokyo    $clock -port :8030 &


sleep 1

$clockwall NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030

# brutal
pkill clock
