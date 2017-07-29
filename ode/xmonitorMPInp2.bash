#!/bin/bash

MAIN="t_vdp_np2_main"

FILES="radau5_mpi.go $MAIN.go"

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    go build -o /tmp/gosl/$MAIN $MAIN.go && mpirun -np 2 /tmp/gosl/$MAIN
done
