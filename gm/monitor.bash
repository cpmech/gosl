#!/bin/bash

FILES="io.go nurbs.go t_nurbs_test.go t_rw_test.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    go test -test.run="rwstep01"
done
