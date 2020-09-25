# Gosl. la/oblas. Wrapper to OpenBLAS

This subpackge implements a light wrapper to OpenBLAS. Therefore, its routines are a little more
_lower level_ than the ones in the parent package `la`.

[Check also OpenBLAS](https://github.com/xianyi/OpenBLAS).

## API

**go doc**

```
package oblas // import "gosl/la/oblas"

Package oblas implements lower-level linear algebra routines using OpenBLAS
for maximum efficiency. This package uses column-major representation for
matrices.

     Example of col-major data:
               _      _
              |  0  3  |
          A = |  1  4  |            ⇒     a = [0, 1, 2, 3, 4, 5]
              |_ 2  5 _|(m x n)

          a[i+j*m] = A[i][j]

    NOTE: the functions here do not check for the limits of indices. Be careful.
          Panic may occur then.

FUNCTIONS

func ColMajorCtoSlice(m, n int, data []complex128) (a [][]complex128)
    ColMajorCtoSlice converts col-major matrix to nested slice

func ColMajorToSlice(m, n int, data []float64) (a [][]float64)
    ColMajorToSlice converts col-major matrix to nested slice

func Daxpy(n int, alpha float64, x []float64, incx int, y []float64, incy int)
    Daxpy computes constant times a vector plus a vector.

        See: http://www.netlib.org/lapack/explore-html/d9/dcd/daxpy_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-axpy

        y += alpha*x + y

func Ddot(n int, x []float64, incx int, y []float64, incy int) (res float64)
    Ddot forms the dot product of two vectors. Uses unrolled loops for
    increments equal to one.

        See: http://www.netlib.org/lapack/explore-html/d5/df6/ddot_8f.html

func Dgeev(calcVl, calcVr bool, n int, a []float64, lda int, wr []float64, wi, vl []float64, ldvl int, vr []float64, ldvr int)
    Dgeev computes for an N-by-N real nonsymmetric matrix A, the eigenvalues
    and, optionally, the left and/or right eigenvectors.

        See: http://www.netlib.org/lapack/explore-html/d9/d28/dgeev_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-geev

        See: https://www.nag.co.uk/numeric/fl/nagdoc_fl26/html/f08/f08naf.html

        The right eigenvector v(j) of A satisfies

                         A * v(j) = lambda(j) * v(j)

        where lambda(j) is its eigenvalue.

        The left eigenvector u(j) of A satisfies

                      u(j)**H * A = lambda(j) * u(j)**H

        where u(j)**H denotes the conjugate-transpose of u(j).

        The computed eigenvectors are normalized to have Euclidean norm
        equal to 1 and largest component real.

func Dgemm(transA, transB bool, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
    Dgemm performs one of the matrix-matrix operations

        false,false:  C_{m,n} := α ⋅ A_{m,k} ⋅ B_{k,n}  +  β ⋅ C_{m,n}
        false,true:   C_{m,n} := α ⋅ A_{m,k} ⋅ B_{n,k}  +  β ⋅ C_{m,n}
        true, false:  C_{m,n} := α ⋅ A_{k,m} ⋅ B_{k,n}  +  β ⋅ C_{m,n}
        true, true:   C_{m,n} := α ⋅ A_{k,m} ⋅ B_{n,k}  +  β ⋅ C_{m,n}

        see: http://www.netlib.org/lapack/explore-html/d7/d2b/dgemm_8f.html

        see: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemm

           C := alpha*op( A )*op( B ) + beta*C,

        where  op( X ) is one of

           op( X ) = X   or   op( X ) = X**T,

        alpha and beta are scalars, and A, B and C are matrices, with op( A )
        an m by k matrix,  op( B )  a  k by n matrix and  C an m by n matrix.

func Dgemv(trans bool, m, n int, alpha float64, a []float64, lda int, x []float64, incx int, beta float64, y []float64, incy int)
    Dgemv performs one of the matrix-vector operations

        See: http://www.netlib.org/lapack/explore-html/dc/da8/dgemv_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemv

           y := alpha*A*x + beta*y,   or   y := alpha*A**T*x + beta*y,

        where alpha and beta are scalars, x and y are vectors and A is an
        m by n matrix.
           trans=false     y := alpha*A*x + beta*y.

           trans=true      y := alpha*A**T*x + beta*y.

func Dger(m, n int, alpha float64, x []float64, incx int, y []float64, incy int, a []float64, lda int)
    Dger performs the rank 1 operation

        See: http://www.netlib.org/lapack/explore-html/dc/da8/dger_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-ger

          A := alpha*x*y**T + A,

    where alpha is a scalar, x is an m element vector, y is an n element vector
    and A is an m by n matrix.

func Dgesv(n, nrhs int, a []float64, lda int, ipiv []int32, b []float64, ldb int)
    Dgesv computes the solution to a real system of linear equations.

        See: http://www.netlib.org/lapack/explore-html/d8/d72/dgesv_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-gesv

        The system is:

           A * X = B,

        where A is an N-by-N matrix and X and B are N-by-NRHS matrices.

        The LU decomposition with partial pivoting and row interchanges is
        used to factor A as

           A = P * L * U,

        where P is a permutation matrix, L is unit lower triangular, and U is
        upper triangular.  The factored form of A is then used to solve the
        system of equations A * X = B.

        NOTE: matrix 'a' will be modified

func Dgesvd(jobu, jobvt rune, m, n int, a []float64, lda int, s []float64, u []float64, ldu int, vt []float64, ldvt int, superb []float64)
    Dgesvd computes the singular value decomposition (SVD) of a real M-by-N
    matrix A, optionally computing the left and/or right singular vectors.

        See: http://www.netlib.org/lapack/explore-html/d8/d2d/dgesvd_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-gesvd

        The SVD is written

             A = U * SIGMA * transpose(V)

        where SIGMA is an M-by-N matrix which is zero except for its
        min(m,n) diagonal elements, U is an M-by-M orthogonal matrix, and
        V is an N-by-N orthogonal matrix.  The diagonal elements of SIGMA
        are the singular values of A; they are real and non-negative, and
        are returned in descending order.  The first min(m,n) columns of
        U and V are the left and right singular vectors of A.

        Note that the routine returns V**T, not V.

        NOTE: matrix 'a' will be modified

func Dgetrf(m, n int, a []float64, lda int, ipiv []int32)
    Dgetrf computes an LU factorization of a general M-by-N matrix A using
    partial pivoting with row interchanges.

        See: http://www.netlib.org/lapack/explore-html/d3/d6a/dgetrf_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-getrf

        The factorization has the form
           A = P * L * U
        where P is a permutation matrix, L is lower triangular with unit
        diagonal elements (lower trapezoidal if m > n), and U is upper
        triangular (upper trapezoidal if m < n).

        This is the right-looking Level 3 BLAS version of the algorithm.

        NOTE: (1) matrix 'a' will be modified
              (2) ipiv indices are 1-based (i.e. Fortran)

func Dgetri(n int, a []float64, lda int, ipiv []int32)
    Dgetri computes the inverse of a matrix using the LU factorization computed
    by DGETRF.

        See: http://www.netlib.org/lapack/explore-html/df/da4/dgetri_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-getri

        This method inverts U and then computes inv(A) by solving the system
        inv(A)*L = inv(U) for inv(A).

func Dpotrf(up bool, n int, a []float64, lda int)
    Dpotrf computes the Cholesky factorization of a real symmetric positive
    definite matrix A.

        See: http://www.netlib.org/lapack/explore-html/d0/d8a/dpotrf_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-potrf

        The factorization has the form

           A = U**T * U,  if UPLO = 'U'

        or

           A = L  * L**T,  if UPLO = 'L'

        where U is an upper triangular matrix and L is lower triangular.

        This is the block version of the algorithm, calling Level 3 BLAS.

func Dscal(n int, alpha float64, x []float64, incx int)
    Dscal scales a vector by a constant. Uses unrolled loops for increment equal
    to 1.

        See: http://www.netlib.org/lapack/explore-html/d4/dd0/dscal_8f.html

func Dsyrk(up, trans bool, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int)
    Dsyrk performs one of the symmetric rank k operations

        See: http://www.netlib.org/lapack/explore-html/dc/d05/dsyrk_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-syrk

           C := alpha*A*A**T + beta*C,

        or

           C := alpha*A**T*A + beta*C,

        where  alpha and beta  are scalars, C is an  n by n  symmetric matrix
        and  A  is an  n by k  matrix in the first case and a  k by n  matrix
        in the second case.

func EigenvecsBuild(vv []complex128, wr, wi, v []float64)
    EigenvecsBuild builds complex eigenvectros created by Dgeev function

        INPUT:
         wr, wi -- real and imag parts of eigenvalues
         v      -- left or right eigenvectors from Dgeev
        OUTPUT:
         vv -- complex version of left or right eigenvector [pre-allocated]
        NOTE (no checks made)
         n = len(wr) = len(wi) = len(v)
         2 * n = len(vv)

func EigenvecsBuildBoth(vvl, vvr []complex128, wr, wi, vl, vr []float64)
    EigenvecsBuildBoth builds complex left and right eigenvectros created by
    Dgeev function

        INPUT:
         wr, wi -- real and imag parts of eigenvalues
         vl, vr -- left and right eigenvectors from Dgeev
        OUTPUT:
         vvl, vvr -- complex version of left and right eigenvectors [pre-allocated]
        NOTE (no checks made)
         n = len(wr) = len(wi) = len(vl) = len(vr)
         2 * n = len(vvl) = len(vvr)

func ExtractCol(j, m, n int, A []float64) (colj []float64)
    ExtractCol extracts j column from (m,n) col-major matrix

func ExtractColC(j, m, n int, A []complex128) (colj []complex128)
    ExtractColC extracts j column from (m,n) col-major matrix (complex version)

func ExtractRow(i, m, n int, A []float64) (rowi []float64)
    ExtractRow extracts i row from (m,n) col-major matrix

func ExtractRowC(i, m, n int, A []complex128) (rowi []complex128)
    ExtractRowC extracts i row from (m,n) col-major matrix (complex version)

func GetJoinComplex(vReal, vImag []float64) (v []complex128)
    GetJoinComplex joins real and imag parts of array

func GetSplitComplex(v []complex128) (vReal, vImag []float64)
    GetSplitComplex splits real and imag parts of array

func JoinComplex(v []complex128, vReal, vImag []float64)
    JoinComplex joins real and imag parts of array

func PrintColMajor(m, n int, data []float64, nfmt string) (l string)
    PrintColMajor prints matrix (without commas or brackets)

func PrintColMajorC(m, n int, data []complex128, nfmtR, nfmtI string) (l string)
    PrintColMajorC prints matrix (without commas or brackets). NOTE: if
    non-empty, nfmtI must have '+' e.g. %+g

func PrintColMajorCgo(m, n int, data []complex128, nfmtR, nfmtI string) (l string)
    PrintColMajorCgo prints matrix in Go format NOTE: if non-empty, nfmtI must
    have '+' e.g. %+g

func PrintColMajorCpy(m, n int, data []complex128, nfmtR, nfmtI string) (l string)
    PrintColMajorCpy prints matrix in Python format NOTE: if non-empty, nfmtI
    must have '+' e.g. %+g

func PrintColMajorGo(m, n int, data []float64, nfmt string) (l string)
    PrintColMajorGo prints matrix in Go format

func PrintColMajorPy(m, n int, data []float64, nfmt string) (l string)
    PrintColMajorPy prints matrix in Python format

func SetNumThreads(n int)
    SetNumThreads sets the number of threads in OpenBLAS

func SliceToColMajor(a [][]float64) (data []float64)
    SliceToColMajor converts nested slice into an array representing a col-major
    matrix

        Example:
                  _      _
                 |  0  3  |
             a = |  1  4  |            ⇒     data = [0, 1, 2, 3, 4, 5]
                 |_ 2  5 _|(m x n)

             data[i+j*m] = a[i][j]

        NOTE: make sure to have at least 1x1 item

func SliceToColMajorC(a [][]complex128) (data []complex128)
    SliceToColMajorC converts nested slice into an array representing a
    col-major matrix of complex numbers.

        Example:
                  _            _
                 |  0+0i  3+3i  |
             a = |  1+1i  4+4i  |          ⇒   data = [0+0i, 1+1i, 2+2i, 3+3i, 4+4i, 5+5i]
                 |_ 2+2i  5+5i _|(m x n)

             data[i+j*m] = a[i][j]

        NOTE: make sure to have at least 1x1 item

func SplitComplex(vReal, vImag []float64, v []complex128)
    SplitComplex splits real and imag parts of array

func Zaxpy(n int, alpha complex128, x []complex128, incx int, y []complex128, incy int)
    Zaxpy computes constant times a vector plus a vector.

        See: http://www.netlib.org/lapack/explore-html/d7/db2/zaxpy_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-axpy

        y += alpha*x + y

func Zgemm(transA, transB bool, m, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
    Zgemm performs one of the matrix-matrix operations

        see: http://www.netlib.org/lapack/explore-html/d7/d76/zgemm_8f.html

        see: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemm

           C := alpha*op( A )*op( B ) + beta*C,

        where  op( X ) is one of

           op( X ) = X   or   op( X ) = X**T   or   op( X ) = X**H,

        alpha and beta are scalars, and A, B and C are matrices, with op( A )
        an m by k matrix,  op( B )  a  k by n matrix and  C an m by n matrix.

func Zgemv(trans bool, m, n int, alpha complex128, a []complex128, lda int, x []complex128, incx int, beta complex128, y []complex128, incy int)
    Zgemv performs one of the matrix-vector operations.

        See: http://www.netlib.org/lapack/explore-html/db/d40/zgemv_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemv

           y := alpha*A*x + beta*y,   or   y := alpha*A**T*x + beta*y,   or

           y := alpha*A**H*x + beta*y,

        where alpha and beta are scalars, x and y are vectors and A is an
        m by n matrix.

func Zgesv(n, nrhs int, a []complex128, lda int, ipiv []int32, b []complex128, ldb int)
    Zgesv computes the solution to a complex system of linear equations.

        See: http://www.netlib.org/lapack/explore-html/d1/ddc/zgesv_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-gesv

        The system is:

           A * X = B,

        where A is an N-by-N matrix and X and B are N-by-NRHS matrices.

        The LU decomposition with partial pivoting and row interchanges is
        used to factor A as

           A = P * L * U,

        where P is a permutation matrix, L is unit lower triangular, and U is
        upper triangular.  The factored form of A is then used to solve the
        system of equations A * X = B.

        NOTE: matrix 'a' will be modified

func Zgesvd(jobu, jobvt rune, m, n int, a []complex128, lda int, s []float64, u []complex128, ldu int, vt []complex128, ldvt int, superb []float64)
    Zgesvd computes the singular value decomposition (SVD) of a complex M-by-N
    matrix A, optionally computing the left and/or right singular vectors.

        See: http://www.netlib.org/lapack/explore-html/d6/d42/zgesvd_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-gesvd

        The SVD is written

             A = U * SIGMA * conjugate-transpose(V)

        where SIGMA is an M-by-N matrix which is zero except for its
        min(m,n) diagonal elements, U is an M-by-M unitary matrix, and
        V is an N-by-N unitary matrix.  The diagonal elements of SIGMA
        are the singular values of A; they are real and non-negative, and
        are returned in descending order.  The first min(m,n) columns of
        U and V are the left and right singular vectors of A.

        Note that the routine returns V**H, not V.

        NOTE: matrix 'a' will be modified

func Zgetrf(m, n int, a []complex128, lda int, ipiv []int32)
    Zgetrf computes an LU factorization of a general M-by-N matrix A using
    partial pivoting with row interchanges.

        See: http://www.netlib.org/lapack/explore-html/dd/dd1/zgetrf_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-getrf

        The factorization has the form
           A = P * L * U
        where P is a permutation matrix, L is lower triangular with unit
        diagonal elements (lower trapezoidal if m > n), and U is upper
        triangular (upper trapezoidal if m < n).

        This is the right-looking Level 3 BLAS version of the algorithm.

        NOTE: (1) matrix 'a' will be modified
              (2) ipiv indices are 1-based (i.e. Fortran)

func Zgetri(n int, a []complex128, lda int, ipiv []int32)
    Zgetri computes the inverse of a matrix using the LU factorization computed
    by Zgetrf.

        See: http://www.netlib.org/lapack/explore-html/d0/db3/zgetri_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-getri

        This method inverts U and then computes inv(A) by solving the system
        inv(A)*L = inv(U) for inv(A).

func Zherk(up, trans bool, n, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128, ldc int)
    Zherk performs one of the hermitian rank k operations

        See: http://www.netlib.org/lapack/explore-html/d1/db1/zherk_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-herk

           C := alpha*A*A**H + beta*C,

        or

           C := alpha*A**H*A + beta*C,

        where  alpha and beta  are  real scalars,  C is an  n by n  hermitian
        matrix and  A  is an  n by k  matrix in the  first case and a  k by n
        matrix in the second case.

func Zpotrf(up bool, n int, a []complex128, lda int)
    Zpotrf computes the Cholesky factorization of a complex Hermitian positive
    definite matrix A.

        See: http://www.netlib.org/lapack/explore-html/d1/db9/zpotrf_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-potrf

        The factorization has the form

           A = U**H * U,  if UPLO = 'U'

        or

           A = L  * L**H,  if UPLO = 'L'

        where U is an upper triangular matrix and L is lower triangular.

        This is the block version of the algorithm, calling Level 3 BLAS.

func Zsyrk(up, trans bool, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int)
    Zsyrk performs one of the symmetric rank k operations

        See: http://www.netlib.org/lapack/explore-html/de/d54/zsyrk_8f.html

        See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-syrk

           C := alpha*A*A**T + beta*C,

        or

           C := alpha*A**T*A + beta*C,

        where  alpha and beta  are scalars,  C is an  n by n symmetric matrix
        and  A  is an  n by k  matrix in the first case and a  k by n  matrix
        in the second case.

```
