#!/bin/bash

go build -o /tmp/gosl/t_ode02_main t_ode02_main.go && mpirun -np 2 /tmp/gosl/t_ode02_main
go build -o /tmp/gosl/t_ode04_main t_ode04_main.go && mpirun -np 3 /tmp/gosl/t_ode04_main
