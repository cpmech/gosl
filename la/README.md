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
```
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

**TODO**
