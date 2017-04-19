#!/bin/bash

set -e

ALL="chk examples fdm fun gm gm/msh gm/rw gm/tri graph io la mpi num ode opt plt rnd tools tsr utl vtk"

runcommand() {
    pkg=$1
    echo
    echo
    echo ">>>>>>>>>>>>>>>> $pkg <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    for f in *.go; do
        echo $f
        goimports -w $f
    done
}

for pkg in $ALL; do
    HERE=`pwd`
    cd $pkg
    runcommand $pkg
    cd $HERE
done
