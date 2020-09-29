# Gosl. la. Linear Algebra: vector, matrix, efficient sparse solvers, eigenvalues, decompositions, etc.

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

<a href="t_vector_test.go">source file</a>
<a href="t_matrix_test.go">source file</a>
<a href="t_matrix_ops_test.go">source file</a>

### BLAS1, 2 and 3 functions

<a href="t_blas1_test.go">source file</a>
<a href="t_blas2_test.go">source file</a>
<a href="t_blas3_test.go">source file</a>

### General dense solver and Cholesky decomposition

<a href="t_densesol_test.go">source file</a>

### Eigenvalues and eigenvectors of general matrix

<a href="t_eigen_test.go">source file</a>

### Eigenvalues of symmetric (3 x 3) matrix

<a href="t_jacobi_test.go">source file</a>

### Sparse BLAS functions

<a href="t_sp_blas_test.go">source file</a>

### Conversions related to sparse matrices

<a href="t_sp_conversions_test.go">source file</a>

### Sparse Triplet and Matrix

<a href="t_sp_matrix_test.go">source file</a>

### Sparse linear solver using MUMPS

<a href="t_sp_solver_mumps_test.go">source file</a>

### Sparse linear solver using UMFPACK

<a href="t_sp_solver_umfpack_test.go">source file</a>

### Solutions using sparse solvers

<a href="t_sp_solver_test.go">source file</a>

## API

**go doc**

```
package la // import "gosl/la"

Package la implements functions and structure for Linear Algebra
computations. It defines a Vector and Matrix types for computations with
dense data and also a Triplet and CCMatrix for sparse data.

FUNCTIONS

func CheckEigenVecL(tst *testing.T, A *Matrix, λ VectorC, u *MatrixC, tol float64)
    CheckEigenVecL checks left eigenvector:

         H                  H
        u [j] ⋅ A = λ[j] ⋅ u [j]    LEFT eigenvectors

func CheckEigenVecR(tst *testing.T, A *Matrix, λ VectorC, v *MatrixC, tol float64)
    CheckEigenVecR checks right eigenvector:

        A ⋅ v[j] = λ[j] ⋅ v[j]      RIGHT eigenvectors

func Cholesky(L, a *Matrix)
    Cholesky returns the Cholesky decomposition of a symmetric positive-definite
    matrix

        a = L * trans(L)

func DenSolve(x Vector, A *Matrix, b Vector, preserveA bool)
    DenSolve solves dense linear system using LAPACK (OpenBLAS)

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

func EigenVal(w VectorC, A *Matrix, preserveA bool)
    EigenVal computes eigenvalues of general matrix

        A ⋅ v[j] = λ[j] ⋅ v[j]

        INPUT:
          a -- general matrix

        OUTPUT:
          w -- eigenvalues [pre-allocated]

func EigenVecL(u *MatrixC, w VectorC, A *Matrix, preserveA bool)
    EigenVecL computes eigenvalues and LEFT eigenvectors of general matrix

         H                  H
        u [j] ⋅ A = λ[j] ⋅ u [j]    LEFT eigenvectors

        INPUT:
          a -- general matrix

        OUTPUT:
          u -- matrix with the eigenvectors; each column contains one eigenvector [pre-allocated]
          w -- eigenvalues [pre-allocated]

func EigenVecLR(u, v *MatrixC, w VectorC, A *Matrix, preserveA bool)
    EigenVecLR computes eigenvalues and LEFT and RIGHT eigenvectors of general
    matrix

        A ⋅ v[j] = λ[j] ⋅ v[j]      RIGHT eigenvectors

         H                  H
        u [j] ⋅ A = λ[j] ⋅ u [j]    LEFT eigenvectors

        INPUT:
          a -- general matrix

        OUTPUT:
          u -- matrix with the LEFT eigenvectors; each column contains one eigenvector [pre-allocated]
          v -- matrix with the RIGHT eigenvectors; each column contains one eigenvector [pre-allocated]
          w -- λ eigenvalues [pre-allocated]

func EigenVecR(v *MatrixC, w VectorC, A *Matrix, preserveA bool)
    EigenVecR computes eigenvalues and RIGHT eigenvectors of general matrix

        A ⋅ v[j] = λ[j] ⋅ v[j]

        INPUT:
          a -- general matrix

        OUTPUT:
          v -- matrix with the eigenvectors; each column contains one eigenvector [pre-allocated]
          w -- eigenvalues [pre-allocated]

func Jacobi(Q *Matrix, v Vector, A *Matrix)
    Jacobi performs the Jacobi transformation of a symmetric matrix to find its
    eigenvectors and eigenvalues.

    The Jacobi method consists of a sequence of orthogonal similarity
    transformations. Each transformation (a Jacobi rotation) is just a plane
    rotation designed to annihilate one of the off-diagonal matrix elements.
    Successive transformations undo previously set zeros, but the off-diagonal
    elements nevertheless get smaller and smaller. Accumulating the product of
    the transformations as you go gives the matrix of eigenvectors (Q), while
    the elements of the final diagonal matrix (A) are the eigenvalues.

    The Jacobi method is absolutely foolproof for all real symmetric matrices.

              A = Q ⋅ L ⋅ Qᵀ

        Input:
         A -- matrix to compute eigenvalues (SYMMETRIC and SQUARE)
        Output:
         A -- modified
         Q -- matrix which columns are the eigenvectors
         v -- vector with the eigenvalues

        NOTE: for matrices of order greater than about 10, say, the algorithm is slower,
              by a significant constant factor, than the QR method.

func MatAdd(res *Matrix, α float64, a *Matrix, β float64, b *Matrix)
    MatAdd adds the scaled components of two matrices

        res := α⋅a + β⋅b   ⇒   result[i][j] := α⋅a[i][j] + β⋅b[i][j]

func MatCondNum(a *Matrix, normtype string) (res float64)
    MatCondNum returns the condition number of a square matrix using the inverse
    of this matrix; thus it is not as efficient as it could be, e.g. by using
    the SV decomposition.

        normtype -- Type of norm to use:
          "F" or "" => Frobenius
          "I"       => Infinite

func MatInv(ai, a *Matrix, calcDet bool) (det float64)
    MatInv computes the inverse of a general matrix (square or not). It also
    computes the pseudo-inverse if the matrix is not square.

        Input:
          a -- input matrix (M x N)
        Output:
          ai -- inverse matrix (N x M)
          det -- determinant of matrix (ONLY if calcDet == true and the matrix is square)
        NOTE: the dimension of the ai matrix must be N x M for the pseudo-inverse

func MatInvSmall(ai, a *Matrix, tol float64) (det float64)
    MatInvSmall computes the inverse of small matrices of size 1x1, 2x2, or 3x3.
    It also returns the determinant.

        Input:
          a   -- the matrix
          tol -- tolerance to assume zero determinant
        Output:
          ai  -- the inverse matrix
          det -- determinant of a

func MatMatMul(c *Matrix, α float64, a, b *Matrix)
    MatMatMul returns the matrix multiplication (scaled)

        c := α⋅a⋅b    ⇒    cij := α * aik * bkj

func MatMatMulAdd(c *Matrix, α float64, a, b *Matrix)
    MatMatMulAdd returns the matrix multiplication (scaled)

        c += α⋅a⋅b    ⇒    cij += α * aik * bkj

func MatMatTrMul(c *Matrix, α float64, a, b *Matrix)
    MatMatTrMul returns the matrix multiplication (scaled) with transposed(b)

        c := α⋅a⋅bᵀ    ⇒    cij := α * aik * bjk

func MatMatTrMulAdd(c *Matrix, α float64, a, b *Matrix)
    MatMatTrMulAdd returns the matrix multiplication (scaled) with transposed(b)

        c += α⋅a⋅bᵀ    ⇒    cij += α * aik * bjk

func MatSvd(s []float64, u, vt, a *Matrix, copyA bool)
    MatSvd performs the SVD decomposition

        Input:
          a     -- matrix a
          copyA -- creates a copy of a; otherwise 'a' is modified
        Output:
          s  -- diagonal terms [must be pre-allocated] len(s) = imin(a.M, a.N)
          u  -- left matrix [must be pre-allocated] u is (a.M x a.M)
          vt -- transposed right matrix [must be pre-allocated] vt is (a.N x a.N)

func MatTrMatMul(c *Matrix, α float64, a, b *Matrix)
    MatTrMatMul returns the matrix multiplication (scaled) with transposed(a)

        c := α⋅aᵀ⋅b    ⇒    cij := α * aki * bkj

func MatTrMatMulAdd(c *Matrix, α float64, a, b *Matrix)
    MatTrMatMulAdd returns the matrix multiplication (scaled) with transposed(a)

        c += α⋅aᵀ⋅b    ⇒    cij += α * aki * bkj

func MatTrMatTrMul(c *Matrix, α float64, a, b *Matrix)
    MatTrMatTrMul returns the matrix multiplication (scaled) with transposed(a)
    and transposed(b)

        c := α⋅aᵀ⋅bᵀ    ⇒    cij := α * aki * bjk

func MatTrMatTrMulAdd(c *Matrix, α float64, a, b *Matrix)
    MatTrMatTrMulAdd returns the matrix multiplication (scaled) with
    transposed(a) and transposed(b)

        c += α⋅aᵀ⋅bᵀ    ⇒    cij += α * aki * bjk

func MatTrVecMul(v Vector, α float64, a *Matrix, u Vector)
    MatTrVecMul returns the transpose(matrix)-vector multiplication

        v = α⋅aᵀ⋅u    ⇒    vi = α * aji * uj = α * uj * aji

func MatVecMul(v Vector, α float64, a *Matrix, u Vector)
    MatVecMul returns the matrix-vector multiplication

        v = α⋅a⋅u    ⇒    vi = α * aij * uj

func MatVecMulAdd(v Vector, α float64, a *Matrix, u Vector)
    MatVecMulAdd returns the matrix-vector multiplication with addition

        v += α⋅a⋅u    ⇒    vi += α * aij * uj

func MatVecMulAddC(v VectorC, α complex128, a *MatrixC, u VectorC)
    MatVecMulAddC returns the matrix-vector multiplication with addition
    (complex version)

        v += α⋅a⋅u    ⇒    vi += α * aij * uj

func MatVecMulC(v VectorC, α complex128, a *MatrixC, u VectorC)
    MatVecMulC returns the matrix-vector multiplication (complex version)

        v = α⋅a⋅u    ⇒    vi = α * aij * uj

func SolveRealLinSysSPD(x Vector, a *Matrix, b Vector)
    SolveRealLinSysSPD solves a linear system with real numbres and a
    Symmetric-Positive-Definite (SPD) matrix

             x := inv(a) * b

        NOTE: this function uses Cholesky decomposition and should be used for small systems

func SolveTwoRealLinSysSPD(x, X Vector, a *Matrix, b, B Vector)
    SolveTwoRealLinSysSPD solves two linear systems with real numbres and
    Symmetric-Positive-Definite (SPD) matrices

             x := inv(a) * b    and    X := inv(a) * B

        NOTE: this function uses Cholesky decomposition and should be used for small systems

func SpCheckDiag(a *CCMatrix) bool
    SpCheckDiag checks if all elements on the diagonal of "a" are present.

        OUTPUT:
         ok -- true if all diagonal elements are present;
               otherwise, ok == false if any diagonal element is missing.

func SpInitRc(R *CCMatrix, C *CCMatrixC, J *CCMatrix)
    SpInitRc initialises two complex sparse matrices (residual correction)
    according to:

        Real:       R :=  γ      *I - J
        Complex:    C := (α + βi)*I - J
        NOTE: "J" must include all diagonal elements

func SpInitSimilar(b *CCMatrix, a *CCMatrix)
    SpInitSimilar initialises another matrix "b" with the same structure (Ap,
    Ai) of sparse matrix "a". The values Ax are not copied though.

func SpInitSimilarR2C(b *CCMatrixC, a *CCMatrix)
    SpInitSimilarR2C initialises another matrix "b" (complex) with the same
    structure (Ap, Ai) of sparse matrix "a" (real). The values Ax are not copied
    though (Bx and Bz are not set).

func SpMatAddI(r *CCMatrix, α float64, a *CCMatrix, β float64)
    SpMatAddI adds an identity matrix I to "a", scaled by α and β according to:

        r := α*a + β*I

func SpMatAddMat(c *CCMatrix, α float64, a *CCMatrix, β float64, b *CCMatrix, a2c, b2c []int)
    SpMatAddMat adds two sparse matrices. The 'c' matrix matrix and the 'a2c'
    and 'b2c' arrays must be pre-allocated by SpAllocMatAddMat. The result is:

        c := α*a + β*b
        NOTE: this routine does not check for the correct sizes, since this is expect to be
              done by SpAllocMatAddMat

func SpMatAddMatC(C *CCMatrixC, R *CCMatrix, α, β, γ float64, a *CCMatrix, μ float64, b *CCMatrix, a2c, b2c []int)
    SpMatAddMatC adds two real sparse matrices with two sets of coefficients in
    such a way that one real matrix (R) and another complex matrix (C) are
    obtained. The results are:

          R :=  γ    *a + μ*b
          C := (α+βi)*a + μ*b
        NOTE: the structure of R and C are the same and can be allocated with SpAllocMatAddMat,
              followed by one call to SpInitSimilarR2C. For example:
                  R, a2c, b2c := SpAllocMatAddMat(a, b)
                  SpInitSimilarR2C(C, R)

func SpMatMatTrMul(b *Matrix, α float64, a *CCMatrix)
    SpMatMatTrMul computes the multiplication of sparse matrix with itself
    transposed:

        b := α * a * aᵀ   b_iq = α * a_ij * a_qj

    / Note: b is symmetric

func SpMatTrVecMul(v Vector, α float64, a *CCMatrix, u Vector)
    SpMatTrVecMul returns the (sparse) matrix-vector multiplication with "a"
    transposed (scaled):

        v := α * transp(a) * u  =>  vj = α * aij * ui
        NOTE: dense vector v will be first initialised with zeros

func SpMatTrVecMulAdd(v Vector, α float64, a *CCMatrix, u Vector)
    SpMatTrVecMulAdd returns the (sparse) matrix-vector multiplication with
    addition and "a" transposed (scaled):

        v += α * transp(a) * u  =>  vj += α * aij * ui

func SpMatTrVecMulAddC(v VectorC, α complex128, a *CCMatrixC, u VectorC)
    SpMatTrVecMulAddC returns the (sparse/complex) matrix-vector multiplication
    with addition and "a" transposed (scaled):

        v += α * transp(a) * u  =>  vj += α * aij * ui

func SpMatTrVecMulC(v VectorC, α complex128, a *CCMatrixC, u VectorC)
    SpMatTrVecMulC returns the (sparse/complex) matrix-vector multiplication
    with "a" transposed (scaled):

        v := α * transp(a) * u  =>  vj = α * aij * ui
        NOTE: dense vector v will be first initialised with zeros

func SpMatVecMul(v Vector, α float64, a *CCMatrix, u Vector)
    SpMatVecMul returns the (sparse) matrix-vector multiplication (scaled):

        v := α * a * u  =>  vi = α * aij * uj
        NOTE: dense vector v will be first initialised with zeros

func SpMatVecMulAdd(v Vector, α float64, a *CCMatrix, u Vector)
    SpMatVecMulAdd returns the (sparse) matrix-vector multiplication with
    addition (scaled):

        v += α * a * u  =>  vi += α * aij * uj

func SpMatVecMulAddC(v VectorC, α complex128, a *CCMatrixC, u VectorC)
    SpMatVecMulAddC returns the (sparse/complex) matrix-vector multiplication
    with addition (scaled):

        v += α * a * u  =>  vi += α * aij * uj

func SpMatVecMulAddX(v Vector, a *CCMatrix, α float64, u Vector, β float64, w Vector)
    SpMatVecMulAddX returns the (sparse) matrix-vector multiplication with
    addition (scaled/extended):

        v += a * (α*u + β*w)  =>  vi += aij * (α*uj + β*wj)

func SpMatVecMulC(v VectorC, α complex128, a *CCMatrixC, u VectorC)
    SpMatVecMulC returns the (sparse/complex) matrix-vector multiplication
    (scaled):

        v := α * a * u  =>  vi = α * aij * uj
        NOTE: dense vector v will be first initialised with zeros

func SpSetRc(R *CCMatrix, C *CCMatrixC, α, β, γ float64, J *CCMatrix)
    SpSetRc sets the values within two complex sparse matrices (residual
    correction) according to:

        Real:       R :=  γ      *I - J
        Complex:    C := (α + βi)*I - J
        NOTE: "J" must include all diagonal elements

func SpTriAdd(c *Triplet, α float64, a *Triplet, β float64, b *Triplet)
    SpTriAdd adds two matrices in Triplet format:

        c := α*a + β*b
        NOTE: the output 'c' triplet must be able to hold all nonzeros of 'a' and 'b'
              actually the 'c' triplet is just expanded

func SpTriAddR2C(c *TripletC, α, β float64, a *Triplet, μ float64, b *Triplet)
    SpTriAddR2C adds two real matrices in Triplet format generating a complex
    triplet according to:

        c := (α+βi)*a + μ*b
        NOTE: the output 'c' triplet must be able to hold all nonzeros of 'a' and 'b'

func SpTriMatTrVecMul(y Vector, a *Triplet, x Vector)
    SpTriMatTrVecMul returns the matrix-vector multiplication with transposed
    matrix a in triplet format and two dense vectors x and y

        y := transpose(a) * x    or    y_I := a_JI * x_J    or     y_j := a_ij * x_i

func SpTriMatVecMul(y Vector, a *Triplet, x Vector)
    SpTriMatVecMul returns the matrix-vector multiplication with matrix a in
    triplet format and two dense vectors x and y

        y := a * x    or    y_i := a_ij * x_j

func SpTriReduce(comm *mpi.Communicator, J *Triplet)
    SpTriReduce joins (MPI) parallel triplets to root (Rank == 0) processor.

        NOTE: J in root is also joined into Jroot

func SpTriSetDiag(a *Triplet, n int, v float64)
    SpTriSetDiag sets a (n x n) real triplet with diagonal values 'v'

func TestSolverResidual(tst *testing.T, a *Matrix, x, b Vector, tolNorm float64)
    TestSolverResidual check the residual of a linear system solution

func TestSolverResidualC(tst *testing.T, a *MatrixC, x, b VectorC, tolNorm float64)
    TestSolverResidualC check the residual of a linear system solution (complex
    version)

func TestSpSolver(tst *testing.T, solverKind string, symmetric bool, t *Triplet, b, xCorrect Vector,
	tolX, tolRes float64, verbose, bIsDistr bool, comm *mpi.Communicator)
    TestSpSolver tests a sparse solver

func TestSpSolverC(tst *testing.T, solverKind string, symmetric bool, t *TripletC, b, xCorrect VectorC,
	tolX, tolRes float64, verbose, bIsDistr bool, comm *mpi.Communicator)
    TestSpSolverC tests a sparse solver (complex version)

func VecAdd(res Vector, α float64, u Vector, β float64, v Vector)
    VecAdd adds the scaled components of two vectors

        res := α⋅u + β⋅v   ⇒   result[i] := α⋅u[i] + β⋅v[i]

func VecDot(u, v Vector) (res float64)
    VecDot returns the dot product between two vectors:

        s := u・v

func VecMaxDiff(u, v Vector) (maxdiff float64)
    VecMaxDiff returns the maximum absolute difference between two vectors

        maxdiff = max(|u - v|)

func VecRmsError(u, v Vector, a, m float64, s Vector) (rms float64)
    VecRmsError returns the scaled root-mean-square of the difference between
    two vectors with components normalised by a scaling factor

                     __________________________
                    /     ————              2
                   /  1   \    /  error[i]  \
        rms =  \  /  ———  /    | —————————— |
                \/    N   ———— \  scale[i]  /

        error[i] = |u[i] - v[i]|

        scale[i] = a + m*|s[i]|

func VecScaleAbs(scale Vector, a, m float64, x Vector)
    VecScaleAbs creates a "scale" vector using the absolute value of another
    vector

        scale := a + m ⋅ |x|     ⇒      scale[i] := a + m ⋅ |x[i]|

func VecVecTrMul(a *Matrix, α float64, u, v Vector)
    VecVecTrMul returns the matrix = vector-transpose(vector) multiplication
    (e.g. dyadic product)

        a = α⋅u⋅vᵀ    ⇒    aij = α * ui * vj


TYPES

type CCMatrix struct {
	// Has unexported fields.
}
    CCMatrix represents a sparse matrix using the so-called "column-compressed
    format".

func SpAllocMatAddMat(a, b *CCMatrix) (c *CCMatrix, a2c, b2c []int)
    SpAllocMatAddMat allocates a matrix 'c' to hold the result of the addition
    of 'a' and 'b'. It also allocates the mapping arrays a2c and b2c, where:

        a2c maps 'k' in 'a' to 'k' in 'c': len(a2c) = a.nnz
        b2c maps 'k' in 'b' to 'k' in 'c': len(b2c) = b.nnz

func (o *CCMatrix) Set(m, n int, Ap, Ai []int, Ax []float64)
    Set sets column-compressed matrix directly

func (o *CCMatrix) ToDense() (res *Matrix)
    ToDense converts a column-compressed matrix to dense form

func (o *CCMatrix) WriteSmat(dirout, fnkey string, tol float64, oneBasedIndices ...bool)
    WriteSmat writes a SMAT file (that can be visualised with vismatrix) or a
    MatrixMarket file

        About the .smat file:
         - lines starting with the percent or hashtag mark (% or #) are ignored (they are comments)
         - the first line contains the number of rows, columns, and non-zero entries
         - the following lines contain the indices of row, column, and the non-zero entry
         - indices are 0-based, unless oneBasedIndices == true
         - with one-based indices, smat is similar to MatrixMarket without the header

        Example of .smat file:

          # this is a comment
          % this is also a comment
          m n nnz
           i j x
            ...
           i j x

        NOTE: CCMatrix must be used to generate the resulting values because
              duplicates must be added before saving file

        dirout -- directory for output. will be created
        fnkey  -- filename key (filename without extension). ".smat" will be added
        tol    -- tolerance to skip zero values
        oneBasedIndices -- if true, 1 is added to all indices (like MatrixMarket file format)

type CCMatrixC struct {
	// Has unexported fields.
}
    CCMatrixC represents a sparse matrix using the so-called "column-compressed
    format". (complex version)

func (o *CCMatrixC) ToDense() (res *MatrixC)
    ToDense converts a column-compressed matrix (complex) to dense form

func (o *CCMatrixC) WriteSmat(dirout, fnkey string, tol float64)
    WriteSmat writes a ".smat" file that can be visualised with vismatrix

        NOTE: CCMatrix must be used to generate the resulting values because
              duplicates must be added before saving file

        dirout -- directory for output. will be created
        fnkey  -- filename key (filename without extension). ".smat" will be added
        tol    -- tolerance to skip zero values

func (o *CCMatrixC) WriteSmatAbs(dirout, fnkey string, tol float64)
    WriteSmatAbs writes a ".smat" file that can be visualised with vismatrix
    (abs(complex) version)

        NOTE: CCMatrix must be used to generate the resulting values because
              duplicates must be added before saving file

        tol -- tolerance to skip zero values

type Equations struct {

	// essential
	N    int   // total number of equations
	Nu   int   // number of unknowns (size of u-system)
	Nk   int   // number of known values (size of k-system)
	UtoF []int // reduced u-system to full-system
	FtoU []int // full-system to reduced u-system
	KtoF []int // reduced k-system to full system
	FtoK []int // full-system to reduced k-system

	// convenience
	Auu, Auk, Aku, Akk *Triplet // the partitioned system in sparse format
	Duu, Duk, Dku, Dkk *Matrix  // the partitioned system in dense format
	Bu, Bk, Xu, Xk     Vector   // partitioned rhs and unknowns vector
}
    Equations organises the identification numbers (IDs) of equations in a
    linear system: [A] ⋅ {x} = {b}. Some of the components of the {x} vector may
    be known in advance—i.e. some {x} components are "prescribed"—therefore the
    system of equations must be reduced first before its solution. The equations
    corresponding to the known {x} values are called "known equations" whereas
    the equations corresponding to unknown {x} values are called "unknown
    equations".

    The right-hand-side vector {b} will also be partitioned into two sets.
    However, the "unknown equations" part will have known values of {b}. If
    needed, the other {b} values can be computed from the "known equations".

    The structure Equations will then hold the IDs of two sets of equations: "u"
    -- unknown {x} values with given {b} values "k" -- known {x} values (the
    corresponding {b} may be post-computed)

        In summary:

         * We define the u-reduced system of equations with the unknown {x} (and known {b})
                         ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄                       ┄┄┄┄┄┄┄
           i.e. the system that needs to be solved is:  [Auu]⋅{xu} = {bu} - [Auk]⋅{xk}

         * We define the k-reduced system of equations with the known {x} values.
                         ┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄                       ┄┄┄┄┄

         * We call the system [A] ⋅ {x} = {b} the "full" system, with all equations.

        The "full" linear system is:

                           [A]⋅{x} = {b}   or   [ Auu Auk ]⋅{ xu ? } = { bu ✓ }
                                                [ Aku Akk ] { xk ✓ }   { bk ? }

        NOTE: the {bu} part is known and the {bk} can be (post-)computed if needed.

        The partitioned system is symbolised as follows:

                                      Auu Auk  u─────> unknown
                                      Aku Akk  k─────> known
                                        u   k
                                        │   └────────> known
                                        └────────────> unknown

        The Equations structure uses arrays to map the indices of equations from the "full"
        system to the reduced systems and vice-versa. We call these arrays "maps" but they
        are not Go maps; they are simple slices of integers and should give better
        performance than maps.

        Four "maps" (slices of integers) are built:
         * UtoF -- reduced u-system to full-system
         * FtoU -- full-system to reduced u-system
         * KtoF -- reduced k-system to full-system
         * FtoK -- full-system to reduced k-system

        EXAMPLE:

         N  = 9          ⇒  total number of equations
         kx = [0  3  6]  ⇒  known x-components (i.e. "prescribed" equations)

           0  1  2  3  4  5  6  7  8                  0   1  2   3   4  5   6   7  8
         0 +--------+--------+------ 0  ◀ known ▶  0 [kk] ku ku [kk] ku ku [kk] ku ku 0
         1 |        |        |       1             1  uk  ○  ○   uk  ○  ○   uk  ○  ○  1
         2 |        |        |       2             2  uk  ○  ○   uk  ○  ○   uk  ○  ○  2
         3 +--------+--------+------ 3  ◀ known ▶  3 [kk] ku ku [kk] ku ku [kk] ku ku 3
         4 |        |        |       4             4  uk  ○  ○   uk  ○  ○   uk  ○  ○  4
         5 |        |        |       5             5  uk  ○  ○   uk  ○  ○   uk  ○  ○  5
         6 +--------+--------+------ 6  ◀ known ▶  6 [kk] ku ku [kk] ku ku [kk] ku ku 6
         7 |        |        |       7             7  uk  ○  ○   uk  ○  ○   uk  ○  ○  7
         8 |        |        |       8             8  uk  ○  ○   uk  ○  ○   uk  ○  ○  8
           0  1  2  3  4  5  6  7  8                  0   1  2   3   4  5   6   7  8
           ▲        ▲        ▲
         known    known    known                             ○ means uu equations

         Nu = 6 => number of "unknown equations"
         Nk = 3 => number of "known equations"
         Nu + Nk = N

         reduced:0  1  2  3  4  5              Example:
         UtoF = [1  2  4  5  7  8]           ⇒ UtoF[3] returns equation # 5 of full system

                 0  1  2  3  4  5  6  7  8
         FtoU = [   0  1     2  3     4  5]  ⇒ FtoU[5] returns equation # 3 of reduced system
                -1       -1       -1         ⇒ -1 indicates 'value not set'

         reduced:0  1  2
         KtoF = [0  3  6]                    ⇒ KtoF[1] returns equation # 3 of full system

                 0  1  2  3  4  5  6  7  8
         FtoK = [0        1        2      ]  ⇒ FtoK[3] returns equation # 1 of reduced system
                   -1 -1    -1 -1    -1 -1   ⇒ -1 indicates 'value not set'

func NewEquations(n int, kx []int) (o *Equations)
    NewEquations creates a new Equations structure

        n  -- total number of equations; i.e. len({x}) in [A]⋅{x}={b}
        kx -- known x-components ⇒ "known equations" [may be unsorted]

func (o *Equations) Alloc(nnz []int, kparts, vectors bool)
    Alloc allocates the A matrices in triplet format (sparse format)

        INPUT:
          nnz -- total number of nonzeros in each part [nnz(Auu), nnz(Auk), nnz(Aku), nnz(Akk)]
                 nnz may be nil, in this case, the following values are assumed:
                          nnz(Auu) = Nu ⋅ Nu
                          nnz(Auk) = Nu ⋅ Nk
                          nnz(Aku) = Nk ⋅ Nu
                          nnz(Akk) = Nk ⋅ Nk
                 Thus, memory is wasted as the size of a fully dense system is considered.
         kparts -- also allocates Aku and Akk
         vectors -- also allocates the partitioned vectors Bu, Bk, Xu, Xk

    OUTPUT:

        The partitioned system (Auu, Auk, Aku, Akk) is stored as member of this object

func (o *Equations) AllocDense(kparts bool)
    AllocDense allocates the A matrices in dense format

        INPUT:
         kparts -- also allocates Aku and Akk
        OUTPUT:
         The partitioned system (Duu, Duk, Dku, Dkk) is stored as member of this object

func (o *Equations) GetAmat() (A *Triplet)
    GetAmat returns the full A matrix (sparse/triplet format) made by Auu, Auk,
    Aku and Akk (e.g. for debugging)

func (o *Equations) Info(full bool)
    Info prints information about Equations

func (o *Equations) JoinVector(b, bu, bk Vector)
    JoinVector joins uknown with known parts of vector

        INPUT:
         bu, bk -- partitioned vectors; e.g. o.Bu, and o.Bk or o.Xu, o.Xk
        OUTPUT:
         b -- pre-allocated b vector such that b = join(bu,bk)

func (o *Equations) Put(I, J int, value float64)
    Put puts component into the right place in partitioned Triplet (Auu, Auk,
    Aku, Akk)

        NOTE: (1) I and J are the equation numbers in the FULL system
              (2) Aku and Akk are ignored if the "kparts" have not been allocated

func (o *Equations) SetDense(A *Matrix, kparts bool)
    SetDense allocates and sets partitioned system in dense format

        kparts -- also computes Aku and Akk

func (o *Equations) Solve(solver SparseSolver, t float64, calcXk, calcBu func(I int, t float64) float64)
    Solve solves linear system (represented by sparse matrices)

        Solve:               {xu} = [Auu]⁻¹ ⋅ ( {bu} - [Auk]⋅{xk} )
        and, if Aku != nil:  {bk} = [Aku]⋅{xu} + [Akk]⋅{xk}

        Input:
         solver -- a pre-configured SparseSolver
         t      -- [optional] a scalar (e.g. time) to be used with calcXk and calcBu
         calcXk -- [optional] a function to calculate the known values of X
         calcBu -- [optional] a function to calculate the right-hand-side Bu

        NOTE: (1) the following must be computed already:
                     [Auu], [Auk] (optionally, [Aku] and [Akk] as well)
              (2) if calcXk is nil, the current values in Xk will be used
              (3) if calcBu is nil, the current values in Bu will be used
              (4) on exit, {bu} is the modified value {bu} - [Auk]⋅{xk} if [Auk] == !nil

        Instead of providing the functions calcXk and calcBu, the vectors {xk} and {bu} can be
        pre-computed. For example, in a FDM grid, the following loops can be used:

              // compute known X values (e.g. Dirichlet boundary conditions)
              for i, I := range KtoF { // i:known, I:full
                  xk[i] = dirichlet(node(I), time)
              }

              // compute RHS
              for i, I := range UtoF { // i:unknown, I:full
                  bu[i] = source(node(I), time)
              }

func (o *Equations) SolveOnce(calcXk, calcBu func(I int, t float64) float64)
    SolveOnce solves linear system just once; thus allocating and discarding a
    linear solver (umfpack) internally. See method Solve() for more details

func (o *Equations) SplitVector(bu, bk, b Vector)
    SplitVector splits full vector b into known and unknown parts

        INPUT:
         b -- full b vector. b = join(bu,bk)
        OUTPUT:
         bu, bk -- partitioned vectors; e.g. o.Bu, and o.Bk or o.Xu, o.Xk

func (o *Equations) Start()
    Start (re)starts index for inserting items using the Put command

type Matrix struct {
	M, N int       // dimensions
	Data []float64 // data array. column-major => Fortran
}
    Matrix implements a column-major representation of a matrix by using a
    linear array that can be passed to Fortran code

        NOTE: the functions related to Matrix do not check for the limits of indices and dimensions.
              Panic may occur then.

        Example:
                   _      _
                  |  0  3  |
              A = |  1  4  |
                  |_ 2  5 _|(m x n)

           data[i+j*m] = A[i][j]

func NewMatrix(m, n int) (o *Matrix)
    NewMatrix allocates a new (empty) Matrix with given (m,n) (row/col sizes)

func NewMatrixDeep2(a [][]float64) (o *Matrix)
    NewMatrixDeep2 allocates a new Matrix from given (Deep2) nested slice. NOTE:
    make sure to have at least 1x1 item

func NewMatrixRaw(m, n int, rawdata []float64) (o *Matrix)
    NewMatrixRaw creates a new Matrix using given raw data

        Input:
          rawdata -- data organized as column-major; e.g. Fortran format
        NOTE:
          (1) rawdata is not copied!
          (2) the external slice rawdata should not be changed or deleted

func (o *Matrix) Add(i, j int, val float64)
    Add adds value to (i,j) location

func (o Matrix) Apply(α float64, another *Matrix)
    Apply sets this matrix with the scaled components of another matrix

        this := α * another   ⇒   this[i] := α * another[i]
        NOTE: "another" may be "this"

func (o *Matrix) ClearBry(diag float64)
    ClearBry clears boundaries

                       _       _                          _       _
        Example:      |  1 2 3  |                        |  1 0 0  |
                  A = |  4 5 6  |  ⇒  clear(1.0)  ⇒  A = |  0 5 0  |
                      |_ 7 8 9 _|                        |_ 0 0 1 _|

func (o *Matrix) ClearRC(rows, cols []int, diag float64)
    ClearRC clear rows and columns and set diagonal components

                       _         _                                     _         _
        Example:      |  1 2 3 4  |                                   |  1 2 3 4  |
                  A = |  5 6 7 8  |  ⇒  clear([1,2], [], 1.0)  ⇒  A = |  0 1 0 0  |
                      |_ 4 3 2 1 _|                                   |_ 0 0 1 0 _|

func (o *Matrix) Col(j int) Vector
    Col access column j of this matrix. No copies are made since the internal
    data are in col-major format already. NOTE: this method can be used to
    modify the columns; e.g. with o.Col(0)[0] = 123

func (o *Matrix) CopyInto(result *Matrix, α float64)
    CopyInto copies the scaled components of this matrix into another one
    (result)

        result := α * this   ⇒   result[ij] := α * this[ij]

func (o *Matrix) Det() (det float64)
    Det computes the determinant of matrix using the LU factorization

        NOTE: this method may fail due to overflow...

func (o *Matrix) ExtractCols(start, endp1 int) (reduced *Matrix)
    ExtractCols returns columns from j=start to j=endp1-1

        start -- first column
        endp1 -- "end-plus-one", the number of the last requested column + 1

func (o *Matrix) Fill(val float64)
    Fill fills this matrix with a single number val

        aij = val

func (o *Matrix) Get(i, j int) float64
    Get gets value

func (o *Matrix) GetCol(j int) (col Vector)
    GetCol returns column j of this matrix

func (o *Matrix) GetComplex() (b *MatrixC)
    GetComplex returns a complex version of this matrix

func (o *Matrix) GetCopy() (clone *Matrix)
    GetCopy returns a copy of this matrix

func (o *Matrix) GetDeep2() (M [][]float64)
    GetDeep2 returns nested slice representation

func (o *Matrix) GetRow(i int) (row Vector)
    GetRow returns row i of this matrix

func (o *Matrix) GetTranspose() (tran *Matrix)
    GetTranspose returns the tranpose matrix

func (o *Matrix) Largest(den float64) (largest float64)
    Largest returns the largest component |a[ij]| of this matrix, normalised by
    den

        largest := |a[ij]| / den

func (o *Matrix) MaxDiff(another *Matrix) (maxdiff float64)
    MaxDiff returns the maximum difference between the components of this and
    another matrix

func (o *Matrix) NormFrob() (nrm float64)
    NormFrob returns the Frobenious norm of this matrix

        nrm := ‖a‖_F = sqrt(Σ_i Σ_j a[ij]⋅a[ij]) = ‖a‖_2

func (o *Matrix) NormInf() (nrm float64)
    NormInf returns the infinite norm of this matrix

        nrm := ‖a‖_∞ = max_i ( Σ_j a[ij] )

func (o *Matrix) Print(nfmt string) (l string)
    Print prints matrix (without commas or brackets)

func (o *Matrix) PrintGo(nfmt string) (l string)
    PrintGo prints matrix in Go format

func (o *Matrix) PrintPy(nfmt string) (l string)
    PrintPy prints matrix in Python format

func (o *Matrix) Set(i, j int, val float64)
    Set sets value

func (o *Matrix) SetCol(j int, value float64)
    SetCol sets the values of a column j with a single value

func (o *Matrix) SetDiag(val float64)
    SetDiag sets diagonal matrix with diagonal components equal to val

func (o *Matrix) SetFromDeep2(a [][]float64)
    SetFromDeep2 sets matrix with data from a nested slice (Deep2) structure

type MatrixC struct {
	M, N int          // dimensions
	Data []complex128 // data array. column-major => Fortran
}
    MatrixC implements a column-major representation of a matrix of complex
    numbers by using a linear array that can be passed to Fortran code.

        NOTE: the functions related to MatrixC do not check for the limits of indices and dimensions.
              Panic may occur then.

        Example:
                   _            _
                  |  0+0i  3+3i  |
              A = |  1+1i  4+4i  |
                  |_ 2+2i  5+5i _|(m x n)

           data[i+j*m] = A[i][j]

func NewMatrixC(m, n int) (o *MatrixC)
    NewMatrixC allocates a new (empty) MatrixC with given (m,n) (row/col sizes)

func NewMatrixDeep2c(a [][]complex128) (o *MatrixC)
    NewMatrixDeep2c allocates a new MatrixC from given (Deep2c) nested slice.
    NOTE: make sure to have at least 1x1 items

func (o *MatrixC) Add(i, j int, val complex128)
    Add adds value to (i,j) location

func (o MatrixC) Apply(α complex128, another *MatrixC)
    Apply sets this matrix with the scaled components of another matrix

        this := α * another   ⇒   this[i] := α * another[i]
        NOTE: "another" may be "this"

func (o *MatrixC) Col(j int) VectorC
    Col access column j of this matrix. No copies are made since the internal
    data are in col-major format already. NOTE: this method can be used to
    modify the columns; e.g. with o.Col(0)[0] = 123

func (o *MatrixC) Fill(val complex128)
    Fill fills this matrix with a single number val

        aij = val

func (o *MatrixC) Get(i, j int) complex128
    Get gets value

func (o *MatrixC) GetCol(j int) (col VectorC)
    GetCol returns column j of this matrix

func (o *MatrixC) GetColReal(j int, checkZeroImag bool) (col Vector)
    GetColReal returns column j of this matrix considering that the imaginary
    part is zero

func (o *MatrixC) GetCopy() (clone *MatrixC)
    GetCopy returns a copy of this matrix

func (o *MatrixC) GetDeep2() (M [][]complex128)
    GetDeep2 returns nested slice representation

func (o *MatrixC) GetRow(i int) (row VectorC)
    GetRow returns row i of this matrix

func (o *MatrixC) GetTranspose() (tran *MatrixC)
    GetTranspose returns the tranpose matrix

func (o *MatrixC) Print(nfmtR, nfmtI string) (l string)
    Print prints matrix (without commas or brackets). NOTE: if non-empty, nfmtI
    must have '+' e.g. %+g

func (o *MatrixC) PrintGo(nfmtR, nfmtI string) (l string)
    PrintGo prints matrix in Go format NOTE: if non-empty, nfmtI must have '+'
    e.g. %+g

func (o *MatrixC) PrintPy(nfmtR, nfmtI string) (l string)
    PrintPy prints matrix in Python format NOTE: if non-empty, nfmtI must have
    '+' e.g. %+g

func (o *MatrixC) Set(i, j int, val complex128)
    Set sets value

func (o *MatrixC) SetFromDeep2c(a [][]complex128)
    SetFromDeep2c sets matrix with data from a nested slice (Deep2c) structure

type Mumps struct {
	// Has unexported fields.
}
    Mumps wraps the MUMPS solver

func (o *Mumps) Fact()
    Fact performs the factorisation

func (o *Mumps) Free()
    Free clears extra memory allocated by MUMPS

func (o *Mumps) Init(t *Triplet, args *SparseConfig)
    Init initialises mumps for sparse linear systems with real numbers args may
    be nil

func (o *Mumps) Solve(x, b Vector, bIsDistr bool)
    Solve solves sparse linear systems using MUMPS or MUMPS

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

        bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.

type MumpsC struct {
	// Has unexported fields.
}
    MumpsC wraps the MUMPS solver (complex version)

func (o *MumpsC) Fact()
    Fact performs the factorisation

func (o *MumpsC) Free()
    Free clears extra memory allocated by MUMPS

func (o *MumpsC) Init(t *TripletC, args *SparseConfig)
    Init initialises mumps for sparse linear systems with real numbers args may
    be nil

func (o *MumpsC) Solve(x, b VectorC, bIsDistr bool)
    Solve solves sparse linear systems using MUMPS or MUMPS

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

        bIsDistr -- this flag tells that the right-hand-side vector 'b' is distributed.

type SparseConfig struct {
	// external
	Symmetric bool   // indicates symmetric system
	SymPosDef bool   // indicates symmetric-positive-defined system
	Verbose   bool   // run on Verbose mode
	Guess     Vector // initial guess for iterative solvers [may be nil]

	// MUMPS control parameters (check MUMPS solver manual)
	MumpsIncreaseOfWorkingSpacePct int // ICNTL(14) default = 100%
	MumpsMaxMemoryPerProcessor     int // ICNTL(23) default = 2000Mb

	// Has unexported fields.
}
    The SparseConfig structure holds configuration arguments for sparse solvers

func NewSparseConfig(comm *mpi.Communicator) (o *SparseConfig)
    NewSparseConfig returns a new SparseConfig comm may be nil

func (o *SparseConfig) SetMumpsOrdering(ordering string)
    SetMumpsOrdering sets ordering for MUMPS solver ordering -- "" or "amf"
    [default]

        "amf", "scotch", "pord", "metis", "qamd", "auto"

    ICNTL(7)

        0: "amd" Approximate Minimum Degree (AMD)
        2: "amf" Approximate Minimum Fill (AMF)
        3: "scotch" SCOTCH5 package is used if previously installed by the user otherwise treated as 7.
        4: "pord" PORD6 is used if previously installed by the user otherwise treated as 7.
        5: "metis" Metis7 package is used if previously installed by the user otherwise treated as 7.
        6: "qamd" Approximate Minimum Degree with automatic quasi-dense row detection (QAMD) is used.
        7: "auto" automatic choice by the software during analysis phase. This choice will depend on the
            ordering packages made available, on the matrix (type and size), and on the number of processors.

func (o *SparseConfig) SetMumpsScaling(scaling string)
    SetMumpsScaling sets scaling for MUMPS solver scaling -- "" or "rcit"
    [default]

        "no", "diag", "col", "rcinf", "rcit", "rrcit", "auto"

    ICNTL(8)

        0: "no" No scaling applied/computed.
        1: "diag" Diagonal scaling computed during the numerical factorization phase,
        3: "col" Column scaling computed during the numerical factorization phase,
        4: "rcinf" Row and column scaling based on infinite row/column norms, computed during the numerical
           factorization phase,
        7: "rcit" Simultaneous row and column iterative scaling based on [41] and [15] computed during the
           numerical factorization phase.
        8: "rrcit" Similar to 7 but more rigorous and expensive to compute; computed during the numerical
           factorization phase.
        77: "auto" Automatic choice of the value of ICNTL(8) done during analy

type SparseSolver interface {
	Init(t *Triplet, args *SparseConfig)
	Free()
	Fact()
	Solve(x, b Vector, bIsDistr bool)
}
    SparseSolver solves sparse linear systems using UMFPACK or MUMPS

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

func NewSparseSolver(kind string) SparseSolver
    NewSparseSolver finds a SparseSolver in database or panic

        kind -- "umfpack" or "mumps"
        NOTE: remember to call Free() to release allocated resources

type SparseSolverC interface {
	Init(t *TripletC, args *SparseConfig)
	Free()
	Fact()
	Solve(x, b VectorC, bIsDistr bool)
}
    SparseSolverC solves sparse linear systems using UMFPACK or MUMPS (complex
    version)

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

func NewSparseSolverC(kind string) SparseSolverC
    NewSparseSolverC finds a SparseSolver in database or panic

        NOTE: remember to call Free() to release allocated resources

type Triplet struct {
	// Has unexported fields.
}
    Triplet is a simple representation of a sparse matrix, where the indices and
    values of this matrix are stored directly.

func NewTriplet(m, n, max int) (o *Triplet)
    NewTriplet returns a new Triplet. This is a wrapper to new(Triplet) followed
    by Init()

func (o *Triplet) Init(m, n, max int)
    Init allocates all memory required to hold a sparse matrix in triplet form

func (o *Triplet) Len() int
    Len returns the number of items just inserted in the triplet

func (o *Triplet) Max() int
    Max returns the maximum number of items that can be inserted in the triplet

func (o *Triplet) Put(i, j int, x float64)
    Put inserts an element to a pre-allocated (with Init) triplet matrix

func (o *Triplet) PutCCMatAndMatT(a *CCMatrix)
    PutCCMatAndMatT adds the content of a compressed-column matrix "a" and its
    transpose "at" to triplet "o" ex: 0 1 2 3 4 5

        [... ... ... a00 a10 ...] 0
        [... ... ... a01 a11 ...] 1
        [... ... ... a02 a12 ...] 2      [. at  .]
        [a00 a01 a02 ... ... ...] 3  =>  [a  .  .]
        [a10 a11 a12 ... ... ...] 4      [.  .  .]
        [... ... ... ... ... ...] 5

func (o *Triplet) PutMatAndMatT(a *Triplet)
    PutMatAndMatT adds the content of a matrix "a" and its transpose "at" to
    triplet "o" ex: 0 1 2 3 4 5

        [... ... ... a00 a10 ...] 0
        [... ... ... a01 a11 ...] 1
        [... ... ... a02 a12 ...] 2      [. at  .]
        [a00 a01 a02 ... ... ...] 3  =>  [a  .  .]
        [a10 a11 a12 ... ... ...] 4      [.  .  .]
        [... ... ... ... ... ...] 5

func (o *Triplet) ReadSmat(filename string, oneBasedIndices ...bool)
    ReadSmat reads a SMAT file or a MatrixMarket file

        About the .smat file:
         - lines starting with the percent or hashtag mark (% or #) are ignored (they are comments)
         - the first line contains the number of rows, columns, and non-zero entries
         - the following lines contain the indices of row, column, and the non-zero entry
         - indices are 0-based, unless oneBasedIndices == true
         - with one-based indices, smat is similar to MatrixMarket without the header

        Example of .smat file:

          # this is a comment
          % this is also a comment
          m n nnz
           i j x
            ...
           i j x

        filename -- the name of the .smat file
        oneBasedIndices -- if true, 1 is removed from all indices (input like MatrixMarket file format)

func (o *Triplet) Size() (m, n int)
    Size returns the row/column size of the matrix represented by the Triplet

func (o *Triplet) Start()
    Start (re)starts index for inserting items using the Put command

func (o *Triplet) ToDense() (a *Matrix)
    ToDense returns the dense matrix corresponding to this Triplet

func (t *Triplet) ToMatrix(a *CCMatrix) *CCMatrix
    ToMatrix converts a sparse matrix in triplet form to column-compressed form
    using Umfpack's routines. "realloc_a" indicates whether the internal "a"
    matrix must be reallocated or not, for instance, in case the structure of
    the triplet has changed.

        INPUT:
         a -- a previous CCMatrix to be filled in; otherwise, "nil" tells to allocate a new one
        OUTPUT:
         the previous "a" matrix or a pointer to a new one

func (o *Triplet) WriteSmat(dirout, fnkey string, tol float64, oneBasedIndices ...bool) (cmat *CCMatrix)
    WriteSmat writes a SMAT file (that can be visualised with vismatrix) or a
    MatrixMarket file

        About the .smat file:
         - lines starting with the percent or hashtag mark (% or #) are ignored (they are comments)
         - the first line contains the number of rows, columns, and non-zero entries
         - the following lines contain the indices of row, column, and the non-zero entry
         - indices are 0-based, unless oneBasedIndices == true
         - with one-based indices, smat is similar to MatrixMarket without the header

        Example of .smat file:

          # this is a comment
          % this is also a comment
          m n nnz
           i j x
            ...
           i j x

        NOTE: this method will create a CCMatrix first because
              duplicates must be added before saving the file

        dirout -- directory for output. will be created
        fnkey  -- filename key (filename without extension). ".smat" will be added
        tol    -- tolerance to skip zero values
        oneBasedIndices -- if true, 1 is added to all indices (like MatrixMarket file format)

type TripletC struct {
	// Has unexported fields.
}
    TripletC is a simple representation of a sparse matrix, where the indices
    and values of this matrix are stored directly. (complex version)

func NewTripletC(m, n, max int) (o *TripletC)
    NewTripletC returns a new TripletC. This is a wrapper to new(TripletC)
    followed by Init()

func (o *TripletC) Init(m, n, max int)
    Init allocates all memory required to hold a sparse matrix in triplet
    (complex) form

func (o *TripletC) Len() int
    Len returns the number of items just inserted in the complex triplet

func (o *TripletC) Max() int
    Max returns the maximum number of items that can be inserted in the complex
    triplet

func (o *TripletC) Put(i, j int, x complex128)
    Put inserts an element to a pre-allocated (with Init) triplet (complex)
    matrix

func (o *TripletC) ReadSmat(filename string)
    ReadSmat reads ".smat" file

        m n nnz
         i j xReal xImag
              ...
         i j xReal xImag

func (o *TripletC) Start()
    Start (re)starts index for inserting items using the Put command

func (o *TripletC) ToDense() (a *MatrixC)
    ToDense returns the dense matrix corresponding to this Triplet

func (t *TripletC) ToMatrix(a *CCMatrixC) *CCMatrixC
    ToMatrix converts a sparse matrix in triplet form with complex numbers to
    column-compressed form. "realloc_a" indicates whether the internal "a"
    matrix must be reallocated or not, for instance, in case the structure of
    the triplet has changed.

        INPUT:
         a -- a previous CCMatrixC to be filled in; otherwise, "nil" tells to allocate a new one
        OUTPUT:
         the previous "a" matrix or a pointer to a new one

func (o *TripletC) WriteSmat(dirout, fnkey string, tol float64) (cmat *CCMatrixC)
    WriteSmat writes a ".smat" file that can be visualised with vismatrix

        NOTE: this method will create a CCMatrixC first because
              duplicates must be added before saving the file

        dirout -- directory for output. will be created
        fnkey  -- filename key (filename without extension). ".smat" will be added
        tol    -- tolerance to skip zero values

type Umfpack struct {
	// Has unexported fields.
}
    Umfpack wraps the UMFPACK solver

func (o *Umfpack) Fact()
    Fact performs the factorisation

func (o *Umfpack) Free()
    Free clears extra memory allocated by UMFPACK

func (o *Umfpack) Init(t *Triplet, args *SparseConfig)
    Init initialises umfpack for sparse linear systems with real numbers args
    may be nil

func (o *Umfpack) Solve(x, b Vector, dummy bool)
    Solve solves sparse linear systems using UMFPACK or MUMPS

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

type UmfpackC struct {
	// Has unexported fields.
}
    UmfpackC wraps the UMFPACK solver (complex version)

func (o *UmfpackC) Fact()
    Fact performs the factorisation

func (o *UmfpackC) Free()
    Free clears extra memory allocated by UMFPACK

func (o *UmfpackC) Init(t *TripletC, args *SparseConfig)
    Init initialises umfpack for sparse linear systems with real numbers args
    may be nil

func (o *UmfpackC) Solve(x, b VectorC, dummy bool)
    Solve solves sparse linear systems using UMFPACK or MUMPS

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

type Vector []float64
    Vector defines the vector type for real numbers simply as a slice of float64

func NewVector(m int) Vector
    NewVector returns a new vector with size m

func NewVectorMapped(m int, f func(i int) float64) (o Vector)
    NewVectorMapped returns a new vector after applying a function over all of
    its components

        new: vi = f(i)

func NewVectorSlice(v []float64) Vector
    NewVectorSlice returns a new vector from given Slice NOTE: This is
    equivalent to cast a slice to Vector as in:

        v := la.Vector([]float64{1,2,3})

func SpSolve(A *Triplet, b Vector) (x Vector)
    SpSolve solves a sparse linear system (using UMFPACK)

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

func (o Vector) Accum() (sum float64)
    Accum sum/accumulates all components in a vector

        sum := Σ_i v[i]

func (o Vector) Apply(α float64, another Vector)
    Apply sets this vector with the scaled components of another vector

        this := α * another   ⇒   this[i] := α * another[i]
        NOTE: "another" may be "this"

func (o Vector) ApplyFunc(f func(i int, x float64) float64)
    ApplyFunc runs a function over all components of a vector

        vi = f(i,vi)

func (o Vector) Fill(val float64)
    Fill fills this vector with a single number val

        vi = val

func (o Vector) GetCopy() (clone Vector)
    GetCopy returns a copy of this vector

        b := a

func (o Vector) GetUnit() (unit Vector)
    GetUnit returns the unit vector parallel to this vector

        b := a / norm(a)

func (o Vector) Largest(den float64) (largest float64)
    Largest returns the largest component |u[i]| of this vector, normalised by
    den

        largest := |u[i]| / den

func (o Vector) Max() (max float64)
    Max returns the maximum component of a vector

func (o Vector) Min() (min float64)
    Min returns the minimum component of a vector

func (o Vector) MinMax() (min, max float64)
    MinMax returns the min and max components of a vector

func (o Vector) Norm() (nrm float64)
    Norm returns the Euclidean norm of a vector:

        nrm := ‖v‖

func (o Vector) NormDiff(v Vector) (nrm float64)
    NormDiff returns the Euclidean norm of the difference:

        nrm := ||u - v||

func (o Vector) Rms() (rms float64)
    Rms returns the root-mean-square of this vector

                     ________________________
                    /     ————            2
                   /  1   \    /         \
        rms =  \  /  ———  /    | this[i] |
                \/    N   ———— \         /

type VectorC []complex128
    VectorC defines the vector type for complex numbers simply as a slice of
    complex128

func NewVectorC(m int) VectorC
    NewVectorC returns a new vector with size m

func NewVectorMappedC(m int, f func(i int) complex128) (o VectorC)
    NewVectorMappedC returns a new vector after applying a function over all of
    its components

        new: vi = f(i)

func SpSolveC(A *TripletC, b VectorC) (x VectorC)
    SpSolveC solves a sparse linear system (using UMFPACK) (complex version)

        Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b

func (o VectorC) Apply(α complex128, another VectorC)
    Apply sets this vector with the scaled components of another vector

        this := α * another   ⇒   this[i] := α * another[i]
        NOTE: "another" may be "this"

func (o VectorC) ApplyFunc(f func(i int, x complex128) complex128)
    ApplyFunc runs a function over all components of a vector

        vi = f(i,vi)

func (o VectorC) Fill(s complex128)
    Fill fills a vector with a single number s:

        v := s*ones(len(v))  =>  vi = s

func (o VectorC) GetCopy() (clone VectorC)
    GetCopy returns a copy of this vector

        b := a

func (o VectorC) JoinRealImag(xR, xI Vector)
    JoinRealImag sets this vector with two vectors having the real and imaginary
    parts

        this := complex(xR, xI)
        NOTE: len(xR) == len(xI) == len(this)

func (o VectorC) MaxDiff(b VectorC) float64
    MaxDiff returns the maximum difference between the components of two vectors

func (o VectorC) Norm() (nrm complex128)
    Norm returns the Euclidean norm of a vector:

        nrm := ‖v‖

func (o VectorC) SplitRealImag(xR, xI Vector)
    SplitRealImag splits this vector into two vectors with the real and
    imaginary parts

        xR := real(this)
        xI := imag(this)
        NOTE: xR and xI must be pre-allocated with length = len(this)

```
