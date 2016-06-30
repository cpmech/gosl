#!/bin/bash

set -e

examples="\
la_HLsparseComplex01 \
la_HLsparseReal01 \
la_sparseComplex01 \
la_sparseReal01 \
num_deriv01 \
vtk_cone01 \
vtk_isosurf01"

for ex in $examples; do
    echo
    echo
    echo "[1;32m>>> running $ex <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    go run "$ex".go
done

