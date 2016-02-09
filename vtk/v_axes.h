// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_AXES_H
#define GOSLVTK_AXES_H

// Std Lib
#include <cmath>

// VTK
#include <vtkPoints.h>
#include <vtkLine.h>
#include <vtkUnstructuredGrid.h>
#include <vtkDataSetMapper.h>
#include <vtkActor.h>
#include <vtkProperty.h>
#include <vtkTextActor3D.h>
#include <vtkTextProperty.h>

// GoslVTK
#include "v_win.h"

namespace GoslVTK
{

class Axes
{
public:
    // Constructor & Destructor
     Axes ();
    ~Axes ();

    // Init method
    void Init (double Scale, bool DrawHydroLine=false, bool Reverse=false, bool Full=false, bool Labels=true);

    // Set Methods
    void SetLabels    (char const * X= "x", char const * Y= "y", char const * Z= "z", double r=0, double g=0, double b=1, int SizePt=22, bool Shadow=true);
    void SetWireWidth (int Width) { if (_initialized) _axes_actor->GetProperty()->SetLineWidth(Width); }

    // Methods
    void AddTo (GoslVTK::Win & win);

private:
    // essential data
    bool _initialized;

    // vtk data
    vtkUnstructuredGrid * _axes;
    vtkDataSetMapper    * _axes_mapper;
    vtkActor            * _axes_actor;
    vtkTextActor3D      * _x_label_actor;
    vtkTextActor3D      * _y_label_actor;
    vtkTextActor3D      * _z_label_actor;
    vtkTextProperty     * _text_prop;
    vtkTextActor3D      * _negx_label_actor;
    vtkTextActor3D      * _negy_label_actor;
    vtkTextActor3D      * _negz_label_actor;
    vtkTextProperty     * _neg_text_prop;
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


Axes::Axes() : _initialized(false) {
    _axes             = NULL;
    _axes_mapper      = NULL;
    _axes_actor       = NULL;
    _x_label_actor    = NULL;
    _y_label_actor    = NULL;
    _z_label_actor    = NULL;
    _text_prop        = NULL;
    _negx_label_actor = NULL;
    _negy_label_actor = NULL;
    _negz_label_actor = NULL;
    _neg_text_prop    = NULL;
}

Axes::~Axes(){
    if (_axes             != NULL) _axes             -> Delete();
    if (_axes_mapper      != NULL) _axes_mapper      -> Delete();
    if (_axes_actor       != NULL) _axes_actor       -> Delete();
    if (_x_label_actor    != NULL) _x_label_actor    -> Delete();
    if (_y_label_actor    != NULL) _y_label_actor    -> Delete();
    if (_z_label_actor    != NULL) _z_label_actor    -> Delete();
    if (_text_prop        != NULL) _text_prop        -> Delete();
    if (_negx_label_actor != NULL) _negx_label_actor -> Delete();
    if (_negy_label_actor != NULL) _negy_label_actor -> Delete();
    if (_negz_label_actor != NULL) _negz_label_actor -> Delete();
    if (_neg_text_prop    != NULL) _neg_text_prop    -> Delete();
    //printf("v_axes.h: Destructor: all deleted\n");
}

void Axes::Init(double Scale, bool DrawHydroLine, bool Reverse, bool Full, bool Labels) {

    if (_initialized) {
        printf("Warning: 'Axes' was already initialised\n");
        return;
    }

    // points
    double cte = ((Reverse && (!Full)) ? -Scale : Scale);
    vtkPoints * points = vtkPoints::New();  points->SetNumberOfPoints(6);
    if (Full) {
        points->InsertPoint(0, -cte, 0.0, 0.0);  points->InsertPoint(1, cte, 0.0, 0.0);
        points->InsertPoint(2,  0.0,-cte, 0.0);  points->InsertPoint(3, 0.0, cte, 0.0);
        points->InsertPoint(4,  0.0, 0.0,-cte);  points->InsertPoint(5, 0.0, 0.0, cte); if (DrawHydroLine) {
        points->InsertPoint(6, -cte,-cte,-cte);  points->InsertPoint(7, cte, cte, cte); }
    } else {
        points->InsertPoint(0, 0,0,0);  points->InsertPoint(1, cte, 0.0, 0.0);
        points->InsertPoint(2, 0,0,0);  points->InsertPoint(3, 0.0, cte, 0.0);
        points->InsertPoint(4, 0,0,0);  points->InsertPoint(5, 0.0, 0.0, cte); if (DrawHydroLine) {
        points->InsertPoint(6, 0,0,0);  points->InsertPoint(7, cte, cte, cte); }
    }

    // lines
    vtkLine * line_X = vtkLine::New(); // X axis
    vtkLine * line_Y = vtkLine::New(); // Y axis
    vtkLine * line_Z = vtkLine::New(); // Z axis
    vtkLine * line_H = vtkLine::New(); // hydro axis
    line_X->GetPointIds()->SetNumberOfIds(2); line_X->GetPointIds()->SetId(0,0);  line_X->GetPointIds()->SetId(1,1);
    line_Y->GetPointIds()->SetNumberOfIds(2); line_Y->GetPointIds()->SetId(0,2);  line_Y->GetPointIds()->SetId(1,3);
    line_Z->GetPointIds()->SetNumberOfIds(2); line_Z->GetPointIds()->SetId(0,4);  line_Z->GetPointIds()->SetId(1,5); if (DrawHydroLine) {
    line_H->GetPointIds()->SetNumberOfIds(2); line_H->GetPointIds()->SetId(0,6);  line_H->GetPointIds()->SetId(1,7); }

    // grid
    _axes = vtkUnstructuredGrid::New();
    _axes->Allocate(3,3);
    _axes->InsertNextCell(line_X->GetCellType(),line_X->GetPointIds());
    _axes->InsertNextCell(line_Y->GetCellType(),line_Y->GetPointIds());
    _axes->InsertNextCell(line_Z->GetCellType(),line_Z->GetPointIds()); if (DrawHydroLine) {
    _axes->InsertNextCell(line_H->GetCellType(),line_H->GetPointIds()); }
    _axes->SetPoints(points);

    // mapper and actor
    _axes_mapper = vtkDataSetMapper ::New();
    _axes_actor  = vtkActor         ::New();
    _axes_mapper -> SetInput        (_axes);
    _axes_actor  -> SetMapper       (_axes_mapper);
    _axes_actor  -> GetProperty     () -> SetColor       (0.0,0.0,0.0); 
    _axes_actor  -> GetProperty     () -> SetDiffuseColor(0.0,0.0,0.0); 
    _axes_actor  -> GetProperty     () -> SetLineWidth   (2);

    // clean up
    points -> Delete();
    line_X -> Delete();
    line_Y -> Delete();
    line_Z -> Delete();
    line_H -> Delete();

    // text
    _text_prop     = vtkTextProperty ::New();
    _x_label_actor = vtkTextActor3D  ::New();
    _y_label_actor = vtkTextActor3D  ::New();
    _z_label_actor = vtkTextActor3D  ::New();
    _x_label_actor->SetTextProperty (_text_prop);
    _y_label_actor->SetTextProperty (_text_prop);
    _z_label_actor->SetTextProperty (_text_prop);
    _x_label_actor->SetPosition     (1.05*cte,0,0);
    _y_label_actor->SetPosition     (0,cte,0);
    _z_label_actor->SetPosition     (0,0,cte);
    _x_label_actor->SetScale        (0.003*Scale);
    _y_label_actor->SetScale        (0.003*Scale);
    _z_label_actor->SetScale        (0.003*Scale);
    if (Reverse) {
        _x_label_actor->SetOrientation (-90,0,-180);
        _y_label_actor->SetOrientation (-90,-90,0);
        _z_label_actor->SetOrientation (-90,-90,45);
    } else {
        _x_label_actor->SetOrientation (90,0,180);
        _y_label_actor->SetOrientation (90,90,0);
        _z_label_actor->SetOrientation (90,90,45);
    }
    if (Full) {
        _neg_text_prop    = vtkTextProperty ::New();
        _negx_label_actor = vtkTextActor3D  ::New();
        _negy_label_actor = vtkTextActor3D  ::New();
        _negz_label_actor = vtkTextActor3D  ::New();
        _negx_label_actor->SetTextProperty (_neg_text_prop);
        _negy_label_actor->SetTextProperty (_neg_text_prop);
        _negz_label_actor->SetTextProperty (_neg_text_prop);
        _negx_label_actor->SetPosition     (-cte,0,0);
        _negy_label_actor->SetPosition     (0,-cte,0);
        _negz_label_actor->SetPosition     (0,0,-cte);
        _negx_label_actor->SetScale        (0.003*Scale);
        _negy_label_actor->SetScale        (0.003*Scale);
        _negz_label_actor->SetScale        (0.003*Scale);
        if (Reverse) {
            _negx_label_actor->SetOrientation (-90,0,-180);
            _negy_label_actor->SetOrientation (-90,-90,0);
            _negz_label_actor->SetOrientation (-90,-90,45);
        } else {
            _negx_label_actor->SetOrientation (90,0,180);
            _negy_label_actor->SetOrientation (90,90,0);
            _negz_label_actor->SetOrientation (90,90,45);
        }
    }

    _initialized = true;

    if (Reverse) {
        if (Full) SetLabels ( "X",  "Y",  "Z");
        else      SetLabels ("-X", "-Y", "-Z");
    } else {
        SetLabels();
    }
}

void Axes::AddTo(GoslVTK::Win & win) {

    if (!_initialized) return;

    win.AddActor(_axes_actor);
    if (   _x_label_actor != NULL) win.AddActor(reinterpret_cast<vtkActor*>(   _x_label_actor));
    if (   _y_label_actor != NULL) win.AddActor(reinterpret_cast<vtkActor*>(   _y_label_actor));
    if (   _z_label_actor != NULL) win.AddActor(reinterpret_cast<vtkActor*>(   _z_label_actor));
    if (_negx_label_actor != NULL) win.AddActor(reinterpret_cast<vtkActor*>(_negx_label_actor));
    if (_negy_label_actor != NULL) win.AddActor(reinterpret_cast<vtkActor*>(_negy_label_actor));
    if (_negz_label_actor != NULL) win.AddActor(reinterpret_cast<vtkActor*>(_negz_label_actor));
}

void Axes::SetLabels(char const * X, char const * Y, char const * Z, double r, double g, double b, int SizePt, bool Shadow) {

    if (!_initialized) return;

    if (_x_label_actor == NULL) return;
    if (_y_label_actor == NULL) return;
    if (_z_label_actor == NULL) return;
    if (_text_prop     == NULL) return;
    _x_label_actor -> SetInput    (X);
    _y_label_actor -> SetInput    (Y);
    _z_label_actor -> SetInput    (Z);
    _text_prop     -> SetFontSize (SizePt);
    _text_prop     -> SetColor    (r, g, b);
    if (Shadow) _text_prop -> ShadowOn ();
    else        _text_prop -> ShadowOff();

    char nX[100];
    char nY[100];
    char nZ[100];
    snprintf(nX, 100, "-%s", X);
    snprintf(nY, 100, "-%s", Y);
    snprintf(nZ, 100, "-%s", Z);

    double nR = 1.0;
    double nG = 0.0;
    double nB = 0.0;
    if (r==0 && g==0 && b==0) {
        nR = 0;
        nG = 0;
        nB = 0;
    }

    if (_negx_label_actor == NULL) return;
    if (_negy_label_actor == NULL) return;
    if (_negz_label_actor == NULL) return;
    if (_neg_text_prop    == NULL) return;
    _negx_label_actor -> SetInput    (nX);
    _negy_label_actor -> SetInput    (nY);
    _negz_label_actor -> SetInput    (nZ);
    _neg_text_prop    -> SetFontSize (SizePt);
    _neg_text_prop    -> SetColor    (1, 0, 0);
    if (Shadow) _neg_text_prop -> ShadowOn ();
    else        _neg_text_prop -> ShadowOff();
}

}; // namespace GoslVTK

#endif // GOSLVTK_AXES_H
