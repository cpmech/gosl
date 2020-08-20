#!/bin/bash

NAME=gosl

docker logout
docker login
docker push gosl/$NAME
