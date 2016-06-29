#!/bin/bash

FILES="io.go nurbs.go search.go t_nurbs_test.go t_rw_test.go t_search_test.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    go test -test.run="bins02"
done
