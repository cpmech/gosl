#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -run Nurbs12
    #go test -run TestNurbsMethods01
    go test -run TestGrid10
done
