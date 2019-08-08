#!/bin/bash

set -e

examples="\
h5_ang_mnist01.go
"

for ex in $examples; do
    echo
    echo
    echo ">>> running $ex <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    go run "$ex"
done

echo
echo "=== SUCCESS! ============================================================"
