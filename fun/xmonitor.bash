#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="LagCardinal01"
    #go test -test.run="LagInterp05"
    #go test -test.run="LagInterp06"
    #go test -test.run="LagCheby01"
    #go test -test.run="LagCheby02"
    #go test -test.run="ChebyInterp06"
    go test -test.run="ChebyInterp07"
done
