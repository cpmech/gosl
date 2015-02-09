// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_ARROW3D_H
#define GOSLVTK_ARROW3D_H

// VTK
#include <vtkConeSource.h>
#include <vtkCylinderSource.h>
#include <vtkAppendPolyData.h>
#include <vtkTransform.h>
#include <vtkTransformFilter.h>
#include <vtkPolyDataMapper.h>
#include <vtkActor.h>
#include <vtkProperty.h>

// GoslVTK
#include "linalg.h"
#include "util.h"
#include "v_win.h"

/*                +-------------------------------+ 
 *                |            length             |
 *                +-----------------------+-------+
 *                |        bod_len        |tip_len|
 *                |                       |       |
 *                                        `.      ----+  
 *                                        | ``.       |
 *             +  +-----------------------|    ``.    |
 *     bod_rad |  |           +           |   +   >   | tip_rad   
 *             +  +-----------|-----------|   |_-'    |
 *                |           |           | _-|       |
 *                |           |           ''  |     --+  
 *                |           |               |
 *                +-----------+---------------+-------> y axis
 *                |           |               |    
 *                y0      y_bod_cen      y_tip_cen
 */

namespace GoslVTK
{

class Arrow
{
public:
    // Constructor & Destructor
     Arrow ();
    ~Arrow ();

    // Init method
    void Init (double const X0[3], double const V[3], double ConPct=0.1, double ConRad=0.03, double CylRad=0.015, int Res=20);

    // Set methods
    void SetGeometry   (double ConPct, double ConRad, double CylRad);
    void SetResolution (int Resolution);
    void SetColor      (double r, double g, double b, double Opacity);
    void SetVector     (double const X0[3], double const V[3]);
    void SetPoints     (double const X0[3], double const X1[3]);

    // Methods
    void AddTo (GoslVTK::Win & win) { if (_initialized) win.AddActor(_arrow_actor); }

private:
    // essential data
    bool _initialized;

    // vtk data
    vtkConeSource      * _cone;
    vtkCylinderSource  * _cylin;
    vtkTransformFilter * _transform;
    vtkAppendPolyData  * _arrow;
    vtkPolyDataMapper  * _arrow_mapper;
    vtkActor           * _arrow_actor;

    // additional data
    double _tot_len;
    double _con_pct;

    // private methods
    void _update_length();
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


Arrow::Arrow() : _initialized(false) {
    _cone         = NULL;
    _cylin        = NULL;
    _transform    = NULL;
    _arrow        = NULL;
    _arrow_mapper = NULL;
    _arrow_actor  = NULL;
}

Arrow::~Arrow() {
    if (_cone         != NULL) _cone         -> Delete();
    if (_cylin        != NULL) _cylin        -> Delete();
    if (_transform    != NULL) _transform    -> Delete();
    if (_arrow        != NULL) _arrow        -> Delete();
    if (_arrow_mapper != NULL) _arrow_mapper -> Delete();
    if (_arrow_actor  != NULL) _arrow_actor  -> Delete();
    //printf("v_arrow.h: Destructor: all deleted\n");
}

void Arrow::Init(double const X0[3], double const V[3], double ConPct, double ConRad, double CylRad, int Res) {

    if (_initialized) {
        printf("Warning: 'Arrow' was already initialised\n");
        return;
    }

    _cone         = vtkConeSource      ::New();
    _cylin        = vtkCylinderSource  ::New();
    _arrow        = vtkAppendPolyData  ::New();
    _transform    = vtkTransformFilter ::New();
    _arrow_mapper = vtkPolyDataMapper  ::New();
    _arrow_actor  = vtkActor           ::New();
    _cone         -> SetDirection      (0.0, 1.0, 0.0);
    _arrow        -> AddInput          (_cone->GetOutput());
    _arrow        -> AddInput          (_cylin->GetOutput());
    _transform    -> SetInput          (_arrow->GetOutput());
    _arrow_mapper -> SetInput          (_transform->GetPolyDataOutput());
    _arrow_actor  -> SetMapper         (_arrow_mapper);

    _tot_len = 1.0;

    _initialized = true;

    SetGeometry(ConPct, ConRad, CylRad);
    SetResolution(Res);
    SetColor(0, 1, 1, 1);
    SetVector(X0, V);
}

void Arrow::SetGeometry(double ConPct, double ConRad, double CylRad) {
    _con_pct = ConPct;
    _update_length ();
    _cone  -> SetRadius (ConRad);
    _cylin -> SetRadius (CylRad);
}

void Arrow::SetResolution(int Res) {
    _cone  -> SetResolution (Res);
    _cylin -> SetResolution (Res);
}

void Arrow::SetColor(double r, double g, double b, double Opacity) {
    _arrow_actor->GetProperty()->SetColor   (r, g, b);
    _arrow_actor->GetProperty()->SetOpacity (Opacity);
}

void Arrow::SetVector(double const X0[3], double const V[3]) {

    // update length
    _tot_len = V3_norm(V);
    _update_length();

    // translate
    double cen[3] = {0, 0, 0};   // center of arrow
    V3_comb(cen, 1, X0, 0.5, V); // cen := 1*X0 + 0.5*v
    vtkTransform * affine = vtkTransform::New();
    affine->Translate(cen);

    // rotate
    double vy[3] = {0, 1, 0}; // direction of cylinder source
    double angle = (180.0/GoslVTK::PI)*acos(V3_dot(vy,V)/V3_norm(V)); // angle of rotation
    if (angle>0.0) {
        double axis[3] = {0, 0, 0};// axis of rotation
        V3_cross(axis, vy, V);

        // not parallel
        if (V3_norm(axis)>0.0) {
            affine->RotateWXYZ(angle, axis);

        // parallel and oposite (alpha=180)
        } else {
            affine->RotateWXYZ(angle, 0.0, 0.0, 1.0); // use z-direction for mirroring
        }
    }

    // tranform
    _transform->SetTransform(affine);

    // clean up
    affine->Delete();
}

void Arrow::SetPoints(double const X0[3], double const X1[3]) {
    double V[3] = {0, 0, 0};
    V3_sub(V, X1, X0); // V := X1 - X0
    SetVector(X0, V);
}

void Arrow::_update_length() {

    double con_len = _con_pct*_tot_len;      // cone length/height
    double cyl_len = _tot_len-con_len;       // cylinder length/height
    double con_cen = (_tot_len-con_len)/2.0; // cone center
    double cyl_cen = -con_len/2.0;           // cylinder center

    _cone  -> SetCenter (0.0, con_cen, 0.0);
    _cone  -> SetHeight (con_len);
    _cylin -> SetCenter (0.0, cyl_cen, 0.0);
    _cylin -> SetHeight (cyl_len);
}

}; // namespace GoslVTK

#endif // GOSLVTK_ARROW3D_H
