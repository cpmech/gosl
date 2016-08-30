#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    go test -test.run="hash01"
    #go test -test.run="nurbs04"
done
