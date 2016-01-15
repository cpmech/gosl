// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CONNECTGOVTK_H
#define CONNECTGOVTK_H

#ifdef __cplusplus
extern "C" {
#endif

double *GOVTK_F;  // scalar
double *GOVTK_VX; // scalar
double *GOVTK_VY; // scalar
double *GOVTK_VZ; // scalar
double *GOVTK_X;  // vector [3];
long   *GOVTK_I;  // scalar

// deallocate memory
void win_dealloc(void * win);
void arrow_dealloc(void * arrow);
void sphere_dealloc(void * sphere);
void spheres_dealloc(void * sset);
void isosurf_dealloc(void * isosurf);

// return non-NULL pointer on success
void * win_alloc(long reverse);

// return 0 on success
int scene_run(
    void       * input_win,
    double       axeslen,
    long         hydroline,
    long         reverse,
    long         fullaxes,
    long         withplanes,
    long         interact,
    long         saveonexit,
    char const * fnk,
    char const * lblX,
    char const * lblY,
    char const * lblZ,
    long         lblSz,
    double     * lblClr);

// return non-NULL pointer on success
void * arrow_addto(
    void   * input_win,
    double * x0,
    double * v,
    double   cone_pct,
    double   cone_rad,
    double   cyli_rad,
    long     resolution,
    double * color);

// return non-NULL pointer on success
void * sphere_addto(
    void   * input_win,
    double * cen,
    double   r,
    double * color);

// return non-NULL pointer on success
void * spheres_addto(
    void   * input_win,
    long     nspheres,
    double * x,
    double * y,
    double * z,
    double * r,
    double * color);

// return non-NULL pointer on success
void * isosurf_addto(
    void       * input_win,
    long         index,
    double     * limits,
    long       * ndiv,
    double     * frange,
    long         octrotate,
    long         nlevels,
    char const * cmaptype,
    long         cmapnclrs,
    long         cmaprangetype,
    double     * cmapfrange,
    double     * color,
    long         showwire,
    long         gridshowpts);

#ifdef __cplusplus
} /* extern "C" */
#endif

#endif // CONNECTGOVTK_H
