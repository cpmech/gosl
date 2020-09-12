#!/bin/bash

# Syntax: ./gosl-clone-and-build.bash [DEV_IMG] [GOSL_BRANCH]

DEV_IMG=${1:-"false"}
GOSL_BRANCH=${2:-"trim-gosl"}

if [ "${DEV_IMG}" = "false" ]; then
    cd /usr/local/go/src
    git clone -b $GOSL_BRANCH --single-branch --depth 1 https://github.com/cpmech/gosl.git
    cd gosl
    bash ./all.bash
fi
