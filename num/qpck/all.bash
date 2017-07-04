#!/bin/bash

HERE=$PWD
TMP=/tmp/gosl/num/build_qpck
mkdir -p $TMP && cd $TMP
cmake $HERE
make

GP="${GOPATH%:*}"
cp libqpck.a $GP/pkg/
echo "file <$GP/pkg/libqpck.a> created"

cd $HERE
