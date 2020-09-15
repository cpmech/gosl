#!/bin/bash

# Syntax: ./gosl-clone-and-build.bash [DEV_IMG] [GOSL_VERSION]

DEV_IMG=${1:-"false"}
GOSL_VERSION=${2:-"2.0.0"}

if [ "${DEV_IMG}" = "true" ]; then
  exit 0
fi

cd /usr/local/go/src
git clone -b "v$GOSL_VERSION" --single-branch --depth 1 https://github.com/cpmech/gosl.git
cd gosl
bash ./all.bash
