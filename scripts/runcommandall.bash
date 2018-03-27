#!/bin/bash

set -e

echo "usage:"
echo "    $0 JOB"
echo "where JOB is:"
echo "    0 -- count lines [default]"
echo "    1 -- execute gofmt"
echo "    2 -- execute golint"
echo "    3 -- execute go vet"
echo "    4 -- execute gocyclo"
echo "    5 -- execute goimports"
echo "    6 -- generate depedency graphs"
echo "    7 -- fix links in README files"

JOB=0
if [[ $# != 0 ]]; then
    JOB=$1
    if [[ $JOB -lt 0 || $JOB -gt 7 ]]; then
        echo
        echo "Job number $1 is invalid"
        echo
        exit 1
    fi
fi

echo "current JOB = $JOB"

if [[ $JOB == 0 ]]; then
    totnfiles=0
    totnlines=0
    for f in `find . -iname "*.go"`; do
        totnfiles=$(($totnfiles+1))
        totnlines=$(($totnlines+`wc -l $f | awk '{print $1}'`))
    done
    echo
    echo "Total number of go files = $totnfiles"
    echo "Total number of go lines = $totnlines"
    exit 0
fi

ALL="
chk \
io \
utl \
utl/al \
plt \
ml \
ml/imgd \
mpi \
la  \
la/mkl \
la/oblas \
num/qpck \
num \
fun \
fun/dbf \
fun/fftw \
gm \
gm/msh \
gm/tri \
gm/rw \
graph \
opt \
rnd \
rnd/sfmt \
rnd/dsfmt \
vtk \
pde \
examples \
tools \
"

rungofmt() {
    pkg=$1
    for f in *.go; do
        gofmt -s -w $f
    done
}

rungolint() {
    golint .
}

rungovet() {
    go vet .
}

rungocyclo() {
    level=`gocyclo -top 1 . | awk '{print $1}'`
    if [[ $level -gt 15 ]]; then
        echo "PROBLEM: level = $level"
    fi
}

rungoimports() {
    pkg=$1
    for f in *.go; do
        echo $f
        goimports -w $f
    done
}

depgraph(){
    if [[ $1 == "examples" || $1 == "tools" ]]; then
        return
    fi
    pkg=$1
    fna="/tmp/gosl/depgraph-${pkg/\//_}-A.png"
    fnb="/tmp/gosl/depgraph-${pkg/\//_}-B.svg"
    godepgraph -s github.com/cpmech/gosl/$pkg | dot -Tpng -o $fna
    graphpkg -stdout -match 'gosl' github.com/cpmech/gosl/$pkg > $fnb
    echo "file <$fna> generated"
    echo "file <$fnb> generated"
}

fixreadme() {
    pkg=$1
    #old="http://rawgit.com/cpmech/gosl/master/doc/xx${pkg/\//-}.html"
    #new="https://godoc.org/github.com/cpmech/gosl/${pkg}"
    #sed -i 's,'"$old"','"$new"',' README.md

    # prepend link before key
    #key="More information is available in"
    #lnk="[![GoDoc](https://godoc.org/github.com/cpmech/gosl/${pkg}?status.svg)](https://godoc.org/github.com/cpmech/gosl/${pkg})"
    #lnk=$(echo "$lnk" | sed 's/\//\\\//g')
    #sed -i "/More information is available in/i $lnk \n" README.md
}

if [[ $JOB == 6 ]]; then # graphs
    mkdir -p /tmp/gosl
fi

idx=1
for pkg in $ALL; do
    HERE=`pwd`
    cd $pkg
    echo
    echo ">>> $idx $pkg <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
    if [[ $JOB == 1 ]]; then
        rungofmt $pkg
    fi
    if [[ $JOB == 2 ]]; then
        rungolint $pkg
    fi
    if [[ $JOB == 3 ]]; then
        rungovet $pkg
    fi
    if [[ $JOB == 4 ]]; then
        rungocyclo $pkg
    fi
    if [[ $JOB == 5 ]]; then
        rungoimports $pkg
    fi
    if [[ $JOB == 6 ]]; then
        depgraph $pkg
    fi
    if [[ $JOB == 7 ]]; then
        fixreadme $pkg
    fi
    cd $HERE
    (( idx++ ))
done
