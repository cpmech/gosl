#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    go install
    go test -test.run="stat01"
done
