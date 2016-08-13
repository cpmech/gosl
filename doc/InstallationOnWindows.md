# Installing Gosl on Windows

Gosl (Go-Scientific Library) is a set of routines for computations involving numerical methods and other mathematical problems in science and engineering.

A number of existent tools are _wrapped_ by Gosl with the two most essential ones being *Lapack* and *Umfpack* (from SuiteSparse by Prof Tim Davis). The former is a _de facto_ set of routines for linear algebra and the latter is one of the most powerful set of tools for computations using sparse matrices and systems.

*Lapack* and *SuiteSparse*  can be readily installed on Debian/Linux systems. On Windows, three ``compilation'' tools for building Lapack and Umfpack are required beforehand:

1. Gcc and Gfortran for Windows.
2. CMake for Windows.
3. Git Bash for Windows

Therefore, these must be first installed, followed by the compilation of Lapack and SuiteSparse.

Lapack and SuiteSparse must installed in a directory named *C:\Gosl*.

This document is organised as follows:

1. INSTALL COMPILATION TOOLS
2. COMPILE LAPACK
    1. Download and uncompress
    2. Create project with CMake
    3. Build Lapack
3. COMPILE SUITESPARSE
4. INSTALL GO (GOLANG)
5. INSTALL GOSL



## 1 INSTALL COMPILATION TOOLS

The videos [published here](https://www.youtube.com/watch?v=9vFODRZTBcc&list=PLk1POg2YgVEI8OMZ-EOlfJGK0YWxY9-sL) may guide the steps listed below; but note that the directory *Gosl* has to be renamed as *Gosl*. Three videos are of relevance to this section:

1. YouTube [Download packages](https://youtu.be/9vFODRZTBcc)
2. YouTube [Install tools](https://youtu.be/dLyoGflSFTI)
3. YouTube [Compile Lapack](https://youtu.be/nsS3C1R3aDw)

Steps:

1. TDM GCC AND GFORTRAN: From http://tdm-gcc.tdragon.net/download download and install *tdm64-gcc-5.1.0-2.exe* (or newer). Make sure to select *fortran* and Keep other default options.
2. CMAKE: From https://cmake.org/download/ download and install *cmake-3.6.1-win64-x64.msi*. Keep default options.
3. GIT AND GIT BASH: From https://git-scm.com/download/ download and install *Git-2.8.1-64-bit.exe*. Keep default options.



## 2 COMPILE LAPACK

Lapack and SuiteSparse will be compiled and installed in a directory named *C:\Gosl*; thus, create this directory first and download lapack-3.6.1.tgz and SuiteSparse-4.5.3.tag.gz into it.

Lapack is one of the earliest _package_ of routines to perform computations in _linear algebra_ (e.g. matrix factorizations, eigenvalues/vectors, linear systems, etc.) and has a widespread usage in computer science. It is for instace called within [Matlab](http://au.mathworks.com/company/newsletters/articles/matlab-incorporates-lapack.html), [Julia](http://docs.julialang.org/en/release-0.4/stdlib/linalg/) and [Numpy](http://docs.scipy.org/doc/numpy-1.10.1/user/install.html). Lapack is also used in several other applications and inspired other tools such as the MKL by Intel. Therefore, calling Lapack from Go is very useful!

### 2.1 Download and uncompress

1. Get *lapack-3.6.1.tgz* from http://www.netlib.org/lapack/lapack-3.6.1.tgz.
2. Save file into *C:\Gosl*
3. Extract files (see commands below)

Start Git Bash and type:

```
cd /c/Gosl
tar xzvf lapack-3.6.1.tgz
mkdir build-lapack
```

### 2.2 Create project with CMake

Start CMake (cmake-gui) and select:

1. Where is the source code = `C:/Gosl/lapack-3.6.1`
2. Where to build the binaries = `C:/Gosl/build-lapack`
3. Hit `[Configure]`
4. Select *MinGW Makefiles* under Specify the generator for this project (leave Use default native compilers on). Hit `[Finish]`
5. Change `CMAKE_INSTALL_PREFIX` = `C:/Gosl`
6. Hit `[Configure]` again
7. Hit `[Generate]` (and close window)

### 2.3 Build Lapack

Continue on Git Bash:

```
cd build-lapack
mingw32-make.exe
mingw32-make.exe test
python lapack_testing.py
mingw32-make.exe install
```

(the error message _python cannot be found_ or _recipe for target test failed_ is OK)



## 3 COMPILE SUITESPARSE

[SuiteSparse](http://faculty.cse.tamu.edu/davis/suitesparse.html) is a collection of routines for computations with sparse matrices and sparse linear systems. Among these routines, *Umfpack* is a powerful (fast) package for dealing with unsymmetrical sparse systems. Umfpack is thus used in Gosl as the main sparse solver. Other routines from Umfpack used in Gosl include the conversion from Triplet to Column-Compressed format.

### 3.1 Download

1. Get SuiteSparse-4.5.3.tar.gz from http://faculty.cse.tamu.edu/davis/SuiteSparse/ and save it into C:\Gosl
2. Get https://github.com/cpmech/gosl/raw/master/scripts/windows/fix-suitesparse.tar.gz and save it into C:\Gosl

### 3.2. Build SuiteSparse

Start Git Bash and type:

```
cd /c/Gosl
tar xzvf SuiteSparse-4.5.3.tar.gz
tar xzvf fix-suitesparse.tar.gz
bash fix-suitesparse/replace-files.bash
cd SuiteSparse
mingw32-make.exe install
cd UMFPACK/Demo
mingw32-make.exe
```

(the small difference causing an error after m test is ok).




## 4 INSTALL GO (GOLANG)

### 4.1 Download Go 1.6.3 (or newer)

Download go1.6.3.windows-amd64.msi from https://golang.org/dl/ and run it. Install into *C:\Go*.

### 4.2 Install from source (patch and compile)

Download, patch and compile Go:

```
cd /c/Gosl
git clone https://go.googlesource.com/go
cd go/src
git fetch https://go.googlesource.com/go refs/changes/70/26670/1
git checkout FETCH_HEAD
export GOROOT_BOOTSTRAP=/c/Go
./all.bat
```

### 4.3 Set environment variables

*Uninstall* go1.6.3, then set the following environment variables:

1. `GOPATH = C:\MyGo`
2. `GOROOT = C:\Gosl\go`
3. Add `C:\Gosl\go\bin` to the `PATH` variable



## 5 INSTALL GOSL

Open Git Bash and type:

```
cd /c/Gosl/go/src
mkdir -p github.com/cpmech
cd github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

*Yay!*
