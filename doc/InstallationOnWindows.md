# Installing Gosl on Windows

Gosl (Go Scientific Library) is a set of routines for computations involving numerical methods and
other mathematical problems in science and engineering.

A number of existent tools are _wrapped_ by Gosl with the two most essential ones being *Lapack* and
*Umfpack* (from SuiteSparse by Prof Tim Davis). The former is a _de facto_ set of routines for
linear algebra and the latter is one of the most powerful tool for computations using sparse
matrices and systems.



## 1 Install PythonXY

In Gosl, Python is used for plotting; but the user does not need to explicitly call Python. Python
is run on the background to generate figures. In Windows, PythonXY is convenient because it has
NumPy, SciPy and Matplotlib (all required) and is easy to install.

Download and install Python(x,y)-2.7.10.0.exe from http://python-xy.github.io/downloads.html. The
default options are OK.



## 2 Install Git and Bash

Install Git and Git Bash from https://git-scm.com/download/

Download and install *Git-2.13.0-64-bit.exe*. Keep default options.



## 3 Install Go Language

### 3.1 Download Go 1.8.1 (or newer)

Download go1.8.1.windows-amd64.msi from https://golang.org/dl/ and proceed with the installation.

Install Go into *C:\Go*.

### 3.2 Set environment variables

Create `C:\MyGo` directory (or any other to be set as `GOPATH`).

Set the following environment variables (press `Windows key + R` to run `sysdm.cpl` then choose
[Advanced] to set Environment Variables):

```
GOPATH = C:\MyGo
```

Remember to close any Git Bash window.


## 4 Download Gosl

Using Git Bash, open a terminal and type:
```bash
mkdir -p /c/MyGo/src/github.com/cpmech
cd /c/MyGo/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
```



## 5a Using Pre-compiled Lapack and SuiteSparse binaries

Download
[Gosl-extra-files-windows.tar.gz](https://sourceforge.net/projects/gosl-installer/files/Gosl-extra-files-windows.tar.gz/download)
into **C:\MyGo\src\github.com\cpmech\gosl**

The, using Git Bash, type:
```bash
cd /c/MyGo/src/github.com/cpmech/gosl
tar xzvf Gosl-extra-files-windows.tar.gz
./all.bash
```

*Finished!*



## 5b [alternative] Compiling Lapack and SuiteSparse from source

On Windows, the following tools for building Lapack and Umfpack are required:

1. Gcc and Gfortran for Windows from http://tdm-gcc.tdragon.net/download. Download and install
   *tdm64-gcc-5.1.0-2.exe* (or newer, as long as it is 64-bit). Make sure to select *fortran* and
   Keep other default options. Click on [Create] and leave the Installation Directory as
   `C:\TDM-GCC-64`
2. CMake for Windows from https://cmake.org/download. Download and install
   *cmake-3.8.1-win64-x64.msi*. Keep default options.
3. Download *lapack-3.7.0.tgz* from http://www.netlib.org/lapack and save into `C:\TDM-GCC-64`
4. Download *SuiteSparse-4.5.5.tar.gz* from http://faculty.cse.tamu.edu/davis/SuiteSparse/ and save
   it into `C:\TDM-GCC-64`

### Compile Lapack

Start Git Bash and type:
```bash
cd /c/TDM-GCC-64
tar xzf lapack-3.7.0.tgz
mkdir build-lapack
```

Start CMake (cmake-gui) and select:

1. Where is the source code = `C:/TDM-GCC-64/lapack-3.7.0`
2. Where to build the binaries = `C:/TDM-GCC-64/build-lapack`
3. Hit `[Configure]`
4. Select *MinGW Makefiles* under Specify the generator for this project (leave Use default native compilers on). Hit `[Finish]`
5. Change `CMAKE_INSTALL_PREFIX` = `C:/TDM-GCC-64`
6. Hit `[Configure]` again
7. Hit `[Generate]` (and close window)

Continue on Git Bash:

```
cd build-lapack
mingw32-make.exe
mingw32-make.exe test
mingw32-make.exe install
```

### Compile SuiteSparse


Start Git Bash and type:
```bash
cd /c/TDM-GCC-64
tar xzf SuiteSparse-4.5.5.tar.gz
bash /c/MyGo/src/github.com/cpmech/gosl/scripts/windows/fix-suitesparse/replace-files.bash
cd SuiteSparse
mingw32-make.exe install
```

TODO: Fix comparison files and check in `cd UMFPACK/Demo` with `mingw32-make.exe`

### Build Gosl

Open Git Bash and type:
```bash
cd /c/MyGo/src/github.com/cpmech/gosl
./all.bash
```

*Finished!*
