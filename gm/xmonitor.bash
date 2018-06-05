#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -run Nurbs14
    #go test -run TestNurbsMethods01
    #go test -run TestGrid06
    #go test -run Transfinite07
    #go test -run bins07
    go test -run npatch02
done
