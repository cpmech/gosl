#!/bin/bash

GP1="${GOPATH%:*}"
GP2="${GOPATH#*:}"
GP=$GP2
if [[ -z "${GP// }" ]]; then
    GP=$GP1
fi
GOSL="$GP/pkg/linux_amd64/github.com/cpmech/gosl/gm.a"

FILES="*.go"

if [ -f $GOSL ]; then
   FILES="$FILES $GOSL"
fi

echo
echo "monitoring:"
echo $FILES
echo
echo "with:"
echo "GP = $GP"
echo
echo

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    go run opt_ipm01.go
done
