#!/bin/bash

set -e

platform='unknown'
unamestr=`uname`
if [[ "$unamestr" == 'Linux' ]]; then
   platform='linux'
elif [[ "$unamestr" == 'MINGW32_NT-6.2' ]]; then
   platform='windows'
elif [[ "$unamestr" == 'MINGW64_NT-10.0' ]]; then
   platform='windows'
elif [[ "$unamestr" == 'Darwin' ]]; then
   platform='darwin'
fi

echo "platform = $platform"

install_and_test(){
    HERE=`pwd`
    PKG=$1
    DOTEST=$2
    HASALLBASH=$3
    echo
    echo
    echo ">>> compiling $PKG <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    cd $PKG
    if [[ ! -z $HASALLBASH ]]; then
        ./all.bash
    fi
    touch *.go
    go install
    if [ "$DOTEST" -eq 1 ]; then
        go test
    fi
    cd $HERE
}

for p in chk io utl plt; do
    install_and_test $p 1
done

if [[ $platform == 'linux' ]]; then
    install_and_test mpi 0
fi

for p in la fdm num fun gm/rw gm/tri gm/msh gm graph ode opt tsr; do
    install_and_test $p 1
done

if [[ $platform != 'windows' ]]; then
    install_and_test rnd/sfmt 1
    install_and_test rnd/dsfmt 1
fi

install_and_test rnd 1

if [[ $platform == 'linux' ]]; then
    install_and_test vtk 0 1
fi

echo
echo ">>> SUCCESS! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
