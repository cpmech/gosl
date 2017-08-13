#!/bin/bash

FILES="*.go"

echo
echo "monitoring:"
echo $FILES
echo
echo

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -test.run="Waterfall01"
    #go test -test.run="draw03"
    go test -test.run="plot01"
done
