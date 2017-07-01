# Installing Gosl on Ubuntu Linux (Debian)

## 1. [Required] Install some dependencies:

Type:
```bash
sudo apt-get install libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev 
sudo apt-get install gfortran libvtk6-dev python-scipy python-matplotlib dvipng
sudo apt-get install libfftw3-dev libfftw3-mpi-dev
```

## 2. [Required] Install OpenBLAS

**Note**: make sure `libopenblas-base` is NOT installed (from apt-get).

Type:
```bash
mkdir -p $HOME/xpkg && cd $HOME/xpkg
git clone https://github.com/xianyi/OpenBLAS.git
cd OpenBLAS
make DYNAMIC_ARCH=1 -j4
sudo make DYNAMIC_ARCH=1 PREFIX=/usr/local install
```

To avoid conflict between BLAS libraries, type:
```bash
sudo apt-get remove libopenblas-base
sudo update-alternatives --remove libblas.so.3 /usr/local/lib/libopenblas.so.0
sudo update-alternatives --list libblas.so.3
sudo update-alternatives --install /usr/lib/libblas.so.3 libblas.so.3 /usr/local/lib/libopenblas.so.0 41
sudo update-alternatives --config libblas.so.3
```
In the last step, check that OpenBLAS is selected in auto mode.

See more [information here](https://github.com/xianyi/OpenBLAS/wiki/faq#debianlts)

## 3. [Optional] Install Intel MKL

Download MKL (~900Mb) from [the intel MKL website](https://software.intel.com/en-us/intel-mkl)
(click on Free Download; need to sign-in), then:
```bash
mkdir -p $HOME/xpkg && cd $HOME/xpkg
tar xzf l_mkl_2017.3.196.tgz
cd l_mkl_2017.3.196/
bash install_GUI.sh
```
and follow the instructions. These options have been tested:
1. Choose _Install as root using sudo_
2. Keep default install location: **/opt/intel**

## 4. [Required] Clone and install Gosl

Type:
```bash
mkdir -p ${GOPATH%:*}/src/github.com/cpmech
cd ${GOPATH%:*}/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

## 5. [Optional] Test la/mkl and img/ocv subpackages

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
