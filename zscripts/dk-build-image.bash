#!/bin/bash

NAME=gosl

echo
echo
echo "... docker .................................................."
echo "............................................................."
echo

docker build -t gosl/$NAME .
docker images -q -f "dangling=true" | xargs docker rmi
