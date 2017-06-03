# Gosl. chk. Check code and unit test tools

[![GoDoc](https://godoc.org/github.com/cpmech/gosl/chk?status.svg)](https://godoc.org/github.com/cpmech/gosl/chk)

[![GoDoc](https://godoc.org/github.com/cpmech/gosl/chk?status.svg)](https://godoc.org/github.com/cpmech/gosl/chk) 

More information is available in [the documentation of this package](https://godoc.org/github.com/cpmech/gosl/chk).

Package `chk` provides tools to check numerical results and to perform unit tests.

This package also contains `assert` functions.

## Examples

Checking that two matrices are equal to each other, by some tolerance (useful in unit tests):

```go
// a hypotetical matrix obtained from a "numerical" solution
// notice "noise" on the last component
Anumerical := [][]float64{
    {76, 35, 64, 8, 35},
    {92, 37, 16, 0, 46},
    {64, 36, 60, 27, 42},
    {41, 20, 21, 35, 45},
    {4, 48, 37, 87, 9 + 1e-9},
}

// a hypotetical matrix obtained from a "closed-form" solution
Aanalytical := [][]float64{
    {76, 35, 64, 8, 35},
    {92, 37, 16, 0, 46},
    {64, 36, 60, 27, 42},
    {41, 20, 21, 35, 45},
    {4, 48, 37, 87, 9},
}

// allocate testing structure, just for this example
tst := &testing.T{}

// tolerance for comparison
//tolerance := 1e-10 // this makes the test to fail
tolerance := 1e-8 // this allows test to pass

// compare matrices
chk.Matrix(tst, "A", tolerance, Anumerical, Aanalytical)
```
