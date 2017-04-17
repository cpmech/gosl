// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <cstdlib>
#include <cstdio>
#include <cmath>
#include <vector>

#include "clrmap.h"
#include "connectgovtk.h"
#include "fatal.h"
#include "linalg.h"
#include "mystring.h"
#include "tsrinvars.h"
#include "util.h"
#include "v_arrow.h"
#include "v_axes.h"
#include "v_isosurf.h"
#include "v_plane.h"
#include "v_sgrid.h"
#include "v_sphere.h"
#include "v_spheres.h"
#include "v_win.h"

extern "C" {

// deallocate memory
void win_dealloc(void * input_win) {
    GoslVTK::Win * win = (GoslVTK::Win*) input_win;
    if (win != NULL) delete win;
}

// deallocate memory
void arrow_dealloc(void * input_arrow) {
    GoslVTK::Arrow * arrow = (GoslVTK::Arrow*) input_arrow;
    if (arrow != NULL) delete arrow;
}

// deallocate memory
void sphere_dealloc(void * input_sphere) {
    GoslVTK::Sphere * sphere = (GoslVTK::Sphere*) input_sphere;
    if (sphere != NULL) delete sphere;
}

// deallocate memory
void spheres_dealloc(void * input_sset) {
    GoslVTK::Spheres * sset = (GoslVTK::Spheres*) input_sset;
    if (sset != NULL) delete sset;
}

// deallocate memory
void isosurf_dealloc(void * input_isosurf) {
    GoslVTK::IsoSurf * isosurf = (GoslVTK::IsoSurf*) input_isosurf;
    if (isosurf != NULL) delete isosurf;
}

// returns non-NULL pointer on success
void * win_alloc(long width, long height, long reverse) {
    try {

        // window
        GoslVTK::Win * win = new GoslVTK::Win();
        win->Init(width, height);
        win->SetViewDefault(reverse>0);

        // success
        return win;
    }

    // fail
    GOSLVTK_CATCH;
    return NULL;
}

// returns 0 on success
int set_camera(void * input_win, double * data) {
    try {
        GoslVTK::Win * win = (GoslVTK::Win*) input_win;
        win->Camera(data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8]);
        return 0; // success
    }
    GOSLVTK_CATCH;
    return 1; // fail
}

// returns 0 on success
int scene_run(
    void       * input_win,
    double       axeslen,
    long         hydroline,
    long         reverse,
    long         fullaxes,
    long         withplanes,
    long         interact,
    long         saveeps,
    long         savepng,
    long         pngmag,
    char const * fnk,
    char const * lblX,
    char const * lblY,
    char const * lblZ,
    long         lblSz,
    double     * lblClr,
    double       zoom) {

    try {

        // window
        GoslVTK::Win * win = (GoslVTK::Win*) input_win;

        // axes
        GoslVTK::Axes axe;
        axe.Init(axeslen, hydroline>0, reverse>0, fullaxes>0);
        axe.SetLabels(lblX, lblY, lblZ, lblClr[0], lblClr[1], lblClr[2], lblSz);
        axe.AddTo(*win);

        // auxiliary planes
        GoslVTK::Plane pxy, pyz, pzx;
        if (withplanes > 0) {
            double al = axeslen;
            double dat[][4][3] = {
                { {-al,-al,0}, {al,-al,0}, {-al,al,0}, {0,0,1} }, // ori, pt1, pt2, normal
                { {0,-al,-al}, {0,al,-al}, {0,-al,al}, {1,0,0} },
                { {-al,0,-al}, {-al,0,al}, {al,0,-al}, {0,1,0} },
            };
            pxy.Init(dat[0][0], dat[0][1], dat[0][2], dat[0][3]);
            pyz.Init(dat[1][0], dat[1][1], dat[1][2], dat[1][3]);
            pzx.Init(dat[2][0], dat[2][1], dat[2][2], dat[2][3]);
            pxy.SetColor(1, 0.5, 0, 0.05);
            pyz.SetColor(1, 0.5, 0, 0.05);
            pzx.SetColor(1, 0.5, 0, 0.05);
            pxy.AddTo(*win);
            pyz.AddTo(*win);
            pzx.AddTo(*win);
        }

        // set zoom
        win->Zoom(zoom);

        // interact
        if (interact > 0) {
            win->Show();
        }

        // save figure
        if (savepng > 0) {
            bool large = (pngmag > 0 ? true : false);
            win->WritePNG(fnk, large, pngmag);
        }
        if (saveeps > 0) {
            win->WriteEPS(fnk);
        }

        // success
        return 0;
    }

    // fail
    GOSLVTK_CATCH;
    return 1;
}

// returns non-NULL pointer on success
void * arrow_addto(
    void   * input_win,
    double * x0,
    double * v,
    double   cone_pct,
    double   cone_rad,
    double   cyli_rad,
    long     resolution,
    double * color) {

    try {

        // window
        GoslVTK::Win * win = (GoslVTK::Win*) input_win;

        // arrow
        GoslVTK::Arrow * arr = new GoslVTK::Arrow();
        arr->Init(x0, v, cone_pct, cone_rad, cyli_rad, resolution);
        arr->AddTo(*win);
        arr->SetColor(color[0], color[1], color[2], color[3]);

        // success
        return arr;
    }

    // fail
    GOSLVTK_CATCH;
    return NULL;
}

// returns non-NULL pointer on success
void * sphere_addto(
    void   * input_win,
    double * cen,
    double   r,
    double * color) {

    try {

        // window
        GoslVTK::Win * win = (GoslVTK::Win*) input_win;

        // sphere
        GoslVTK::Sphere * sph = new GoslVTK::Sphere();
        sph->Init(cen, r);
        sph->AddTo(*win);
        sph->SetColor(color[0], color[1], color[2], color[3]);

        // success
        return sph;
    }

    // fail
    GOSLVTK_CATCH;
    return NULL;
}

// returns non-NULL pointer on success
void * spheres_addto(
    void   * input_win,
    long     nspheres,
    double * x,
    double * y,
    double * z,
    double * r,
    double * color) {

    try {

        // window
        GoslVTK::Win * win = (GoslVTK::Win*) input_win;

        // spheres
        GoslVTK::Spheres * S = new GoslVTK::Spheres();
        S->Init(nspheres, x, y, z, r);
        S->AddTo(*win);
        S->SetColor(color[0], color[1], color[2], color[3]);

        // success
        return S;
    }

    // fail
    GOSLVTK_CATCH;
    return NULL;
}

// call Go
void isosurf_fcn(double & f, double v[3], const double x[3], int index) {

    // set global data
    GOVTK_X[0] = x[0];
    GOVTK_X[1] = x[1];
    GOVTK_X[2] = x[2];
    *GOVTK_I   = index;

    // call go callback function
    extern void govtk_isosurf_fcn();
    govtk_isosurf_fcn();

    // read global data
    f    = *GOVTK_F;
    v[0] = *GOVTK_VX;
    v[1] = *GOVTK_VY;
    v[2] = *GOVTK_VZ;
}

// returns non-NULL pointer on success
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
    long         gridshowpts) {

    try {

        // window
        GoslVTK::Win * win = (GoslVTK::Win*) input_win;

        // isosurface
        int N[3] = {ndiv[0], ndiv[1], ndiv[2]};
        GoslVTK::IsoSurf * isf = new GoslVTK::IsoSurf();
        isf->Init(N, limits, &isosurf_fcn, index, octrotate>0);
        isf->AddTo(*win);

        // structured grid
        GoslVTK::SGrid * grd = isf->GetSGrid();
        if (gridshowpts > 0) {
            grd->ShowPoints();
            grd->AddTo(*win);
        }

        // levels
        if (nlevels == 0 || nlevels == 1) {
            isf->SetValue(frange[0]);
        } else {
            if (abs(frange[1] - frange[0]) > 1e-10) {
                isf->GenValues(nlevels, frange[0], frange[1]);
            } else {
                isf->GenValues(nlevels);
            }
        }
        if (showwire > 0) {
            isf->SetWire();
        }

        // set colors and opacity
        if (cmapnclrs > 0) {
            switch (cmaprangetype) {
            case 1: // use Frange
                isf->SetCmap(cmaptype, cmapnclrs, frange[0], frange[1]);
                break;
            case 2: // use CmapFrange
                isf->SetCmap(cmaptype, cmapnclrs, cmapfrange[0], cmapfrange[1]);
                break;
            default: // use sgrid range values (automatic)
                isf->SetCmap(cmaptype, cmapnclrs);
                break;
            }

        // fixed color
        } else {
            isf->SetColor(color[0], color[1], color[2], color[3]);
        }
        isf->SetOpac(color[3]);

        // success
        return isf;
    }

    // fail
    GOSLVTK_CATCH;
    return NULL;
}

} // extern "C"
