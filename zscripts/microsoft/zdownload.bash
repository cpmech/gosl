#!/bin/bash

SCRIPTS="\
  common-debian.sh \
  go-debian.sh
"

for s in $SCRIPTS; do
  curl https://raw.githubusercontent.com/microsoft/vscode-dev-containers/master/script-library/$s -o $s
done