# Installing Gosl on Windows

Four tools are neccessary/recommended to work with Go/Gosl:
* Python is used for plotting but the user does not need to explicitly call Python; it
  is run on the background to generate figures.
* Lapack and SuiteSparse are required by the Linear Algebra package.
* Git and Git Bash are convenient to work with Git Version Control System
* Visual Studio Code (VS-Code) is a great tool to develop Go code (and more).

Below, two options to install Gosl (and Gcc) are presented:
* __Option A__: Quick install using pre-compiled code (Windows 10 Installer)
* __Option B__: Installing from sources (i.e. compiling Lapack and SuiteSparse)

## [Strongly Recommended] Install PythonXY

1. Download and install *Python(x,y)-2.7.10.0.exe* from http://python-xy.github.io/downloads.html

## [Recommended] Install Git and Bash

2. Download and install *Git-2.13.0-64-bit.exe* from https://git-scm.com/download/

## [Recommended] Install Visual Studio Code

3. Download and install VS-Code from https://code.visualstudio.com/docs/?dv=win

Steps to Install [Go extension for VS Code](https://marketplace.visualstudio.com/items?itemName=lukehoban.Go):
* Press Ctrl+P
* Type `ext install Go`
* Select the version by *lukehoban*, click install, and click reload

## [Option A] Quick Install (Windows 10 Installer)

The _Windows 10 Installer_ will install Gosl, TDM Gcc64, compiled Lapack and SuiteSparse,
and will set the required environment variables automatically.

4. Download and install *go1.8.3.windows-amd64.msi* (or newer) from https://golang.org/dl/
5. Download and install the [gosl installer from here](https://sourceforge.net/projects/gosl-installer/files/)

After the installation, the Gosl code can be updated with the current version from GitHub.
Also the environment variables can be replaced as desired.

*Finished!*



## [Option B] Installing from Sources

Skip these steps if you have used the _installer_ already.

### Install Go and Set Environment Variable

Download and install *go1.8.3.windows-amd64.msi* (or newer) from https://golang.org/dl/

Create `C:\MyGo` directory (or any other to be set as `GOPATH`).

Set the following environment variables (press `Windows key + R` to run `sysdm.cpl` then choose
[Advanced] to set Environment Variables):

```
GOPATH = C:\MyGo
```

Remember to close any Git Bash window.

### Download Gosl

Using Git Bash, open a terminal and type:
```bash
mkdir -p /c/MyGo/src/github.com/cpmech
cd /c/MyGo/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
```

### Download and Install Compilation Tools

Download and Install:
1. TDM-GCC-64 with Gcc and Gfortran for Windows from http://tdm-gcc.tdragon.net/download
   1. Download and install *tdm64-gcc-5.1.0-2.exe* (or newer, as long as it is 64-bit)
   2. Make sure to select *fortran* and keep other default options
   3. Click on [Create] and leave the Installation Directory as `C:\TDM-GCC-64`
2. CMake for Windows from https://cmake.org/download
   1. Download and install *cmake-3.8.1-win64-x64.msi*. Keep default options

### Download and Compile Lapack

Download *lapack-3.7.0.tgz* from http://www.netlib.org/lapack and save it into `C:\TDM-GCC-64`

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

### Download and Compile SuiteSparse

Download *SuiteSparse-4.5.5.tar.gz* from http://faculty.cse.tamu.edu/davis/SuiteSparse/ and save it into `C:\TDM-GCC-64`

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
