#!/bin/bash

tests="t_mpi01_main t_mpi02_main t_mpi03_main t_mpi04_main"
for t in $tests; do
    go build -o /tmp/gosl/$t "$t".go && mpirun -np 3 /tmp/gosl/$t
done
