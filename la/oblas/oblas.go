// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package oblas implements lower level linear algebra using OpenBLAS for max efficiency
package oblas

/*
#cgo linux   CFLAGS: -DOPENBLAS_USE64BITINT -O2 -I/usr/include
#cgo linux   CFLAGS: -DOPENBLAS_USE64BITINT -O2 -I/usr/local/include
#cgo windows CFLAGS: -DOPENBLAS_USE64BITINT -O2 -IC:/Gosl/include
#cgo linux   LDFLAGS: -lopenblas -L/local/lib
#cgo darwin  LDFLAGS: -lopenblas -L/usr/local/lib
#cgo windows LDFLAGS: -lopenblas -LC:/Gosl/lib
#ifdef WIN32
#define LONG long long
#else
#define LONG long
#endif

#include <cblas.h>
#include <lapacke.h>
static inline double* cpt(double complex* p) { return (double*)p; }
*/
import "C"

import (
	"unsafe"

	"github.com/cpmech/gosl/chk"
)

// Daxpy computes constant times a vector plus a vector.
//  See: http://www.netlib.org/lapack/explore-html/d9/dcd/daxpy_8f.html
//
//  y += alpha*x + y
//
func Daxpy(n int, alpha float64, x []float64, incx int, y []float64, incy int) (err error) {
	nmin := imin(len(x), len(y))
	if n > nmin {
		return chk.Err("n must not be greater than %d", n, nmin)
	}
	C.cblas_daxpy(
		C.blasint(n),
		C.double(alpha),
		(*C.double)(unsafe.Pointer(&x[0])),
		C.blasint(incx),
		(*C.double)(unsafe.Pointer(&y[0])),
		C.blasint(incy),
	)
	return
}

// Zaxpy computes constant times a vector plus a vector.
//  See: http://www.netlib.org/lapack/explore-html/d7/db2/zaxpy_8f.html
//
//  y += alpha*x + y
//
func Zaxpy(n int, alpha complex128, x []complex128, incx int, y []complex128, incy int) (err error) {
	nmin := imin(len(x), len(y))
	if n > nmin {
		return chk.Err("n must not be greater than %d", n, nmin)
	}
	C.cblas_zaxpy(
		C.blasint(n),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&alpha))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&x[0]))),
		C.blasint(incx),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&y[0]))),
		C.blasint(incy),
	)
	return
}

// Dgemv performs one of the matrix-vector operations
//  See: http://www.netlib.org/lapack/explore-html/dc/da8/dgemv_8f.html
//
//     y := alpha*A*x + beta*y,   or   y := alpha*A**T*x + beta*y,
//
//  where alpha and beta are scalars, x and y are vectors and A is an
//  m by n matrix.
//     trans=false     y := alpha*A*x + beta*y.
//
//     trans=true      y := alpha*A**T*x + beta*y.
func Dgemv(trans bool, m, n int, alpha float64, a *Matrix, lda int, x []float64, incx int, beta float64, y []float64, incy int) (err error) {
	if trans {
		if len(x) != a.m {
			return chk.Err("len(x)=%d must be equal to m=%d", len(x), a.m)
		}
		if len(y) != a.n {
			return chk.Err("len(y)=%d must be equal to n=%d", len(y), a.n)
		}
	} else {
		if len(x) != a.n {
			return chk.Err("len(x)=%d must be equal to n=%d", len(x), a.n)
		}
		if len(y) != a.m {
			return chk.Err("len(y)=%d must be equal to m=%d", len(y), a.m)
		}
	}
	C.cblas_dgemv(
		cblasColMajor,
		cTrans(trans),
		C.blasint(m),
		C.blasint(n),
		C.double(alpha),
		(*C.double)(unsafe.Pointer(&a.data[0])),
		C.blasint(lda),
		(*C.double)(unsafe.Pointer(&x[0])),
		C.blasint(incx),
		C.double(beta),
		(*C.double)(unsafe.Pointer(&y[0])),
		C.blasint(incy),
	)
	return
}

// Zgemv performs one of the matrix-vector operations.
//  See: http://www.netlib.org/lapack/explore-html/db/d40/zgemv_8f.html
//
//     y := alpha*A*x + beta*y,   or   y := alpha*A**T*x + beta*y,   or
//
//     y := alpha*A**H*x + beta*y,
//
//  where alpha and beta are scalars, x and y are vectors and A is an
//  m by n matrix.
func Zgemv(trans bool, m, n int, alpha complex128, a *MatrixC, lda int, x []complex128, incx int, beta complex128, y []complex128, incy int) (err error) {
	if trans {
		if len(x) != a.m {
			return chk.Err("len(x)=%d must be equal to m=%d", len(x), a.m)
		}
		if len(y) != a.n {
			return chk.Err("len(y)=%d must be equal to n=%d", len(y), a.n)
		}
	} else {
		if len(x) != a.n {
			return chk.Err("len(x)=%d must be equal to n=%d", len(x), a.n)
		}
		if len(y) != a.m {
			return chk.Err("len(y)=%d must be equal to m=%d", len(y), a.m)
		}
	}
	C.cblas_zgemv(
		cblasColMajor,
		cTrans(trans),
		C.blasint(m),
		C.blasint(n),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&alpha))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&a.data[0]))),
		C.blasint(lda),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&x[0]))),
		C.blasint(incx),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&beta))),
		C.cpt((*C.complexdouble)(unsafe.Pointer(&y[0]))),
		C.blasint(incy),
	)
	return
}

// Dgesv computes the solution to a real system of linear equations.
//  See: http://www.netlib.org/lapack/explore-html/d8/d72/dgesv_8f.html
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
func Dgesv(n, nrhs int, a *Matrix, lda int, ipiv []int32, b []float64, ldb int) (err error) {
	if len(ipiv) != n {
		return chk.Err("len(ipiv) must be equal to n. %d != %d\n", len(ipiv), n)
	}
	info := C.LAPACKE_dgesv(
		C.int(lapackColMajor),
		C.lapack_int(n),
		C.lapack_int(nrhs),
		(*C.double)(unsafe.Pointer(&a.data[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
		(*C.double)(unsafe.Pointer(&b[0])),
		C.lapack_int(ldb),
	)
	if info != 0 {
		err = chk.Err("lapack failed\n")
	}
	return
}

// Zgesv computes the solution to a complex system of linear equations.
//  See: http://www.netlib.org/lapack/explore-html/d1/ddc/zgesv_8f.html
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
func Zgesv(n, nrhs int, a *MatrixC, lda int, ipiv []int32, b []complex128, ldb int) (err error) {
	if len(ipiv) != n {
		return chk.Err("len(ipiv) must be equal to n. %d != %d\n", len(ipiv), n)
	}
	info := C.LAPACKE_zgesv(
		C.int(lapackColMajor),
		C.lapack_int(n),
		C.lapack_int(nrhs),
		(*C.lapack_complex_double)(unsafe.Pointer(&a.data[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
		(*C.lapack_complex_double)(unsafe.Pointer(&b[0])),
		C.lapack_int(ldb),
	)
	if info != 0 {
		err = chk.Err("lapack failed\n")
	}
	return
}

// Dgesvd computes the singular value decomposition (SVD) of a real M-by-N matrix A, optionally computing the left and/or right singular vectors.
//  See: http://www.netlib.org/lapack/explore-html/d8/d2d/dgesvd_8f.html
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
func Dgesvd(jobu, jobvt rune, m, n int, a *Matrix, lda int, s []float64, u *Matrix, ldu int, vt *Matrix, ldvt int, work []float64, lwork int) (err error) {
	info := C.LAPACKE_dgesvd_work(
		C.int(lapackColMajor),
		C.char(jobu),
		C.char(jobvt),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a.data[0])),
		C.lapack_int(lda),
		(*C.double)(unsafe.Pointer(&s[0])),
		(*C.double)(unsafe.Pointer(&u.data[0])),
		C.lapack_int(ldu),
		(*C.double)(unsafe.Pointer(&vt.data[0])),
		C.lapack_int(ldvt),
		(*C.double)(unsafe.Pointer(&work[0])),
		C.lapack_int(lwork),
	)
	if info != 0 {
		err = chk.Err("lapack failed\n")
	}
	return
}

// Zgesvd computes the singular value decomposition (SVD) of a complex M-by-N matrix A, optionally computing the left and/or right singular vectors.
//  See: http://www.netlib.org/lapack/explore-html/d6/d42/zgesvd_8f.html
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
func Zgesvd(jobu, jobvt rune, m, n int, a *MatrixC, lda int, s []float64, u *MatrixC, ldu int, vt *MatrixC, ldvt int, work []complex128, lwork int, rwork []float64) (err error) {
	info := C.LAPACKE_zgesvd_work(
		C.int(lapackColMajor),
		C.char(jobu),
		C.char(jobvt),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.lapack_complex_double)(unsafe.Pointer(&a.data[0])),
		C.lapack_int(lda),
		(*C.double)(unsafe.Pointer(&s[0])),
		(*C.lapack_complex_double)(unsafe.Pointer(&u.data[0])),
		C.lapack_int(ldu),
		(*C.lapack_complex_double)(unsafe.Pointer(&vt.data[0])),
		C.lapack_int(ldvt),
		(*C.lapack_complex_double)(unsafe.Pointer(&work[0])),
		C.lapack_int(lwork),
		(*C.double)(unsafe.Pointer(&rwork[0])),
	)
	if info != 0 {
		err = chk.Err("lapack failed\n")
	}
	return
}

// Dgetrf computes an LU factorization of a general M-by-N matrix A using partial pivoting with row interchanges.
//  See: http://www.netlib.org/lapack/explore-html/d3/d6a/dgetrf_8f.html
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
func Dgetrf(m, n int, a *Matrix, lda int, ipiv []int32) (err error) {
	info := C.LAPACKE_dgetrf_work(
		C.int(lapackColMajor),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.double)(unsafe.Pointer(&a.data[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
	)
	if info != 0 {
		err = chk.Err("lapack failed\n")
	}
	return
}

// Zgetrf computes an LU factorization of a general M-by-N matrix A using partial pivoting with row interchanges.
//  See: http://www.netlib.org/lapack/explore-html/dd/dd1/zgetrf_8f.html
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
func Zgetrf(m, n int, a *MatrixC, lda int, ipiv []int32) (err error) {
	info := C.LAPACKE_zgetrf_work(
		C.int(lapackColMajor),
		C.lapack_int(m),
		C.lapack_int(n),
		(*C.lapack_complex_double)(unsafe.Pointer(&a.data[0])),
		C.lapack_int(lda),
		(*C.lapack_int)(unsafe.Pointer(&ipiv[0])),
	)
	if info != 0 {
		err = chk.Err("lapack failed\n")
	}
	return
}

// Dgetri computes the inverse of a matrix using the LU factorization computed by DGETRF.
//  See: http://www.netlib.org/lapack/explore-html/df/da4/dgetri_8f.html
//
//  This method inverts U and then computes inv(A) by solving the system
//  inv(A)*L = inv(U) for inv(A).
func Dgetri(n int, a []float64, lda int, ipiv []int32, work []float64, lwork int) (err error) {
	chk.Panic("TODO: Dgetri")
	return
}

// Zgetri computes the inverse of a matrix using the LU factorization computed by Zgetrf.
//  See: http://www.netlib.org/lapack/explore-html/d0/db3/zgetri_8f.html
//
//  This method inverts U and then computes inv(A) by solving the system
//  inv(A)*L = inv(U) for inv(A).
func Zgetri(n int, a []complex128, lda int, ipiv []int32, work []complex128, lwork int) (err error) {
	chk.Panic("TODO: Zgetri")
	return
}

// Dsyrk performs one of the symmetric rank k operations
//  See: http://www.netlib.org/lapack/explore-html/dc/d05/dsyrk_8f.html
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
func Dsyrk(up, trans bool, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) (err error) {
	chk.Panic("TODO: Dsyrk")
	return
}

// Zsyrk performs one of the symmetric rank k operations
//  See: http://www.netlib.org/lapack/explore-html/de/d54/zsyrk_8f.html
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
func Zsyrk(up, trans bool, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int) (err error) {
	chk.Panic("TODO: Zsyrk")
	return
}

// Zherk performs one of the hermitian rank k operations
//  See: http://www.netlib.org/lapack/explore-html/d1/db1/zherk_8f.html
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
func Zherk(up, trans bool, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int) (err error) {
	chk.Panic("TODO: Zherk")
	return
}

// Dpotrf computes the Cholesky factorization of a real symmetric positive definite matrix A.
//  See: http://www.netlib.org/lapack/explore-html/d0/d8a/dpotrf_8f.html
//
//  The factorization has the form
//     A = U**T * U,  if UPLO = 'U', or
//     A = L  * L**T,  if UPLO = 'L',
//  where U is an upper triangular matrix and L is lower triangular.
//
//  This is the block version of the algorithm, calling Level 3 BLAS.
func Dpotrf(up bool, n int, a []float64, lda int) (err error) {
	chk.Panic("TODO: Dpotrf")
	return
}

// Zpotrf computes the Cholesky factorization of a complex Hermitian positive definite matrix A.
//  See: http://www.netlib.org/lapack/explore-html/d1/db9/zpotrf_8f.html
//
//  The factorization has the form
//     A = U**H * U,  if UPLO = 'U', or
//     A = L  * L**H,  if UPLO = 'L',
//  where U is an upper triangular matrix and L is lower triangular.
//
//  This is the block version of the algorithm, calling Level 3 BLAS.
func Zpotrf(up bool, n int, a []complex128, lda int) (err error) {
	chk.Panic("TODO: Zpotrf")
	return
}

// Dcholesky performs the Cholesky factorization
func Dcholesky() (err error) {
	chk.Panic("TODO: Dcholesky")
	return
}

// auxiliary //////////////////////////////////////////////////////////////////////////////////////

// constants
const (
	// Lapack matrix layout
	lapackRowMajor int = 101
	lapackColMajor int = 102

	// CBLAS_ORDER;
	cblasRowMajor uint32 = 101
	cblasColMajor uint32 = 102

	// CBLAS_TRANSPOSE;
	cblasNoTrans     uint32 = 111
	cblasTrans       uint32 = 112
	cblasConjTrans   uint32 = 113
	cblasConjNoTrans uint32 = 114

	// CBLAS_UPLO;
	cblasUpper uint32 = 121
	cblasLower uint32 = 122

	// CBLAS_DIAG;
	cblasNonUnit uint32 = 131
	cblasUnit    uint32 = 132

	// CBLAS_SIDE;
	cblasLeft  uint32 = 141
	cblasRight uint32 = 142
)

func cTrans(trans bool) uint32 {
	if trans {
		return cblasTrans
	}
	return cblasNoTrans
}
