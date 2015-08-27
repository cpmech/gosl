#!/bin/bash

FILES="auxiliary.go graph.go t_graph_test.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    go test -test.run="graph04"
done
