# Gosl. la/oblas. Wrapper to OpenBLAS

[![GoDoc](https://godoc.org/github.com/cpmech/gosl/la/oblas?status.svg)](https://godoc.org/github.com/cpmech/gosl/la/oblas) 

More information is available in **[the documentation of this package](https://godoc.org/github.com/cpmech/gosl/la/oblas).**

This subpackge implements a light wrapper to OpenBLAS. Therefore, its routines are a little more
_lower level_ than the ones in the parent package `la`.

[Check also OpenBLAS](https://github.com/xianyi/OpenBLAS).



## Examples

### Axpy

```go
α := 0.5
x := []float64{20, 10, 30, 123, 123}
y := []float64{-15, -5, -24, 666, 666, 666}
n, incx, incy := 3, 1, 1
err := Daxpy(n, α, x, incx, y, incy)
if err != nil {
    tst.Errorf("Daxpy failed:\n%v\n", err)
    return
}

chk.Vector(tst, "x", 1e-15, x, []float64{20, 10, 30, 123, 123})
chk.Vector(tst, "y", 1e-15, y, []float64{-5, 0, -9, 666, 666, 666})
```

### Zaxpy

```go
α := 1.0 + 0i
x := []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i}
y := []complex128{-15 + 1.5i, -5 - 2i, -24 + 1i, 666 - 0.5i, 666 + 5i}
n, incx, incy := len(x), 1, 1
err := Zaxpy(n, α, x, incx, y, incy)
if err != nil {
    tst.Errorf("Daxpy failed:\n%v\n", err)
    return
}

chk.VectorC(tst, "x", 1e-15, x, []complex128{20 + 1i, 10 + 2i, 30 + 1.5i, -123 + 0.5i, -123 + 0.5i})
chk.VectorC(tst, "y", 1e-15, y, []complex128{5 + 2.5i, 5, 6 + 2.5i, 543, 543 + 5.5i})
```

### Dgemv

```go
m, n := 4, 3
a := NewMatrix(m, n)
a.SetFromMat([][]float64{
    {0.1, 0.2, 0.3},
    {1.0, 0.2, 0.3},
    {2.0, 0.2, 0.3},
    {3.0, 0.2, 0.3},
})

α, β := 0.5, 2.0
x := []float64{20, 10, 30}
y := []float64{3, 1, 2, 4}
lda, incx, incy := m, 1, 1
err := Dgemv(false, m, n, α, a, lda, x, incx, β, y, incy)
if err != nil {
    tst.Errorf("Dgemv failed:\n%v\n", err)
    return
}
chk.Vector(tst, "y", 1e-15, y, []float64{12.5, 17.5, 29.5, 43.5})
```
