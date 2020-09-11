#!/bin/bash

set -e

PKGS_ESSENTIAL="chk io"

PKGS_NEED_FLAGS="mpi"

PKGS_ALL=" \
fun/dbf fun/fftw fun \
gm/rw gm/tri gm/msh gm \
graph \
la/oblas la \
ml/imgd ml \
num/qpck num \
ode \
opt \
pde \
plt \
rnd/sfmt rnd/dsfmt rnd \
utl \
"

install_and_test(){
    HERE=`pwd`
    PKG=$1
    DOTEST=$2
    HASGENBASH=$3
    echo
    echo
    echo "=== compiling $PKG ============================================================="
    cd $PKG
    if [[ ! -z $HASGENBASH ]]; then
        bash xgenflagsfile.bash
    fi
    touch *.go
    go install
    if [ "$DOTEST" -eq 1 ]; then
        go test
    fi
    cd $HERE
}

for p in $PKGS_ESSENTIAL; do
    install_and_test $p 1
done

for p in $PKGS_NEED_FLAGS; do
    install_and_test $p 1 1
done

for p in $PKGS_ALL; do
    install_and_test $p 1
done

echo
echo "=== SUCCESS! ============================================================"
