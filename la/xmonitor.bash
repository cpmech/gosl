#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    go test -test.run="SpSolver06M"
    #go test
    #mpirun -np 3 go run t_sp_mpi_main.go
    #mpirun -np 2 go run t_mumpssol01a_main.go
done
