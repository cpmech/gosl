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
    echo "=== compiling $PKG ============================================================="
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

for p in chk io utl/al utl plt; do
    install_and_test $p 1
done

if [[ $platform != 'windows' ]]; then
    install_and_test mpi 1
    install_and_test io/h5 1
else
    install_and_test mpi 0
fi

for p in la/oblas la fun/dbf fun/fftw fun num/qpck num gm/rw gm/tri gm/msh gm graph opt; do
    install_and_test $p 1
done

if [[ $platform != 'windows' ]]; then
    install_and_test rnd/sfmt 1
    install_and_test rnd/dsfmt 1
fi

for p in rnd ml/imgd ml ode pde tsr; do
    install_and_test $p 1
done

echo
echo "=== SUCCESS! ============================================================"
