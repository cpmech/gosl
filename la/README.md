# Gosl. la. Linear Algebra and efficient sparse solvers

[![GoDoc](https://godoc.org/github.com/cpmech/gosl/la?status.svg)](https://godoc.org/github.com/cpmech/gosl/la) 

More information is available in **[the documentation of this package](https://godoc.org/github.com/cpmech/gosl/la).**

The `la` subpackage implements functions to perform linear algebra operations such as vector
addition or matrix-vector multiplications. It also wraps some efficient sparse linear systems
solvers such as [Umfpack](http://faculty.cse.tamu.edu/davis/suitesparse.html) and
[MUMPS](http://mumps.enseeiht.fr) (**not** the _parotitis_ disease!).

Both [Umfpack](http://faculty.cse.tamu.edu/davis/suitesparse.html) and
[MUMPS](http://mumps.enseeiht.fr) solvers are very efficient!

The other subpackages [la/oblas](https://github.com/cpmech/gosl/tree/master/la/oblas) and
[la/mkl](https://github.com/cpmech/gosl/tree/master/la/mkl) are sometimes called by `la` to improve
performance.


## Structures for sparse problems

In `la`, there are two types of structures to hold data for solving a sparse linear system:
_triplet_ and _column-compressed matrix_. Both have a complex version and are named as follows:
1. `Triplet` and `TripletC` (complex)
2. `CCMatrix` and `CCMatrixC` (complex)

The `Triplet` is more convenient for data input; e.g. during the assemblage step in a finite element
code. The `CCMatrix` is better (faster) for computations and is the structure used by Umfpack.

Triplets are initialised by giving the size of the corresponding matrix and the number of non-zero
entries (components). Thus, only space for the non-zero entries is allocated. Afterwards, each
component can be set by calling the `Put` method after calling the `Start` method. The `Start`
method can be called several times. Therefore, this structure helps with the efficient
implementation of codes requiring multiple solutions of linear systems as in finite element
analyses.

To convert the Triplet into a Column-Compressed Matrix, call the `ToMatrix` method of `Triplet`. And
to convert the sparse matrix to a dense version (e.g. for reporting/printing), call the `ToDense`
method of `CCMatrix`.

The `Triplet` has also a very convenient method to copy the contents of another Triplet (i.e. sparse
matrix) into the positions starting with the maximum number of columns of this second matrix and the
positions starting with the maximum number of rows of this second matrix. This is done with the
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


## Linear solvers for sparse problems

`SparseSolver` defines an interface for linear solvers in `la`. Two implementations satisfying this
interface are:
1. `Umfpack` wrapper to Umfpack; and
2. `Mumps` wrapper to MUMPS

There are also _high level_ functions to solve linear systems with Umfpack:
1. `SpSolve`; and
2. `SpSolveC` with complex numbers

Note however that the high level functions shouldn't be used for repeated executions because memory
would be constantly allocated and deallocated.



## Examples


### Eigenvalues of (3 x 3) matrix

See: <a href="t_jacobi_test.go">t_jacobi_test.go</a>


### Cholesky decomposition

See: <a href="t_densesol_test.go">t_densesol_test.go</a>
