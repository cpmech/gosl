#!/bin/bash

tests="t_mpi01_main t_mpi02_main t_mpi03_main t_mpi04_main"
for t in $tests; do
    go build "$t".go
    mpirun -np 3 ./$t
    rm "$t"
done
