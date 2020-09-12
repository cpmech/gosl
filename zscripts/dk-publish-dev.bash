#!/bin/bash

NAME=gosl_dev

docker logout
docker login
docker push gosl/$NAME
