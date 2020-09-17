#!/bin/bash

NAME="gosl"
VERSION="latest"

docker logout
docker login
docker push gosl/$NAME:$VERSION
