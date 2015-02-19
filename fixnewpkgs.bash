#!/bin/bash

set -e

fix_pkgs() {
    HERE=`pwd`
    PKG=$1
    DOTEST=$2
    echo
    echo
    echo "[1;32m>>> fixing $PKG <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    cd $PKG

    for f in *.go; do
        sed -i -e 's/utl.Panic/chk.Panic/g' \
               -e 's/utl.Err\>/chk.Err/g' \
               -e 's/utl.CheckScalar/chk.Scalar/g' \
               -e 's/utl.CheckString/chk.String/g' \
               -e 's/utl.CheckVector/chk.Vector/g' \
               -e 's/utl.CheckMatrix/chk.Matrix/g' \
               -e 's/utl.CompareStrs/chk.Strings/g' \
               -e 's/utl.CompareInts/chk.Ints/g' \
               -e 's/utl.CompareDbls/chk.Vector/g' \
               -e 's/utl.TTitle/chk.PrintTitle/g' \
               -e 's/utl.AnaNum/chk.PrintAnaNum/g' \
               -e 's/utl.CheckAnaNum/chk.AnaNum/g' \
               -e 's/utl.Pf/io.Pf/g' \
               -e 's/utl.Sf/io.Sf/g' \
               -e 's/utl.Ff/io.Ff/g' \
               -e 's/utl.Write/io.Write/g' \
               -e 's/utl.Read/io.Read/g' \
               -e 's/utl.Atoi/io.Atoi/g' \
               -e 's/utl.Atof/io.Atof/g' \
               -e 's/utl.Atob/io.Atob/g' \
               -e 's/utl.Itob/io.Itob/g' \
               -e 's/utl.Btoi/io.Btoi/g' \
               -e 's/utl.Btoa/io.Btoa/g' $f
        goimports -w $f
    done

    cd $HERE
}

fix_pkgs_simple() {
    HERE=`pwd`
    PKG=$1
    DOTEST=$2
    echo
    echo
    echo "[1;32m>>> fixing $PKG <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<[0m"
    cd $PKG

    for f in t_*.go; do
        sed -i -e 's/utl.Sramp/fun.Sramp/g' $f
        #goimports -w $f
    done

    cd $HERE
}

for p in mpi la plt fdm num fun ode gm tsr vtk; do
#    fix_pkgs $p 1
    fix_pkgs_simple $p 1
done
