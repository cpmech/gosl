#!/bin/bash

FILES="*.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -run TestGrid01
    go test -run TestRectGrid01
    #go test -run TestCurvGrid02
done
