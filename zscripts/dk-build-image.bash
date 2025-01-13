#!/bin/bash

NAME="gosl/dk-gosl-gofem"

echo
echo
echo "... docker .................................................."
echo "............................................................."
echo

docker build -t $NAME .
#docker images -q -f "dangling=true" | xargs docker rmi
