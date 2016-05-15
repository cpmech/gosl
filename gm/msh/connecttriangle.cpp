// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdlib.h> // for malloc and free
#include "triangle.h"

extern "C" {

void trisetnull(triangulateio * t) {
    // points
    t->pointlist               = NULL;
    t->pointattributelist      = NULL;
    t->pointmarkerlist         = NULL;
    t->numberofpoints          = 0;
    t->numberofpointattributes = 0;

    // triangles
    t->trianglelist               = NULL;
    t->triangleattributelist      = NULL;
    t->trianglearealist           = NULL;
    t->neighborlist               = NULL;
    t->numberoftriangles          = 0;
    t->numberofcorners            = 0;
    t->numberoftriangleattributes = 0;
    t->triedgemarks               = NULL;

    // segments
    t->segmentlist       = NULL;
    t->segmentmarkerlist = NULL;
    t->numberofsegments  = 0;

    // holes
    t->holelist      = NULL;
    t->numberofholes = 0;

    // regions
    t->regionlist      = NULL;
    t->numberofregions = 0;

    // edges
    t->edgelist       = NULL;
    t->edgemarkerlist = NULL;
    t->normlist       = NULL;
    t->numberofedges  = 0;
}

void trifree(triangulateio * t) {
    // points
    if (t->pointlist          != NULL) free(t->pointlist);
    if (t->pointattributelist != NULL) free(t->pointattributelist);
    if (t->pointmarkerlist    != NULL) free(t->pointmarkerlist);

    // triangles
    if (t->trianglelist          != NULL) free(t->trianglelist);
    if (t->triangleattributelist != NULL) free(t->triangleattributelist);
    if (t->trianglearealist      != NULL) free(t->trianglearealist);
    if (t->neighborlist          != NULL) free(t->neighborlist);
    if (t->triedgemarks          != NULL) free(t->triedgemarks);

    // segments
    if (t->segmentlist       != NULL) free(t->segmentlist);
    if (t->segmentmarkerlist != NULL) free(t->segmentmarkerlist);

    // holes
    if (t->holelist != NULL) free(t->holelist);

    // regions
    if (t->regionlist != NULL) free(t->regionlist);

    // edges
    if (t->edgelist       != NULL) free(t->edgelist);
    if (t->edgemarkerlist != NULL) free(t->edgemarkerlist);
    if (t->normlist       != NULL) free(t->normlist);

    // clear all
    trisetnull(t);
}

long delaunay2d(triangulateio *out, long npoints, double *X, double *Y, long verbose) {

    // input structure
    triangulateio tin;
    trisetnull(&tin);

    // set points
	tin.pointlist = (double*)malloc(npoints*2*sizeof(double));
	tin.numberofpoints = npoints;
    for (long i=0; i<npoints; ++i) {
        tin.pointlist[0+i*2] = X[i];
        tin.pointlist[1+i*2] = Y[i];
    }

    // triangulate
    if (verbose == 1) {
        triangulate("z", &tin, out, NULL); // zero-based
    } else {
        triangulate("Qz", &tin, out, NULL); // quiet, zero-based
    }

    // clean up
    trifree(&tin);
    return 0; // success
}

} // extern "C"
