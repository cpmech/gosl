#!/bin/bash

NAME="gosl"
VERSION="2.0.0"

docker logout
docker login
docker push gosl/$NAME:$VERSION
