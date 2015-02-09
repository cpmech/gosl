#!/bin/bash

go build t_jacobian01b_main.go
mpirun -np 2 ./t_jacobian01b_main
rm t_jacobian01b_main

go build t_jacobian02b_main.go
mpirun -np 4 ./t_jacobian02b_main
rm t_jacobian02b_main
