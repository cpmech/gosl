#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -run TestGrid06
    go test -run Nurbs09
done
