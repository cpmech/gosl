# Compiling on Windows

# 1 Git+Bash
<ol>
<li>Install *Git+Bash* from https://git-scm.com/download/win **IMPORTANT** install Git+Bash on `C:\Git` or `D:\Git` but NOT on "C:\Program Files" because MinGW fails due to the space in "Program Files."</li>
</ol>

# 2 TDM64-GCC and Gfortran
<ol>
<li>Download and Install TDM64-GCC (tdm64-gcc-5.1.0-2.exe) from http://tdm-gcc.tdragon.net/download</li>
<li>IMPORTANT: Make sure to select the 64-bit version</li>
<li>IMPORTANT: Make sure to select **fortran** under **gcc** (TDM64 Current: 5.1.0-tdm64-1)</li>
<li>Install TDM64-GCC to `C:\TDM-GCC-64`</li>
<li>IMPORTANT: Duplicate **mingw32-make.exe** located in C:\TDM-GCC-64\bin an rename it to **make** (without .exe)</li>
</ol>

# 3 Go
<ol>
<li>Download and Install Go (go1.10.windows-amd64.msi) from https://golang.org/dl/</li>
<li>Create the directory `C:\MyGo`</li>
<li>Set environment variable `GOPATH = C:\MyGo` (press WindowsKey + R or run sysdm.cpl)</li>
<li>Close any Git+Bash terminals</li>
</ol>

# 4. OpenBLAS
<ol>
<li>Download OpenBLAS source code (OpenBLAS-v0.2.19-Win64-int32.zip) from https://sourceforge.net/projects/openblas/files/v0.2.19/ and save it into your Downloads directory</li>
<li>Extract all files in a temporary directory; e.g. in the Downloads directory</li>
<li>Copy all files in the temporary `include` directory of OpenBLAS...zip to `C:\TDM-GCC-64\include`</li>
<li>Copy **libopenblas.a** file from the `lib` directory of OpenBLAS...zip to `C:\TDM-GCC-64\lib`</li>
</ol>

# 5. Gosl - Part 1: scripts to patch SuiteSparse
<ol>
<li>Assuming that your Go code will be located at `C:\MyGo`, download Gosl using the following commands (in a Git+Bash terminal)</li>
</ol>
```bash
mkdir -p /c/MyGo/src/github.com/cpmech
cd /c/MyGo/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
```

# 6. SuiteSparse
<ol>
<li>Download SuiteSparse source code (SuiteSparse-5.2.0.tar.gz) from http://faculty.cse.tamu.edu/davis/SuiteSparse/ and save it into `C:\TDM-GCC-64`</li>
<li>Run the following commands (in a Git+Bash terminal)</li>
</ol>
```bash
cd /c/TDM-GCC-64
tar xzf SuiteSparse-5.2.0.tar.gz
bash /c/MyGo/src/github.com/cpmech/gosl/scripts/windows/fix-suitesparse/replace-files.bash
cd SuiteSparse
make install
```

# 7. FFTW
<ol>
<li>Download FFTW source code (fftw-3.3.7.tar.gz) from http://www.fftw.org/download.html and save it into `C:\TDM-GCC-64`</li>
<li>Run the following commands (in a Git+Bash terminal)</li>
</ol>
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
<ol>
<li>Download and install Cmake (cmake-3.11.0-win64-x64.msi) from https://cmake.org/download/</li>
<li>Proceed with the straightforward installation of CMake</li>
</ol>

# 9 HDF5
<ol>
<li>Download HDF5 source code (hdf5-1.8.16.tar.bz2) from https://support.hdfgroup.org/ftp/HDF5/releases/hdf5-1.8/hdf5-1.8.16/src/ and save it into `C:\TDM-GCC-64`</li>
<li>Run the following commands (in a Git+Bash terminal)</li>
</ol>
```bash
cd /c/TDM-GCC-64
tar xjf hdf5-1.8.16.tar.bz2
```
<ol>
<li>Load CMake by clicking on its icon (Start Menu)</li>
<li>Set `Where is the source code = C:/TDM-GCC-64/hdf5-1.8.16`</li>
<li>Set `Where to build the binaries = C:/TDM-GCC-64/build_hdf5`</li>
<li>Click **Configure**</li>
<li>Select **MinGW Makefiles** under _Specify the generator for this project_</li>
<li>Leave _Use default native compilers_ selected and click Finish</li>
<li>Unselect BUILD_SHARED_LIBS, BUILD_TESTING, and HDF5_BUILD_TOOLS</li>
<li>Set `CMAKE_INSTALL_PREFIX = C:/TDM-GCC-64`</li>
<li>Click **Configure** again and then click **Generate**</li>
<li>Run the following commands (in a Git+Bash terminal)</li>
</ol>
```bash
cd /c/TDM-GCC-64/build_hdf5
make install
```

# 10 Gosl - Part 2: compilation and testing
<ol>
<li>Now we can compile Gosl. Type the following commands (in a Git+Bash terminal)</li>
</ol>
```bash
cd /c/MyGo/src/github.com/cpmech/gosl
./all.bash
```

*Finished!*