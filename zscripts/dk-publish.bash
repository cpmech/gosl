#!/bin/bash

NAME="gosl/dk-gosl-gofem"

docker logout
docker login
docker push $NAME
