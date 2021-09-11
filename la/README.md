# Gosl. la. Linear Algebra: vector, matrix, efficient sparse solvers, eigenvalues, decompositions, etc.

[![Go Reference](https://pkg.go.dev/badge/github.com/cpmech/gosl/la.svg)](https://pkg.go.dev/github.com/cpmech/gosl/la)

The `la` subpackage implements functions to perform linear algebra operations such as vector
addition or matrix-vector multiplications. It also wraps some efficient sparse linear systems
solvers such as [Umfpack](http://faculty.cse.tamu.edu/davis/suitesparse.html) and
[MUMPS](http://mumps.enseeiht.fr) (**not** the _parotitis_ disease!).

Both [Umfpack](http://faculty.cse.tamu.edu/davis/suitesparse.html) and
[MUMPS](http://mumps.enseeiht.fr) solvers are very efficient!

Sometimes, we call the _lower level_ functions in [la/oblas](https://github.com/cpmech/gosl/tree/master/la/oblas)
to improve performance.

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

### Vectors and matrices

* <a href="t_vector_test.go">source file</a> Test Vector
* <a href="t_matrix_test.go">source file</a> Test Matrix
* <a href="t_matrix_ops_test.go">source file</a> Test Matrix operations

### BLAS1, 2 and 3 functions

* <a href="t_blas1_test.go">source file</a> Test BLAS1 routines
* <a href="t_blas2_test.go">source file</a> Test BLAS2 routines
* <a href="t_blas3_test.go">source file</a> Test BLAS3 routines

### General dense solver and Cholesky decomposition

* <a href="t_densesol_test.go">source file</a> Test Dense Solver

### Eigenvalues and eigenvectors of general matrix

* <a href="t_eigen_test.go">source file</a> Test Eigenvalues/Eigenvectors

### Eigenvalues of symmetric (3 x 3) matrix

* <a href="t_jacobi_test.go">source file</a> Test Jacobi iteration

### Sparse BLAS functions

* <a href="t_sp_blas_test.go">source file</a> Test sparse BLAS

### Conversions related to sparse matrices

* <a href="t_sp_conversions_test.go">source file</a> Test conversions

### Sparse Triplet and Matrix

* <a href="t_sp_matrix_test.go">source file</a> Test sparse Triplet and Matrix

### Sparse linear solver using MUMPS

* <a href="t_sp_solver_mumps_test.go">source file</a> Test sparse solver MUMPS

### Sparse linear solver using UMFPACK

* <a href="t_sp_solver_umfpack_test.go">source file</a> Test sparse solver UMFPACK

### Solutions using sparse solvers

* <a href="t_sp_solver_test.go">source file</a> Test solutions of sparse linear systems

## API

[Please see the documentation here](https://pkg.go.dev/github.com/cpmech/gosl/la)
