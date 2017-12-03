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

void tiosetnull(struct triangulateio *t);

void tiofree(struct triangulateio *t);

void tioalloc(struct triangulateio *t, int npoints, int nsegments, int nregions, int nholes);

void triangulate(char const *switches, struct triangulateio *in, struct triangulateio *out,
    struct triangulateio *vorout);

int delaunay2d(struct triangulateio *out, int npoints, double *X, double *Y, int verbose);

double getpoint(int pointId, int dimIdx, struct triangulateio *t);

int getcorner(int cellId, int pointIdx, struct triangulateio *T);

int getcelltag(int cellId, struct triangulateio *T);

int getedgetag(int cellId, int edgeIdx, struct triangulateio *T);

void setpoint(struct triangulateio *t, int i, int tag, double x, double y);

void setsegment(struct triangulateio *t, int iSeg, int tag, int l, int r);

void setregion(struct triangulateio *t, int i, int tag, double maxarea, double x, double y);

void sethole(struct triangulateio *t, int i, double x, double y);

#endif /* TRIANGLE_H_INCLUDED */
