#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="LagCardinal01"
    #go test -test.run="LagInterp08"
    #go test -test.run="ChebyInterp07"
    go test -test.run="LagCheby03"
done
