#!/bin/bash

go build -o /tmp/gosl/t_jacobian01b_main t_jacobian01b_main.go && mpirun -np 2 /tmp/gosl/t_jacobian01b_main
go build -o /tmp/gosl/t_jacobian02b_main t_jacobian02b_main.go && mpirun -np 4 /tmp/gosl/t_jacobian02b_main
