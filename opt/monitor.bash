#!/bin/bash

#FILE="linipm.go"
FILE="t_linipm_test.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    go test -test.run="linipm02"
done
