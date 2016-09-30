# Installing Gosl on Mac OS X

The installation on OSX includes:

1. Xcode
2. MacPorts
3. CMake and Lapack
4. SuiteSparse
5. Go (golang)
6. Gosl 
 
## 1 Xcode

In a terminal, type the following commands and click install to install the *command line developer tools*:

```
xcode-select --install
```

## 2 MacPorts

Download and install MacPorts from https://www.macports.org/install.php

Update the list of available packages with the following command:

```
sudo port selfupdate
```

## 3 CMake and Lapack

To install CMake and Lapack, type:

```
sudo port install lapack cmake
```

(yes: it's a bit slow, especially the llvm package)

## 4 SuiteSparse

Download SuiteSparse from http://faculty.cse.tamu.edu/davis/suitesparse.html and extract the files into `/Users/name/pkg/SuiteSparse`.

Type the following commands:

```
cd /Users/name/pkg/SuiteSparse
make
make install
sudo mkdir -p /usr/local/lib
sudo cp lib/*.dylib /usr/local/lib
sudo mkdir -p /usr/local/include/suitesparse
sudo cp include/*.h /usr/local/include/
```

## 5 Go (golang)

Install Go (golang) for Mac OS X from https://golang.org

Edit `/Users/name/.profile` and add:

```
export GOPATH=/Users/name/mygo
export PYTHONPATH=$PYTHONPATH:$GOPATH/src/github.com/cpmech/gosl/scripts
```

## 6 Gosl
 
Type the following commands:

```
cd mygo
mkdir -p src/github.com/cpmech.com
cd src/github.com/cpmech/
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

Yay!
