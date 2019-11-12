#!/bin/bash

set -e

arch=`uname -m`
if [[ "$arch" == 'x86_64' ]]; then
    sse="-msse2 -DHAVE_SSE2"
fi

SFMT_GO="sfmt.go"
SFMT_GO_IN="${SFMT_GO}.in"
echo "generating $SFMT_GO ($arch)"

CGO="\/*\n\
#cgo CFLAGS: -O3 -fomit-frame-pointer -DNDEBUG -fno-strict-aliasing -std=c99 $sse -DSFMT_MEXP=19937\n\
#include \"connectsfmt.h\"\n\
#ifdef WIN32\n\
#define LONG long long\n\
#else\n\
#define LONG long\n\
#endif\n\
*\/"

cat $SFMT_GO_IN | sed "s/@@CGO@@/$CGO/" > $SFMT_GO
