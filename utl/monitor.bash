#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo

    #go test -test.run="pareto04"
    #python /tmp/gosl/test_pareto04.py

    go test -test.run="mylab10"
done
