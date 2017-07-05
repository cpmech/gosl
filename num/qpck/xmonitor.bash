#!/bin/bash

FILE="*.go *.c *.f"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test
    go test -test.run="SpecProb01"
done
