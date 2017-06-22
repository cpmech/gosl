# Installing Gosl on Ubuntu Linux (Debian)

## 1. [Required] Install some dependencies:

```bash
sudo apt-get install libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev 
sudo apt-get install gfortran libvtk6-dev python-scipy python-matplotlib dvipng
sudo apt-get install libfftw3-dev libfftw3-mpi-dev
```

## 2. [Required] Set dynamic library flags

To set LD\_LIBRARY\_PATH, add the following line to `.bashrc` or `.bash_aliases`:
```bash
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```
Alternatively, change `/etc/ld.so.conf` file as appropriate.

## 3. [Required] Install OpenBLAS

Type:
```bash
mkdir -p $HOME/xpkg && cd $HOME/xpkg
git clone https://github.com/xianyi/OpenBLAS.git
cd OpenBLAS
make -j4
sudo make PREFIX=/usr/local install
```

**Note**: Make sure to set the `/usr/local/lib` directory as a searchable LD\_LIBRARY\_PATH.
Otherwise, the following error may happen:
```
[...] libopenblas.so.0: cannot open shared object file: [...]
```

## 4. [Optional] Install Intel MKL

Download MKL (~900Mb) from [the intel MKL website](https://software.intel.com/en-us/intel-mkl)
(click on Free Download; need to sign-in), then:
```bash
mkdir -p $HOME/xpkg && cd $HOME/xpkg
tar xzvf l_mkl_2017.2.174.tgz
cd l_mkl_2017.2.174/
bash install_GUI.sh
```
and follow the instructions. These options have been tested:
1. Choose _Install as root using sudo_
2. Keep default install location: **/opt/intel**

## 5. [Optional] Install OpenCV

Type:
```bash
sudo apt-get install libgtk2.0-dev pkg-config libavcodec-dev libavformat-dev libswscale-dev
mkdir -p $HOME/xpkg && cd $HOME/xpkg
git clone https://github.com/opencv/opencv.git
mkdir build_opencv
cd build_opencv
ccmake ../opencv
```
press `[c][c][g]`
```
make -j4
sudo make install
```

## 6. [Required] Clone and install Gosl

Type:
```bash
mkdir -p ${GOPATH%:*}/src/github.com/cpmech
cd ${GOPATH%:*}/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

## 7. [Optional] Test la/mkl and la/opencv subpackages

Install and test subpackage `la/mkl`:
```bash
cd ${GOPATH%:*}/src/github.com/cpmech/la/mkl
go install
go test
```

Install and test subpackge `img/ocv`:
```bash
cd ${GOPATH%:*}/src/github.com/cpmech/img/ocv
go install
go test
```
