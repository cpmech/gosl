#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    go test -test.run="conv07"
done
