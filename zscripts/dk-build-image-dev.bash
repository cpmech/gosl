#!/bin/bash

NAME=gosl_dev

echo
echo
echo "... docker .................................................."
echo "............................................................."
echo

docker build --no-cache -t gosl/$NAME . --build-arg DEV_IMG=true
docker images -q -f "dangling=true" | xargs docker rmi
