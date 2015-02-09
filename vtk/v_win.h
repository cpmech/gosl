// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_WIN_H
#define GOSLVTK_WIN_H

// VTK
#include <vtkCamera.h>
#include <vtkRenderer.h>
#include <vtkRenderWindow.h>
#include <vtkRenderWindowInteractor.h>
#include <vtkInteractorStyleSwitch.h>
#include <vtkActor.h>
#include <vtkWindowToImageFilter.h>
#include <vtkPNGWriter.h>
#include <vtkRenderLargeImage.h>

// GoslVTK
#include "mystring.h"

namespace GoslVTK
{

class Win
{
public:
    // Constructor & Destructor
     Win ();
    ~Win ();

    // Init method
    void Init (int Width=600, int Height=600, double BgRed=1.0, double BgGreen=1.0, double BgBlue=1.0);

    // Methods
    void AddActor (vtkActor * TheActor, bool RstCam=true);
    void DelActor (vtkActor * TheActor) { if (_initialized) _renderer -> RemoveActor(TheActor); }
    void AddLight (vtkLight * Light)    { if (_initialized) _renderer -> AddLight(Light);       }
    void Render   ()                    { if (_initialized) _ren_win  -> Render();              }
    void Show     ();
    void WritePNG (char const * Filename, bool Large=false, int Magnification=1);
    void Camera   (double xUp, double yUp, double zUp, double xFoc, double yFoc, double zFoc, double xPos, double yPos, double zPos);
    void Parallel (bool ParallelProjection=true) { if (_initialized) _camera->SetParallelProjection(ParallelProjection); }

    // Set methods
    void SetViewDefault (bool RevCam=false);
    void SetViewPIplane (bool RevCam=false);
    void SetBgColor     (double r, double g, double b);

private:
    // essential data
    bool _initialized;

    // vtk data
    vtkCamera                 * _camera;
    vtkRenderer               * _renderer;
    vtkRenderWindow           * _ren_win;
    vtkRenderWindowInteractor * _interactor;
    vtkInteractorStyleSwitch  * _int_switch;
};


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


Win::Win() : _initialized(false) {
    _camera     = NULL;
    _renderer   = NULL;
    _ren_win    = NULL;
    _interactor = NULL;
    _int_switch = NULL;
}

Win::~Win() {
    if (_camera     != NULL) _camera     -> Delete();
    if (_renderer   != NULL) _renderer   -> Delete();
    if (_ren_win    != NULL) _ren_win    -> Delete();
    if (_interactor != NULL) _interactor -> Delete();
    if (_int_switch != NULL) _int_switch -> Delete();
    //printf("v_win.h: Destructor: all deleted\n");
}

void Win::Init(int Width, int Height, double BgRed, double BgGreen, double BgBlue) {

    if (_initialized) {
        printf("Warning: 'Win' was already initialised\n");
        return;
    }

    _camera = vtkCamera::New();

    _renderer = vtkRenderer::New();
    _ren_win  = vtkRenderWindow::New();
    _ren_win  -> AddRenderer   (_renderer);
    _ren_win  -> SetSize       (Width, Height);
    _renderer -> SetBackground (BgRed, BgGreen, BgBlue);

    _interactor = vtkRenderWindowInteractor::New();
    _int_switch = vtkInteractorStyleSwitch::New();
    _interactor -> SetRenderWindow    (_ren_win);
    _interactor -> SetInteractorStyle (_int_switch);
    _int_switch -> SetCurrentStyleToTrackballCamera();

    _initialized = true;

    SetViewDefault();
}

void Win::AddActor(vtkActor * TheActor, bool RstCam) {

    if (!_initialized) return;

    _renderer -> AddActor(TheActor);
    if (RstCam) {
        _renderer -> SetActiveCamera (_camera);
        _renderer -> ResetCamera     ();
    }
}

void Win::Show() {

    if (!_initialized) return;

    Render();
    _interactor->Start ();
}

void Win::WritePNG(char const * Filekey, bool Large, int Magnification) {

    if (!_initialized) return;

    // png writer
    String fname(Filekey);   fname += ".png";
    vtkPNGWriter * writer = vtkPNGWriter::New();
    writer -> SetFileName (fname.CStr());

    // re-render window
    Render();

    // write
    if (Large) {
        vtkRenderLargeImage * large_img = vtkRenderLargeImage::New();
        large_img -> SetInput           (_renderer);
        large_img -> Update             ();
        large_img -> SetMagnification   (Magnification);
        writer    -> SetInputConnection (large_img->GetOutputPort());
        writer    -> Write              ();
        large_img -> Delete             ();
    } else {
        vtkWindowToImageFilter * win_to_img = vtkWindowToImageFilter::New();
        win_to_img -> SetInput (_ren_win);
        win_to_img -> Update   ();
        writer     -> SetInput (win_to_img->GetOutput());
        writer     -> Write    ();
        win_to_img -> Delete   ();
    }

    // clean up
    writer -> Delete();

    // Notification
    printf("File <%s.png%s> written\n", Filekey);
}

void Win::Camera(double xUp, double yUp, double zUp, double xFoc, double yFoc, double zFoc, double xPos, double yPos, double zPos) {

    if (!_initialized) return;

    _camera->SetViewUp     (xUp, yUp, zUp);
    _camera->SetFocalPoint (xFoc,yFoc,zFoc);
    _camera->SetPosition   (xPos,yPos,zPos);
    _renderer->ResetCamera ();
}

void Win::SetViewDefault(bool RevCam) {

    if (!_initialized) return;

    double c = (RevCam ? -1 : 1);
    _camera->SetViewUp     (0,0,c);
    _camera->SetPosition   (2*c,c,c);
    _camera->SetFocalPoint (0,0,0);
}

void Win::SetViewPIplane(bool RevCam) {

    if (!_initialized) return;

    double c = (RevCam ? -1 : 1);
    _camera->SetViewUp     (0,0,c);
    _camera->SetPosition   (c,c,c);
    _camera->SetFocalPoint (0,0,0);
}

void Win::SetBgColor(double r, double g, double b) {

    if (!_initialized) return;

    _renderer -> SetBackground(r, g, b);
}

}; // namespace GoslVTK

#endif // GOSLVTK_WIN_H
