#!/bin/bash

NP1="
t_eq11_np1_main \
"

NP2="
t_vdp_np2_main \
"

NP3="
t_amp_np3_main
"

for main in $NP1; do
    go build -o /tmp/gosl/$main $main.go && mpirun -np 1 /tmp/gosl/$main
done

for main in $NP2; do
    go build -o /tmp/gosl/$main $main.go && mpirun -np 2 /tmp/gosl/$main
done

for main in $NP3; do
    go build -o /tmp/gosl/$main $main.go && mpirun -np 3 /tmp/gosl/$main
done
