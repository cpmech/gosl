#!/bin/bash

NAME="gosl"
VERSION="1.1.3"

docker logout
docker login
docker push gosl/$NAME:$VERSION
