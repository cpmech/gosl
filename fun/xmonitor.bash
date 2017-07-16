#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="Chebyshev01"
    go test -test.run="ChebyInterp01"
done
