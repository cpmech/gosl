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
    #go test -test.run="LagCheby03"
    go test -test.run="FourierInterp05"
    #go test -test.run="Dft01"
    #go test -test.run="Hat01"
    #go test -test.run="ImagPowN"
done
