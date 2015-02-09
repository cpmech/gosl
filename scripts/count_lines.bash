#!/bin/bash

set -e

totnfiles=0
totnlines=0
for f in `find . -iname "*.go"`; do
	totnfiles=$(($totnfiles+1))
	totnlines=$(($totnlines+`wc -l $f | awk '{print $1}'`))
done

echo "Total number of files = $totnfiles"
echo "Total number of lines = $totnlines"
