// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CONNECTTRIANGLE_H
#define CONNECTTRIANGLE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
  double *pointlist;                                             /* In / out */
  double *pointattributelist;                                    /* In / out */
  int *pointmarkerlist;                                          /* In / out */
  int numberofpoints;                                            /* In / out */
  int numberofpointattributes;                                   /* In / out */

  int *trianglelist;                                             /* In / out */
  double *triangleattributelist;                                 /* In / out */
  double *trianglearealist;                                       /* In only */
  int *neighborlist;                                             /* Out only */
  int numberoftriangles;                                         /* In / out */
  int numberofcorners;                                           /* In / out */
  int numberoftriangleattributes;                                /* In / out */
  int *triedgemarks;                     /* Out (size = 3*numberoftriangles) */

  int *segmentlist;                                              /* In / out */
  int *segmentmarkerlist;                                        /* In / out */
  int numberofsegments;                                          /* In / out */

  double *holelist;                      /* In / pointer to array copied out */
  int numberofholes;                                      /* In / copied out */

  double *regionlist;                    /* In / pointer to array copied out */
  int numberofregions;                                    /* In / copied out */

  int *edgelist;                                                 /* Out only */
  int *edgemarkerlist;            /* Not used with Voronoi diagram; out only */
  double *normlist;              /* Used only with Voronoi diagram; out only */
  int numberofedges;                                             /* Out only */
} triangulateio;

long delaunay2d(triangulateio *out, long npoints, double *X, double *Y, long verbose);

void trifree(triangulateio * t);

double getpoint(long pointId, long dimIdx, triangulateio *T) {
    return T->pointlist[pointId*2 + dimIdx];
}

long ptmap[] = {0,1,2,5,3,4};

long getcorner(long cellId, long pointIdx, triangulateio *T) {
    return T->trianglelist[cellId*T->numberofcorners + ptmap[pointIdx]];
}

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif // CONNECTTRIANGLE_H
