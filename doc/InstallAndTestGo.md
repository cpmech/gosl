# Installing and testing Go

## Official instructions

First, please [read official instructions here](https://golang.org/doc/install)

## Set GOROOT and GOPATH

Assuming that Go is installed in `$HOME/xpkg` and your Go code is in `$HOME/mygo`,
add these to your `.bash_aliases`

```
export GOROOT=$HOME/xpkg/go
export GOPATH=$HOME/mygo
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

## Create a new library

Create this directory

```
mkdir -p $HOME/mygo/src/mylib
cd $HOME/mygo/src/mylib/
```

Create these files in `mylib`

**mylib.go**
```
package mylib

import (
	"fmt"
)

func SayHello() bool {
	fmt.Printf("Hello World from mylib!\n")
	return true
}
```

**t_mylib_test**

```
package mylib

import (
	"testing"
)

func TestMylib(tst *testing.T) {
	res := SayHello()
	if !res {
		tst.Errorf("test failed\n")
	}
}
```

## Run test

Run test

```
cd $HOME/mygo/src/mylib/
go test
```

Output

```
Hello World from mylib!
PASS
ok      mylib    0.001s
```
