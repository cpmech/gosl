# Installing Gosl on Windows

Below, two options to install Gosl (and Gcc) are presented:
* __Option A__: Quick install using pre-compiled code (Windows 10 Installer)
* __Option B__: Installing from sources (i.e. compiling Lapack and SuiteSparse)

### [Strongly Recommended] Install PythonXY

1. Download and install *Python(x,y)-2.7.10.0.exe* from http://python-xy.github.io/downloads.html

### [Recommended] Install Git and Bash

2. Download and install *Git-2.13.1.2-64-bit.exe* from https://git-scm.com/download/

**NOTE**: You must install Git in `C:\Git` or `D:\Git` but **NOT** within a "Program Files" directory (the problem here are the spaces between "Program" and "Files" that FFTW doesn't like...).

All the other default options are OK.

### [Recommended] Install Visual Studio Code

3. Download and install VS-Code from https://code.visualstudio.com/docs/?dv=win

Steps to Install [Go extension for VS Code](https://marketplace.visualstudio.com/items?itemName=lukehoban.Go):
* Press Ctrl+P
* Type `ext install Go`
* Select the version by *lukehoban*, click install, and click reload

### [Required] Install Go

4. Download and install *go1.9.1.windows-amd64.msi* (or newer) from https://golang.org/dl/

### [Option A] Quick Install (Windows 10 Installer)

The _Windows 10 Installer_ will install Gosl, TDM Gcc64, compiled SuiteSparse and FFTW,
and will set the required environment variables automatically.

5. Download and install the [gosl installer from here](https://sourceforge.net/projects/gosl-installer/files/)

After the installation, the Gosl code can be updated with the current version from GitHub.
Also the environment variables can be replaced as desired.

*Finished!*



------------------------------------------------------------------------------------------------------------------------------------

## [Option B] Installing from Sources

Four tools are neccessary/recommended to work with Go/Gosl:
* Python is used for plotting but the user does not need to explicitly call Python; it
  is run on the background to generate figures.
* Lapack and SuiteSparse are required by the Linear Algebra package.
* Git and Git Bash are convenient to work with Git Version Control System
* Visual Studio Code (VS-Code) is a great tool to develop Go code (and more).

*NOTE*: Git Bash is required for the commands below. Thus, see the section above regarding the installation of Git and Bash.

### Install Go and Set Environment Variable

Download and install *go1.9.1.windows-amd64.msi* (or newer) from https://golang.org/dl/

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

Download and Install TDM-GCC-64 with Gcc and Gfortran for Windows from http://tdm-gcc.tdragon.net/download

1. Download and install *tdm64-gcc-5.1.0-2.exe* (or newer, it must be **64-bit**)
2. Make sure to select **fortran** (under gcc) and keep other default options
3. Click on [Create] and leave the Installation Directory as `C:\TDM-GCC-64`
4. Duplicate `mingw32-make.exe` located in `C:\TDM-GCC-64\bin` into a file named `make` (without .exe).

### Download and Install OpenBLAS

Download OpenBLAS binaries for Windows from here https://sourceforge.net/projects/openblas/files/

For example, download OpenBLAS-v0.2.19-Win64-int32.zip or newer.

1. Extract all files from the **include** directory of OpenBLAS...zip into `C:\TDM-GCC-64\include`
2. Extract the `libopenblas.a` file from the **lib** directory of OpenBLAS...zip into `C:\TDM-GCC-64\lib`

### Download and Compile SuiteSparse

Download *SuiteSparse-4.5.5.tar.gz* from http://faculty.cse.tamu.edu/davis/SuiteSparse/ and save it into `C:\TDM-GCC-64`

Start Git Bash and type:
```bash
cd /c/TDM-GCC-64
tar xzf SuiteSparse-4.5.5.tar.gz
bash /c/MyGo/src/github.com/cpmech/gosl/scripts/windows/fix-suitesparse/replace-files.bash
cd SuiteSparse
make install
```

### Download and Compile FFTW

Download *fftw-3.3.6-pl2.tar.gz* from http://www.fftw.org/ and save it into `C:\TDM-GCC-64`

Configure and compile FFTW with the following commands:
```bash
cd /c/TDM-GCC-64
tar xzf fftw-3.3.6-pl2.tar.gz
cd fftw-3.3.6-pl2
./configure --disable-alloca --with-our-malloc --disable-shared --enable-static --enable-sse2 --with-incoming-stack-boundary=2
make
```

If an error such as "C:\Program cannot be found" happens, probably you forgot to install Git Bash on C:\Git or D:\Git (Git must not be in "Program Files"). See instructions above.

Copy the following files to the `C:\TDM-GCC-64` directories:
```bash
cp api/fftw3.h ../include/
cp .libs/libfftw3.a ../lib/
```

### Build Gosl

Open Git Bash and type:
```bash
cd /c/MyGo/src/github.com/cpmech/gosl
./all.bash
```

*Finished!*
