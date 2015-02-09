#!/bin/bash

go build t_ODE01b_main.go
mpirun -np 2 ./t_ODE01b_main
rm t_ODE01b_main

go build t_ODE04b_main.go
mpirun -np 3 ./t_ODE04b_main
rm t_ODE04b_main
