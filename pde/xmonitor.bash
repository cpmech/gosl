#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="EssentialBcs02"
    #go test -test.run="Fdm02"
    go test -test.run="Spc05"
done
