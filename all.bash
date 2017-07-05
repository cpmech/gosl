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
    HASGENBASH=$3
    echo
    echo
    echo ">>> compiling $PKG <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    cd $PKG
    if [[ ! -z $HASGENBASH ]]; then
        ./xgenflagsfile.bash
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
    install_and_test mpi 1
else
    install_and_test mpi 0
fi

for p in la/oblas la fun/dbf fun/fftw fun num/qpck num gm/rw gm/msh gm graph opt; do
    install_and_test $p 1
done

if [[ $platform != 'windows' ]]; then
    install_and_test gm/tri 1
    install_and_test rnd/sfmt 1
    install_and_test rnd/dsfmt 1
fi

install_and_test rnd 1

if [[ $platform == 'linux' ]]; then
    install_and_test vtk 0 1
fi

echo
echo ">>> SUCCESS! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
