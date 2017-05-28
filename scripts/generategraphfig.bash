#!/bin/bash

# Install:
#   go get github.com/kisielk/godepgraph
#   go get github.com/davecheney/graphpkg

generate_depgraph(){
    IDX=$1
    PKG=$2
    echo
    echo
    echo ">>> gen dep graph $IDX $PKG <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    FNA="/tmp/gosl/depgraph-${PKG/\//_}-A.png"
    FNB="/tmp/gosl/depgraph-${PKG/\//_}-B.svg"
    godepgraph -s github.com/cpmech/gosl/$PKG | dot -Tpng -o $FNA
    graphpkg -stdout -match 'gosl' github.com/cpmech/gosl/$PKG > $FNB
    echo "file <$FNA> generated"
    echo "file <$FNB> generated"
}

ALL="
chk \
io \
utl \
plt \
mpi \
la  \
la/mkl \
la/oblas \
fdm \
num \
fun \
gm \
gm/msh \
gm/tri \
gm/rw \
graph \
ode \
opt \
rnd \
rnd/sfmt \
rnd/dsfmt \
tsr \
vtk \
img \
img/ocv \
"

#ALL="la/mkl"

idx=1
for p in $ALL; do
    generate_depgraph $idx $p
    (( idx++ ))
done
