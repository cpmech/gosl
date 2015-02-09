// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_SPHERE_H
#define GOSLVTK_SPHERE_H

// VTK
#include <vtkSphereSource.h>
#include <vtkPolyDataMapper.h>
#include <vtkActor.h>
#include <vtkProperty.h>

// GoslVTK
#include "linalg.h"
#include "v_win.h"

namespace GoslVTK
{

class Sphere
{
public:
    // Constructor & Destructor
     Sphere ();
    ~Sphere ();

    // Init method
    void Init (double X[3], double R, int ThetaRes=20, int PhiRes=20);

    // Set methods
    void SetCenter     (double X[3]);
    void SetRadius     (double R);
    void SetResolution (int ThetaRes=20, int PhiRes=20);
    void SetColor      (double r, double g, double b, double Opacity);

    // Methods
    void AddTo (GoslVTK::Win & win, bool RstCam=true) { if (_initialized) win.AddActor(_sphere_actor, RstCam); }

private:
    // essential data
    bool _initialized;

    // vtk data
    vtkSphereSource   * _sphere;
    vtkPolyDataMapper * _sphere_mapper;
    vtkActor          * _sphere_actor;
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


Sphere::Sphere() : _initialized(false) {
    _sphere        = NULL;
    _sphere_mapper = NULL;
    _sphere_actor  = NULL;
}

Sphere::~Sphere() {
    if (_sphere        != NULL) _sphere        -> Delete();
    if (_sphere_mapper != NULL) _sphere_mapper -> Delete();
    if (_sphere_actor  != NULL) _sphere_actor  -> Delete();
    //printf("v_sphere.h: Destructor: all deleted\n");
}

void Sphere::Init(double X[3], double R, int ThetaRes, int PhiRes) {

    if (_initialized) {
        printf("Warning: 'Sphere' was already initialised\n");
        return;
    }

    _sphere        = vtkSphereSource     ::New();
    _sphere_mapper = vtkPolyDataMapper   ::New();
    _sphere_actor  = vtkActor            ::New();
    _sphere_mapper -> SetInputConnection (_sphere->GetOutputPort());
    _sphere_actor  -> SetMapper          (_sphere_mapper);

    _initialized = true;

    SetCenter(X);
    SetRadius(R);
    SetResolution(ThetaRes, PhiRes);
    SetColor(1, 1, 0, 1);
}

void Sphere::SetCenter(double X[3]) {

    if (!_initialized) return;

    _sphere -> SetCenter(X[0], X[1], X[2]);
}

void Sphere::SetRadius(double R) {

    if (!_initialized) return;

    _sphere -> SetRadius(R);
}

void Sphere::SetResolution(int ThetaRes, int PhiRes) {

    if (!_initialized) return;

    _sphere -> SetThetaResolution (ThetaRes);
    _sphere -> SetPhiResolution   (PhiRes);
}

void Sphere::SetColor(double r, double g, double b, double Opacity) {

    if (!_initialized) return;

    _sphere_actor->GetProperty()->SetColor   (r, g, b);
    _sphere_actor->GetProperty()->SetOpacity (Opacity);
}

}; // namespace GoslVTK

#endif // GOSLVTK_SPHERE_H
