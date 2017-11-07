#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    go install
    #go test -test.run="octree02"
    #go test -test.run="Grid03"
    #go test -test.run="Nurbs08"
    #go test -test.run="Transfinite01"
    #go test -test.run="Transfinite02"
    #go test -test.run="Transfinite03"
    #go test -test.run="Transfinite04"
    #go test -test.run="Transfinite05"
    go test -test.run="CurvGrid02"
done
