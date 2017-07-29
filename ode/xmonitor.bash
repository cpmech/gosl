#!/bin/bash

FILE="*.go"

while true; do
    inotifywait -q -e modify $FILE
    echo
    echo
    echo
    echo
    #go test -test.run="ode04"
    go build -o /tmp/gosl/t_ode02_main t_ode02_main.go && mpirun -np 2 /tmp/gosl/t_ode02_main
done
