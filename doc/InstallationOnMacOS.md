# Installing Gosl on macOS

<div id="container">
<p>
<a href="https://github.com/cpmech/gosl/blob/master/doc/InstallationOnMacOS.md"><img src="icon-macos.png"></a>
</p>
</div>

Tested on macOS Mojave 10.14 with Xcode 10.3

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

Install Go for macOS from https://golang.org (Tested with Go 1.12.7)

## Gosl

Type the following commands:

```
cd ~
mkdir -p mygo/src/github.com/cpmech
cd ~/mygo/src/github.com/cpmech/
git clone https://github.com/cpmech/gosl.git
cd gosl
bash all.bash
```

Enjoy!
