#!/bin/bash

FILE="t_pareto_test.go"
#FILE="pareto.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    go test -test.run="pareto04"
    python /tmp/gosl/test_pareto04.py
done
