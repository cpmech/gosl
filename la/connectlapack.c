// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdlib.h> // for malloc and free

// DGETRF - compute an LU factorization of a general M-by-N matrix A using partial pivoting with row interchanges
void dgetrf_(const int* m, const int* n, double* a, const int* lda, int* ipiv, int* info);

// DGETRI - compute the inverse of a matrix using the LU factorization computed by DGETRF
void dgetri_(const int* n, double* a, const int* lda, int* ipiv, double* work, const int* lwork, int* info);

// DGESVD - computes the singular value decomposition of a real M-by-N matrix A, optionally computing the left and/or right singular vectors
void dgesvd_(const char* jobu, const char* jobvt, const int* M, const int* N, double* A, const int* lda, double* S, double* U, const int* ldu, double* VT, const int* ldvt, double* work,const int* lwork, const int* info);

// DGESV - double-general-solver
void dgesv_(int *Np, int *NRHSp, double *A, int *LDAp, int *IPIVp, double *B, int *LDBp, int *INFOp);

// make_int converts long (a_long) to int (a_int)
// Return codes:
//   0 : no problem
//   1 : conversion failed
int make_int(int * a_int, long a_long) {
    *a_int  = (int)(a_long);
    long a2 = (long)(*a_int);
    if (a2 != a_long) return 1;
    return 0;
}

// lapack_square_inverse inverts a square matrix.
// Return codes:
//   0 : no problems
//   1 : make_int failed
//   2 : LU factorization failed
//   3 : inversion failed
int lapack_square_inverse(double *Ai, long m_long, double *A) {

    // matrix size
    int m;
    int info = make_int(&m, m_long);
    if (info != 0) return 1;

    // copy A into Ai
    int i = 0;
    for (i=0; i<m*m; ++i) Ai[i] = A[i];

    // auxiliary variable
    int * ipiv = (int*)malloc(m * sizeof(int));

    // factorization
    dgetrf_(&m,     // M
            &m,     // N
            Ai,     // double * A
            &m,     // LDA
            ipiv,   // Pivot indices
            &info); // INFO

    // clean-up and check
    if (info != 0) {
        free(ipiv);
        return 2;
    }

    // auxiliary variables
    int      NB    = 8;    // Optimal blocksize ?
    int      lwork = m*NB; // Dimension of work >= max(1,m), optimal=m*NB
    double * work  = (double*)malloc(lwork * sizeof(double)); // Work

    // inversion
    dgetri_(&m,     // N
            Ai,     // double * A
            &m,     // LDA
            ipiv,   // Pivot indices
            work,   // work
            &lwork, // dimension of work
            &info); // INFO

    // clean up
    free(ipiv);
    free(work);

    // check
    if (info != 0) return 3;
    return 0;
}

int min(int a, int b) { return (a < b ? a : b); }
int max(int a, int b) { return (a > b ? a : b); }

// lapack_svd computes the singular value decomposition: A = U_mxm * D_mxn * Vt_nxn
// Note: the output arrays must have the following sizes:
//   U  [m * m]
//   S  [min(m,n)]
//   Vt [n * n]
// Note: M matrix will be modified in this method
// Return codes:
//   0 : no problems
//   1 : make_int failed
//   2 : svd failed
int lapack_svd(double *U, double *S, double *Vt, long m_long, long n_long, double *A) {

    // matrix size
    int m, n;
    int info = make_int(&m, m_long);
    if (info != 0) return 1;
    info = make_int(&n, n_long);
    if (info != 0) return 1;

    // auxiliary variables
    char job    = 'A';
    int  min_mn = min(m, n);
    int  max_mn = max(m, n);
    int  lwork  = 2.0 * max(3 * min_mn + max_mn, 5 * min_mn);

    // auxiliary arrays
    double * work = (double*)malloc(lwork * sizeof(double));

    // decomposition
    dgesvd_(&job,   // JOBU
            &job,   // JOBVT
            &m,     // M
            &n,     // N
            A,      // A
            &m,     // LDA
            S,      // S
            U,      // U
            &m,     // LDU
            Vt,     // VT
            &n,     // LDVT
            work,   // WORK
            &lwork, // LWORK
            &info); // INFO

    // clean-up
    free(work);

    // check
    if (info != 0) {
        return 2;
    }
    return 0;
}

// lapack_pseudo_inverse inverts a non-square matrix
// Note: Ai must have the following sizes:
//   Ai [n * m]
// Return codes:
//   0 : no problems
//   1 : make_int failed
//   2 : svd failed
//   3 : pseudo inverse failed
int lapack_pseudo_inverse(double *Ai, long m_long, long n_long, double *A, double tol) {

    // matrix size
    int m, n;
    int info = make_int(&m, m_long);
    if (info != 0) return 1;
    info = make_int(&n, n_long);
    if (info != 0) return 1;

    // auxiliary variables
    int ns = min(m, n);
    double * U  = (double*)malloc(m * m * sizeof(double));
    double * S  = (double*)malloc(ns * sizeof(double));
    double * Vt = (double*)malloc(n * n * sizeof(double));

    // perform singular value decomposition
    info = lapack_svd(U, S, Vt, m, n, A);

    // clean-up and check
    if (info != 0) {
        free(U);
        free(S);
        free(Vt);
        return 2;
    }

    // compute inverse
    int i = 0;
    int j = 0;
    int k = 0;
    for (i=0; i<n; ++i) {
        for (j=0; j<m; ++j) {
            Ai[i+j*n] = 0.0;
            for (k=0; k<ns; ++k) {
                if (S[k] > tol) {
                    Ai[i+j*n] += Vt[k+i*n] * (1.0 / S[k]) * U[j+k*m];
                }
            }
        }
    }

    // clean-up
    free(U);
    free(S);
    free(Vt);
    return 0;
}
