#!/bin/bash

GP1="${GOPATH%:*}"
GP2="${GOPATH#*:}"
GP=$GP2
if [[ -z "${GP// }" ]]; then
    GP=$GP1
fi
GOSL="$GP/pkg/linux_amd64/github.com/cpmech/gosl/gm.a"

FILES="*.go"

if [ -f $GOSL ]; then
   FILES="$FILES $GOSL"
fi

echo
echo "monitoring:"
echo $FILES
echo
echo "with:"
echo "GP = $GP"
echo
echo

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    #go run graph_munkres01.go
    #go run graph_shortestpaths01.go
    #go run graph_siouxfalls01.go
    #go run rnd_ints01.go
    go run rnd_lognormalDistribution.go
done
