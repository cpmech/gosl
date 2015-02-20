// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_ISOSURF_H
#define GOSLVTK_ISOSURF_H

// VTK
#include <vtkMarchingContourFilter.h>
#include <vtkHedgeHog.h>
#include <vtkDataSetMapper.h>
#include <vtkPolyDataMapper.h>
#include <vtkLookupTable.h>
#include <vtkProperty.h>

// GoslVTK
#include "clrmap.h"
#include "v_sgrid.h"
#include "v_win.h"

namespace GoslVTK
{

class IsoSurf
{
public:
    // Constructor & Destructor
     IsoSurf ();
    ~IsoSurf ();

    // Init method
    void Init (int N[3], double L[6], GridCallBack Func, int index, bool octrotate);

    // Set methods
    void SetColor    (double r, double g, double b, double Opacity);
    void SetOpac     (double Opacity)                       { if (_initialized) _isosurf_actor->GetProperty()->SetOpacity(Opacity); }
    void SetValue    (double F)                             { if (_initialized) _isosurf->SetValue(0, F); }
    void GenValues   (int NSurfs, double fMin, double fMax) { if (_initialized) _isosurf->GenerateValues(NSurfs,fMin,fMax); }
    void GenValues   (int NSurfs);
    void SetVecScale (double Factor) { if (_initialized) _hedgehog->SetScaleFactor(Factor); }
    void SetWire     ()              { if (_initialized) _isosurf_actor->GetProperty()->SetRepresentationToWireframe(); }
    void SetCmap     (char const *Name, int NumClrs);
    void SetCmap     (char const *Name, int NumClrs, double Fmin, double Fmax);
    void SetVecCmap  (char const *Name, int NumClrs);
    void SetVecCmap  (char const *Name, int NumClrs, double Fmin, double Fmax);

    // Extra methods
    void SetMaterial (double Ambient, double Diffuse, double Specular, double SpecularPower);

    // Access methods
    SGrid * GetSGrid() { if (_initialized) return _sgrid; return NULL; }

    // Methods
    void AddTo (GoslVTK::Win & win);

    // Data
    bool ShowIsoSurf;
    bool ShowVectors;

private:
    // essential data
    bool _initialized;

    // vtk data
    GoslVTK::SGrid           * _sgrid;
    vtkMarchingContourFilter * _isosurf;
    vtkPolyDataMapper        * _isosurf_mapper;
    vtkActor                 * _isosurf_actor;
    vtkLookupTable           * _isosurf_lt;
    vtkHedgeHog              * _hedgehog;
    vtkPolyDataMapper        * _hedgehog_mapper;
    vtkActor                 * _hedgehog_actor;
    vtkLookupTable           * _hedgehog_lt;
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


IsoSurf::IsoSurf() : _initialized(false) {
    _sgrid           = NULL;
    _isosurf         = NULL;
    _isosurf_mapper  = NULL;
    _isosurf_actor   = NULL;
    _isosurf_lt      = NULL;
    _hedgehog        = NULL;
    _hedgehog_mapper = NULL;
    _hedgehog_actor  = NULL;
    _hedgehog_lt     = NULL;
}

IsoSurf::~IsoSurf() {
    if (_sgrid           != NULL) delete _sgrid;
    if (_isosurf         != NULL) _isosurf         -> Delete();
    if (_isosurf_mapper  != NULL) _isosurf_mapper  -> Delete();
    if (_isosurf_actor   != NULL) _isosurf_actor   -> Delete();
    if (_isosurf_lt      != NULL) _isosurf_lt      -> Delete();
    if (_hedgehog        != NULL) _hedgehog        -> Delete();
    if (_hedgehog_mapper != NULL) _hedgehog_mapper -> Delete();
    if (_hedgehog_actor  != NULL) _hedgehog_actor  -> Delete();
    if (_hedgehog_lt     != NULL) _hedgehog_lt     -> Delete();
    //printf("v_isosurf.h: Destructor: all deleted\n");
}

void IsoSurf::Init(int N[3], double L[6], GridCallBack Func, int index, bool octrotate) {

    if (_initialized) {
        printf("Warning: 'IsoSurf' was already initialised\n");
        return;
    }

    // flags
    ShowIsoSurf = true;
    ShowVectors = false;

    // grid
    _sgrid = new SGrid();
    _sgrid -> Init (N, L, Func, index, octrotate);

    // isosurf
    _isosurf        = vtkMarchingContourFilter ::New();
    _isosurf_mapper = vtkPolyDataMapper        ::New();
    _isosurf_actor  = vtkActor                 ::New();
    _isosurf_lt     = vtkLookupTable           ::New();
    _isosurf        -> SetInput                (_sgrid->GetGrid());
    _isosurf        -> ComputeNormalsOff       ();
    _isosurf        -> ComputeGradientsOff     ();
    _isosurf_mapper -> SetInputConnection      (_isosurf->GetOutputPort());
    _isosurf_mapper -> SetLookupTable          (_isosurf_lt);
    _isosurf_actor  -> SetMapper               (_isosurf_mapper);

    // hedgehog
    _hedgehog        = vtkHedgeHog         ::New();
    _hedgehog_mapper = vtkPolyDataMapper   ::New();
    _hedgehog_actor  = vtkActor            ::New();
    _hedgehog_lt     = vtkLookupTable      ::New();
    _hedgehog        -> SetInput           (_sgrid->GetGrid());
    _hedgehog_mapper -> SetInputConnection (_hedgehog->GetOutputPort());
    _hedgehog_mapper -> SetLookupTable     (_hedgehog_lt);
    _hedgehog_actor  -> SetMapper          (_hedgehog_mapper);

    _initialized = true;

    SetColor(0, 1, 1, 1);
    SetValue(0.0);
    SetVecScale(1);
}

void IsoSurf::SetColor(double r, double g, double b, double Opacity) {

    if (!_initialized) return;

    _isosurf_lt->SetNumberOfColors (1);
    _isosurf_lt->Build             ();
    _isosurf_lt->SetTableValue (0, r, g, b);
    _isosurf_actor->GetProperty()->SetOpacity(Opacity);
}

void IsoSurf::AddTo(GoslVTK::Win & win) {

    if (!_initialized) return;

    if (ShowIsoSurf) win.AddActor (_isosurf_actor);
    if (ShowVectors) win.AddActor (_hedgehog_actor);
}

void IsoSurf::SetMaterial(double Ambient, double Diffuse, double Specular, double SpecularPower) {

    if (!_initialized) return;

    _isosurf_actor->GetProperty()->SetAmbient       (Ambient);
    _isosurf_actor->GetProperty()->SetDiffuse       (Diffuse);
    _isosurf_actor->GetProperty()->SetSpecular      (Specular);
    _isosurf_actor->GetProperty()->SetSpecularPower (SpecularPower);
}

void IsoSurf::SetCmap(char const *Name, int NumClrs) {

    if (!_initialized) return;

    double fmin, fmax;
    _sgrid->GetFrange(fmin, fmax);
    SetCmap(Name, NumClrs, fmin, fmax);
}

void IsoSurf::SetCmap(char const *Name, int NumClrs, double Fmin, double Fmax) {

    if (!_initialized) return;

    CmapSetTable(_isosurf_lt, Name, NumClrs, Fmin, Fmax);
    _isosurf_mapper->UseLookupTableScalarRangeOn(); // without this, the mapper changes the lt range
}

void IsoSurf::SetVecCmap(char const *Name, int NumClrs) {

    if (!_initialized) return;

    double fmin, fmax;
    _sgrid->GetFrange(fmin, fmax);
    SetVecCmap(Name, NumClrs, fmin, fmax);
}

void IsoSurf::SetVecCmap(char const *Name, int NumClrs, double Fmin, double Fmax) {

    if (!_initialized) return;

    CmapSetTable(_hedgehog_lt, Name, NumClrs, Fmin, Fmax);
    _hedgehog_mapper->UseLookupTableScalarRangeOn(); // without this, the mapper changes the lt range
}

void IsoSurf::GenValues(int NSurfs) {

    if (!_initialized) return;

    double fmin, fmax;
    _sgrid->GetFrange(fmin, fmax);
    _isosurf->GenerateValues(NSurfs, fmin, fmax);
}

}; // namespace GoslVTK

#endif // GOSLVTK_ISOSURF_H
