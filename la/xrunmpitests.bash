#!/bin/bash

tests="\
z_mumpssol01a_main \
z_mumpssol01b_main \
z_mumpssol02_main \
z_mumpssol03_main \
z_mumpssol04_main \
z_mumpssol05_main \
z_sp_matrix01_main \
z_sp_matrix02_main \
z_sp_matrix03_main \
z_sp_mpi_main
"

for t in $tests; do
    go build -o /tmp/gosl/$t "$t".go && mpirun -np 2 /tmp/gosl/$t
done
