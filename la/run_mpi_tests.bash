#!/bin/bash

tests="t_sumtoroot_main t_mumpssol01a_main t_mumpssol01b_main t_mumpssol02_main t_mumpssol03_main t_mumpssol04_main t_mumpssol05_main"
for t in $tests; do
    go build "$t".go
    mpirun -np 2 ./$t
    rm "$t"
done
