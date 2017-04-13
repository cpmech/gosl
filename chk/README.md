# Gosl. chk. Check code and unit test tools

More information is available in [the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxchk.html).

Package `chk` provides tools to check numerical results and to perform unit tests.

This package also contains `assert` functions.

## Examples

Checking that a matrix, allocated in the `la` package, is correct.
```go
a := la.MatAlloc(3, 5)
a[0][0] = 1
a[0][1] = 2
a[0][2] = 3
a[0][3] = 4
a[0][4] = 5
a[1][0] = 0.1
a[1][1] = 0.2
a[1][2] = 0.3
a[1][3] = 0.4
a[1][4] = 0.5
a[2][0] = 10
a[2][1] = 20
a[2][2] = 30
a[2][3] = 40
a[2][4] = 50
chk.Matrix(tst, "a", 1e-17, a, [][]float64{
    {1, 2, 3, 4, 5},
    {0.1, 0.2, 0.3, 0.4, 0.5},
    {10, 20, 30, 40, 50},
})
```
