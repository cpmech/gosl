// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef TRIANGLE_H_INCLUDED
#define TRIANGLE_H_INCLUDED

#define REAL double
#define VOID int

struct triangulateio {
  REAL *pointlist;                                               /* In / out */
  REAL *pointattributelist;                                      /* In / out */
  int *pointmarkerlist;                                          /* In / out */
  int numberofpoints;                                            /* In / out */
  int numberofpointattributes;                                   /* In / out */

  int *trianglelist;                                             /* In / out */
  REAL *triangleattributelist;                                   /* In / out */
  REAL *trianglearealist;                                         /* In only */
  int *neighborlist;                                             /* Out only */
  int numberoftriangles;                                         /* In / out */
  int numberofcorners;                                           /* In / out */
  int numberoftriangleattributes;                                /* In / out */
  int *triedgemarks;                     /* Out (size = 3*numberoftriangles) */

  int *segmentlist;                                              /* In / out */
  int *segmentmarkerlist;                                        /* In / out */
  int numberofsegments;                                          /* In / out */

  REAL *holelist;                        /* In / pointer to array copied out */
  int numberofholes;                                      /* In / copied out */

  REAL *regionlist;                      /* In / pointer to array copied out */
  int numberofregions;                                    /* In / copied out */

  int *edgelist;                                                 /* Out only */
  int *edgemarkerlist;            /* Not used with Voronoi diagram; out only */
  REAL *normlist;                /* Used only with Voronoi diagram; out only */
  int numberofedges;                                             /* Out only */
};

void tiofree(struct triangulateio * t);

void triangulate(char const *switches, struct triangulateio *in, struct triangulateio *out,
    struct triangulateio *vorout);

long delaunay2d(struct triangulateio *out, long npoints, double *X, double *Y, long verbose);

double getpoint(long pointId, long dimIdx, struct triangulateio *T);

long getcorner(long cellId, long pointIdx, struct triangulateio *T);

#endif /* TRIANGLE_H_INCLUDED */
