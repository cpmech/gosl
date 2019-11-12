#!/bin/bash

set -e

arch=`uname -m`
if [[ "$arch" == 'x86_64' ]]; then
    sse="-msse2 -DHAVE_SSE2"
fi

DSFMT_GO="dsfmt.go"
DSFMT_GO_IN="${DSFMT_GO}.in"
echo "generating $DSFMT_GO ($arch)"

CGO="\/*\n\
#cgo CFLAGS: -O3 -fomit-frame-pointer -DNDEBUG -fno-strict-aliasing -std=c99 $sse -DDSFMT_MEXP=19937\n\
#include \"connectdsfmt.h\"\n\
*\/"

cat $DSFMT_GO_IN | sed "s/@@CGO@@/$CGO/" > $DSFMT_GO
