// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mkl implements lower-level linear algebra routines using MKL
// for maximum efficiency. This package uses column-major representation for matrices.
//
//   Example of col-major data:
//             _      _
//            |  0  3  |
//        A = |  1  4  |            ⇒     a = [0, 1, 2, 3, 4, 5]
//            |_ 2  5 _|(m x n)
//
//        a[i+j*m] = A[i][j]
//
//  NOTE: the functions here do not check for the limits of indices. Be careful.
//        Panic may occur then.
//
package mkl

/*
#include <mkl.h>

#include <complex.h>
static inline void* cpt(double complex* p) { return (void*)p; }
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

// SetNumThreads sets the number of threads in OpenBLAS
func SetNumThreads(n int) {
	C.mkl_set_num_threads(C.int(n))
}

// Ddot forms the dot product of two vectors. uses unrolled loops for increments equal to one.
//
//  See: http://www.netlib.org/lapack/explore-html/d5/df6/ddot_8f.html
func Ddot(n int, x []float64, incx int, y []float64, incy int) (res float64) {
	cres := C.cblas_ddot(
		C.MKL_INT(n),
		(*C.double)(unsafe.Pointer(&x[0])),
		C.MKL_INT(incx),
		(*C.double)(unsafe.Pointer(&y[0])),
		C.MKL_INT(incy),
	)
	return float64(cres)
}

// Daxpy computes constant times a vector plus a vector.
//
//  See: http://www.netlib.org/lapack/explore-html/d9/dcd/daxpy_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-axpy
//
//  y += alpha*x + y
//
func Daxpy(n int, alpha float64, x []float64, incx int, y []float64, incy int) {
	C.cblas_daxpy(
		C.MKL_INT(n),
		C.double(alpha),
		(*C.double)(unsafe.Pointer(&x[0])),
		C.MKL_INT(incx),
		(*C.double)(unsafe.Pointer(&y[0])),
		C.MKL_INT(incy),
	)
}

// Zaxpy computes constant times a vector plus a vector.
//
//  See: http://www.netlib.org/lapack/explore-html/d7/db2/zaxpy_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-axpy
//
//  y += alpha*x + y
//
func Zaxpy(n int, alpha complex128, x []complex128, incx int, y []complex128, incy int) {
	C.cblas_zaxpy(
		C.MKL_INT(n),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&alpha))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&x[0]))),
		C.MKL_INT(incx),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&y[0]))),
		C.MKL_INT(incy),
	)
}

// Dgemv performs one of the matrix-vector operations
//
//  See: http://www.netlib.org/lapack/explore-html/dc/da8/dgemv_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemv
//
//     y := alpha*A*x + beta*y,   or   y := alpha*A**T*x + beta*y,
//
//  where alpha and beta are scalars, x and y are vectors and A is an
//  m by n matrix.
//     trans=false     y := alpha*A*x + beta*y.
//
//     trans=true      y := alpha*A**T*x + beta*y.
func Dgemv(trans bool, m, n int, alpha float64, a []float64, lda int, x []float64, incx int, beta float64, y []float64, incy int) {
	C.cblas_dgemv(
		cblasColMajor,
		cTrans(trans),
		C.MKL_INT(m),
		C.MKL_INT(n),
		C.double(alpha),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.MKL_INT(lda),
		(*C.double)(unsafe.Pointer(&x[0])),
		C.MKL_INT(incx),
		C.double(beta),
		(*C.double)(unsafe.Pointer(&y[0])),
		C.MKL_INT(incy),
	)
}

// Zgemv performs one of the matrix-vector operations.
//
//  See: http://www.netlib.org/lapack/explore-html/db/d40/zgemv_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemv
//
//     y := alpha*A*x + beta*y,   or   y := alpha*A**T*x + beta*y,   or
//
//     y := alpha*A**H*x + beta*y,
//
//  where alpha and beta are scalars, x and y are vectors and A is an
//  m by n matrix.
func Zgemv(trans bool, m, n int, alpha complex128, a []complex128, lda int, x []complex128, incx int, beta complex128, y []complex128, incy int) {
	C.cblas_zgemv(
		cblasColMajor,
		cTrans(trans),
		C.MKL_INT(m),
		C.MKL_INT(n),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&alpha))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&a[0]))),
		C.MKL_INT(lda),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&x[0]))),
		C.MKL_INT(incx),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&beta))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&y[0]))),
		C.MKL_INT(incy),
	)
}

// Dgemm performs one of the matrix-matrix operations
//
//  false,false:  C_{m,n} := α ⋅ A_{m,k} ⋅ B_{k,n}  +  β ⋅ C_{m,n}
//  false,true:   C_{m,n} := α ⋅ A_{m,k} ⋅ B_{n,k}  +  β ⋅ C_{m,n}
//  true, false:  C_{m,n} := α ⋅ A_{k,m} ⋅ B_{k,n}  +  β ⋅ C_{m,n}
//  true, true:   C_{m,n} := α ⋅ A_{k,m} ⋅ B_{n,k}  +  β ⋅ C_{m,n}
//
//  see: http://www.netlib.org/lapack/explore-html/d7/d2b/dgemm_8f.html
//
//  see: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemm
//
//     C := alpha*op( A )*op( B ) + beta*C,
//
//  where  op( X ) is one of
//
//     op( X ) = X   or   op( X ) = X**T,
//
//  alpha and beta are scalars, and A, B and C are matrices, with op( A )
//  an m by k matrix,  op( B )  a  k by n matrix and  C an m by n matrix.
func Dgemm(transA, transB bool, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	C.cblas_dgemm(
		cblasColMajor,
		cTrans(transA),
		cTrans(transB),
		C.MKL_INT(m),
		C.MKL_INT(n),
		C.MKL_INT(k),
		C.double(alpha),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.MKL_INT(lda),
		(*C.double)(unsafe.Pointer(&b[0])),
		C.MKL_INT(ldb),
		C.double(beta),
		(*C.double)(unsafe.Pointer(&c[0])),
		C.MKL_INT(ldc),
	)
}

// Zgemm performs one of the matrix-matrix operations
//
//  see: http://www.netlib.org/lapack/explore-html/d7/d76/zgemm_8f.html
//
//  see: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-gemm
//
//     C := alpha*op( A )*op( B ) + beta*C,
//
//  where  op( X ) is one of
//
//     op( X ) = X   or   op( X ) = X**T   or   op( X ) = X**H,
//
//  alpha and beta are scalars, and A, B and C are matrices, with op( A )
//  an m by k matrix,  op( B )  a  k by n matrix and  C an m by n matrix.
func Zgemm(transA, transB bool, m, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int) {
	C.cblas_zgemm(
		cblasColMajor,
		cTrans(transA),
		cTrans(transB),
		C.MKL_INT(m),
		C.MKL_INT(n),
		C.MKL_INT(k),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&alpha))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&a[0]))),
		C.MKL_INT(lda),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&b[0]))),
		C.MKL_INT(ldb),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&beta))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&c[0]))),
		C.MKL_INT(ldc),
	)
}

// Dgesv computes the solution to a real system of linear equations.
//
//  See: http://www.netlib.org/lapack/explore-html/d8/d72/dgesv_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-gesv
//
//  The system is:
//
//     A * X = B,
//
//  where A is an N-by-N matrix and X and B are N-by-NRHS matrices.
//
//  The LU decomposition with partial pivoting and row interchanges is
//  used to factor A as
//
//     A = P * L * U,
//
//  where P is a permutation matrix, L is unit lower triangular, and U is
//  upper triangular.  The factored form of A is then used to solve the
//  system of equations A * X = B.
//
//  NOTE: matrix 'a' will be modified
func Dgesv(n, nrhs int, a []float64, lda int, ipiv []int64, b []float64, ldb int) {
	if len(ipiv) != n {
		chk.Panic("len(ipiv) must be equal to n. %d != %d\n", len(ipiv), n)
	}
	info := C.LAPACKE_dgesv(
		C.int(lapackColMajor),
		C.lapack_int(n),
		C.lapack_int(nrhs),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
		(*C.double)(unsafe.Pointer(&b[0])),
		C.lapack_int(ldb),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Zgesv computes the solution to a complex system of linear equations.
//
//  See: http://www.netlib.org/lapack/explore-html/d1/ddc/zgesv_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-gesv
//
//  The system is:
//
//     A * X = B,
//
//  where A is an N-by-N matrix and X and B are N-by-NRHS matrices.
//
//  The LU decomposition with partial pivoting and row interchanges is
//  used to factor A as
//
//     A = P * L * U,
//
//  where P is a permutation matrix, L is unit lower triangular, and U is
//  upper triangular.  The factored form of A is then used to solve the
//  system of equations A * X = B.
//
//  NOTE: matrix 'a' will be modified
func Zgesv(n, nrhs int, a []complex128, lda int, ipiv []int64, b []complex128, ldb int) {
	if len(ipiv) != n {
		chk.Panic("len(ipiv) must be equal to n. %d != %d\n", len(ipiv), n)
	}
	info := C.LAPACKE_zgesv(
		C.int(lapackColMajor),
		C.lapack_int(n),
		C.lapack_int(nrhs),
		(*C.lapack_complex_double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
		(*C.lapack_complex_double)(unsafe.Pointer(&b[0])),
		C.lapack_int(ldb),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Dgesvd computes the singular value decomposition (SVD) of a real M-by-N matrix A, optionally computing the left and/or right singular vectors.
//
//  See: http://www.netlib.org/lapack/explore-html/d8/d2d/dgesvd_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-gesvd
//
//  The SVD is written
//
//       A = U * SIGMA * transpose(V)
//
//  where SIGMA is an M-by-N matrix which is zero except for its
//  min(m,n) diagonal elements, U is an M-by-M orthogonal matrix, and
//  V is an N-by-N orthogonal matrix.  The diagonal elements of SIGMA
//  are the singular values of A; they are real and non-negative, and
//  are returned in descending order.  The first min(m,n) columns of
//  U and V are the left and right singular vectors of A.
//
//  Note that the routine returns V**T, not V.
//
//  NOTE: matrix 'a' will be modified
func Dgesvd(jobu, jobvt rune, m, n int, a []float64, lda int, s []float64, u []float64, ldu int, vt []float64, ldvt int, superb []float64) {
	info := C.LAPACKE_dgesvd(
		C.int(lapackColMajor),
		C.char(jobu),
		C.char(jobvt),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.double)(unsafe.Pointer(&s[0])),
		(*C.double)(unsafe.Pointer(&u[0])),
		C.lapack_int(ldu),
		(*C.double)(unsafe.Pointer(&vt[0])),
		C.lapack_int(ldvt),
		(*C.double)(unsafe.Pointer(&superb[0])),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Zgesvd computes the singular value decomposition (SVD) of a complex M-by-N matrix A, optionally computing the left and/or right singular vectors.
//
//  See: http://www.netlib.org/lapack/explore-html/d6/d42/zgesvd_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-gesvd
//
//  The SVD is written
//
//       A = U * SIGMA * conjugate-transpose(V)
//
//  where SIGMA is an M-by-N matrix which is zero except for its
//  min(m,n) diagonal elements, U is an M-by-M unitary matrix, and
//  V is an N-by-N unitary matrix.  The diagonal elements of SIGMA
//  are the singular values of A; they are real and non-negative, and
//  are returned in descending order.  The first min(m,n) columns of
//  U and V are the left and right singular vectors of A.
//
//  Note that the routine returns V**H, not V.
//
//  NOTE: matrix 'a' will be modified
func Zgesvd(jobu, jobvt rune, m, n int, a []complex128, lda int, s []float64, u []complex128, ldu int, vt []complex128, ldvt int, superb []float64) {
	info := C.LAPACKE_zgesvd(
		C.int(lapackColMajor),
		C.char(jobu),
		C.char(jobvt),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.lapack_complex_double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.double)(unsafe.Pointer(&s[0])),
		(*C.lapack_complex_double)(unsafe.Pointer(&u[0])),
		C.lapack_int(ldu),
		(*C.lapack_complex_double)(unsafe.Pointer(&vt[0])),
		C.lapack_int(ldvt),
		(*C.double)(unsafe.Pointer(&superb[0])),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Dgetrf computes an LU factorization of a general M-by-N matrix A using partial pivoting with row interchanges.
//
//  See: http://www.netlib.org/lapack/explore-html/d3/d6a/dgetrf_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-getrf
//
//  The factorization has the form
//     A = P * L * U
//  where P is a permutation matrix, L is lower triangular with unit
//  diagonal elements (lower trapezoidal if m > n), and U is upper
//  triangular (upper trapezoidal if m < n).
//
//  This is the right-looking Level 3 BLAS version of the algorithm.
//
//  NOTE: (1) matrix 'a' will be modified
//        (2) ipiv indices are 1-based (i.e. Fortran)
func Dgetrf(m, n int, a []float64, lda int, ipiv []int64) {
	info := C.LAPACKE_dgetrf(
		C.int(lapackColMajor),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Zgetrf computes an LU factorization of a general M-by-N matrix A using partial pivoting with row interchanges.
//
//  See: http://www.netlib.org/lapack/explore-html/dd/dd1/zgetrf_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-getrf
//
//  The factorization has the form
//     A = P * L * U
//  where P is a permutation matrix, L is lower triangular with unit
//  diagonal elements (lower trapezoidal if m > n), and U is upper
//  triangular (upper trapezoidal if m < n).
//
//  This is the right-looking Level 3 BLAS version of the algorithm.
//
//  NOTE: (1) matrix 'a' will be modified
//        (2) ipiv indices are 1-based (i.e. Fortran)
func Zgetrf(m, n int, a []complex128, lda int, ipiv []int64) {
	info := C.LAPACKE_zgetrf(
		C.int(lapackColMajor),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.lapack_complex_double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Dgetri computes the inverse of a matrix using the LU factorization computed by DGETRF.
//
//  See: http://www.netlib.org/lapack/explore-html/df/da4/dgetri_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-getri
//
//  This method inverts U and then computes inv(A) by solving the system
//  inv(A)*L = inv(U) for inv(A).
func Dgetri(n int, a []float64, lda int, ipiv []int64) {
	info := C.LAPACKE_dgetri(
		C.int(lapackColMajor),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Zgetri computes the inverse of a matrix using the LU factorization computed by Zgetrf.
//
//  See: http://www.netlib.org/lapack/explore-html/d0/db3/zgetri_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-getri
//
//  This method inverts U and then computes inv(A) by solving the system
//  inv(A)*L = inv(U) for inv(A).
func Zgetri(n int, a []complex128, lda int, ipiv []int64) {
	info := C.LAPACKE_zgetri(
		C.int(lapackColMajor),
		C.lapack_int(n),
		(*C.lapack_complex_double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Dsyrk performs one of the symmetric rank k operations
//
//  See: http://www.netlib.org/lapack/explore-html/dc/d05/dsyrk_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-syrk
//
//     C := alpha*A*A**T + beta*C,
//
//  or
//
//     C := alpha*A**T*A + beta*C,
//
//  where  alpha and beta  are scalars, C is an  n by n  symmetric matrix
//  and  A  is an  n by k  matrix in the first case and a  k by n  matrix
//  in the second case.
func Dsyrk(up, trans bool, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) {
	C.cblas_dsyrk(
		cblasColMajor,
		cUplo(up),
		cTrans(trans),
		C.MKL_INT(n),
		C.MKL_INT(k),
		C.double(alpha),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.MKL_INT(lda),
		C.double(beta),
		(*C.double)(unsafe.Pointer(&c[0])),
		C.MKL_INT(ldc),
	)
}

// Zsyrk performs one of the symmetric rank k operations
//
//  See: http://www.netlib.org/lapack/explore-html/de/d54/zsyrk_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-syrk
//
//     C := alpha*A*A**T + beta*C,
//
//  or
//
//     C := alpha*A**T*A + beta*C,
//
//  where  alpha and beta  are scalars,  C is an  n by n symmetric matrix
//  and  A  is an  n by k  matrix in the first case and a  k by n  matrix
//  in the second case.
func Zsyrk(up, trans bool, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int) {
	C.cblas_zsyrk(
		cblasColMajor,
		cUplo(up),
		cTrans(trans),
		C.MKL_INT(n),
		C.MKL_INT(k),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&alpha))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&a[0]))),
		C.MKL_INT(lda),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&beta))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&c[0]))),
		C.MKL_INT(ldc),
	)
}

// Zherk performs one of the hermitian rank k operations
//
//  See: http://www.netlib.org/lapack/explore-html/d1/db1/zherk_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-cblas-herk
//
//     C := alpha*A*A**H + beta*C,
//
//  or
//
//     C := alpha*A**H*A + beta*C,
//
//  where  alpha and beta  are  real scalars,  C is an  n by n  hermitian
//  matrix and  A  is an  n by k  matrix in the  first case and a  k by n
//  matrix in the second case.
func Zherk(up, trans bool, n, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128, ldc int) {
	C.cblas_zherk(
		cblasColMajor,
		cUplo(up),
		cTrans(trans),
		C.MKL_INT(n),
		C.MKL_INT(k),
		C.double(alpha),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&a[0]))),
		C.MKL_INT(lda),
		C.double(beta),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&c[0]))),
		C.MKL_INT(ldc),
	)
}

// Dpotrf computes the Cholesky factorization of a real symmetric positive definite matrix A.
//
//  See: http://www.netlib.org/lapack/explore-html/d0/d8a/dpotrf_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-potrf
//
//  The factorization has the form
//
//     A = U**T * U,  if UPLO = 'U'
//
//  or
//
//     A = L  * L**T,  if UPLO = 'L'
//
//  where U is an upper triangular matrix and L is lower triangular.
//
//  This is the block version of the algorithm, calling Level 3 BLAS.
func Dpotrf(up bool, n int, a []float64, lda int) {
	info := C.LAPACKE_dpotrf(
		C.int(lapackColMajor),
		lUplo(up),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Zpotrf computes the Cholesky factorization of a complex Hermitian positive definite matrix A.
//
//  See: http://www.netlib.org/lapack/explore-html/d1/db9/zpotrf_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-potrf
//
//  The factorization has the form
//
//     A = U**H * U,  if UPLO = 'U'
//
//  or
//
//     A = L  * L**H,  if UPLO = 'L'
//
//  where U is an upper triangular matrix and L is lower triangular.
//
//  This is the block version of the algorithm, calling Level 3 BLAS.
func Zpotrf(up bool, n int, a []complex128, lda int) {
	info := C.LAPACKE_zpotrf(
		C.int(lapackColMajor),
		lUplo(up),
		C.lapack_int(n),
		(*C.lapack_complex_double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// Dgeev computes for an N-by-N real nonsymmetric matrix A, the
// eigenvalues and, optionally, the left and/or right eigenvectors.
//
//  See: http://www.netlib.org/lapack/explore-html/d9/d28/dgeev_8f.html
//
//  See: https://software.intel.com/en-us/mkl-developer-reference-c-geev
//
//  See: https://www.nag.co.uk/numeric/fl/nagdoc_fl26/html/f08/f08naf.html
//
//  The right eigenvector v(j) of A satisfies
//
//                   A * v(j) = lambda(j) * v(j)
//
//  where lambda(j) is its eigenvalue.
//
//  The left eigenvector u(j) of A satisfies
//
//                u(j)**H * A = lambda(j) * u(j)**H
//
//  where u(j)**H denotes the conjugate-transpose of u(j).
//
//  The computed eigenvectors are normalized to have Euclidean norm
//  equal to 1 and largest component real.
func Dgeev(calcVl, calcVr bool, n int, a []float64, lda int, wr []float64, wi, vl []float64, ldvl int, vr []float64, ldvr int) {
	var vvl, vvr *C.double
	if calcVl {
		vvl = (*C.double)(unsafe.Pointer(&vl[0]))
	} else {
		ldvl = 1
	}
	if calcVr {
		vvr = (*C.double)(unsafe.Pointer(&vr[0]))
	} else {
		ldvr = 1
	}
	info := C.LAPACKE_dgeev(
		C.int(lapackColMajor),
		jobVlr(calcVl),
		jobVlr(calcVr),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a[0])),
		C.lapack_int(lda),
		(*C.double)(unsafe.Pointer(&wr[0])),
		(*C.double)(unsafe.Pointer(&wi[0])),
		vvl,
		C.lapack_int(ldvl),
		vvr,
		C.lapack_int(ldvr),
	)
	if info != 0 {
		chk.Panic("lapack failed\n")
	}
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// constants
const (
	// Lapack matrix layout
	lapackRowMajor = 101
	lapackColMajor = 102

	// CBLAS_LAYOUT == CBLAS_ORDER
	cblasRowMajor = 101
	cblasColMajor = 102

	// CBLAS_TRANSPOSE
	cblasNoTrans   = 111
	cblasTrans     = 112
	cblasConjTrans = 113

	// CBLAS_UPLO
	cblasUpper = 121
	cblasLower = 122

	// CBLAS_DIAG
	cblasNonUnit = 131
	cblasUnit    = 132

	// CBLAS_SIDE
	cblasLeft  = 141
	cblasRight = 142

	// CBLAS_STORAGE
	cblasPacked = 151

	// CBLAS_IDENTIFIER
	cblasAMatrix = 161
	cblasBMatrix = 162
)

func cTrans(trans bool) C.CBLAS_TRANSPOSE {
	if trans {
		return cblasTrans
	}
	return cblasNoTrans
}

func cUplo(up bool) C.CBLAS_UPLO {
	if up {
		return cblasUpper
	}
	return cblasLower
}

func lUplo(up bool) C.char {
	if up {
		return 'U'
	}
	return 'L'
}

func jobVlr(doCalc bool) C.char {
	if doCalc {
		return 'V'
	}
	return 'N'
}
