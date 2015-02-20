// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_PLANE_H
#define GOSLVTK_PLANE_H

// Std Lib
#include <string>

// VTK
#include <vtkPlaneSource.h>
#include <vtkPolyDataMapper.h>
#include <vtkActor.h>
#include <vtkProperty.h>

// GoslVTK
#include "linalg.h"
#include "v_win.h"

namespace GoslVTK
{

class Plane
{
public:
    // Constructor and destructor
     Plane ();
    ~Plane ();

    // Init methods
    void Init (const double Ori[3], const double Pt1[3], const double Pt2[3], const double n[3]); // origin, point1, point2, normal

    // Set methods
    void SetCen       (const double Cen[3]) { if (_initialized) _plane->SetCenter(Cen[0], Cen[1], Cen[2]); }
    void SetNormal    (const double n[3])   { if (_initialized) _plane->SetNormal(  n[0],   n[1],   n[2]); }
    void SetColor     (double r, double g, double b, double Opacity);
    void SetWireColor (double r, double g, double b);
    void SetWireWidth (int Width);

    // Methods
    void AddTo (GoslVTK::Win & win);

private:
    // essential data
    bool _initialized;

    // vtk data
    vtkPlaneSource    * _plane;
    vtkPolyDataMapper * _plane_mapper;
    vtkActor          * _plane_actor;
    vtkPolyDataMapper * _wire_mapper;
    vtkActor          * _wire_actor;
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


Plane::Plane() : _initialized(false) {
    _plane        = NULL;
    _plane_mapper = NULL;
    _plane_actor  = NULL;
    _wire_mapper  = NULL;
    _wire_actor   = NULL;
}

Plane::~Plane(){
    if (_plane        != NULL) _plane        -> Delete();
    if (_plane_mapper != NULL) _plane_mapper -> Delete();
    if (_plane_actor  != NULL) _plane_actor  -> Delete();
    if (_wire_mapper  != NULL) _wire_mapper  -> Delete();
    if (_wire_actor   != NULL) _wire_actor   -> Delete();
    //printf("v_plane.h: Destructor: all deleted\n");
}

void Plane::Init(const double Ori[3], const double Pt1[3], const double Pt2[3], const double n[3]) {

    if (_initialized) {
        printf("Warning: 'Plane' was already initialised\n");
        return;
    }

    // plane
    _plane        = vtkPlaneSource    ::New();
    _plane_mapper = vtkPolyDataMapper ::New();
    _plane_actor  = vtkActor          ::New();
    _plane_mapper -> SetInputConnection (_plane->GetOutputPort());
    _plane_actor  -> SetMapper          (_plane_mapper);

    // set plane
    _plane -> SetOrigin(Ori[0], Ori[1], Ori[2]);
    _plane -> SetPoint1(Pt1[0], Pt1[1], Pt1[2]);
    _plane -> SetPoint2(Pt2[0], Pt2[1], Pt2[2]);
    _plane -> SetNormal(  n[0],   n[1],   n[2]);

    // borders
    _wire_mapper = vtkPolyDataMapper    ::New();
    _wire_actor  = vtkActor             ::New();
    _wire_mapper -> SetInput            (_plane->GetOutput());
    _wire_mapper -> ScalarVisibilityOff ();
    _wire_actor  -> SetMapper           (_wire_mapper);
    _wire_actor  -> GetProperty         ()->SetRepresentationToWireframe();

    // set mapper
    _plane_mapper -> SetResolveCoincidentTopologyPolygonOffsetParameters (0,1);
    _plane_mapper -> SetResolveCoincidentTopologyToPolygonOffset         ();
    _wire_mapper  -> SetResolveCoincidentTopologyPolygonOffsetParameters (1,1);
    _wire_mapper  -> SetResolveCoincidentTopologyToPolygonOffset         ();

    // same color for inside and outside edges
    _wire_mapper -> ScalarVisibilityOff          ();
    _wire_actor  -> GetProperty() -> SetAmbient  (1.0);
    _wire_actor  -> GetProperty() -> SetDiffuse  (0.0);
    _wire_actor  -> GetProperty() -> SetSpecular (0.0);

    _initialized = true;

    SetColor(1, 0, 0, 1);
    SetWireColor(0, 0, 1);
    SetWireWidth(1);
}

void Plane::SetColor(double r, double g, double b, double Opacity) {

    if (!_initialized) return;

    _plane_actor->GetProperty()->SetColor   (r, g, b);
    _plane_actor->GetProperty()->SetOpacity (Opacity);
}

void Plane::SetWireColor(double r, double g, double b) {

    if (!_initialized) return;

    _wire_actor->GetProperty()->SetColor (r, g, b);
}

void Plane::SetWireWidth(int Width) {

    if (!_initialized) return;

    _wire_actor->GetProperty()->SetLineWidth(Width);
}

void Plane::AddTo(GoslVTK::Win & win) {

    if (!_initialized) return;

    win.AddActor(_plane_actor);
    win.AddActor(_wire_actor);
}

}; // namespace GoslVTK

#endif // GOSLVTK_PLANE_H
