#!/bin/bash

set -e

PKGS_ESSENTIAL="chk io utl la/oblas la"

PKGS_ALL=" \
fun/fftw fun \
gm/tri gm/msh gm \
hb \
num/qpck num \
ode \
opt \
pde \
rnd/sfmt rnd/dsfmt rnd \
"

install_and_test(){
    HERE=`pwd`
    PKG=$1
    DOTEST=$2
    echo
    echo
    echo "=== compiling $PKG ============================================================="
    cd $PKG
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

for p in $PKGS_ALL; do
    install_and_test $p 1
done

echo
echo "=== SUCCESS! ============================================================"
