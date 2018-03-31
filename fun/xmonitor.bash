#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    echo
    echo
    #go test -run MultiInterp01
    go test -run Pow2
done
