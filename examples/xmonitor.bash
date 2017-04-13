#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    #go run num_brent01.go
    go run num_newton01.go
done
