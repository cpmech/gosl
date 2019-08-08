#!/bin/bash

GOVER=1.12
GOFN=go$GOVER.linux-amd64

# deps
echo "... installing deps ..."
sudo apt-get -y install wget git gcc \
    libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev \
    gfortran python-scipy python-matplotlib dvipng \
    libfftw3-dev libfftw3-mpi-dev libmetis-dev \
    liblapacke-dev libopenblas-dev libhdf5-dev

# vtk
if [[ ! -z "$USE_VTK" ]]; then
    echo "... installing VTK ..."
    sudo apt-get -y install libvtk7-dev
fi

# go
mkdir -p ~/xpkg
cd ~/xpkg
rm -rf go
wget https://dl.google.com/go/$GOFN.tar.gz -O ~/xpkg/$GOFN.tar.gz
tar xf $GOFN.tar.gz
go get -u all

# output
echo
echo "go version"
go version
