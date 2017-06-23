#!/bin/bash

set -e

keys="\
oblas_dgemm01
"

for k in $keys; do
    echo
    echo
    echo "[1;32m>>> running $f <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    go build -o /tmp/$k $k.go && /tmp/$k
done
