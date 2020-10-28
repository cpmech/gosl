#!/bin/bash

# Syntax: ./gosl-clone-and-build.bash [DEV_IMG]

DEV_IMG=${1:-"false"}

BRANCH_OR_TAG="stable-1.1.3"

if [ "${DEV_IMG}" = "true" ]; then
  exit 0
fi

cd /usr/local/go/src
git clone -b $BRANCH_OR_TAG --single-branch --depth 1 https://github.com/cpmech/gosl.git
cd gosl
bash ./all.bash
