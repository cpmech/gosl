# Installing Gosl on Ubuntu Linux (Debian)

## 0. [Required] Install Go language

[See instructions here](https://github.com/cpmech/gosl/blob/master/doc/InstallAndTestGo.md)

## 1. [Required] Install some dependencies:

Type:
```bash
sudo apt-get install libopenmpi-dev libhwloc-dev libsuitesparse-dev libmumps-dev 
sudo apt-get install gfortran libvtk6-dev python-scipy python-matplotlib dvipng
sudo apt-get install libfftw3-dev libfftw3-mpi-dev libmetis-dev
sudo apt-get install liblapacke-dev libopenblas-dev
```

## 2. [Optional] Install Intel MKL

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

## 3. [Required] Clone and install Gosl

Type:
```bash
mkdir -p ${GOPATH%:*}/src/github.com/cpmech
cd ${GOPATH%:*}/src/github.com/cpmech
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

## 4. [Optional] Test la/mkl subpackage

Install and test subpackage `la/mkl`:
```bash
cd ${GOPATH%:*}/src/github.com/cpmech/la/mkl
go install
go test
```
