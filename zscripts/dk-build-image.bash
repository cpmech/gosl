#!/bin/bash

NAME="gosl"
VERSION="1.1.3"

echo
echo
echo "... docker .................................................."
echo "............................................................."
echo

docker build --no-cache -t gosl/$NAME:$VERSION . --build-arg GOSL_VERSION=$VERSION
docker images -q -f "dangling=true" | xargs docker rmi
