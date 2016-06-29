#!/bin/bash

FILES="parsing.go t_parsing_test.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    go test -test.run="parsing07"
done
