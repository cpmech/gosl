# include <complex.h>
# include <stdlib.h>
# include <stdio.h>
# include <math.h>
# include <string.h>

# include"test_mat.h"

int main ( );
void bvec_next_grlex_test ( );
void legendre_zeros_test ( );
void mertens_test ( );
void moebius_test ( );
void r8mat_is_eigen_left_test ( );
void r8mat_is_llt_test ( );
void r8mat_is_null_left_test ( );
void r8mat_is_null_right_test ( );
void r8mat_is_solution_test ( );
void r8mat_norm_fro_test ( );
void test_condition ( );
void test_determinant ( );
void test_eigen_left ( );
void test_eigen_right ( );
void test_inverse ( );
void test_llt ( );
void test_null_left ( );
void test_null_right ( );
void test_plu ( );
void test_solution ( );
void test_type ( );

/******************************************************************************/

int main ( )

/******************************************************************************/
/*
  Purpose:

    MAIN is the main program for TEST_MAT_PRB.

  Discussion:

    TEST_MAT_PRB tests the TEST_MAT library.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    15 March 2015

  Author:

    John Burkardt
*/
{
  timestamp ( );
  printf ( "\n" );
  printf ( "TEST_MAT_PRB\n" );
  printf ( "  C version\n" );
  printf ( "  Test the TEST_MAT library.\n" );
/*
  Utilities.
*/
  bvec_next_grlex_test ( );
  legendre_zeros_test ( );
  mertens_test ( );
  moebius_test ( );
  r8mat_is_eigen_left_test ( );
  r8mat_is_llt_test ( );
  r8mat_is_null_left_test ( );
  r8mat_is_null_right_test ( );
  r8mat_is_solution_test ( );
  r8mat_norm_fro_test ( );
/*
  Important things.
*/
  test_condition ( );
  test_determinant ( );
  test_eigen_left ( );
  test_eigen_right ( );
  test_inverse ( );
  test_llt ( );
  test_null_left ( );
  test_null_right ( );
  test_plu ( );
  test_solution ( );
  test_type ( );
/*
  Terminate.
*/
  printf ( "\n" );
  printf ( "TEST_MAT_PRB\n" );
  printf ( "  Normal end of execution.\n" );
  printf ( "\n" );
  timestamp ( );

  return 0;
}
/******************************************************************************/

void bvec_next_grlex_test ( )

/******************************************************************************/
/*
  Purpose:

    BVEC_NEXT_GRLEX_TEST tests BVEC_NEXT_GRLEX.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    13 March 2015

  Author:

    John Burkardt
*/
{ 
  int *b;
  int i;
  int j;
  int n = 4;

  printf ( "\n" );
  printf ( "BVEC_NEXT_GRLEX_TEST\n" );
  printf ( "  BVEC_NEXT_GRLEX computes binary vectors in GRLEX order.\n" );
  printf ( "\n" );

  b = ( int * ) malloc ( n * sizeof ( int ) );

  for ( j = 0; j < n; j++ )
  {
    b[j] = 0;
  }

  for ( i = 0; i <= 16; i++ )
  {
    printf ( "  %2d:  ", i );
    for ( j = 0; j < n; j++ )
    {
      printf ( "%d", b[j] );
    }
    printf ( "\n" );
    bvec_next_grlex ( n, b );
  }

  free ( b );

  printf ( "\n" );
  printf ( "BVEC_NEXT_GRLEX_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void legendre_zeros_test ( )

/******************************************************************************/
/*
  Purpose:

    LEGENDRE_ZEROS_TEST tests LEGENDRE_ZEROS.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    18 February 2015

  Author:

    John Burkardt
*/
{
  double *l;
  int n;

  printf ( "\n" );
  printf ( "LEGENDRE_ZEROS_TEST:\n" );
  printf ( "  LEGENDRE_ZEROS computes the zeros of the N-th Legendre\n" );
  printf ( "  polynomial.\n" );

  for ( n = 1; n <= 7; n++ )
  {
    l = legendre_zeros ( n );
    r8vec_print ( n, l, "  Legendre zeros" );
    free ( l );
  }

  printf ( "\n" );
  printf ( "LEGENDRE_ZEROS_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void mertens_test ( )

/******************************************************************************/
/*
  Purpose:

    MERTENS_TEST tests MERTENS.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    13 May 2012

  Author:

    John Burkardt
*/
{
  int c;
  int n;
  int n_data;

  printf ( "\n" );
  printf ( "MERTENS_TEST\n" );
  printf ( "  MERTENS computes the Mertens function.\n" );
  printf ( "\n" );
  printf ( "      N   Exact   MERTENS(N)\n" );
  printf ( "\n" );
 
  n_data = 0;

  for ( ; ; )
  {
     mertens_values ( &n_data, &n, &c );

    if ( n_data == 0 )
    {
      break;
    }
    printf ( "  %8d  %10d  %10d\n", n, c, mertens ( n ) );
  }
 
  printf ( "\n" );
  printf ( "MERTENS_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void moebius_test ( )

/******************************************************************************/
/*
  Purpose:

    MOEBIUS_TEST tests MOEBIUS.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    13 May 2012

  Author:

    John Burkardt
*/
{
  int c;
  int n;
  int n_data;

  printf ( "\n" );
  printf ( "MOEBIUS_TEST\n" );
  printf ( "  MOEBIUS computes the Moebius function.\n" );
  printf ( "\n" );
  printf ( "      N   Exact   MOEBIUS(N)\n" );
  printf ( "\n" );
 
  n_data = 0;

  for ( ; ; )
  {
     moebius_values ( &n_data, &n, &c );

    if ( n_data == 0 )
    {
      break;
    }

    printf ( "  %8d  %10d  %10d\n", n, c, moebius ( n ) );
  }
 
  printf ( "\n" );
  printf ( "MOEBIUS_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void r8mat_is_eigen_left_test ( )

/******************************************************************************/
/*
  Purpose:

    R8MAT_IS_EIGEN_LEFT_TEST tests R8MAT_IS_EIGEN_LEFT.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    13 March 2015

  Author:

    John Burkardt
*/
{
/*
  This is the CARRY ( 4.0, 4 ) matrix.
*/
  double a[4*4] = {
   0.13671875,   0.05859375,   0.01953125,   0.00390625, 
   0.60546875,   0.52734375,   0.39453125,   0.25390625, 
   0.25390625,   0.39453125,   0.52734375,   0.60546875, 
   0.00390625,   0.01953125,   0.05859375,   0.13671875 };
  int k = 4;
  double lam[4] = {
     1.000000000000000, 
     0.250000000000000, 
     0.062500000000000, 
     0.015625000000000 };
  int n = 4;
  double value;
  double x[4*4] = {
       1.0, 11.0, 11.0,  1.0, 
       1.0,  3.0, -3.0, -1.0, 
       1.0, -1.0, -1.0,  1.0, 
       1.0, -3.0,  3.0, -1.0 };

  printf ( "\n" );
  printf ( "R8MAT_IS_EIGEN_LEFT_TEST:\n" );
  printf ( "  R8MAT_IS_EIGEN_LEFT tests the error in the left eigensystem\n" );
  printf ( "    A' * X - X * LAMBDA = 0\n" );

  r8mat_print ( n, n, a, "  Matrix A:" );
  r8mat_print ( n, k, x, "  Eigenmatrix X:" );
  r8vec_print ( n, lam, "  Eigenvalues LAM:" );

  value = r8mat_is_eigen_left ( n, k, a, x, lam );

  printf ( "\n" );
  printf ( "  Frobenius norm of A'*X-X*LAMBDA is %g\n", value );

  printf ( "\n" );
  printf ( "R8MAT_IS_EIGEN_LEFT_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void r8mat_is_llt_test ( )

/******************************************************************************/
/*
  Purpose:

    R8MAT_IS_LLT_TEST tests R8MAT_IS_LLT.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    15 March 2015

  Author:

    John Burkardt
*/
{
/*
  Matrix is listed by columns.
*/
  double a[4*4] = { 
    2.0, 1.0, 0.0, 0.0, 
    1.0, 2.0, 1.0, 0.0, 
    0.0, 1.0, 2.0, 1.0, 
    0.0, 0.0, 1.0, 2.0 };
  double enorm;
  double l[4*4] = { 
    1.414213562373095, 0.707106781186547, 
    0.0,               0.0,               
    0.0,               1.224744871391589, 
    0.816496580927726, 0.0,               
    0.0,               0.0,               
    1.154700538379251, 0.866025403784439, 
    0.0,               0.0,               
    0.0,               1.118033988749895  };
  int m = 4;
  int n = 4;

  printf ( "\n" );
  printf ( "R8MAT_IS_LLT_TEST:\n" );
  printf ( "  R8MAT_IS_LLT tests the error in a lower triangular\n" );
  printf ( "  Cholesky factorization A = L * L' by looking at A-L*L'\n" );

  r8mat_print ( m, m, a, "  Matrix A:" );
  r8mat_print ( m, n, l, "  Cholesky factor L:" );

  enorm = r8mat_is_llt ( m, n, a, l );

  printf ( "\n" );
  printf ( "  Frobenius norm of A-L*L' is %g\n", enorm );

  printf ( "\n" );
  printf ( "R8MAT_IS_LLT_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void r8mat_is_null_left_test ( )

/******************************************************************************/
/*
  Purpose:

    R8MAT_IS_NULL_LEFT_TEST tests R8MAT_IS_NULL_LEFT.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    05 March 2015

  Author:

    John Burkardt
*/
{
/*
  Matrix is listed by columns.
*/
  double a[3*3] = { 
    1.0, 4.0, 7.0, 
    2.0, 5.0, 8.0,
    3.0, 6.0, 9.0 };
  double enorm;
  int m = 3;
  int n = 3;
  double x[3] = { 1.0, -2.0, 1.0 };

  printf ( "\n" );
  printf ( "R8MAT_IS_NULL_LEFT_TEST:\n" );
  printf ( "  R8MAT_IS_NULL_LEFT tests whether the M vector X\n" );
  printf ( "  is a left null vector of A, that is, x'*A=0.\n" );

  r8mat_print ( m, n, a, "  Matrix A:" );
  r8vec_print ( m, x, "  Vector X:" );

  enorm = r8mat_is_null_left ( m, n, a, x );

  printf ( "\n" );
  printf ( "  Frobenius norm of X'*A is %g\n", enorm );

  printf ( "\n" );
  printf ( "R8MAT_IS_NULL_LEFT_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void r8mat_is_null_right_test ( )

/******************************************************************************/
/*
  Purpose:

    R8MAT_IS_NULL_RIGHT_TEST tests R8MAT_IS_NULL_RIGHT.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    05 March 2015

  Author:

    John Burkardt
*/
{
/*
  Matrix is listed by columns.
*/
  double a[3*3] = { 
    1.0, 4.0, 7.0, 
    2.0, 5.0, 8.0,
    3.0, 6.0, 9.0 };
  double enorm;
  int m = 3;
  int n = 3;
  double x[3] = { 1.0, -2.0, 1.0 };

  printf ( "\n" );
  printf ( "R8MAT_IS_NULL_RIGHT_TEST:\n" );
  printf ( "  R8MAT_IS_NULL_RIGHT tests whether the N vector X\n" );
  printf ( "  is a right null vector of A, that is, A*x=0.\n" );

  r8mat_print ( m, n, a, "  Matrix A:" );
  r8vec_print ( n, x, "  Vector X:" );

  enorm = r8mat_is_null_right ( m, n, a, x );

  printf ( "\n" );
  printf ( "  Frobenius norm of A*x is %g\n", enorm );

  printf ( "\n" );
  printf ( "R8MAT_IS_NULL_RIGHT_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void r8mat_is_solution_test ( )

/******************************************************************************/
/*
  Purpose:

    R8MAT_IS_SOLUTION_TEST tests R8MAT_IS_SOLUTION.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    01 March 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double *b;
  double enorm;
  int i4_hi;
  int i4_lo;
  int k;
  int m;
  int n;
  double r8_hi;
  double r8_lo;
  int seed;
  double *x;

  printf ( "\n" );
  printf ( "R8MAT_IS_SOLUTION_TEST:\n" );
  printf ( "  R8MAT_IS_SOLUTION tests whether X is the solution of\n" );
  printf ( "  A*X=B by computing the Frobenius norm of the residual.\n" );
/*
  Get random shapes.
*/
  i4_lo = 1;
  i4_hi = 10;
  seed = 123456789;
  m = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  n = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  k = i4_uniform_ab ( i4_lo, i4_hi, &seed );
/*
  Get a random A.
*/
  r8_lo = -5.0;
  r8_hi = +5.0;
  a = r8mat_uniform_ab_new ( m, n, r8_lo, r8_hi, &seed );
/*
  Get a random X.
*/
  r8_lo = -5.0;
  r8_hi = +5.0;
  x = r8mat_uniform_ab_new ( n, k, r8_lo, r8_hi, &seed );
/*
  Compute B = A * X.
*/
  b = r8mat_mm_new ( m, n, k, a, x );
/*
  Compute || A*X-B||
*/
  enorm = r8mat_is_solution ( m, n, k, a, x, b );
  
  printf ( "\n" );
  printf ( "  A is %d by %d\n", m, n );
  printf ( "  X is %d by %d\n", n, k );
  printf ( "  B is %d by %d\n", m, k );
  printf ( "  Frobenius error in A*X-B is %g\n", enorm );
/*
  Free memory.
*/
  free ( a );
  free ( b );
  free ( x );

  printf ( "\n" );
  printf ( "R8MAT_IS_SOLUTION_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void r8mat_norm_fro_test ( )

/******************************************************************************/
/*
  Purpose:

    R8MAT_NORM_FRO_TEST tests R8MAT_NORM_FRO.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    04 December 2014

  Author:

    John Burkardt
*/
{
  double *a;
  int i;
  int j;
  int k;
  int m = 5;
  int n = 4;
  double t1;
  double t2;

  printf ( "\n" );
  printf ( "R8MAT_NORM_FRO_TEST\n" );
  printf ( "  R8MAT_NORM_FRO computes the Frobenius norm of a matrix.\n" );

  a = ( double * ) malloc ( m * n * sizeof ( double ) );

  k = 0;
  t1 = 0.0;
  for ( i = 0; i < m; i++ )
  {
    for ( j = 0; j < n; j++ )
    {
      k = k + 1;
      a[i+j*m] = ( double ) ( k );
      t1 = t1 + k * k;
    }
  }
  t1 = sqrt ( t1 );

  r8mat_print ( m, n, a, "  Matrix A:" );

  t2 = r8mat_norm_fro ( m, n, a );

  printf ( "\n" );
  printf ( "  Expected Frobenius norm = %g\n", t1 );
  printf ( "  Computed Frobenius norm = %g\n", t2 );

  free ( a );

  printf ( "\n" );
  printf ( "R8MAT_NORM_FRO_TEST\n" );
  printf ( "  Normal end of execution.\n" );

  return;
}
/******************************************************************************/

void test_condition ( )

/******************************************************************************/
/*
  Purpose:

    TEST_CONDITION tests the condition number computations.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    11 April 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double a_norm;
  double alpha;
  double *b;
  double b_norm;
  double beta;
  double cond1;
  double cond2;
  int i;
  int n;
  double r8_hi;
  double r8_lo;
  int seed;
  char title[21];
  double *x;

  printf ( "\n" );
  printf ( "TEST_CONDITION\n" );
  printf ( "  Compute the condition number of an example of each\n" );
  printf ( "  test matrix\n" );
  printf ( "\n" );
  printf ( "  Title                    N            COND            COND\n" );
  printf ( "\n" );
/*
  AEGERTER
*/
  strcpy ( title, "AEGERTER" );
  n = 5;
  cond1 = aegerter_condition ( n );

  a = aegerter ( n );
  b = aegerter_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  BAB
*/
  strcpy ( title, "BAB" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = bab_condition ( n, alpha, beta );

  a = bab ( n, alpha, beta );
  b = bab_inverse ( n, alpha, beta );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  BAUER
*/
  strcpy ( title, "BAUER" );
  n = 6;
  cond1 = bauer_condition ( );

  a = bauer ( );
  b = bauer_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  BIS
*/
  strcpy ( title, "BIS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = bis_condition ( alpha, beta, n );

  a = bis ( alpha, beta, n, n );
  b = bis_inverse ( alpha, beta, n  );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  BIW
*/
  strcpy ( title, "BIW" );
  n = 5;
  cond1 = biw_condition ( n );

  a = biw ( n );
  b = biw_inverse ( n  );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  BODEWIG
*/
  strcpy ( title, "BODEWIG" );
  n = 4;
  cond1 = bodewig_condition ( );

  a = bodewig ( );
  b = bodewig_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  BOOTHROYD
*/
  strcpy ( title, "BOOTHROYD" );
  n = 5;
  cond1 = boothroyd_condition ( n );

  a = boothroyd ( n );
  b = boothroyd_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  COMBIN
*/
  strcpy ( title, "COMBIN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = combin_condition ( alpha, beta, n );

  a = combin ( alpha, beta, n );
  b = combin_inverse ( alpha, beta, n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  COMPANION.
*/
  strcpy ( title, "COMPANION" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +10.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  cond1 = companion_condition ( n, x );

  a = companion ( n, x );
  b = companion_inverse ( n, x );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
  free ( x );
/*
  CONEX1
*/
  strcpy ( title, "CONEX1" );
  n = 4;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = conex1_condition ( alpha );

  a = conex1 ( alpha );
  b = conex1_inverse ( alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  CONEX2
*/
  strcpy ( title, "CONEX2" );
  n = 3;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = conex2_condition ( alpha );

  a = conex2 ( alpha );
  b = conex2_inverse ( alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  CONEX1
*/
  strcpy ( title, "CONEX1" );
  n = 4;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = conex1_condition ( alpha );

  a = conex1 ( alpha );
  b = conex1_inverse ( alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  CONEX2
*/
  strcpy ( title, "CONEX2" );
  n = 3;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = conex2_condition ( alpha );

  a = conex2 ( alpha );
  b = conex2_inverse ( alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  CONEX3
*/
  strcpy ( title, "CONEX3" );
  n = 5;
  cond1 = conex3_condition ( n );

  a = conex3 ( n );
  b = conex3_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  CONEX4
*/
  strcpy ( title, "CONEX4" );
  n = 4;
  cond1 = conex4_condition ( );

  a = conex4 ( );
  b = conex4_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DAUB2
*/
  strcpy ( title, "DAUB2" );
  n = 4;
  cond1 = daub2_condition ( n );

  a = daub2 ( n );
  b = daub2_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DAUB4
*/
  strcpy ( title, "DAUB4" );
  n = 8;
  cond1 = daub4_condition ( n );

  a = daub4 ( n );
  b = daub4_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DAUB6
*/
  strcpy ( title, "DAUB6" );
  n = 12;
  cond1 = daub6_condition ( n );

  a = daub6 ( n );
  b = daub6_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DAUB8
*/
  strcpy ( title, "DAUB8" );
  n = 16;
  cond1 = daub8_condition ( n );

  a = daub8 ( n );
  b = daub8_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DAUB10
*/
  strcpy ( title, "DAUB10" );
  n = 20;
  cond1 = daub10_condition ( n );

  a = daub10 ( n );
  b = daub10_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DAUB12
*/
  strcpy ( title, "DAUB12" );
  n = 24;
  cond1 = daub12_condition ( n );

  a = daub12 ( n );
  b = daub12_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DIAGONAL
*/
  strcpy ( title, "DIAGONAL" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  cond1 = diagonal_condition ( n, x );

  a = diagonal ( n, n, x );
  b = diagonal_inverse ( n, x );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
  free ( x );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  n = 5;
  cond1 = dif2_condition ( n );

  a = dif2 ( n, n );
  b = dif2_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  DOWNSHIFT
*/
  strcpy ( title, "DOWNSHIFT" );
  n = 5;
  cond1 = downshift_condition ( n );

  a = downshift ( n );
  b = downshift_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  EXCHANGE
*/
  strcpy ( title, "EXCHANGE" );
  n = 5;
  cond1 = exchange_condition ( n );

  a = exchange ( n, n );
  b = exchange_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  FIBONACCI2
*/
  strcpy ( title, "FIBONACCI2" );
  n = 5;
  cond1 = fibonacci2_condition ( n );

  a = fibonacci2 ( n );
  b = fibonacci2_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  GFPP
*/
  strcpy ( title, "GFPP" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = gfpp_condition ( n, alpha );

  a = gfpp ( n, alpha );
  b = gfpp_inverse ( n, alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  GIVENS
*/
  strcpy ( title, "GIVENS" );
  n = 5;
  cond1 = givens_condition ( n );

  a = givens ( n, n );
  b = givens_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  HANKEL_N
*/
  strcpy ( title, "HANKEL_N" );
  n = 5;
  cond1 = hankel_n_condition ( n );

  a = hankel_n ( n );
  b = hankel_n_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  HARMAN
*/
  strcpy ( title, "HARMAN" );
  n = 8;
  cond1 = harman_condition ( );

  a = harman ( );
  b = harman_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  HARTLEY
*/
  strcpy ( title, "HARTLEY" );
  n = 5;
  cond1 = hartley_condition ( n );

  a = hartley ( n );
  b = hartley_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  IDENTITY
*/
  strcpy ( title, "IDENTITY" );
  n = 5;
  cond1 = identity_condition ( n );

  a = identity ( n, n );
  b = identity_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  ILL3
*/
  strcpy ( title, "ILL3" );
  n = 3;
  cond1 = ill3_condition ( );

  a = ill3 ( );
  b = ill3_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  JORDAN
*/
  strcpy ( title, "JORDAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = jordan_condition ( n, alpha );

  a = jordan ( n, n, alpha );
  b = jordan_inverse ( n, alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  KERSHAW
*/
  strcpy ( title, "KERSHAW" );
  n = 4;
  cond1 = kershaw_condition ( );

  a = kershaw ( );
  b = kershaw_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  LIETZKE
*/
  strcpy ( title, "LIETZKE" );
  n = 5;
  cond1 = lietzke_condition ( n );

  a = lietzke ( n );
  b = lietzke_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  MAXIJ
*/
  strcpy ( title, "MAXIJ" );
  n = 5;
  cond1 = maxij_condition ( n );

  a = maxij ( n, n );
  b = maxij_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  MINIJ
*/
  strcpy ( title, "MINIJ" );
  n = 5;
  cond1 = minij_condition ( n );

  a = minij ( n, n );
  b = minij_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  ORTH_SYMM
*/
  strcpy ( title, "ORTH_SYMM" );
  n = 5;
  cond1 = orth_symm_condition ( n );

  a = orth_symm ( n );
  b = orth_symm_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  OTO
*/
  strcpy ( title, "OTO" );
  n = 5;
  cond1 = oto_condition ( n );

  a = oto ( n, n );
  b = oto_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  PASCAL1
*/
  strcpy ( title, "PASCAL1" );
  n = 5;
  cond1 = pascal1_condition ( n );

  a = pascal1 ( n );
  b = pascal1_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  PASCAL3
*/
  strcpy ( title, "PASCAL3" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = pascal3_condition ( n, alpha );

  a = pascal3 ( n, alpha );
  b = pascal3_inverse ( n, alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  PEI
*/
  strcpy ( title, "PEI" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = pei_condition ( alpha, n );

  a = pei ( alpha, n );
  b = pei_inverse ( alpha, n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  RODMAN
*/
  strcpy ( title, "RODMAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = rodman_condition ( alpha, n );

  a = rodman ( n, n, alpha );
  b = rodman_inverse ( n, alpha );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  RUTIS1
*/
  strcpy ( title, "RUTIS1" );
  n = 4;
  cond1 = rutis1_condition ( );

  a = rutis1 ( );
  b = rutis1_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  RUTIS2
*/
  strcpy ( title, "RUTIS2" );
  n = 4;
  cond1 = rutis2_condition ( );

  a = rutis2 ( );
  b = rutis2_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  RUTIS3
*/
  strcpy ( title, "RUTIS3" );
  n = 4;
  cond1 = rutis3_condition ( );

  a = rutis3 ( );
  b = rutis3_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  RUTIS5
*/
  strcpy ( title, "RUTIS5" );
  n = 4;
  cond1 = rutis5_condition ( );

  a = rutis5 ( );
  b = rutis5_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  SUMMATION
*/
  strcpy ( title, "SUMMATION" );
  n = 5;
  cond1 = summation_condition ( n );

  a = summation ( n, n );
  b = summation_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  SWEET1
*/
  strcpy ( title, "SWEET1" );
  n = 6;
  cond1 = sweet1_condition ( );

  a = sweet1 ( );
  b = sweet1_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  SWEET2
*/
  strcpy ( title, "SWEET2" );
  n = 6;
  cond1 = sweet2_condition ( );

  a = sweet2 ( );
  b = sweet2_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  SWEET3
*/
  strcpy ( title, "SWEET3" );
  n = 6;
  cond1 = sweet3_condition ( );

  a = sweet3 ( );
  b = sweet3_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  SWEET4
*/
  strcpy ( title, "SWEET4" );
  n = 13;
  cond1 = sweet4_condition ( );

  a = sweet4 ( );
  b = sweet4_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  TRI_UPPER
*/
  strcpy ( title, "TRI_UPPER" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  cond1 = tri_upper_condition ( alpha, n );

  a = tri_upper ( alpha, n );
  b = tri_upper_inverse ( alpha, n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  UPSHIFT
*/
  strcpy ( title, "UPSHIFT" );
  n = 5;
  cond1 = upshift_condition ( n );

  a = upshift ( n );
  b = upshift_inverse ( n );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  WILK03
*/
  strcpy ( title, "WILK03" );
  n = 3;
  cond1 = wilk03_condition ( );

  a = wilk03 ( );
  b = wilk03_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  WILK04
*/
  strcpy ( title, "WILK04" );
  n = 4;
  cond1 = wilk04_condition ( );

  a = wilk04 ( );
  b = wilk04_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  WILK05
*/
  strcpy ( title, "WILK05" );
  n = 5;
  cond1 = wilk05_condition ( );

  a = wilk05 ( );
  b = wilk05_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  n = 4;
  cond1 = wilson_condition ( );

  a = wilson ( );
  b = wilson_inverse ( );
  a_norm = r8mat_norm_l1 ( n, n, a );
  b_norm = r8mat_norm_l1 ( n, n, b );
  cond2 = a_norm * b_norm;

  printf ( "  %-20s  %4d  %14.6g  %14.6g\n", title, n, cond1, cond2 );
  free ( a );
  free ( b );

  return;
}
/******************************************************************************/

void test_determinant ( )

/******************************************************************************/
/*
  Purpose:

    TEST_DETERMINANT tests the determinant computations.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    13 April 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double b;
  double beta;
  int col_num;
  double *d;
  double d1;
  double d2;
  double d3;
  double d4;
  double d5;
  double da;
  double determ1;
  double determ2;
  double di;
  double gamma;
  int i;
  int i4_hi;
  int i4_lo;
  int ii;
  int jj;
  int k;
  int key;
  double *l;
  int m;
  int n;
  double norm_frobenius;
  double *p;
  int *pivot;
  double prob;
  int rank;
  double r8_hi;
  double r8_lo;
  int row_num;
  int seed;
  char title[21];
  double *u;
  double *v1;
  double *v2;
  double *v3;
  double *w;
  double *x;
  double x_hi;
  double x_lo;
  double x1;
  double x2;
  int x_n;
  double *y;
  int y_n;
  double y_sum;
  double *z;
  
  printf ( "\n" );
  printf ( "TEST_DETERMINANT\n" );
  printf ( "  Compute the determinants of an example of each\n" );
  printf ( "  test matrix; compare with the determinant routine,\n" );
  printf ( "  if available.  Print the matrix Frobenius norm\n" );
  printf ( "  for an estimate of magnitude.\n" );
  printf ( "\n" );
  printf ( "  Title                    N          Determ          Determ          ||A||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  n = 3;
  a = a123 ( );
  determ1 = a123_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  AEGERTER
*/
  strcpy ( title, "AEGERTER" );
  n = 5;
  a = aegerter ( n );
  determ1 = aegerter_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ANTICIRCULANT
*/
  strcpy ( title, "ANTICIRCULANT" );
  n = 3;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = anticirculant ( n, n, x );
  determ1 = anticirculant_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  ANTICIRCULANT
*/
  strcpy ( title, "ANTICIRCULANT" );
  n = 4;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = anticirculant ( n, n, x );
  determ1 = anticirculant_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  ANTICIRCULANT
*/
  strcpy ( title, "ANTICIRCULANT" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = anticirculant ( n, n, x );
  determ1 = anticirculant_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  ANTIHADAMARD
*/
  strcpy ( title, "ANTIHADAMARD" );
  n = 5;
  a = antihadamard ( n );
  determ1 = antihadamard_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ANTISYMM_RANDOM
*/
  strcpy ( title, "ANTISYMM_RANDOM" );
  n = 5;
  key = 123456789;
  a = antisymm_random ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  ANTISYMM_RANDOM
*/
  strcpy ( title, "ANTISYMM_RANDOM" );
  n = 6;
  key = 123456789;
  a = antisymm_random ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  BAB
*/
  strcpy ( title, "BAB" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = bab ( n, alpha, beta );
  determ1 = bab_determinant ( n, alpha, beta );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BAUER
*/
  strcpy ( title, "BAUER" );
  n = 6;
  a = bauer ( );
  determ1 = bauer_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BERNSTEIN
*/
  strcpy ( title, "BERNSTEIN" );
  n = 5;
  a = bernstein ( n );
  determ1 = bernstein_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BIMARKOV_RANDOM
*/
  strcpy ( title, "BIMARKOV_RANDOM" );
  n = 5;
  key = 123456789;
  a = bimarkov_random ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  BIS
*/
  strcpy ( title, "BIS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = bis ( alpha, beta, n, n );
  determ1 = bis_determinant ( alpha, beta, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BIW
*/
  strcpy ( title, "BIW" );
  n = 5;
  a = biw ( n );
  determ1 = biw_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BODEWIG
*/
  strcpy ( title, "BODEWIG" );
  n = 4;
  a = bodewig ( );
  determ1 = bodewig_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BOOTHROYD
*/
  strcpy ( title, "BOOTHROYD" );
  n = 5;
  a = boothroyd ( n );
  determ1 = boothroyd_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  BORDERBAND
*/
  strcpy ( title, "BORDERBAND" );
  n = 5;
  a = borderband ( n );
  determ1 = borderband_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CARRY
*/
  strcpy ( title, "CARRY" );
  n = 5;
  i4_lo = 2;
  i4_hi = 20;
  seed = 123456789;
  k = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  a = carry ( n, k );
  determ1 = carry_determinant ( n, k );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CAUCHY
*/
  strcpy ( title, "CAUCHY" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = cauchy ( n, x, y );
  determ1 = cauchy_determinant ( n, x, y );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
/*
  CHEBY_DIFF1
*/
  strcpy ( title, "CHEBY_DIFF1" );
  n = 5;
  a = cheby_diff1 ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  CHEBY_DIFF1
*/
  strcpy ( title, "CHEBY_DIFF1" );
  n = 5;
  a = cheby_diff1 ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  CHEBY_T
*/
  strcpy ( title, "CHEBY_T" );
  n = 5;
  a = cheby_t ( n );
  determ1 = cheby_t_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CHEBY_U
*/
  strcpy ( title, "CHEBY_U" );
  n = 5;
  a = cheby_u ( n );
  determ1 = cheby_u_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CHEBY_VAN1
*/
  strcpy ( title, "CHEBY_VAN1" );
  n = 5;
  x_lo = -1.0;
  x_hi = +1.0;
  x = r8vec_linspace_new ( n, x_lo, x_hi );
  a = cheby_van1 ( n, x_lo, x_hi, n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  CHEBY_VAN2
*/
  strcpy ( title, "CHEBY_VAN2" );
  for ( n = 2; n <= 10; n++ )
  {
    a = cheby_van2 ( n );
    determ1 = cheby_van2_determinant ( n );
    determ2 = r8mat_determinant ( n, a );
    norm_frobenius = r8mat_norm_fro ( n, n, a );
    printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
      title, n, determ1, determ2, norm_frobenius );
    free ( a );
  }
/*
  CHEBY_VAN3
*/
  strcpy ( title, "CHEBY_VAN3" );
  n = 5;
  a = cheby_van3 ( n );
  determ1 = cheby_van3_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CHOW
*/
  strcpy ( title, "CHOW" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = chow ( alpha, beta, n, n );
  determ1 = chow_determinant ( alpha, beta, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CIRCULANT
*/
  strcpy ( title, "CIRCULANT" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = circulant ( n, n, x );
  determ1 = circulant_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  CIRCULANT2
*/
  strcpy ( title, "CIRCULANT2" );
  n = 3;
  a = circulant2 ( n );
  determ1 = circulant2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CIRCULANT2
*/
  strcpy ( title, "CIRCULANT2" );
  n = 4;
  a = circulant2 ( n );
  determ1 = circulant2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CIRCULANT2
*/
  strcpy ( title, "CIRCULANT2" );
  n = 5;
  a = circulant2 ( n );
  determ1 = circulant2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CLEMENT1
*/
  strcpy ( title, "CLEMENT1" );
  n = 5;
  a = clement1 ( n );
  determ1 = clement1_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CLEMENT1
*/
  strcpy ( title, "CLEMENT1" );
  n = 6;
  a = clement1 ( n );
  determ1 = clement1_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CLEMENT2
*/
  strcpy ( title, "CLEMENT2" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = clement2 ( n, x, y );
  determ1 = clement2_determinant ( n, x, y );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
/*
  CLEMENT2
*/
  strcpy ( title, "CLEMENT2" );
  n = 6;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = clement2 ( n, x, y );
  determ1 = clement2_determinant ( n, x, y );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
/*
  COMBIN
*/
  strcpy ( title, "COMBIN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = combin ( alpha, beta, n );
  determ1 = combin_determinant ( alpha, beta, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  COMPANION
*/
  strcpy ( title, "COMPANION" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = companion ( n, x );
  determ1 = companion_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  COMPLEX_I
*/
  strcpy ( title, "COMPLEX_I" );
  n = 2;
  a = complex_i ( );
  determ1 = complex_i_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CONEX1
*/
  strcpy ( title, "CONEX1" );
  n = 4;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = conex1 ( alpha );
  determ1 = conex1_determinant ( alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CONEX2
*/
  strcpy ( title, "CONEX2" );
  n = 3;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = conex2 ( alpha );
  determ1 = conex2_determinant ( alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CONEX3
*/
  strcpy ( title, "CONEX3" );
  n = 5;
  a = conex3 ( n );
  determ1 = conex3_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CONEX4
*/
  strcpy ( title, "CONEX4" );
  n = 4;
  a = conex4 ( );
  determ1 = conex4_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CONFERENCE
  N-1 must be an odd prime or a power of an odd prime.
*/
  strcpy ( title, "CONFERENCE" );
  n = 6;
  a = conference ( n );
  determ1 = conference_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  CREATION
*/
  strcpy ( title, "CREATION" );
  n = 5;
  a = creation ( n, n );
  determ1 = creation_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DAUB2
*/
  strcpy ( title, "DAUB2" );
  n = 4;
  a = daub2 ( n );
  determ1 = daub2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DAUB4
*/
  strcpy ( title, "DAUB4" );
  n = 8;
  a = daub4 ( n );
  determ1 = daub4_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DAUB6
*/
  strcpy ( title, "DAUB6" );
  n = 12;
  a = daub6 ( n );
  determ1 = daub6_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DAUB8
*/
  strcpy ( title, "DAUB8" );
  n = 16;
  a = daub8 ( n );
  determ1 = daub8_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DAUB10
*/
  strcpy ( title, "DAUB10" );
  n = 20;
  a = daub10 ( n );
  determ1 = daub10_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DAUB12
*/
  strcpy ( title, "DAUB12" );
  n = 24;
  a = daub12 ( n );
  determ1 = daub12_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DIAGONAL
*/
  strcpy ( title, "DIAGONAL" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = diagonal ( n, n, x );
  determ1 = diagonal_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  DIF1
*/
  strcpy ( title, "DIF1" );
  n = 5;
  a = dif1 ( n, n );
  determ1 = dif1_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DIF1
*/
  strcpy ( title, "DIF1" );
  n = 6;
  a = dif1 ( n, n );
  determ1 = dif1_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DIF1CYCLIC
*/
  strcpy ( title, "DIF1CYCLIC" );
  n = 5;
  a = dif1cyclic ( n );
  determ1 = dif1cyclic_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  n = 5;
  a = dif2 ( n, n );
  determ1 = dif2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DIF2CYCLIC
*/
  strcpy ( title, "DIF2CYCLIC" );
  n = 5;
  a = dif2cyclic ( n );
  determ1 = dif2cyclic_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DORR
*/
  strcpy ( title, "DORR" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = dorr ( alpha, n );
  determ1 = dorr_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  DOWNSHIFT
*/
  strcpy ( title, "DOWNSHIFT" );
  n = 5;
  a = downshift ( n );
  determ1 = downshift_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  EBERLEIN
*/
  strcpy ( title, "EBERLEIN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = eberlein ( alpha, n );
  determ1 = eberlein_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  EULERIAN
*/
  strcpy ( title, "EULERIAN" );
  n = 5;
  a = eulerian ( n, n );
  determ1 = eulerian_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  EXCHANGE
*/
  strcpy ( title, "EXCHANGE" );
  n = 5;
  a = exchange ( n, n );
  determ1 = exchange_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FIBONACCI1
*/
  strcpy ( title, "FIBONACCI1" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = fibonacci1 ( n, alpha, beta );
  determ1 = fibonacci1_determinant ( n, alpha, beta );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FIBONACCI2
*/
  strcpy ( title, "FIBONACCI2" );
  n = 5;
  a = fibonacci2 ( n );
  determ1 = fibonacci2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FIBONACCI3
*/
  strcpy ( title, "FIBONACCI3" );
  n = 5;
  a = fibonacci3 ( n );
  determ1 = fibonacci3_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FIEDLER
*/
  strcpy ( title, "FIEDLER" );
  n = 7;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = fiedler ( n, n, x );
  determ1 = fiedler_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  FORSYTHE
*/
  strcpy ( title, "FORSYTHE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = forsythe ( alpha, beta, n );
  determ1 = forsythe_determinant ( alpha, beta, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FORSYTHE
*/
  strcpy ( title, "FORSYTHE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = forsythe ( alpha, beta, n );
  determ1 = forsythe_determinant ( alpha, beta, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FOURIER_COSINE
*/
  strcpy ( title, "FOURIER_COSINE" );
  n = 5;
  a = fourier_cosine ( n );
  determ1 = fourier_cosine_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FOURIER_SINE
*/
  strcpy ( title, "FOURIER_SINE" );
  n = 5;
  a = fourier_sine ( n );
  determ1 = fourier_sine_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  FRANK
*/
  strcpy ( title, "FRANK" );
  n = 5;
  a = frank ( n );
  determ1 = frank_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  GEAR
*/
  for ( n = 4; n <= 8; n++ )
  {
    strcpy ( title, "GEAR" );
    i4_lo = -n;
    i4_hi = +n;
    seed = 123456789;
    ii = i4_uniform_ab ( i4_lo, i4_hi, &seed );
    jj = i4_uniform_ab ( i4_lo, i4_hi, &seed );
    a = gear ( ii, jj, n );
    determ1 = gear_determinant ( ii, jj, n );
    determ2 = r8mat_determinant ( n, a );
    norm_frobenius = r8mat_norm_fro ( n, n, a );
    printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
      title, n, determ1, determ2, norm_frobenius );
    free ( a );
  }
/*
  GFPP
*/
  strcpy ( title, "GFPP" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = gfpp ( n, alpha );
  determ1 = gfpp_determinant ( n, alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  GIVENS
*/
  strcpy ( title, "GIVENS" );
  n = 5;
  a = givens ( n, n );
  determ1 = givens_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  GK316
*/
  strcpy ( title, "GK316" );
  n = 5;
  a = gk316 ( n );
  determ1 = gk316_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  GK323
*/
  strcpy ( title, "GK323" );
  n = 5;
  a = gk323 ( n, n );
  determ1 = gk323_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  GK324
*/
  strcpy ( title, "GK324" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = gk324 ( n, n, x );
  determ1 = gk324_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  GRCAR
*/
  strcpy ( title, "GRCAR" );
  n = 5;
  i4_lo = 1;
  i4_hi = n - 1;
  seed = 123456789;
  k = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  a = grcar ( n, n, k );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  HADAMARD
*/
  strcpy ( title, "HADAMARD" );
  n = 5;
  a = hadamard ( n, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  HANKEL
*/
  strcpy ( title, "HANKEL" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( 2 * n - 1, r8_lo, r8_hi, &seed );
  a = hankel ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  HANKEL_N
*/
  strcpy ( title, "HANKEL_N" );
  n = 5;
  a = hankel_n ( n );
  determ1 = hankel_n_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HANOWA
*/
  strcpy ( title, "HANOWA" );
  n = 6;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = hanowa ( alpha, n );
  determ1 = hanowa_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HARMAN
*/
  strcpy ( title, "HARMAN" );
  n = 8;
  a = harman ( );
  determ1 = harman_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HARTLEY
*/
  strcpy ( title, "HARTLEY" );
  for ( n = 5; n <= 8; n++ )
  {
    a = hartley ( n );
    determ1 = hartley_determinant ( n );
    determ2 = r8mat_determinant ( n, a );
    norm_frobenius = r8mat_norm_fro ( n, n, a );
    printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
      title, n, determ1, determ2, norm_frobenius );
    free ( a );
  }
/*
  HELMERT
*/
  strcpy ( title, "HELMERT" );
  n = 5;
  a = helmert ( n );
  determ1 = helmert_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HELMERT2
*/
  strcpy ( title, "HELMERT2" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = helmert2 ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  HERMITE
*/
  strcpy ( title, "HERMITE" );
  n = 5;
  a = hermite ( n );
  determ1 = hermite_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HERNDON
*/
  strcpy ( title, "HERNDON" );
  n = 5;
  a = herndon ( n );
  determ1 = herndon_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HILBERT
*/
  strcpy ( title, "HILBERT" );
  n = 5;
  a = hilbert ( n, n );
  determ1 = hilbert_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  HOUSEHOLDER
*/
  strcpy ( title, "HOUSEHOLDER" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = householder ( n, x );
  determ1 = householder_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  IDEM_RANDOM
*/
  strcpy ( title, "IDEM_RANDOM" );
  n = 5;
  i4_lo = 0;
  i4_hi = n;
  seed = 123456789;
  rank = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  key = 123456789;
  a = idem_random ( n, rank, key );
  determ1 = idem_random_determinant ( n, rank, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  IDENTITY
*/
  strcpy ( title, "IDENTITY" );
  n = 5;
  a = identity ( n, n );
  determ1 = identity_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  IJFACT1
*/
  strcpy ( title, "IJFACT1" );
  n = 5;
  a = ijfact1 ( n );
  determ1 = ijfact1_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  IJFACT2
*/
  strcpy ( title, "IJFACT2" );
  n = 5;
  a = ijfact2 ( n );
  determ1 = ijfact2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ILL3
*/
  strcpy ( title, "ILL3" );
  n = 3;
  a = ill3 ( );
  determ1 = ill3_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  INTEGRATION
*/
  strcpy ( title, "INTEGRATION" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = integration ( alpha, n );
  determ1 = integration_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  INVOL
*/
  strcpy ( title, "INVOL" );
  n = 5;
  a = invol ( n );
  determ1 = invol_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  INVOL_RANDOM
*/
  strcpy ( title, "INVOL_RANDOM" );
  n = 5;
  i4_lo = 0;
  i4_hi = n;
  seed = 123456789;
  rank = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  key = 123456789;
  a = invol_random ( n, rank, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  JACOBI
*/
  strcpy ( title, "JACOBI" );
  n = 5;
  a = jacobi ( n, n );
  determ1 = jacobi_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  JACOBI
*/
  strcpy ( title, "JACOBI" );
  n = 6;
  a = jacobi ( n, n );
  determ1 = jacobi_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  JORDAN
*/
  strcpy ( title, "JORDAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = jordan ( n, n, alpha );
  determ1 = jordan_determinant ( n, alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  KAHAN
*/
  strcpy ( title, "KAHAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = kahan ( alpha, n, n );
  determ1 = kahan_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  KERSHAW
*/
  strcpy ( title, "KERSHAW" );
  n = 4;
  a = kershaw ( );
  determ1 = kershaw_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  KERSHAWTRI
*/
  strcpy ( title, "KERSHAWTRI" );
  n = 5;
  x_n = ( n + 1 ) / 2;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( x_n, r8_lo, r8_hi, &seed );
  a = kershawtri ( n, x );
  determ1 = kershawtri_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  KMS
*/
  strcpy ( title, "KMS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = kms ( alpha, n, n );
  determ1 = kms_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LAGUERRE
*/
  strcpy ( title, "LAGUERRE" );
  n = 5;
  a = laguerre ( n );
  determ1 = laguerre_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LEGENDRE
*/
  strcpy ( title, "LEGENDRE" );
  n = 5;
  a = legendre ( n );
  determ1 = legendre_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LEHMER
*/
  strcpy ( title, "LEHMER" );
  n = 5;
  a = lehmer ( n, n );
  determ1 = lehmer_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LESLIE
*/
  strcpy ( title, "LESLIE" );
  n = 4;
  b = 0.025;
  di = 0.010;
  da = 0.100;
  a = leslie ( b, di, da );
  determ1 = leslie_determinant ( b, di, da );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LESP
*/
  strcpy ( title, "LESP" );
  n = 5;
  a = lesp ( n, n );
  determ1 = lesp_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LIETZKE
*/
  strcpy ( title, "LIETZKE" );
  n = 5;
  a = lietzke ( n );
  determ1 = lietzke_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LIGHTS_OUT
*/
  strcpy ( title, "LIGHTS_OUT" );
  if ( 0 )
  {
    row_num = 5;
    col_num = 5;
    n = row_num * col_num;
/*
    a = lights_out ( row_num, col_num, n );
*/
    determ2 = r8mat_determinant ( n, a );
    norm_frobenius = r8mat_norm_fro ( n, n, a );
    printf ( "  %-20s  %4d                  %14g  %14f\n", 
      title, n,          determ2, norm_frobenius );
    }
  else
  {
    printf ( "  LIGHTS_OUT          -----Not ready----\n" );
  }
/*
  LINE_ADJ
*/
  strcpy ( title, "LINE_ADJ" );
  n = 5;
  a = line_adj ( n );
  determ1 = line_adj_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LINE_ADJ
*/
  strcpy ( title, "LINE_ADJ" );
  n = 6;
  a = line_adj ( n );
  determ1 = line_adj_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LINE_LOOP_ADJ
*/
  strcpy ( title, "LINE_LOOP_ADJ" );
  n = 5;
  a = line_loop_adj ( n );
  determ1 = line_loop_adj_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  LOEWNER
*/
  strcpy ( title, "LOEWNER" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  w = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  z = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = loewner ( w, x, y, z, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( w );
  free ( x );
  free ( y );
  free ( z );
/*
  LOTKIN
*/
  strcpy ( title, "LOTKIN" );
  n = 5;
  a = lotkin ( n, n );
  determ1 = lotkin_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  MARKOV_RANDOM
*/
  strcpy ( title, "MARKOV_RANDOM" );
  n = 5;
  key = 123456789;
  a = markov_random ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  MAXIJ
*/
  strcpy ( title, "MAXIJ" );
  n = 5;
  a = maxij ( n, n );
  determ1 = maxij_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  MILNES
*/
  strcpy ( title, "MILNES" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = milnes ( n, n, x );
  determ1 = milnes_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  MINIJ
*/
  strcpy ( title, "MINIJ" );
  n = 5;
  a = minij ( n, n );
  determ1 = minij_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  MOLER1
*/
  strcpy ( title, "MOLER1" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = moler1 ( alpha, n, n );
  determ1 = moler1_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  MOLER2
*/
  strcpy ( title, "MOLER2" );
  n = 5;
  a = moler2 ( );
  determ1 = moler2_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  MOLER3
*/
  strcpy ( title, "MOLER3" );
  n = 5;
  a = moler3 ( n, n );
  determ1 = moler3_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  MOLER4
*/
  strcpy ( title, "MOLER4" );
  n = 4;
  a = moler4 ( );
  determ1 = moler4_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  NEUMANN
*/
  strcpy ( title, "NEUMANN" );
  row_num = 5;
  col_num = 5;
  n = row_num * col_num;
  a = neumann ( row_num, col_num );
  determ1 = neumann_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ONE
*/
  strcpy ( title, "ONE" );
  n = 5;
  a = one ( n, n );
  determ1 = one_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ORTEGA
*/
  strcpy ( title, "ORTEGA" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  v1 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  v2 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  v3 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = ortega ( n, v1, v2, v3 );
  determ1 = ortega_determinant ( n, v1, v2, v3 );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( v1 );
  free ( v2 );
  free ( v3 );
/*
  ORTH_RANDOM
*/
  strcpy ( title, "ORTH_RANDOM" );
  n = 5;
  key = 123456789;
  a = orth_random ( n, key );
  determ1 = orth_random_determinant ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ORTH_SYMM
*/
  strcpy ( title, "ORTH_SYMM" );
  n = 5;
  a = orth_symm ( n );
  determ1 = orth_symm_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  OTO
*/
  strcpy ( title, "OTO" );
  n = 5;
  a = oto ( n, n );
  determ1 = oto_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PARTER
*/
  strcpy ( title, "PARTER" );
  n = 5;
  a = parter ( n, n );
  determ1 = parter_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PASCAL1
*/
  strcpy ( title, "PASCAL1" );
  n = 5;
  a = pascal1 ( n );
  determ1 = pascal1_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PASCAL2
*/
  strcpy ( title, "PASCAL2" );
  n = 5;
  a = pascal2 ( n );
  determ1 = pascal2_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PASCAL3
*/
  strcpy ( title, "PASCAL3" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = pascal3 ( n, alpha );
  determ1 = pascal3_determinant ( n, alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PDS_RANDOM
*/
  strcpy ( title, "PDS_RANDOM" );
  n = 5;
  key = 123456789;
  a = pds_random ( n, key );
  determ1 = pds_random_determinant ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PEI
*/
  strcpy ( title, "PEI" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = pei ( alpha, n );
  determ1 = pei_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PERMUTATION_RANDOM
*/
  strcpy ( title, "PERMUTATION_RANDOM" );
  n = 5;
  key = 123456789;
  a = permutation_random ( n, key );
  determ1 = permutation_random_determinant ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PLU
*/
  strcpy ( title, "PLU" );
  n = 5;
  pivot = ( int * ) malloc ( n * sizeof ( int ) );
  seed = 123456789;
  for ( i = 0; i < n; i++ )
  {
    i4_lo = i;
    i4_hi = n - 1;
    pivot[i] = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  }
  a = plu ( n, pivot );
  determ1 = plu_determinant ( n, pivot );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( pivot );
/*
  POISSON
*/
  strcpy ( title, "POISSON" );
  row_num = 5;
  col_num = 5;
  n = row_num * col_num;
  a = poisson ( row_num, col_num );
  determ1 = poisson_determinant ( row_num, col_num );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  PROLATE
*/
  strcpy ( title, "PROLATE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = prolate ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  REDHEFFER
*/
  strcpy ( title, "REDHEFFER" );
  n = 5;
  a = redheffer ( n );
  determ1 = redheffer_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  REF_RANDOM
*/
  strcpy ( title, "REF_RANDOM" );
  n = 5;
  prob = 0.65;
  key = 123456789;
  a = ref_random ( n, n, prob, key );
  determ1 = ref_random_determinant ( n, prob, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  REF_RANDOM
*/
  strcpy ( title, "REF_RANDOM" );
  n = 5;
  prob = 0.85;
  key = 123456789;
  a = ref_random ( n, n, prob, key );
  determ1 = ref_random_determinant ( n, prob, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  RIEMANN
*/
  strcpy ( title, "RIEMANN" );
  n = 5;
  a = riemann ( n, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  RING_ADJ
*/
  strcpy ( title, "RING_ADJ" );
  for ( n = 1; n <= 8; n++ )
  {
    a = ring_adj ( n );
    determ1 = ring_adj_determinant ( n );
    determ2 = r8mat_determinant ( n, a );
    norm_frobenius = r8mat_norm_fro ( n, n, a );
    printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
      title, n, determ1, determ2, norm_frobenius );
    free ( a );
  }
/*
  RIS
*/
  strcpy ( title, "RIS" );
  n = 5;
  a = ris ( n );
  determ1 = ris_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  RODMAN
*/
  strcpy ( title, "RODMAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = rodman ( n, n, alpha );
  determ1 = rodman_determinant ( n, alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ROSSER1
*/
  strcpy ( title, "ROSSER1" );
  n = 8;
  a = rosser1 ( );
  determ1 = rosser1_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ROUTH
*/
  strcpy ( title, "ROUTH" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = routh ( n, x );
  determ1 = routh_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  RUTIS1
*/
  strcpy ( title, "RUTIS1" );
  n = 4;
  a = rutis1 ( );
  determ1 = rutis1_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  RUTIS2
*/
  strcpy ( title, "RUTIS2" );
  n = 4;
  a = rutis2 ( );
  determ1 = rutis2_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  RUTIS3
*/
  strcpy ( title, "RUTIS3" );
  n = 4;
  a = rutis3 ( );
  determ1 = rutis3_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  RUTIS4
*/
  strcpy ( title, "RUTIS4" );
  n = 4;
  a = rutis4 ( n );
  determ1 = rutis4_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  RUTIS5
*/
  strcpy ( title, "RUTIS5" );
  n = 4;
  a = rutis5 ( );
  determ1 = rutis5_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SCHUR_BLOCK
*/
  strcpy ( title, "SCHUR_BLOCK" );
  n = 5;
  x_n = ( n + 1 ) / 2;
  y_n = n / 2;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( x_n, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( y_n, r8_lo, r8_hi, &seed );
  a = schur_block ( n, x, y );
  determ1 = schur_block_determinant ( n, x, y );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
/*
  SKEW_CIRCULANT
*/
  strcpy ( title, "SKEW_CIRCULANT" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = skew_circulant ( n, n, x );
  determ1 = skew_circulant_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  SPLINE
*/
  strcpy ( title, "SPLINE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = spline ( n, x );
  determ1 = spline_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  STIRLING
*/
  strcpy ( title, "STIRLING" );
  n = 5;
  a = stirling ( n, n );
  determ1 = stirling_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  STRIPE
*/
  strcpy ( title, "STRIPE" );
  n = 5;
  a = stripe ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  SUMMATION
*/
  strcpy ( title, "SUMMATION" );
  n = 5;
  a = summation ( n, n );
  determ1 = summation_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SWEET1
*/
  strcpy ( title, "SWEET1" );
  n = 6;
  determ1 = sweet1_determinant ( );
  a = sweet1 ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SWEET2
*/
  strcpy ( title, "SWEET2" );
  n = 6;
  determ1 = sweet2_determinant ( );
  a = sweet2 ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SWEET3
*/
  strcpy ( title, "SWEET3" );
  n = 6;
  determ1 = sweet3_determinant ( );
  a = sweet3 ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SWEET4
*/
  strcpy ( title, "SWEET4" );
  n = 13;
  determ1 = sweet4_determinant ( );
  a = sweet4 ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SYLVESTER
*/
  strcpy ( title, "SYLVESTER" );
  n = 5;
  x_n = 3 + 1;
  y_n = 2 + 1;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( x_n, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( y_n, r8_lo, r8_hi, &seed );
  a = sylvester ( n, x_n - 1, x, y_n - 1, y );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
/*
  SYLVESTER_KAC
*/
  strcpy ( title, "SYLVESTER_KAC" );
  n = 5;
  a = sylvester_kac ( n );
  determ1 = sylvester_kac_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SYLVESTER_KAC
*/
  strcpy ( title, "SYLVESTER_KAC" );
  n = 6;
  a = sylvester_kac ( n );
  determ1 = sylvester_kac_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  SYMM_RANDOM
*/
  strcpy ( title, "SYMM_RANDOM" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  key = 123456789;
  a = symm_random ( n, d, key );
  determ1 = symm_random_determinant ( n, d, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( d );
/*
  TOEPLITZ
*/
  strcpy ( title, "TOEPLITZ" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( 2 * n - 1, r8_lo, r8_hi, &seed );
  a = toeplitz ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  TOEPLITZ_5DIAG
*/
  strcpy ( title, "TOEPLITZ_5DIAG" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d1 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  d2 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  d3 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  d4 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  d5 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = toeplitz_5diag ( n, d1, d2, d3, d4, d5 );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  TOEPLITZ_5S
*/
  strcpy ( title, "TOEPLITZ_5S" );
  row_num = 5;
  col_num = 5;
  n = row_num * col_num;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  gamma = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = toeplitz_5s ( row_num, col_num, alpha, beta, gamma, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  TOEPLITZ_PDS
*/
  strcpy ( title, "TOEPLITZ_PDS" );
  m = 3;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( m, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( m, r8_lo, r8_hi, &seed );
  y_sum = r8vec_sum ( m, y );
  for ( i = 0; i < m; i++ )
  {
    y[i] = y[i] / y_sum;
  }
  a = toeplitz_pds ( m, n, x, y );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
/*
  TOURNAMENT_RANDOM
*/
  strcpy ( title, "TOURNAMENT_RANDOM" );
  n = 5;
  key = 123456789;
  a = tournament_random ( n, key );
  determ1 = tournament_random_determinant ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  TRANSITION_RANDOM
*/
  strcpy ( title, "TRANSITION_RANDOM" );
  n = 5;
  key = 123456789;
  a = transition_random ( n, key );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  TRENCH
*/
  strcpy ( title, "TRENCH" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = trench ( alpha, n, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  TRI_UPPER
*/
  strcpy ( title, "TRI_UPPER" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = tri_upper ( alpha, n );
  determ1 = tri_upper_determinant ( alpha, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  TRIS
*/
  strcpy ( title, "TRIS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  gamma = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = tris ( n, n, alpha, beta, gamma );
  determ1 = tris_determinant ( n, alpha, beta, gamma );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  TRIV
*/
  strcpy ( title, "TRIV" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  z = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = triv ( n, x, y, z );
  determ1 = triv_determinant ( n, x, y, z );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
  free ( y );
  free ( z );
/*
  TRIW
*/
  strcpy ( title, "TRIW" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  i4_lo = 0;
  i4_hi = n - 1;
  k = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = triw ( alpha, k, n );
  determ1 = triw_determinant ( alpha, k, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  UPSHIFT
*/
  strcpy ( title, "UPSHIFT" );
  n = 5;
  a = upshift ( n );
  determ1 = upshift_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  VAND1
*/
  strcpy ( title, "VAND1" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = vand1 ( n, x );
  determ1 = vand1_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  VAND2
*/
  strcpy ( title, "VAND2" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = vand2 ( n, x );
  determ1 = vand2_determinant ( n, x );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
  free ( x );
/*
  WATHEN
*/
  strcpy ( title, "WATHEN" );
  if ( 0 )
  {
  row_num = 5;
  col_num = 5;
  n = wathen_order ( row_num, col_num );
  a = wathen ( row_num, col_num, n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
  }
  else
  {
  printf ( "  WATHEN             -----Not ready-----\n" );
  }
/*
  WILK03
*/
  strcpy ( title, "WILK03" );
  n = 3;
  a = wilk03 ( );
  determ1 = wilk03_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  WILK04
*/
  strcpy ( title, "WILK04" );
  n = 4;
  a = wilk04 ( );
  determ1 = wilk04_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  WILK05
*/
  strcpy ( title, "WILK05" );
  n = 5;
  a = wilk05 ( );
  determ1 = wilk05_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  WILK12
*/
  strcpy ( title, "WILK12" );
  n = 12;
  a = wilk12 ( );
  determ1 = wilk12_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  WILK20
*/
  strcpy ( title, "WILK20" );
  n = 20;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = wilk20 ( alpha );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );
/*
  WILK21
*/
  strcpy ( title, "WILK21" );
  n = 21;
  a = wilk21 ( n );
  determ1 = wilk21_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  n = 4;
  a = wilson ( );
  determ1 = wilson_determinant ( );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ZERO
*/
  strcpy ( title, "ZERO" );
  n = 5;
  a = zero ( n, n );
  determ1 = zero_determinant ( n );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %14g  %14g  %14f\n", 
    title, n, determ1, determ2, norm_frobenius );
  free ( a );
/*
  ZIELKE
*/
  strcpy ( title, "ZIELKE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d1 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  d2 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  d3 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = zielke ( n, d1, d2, d3 );
  determ2 = r8mat_determinant ( n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d                  %14g  %14f\n", 
    title, n,          determ2, norm_frobenius );
  free ( a );

  return;
}
/******************************************************************************/

void test_eigen_left ( )

/******************************************************************************/
/*
  Purpose:

    TEST_EIGEN_LEFT tests left eigensystems.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    15 March 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double beta;
  double *d;
  double error_frobenius;
  double gamma;
  int i;
  int i1;
  int i4_hi;
  int i4_lo;
  int k;
  int key;
  double *lambda;
  int n;
  double norm_frobenius;
  double r8_hi;
  double r8_lo;
  int rank;
  int seed;
  char title[21];
  double *v1;
  double *v2;
  double *v3;
  double *x;

  printf ( "\n" );
  printf ( "TEST_EIGEN_LEFT\n" );
  printf ( "  Compute the Frobenius norm of the eigenvalue error:\n" );
  printf ( "    X * A - LAMBDA * X\n" );
  printf ( "  given K left eigenvectors X and eigenvalues LAMBDA.\n" );
  printf ( "\n" );
  printf ( "  Title                   N     K          ||A||          ||X*A-LAMBDA*X||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  n = 3;
  k = 3;
  a = a123 ( );
  lambda = a123_eigenvalues ( );
  x = a123_eigen_left ( );
  error_frobenius = r8mat_is_eigen_left ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  CARRY
*/
  strcpy ( title, "CARRY" );
  n = 5;
  k = 5;
  i4_lo = 2;
  i4_hi = 20;
  seed = 123456789;
  i1 = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  a = carry ( n, i1 );
  lambda = carry_eigenvalues ( n, i1 );
  x = carry_eigen_left ( n, i1 );
  error_frobenius = r8mat_is_eigen_left ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  CHOW
*/
  strcpy ( title, "CHOW" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = chow ( alpha, beta, n, n );
  lambda = chow_eigenvalues ( alpha, beta, n );
  x = chow_eigen_left ( alpha, beta, n );
  error_frobenius = r8mat_is_eigen_left ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  DIAGONAL
*/
  strcpy ( title, "DIAGONAL" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = diagonal ( n, n, d );
  lambda = diagonal_eigenvalues ( n, d );
  x = diagonal_eigen_left ( n, d );
  error_frobenius = r8mat_is_eigen_left ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( d );
  free ( lambda );
  free ( x );
/*
  ROSSER1
*/
  strcpy ( title, "ROSSER1" );
  n = 8;
  k = 8;
  a = rosser1 ( );
  lambda = rosser1_eigenvalues ( );
  x = rosser1_eigen_left ( );
  error_frobenius = r8mat_is_eigen_left ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  SYMM_RANDOM
*/
  strcpy ( title, "SYMM_RANDOM" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  key = 123456789;
  a = symm_random ( n, d, key );
  lambda = symm_random_eigenvalues ( n, d, key );
  x = symm_random_eigen_left ( n, d, key );
  error_frobenius = r8mat_is_eigen_left ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( d );
  free ( lambda );
  free ( x );

  return;
}
/******************************************************************************/

void test_eigen_right ( )

/******************************************************************************/
/*
  Purpose:

    TEST_EIGEN_RIGHT tests right eigensystems.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    15 April 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double beta;
  double *d;
  double error_frobenius;
  double gamma;
  int i;
  int i1;
  int i4_hi;
  int i4_lo;
  int k;
  int key;
  double *lambda;
  int n;
  double norm_frobenius;
  double r8_hi;
  double r8_lo;
  int rank;
  int seed;
  char title[21];
  double *v1;
  double *v2;
  double *v3;
  double *x;

  printf ( "\n" );
  printf ( "TEST_EIGEN_RIGHT\n" );
  printf ( "  Compute the Frobenius norm of the eigenvalue error:\n" );
  printf ( "    A * X - X * LAMBDA\n" );
  printf ( "  given K right eigenvectors X and eigenvalues LAMBDA.\n" );
  printf ( "\n" );
  printf ( "  Title                   N     K          ||A||          ||(A*X-X*Lambda||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  n = 3;
  k = 3;
  a = a123 ( );
  lambda = a123_eigenvalues ( );
  x = a123_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  BAB
*/
  strcpy ( title, "BAB" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = bab ( n, alpha, beta );
  lambda = bab_eigenvalues ( n, alpha, beta );
  x = bab_eigen_right ( n, alpha, beta );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  BODEWIG
*/
  strcpy ( title, "BODEWIG" );
  n = 4;
  k = 4;
  a = bodewig ( );
  lambda = bodewig_eigenvalues ( );
  x = bodewig_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  CARRY
*/
  strcpy ( title, "CARRY" );
  n = 5;
  k = 5;
  i4_lo = 2;
  i4_hi = 20;
  seed = 123456789;
  i1 = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  a = carry ( n, i1 );
  lambda = carry_eigenvalues ( n, i1 );
  x = carry_eigen_right ( n, i1 );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  CHOW
*/
  strcpy ( title, "CHOW" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = chow ( alpha, beta, n, n );
  lambda = chow_eigenvalues ( alpha, beta, n );
  x = chow_eigen_right ( alpha, beta, n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  COMBIN
*/
  strcpy ( title, "COMBIN" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = combin ( alpha, beta, n );
  lambda = combin_eigenvalues ( alpha, beta, n );
  x = combin_eigen_right ( alpha, beta, n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  n = 5;
  k = 5;
  a = dif2 ( n, n );
  lambda = dif2_eigenvalues ( n );
  x = dif2_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  EXCHANGE
*/
  strcpy ( title, "EXCHANGE" );
  n = 5;
  k = 5;
  a = exchange ( n, n );
  lambda = exchange_eigenvalues ( n );
  x = exchange_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  IDEM_RANDOM
*/
  strcpy ( title, "IDEM_RANDOM" );
  n = 5;
  k = 5;
  rank = 3;
  key = 123456789;
  a = idem_random ( n, rank, key );
  lambda = idem_random_eigenvalues ( n, rank, key );
  x = idem_random_eigen_right ( n, rank, key );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  IDENTITY
*/
  strcpy ( title, "IDENTITY" );
  n = 5;
  k = 5;
  a = identity ( n, n );
  lambda = identity_eigenvalues ( n );
  x = identity_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  ILL3
*/
  strcpy ( title, "ILL3" );
  n = 3;
  k = 3;
  a = ill3 ( );
  lambda = ill3_eigenvalues ( );
  x = ill3_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  KERSHAW
*/
  strcpy ( title, "KERSHAW" );
  n = 4;
  k = 4;
  a = kershaw ( );
  lambda = kershaw_eigenvalues ( );
  x = kershaw_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  KMS
  Eigenvalue information requires 0 <= ALPHA <= 1.
*/
  strcpy ( title, "KMS" );
  n = 5;
  k = 5;
  r8_lo = 0.0;
  r8_hi = 1.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = kms ( alpha, n, n );
  lambda = kms_eigenvalues ( alpha, n );
  x = kms_eigen_right ( alpha, n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  LINE_ADJ
*/
  strcpy ( title, "LINE_ADJ" );
  n = 5;
  k = 5;
  a = line_adj ( n );
  lambda = line_adj_eigenvalues ( n );
  x = line_adj_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  LINE_LOOP_ADJ
*/
  strcpy ( title, "LINE_LOOP_ADJ" );
  n = 5;
  k = 5;
  a = line_loop_adj ( n );
  lambda = line_loop_adj_eigenvalues ( n );
  x = line_loop_adj_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  ONE
*/
  strcpy ( title, "ONE" );
  n = 5;
  k = 5;
  a = one ( n, n );
  lambda = one_eigenvalues ( n );
  x = one_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  ORTEGA
*/
  strcpy ( title, "ORTEGA" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  v1 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  v2 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  v3 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = ortega ( n, v1, v2, v3 );
  lambda = ortega_eigenvalues ( n, v1, v2, v3 );
  x = ortega_eigen_right ( n, v1, v2, v3 );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( v1 );
  free ( v2 );
  free ( v3 );
  free ( x );
/*
  OTO
*/
  strcpy ( title, "OTO" );
  n = 5;
  k = 5;
  a = oto ( n, n );
  lambda = oto_eigenvalues ( n );
  x = oto_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  PDS_RANDOM
*/
  strcpy ( title, "PDS_RANDOM" );
  n = 5;
  k = 5;
  key = 123456789;
  a = pds_random ( n, key );
  lambda = pds_random_eigenvalues ( n, key );
  x = pds_random_eigen_right ( n, key );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  PEI
*/
  strcpy ( title, "PEI" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = pei ( alpha, n );
  lambda = pei_eigenvalues ( alpha, n );
  x = pei_eigen_right ( alpha, n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  RODMAN
*/
  strcpy ( title, "RODMAN" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = rodman ( n, n, alpha );
  lambda = rodman_eigenvalues ( n, alpha );
  x = rodman_eigen_right ( n, alpha );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  ROSSER1
*/
  strcpy ( title, "ROSSER1" );
  n = 8;
  k = 8;
  a = rosser1 ( );
  lambda = rosser1_eigenvalues ( );
  x = rosser1_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  RUTIS1
*/
  strcpy ( title, "RUTIS1" );
  n = 4;
  k = 4;
  a = rutis1 ( );
  lambda = rutis1_eigenvalues ( );
  x = rutis1_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  RUTIS2
*/
  strcpy ( title, "RUTIS2" );
  n = 4;
  k = 4;
  a = rutis2 ( );
  lambda = rutis2_eigenvalues ( );
  x = rutis2_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  RUTIS5
*/
  strcpy ( title, "RUTIS5" );
  n = 4;
  k = 4;
  a = rutis5 ( );
  lambda = rutis5_eigenvalues ( );
  x = rutis5_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  SYLVESTER_KAC
*/
  strcpy ( title, "SYLVESTER_KAC" );
  n = 5;
  k = 5;
  a = sylvester_kac ( n );
  lambda = sylvester_kac_eigenvalues ( n );
  x = sylvester_kac_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  SYMM_RANDOM
*/
  strcpy ( title, "SYMM_RANDOM" );
  n = 5;
  k = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  key = 123456789;
  a = symm_random ( n, d, key );
  lambda = symm_random_eigenvalues ( n, d, key );
  x = symm_random_eigen_right ( n, d, key );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( d );
  free ( lambda );
  free ( x );
/*
  WILK12
*/
  strcpy ( title, "WILK12" );
  n = 12;
  k = 12;
  a = wilk12 ( );
  lambda = wilk12_eigenvalues ( );
  x = wilk12_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  n = 4;
  k = 4;
  a = wilson ( );
  lambda = wilson_eigenvalues ( );
  x = wilson_eigen_right ( );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );
/*
  ZERO
*/
  strcpy ( title, "ZERO" );
  n = 5;
  k = 5;
  a = zero ( n, n );
  lambda = zero_eigenvalues ( n );
  x = zero_eigen_right ( n );
  error_frobenius = r8mat_is_eigen_right ( n, k, a, x, lambda );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s %4d  %4d  %14g  %14g\n",
    title, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( lambda );
  free ( x );

  return;
}
/******************************************************************************/

void test_inverse ( )

/******************************************************************************/
/*
  Purpose:

    TEST_INVERSE tests the inverse computations.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    17 April 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double *b;;
  double beta;
  double *c;
  double *d;
  double error_ab;
  double error_ac;;
  double gamma;
  int i;
  int i4_hi;
  int i4_lo;
  int ii;
  int jj;
  int k;
  int key;
  double *l;
  int n;
  double norma_frobenius;
  double normc_frobenius;
  double *p;
  int *pivot;
  double r8_hi;
  double r8_lo;
  int seed;
  char title[21];
  double *u;
  double *v1;
  double *v2;
  double *v3;
  double *w;
  double *x;
  int x_n;
  double *y;
  int y_n;
  double *z;

  printf ( "\n" );
  printf ( "TEST_INVERSE\n" );
  printf ( "  A = a test matrix of order N;\n" );
  printf ( "  B = inverse as computed by a routine.\n" );
  printf ( "  C = inverse as computed by R8MAT_INVERSE.\n" );
  printf ( "\n" );
  printf ( "  ||I-AB|| = Frobenius norm of I-A*B.\n" );
  printf ( "  ||I-AC|| = Frobenius norm of I-A*C.\n" );
  printf ( "  ||I-AB|| = Frobenius norm of I-A*B.\n" );
  printf ( "\n" );
  printf ( "  Title                    N   " );
  printf ( "   ||A||      ||C||  ||I-AC||    ||I-AB||\n" );
  printf ( "\n" );
/*
  AEGERTER
*/
  strcpy ( title, "AEGERTER" );
  n = 5;
  a = aegerter ( n );
  b = aegerter_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BAB
*/
  strcpy ( title, "BAB" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = bab ( n, alpha, beta );
  b = bab_inverse ( n, alpha, beta );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BAUER
*/
  strcpy ( title, "BAUER" );
  n = 6;
  a = bauer ( );
  b = bauer_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BERNSTEIN
*/
  strcpy ( title, "BERNSTEIN" );
  n = 5;
  a = bernstein ( n );
  b = bernstein_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BIS
*/
  strcpy ( title, "BIS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = bis ( alpha, beta, n, n );
  b = bis_inverse ( alpha, beta, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BIW
*/
  strcpy ( title, "BIW" );
  n = 5;
  a = biw ( n );
  b = biw_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BODEWIG
*/
  strcpy ( title, "BODEWIG" );
  n = 4;
  a = bodewig ( );
  b = bodewig_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BOOTHROYD
*/
  strcpy ( title, "BOOTHROYD" );
  n = 5;
  a = boothroyd ( n );
  b = boothroyd_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  BORDERBAND
*/
  strcpy ( title, "BORDERBAND" );
  n = 5;
  a = borderband ( n );
  b = borderband_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CARRY
*/
  strcpy ( title, "CARRY" );
  n = 5;
  seed = 123456789;
  i4_lo = 2;
  i4_hi = 20;
  k = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  a = carry ( n, k );
  b = carry_inverse ( n, k );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CAUCHY
*/
  strcpy ( title, "CAUCHY" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = cauchy ( n, x, y );
  b = cauchy_inverse ( n, x, y );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
  free ( y );
/*
  CHEBY_T
*/
  strcpy ( title, "CHEBY_T" );
  n = 5;
  a = cheby_t ( n );
  b = cheby_t_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CHEBY_U
*/
  strcpy ( title, "CHEBY_U" );
  n = 5;
  a = cheby_u ( n );
  b = cheby_u_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CHEBY_VAN2
*/
  strcpy ( title, "CHEBY_VAN2" );
  n = 5;
  a = cheby_van2 ( n );
  b = cheby_van2_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CHEBY_VAN3
*/
  strcpy ( title, "CHEBY_VAN3" );
  n = 5;
  a = cheby_van3 ( n );
  b = cheby_van3_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CHOW
*/
  strcpy ( title, "CHOW" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = chow ( alpha, beta, n, n );
  b = chow_inverse ( alpha, beta, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CIRCULANT
*/
  strcpy ( title, "CIRCULANT" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = circulant ( n, n, x );
  b = circulant_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  CIRCULANT2
*/
  strcpy ( title, "CIRCULANT2" );
  n = 5;
  a = circulant2 ( n );
  b = circulant2_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CLEMENT1
  N must be even.
*/
  strcpy ( title, "CLEMENT1" );
  n = 6;
  a = clement1 ( n );
  b = clement1_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CLEMENT2
  N must be even.
*/
  strcpy ( title, "CLEMENT2" );
  n = 6;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = clement2 ( n, x, y );
  b = clement2_inverse ( n, x, y );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
  free ( y );
/*
  COMBIN
*/
  strcpy ( title, "COMBIN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = combin ( alpha, beta, n );
  b = combin_inverse ( alpha, beta, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  COMPANION
*/
  strcpy ( title, "COMPANION" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = companion ( n, x );
  b = companion_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  COMPLEX_I
*/
  strcpy ( title, "COMPLEX_I" );
  n = 2;
  a = complex_i ( );
  b = complex_i_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CONEX1
*/
  strcpy ( title, "CONEX1" );
  n = 4;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = conex1 ( alpha );
  b = conex1_inverse ( alpha );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CONEX2
*/
  strcpy ( title, "CONEX2" );
  n = 3;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = conex2 ( alpha );
  b = conex2_inverse ( alpha );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CONEX3
*/
  strcpy ( title, "CONEX3" );
  n = 5;
  a = conex3 ( n );
  b = conex3_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  CONFERENCE
  N-1 must be an odd prime or a power of an odd prime.
*/
  strcpy ( title, "CONFERENCE" );
  n = 6;
  a = conference ( n );
  b = conference_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DAUB2
*/
  strcpy ( title, "DAUB2" );
  n = 4;
  a = daub2 ( n );
  b = daub2_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DAUB4
*/
  strcpy ( title, "DAUB4" );
  n = 8;
  a = daub4 ( n );
  b = daub4_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DAUB6
*/
  strcpy ( title, "DAUB6" );
  n = 12;
  a = daub6 ( n );
  b = daub6_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DAUB8
*/
  strcpy ( title, "DAUB8" );
  n = 16;
  a = daub8 ( n );
  b = daub8_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DAUB10
*/
  strcpy ( title, "DAUB10" );
  n = 20;
  a = daub10 ( n );
  b = daub10_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DAUB12
*/
  strcpy ( title, "DAUB12" );
  n = 24;
  a = daub12 ( n );
  b = daub12_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DIAGONAL
*/
  strcpy ( title, "DIAGONAL" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = diagonal ( n, n, x );
  b = diagonal_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  DIF1
  N must be even.
*/
  strcpy ( title, "DIF1" );
  n = 6;
  a = dif1 ( n, n );
  b = dif1_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  n = 5;
  a = dif2 ( n, n );
  b = dif2_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DORR
*/
  strcpy ( title, "DORR" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = dorr ( alpha, n );
  b = dorr_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  DOWNSHIFT
*/
  strcpy ( title, "DOWNSHIFT" );
  n = 5;
  a = downshift ( n );
  b = downshift_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  EULERIAN
*/
  strcpy ( title, "EULERIAN" );
  n = 5;
  a = eulerian ( n, n );
  b = eulerian_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  EXCHANGE
*/
  strcpy ( title, "EXCHANGE" );
  n = 5;
  a = exchange ( n, n );
  b = exchange_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  FIBONACCI2
*/
  strcpy ( title, "FIBONACCI2" );
  n = 5;
  a = fibonacci2 ( n );
  b = fibonacci2_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  FIBONACCI3
*/
  strcpy ( title, "FIBONACCI3" );
  n = 5;
  a = fibonacci3 ( n );
  b = fibonacci3_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  FIEDLER.
  The FIEDLER_INVERSE routine assumes the X vector is sorted.
*/
  strcpy ( title, "FIEDLER" );
  n = 7;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  r8vec_sort_bubble_a ( n, x );
  a = fiedler ( n, n, x );
  b = fiedler_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  FORSYTHE
*/
  strcpy ( title, "FORSYTHE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = forsythe ( alpha, beta, n );
  b = forsythe_inverse ( alpha, beta, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  FOURIER_COSINE
*/
  strcpy ( title, "FOURIER_COSINE" );
  n = 5;
  a = fourier_cosine ( n );
  b = fourier_cosine_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  FOURIER_SINE
*/
  strcpy ( title, "FOURIER_SINE" );
  n = 5;
  a = fourier_sine ( n );
  b = fourier_sine_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  FRANK
*/
  strcpy ( title, "FRANK" );
  n = 5;
  a = frank ( n );
  b = frank_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  GFPP
*/
  strcpy ( title, "GFPP" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = gfpp ( n, alpha );
  b = gfpp_inverse ( n, alpha );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  GIVENS
*/
  strcpy ( title, "GIVENS" );
  n = 5;
  a = givens ( n, n );
  b = givens_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  GK316
*/
  strcpy ( title, "GK316" );
  n = 5;
  a = gk316 ( n );
  b = gk316_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  GK323
*/
  strcpy ( title, "GK323" );
  n = 5;
  a = gk323 ( n, n );
  b = gk323_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  GK324
*/
  strcpy ( title, "GK324" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = gk324 ( n, n, x );
  b = gk324_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  HANKEL_N
*/
  strcpy ( title, "HANKEL_N" );
  n = 5;
  a = hankel_n ( n );
  b = hankel_n_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HANOWA
  N must be even.
*/
  strcpy ( title, "HANOWA" );
  n = 6;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = hanowa ( alpha, n );
  b = hanowa_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HARMAN
*/
  strcpy ( title, "HARMAN" );
  n = 8;
  a = harman (  );
  b = harman_inverse (  );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HARTLEY
*/
  strcpy ( title, "HARTLEY" );
  n = 5;
  a = hartley ( n );
  b = hartley_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HELMERT
*/
  strcpy ( title, "HELMERT" );
  n = 5;
  a = helmert ( n );
  b = helmert_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HELMERT2
*/
  strcpy ( title, "HELMERT2" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = helmert2 ( n, x );
  b = helmert2_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  HERMITE
*/
  strcpy ( title, "HERMITE" );
  n = 5;
  a = hermite ( n );
  b = hermite_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HERNDON
*/
  strcpy ( title, "HERNDON" );
  n = 5;
  a = herndon ( n );
  b = herndon_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HILBERT
*/
  strcpy ( title, "HILBERT" );
  n = 5;
  a = hilbert ( n, n );
  b = hilbert_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  HOUSEHOLDER
*/
  strcpy ( title, "HOUSEHOLDER" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = householder ( n, x );
  b = householder_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  IDENTITY
*/
  strcpy ( title, "IDENTITY" );
  n = 5;
  a = identity ( n, n );
  b = identity_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  ILL3
*/
  strcpy ( title, "ILL3" );
  n = 3;
  a = ill3 ( );
  b = ill3_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  INTEGRATION
*/
  strcpy ( title, "INTEGRATION" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = integration ( alpha, n );
  b = integration_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  INVOL
*/
  strcpy ( title, "INVOL" );
  n = 5;
  a = invol ( n );
  b = invol_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  JACOBI
  N must be even.
*/
  strcpy ( title, "JACOBI" );
  n = 6;
  a = jacobi ( n, n );
  b = jacobi_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  JORDAN
*/
  strcpy ( title, "JORDAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = jordan ( n, n, alpha );
  b = jordan_inverse ( n, alpha );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  KAHAN
*/
  strcpy ( title, "KAHAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = kahan ( alpha, n, n );
  b = kahan_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  KERSHAW
*/
  strcpy ( title, "KERSHAW" );
  n = 4;
  a = kershaw ( );
  b = kershaw_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  KERSHAWTRI
*/
  strcpy ( title, "KERSHAWTRI" );
  n = 5;
  x_n = ( n + 1 ) / 2;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( x_n, r8_lo, r8_hi, &seed );
  a = kershawtri ( n, x );
  b = kershawtri_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  KMS
*/
  strcpy ( title, "KMS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = kms ( alpha, n, n );
  b = kms_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LAGUERRE
*/
  strcpy ( title, "LAGUERRE" );
  n = 5;
  a = laguerre ( n );
  b = laguerre_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LEGENDRE
*/
  strcpy ( title, "LEGENDRE" );
  n = 5;
  a = legendre ( n );
  b = legendre_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LEHMER
*/
  strcpy ( title, "LEHMER" );
  n = 5;
  a = lehmer ( n, n );
  b = lehmer_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LESP
*/
  strcpy ( title, "LESP" );
  n = 5;
  a = lesp ( n, n );
  b = lesp_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LIETZKE
*/
  strcpy ( title, "LIETZKE" );
  n = 5;
  a = lietzke ( n );
  b = lietzke_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LINE_ADJ
  N must be even.
*/
  strcpy ( title, "LINE_ADJ" );
  n = 6;
  a = line_adj ( n );
  b = line_adj_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  LOTKIN
*/
  strcpy ( title, "LOTKIN" );
  n = 5;
  a = lotkin ( n, n );
  b = lotkin_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  MAXIJ
*/
  strcpy ( title, "MAXIJ" );
  n = 5;
  a = maxij ( n, n );
  b = maxij_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  MILNES
*/
  strcpy ( title, "MILNES" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = milnes ( n, n, x );
  b = milnes_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  MINIJ
*/
  strcpy ( title, "MINIJ" );
  n = 5;
  a = minij ( n, n );
  b = minij_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  MOLER1
*/
  strcpy ( title, "MOLER1" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = moler1 ( alpha, n, n );
  b = moler1_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  MOLER3
*/
  strcpy ( title, "MOLER3" );
  n = 5;
  a = moler3 ( n, n );
  b = moler3_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  ORTEGA
*/
  strcpy ( title, "ORTEGA" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  v1 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  v2 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  v3 = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = ortega ( n, v1, v2, v3 );
  b = ortega_inverse ( n, v1, v2, v3 );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( v1 );
  free ( v2 );
  free ( v3 );
/*
  ORTH_SYMM
*/
  strcpy ( title, "ORTH_SYMM" );
  n = 5;
  a = orth_symm ( n );
  b = orth_symm_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  OTO
*/
  strcpy ( title, "OTO" );
  n = 5;
  a = oto ( n, n );
  b = oto_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PARTER
*/
  strcpy ( title, "PARTER" );
  n = 5;
  a = parter ( n, n );
  b = parter_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PASCAL1
*/
  strcpy ( title, "PASCAL1" );
  n = 5;
  a = pascal1 ( n );
  b = pascal1_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PASCAL2
*/
  strcpy ( title, "PASCAL2" );
  n = 5;
  a = pascal2 ( n );
  b = pascal2_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PASCAL3
*/
  strcpy ( title, "PASCAL3" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = pascal3 ( n, alpha );
  b = pascal3_inverse ( n, alpha );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PDS_RANDOM
*/
  strcpy ( title, "PDS_RANDOM" );
  n = 5;
  key = 123456789;
  a = pds_random ( n, key );
  b = pds_random_inverse ( n, key );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PEI
*/
  strcpy ( title, "PEI" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = pei ( alpha, n );
  b = pei_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PERMUTATION_RANDOM
*/
  strcpy ( title, "PERMUTATION_RANDOM" );
  n = 5;
  key = 123456789;
  a = permutation_random ( n, key );
  b = permutation_random_inverse ( n, key );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  PLU
*/
  strcpy ( title, "PLU" );
  n = 5;
  pivot = ( int * ) malloc ( n * sizeof ( int ) );
  seed = 123456789;
  for ( i = 0; i < n; i++ )
  {
    i4_lo = i;
    i4_hi = n - 1;
    pivot[i] = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  }
  a = plu ( n, pivot );
  b = plu_inverse ( n, pivot );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( pivot );
/*
  RIS
*/
  strcpy ( title, "RIS" );
  n = 5;
  a = ris ( n );
  b = ris_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  RODMAN
*/
  strcpy ( title, "RODMAN" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = rodman ( n, n, alpha );
  b = rodman_inverse ( n, alpha );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  RUTIS1
*/
  strcpy ( title, "RUTIS1" );
  n = 4;
  a = rutis1 ( );
  b = rutis1_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  RUTIS2
*/
  strcpy ( title, "RUTIS2" );
  n = 4;
  a = rutis2 ( );
  b = rutis2_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  RUTIS3
*/
  strcpy ( title, "RUTIS3" );
  n = 4;
  a = rutis3 ( );
  b = rutis3_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  RUTIS4
*/
  strcpy ( title, "RUTIS4" );
  n = 5;
  a = rutis4 ( n );
  b = rutis4_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  RUTIS5
*/
  strcpy ( title, "RUTIS5" );
  n = 4;
  a = rutis5 ( );
  b = rutis5_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SCHUR_BLOCK
*/
  strcpy ( title, "SCHUR_BLOCK" );
  n = 5;
  x_n = ( n + 1 ) / 2;
  y_n = n / 2;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( x_n, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( y_n, r8_lo, r8_hi, &seed );
  a = schur_block ( n, x, y );
  b = schur_block_inverse ( n, x, y );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
  free ( y );
/*
  SPLINE
*/
  strcpy ( title, "SPLINE" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = spline ( n, x );
  b = spline_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  STIRLING
*/
  strcpy ( title, "STIRLING" );
  n = 5;
  a = stirling ( n, n );
  b = stirling_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SUMMATION
*/
  strcpy ( title, "SUMMATION" );
  n = 5;
  a = summation ( n, n );
  b = summation_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SWEET1
*/
  strcpy ( title, "SWEET1" );
  n = 6;
  a = sweet1 ( );
  b = sweet1_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SWEET2
*/
  strcpy ( title, "SWEET2" );
  n = 6;
  a = sweet2 ( );
  b = sweet2_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SWEET3
*/
  strcpy ( title, "SWEET3" );
  n = 6;
  a = sweet3 ( );
  b = sweet3_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SWEET4
*/
  strcpy ( title, "SWEET4" );
  n = 13;
  a = sweet4 ( );
  b = sweet4_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SYLVESTER_KAC
  N must be even.
*/
  strcpy ( title, "SYLVESTER_KAC" );
  n = 6;
  a = sylvester_kac ( n );
  b = sylvester_kac_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  SYMM_RANDOM
*/
  strcpy ( title, "SYMM_RANDOM" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  d = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  key = 123456789;
  a = symm_random ( n, d, key );
  b = symm_random_inverse ( n, d, key );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( d );
/*
  TRI_UPPER
*/
  strcpy ( title, "TRI_UPPER" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = tri_upper ( alpha, n );
  b = tri_upper_inverse ( alpha, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  TRIS
*/
  strcpy ( title, "TRIS" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  beta = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  gamma = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = tris ( n, n, alpha, beta, gamma );
  b = tris_inverse ( n, alpha, beta, gamma );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  TRIV
*/
  strcpy ( title, "TRIV" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  y = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  z = r8vec_uniform_ab_new ( n - 1, r8_lo, r8_hi, &seed );
  a = triv ( n, x, y, z );
  b = triv_inverse ( n, x, y, z );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
  free ( y );
  free ( z );
/*
  TRIW
*/
  strcpy ( title, "TRIW" );
  n = 5;
  i4_lo = 0;
  i4_hi = n - 1;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  k = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = triw ( alpha, k, n );
  b = triw_inverse ( alpha, k, n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  UPSHIFT
*/
  strcpy ( title, "UPSHIFT" );
  n = 5;
  a = upshift ( n );
  b = upshift_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  VAND1
*/
  strcpy ( title, "VAND1" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = vand1 ( n, x );
  b = vand1_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
  free ( x );
/*
  VAND2
*/
  strcpy ( title, "VAND2" );
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( n, r8_lo, r8_hi, &seed );
  a = vand2 ( n, x );
  b = vand2_inverse ( n, x );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  WILK03
*/
  strcpy ( title, "WILK03" );
  n = 3;
  a = wilk03 ( );
  b = wilk03_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  WILK04
*/
  strcpy ( title, "WILK04" );
  n = 4;
  a = wilk04 ( );
  b = wilk04_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  WILK05
*/
  strcpy ( title, "WILK05" );
  n = 5;
  a = wilk05 ( );
  b = wilk05_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  WILK21
*/
  strcpy ( title, "WILK21" );
  n = 21;
  a = wilk21 ( n );
  b = wilk21_inverse ( n );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  n = 4;
  a = wilson ( );
  b = wilson_inverse ( );
  c = r8mat_inverse ( n, a );
  error_ab = r8mat_is_inverse ( n, a, b );
  error_ac = r8mat_is_inverse ( n, a, c );
  norma_frobenius = r8mat_norm_fro ( n, n, a );
  normc_frobenius = r8mat_norm_fro ( n, n, c );
  printf ( "  %-20s  %4d  %10g  %10g  %10g  %10g\n",
    title, n, norma_frobenius, normc_frobenius, error_ac, error_ab );
  free ( a );
  free ( b );
  free ( c );

  return;
}
/******************************************************************************/

void test_llt ( )

/******************************************************************************/
/*
  Purpose:

    TEST_LLT tests LLT factors.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    09 April 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double error_frobenius;
  double *l;
  int m;
  int n;
  double norm_a_frobenius;
  double r8_hi;
  double r8_lo;
  int seed;
  char title[21];

  printf ( "\n" );
  printf ( "TEST_LLT\n" );
  printf ( "  A = a test matrix of order M by M\n" );
  printf ( "  L is an M by N lower triangular Cholesky factor.\n" );
  printf ( "\n" );
  printf ( "  ||A|| = Frobenius norm of A.\n" );
  printf ( "  ||A-LLT|| = Frobenius norm of A-L*L'.\n" );
  printf ( "\n" );
  printf ( "  Title                    M     N      ||A||            ||A-LLT||\n" );
  printf ( "\n" );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  m = 5;
  n = 5;
  a = dif2 ( m, n );
  l = dif2_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  GIVENS
*/
  strcpy ( title, "GIVENS" );
  m = 5;
  n = 5;
  a = givens ( m, n );
  l = givens_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  KERSHAW
*/
  strcpy ( title, "KERSHAW" );
  m = 4;
  n = 4;
  a = kershaw ( );
  l = kershaw_llt ( );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  LEHMER
*/
  strcpy ( title, "LEHMER" );
  m = 5;
  n = 5;
  a = lehmer ( n, n );
  l = lehmer_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  MINIJ
*/
  strcpy ( title, "MINIJ" );
  m = 5;
  n = 5;
  a = minij ( n, n );
  l = minij_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  MOLER1
*/
  strcpy ( title, "MOLER1" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = moler1 ( alpha, m, n );
  l = moler1_llt ( alpha, n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  MOLER3
*/
  strcpy ( title, "MOLER3" );
  m = 5;
  n = 5;
  a = moler3 ( m, n );
  l = moler3_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  OTO
*/
  strcpy ( title, "OTO" );
  m = 5;
  n = 5;
  a = oto ( m, n );
  l = oto_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  PASCAL2
*/
  strcpy ( title, "PASCAL2" );
  m = 5;
  n = 5;
  a = pascal2 ( n );
  l = pascal2_llt ( n );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  m = 4;
  n = 4;
  a = wilson ( );
  l = wilson_llt ( );
  error_frobenius = r8mat_is_llt ( m, n, a, l );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14.6g  %14.6g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );

  return;
}
/******************************************************************************/

void test_null_left ( )

/******************************************************************************/
/*
  Purpose:

    TEST_NULL_LEFT tests the left null vectors.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    13 March 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  int col_num;
  double error_l2;
  double f1;
  double f2;
  int m;
  int n;
  double norm_a_frobenius;
  double norm_x_l2;
  double r8_hi;
  double r8_lo;
  int seed;
  char title[21];
  double *x;

  printf ( "\n" );
  printf ( "TEST_NULL_LEFT\n" );
  printf ( "  A = a test matrix of order M by N\n" );
  printf ( "  x = an M vector, candidate for a left null vector.\n" );
  printf ( "\n" );
  printf ( "  ||A|| = Frobenius norm of A.\n" );
  printf ( "  ||x|| = L2 norm of x.\n" );
  printf ( "  ||A'*x||/||x|| = L2 norm of A'*x over L2 norm of x.\n" );
  printf ( "\n" );
  printf ( "  Title                    M     N       " );
  printf ( "||A||            ||x||        ||A'*x||/||x||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  m = 3;
  n = 3;
  a = a123 ( );
  x = a123_null_left ( );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  CHEBY_DIFF1
*/
  strcpy ( title, "CHEBY_DIFF1" );
  m = 5;
  n = 5;
  a = cheby_diff1 ( n );
  x = cheby_diff1_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  CREATION
*/
  strcpy ( title, "CREATION" );
  m = 5;
  n = 5;
  a = creation ( m, n );
  x = creation_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  DIF1
  Only has null vectors for M odd.
*/
  strcpy ( title, "DIF1" );
  m = 5;
  n = 5;
  a = dif1 ( m, n );
  x = dif1_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  DIF1CYCLIC
*/
  strcpy ( title, "DIF1CYCLIC" );
  m = 5;
  n = 5;
  a = dif1cyclic ( n );
  x = dif1cyclic_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  DIF2CYCLIC
*/
  strcpy ( title, "DIF2CYCLIC" );
  m = 5;
  n = 5;
  a = dif2cyclic ( n );
  x = dif2cyclic_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  EBERLEIN
*/
  strcpy ( title, "EBERLEIN" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = eberlein ( alpha, n );
  x = eberlein_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  FIBONACCI1
*/
  strcpy ( title, "FIBONACCI1" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  f1 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  f2 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = fibonacci1 ( n, f1, f2 );
  x = fibonacci1_null_left ( m, n, f1, f2 );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  LAUCHLI
*/
  strcpy ( title, "LAUCHLI" );
  m = 6;
  n = m - 1;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = lauchli ( alpha, m, n );
  x = lauchli_null_left ( alpha, m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  LINE_ADJ
*/
  strcpy ( title, "LINE_ADJ" );
  m = 7;
  n = 7;
  a = line_adj ( n );
  x = line_adj_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  MOLER2
*/
  strcpy ( title, "MOLER2" );
  m = 5;
  n = 5;
  a = moler2 ( );
  x = moler2_null_left ( );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ONE
*/
  strcpy ( title, "ONE" );
  m = 5;
  n = 5;
  a = one ( n, n );
  x = one_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  RING_ADJ
  M must be a multiple of 4 for there to be a null vector.
*/
  strcpy ( title, "RING_ADJ" );
  m = 12;
  n = 12;
  a = ring_adj ( n );
  x = ring_adj_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ROSSER1
*/
  strcpy ( title, "ROSSER1" );
  m = 8;
  n = 8;
  a = rosser1 ( );
  x = rosser1_null_left ( );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ZERO
*/
  strcpy ( title, "ZERO" );
  m = 5;
  n = 5;
  a = zero ( m, n );
  x = zero_null_left ( m, n );
  error_l2 = r8mat_is_null_left ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( m, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );

  return;
}
/******************************************************************************/

void test_null_right ( )

/******************************************************************************/
/*
  Purpose:

    TEST_NULL_RIGHT tests the right null vectors.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    11 March 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  int col_num;
  double error_l2;
  double f1;
  double f2;
  int m;
  int n;
  double norm_a_frobenius;
  double norm_x_l2;
  double r8_hi;
  double r8_lo;
  int row_num;
  int seed;
  char title[21];
  double *x;

  printf ( "\n" );
  printf ( "TEST_NULL_RIGHT\n" );
  printf ( "  A = a test matrix of order M by N\n" );
  printf ( "  x = an N vector, candidate for a right null vector.\n" );
  printf ( "\n" );
  printf ( "  ||A|| = Frobenius norm of A.\n" );
  printf ( "  ||x|| = L2 norm of x.\n" );
  printf ( "  ||A*x||/||x|| = L2 norm of A*x over L2 norm of x.\n" );
  printf ( "\n" );
  printf ( "  Title                    M     N         " );
  printf ( "||A||            ||x||        ||A*x||/||x||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  m = 3;
  n = 3;
  a = a123 ( );
  x = a123_null_right ( );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ARCHIMEDES
*/
  strcpy ( title, "ARCHIMEDES" );
  m = 7;
  n = 8;
  a = archimedes ( );
  x = archimedes_null_right ( );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  CHEBY_DIFF1
*/
  strcpy ( title, "CHEBY_DIFF1" );
  m = 5;
  n = 5;
  a = cheby_diff1 ( n );
  x = cheby_diff1_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  CREATION
*/
  strcpy ( title, "CREATION" );
  m = 5;
  n = 5;
  a = creation ( m, n );
  x = creation_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  DIF1
  Only has null vectors for N odd.
*/
  strcpy ( title, "DIF1" );
  m = 5;
  n = 5;
  a = dif1 ( m, n );
  x = dif1_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  DIF1CYCLIC
*/
  strcpy ( title, "DIF1CYCLIC" );
  m = 5;
  n = 5;
  a = dif1cyclic ( n );
  x = dif1cyclic_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  DIF2CYCLIC
*/
  strcpy ( title, "DIF2CYCLIC" );
  m = 5;
  n = 5;
  a = dif2cyclic ( n );
  x = dif2cyclic_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  FIBONACCI1
*/
  strcpy ( title, "FIBONACCI1" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  f1 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  f2 = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  a = fibonacci1 ( n, f1, f2 );
  x = fibonacci1_null_right ( m, n, f1, f2 );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  HAMMING
*/
  strcpy ( title, "HAMMING" );
  m = 5;
  n = i4_power ( 2, 5 ) - 1;
  a = hamming ( m, n );
  x = hamming_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  LINE_ADJ
*/
  strcpy ( title, "LINE_ADJ" );
  m = 7;
  n = 7;
  a = line_adj ( n );
  x = line_adj_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  MOLER2
*/
  strcpy ( title, "MOLER2" );
  m = 5;
  n = 5;
  a = moler2 ( );
  x = moler2_null_right ( );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  NEUMANN
*/
  strcpy ( title, "NEUMANN" );
  row_num = 5;
  col_num = 5;
  m = row_num * col_num;
  n = row_num * col_num;
  a = neumann ( row_num, col_num );
  x = neumann_null_right ( row_num, col_num );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ONE
*/
  strcpy ( title, "ONE" );
  m = 5;
  n = 5;
  a = one ( n, n );
  x = one_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  RING_ADJ
  N must be a multiple of 4 for there to be a null vector.
*/
  strcpy ( title, "RING_ADJ" );
  m = 12;
  n = 12;
  a = ring_adj ( n );
  x = ring_adj_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ROSSER1
*/
  strcpy ( title, "ROSSER1" );
  m = 8;
  n = 8;
  a = rosser1 ( );
  x = rosser1_null_right ( );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );
/*
  ZERO
*/
  strcpy ( title, "ZERO" );
  m = 5;
  n = 5;
  a = zero ( m, n );
  x = zero_null_right ( m, n );
  error_l2 = r8mat_is_null_right ( m, n, a, x );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  norm_x_l2 = r8vec_norm_l2 ( n, x );
  printf ( "  %-20s  %4d  %4d  %14g  %14g  %14g\n",
    title, m, n, norm_a_frobenius, norm_x_l2, error_l2 );
  free ( a );
  free ( x );

  return;
}
/******************************************************************************/

void test_plu ( )

/******************************************************************************/
/*
  Purpose:

    TEST_PLU tests the PLU factors.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    07 April 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double error_frobenius;
  int i;
  int i4_hi;
  int i4_lo;
  double *l;
  int m;
  int n;
  double norm_a_frobenius;
  double *p;
  int *pivot;
  double r8_hi;
  double r8_lo;
  int seed;
  char title[21];
  double *u;
  double *x;

  printf ( "\n" );
  printf ( "TEST_PLU\n" );
  printf ( "  A = a test matrix of order M by N\n" );
  printf ( "  P, L, U are the PLU factors.\n" );
  printf ( "\n" );
  printf ( "  ||A|| = Frobenius norm of A.\n" );
  printf ( "  ||A-PLU|| = Frobenius norm of A-P*L*U.\n" );
  printf ( "\n" );
  printf ( "  Title                    M     N  " );
  printf ( "    ||A||        ||A-PLU||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  m = 3;
  n = 3;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = a123 ( );
  a123_plu ( p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  BODEWIG
*/
  strcpy ( title, "BODEWIG" );
  m = 4;
  n = 4;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = bodewig ( );
  bodewig_plu ( p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  BORDERBAND
*/
  strcpy ( title, "BORDERBAND" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = borderband ( n );
  borderband_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = dif2 ( m, n );
  dif2_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  GFPP
*/
  strcpy ( title, "GFPP" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = gfpp ( n, alpha );
  gfpp_plu ( n, alpha, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  GIVENS
*/
  strcpy ( title, "GIVENS" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = givens ( n, n );
  givens_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  KMS
*/
  strcpy ( title, "KMS" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = kms ( alpha, m, n );
  kms_plu ( alpha, n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  LEHMER
*/
  strcpy ( title, "LEHMER" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = lehmer ( n, n );
  lehmer_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  MAXIJ
*/
  strcpy ( title, "MAXIJ" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = maxij ( n, n );
  maxij_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  MINIJ
*/
  strcpy ( title, "MINIJ" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = minij ( m, n );
  minij_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  MOLER1
*/
  strcpy ( title, "MOLER1" );
  m = 5;
  n = 5;
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  alpha = r8_uniform_ab ( r8_lo, r8_hi, &seed );
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = moler1 ( alpha, m, n );
  moler1_plu ( alpha, n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  MOLER3
*/
  strcpy ( title, "MOLER3" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = moler3 ( m, n );
  moler3_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  OTO
*/
  strcpy ( title, "OTO" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = oto ( m, n );
  oto_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  PASCAL2
*/
  strcpy ( title, "PASCAL2" );
  m = 5;
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = pascal2 ( n );
  pascal2_plu ( n, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
/*
  PLU
*/
  strcpy ( title, "PLU" );
  n = 5;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  pivot = ( int * ) malloc ( n * sizeof ( int ) );
  seed = 123456789;
  for ( i = 0; i < n; i++ )
  {
    i4_lo = i;
    i4_hi = n - 1;
    pivot[i] = i4_uniform_ab ( i4_lo, i4_hi, &seed );
  }
  a = plu ( n, pivot );
  plu_plu ( n, pivot, p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( pivot );
  free ( u );
/*
  VAND2
*/
  strcpy ( title, "VAND2" );
  m = 4;
  n = 4;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  r8_lo = -5.0;
  r8_hi = +5.0;
  seed = 123456789;
  x = r8vec_uniform_ab_new ( m, r8_lo, r8_hi, &seed );
  a = vand2 ( m, x );
  vand2_plu ( m, x, p, l, u  );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );
  free ( x );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  m = 4;
  n = 4;
  p = ( double * ) malloc ( m * m * sizeof ( double ) );
  l = ( double * ) malloc ( m * m * sizeof ( double ) );
  u = ( double * ) malloc ( m * n * sizeof ( double ) );
  a = wilson ( );
  wilson_plu ( p, l, u );
  error_frobenius = r8mat_is_plu ( m, n, a, p, l, u );
  norm_a_frobenius = r8mat_norm_fro ( m, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_a_frobenius, error_frobenius );
  free ( a );
  free ( l );
  free ( p );
  free ( u );

  return;
}
/******************************************************************************/

void test_solution ( )

/******************************************************************************/
/*
  Purpose:

    TEST_SOLUTION tests the linear solution computations.

  Licensing:

    This code is distributed under the GNU LGPL license. 

  Modified:

    09 March 2015

  Author:

    John Burkardt
*/
{
  double *a;
  double alpha;
  double *b;
  double beta;
  double error_frobenius;
  double gamma;
  int i1;
  int k;
  int m;
  int n;
  int ncol;
  double norm_frobenius;
  int nrow;
  int seed;
  char title[21];
  double *x;

  printf ( "\n" );
  printf ( "TEST_SOLUTION\n" );
  printf ( "  Compute the Frobenius norm of the solution error:\n" );
  printf ( "    A * X - B\n" );
  printf ( "  given MxN matrix A, NxK solution X, MxK right hand side B.\n" );
  printf ( "\n" );
  printf ( "  Title                    M     N     K      ||A||         ||A*X-B||\n" );
  printf ( "\n" );
/*
  A123
*/
  strcpy ( title, "A123" );
  m = 3;
  n = 3;
  k = 1;
  a = a123 ( );
  b = a123_rhs ( );
  x = a123_solution ( );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  BODEWIG
*/
  strcpy ( title, "BODEWIG" );
  m = 4;
  n = 4;
  k = 1;
  a = bodewig ( );
  b = bodewig_rhs ( );
  x = bodewig_solution ( );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  DIF2
*/
  strcpy ( title, "DIF2" );
  m = 10;
  n = 10;
  k = 2;
  a = dif2 ( m, n );
  b = dif2_rhs ( m, k );
  x = dif2_solution ( n, k );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  FRANK
*/
  strcpy ( title, "FRANK" );
  m = 10;
  n = 10;
  k = 2;
  a = frank ( n );
  b = frank_rhs ( m, k );
  x = frank_solution ( n, k );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  POISSON
*/
  strcpy ( title, "POISSON" );
  nrow = 4;
  ncol = 5;
  m = nrow * ncol;
  n = nrow * ncol;
  k = 1;
  a = poisson ( nrow, ncol );
  b = poisson_rhs ( nrow, ncol );
  x = poisson_solution ( nrow, ncol );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  WILK03
*/
  strcpy ( title, "WILK03" );
  m = 3;
  n = 3;
  k = 1;
  a = wilk03 ( );
  b = wilk03_rhs ( );
  x = wilk03_solution ( );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  WILK04
*/
  strcpy ( title, "WILK04" );
  m = 4;
  n = 4;
  k = 1;
  a = wilk04 ( );
  b = wilk04_rhs ( );
  x = wilk04_solution ( );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );
/*
  WILSON
*/
  strcpy ( title, "WILSON" );
  m = 4;
  n = 4;
  k = 1;
  a = wilson ( );
  b = wilson_rhs ( );
  x = wilson_solution ( );
  error_frobenius = r8mat_is_solution ( m, n, k, a, x, b );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %4d  %14g  %14g\n",
    title, m, n, k, norm_frobenius, error_frobenius );
  free ( a );
  free ( b );
  free ( x );

  return;
}
/******************************************************************************/

void test_type ( )

/******************************************************************************/
/*
  Purpose:

    TEST_TYPE tests functions which test the type of a matrix.

  Licensing:

    This code is distributed under the GNU LGPL license.

  Modified:

    15 July 2013

  Author:

    John Burkardt
*/
{
  double *a;
  double error_frobenius;
  int key;
  int m;
  int n;
  double norm_frobenius;
  int seed;
  char title[21];

  printf ( "\n" );
  printf ( "TEST_TYPE\n" );
  printf ( "  Demonstrate functions which test the type of a matrix.\n" );
/*
  TRANSITION.
*/
  printf ( "\n" );
  printf ( "  Title                    M     N     ||A||" );
  printf ( "            ||Transition Error||\n" );
  printf ( "\n" );

  strcpy ( title, "BODEWIG" );
  m = 4;
  n = 4;
  a = bodewig ( );
  error_frobenius = r8mat_is_transition ( m, n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_frobenius, error_frobenius );
  free ( a );

  strcpy ( title, "SNAKES" );
  m = 101;
  n = 101;
  a = snakes ( );
  error_frobenius = r8mat_is_transition ( m, n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_frobenius, error_frobenius );
  free ( a );

  strcpy ( title, "TRANSITION_RANDOM" );
  m = 5;
  n = 5;
  key = 123456789;
  a = transition_random ( n, key );
  error_frobenius = r8mat_is_transition ( m, n, a );
  norm_frobenius = r8mat_norm_fro ( n, n, a );
  printf ( "  %-20s  %4d  %4d  %14g  %14g\n",
    title, m, n, norm_frobenius, error_frobenius );
  free ( a );

  return;
}
