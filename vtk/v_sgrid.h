// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_SGRID_H
#define GOSLVTK_SGRID_H

// VTK
#include <vtkPoints.h>
#include <vtkDoubleArray.h>
#include <vtkStructuredGrid.h>
#include <vtkDataSetMapper.h>
#include <vtkActor.h>
#include <vtkTextActor3D.h>
#include <vtkProperty.h>
#include <vtkTextProperty.h>
#include <vtkPointData.h>
#include <vtkStructuredGridWriter.h>
#include <vtkColorTransferFunction.h>

// GoslVTK
#include "mystring.h"
#include "linalg.h"
#include "v_win.h"

namespace GoslVTK
{

// f and v=dfdx are output
typedef void (*GridCallBack)(double & f, double v[3], const double x[3], int index);

class SGrid
{
public:
    // Constructor & Destructor
     SGrid ();
    ~SGrid ();

    // Init method
    void Init (int N[3], double L[6], GridCallBack Func, int index, bool octrotate);

    // Set methods
    void Resize   (int N[3], double L[6], bool octrotate);
    void SetColor (double r, double g, double b, double Opacity);

    // Additional methods
    double GetF        (int i, int j, int k) const     { if (_initialized) return _scalars->GetTuple1(i+j*_Nx+k*_Nx*_Ny); return 0; }
    void   SetF        (int i, int j, int k, double F) { if (_initialized) _scalars->SetTuple1(i+j*_Nx+k*_Nx*_Ny, F); }
    void   SetCMap     (double Fmin, double Fmax, char const * Name="Diverging");
    void   SetCMap     (                          char const * Name="Diverging") { SetCMap(_Fmin, _Fmax, Name); }
    void   RescaleCMap ();
    void   GetFrange   (double &Fmin, double &Fmax) const { Fmin = _Fmin; Fmax = _Fmax; }

    // Access methods
    int  Size     ()                   const { if (_initialized) return _points->GetNumberOfPoints(); return 0; }
    void GetPoint (int i, double x[3]) const { if (_initialized) _points->GetPoint(i, x); }
    void SetPoint (int i, const double x[3]) { if (_initialized) _points->SetPoint(i, x); }
    
    // Additional access methods
    vtkStructuredGrid * GetGrid() const { if (_initialized) return _sgrid; }

    // Methods
    void ShowWire    ()             { if (_initialized) _sgrid_actor->GetProperty()->SetRepresentationToWireframe(); }
    void ShowSurface ()             { if (_initialized) _sgrid_actor->GetProperty()->SetRepresentationToSurface(); }
    void ShowPoints  (int PtSize=4) { if (_initialized) _sgrid_actor->GetProperty()->SetRepresentationToPoints();  _sgrid_actor->GetProperty()->SetPointSize(PtSize); }
    void AddTo       (GoslVTK::Win & win) { if (_initialized) win.AddActor(_sgrid_actor); }
    void WriteVTK    (char const * Filekey);
    void FilterV     (double F=0.0, double Tol=1.0e-3, bool Normalize=false);

private:
    // essential data
    bool _initialized;

    // additional data
    GridCallBack _func;
    int          _index;

    // vtk data
    vtkPoints                * _points;
    vtkDoubleArray           * _scalars;
    vtkDoubleArray           * _vectors;
    vtkStructuredGrid        * _sgrid;
    vtkDataSetMapper         * _sgrid_mapper;
    vtkActor                 * _sgrid_actor;
    vtkColorTransferFunction * _color_func;

    // additional data
    double _Fmin;
    double _Fmax;
    int    _Nx, _Ny, _Nz;
    String _cmap_name;

    // private methods
    void _calc_f();
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


SGrid::SGrid() : _initialized(false) {
    _points       = NULL;
    _scalars      = NULL;
    _vectors      = NULL;
    _sgrid        = NULL;
    _sgrid_mapper = NULL;
    _sgrid_actor  = NULL;
    _color_func   = NULL;
}

SGrid::~SGrid() {
    if (_points       !=NULL) _points       -> Delete();
    if (_scalars      !=NULL) _scalars      -> Delete();
    if (_vectors      !=NULL) _vectors      -> Delete();
    if (_sgrid        !=NULL) _sgrid        -> Delete();
    if (_sgrid_mapper !=NULL) _sgrid_mapper -> Delete();
    if (_sgrid_actor  !=NULL) _sgrid_actor  -> Delete();
    if (_color_func   !=NULL) _color_func   -> Delete();
    //printf("v_sgrid.h: Destructor: all deleted\n");
}

void SGrid::Init(int N[3], double L[6], GridCallBack Func, int index, bool octrotate) {

    if (_initialized) {
        printf("Warning: 'SGrid' was already initialised\n");
        return;
    }

    // essential data
    _func  = Func;
    _index = index;

    // vtk data
    _points       = vtkPoints                ::New();
    _scalars      = vtkDoubleArray           ::New();
    _vectors      = vtkDoubleArray           ::New();
    _sgrid        = vtkStructuredGrid        ::New();
    _sgrid_mapper = vtkDataSetMapper         ::New();
    _sgrid_actor  = vtkActor                 ::New();
    _color_func   = vtkColorTransferFunction ::New();
    _sgrid        -> SetPoints        (_points);
    _sgrid_mapper -> SetInput         (_sgrid);
    _sgrid_mapper -> SetLookupTable   (_color_func);
    _sgrid_actor  -> SetMapper        (_sgrid_mapper);
    _sgrid_actor  -> GetProperty() -> SetPointSize (4);

    _cmap_name = "rainbow";
    _Fmin      = 0.0;
    _Fmax      = 1.0;

    _initialized = true;

    ShowWire();
    SetColor(0, 0, 0, 1);
    SetCMap(_Fmin, _Fmax);
    Resize(N, L, octrotate); // also calculates F and set _Fmin and _Fmax
}

void SGrid::Resize(int N[3], double L[6], bool octrotate) {

    if (!_initialized) return;

    if (N[0]<2) throw new Fatal("SGrid::Resize: Nx==N[0]=%d must be greater than 1", N[0]);
    if (N[1]<2) throw new Fatal("SGrid::Resize: Ny==N[1]=%d must be greater than 1", N[1]);
    if (N[2]<1) throw new Fatal("SGrid::Resize: Nz==N[2]=%d must be greater than 1", N[2]);
    _Nx = N[0];
    _Ny = N[1];
    _Nz = N[2];
    _points   -> Reset                 ();
    _points   -> Allocate              (_Nx*_Ny*_Nz);
    _scalars  -> Reset                 ();
    _scalars  -> Allocate              (_Nx*_Ny*_Nz);
    _vectors  -> Reset                 ();
    _vectors  -> SetNumberOfComponents (3);
    _vectors  -> SetNumberOfTuples     (_Nx*_Ny*_Nz);
    _sgrid    -> SetDimensions         (_Nx,_Ny,_Nz);
    double dx  = (L[1]-L[0])/(_Nx-1.0);
    double dy  = (L[3]-L[2])/(_Ny-1.0);
    double dz  = (_Nz>1 ? (L[5]-L[4])/(_Nz-1.0) : 0.0);
    double f   = 0.0;
    double x[3] = {0, 0, 0};
    double v[3] = {1, 1, 1};
    double l[3] = {0, 0, 0};
    _Fmin = 0.0;
    _Fmax = 0.0;
    for (int k=0; k<_Nz; ++k)
    for (int j=0; j<_Ny; ++j)
    for (int i=0; i<_Nx; ++i)
    {
        int idx = i + j*_Nx + k*_Nx*_Ny;
        x[0] = L[0] + (double)i * dx;
        x[1] = L[2] + (double)j * dy;
        x[2] = L[4] + (double)k * dz;
        if (octrotate) { // x = {p, q, th}
            GoslVTK::pqth2L(l, x[0], x[1], x[2]);
            x[0] = l[0];
            x[1] = l[1];
            x[2] = l[2];
        }
        _points -> InsertPoint (idx, x);
        if (_func==NULL) {
            _scalars -> InsertTuple1 (idx, 0);
            _vectors -> InsertTuple3 (idx, 0,0,0);
        } else {
            (*_func)(f, v, x, _index);
            _scalars -> InsertTuple1 (idx, f);
            _vectors -> InsertTuple3 (idx, v[0], v[1], v[2]);
            if (f<_Fmin) _Fmin = f;
            if (f>_Fmax) _Fmax = f;
        }
    }
    _sgrid -> GetPointData() -> SetScalars (_scalars);
    _sgrid -> GetPointData() -> SetVectors (_vectors);
}

void SGrid::SetColor(double r, double g, double b, double Opacity) {

    if (!_initialized) return;

    _sgrid_actor->GetProperty()->SetColor   (r, g, b);
    _sgrid_actor->GetProperty()->SetOpacity (Opacity);
}

void SGrid::SetCMap(double Fmin, double Fmax, char const * Name) {

    if (!_initialized) return;

    _cmap_name = Name;
    double midpoint  = 0.5; // halfway between the control points
    double sharpness = 0.0; // linear
    if (_color_func->GetSize() > 0) _color_func->RemoveAllPoints(); // existent
    if (_cmap_name=="rainbow")
    {
        _color_func -> SetColorSpaceToHSV ();
        _color_func -> HSVWrapOff         ();
        _color_func -> AddHSVPoint        (Fmin, 2.0/3.0, 1.0, 1.0, midpoint, sharpness);
        _color_func -> AddHSVPoint        (Fmax, 0.0,     1.0, 1.0, midpoint, sharpness);
    }
    else
    {
        _color_func -> SetColorSpaceToDiverging ();
        _color_func -> HSVWrapOn                ();
        _color_func -> AddRGBPoint              (Fmin, 0.230, 0.299, 0.754, midpoint, sharpness);
        _color_func -> AddRGBPoint              (Fmax, 0.706, 0.016, 0.150, midpoint, sharpness);
    }
}

void SGrid::RescaleCMap() {

    if (!_initialized) return;

    _Fmin = _scalars->GetTuple1 (0);
    _Fmax = _Fmin;
    for (int i=0; i<_scalars->GetNumberOfTuples(); ++i)
    {
        double f = _scalars->GetTuple1(i);
        if (f<_Fmin) _Fmin = f;
        if (f>_Fmax) _Fmax = f;
    }
    SetCMap (_Fmin, _Fmax, _cmap_name.CStr());
}

void SGrid::WriteVTK(char const * Filekey) {

    if (!_initialized) return;

    String buf(Filekey);
    buf.append(".vtk");
    vtkStructuredGridWriter * writer = vtkStructuredGridWriter::New();
    writer -> SetInput    (_sgrid);
    writer -> SetFileName (buf.CStr());
    writer -> Write       ();
    writer -> Delete      ();
    printf("File <%s.vtk> written\n", Filekey);
}

void SGrid::FilterV(double F, double Tol, bool Normalize) {

    if (!_initialized) return;

    for (int i=0; i<_scalars->GetNumberOfTuples(); ++i) {
        double f = _scalars->GetTuple1(i);
        if (Normalize) {
            double * v = _vectors->GetTuple3(i);
            double norm = sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]);
            if (norm>0.0) _vectors->SetTuple3 (i, v[0]/norm, v[1]/norm, v[2]/norm);
        }
        if (fabs(f-F)>Tol) _vectors->SetTuple3 (i, 0.0, 0.0, 0.0);
    }
}

void SGrid::_calc_f() {

    if (!_initialized) return;

    double f;
    double x[3] = {0, 0, 0};
    double v[3] = {1, 1, 1};
    if (_func==NULL) {
        _Fmin = 0.0;
        _Fmax = 1.0;
    } else {
        GetPoint(0, x);
        (*_func)(_Fmin, v, x, _index);
        (*_func)(_Fmax, v, x, _index);
    }
    for (int i=0; i<Size(); ++i) {
        GetPoint(i, x);
        if (_func==NULL) {
            _scalars -> SetTuple1 (i, 0);
            _vectors -> SetTuple3 (i, 0,0,0);
        } else {
            (*_func)(f, v, x, _index);
            _scalars -> SetTuple1 (i, f);
            _vectors -> SetTuple3 (i, v[0], v[1], v[2]);
            if (f<_Fmin) _Fmin = f;
            if (f>_Fmax) _Fmax = f;
        }
    }
}

}; // namespace GoslVTK

#endif // GOSLVTK_SGRID_H
