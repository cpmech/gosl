#!/bin/bash

set -e

examples="\
la_umfpack01 \
num_deriv01 \
vtk_isosurf01 "

for ex in $examples; do
    echo
    echo
    echo "[1;32m>>> running $ex <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    go run "$ex".go
done
