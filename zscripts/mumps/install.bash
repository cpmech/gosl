#!/bin/bash

set -e

# constants
MUMPS_VERSION="5.3.5"
MUMPS_GZ=mumps_${MUMPS_VERSION}.orig.tar.gz
MUMPS_TMP=/tmp/MUMPS_${MUMPS_VERSION}
PDIR=/tmp/zscripts/mumps/patch

# download and exctract the source code
curl http://deb.debian.org/debian/pool/main/m/mumps/${MUMPS_GZ} -o /tmp/${MUMPS_GZ}
cd /tmp && tar xzf ${MUMPS_GZ} && rm ${MUMPS_GZ}

# patch and compile
cd ${MUMPS_TMP}
patch -u PORD/lib/Makefile ${PDIR}/PORD/lib/Makefile.diff
patch -u src/Makefile ${PDIR}/src/Makefile.diff
patch -u Makefile ${PDIR}/Makefile.diff
cp ${PDIR}/Makefile.inc .
make all

# copy include and lib files to the right places
cp include/*.h /usr/include/
cp -av lib/*.so /usr/lib/

# clean up
rm -rf ${MUMPS_TMP}
