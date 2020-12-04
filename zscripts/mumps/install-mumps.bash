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
cp ${PDIR}/Makefile.inc .
make d
make z
chmod -x lib/*

# copy include and lib files to the right places
mkdir -p /usr/include/mumps
cp -av include/*.h /usr/include/mumps/
cp -av lib/libpord.so /usr/lib/
cp -av lib/libdmumps.so /usr/lib/
cp -av lib/libzmumps.so /usr/lib/
cp -av lib/libmumps_common.so /usr/lib/

# clean up
rm -rf ${MUMPS_TMP}
