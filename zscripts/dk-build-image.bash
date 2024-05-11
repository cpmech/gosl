#!/bin/bash

NAME="dk-gosl-gofem"
VERSION="latest"

echo
echo
echo "... docker .................................................."
echo "............................................................."
echo

docker buildx build -t $NAME:$VERSION .
#docker images -q -f "dangling=true" | xargs docker rmi
