// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_SPHERES_H
#define GOSLVTK_SPHERES_H

// VTK
#include <vtkPoints.h>
#include <vtkPointData.h>
#include <vtkDoubleArray.h>
#include <vtkSphereSource.h>
#include <vtkGlyph3D.h>
#include <vtkPolyDataMapper.h>
#include <vtkLODActor.h>
#include <vtkProperty.h>
#include <vtkLookupTable.h>
#include <vtkTextActor3D.h>

// GoslVTK
#include "linalg.h"
#include "v_win.h"

namespace GoslVTK
{

class Spheres
{
public:
    // Constructor & Destructor
     Spheres ();
    ~Spheres ();

    // Init method
    void Init (int N, double * X, double * Y, double * Z, double * R, int ThetaRes=20, int PhiRes=20);

    // Set methods
    void SetResolution (int ThetaRes=20, int PhiRes=20);
    void SetColor      (double r, double g, double b, double Opacity);

    // Methods
    void AddTo (GoslVTK::Win & win, bool RstCam=true) { if (_initialized) win.AddActor(_spheres_actor, RstCam); }

private:
    // essential data
    bool _initialized;

    // vtk data
    vtkPoints              * _points;
    vtkDoubleArray         * _scalars;
    vtkSphereSource        * _sphere;
    vtkGlyph3D             * _spheres;
    vtkPolyDataMapper      * _spheres_mapper;
    vtkLODActor            * _spheres_actor;
    vtkLookupTable         * _ltable;
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


Spheres::Spheres() : _initialized(false) {
    _points         = NULL;
    _scalars        = NULL;
    _sphere         = NULL;
    _spheres        = NULL;
    _spheres_mapper = NULL;
    _spheres_actor  = NULL;
    _ltable         = NULL;
}

Spheres::~Spheres () {
    if (_points         != NULL) _points         -> Delete();
    if (_scalars        != NULL) _scalars        -> Delete();
    if (_sphere         != NULL) _sphere         -> Delete();
    if (_spheres        != NULL) _spheres        -> Delete();
    if (_spheres_mapper != NULL) _spheres_mapper -> Delete();
    if (_spheres_actor  != NULL) _spheres_actor  -> Delete();
    if (_ltable         != NULL) _ltable         -> Delete();
    //printf("v_spheres.h: Destructor: all deleted\n");
}

void Spheres::Init(int N, double * X, double * Y, double * Z, double * R, int ThetaRes, int PhiRes) {

    // check
    if (N < 1) return ;
    if (_initialized) {
        printf("Warning: 'Spheres' was already initialised\n");
        return;
    }

    // points and scalars
    _points  = vtkPoints      ::New();
    _scalars = vtkDoubleArray ::New();
    _scalars -> SetNumberOfComponents (1);

    // polydata
    vtkPolyData * polydata = vtkPolyData   ::New();
    polydata -> SetPoints                  (_points);
    polydata -> GetPointData()->SetScalars (_scalars);

    // spheres
    _sphere         = vtkSphereSource              ::New();
    _spheres        = vtkGlyph3D                   ::New();
    _spheres_mapper = vtkPolyDataMapper            ::New();
    _spheres_actor  = vtkLODActor                  ::New();
    _ltable         = vtkLookupTable               ::New();
    _spheres        -> SetInputData                (polydata);
    _spheres        -> SetSourceConnection         (_sphere->GetOutputPort());
    _spheres        -> SetScaleModeToScaleByScalar ();
    _spheres        -> SetColorModeToColorByScalar ();
    _spheres        -> SetScaleFactor              (1.0);
    _spheres_mapper -> SetInputConnection          (_spheres->GetOutputPort());
    _spheres_mapper -> SetLookupTable              (_ltable);
    _spheres_actor  -> SetMapper                   (_spheres_mapper);
    SetColor (0.8, 0.6, 0.4, 1.0);

    // flag initialised
    _initialized = true;

    // resolution
    SetResolution(ThetaRes, PhiRes);

    // set points
    _points  -> SetNumberOfPoints (N);
    _scalars -> SetNumberOfTuples (N);
    for (size_t i=0; i<N; ++i) {
        _points -> InsertPoint (i, X[i], Y[i], Z[i]);
        if (R==NULL) _scalars -> InsertTuple1 (i, 1.0);
        else         _scalars -> InsertTuple1 (i, 2.0*R[i]);
    }
}

void Spheres::SetResolution(int ThetaRes, int PhiRes) {

    if (!_initialized) return;

    _sphere -> SetThetaResolution (ThetaRes);
    _sphere -> SetPhiResolution   (PhiRes);
}

void Spheres::SetColor(double r, double g, double b, double Opacity) {

    if (!_initialized) return;

    _ltable->SetNumberOfColors (2);
    _ltable->Build             ();
    _ltable->SetTableValue     (0, r, g, b);
    _ltable->SetTableValue     (1, r, g, b);
    _spheres_actor->GetProperty()->SetOpacity (Opacity);
}

}; // namespace GoslVTK

#endif // GOSLVTK_SPHERE_H
