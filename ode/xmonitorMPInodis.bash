#!/bin/bash

MAIN="t_eq11_nodis_main"

FILES="radau5.go $MAIN.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    go build -o /tmp/gosl/$MAIN $MAIN.go && mpirun -np 1 /tmp/gosl/$MAIN
done
