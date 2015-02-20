// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CONNECTLAPACK_H
#define CONNECTLAPACK_H

int lapack_square_inverse(double *Ai, long m, double *A);
int lapack_pseudo_inverse(double *Ai, long m, long n, double *A, double tol);
int lapack_svd(double *U, double *S, double *Vt, long m_long, long n_long, double *A);

#endif // CONNECTLAPACK_H
