#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -test.run="cubeandtet"
    #go test -test.run="mesh01"
    #go test -test.run="singleq4"
    #go test -test.run="Quadpts01"
    #go test -test.run="Integ04"
    go test -test.run="Gen01"
done
