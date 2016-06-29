#!/bin/bash

FILES="io.go step.go t_step_test.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    go test -test.run="step02"
done
