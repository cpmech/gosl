#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    go test -test.run="BryConds02"
    #go test -test.run="Fdm02"
    #go test -test.run="Spc05"
done
