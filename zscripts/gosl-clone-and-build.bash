#!/bin/bash

# Syntax: ./gosl-clone-and-build.bash [DEV_IMG] [GOSL_VERSION]

DEV_IMG=${1:-"false"}
GOSL_VERSION=${2:-"2.0.0"}

if [ "${DEV_IMG}" = "true" ]; then
  exit 0
fi

BRANCH="v$GOSL_VERSION"
if [ "${GOSL_VERSION}" = "latest" ]; then
  BRANCH="master"
fi

cd /usr/local/go/src
mkdir -p github.com/cpmech
cd github.com/cpmech
git clone -b $BRANCH --single-branch --depth 1 https://github.com/cpmech/gosl.git
cd gosl
bash ./all.bash
