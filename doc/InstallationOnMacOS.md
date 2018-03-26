# Installing Gosl on macOS

The installation on macOS includes:

1. Xcode
2. HomeBrew and dependencies
3. Go language
4. Gosl 
 
## 1 Xcode

In a terminal, type the following commands and click install to install the *command line developer tools*:

```
xcode-select --install
```

## 2 HomeBrew and dependencies

Download and install HomeBrew from https://brew.sh

```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Install dependencies with:

```
brew install fftw
brew install homebrew/science/openblas
brew install suite-sparse
brew install openmpi
```

## 3 Go language

Install Go for macOS from https://golang.org

Edit `/Users/name/.bash_profile` (or `/Users/name/.profile`) and add:

```
export GOPATH=/Users/name/mygo
```

(or another location as desired)

## 4 Gosl
 
Type the following commands:

```
cd mygo
mkdir -p src/github.com/cpmech
cd src/github.com/cpmech/
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

Yay!
