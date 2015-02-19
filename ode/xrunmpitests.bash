#!/bin/bash

go build -o /tmp/gosl/t_ODE01b_main t_ODE01b_main.go && mpirun -np 2 /tmp/gosl/t_ODE01b_main
go build -o /tmp/gosl/t_ODE04b_main t_ODE04b_main.go && mpirun -np 3 /tmp/gosl/t_ODE04b_main
