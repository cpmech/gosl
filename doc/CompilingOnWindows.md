# Compiling on Windows

# 1 Git+Bash

1. Install *Git+Bash* from https://git-scm.com/download/win **IMPORTANT** install Git+Bash on `C:\Git` or `D:\Git` but NOT on "C:\Program Files" because MinGW fails due to the space in "Program Files." [Screenshot of suggested "Select Destination Location"](https://github.com/cpmech/gosl/blob/master/doc/git-installation-windows-must-be-on-c-drive.png)



# 2 TDM64-GCC and Gfortran

1. Download and Install TDM64-GCC (tdm64-gcc-5.1.0-2.exe) from http://tdm-gcc.tdragon.net/download
2. IMPORTANT: Make sure to select the 64-bit version
3. IMPORTANT: Make sure to select **fortran** under **gcc** (TDM64 Current: 5.1.0-tdm64-1) [Screenshot of "New Installation: Choose Components"](https://github.com/cpmech/gosl/blob/master/doc/tdm64-gcc-with-fortran.png)
4. Install TDM64-GCC to `C:\TDM-GCC-64`
5. IMPORTANT: Duplicate **mingw32-make.exe** located in C:\TDM-GCC-64\bin an rename it to **make** (without .exe)



# 3 Go

1. Download and Install Go (go1.10.windows-amd64.msi) from https://golang.org/dl/
2. Create the directory `C:\MyGo`
3. Set environment variable `GOPATH = C:\MyGo` (press WindowsKey + R or run sysdm.cpl) [Screenshot of "Setting Environment Variables"](https://github.com/cpmech/gosl/blob/master/doc/go-environment-variables.png)
4. Close any Git+Bash terminals



# 4 OpenBLAS

1. Download OpenBLAS source code (OpenBLAS-v0.2.19-Win64-int32.zip) from https://sourceforge.net/projects/openblas/files/v0.2.19/ and save it into your Downloads directory
2. Extract all files in a temporary directory; e.g. in the Downloads directory
3. Copy all files in the temporary `include` directory of OpenBLAS...zip to `C:\TDM-GCC-64\include`
4. Copy **libopenblas.a** file from the `lib` directory of OpenBLAS...zip to `C:\TDM-GCC-64\lib`



# 5 Gosl - Part 1: scripts to patch SuiteSparse

1. Assuming that your Go code will be located at `C:\MyGo`, download Gosl using the following commands (in a Git+Bash terminal)

```bash
mkdir -p /c/MyGo/src/github.com/cpmech
cd /c/MyGo/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
```



# 6 SuiteSparse

1. Download SuiteSparse source code (SuiteSparse-5.2.0.tar.gz) from http://faculty.cse.tamu.edu/davis/SuiteSparse/ and save it into `C:\TDM-GCC-64`
2. Run the following commands (in a Git+Bash terminal)

```bash
cd /c/TDM-GCC-64
tar xzf SuiteSparse-5.2.0.tar.gz
bash /c/MyGo/src/github.com/cpmech/gosl/scripts/windows/fix-suitesparse/replace-files.bash
cd SuiteSparse
make install
```



# 7 FFTW

1. Download FFTW source code (fftw-3.3.7.tar.gz) from http://www.fftw.org/download.html and save it into `C:\TDM-GCC-64`
2. Run the following commands (in a Git+Bash terminal)

```bash
cd /c/TDM-GCC-64
tar xzf fftw-3.3.7.tar.gz
cd fftw-3.3.7
./configure --disable-alloca --with-our-malloc --disable-shared
--enable-static --enable-sse2 --with-incoming-stack-boundary=2
make
cp api/fftw3.h ../include/
cp .libs/libfftw3.a ../lib/
```



# 8 CMake

1. Download and install Cmake (cmake-3.11.0-win64-x64.msi) from https://cmake.org/download/
2. Proceed with the straightforward installation of CMake



# 9 HDF5

1. Download HDF5 source code (hdf5-1.8.16.tar.bz2) from https://support.hdfgroup.org/ftp/HDF5/releases/hdf5-1.8/hdf5-1.8.16/src/ and save it into `C:\TDM-GCC-64`
2. Run the following commands (in a Git+Bash terminal)

```bash
cd /c/TDM-GCC-64
tar xjf hdf5-1.8.16.tar.bz2
```

1. Load CMake by clicking on its icon (Start Menu)
2. Set `Where is the source code = C:/TDM-GCC-64/hdf5-1.8.16`
3. Set `Where to build the binaries = C:/TDM-GCC-64/build_hdf5`
4. Click **Configure**
5. Select **MinGW Makefiles** under _Specify the generator for this project_
6. Leave _Use default native compilers_ selected and click Finish
7. Unselect BUILD_SHARED_LIBS, BUILD_TESTING, and HDF5_BUILD_TOOLS
8. Set `CMAKE_INSTALL_PREFIX = C:/TDM-GCC-64`
9. Click **Configure** again and then click **Generate**
10. Run the following commands (in a Git+Bash terminal)

```bash
cd /c/TDM-GCC-64/build_hdf5
make install
```



# 10 Gosl - Part 2: compilation and testing

1. Now we can compile Gosl. Type the following commands (in a Git+Bash terminal)

```bash
cd /c/MyGo/src/github.com/cpmech/gosl
./all.bash
```



*Finished!*
