#!/bin/bash

set -e

platform='unknown'
unamestr=`uname`
if [[ "$unamestr" == 'Linux' ]]; then
   platform='linux'
elif [[ "$unamestr" == 'MINGW32_NT-6.2' ]]; then
   platform='windows'
fi

echo "platform = $platform"

install_and_test(){
    HERE=`pwd`
    PKG=$1
    DOTEST=$2
    echo
    echo
    echo "[1;32m>>> compiling $PKG <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    cd $PKG
    go install
    if [ "$DOTEST" -eq 1 ]; then
        go test
    fi
    cd $HERE
}

for p in chk io utl; do
    install_and_test $p 1
done

if [[ $platform == 'linux' ]]; then
    install_and_test mpi 0
fi

for p in la plt fdm num fun ode gm tsr; do
    install_and_test $p 1
done

if [[ $platform == 'linux' ]]; then
    install_and_test vtk 0
fi
