#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="FdmLaplace02"
    go test -test.run="EssentialBcs01"
done
