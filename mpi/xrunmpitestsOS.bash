#!/bin/bash

go build -o /tmp/gosl/t_mpi00_main t_mpi00_main.go && mpirun --oversubscribe -np 8 /tmp/gosl/t_mpi00_main

tests="t_mpi01_main t_mpi02_main t_mpi03_main t_mpi04_main"
for t in $tests; do
    go build -o /tmp/gosl/$t "$t".go && mpirun --oversubscribe -np 3 /tmp/gosl/$t
done
