// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CONNECTMPI_H
#define CONNECTMPI_H

#ifdef __cplusplus
extern "C" {
#endif

#include <complex.h>

void abortmpi();
int  ison();
void startmpi(int debug);
void stopmpi(int debug);
int  mpirank();
int  mpisize();
void barrier();
void sumtoroot(double *dest, double *orig, int n);
void sumtorootC(double complex *dest, double complex *orig, int n);
void bcastfromroot(double *x, int n);
void bcastfromrootC(double complex *x, int n);
void allreducesum(double *dest, double *orig, int n);
void allreducemin(double *dest, double *orig, int n);
void allreducemax(double *dest, double *orig, int n);
void intallreducemax(int *dest, int *orig, int n);
void intsend(int *vals, int n, int to_proc);
void intrecv(int *vals, int n, int from_proc);
void dblsend(double *vals, int n, int to_proc);
void dblrecv(double *vals, int n, int from_proc);

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif // CONNECTMPI_H
