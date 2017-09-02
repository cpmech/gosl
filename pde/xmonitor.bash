#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="EssentialBcs01"
    #go test -test.run="FdmLaplace03"
    go test -test.run="Spc01"
done
