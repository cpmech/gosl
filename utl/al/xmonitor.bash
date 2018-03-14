#!/bin/bash

FILES="*.go"

echo
echo "monitoring:"
echo $FILES
echo
echo

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    #go test -run Adata01
    #go test -run Queue01
    #go test -run LinkedList01
    go test -run LinkedList02
done
