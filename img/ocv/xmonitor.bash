#!/bin/bash

FILE="*.go *.h *.cpp"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    go install
    go test -test.run="ocv01"
done
