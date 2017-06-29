#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    mpirun -np 8 go run t_mpi00_main.go
    #mpirun -np 3 go run t_mpi01_main.go
    #mpirun -np 3 go run t_mpi02_main.go
    #mpirun -np 3 go run t_mpi03_main.go
    #mpirun -np 3 go run t_mpi04_main.go
done
