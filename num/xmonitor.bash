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
    #go test -run Bracket01
    #go test -run Bracket02
    #go test -run Brent01
    #go test -run Brent02
    #go test -run Brent03
    #go test -run LineSolver01
    #go test -run LineSolver02
    go test -run ConjGrad01
done
