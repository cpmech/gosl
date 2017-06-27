#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    go test -test.run="Linsol03"
    #go test
done
