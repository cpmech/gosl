# Gosl. la. Linear Algebra and efficient sparse solvers

More information is available in [the documentation of this package](http://rawgit.com/cpmech/gosl/master/doc/xxla.html).

The `la` subpackage implements functions to perform linear algebra operations such as vector
addition or matrix-vector multiplications. It also wraps some efficient sparse linear systems
solvers such as [Umfpack](http://faculty.cse.tamu.edu/davis/suitesparse.html) and
[MUMPS](http://mumps.enseeiht.fr) (**not** the _parotitis_ disease!).

Both [Umfpack](http://faculty.cse.tamu.edu/davis/suitesparse.html) and
[MUMPS](http://mumps.enseeiht.fr) solvers are very efficient!

The vector and matrix operations implemented in the `la` package do not take advantage of the best
capabilities of the computer hardware as in other libraries such as
[OpenBlas](http://www.openblas.net) or [IntelMKL](https://software.intel.com/en-us/intel-mkl).

Therefore, `la` should be used for _small_ to _normal-sized_ vectors and matrices. The other
subpackage `lax` provides a wrapper to [OpenBlas](http://www.openblas.net) and should be used for
_very-large_ or _large_ vectors and matrices (TODO).

TODO: The reason for having these two packages is that `la` is probably more efficient for smaller systems
such as the assemblage of stiffness matrices in finite element computations. On the other hand,
`lax` involves copies of data to be passed to the lower level routines and might add some overhead;
however the greater efficiency compensates this overhead in the end.

There is only one type of integer and one type of real number (float point number, FPN) in this package.
Integers are `Int` and FPNs are `float64` (i.e. doubles). Also, there is only one type of complex
number corresponding to two doubles (i.e. both the real and imaginary numbers are float64).

Function names to compute on vectors begin with `Vec` whereas functions to compute on matrices have
a `Mat` prefix. Function names for sparse matrices begin with `Sp`. Some functions have a suffix `C` for
complex numbers. Functions that operate on vectors and matrices are named with a `MatVec` prefix.

For example, the following functions operate on vectors:
1. `VecAccum`, `VecAdd`, `VecAdd2`, `VecClone`, `VecCopy`
2. `VecDot`, `VecFill`, `VecFillC`, `VecLargest`, `VecMax`
3. `VecMaxDiff`, `VecMaxDiffC`, `VecMin`, `VecMinMax`, `VecNorm`
4. `VecNormDiff`, `VecOuterAdd`, `VecRms`, `VecRmsErr`,
5. `VecRmsError`, `VecScale`, `VecScaleAbs`

The following functions operate on matrices:
1. `MatAlloc` `MatClone` `MatCopy` `MatFill`
2. `MatGetCol` `MatLargest` `MatMaxDiff`
3. `MatMul` `MatMul3` `MatNormF` `MatNormI`
4. `MatScale` `MatSetDiag` `MatToColMaj`

And the following functions operate on matrices and vectors (note the use of `Tr` for the transposed
version):
1. `MatTrMul3`, `MatTrMulAdd3`, `MatTrVecMul`
2. `MatTrVecMulAdd`, `MatVecMul`, `MatVecMulAdd`
3. `MatVecMulAddC`, `MatVecMulCopyAdd`

There are two kinds of matrix inversion functions:
1. `MatInv` for small and square matrices (currently smaller than 4x4)
2. `MatInvG` for matrices of any size, including rectangular ones (pseudo-inverse). This version
   shouldn't be too slow since it calls [LAPACK](http://www.netlib.org/lapack)

Other useful functions are:
1. `Cholesky` native Go implementation of the Cholesky decomposition
2. `SPDsolve` to solve a Symmetric/Positive-Definite system after the Cholesky decomposition
3. `MatSvd` wrapper to [LAPACK](http://www.netlib.org/lapack) SVD decomposition
4. `MatCondG` to compute the condition number of a matrix


## Structures for sparse problems

In `la`, there are two types of structures to hold data for solving a sparse linear system:
_triplet_ and _column-compressed matrix_. Both have a complex version and all are named as follows:
1. `Triplet` and `TripletC` (complex)
2. `CCMatrix` and `CCMatrixC` (complex)

The `Triplet` is more convenient for data input; e.g. during the assemblage step in a finite element
code. The `CCMatrix` is better (faster) for computations and is the structure given to Umfpack.

Triplets are initiated by giving the size of the corresponding matrix and the number of non-zero
entries (components). Thus, only space for the non-zero entries is allocated. Afterwards, each
component can be set by calling the `Put` method after calling the `Start` method. The `Start`
method can be called several times. Therefore, this structure helps with the efficient
implementation of codes requiring multiple solutions of linear systems as in a finite element
analysis.

To convert the Triplet into a Column-Compressed Matrix, call the `ToMatrix` method of `Triplet`. And
to convert the sparse matrix to a dense version (e.g. for reporting/printing), call the `ToDense`
method of `CCMatrix`.

The `Triplet` has also a very convenient method to copy the contents of another Triplet (i.e. sparse
matrix) into the positions starting with the maximum number of columns of this second matrix and the
positions starting with the maximum number of rows of this second matrix. To this end, we use the
method `PutMatAndMatT`. This operation is particularly common when solving mechanical problems with
Lagrange multipliers and can be illustrated as follows:
```
[... ... ... a00 a10 ...] 0
[... ... ... a01 a11 ...] 1
[... ... ... a02 a12 ...] 2      [. at  .]
[a00 a01 a02 ... ... ...] 3  =>  [a  .  .]
[a10 a11 a12 ... ... ...] 4      [.  .  .]
[... ... ... ... ... ...] 5
```
The version in which the second matrix is a column-compressed matrix is named `PutCCMatAndMatT`.


## Linear solvers

`LinSol` defines an interface for linear solvers in `la`. Two implementations satisfying this
interface are:
1. `LinSolUmfpack` wrapper to Umfpack; and
2. `LinSolMumps` wrapper to MUMPS

There are also two _high level_ functions to solve linear systems with Umfpack:
1. `SolveRealLinSys`; and
2. `SolveComplexLinSys`
that are very convenient, for example:
```go
// sparse matrix
var A Triplet
A.Init(5, 5, 13); A.Put(0, 0, 1.0); A.Put(0, 0, 1.0); A.Put(1, 0, 3.0);
A.Put(0, 1, 3.0); A.Put(2, 1, -1.0); A.Put(4, 1, 4.0); A.Put(1, 2, 4.0);
A.Put(2, 2, -3.0); A.Put(3, 2, 1.0); A.Put(4, 2, 2.0); A.Put(2, 3, 2.0);
A.Put(1, 4, 6.0); A.Put(4, 4, 1.0) 

// right-hand-side
b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}

// solve
x, _ := SolveRealLinSys(&A, b)
```

Note however that the functions `SolveRealLinSys` and `SolveComplexLinSys` shouldn't be used for
repeated executions because memory would be constantly allocated and deallocated. In this case, it
is better to call the `LinSol` methods directly.



## Examples


### Eigenvalues of (3 x 3) matrix

See: <a href="t_jacobi_test.go">t_jacobi_test.go</a>

```go
A := [][]float64{
    {1, 2, 3},
    {2, 3, 2},
    {3, 2, 2},
}
Q := la.MatAlloc(3, 3)
v := make([]float64, 3)
nit, err := la.Jacobi(Q, v, A)
if err != nil {
    chk.Panic("Jacobi failed:\n%v", err)
}
```


### Cholesky decomposition

See: <a href="t_densesol_test.go">t_densesol_test.go</a>

```go
a := [][]float64{
    {25.0, 15.0, -5.0},
    {15.0, 18.0, 0.0},
    {-5.0, 0.0, 11.0},
}

L := la.MatAlloc(3, 3)
la.Cholesky(L, a) // L is such that: A = L * transp(L)
la.PrintMat("a", a, "%6g", false)
la.PrintMat("L", L, "%6g", false)
chk.Matrix(tst, "L", 1e-17, L, [][]float64{
    {5, 0, -0},
    {3, 3, 0},
    {-1, 1, 3},
})
```


### Linear Solver with real numbers

See: <a href="t_linsol_test.go">t_linsol_test.go</a>

```go
// input matrix data into Triplet
var t la.Triplet
t.Init(5, 5, 13)
t.Put(0, 0, 1.0)
t.Put(0, 0, 1.0)
t.Put(1, 0, 3.0)
t.Put(0, 1, 3.0)
t.Put(2, 1, -1.0)
t.Put(4, 1, 4.0)
t.Put(1, 2, 4.0)
t.Put(2, 2, -3.0)
t.Put(3, 2, 1.0)
t.Put(4, 2, 2.0)
t.Put(2, 3, 2.0)
t.Put(1, 4, 6.0)
t.Put(4, 4, 1.0)

// right-hand-side and solution
b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
x_correct := []float64{1, 2, 3, 4, 5}

// info
symmetric := false
timing := false

// allocate solver
lis := la.GetSolver("umfpack")
defer lis.Free()

// initialise solver
err := lis.InitR(t, symmetric, verbose, timing)
if err != nil {
    chk.Panic("%v", err.Error())
}

// factorise
err = lis.Fact()
if err != nil {
    chk.Panic("%v", err.Error())
}

// solve
var dummy bool
x := make([]float64, len(b))
err = lis.SolveR(x, b, dummy) // x := inv(A) * b
if err != nil {
    chk.Panic("%v", err.Error())
}

// output
A := t.ToMatrix(nil)
io.Pforan("A.x = b\n")
la.PrintMat("A", A.ToDense(), "%5g", false)
la.PrintVec("x", x, "%g ", false)
la.PrintVec("b", b, "%g ", false)

// check
chk.Vector(tst, "x", tol_cmp, x, x_correct)
la.CheckResidR(tst, tol_res, A.ToDense(), x, b)
```


### Linear Solver with complex numbers

See: <a href="t_linsol_test.go">t_linsol_test.go</a>

```go
// given the following matrix of complex numbers:
//      _                                                  _
//     |  19.73    12.11-i      5i        0          0      |
//     |  -0.51i   32.3+7i    23.07       i          0      |
// A = |    0      -0.51i    70+7.3i     3.95    19+31.83i  |
//     |    0        0        1+1.1i    50.17      45.51    |
//     |_   0        0          0      -9.351i       55    _|
//
// and the following vector:
//      _                  _
//     |    77.38+8.82i     |
//     |   157.48+19.8i     |
// b = |  1175.62+20.69i    |
//     |   912.12-801.75i   |
//     |_     550-1060.4i  _|
//
// solve:
//         A.x = b
//
// the solution is:
//      _            _
//     |     3.3-i    |
//     |    1+0.17i   |
// x = |      5.5     |
//     |       9      |
//     |_  10-17.75i _|

// flag indicating to store (real,complex) values in monolithic form => 1D array
xzmono := false

// input matrix in Complex Triplet format
var t la.TripletC
t.Init(5, 5, 16, xzmono) // 5 x 5 matrix with 16 non-zeros

// first column
t.Put(0, 0, 19.73, 0) // i=0, j=0, real=19.73, complex=0
t.Put(1, 0, 0, -0.51) // i=1, j=0, real=0, complex=-0.51

// second column
t.Put(0, 1, 12.11, -1) // i=0, j=1, real=12.11, complex=-1
t.Put(1, 1, 32.3, 7)
t.Put(2, 1, 0, -0.51)

// third column
t.Put(0, 2, 0, 5)
t.Put(1, 2, 23.07, 0)
t.Put(2, 2, 70, 7.3)
t.Put(3, 2, 1, 1.1)

// fourth column
t.Put(1, 3, 0, 1)
t.Put(2, 3, 3.95, 0)
t.Put(3, 3, 50.17, 0)
t.Put(4, 3, 0, -9.351)

// fifth column
t.Put(2, 4, 19, 31.83)
t.Put(3, 4, 45.51, 0)
t.Put(4, 4, 55, 0)

// right-hand-side
b := []complex128{
    77.38 + 8.82i,
    157.48 + 19.8i,
    1175.62 + 20.69i,
    912.12 - 801.75i,
    550 - 1060.4i,
}

// solution
x_correct := []complex128{
    3.3 - 1i,
    1 + 0.17i,
    5.5,
    9,
    10 - 17.75i,
}

// info
symmetric := false
timing := false

// allocate solver
lis := la.GetSolver("umfpack")
defer lis.Free()

// initialise solver
err := lis.InitC(t, symmetric, verbose, timing)
if err != nil {
    chk.Panic("%v", err.Error())
}

// factorise
err = lis.Fact()
if err != nil {
    chk.Panic("%v", err.Error())
}

// solve
var dummy bool
bR, bC := la.ComplexToRC(b)
xR := make([]float64, len(b))
xC := make([]float64, len(b))
err = lis.SolveC(xR, xC, bR, bC, dummy) // x := inv(A) * b
if err != nil {
    chk.Panic("%v", err.Error())
}
x := la.RCtoComplex(xR, xC)

// output
A := t.ToMatrix(nil)
io.Pforan("A.x = b\n")
la.PrintMatC("A", A.ToDense(), "(%g", "%+gi) ", false)
la.PrintVecC("x", x, "(%g", "%+gi) ", false)
la.PrintVecC("b", b, "(%g", "%+gi) ", false)

// check
chk.VectorC(tst, "x", tol_cmp, x, x_correct)
la.CheckResidC(tst, tol_res, A.ToDense(), x, b)
```


### Parallel solver with distributed right-hand-side

See: <a href="t_mumpssol01b_main.go">t_mumpssol01b_main.go</a>

```go
mpi.Start(false)
defer mpi.Stop(false)

myrank := mpi.Rank()
if myrank == 0 {
    chk.PrintTitle("Test MUMPS Sol 01b")
}

var t la.Triplet
b := []float64{8.0, 45.0, -3.0, 3.0, 19.0}
switch mpi.Size() {
case 1:
    t.Init(5, 5, 13)
    t.Put(0, 0, 1.0)
    t.Put(0, 0, 1.0)
    t.Put(1, 0, 3.0)
    t.Put(0, 1, 3.0)
    t.Put(2, 1, -1.0)
    t.Put(4, 1, 4.0)
    t.Put(1, 2, 4.0)
    t.Put(2, 2, -3.0)
    t.Put(3, 2, 1.0)
    t.Put(4, 2, 2.0)
    t.Put(2, 3, 2.0)
    t.Put(1, 4, 6.0)
    t.Put(4, 4, 1.0)
case 2:
    la.VecFill(b, 0)
    if myrank == 0 {
        t.Init(5, 5, 8)
        t.Put(0, 0, 1.0)
        t.Put(0, 0, 1.0)
        t.Put(1, 0, 3.0)
        t.Put(0, 1, 3.0)
        t.Put(2, 1, -1.0)
        t.Put(4, 1, 1.0)
        t.Put(4, 1, 1.5)
        t.Put(4, 1, 1.5)
        b[0] = 8.0
        b[1] = 40.0
        b[2] = 1.5
    } else {
        t.Init(5, 5, 8)
        t.Put(1, 2, 4.0)
        t.Put(2, 2, -3.0)
        t.Put(3, 2, 1.0)
        t.Put(4, 2, 2.0)
        t.Put(2, 3, 2.0)
        t.Put(1, 4, 6.0)
        t.Put(4, 4, 0.5)
        t.Put(4, 4, 0.5)
        b[1] = 5.0
        b[2] = -4.5
        b[3] = 3.0
        b[4] = 19.0
    }
default:
    chk.Panic("this test needs 1 or 2 procs")
}

x_correct := []float64{1, 2, 3, 4, 5}
sum_b_to_root := true
la.RunMumpsTestR(&t, 1e-14, b, x_correct, sum_b_to_root)
```


### Vectors

See: <a href="t_vector_test.go">t_vector_test.go</a>

```go
v := make([]float64, 5)
la.VecFill(v, 666)

vc := make([]complex128, 5)
la.VecFillC(vc, 666+666i)

va := []float64{1, 2, 3, 4, 5, 6}
vb := la.VecClone(va)

sum := la.VecAccum(v)

nrm := la.VecNorm(v)

u := []float64{333, 333, 333, 333, 333}
nrm = la.VecNormDiff(u, v)

u = []float64{0.1, 0.2, 0.3, 0.4, 0.5}
udotv := la.VecDot(u, v)

a := make([]float64, len(u))
la.VecCopy(a, 1, u)

b := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
la.VecAdd(b, 10, b) // b += 10.0*b

c := make([]float64, len(a))
la.VecAdd2(c, 1, a, 10, b) // c = 1.0*a + 10.0*b

mina := la.VecMin(a)

maxa := la.VecMax(a)

min2a, max2a := la.VecMinMax(a)

bdiv11 := []float64{b[0] / 11.0, b[1] / 11.0, b[2] / 11.0, b[3] / 11.0, b[4] / 11.0}
maxbdiv11 := la.VecLargest(b, 11)

amb1 := []float64{a[0] - b[0], a[1] - b[1], a[2] - b[2], a[3] - b[3], a[4] - b[4]}
amb2 := make([]float64, len(a))
la.VecAdd2(amb2, 1, a, -1, b)
maxdiffab := la.VecMaxDiff(a, b)

az := []complex128{complex(a[0], 1), complex(a[1], 3), complex(a[2], 0.5), complex(a[3], 1), complex(a[4], 0)}
bz := []complex128{complex(b[0], 1), complex(b[1], 6), complex(b[2], 0.8), complex(b[3], -3), complex(b[4], 1)}
ambz := []complex128{az[0] - bz[0], az[1] - bz[1], az[2] - bz[2], az[3] - bz[3], az[4] - bz[4]}
maxdiffabz := la.VecMaxDiffC(az, bz)

scal1 := make([]float64, len(a))
la.VecScale(scal1, 0.5, 0.1, amb1)

scal2 := make([]float64, len(a))
la.VecScaleAbs(scal2, 0.5, 0.1, amb1)

rms := la.VecRms(v)

rmserr := la.VecRmsErr(v, 0, 1, v)

w := []float64{333, 333, 333, 333, 333}
rmserr = la.VecRmsError(v, w, 0, 1, v)
```


### Matrices

See: <a href="t_matrix_test.go">t_matrix_test.go</a>

```go
// MatAlloc
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

aclone := la.MatClone(a)

// MatFill
b := la.MatAlloc(5, 3)
la.MatFill(b, 2)

// MatScale
c := la.MatAlloc(5, 3)
la.MatFill(c, 2)
la.MatScale(c, 1.0/4.0)

// MatCopy
d := la.MatAlloc(3, 5)
la.MatCopy(d, 1, a)

// MatSetDiag
e := la.MatAlloc(3, 3)
la.MatSetDiag(e, 1)

// MatMaxDiff
f := [][]float64{
    {1.1, 2.2, 3.3, 4.4, 5.5},
    {0.1, 0.2, 0.3, 0.4, 0.5},
    {1, 2, 3, 4, 5},
}
fma := [][]float64{
    {f[0][0] - a[0][0], f[0][1] - a[0][1], f[0][2] - a[0][2], f[0][3] - a[0][3], f[0][4] - a[0][4]},
    {f[1][0] - a[1][0], f[1][1] - a[1][1], f[1][2] - a[1][2], f[1][3] - a[1][3], f[1][4] - a[1][4]},
    {f[2][0] - a[2][0], f[2][1] - a[2][1], f[2][2] - a[2][2], f[2][3] - a[2][3], f[2][4] - a[2][4]},
}
maxdiff := la.MatMaxDiff(f, a)

// MatLargest
largest := la.MatLargest(a, 1)

// MatGetCol
col := la.MatGetCol(0, a)

// MatNormF
Pll = true
NCPU = 3
A := [][]float64{
    {-3, 5, 7},
    {2, 6, 4},
    {0, 2, 8},
}
normFA := la.MatNormF(A)

// MatNormI
normIA := la.MatNormI(A)

// MatInv
g := [][]float64{
    {1, 2, 3},
    {0, 4, 5},
    {1, 0, 6},
}
gi := la.MatAlloc(3, 3)
detg, err := la.MatInv(gi, g, 1e-17)
if err != nil {
    chk.Panic("%v", err.Error())
}
gi22 := [][]float64{
    {gi[0][0] * 22.0, gi[0][1] * 22.0, gi[0][2] * 22.0},
    {gi[1][0] * 22.0, gi[1][1] * 22.0, gi[1][2] * 22.0},
    {gi[2][0] * 22.0, gi[2][1] * 22.0, gi[2][2] * 22.0},
}
```


### Matrices and Vectors

See: <a href="t_matvec_test.go">t_matvec_test.go</a>

```go
a := [][]float64{
    {1, 2, 3, 4, 5},
    {0.1, 0.2, 0.3, 0.4, 0.5},
    {10, 20, 30, 40, 50},
}
b := [][]float64{
    {1.1, 2.2, 3.3, 4.4, 5.5},
    {0.1, 0.2, 0.3, 0.4, 0.5},
    {1, 2, 3, 4, 5},
}
c := [][]float64{
    {1, 2, 3},
    {0, 4, 5},
    {1, 0, 6},
}
d := [][]float64{
    {0.1, 10, 0.1},
    {0.2, 20, 0.2},
    {0.3, 30, 0.3},
    {0.4, 40, 0.4},
    {0.5, 50, 0.5},
}
e := [][]float64{
    {0.01, 0.02, 1},
    {-0.02, 0.03, 0.01},
    {0.5, -0.01, 0.02},
}
u := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
v := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
w := []float64{10.0, 20.0, 30.0}
r := []float64{1000, 1000, 1000, 1000, 1000}
s := []float64{1000, 1000, 1000}

p := make([]float64, len(a))
la.MatVecMul(p, 1, a, u) // p := 1*a*u

la.MatVecMulAdd(s, 1, a, u) // s += dot(a, u)

q := make([]float64, 5)
la.MatTrVecMul(q, 1, a, w) // q = dot(transpose(a), w)

la.MatTrVecMulAdd(r, 1, a, w) // r += dot(transpose(a), w)

la.MatVecMulCopyAdd(p, 0.5, w, 2.0, a, u) // p := 0.5*w + 2*a*u

f := la.MatAlloc(3, 3)
la.MatMul(f, 1, a, d) // f = dot(a, d)

g := la.MatAlloc(3, 3)
la.MatMul3(g, 1, c, e, f) // g = dot(c, dot(e, f))

h := la.MatAlloc(5, 3)
la.MatTrMul3(h, 1, a, e, f) // h = dot(transpose(a), dot(e, f))

m := la.MatAlloc(5, 3)
la.MatFill(m, 10000)
n := la.MatAlloc(5, 3)
la.MatCopy(n, 1, m)            // n := m
la.MatTrMulAdd3(n, 1, a, e, f) // n += dot(transpose(a), dot(e, f))

udyv := la.MatAlloc(5, 5)
la.MatFill(udyv, 1000)
la.VecOuterAdd(udyv, 0.5, u, v)
```


### Matrix-Vector multiplication

See: <a href="t_matvecmul_test.go">t_matvecmul_test.go</a>

```go
a := [][]float64{
    {10.0, 20.0, 30.0, 40.0, 50.0},
    {1.0, 20.0, 3.0, 40.0, 5.0},
    {10.0, 2.0, 30.0, 4.0, 50.0},
}
u := []float64{0.5, 0.4, 0.3, 0.2, 0.1}
r := []float64{0.5, 0.4, 0.3}
z := []float64{1000, 1000, 1000}
w := []float64{1000, 2000, 3000, 4000, 5000}
au := make([]float64, 3)
atr := make([]float64, 5)
au_cor := []float64{35, 17.9, 20.6}
atr_cor := []float64{8.4, 18.6, 25.2, 37.2, 42}
zpau_cor := []float64{1035, 1017.9, 1020.6}
wpar_cor := []float64{1008.4, 2018.6, 3025.2, 4037.2, 5042}
la.MatVecMul(au, 1, a, u)     // au  = 1*a*u
la.MatTrVecMul(atr, 1, a, r)  // atr = 1*transp(a)*r
la.MatVecMulAdd(z, 1, a, u)   // z  += 1*a*u
la.MatTrVecMulAdd(w, 1, a, r) // w  += 1*transp(a)*r
chk.Vector(tst, "au", 1.0e-17, au, au_cor)
chk.Vector(tst, "atr", 1.0e-17, atr, atr_cor)
chk.Vector(tst, "zpau", 1.0e-12, z, zpau_cor)
chk.Vector(tst, "wpar", 1.0e-12, w, wpar_cor)
```
