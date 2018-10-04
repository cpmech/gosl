# Installing Gosl on macOS

<div id="container">
<p>
<a href="https://github.com/cpmech/gosl/blob/master/doc/InstallationOnMacOS.md"><img src="icon-macos.png"></a>
</p>
</div>

Tested on macOS Mojave 10.14 with Xcode 10.1 beta 2
 
## HomeBrew and dependencies

Download and install HomeBrew from https://brew.sh

```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Install dependencies with:

```
brew install gcc fftw hdf5
brew install openblas suite-sparse openmpi
```

Useful commands:
```
brew update
brew upgrade
brew cleanup
```

## Go language

Install Go for macOS from https://golang.org (Tested with Go 1.11.1)

Edit `/Users/name/.bash_profile` (or `/Users/name/.profile`) and add:

```
export GOPATH=/Users/name/mygo
```

(or another location as desired)

[More information is available here](https://github.com/cpmech/gosl/blob/master/doc/InstallAndTestGo.md)

## Gosl
 
Type the following commands:

```
cd ~
mkdir -p mygo/src/github.com/cpmech
cd ~/mygo/src/github.com/cpmech/
git clone https://github.com/cpmech/gosl.git
cd gosl
./all.bash
```

Enjoy!
