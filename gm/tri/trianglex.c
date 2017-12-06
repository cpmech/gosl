// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdlib.h> // for malloc and free
#include "triangle.h"

void tiosetnull(struct triangulateio * t) {
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

void tiofree(struct triangulateio * t) {
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
    tiosetnull(t);
}

void tioalloc(struct triangulateio *t, int npoints, int nsegments, int nregions, int nholes) {

    // points
    t->pointlist       = (double*)malloc(npoints*2*sizeof(double));
    t->numberofpoints  = npoints;
    t->pointmarkerlist = (int*)malloc(npoints*sizeof(int));

    // segments
    t->segmentlist       = (int*)malloc(nsegments*2*sizeof(int));
    t->segmentmarkerlist = (int*)malloc(nsegments * sizeof(int));
    t->numberofsegments  = nsegments;
    int i;
    for (i=0; i<nsegments; ++i) {
        t->segmentmarkerlist[i]=0;
    }

    // regions
    if (nregions>0) {
        t->regionlist      = (double*)malloc(nregions*4*sizeof(double));
        t->numberofregions = nregions;
    }

    // holes
    if (nholes>0) {
        t->holelist      = (double*)malloc(nholes*2*sizeof(double));
        t->numberofholes = nholes;
    }
}

int delaunay2d(struct triangulateio *out, int npoints, double *X, double *Y, int verbose) {

    // input structure
    struct triangulateio tin;
    tiosetnull(&tin);

    // set points
	tin.pointlist = (double*)malloc(npoints*2*sizeof(double));
	tin.numberofpoints = npoints;
    int i;
    for (i=0; i<npoints; ++i) {
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
    tiofree(&tin);
    return 0; // success
}

double getpoint(int pointId, int dimIdx, struct triangulateio *T) {
    return T->pointlist[pointId*2 + dimIdx];
}

int ptmap[] = {0,1,2,5,3,4};

int getcorner(int cellId, int pointIdx, struct triangulateio *T) {
    return T->trianglelist[cellId*T->numberofcorners + ptmap[pointIdx]];
}

int getcelltag(int cellId, struct triangulateio *T) {
    return T->triangleattributelist[cellId*T->numberoftriangleattributes];
}

int getedgetag(int cellId, int edgeIdx, struct triangulateio *T) {
    return T->triedgemarks[cellId*3+edgeIdx];
}

void setpoint(struct triangulateio *t, int i, int tag, double x, double y) {
    t->pointlist[i*2  ]   = x;
    t->pointlist[i*2+1]   = y;
    t->pointmarkerlist[i] = tag;
}

void setsegment(struct triangulateio *t, int iSeg, int tag, int l, int r) {
    t->segmentlist[iSeg*2  ]   = l;
    t->segmentlist[iSeg*2+1]   = r;
    t->segmentmarkerlist[iSeg] = tag;
}

void setregion(struct triangulateio *t, int i, int tag, double maxarea, double x, double y) {
    t->regionlist[i*4  ] = x;
    t->regionlist[i*4+1] = y;
    t->regionlist[i*4+2] = tag;
    t->regionlist[i*4+3] = maxarea;
}

void sethole(struct triangulateio *t, int i, double x, double y) {
    t->holelist[i*2  ] = x;
    t->holelist[i*2+1] = y;
}
