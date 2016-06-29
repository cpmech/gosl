// Copyright 2015 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_CLRMAP_H
#define GOSLVTK_CLRMAP_H

// STL
#include <cstring> // for strcmp

// VTK
#include <vtkLookupTable.h>
#include <vtkColorTransferFunction.h>

namespace GoslVTK {

// Use a color transfer function to generate the colors in the lookup table
void CmapSetTable(vtkLookupTable *lt, char const *name, int ncolors, double fmin, double fmax) {

    // fix range
    if (fmax < fmin) {
        double tmp = fmax;
        fmax = fmin;
        fmin = tmp;
    }
    double df = fmax - fmin;
    if (abs(df)<1e-10) {
        fmin = 0.0;
        fmax = 1.0;
        df   = 1.0;
    }

    // transfer function
    vtkColorTransferFunction *tf = vtkColorTransferFunction::New();

    // green to tan
    if (strcmp(name,"green-tan")==0) {
        tf->SetColorSpaceToDiverging();
        tf->AddRGBPoint(fmin+0.0*df, 0.085, 0.532, 0.201);
        tf->AddRGBPoint(fmin+0.5*df, 0.865, 0.865, 0.865);
        tf->AddRGBPoint(fmin+1.0*df, 0.677, 0.492, 0.093);

    // rainbow
    } else if (strcmp(name,"rainbow")==0) {
        tf->SetColorSpaceToHSV();
        tf->HSVWrapOff();
        tf->AddHSVPoint(fmin+0.0*df, 0.66667, 1.0, 1.0);
        tf->AddHSVPoint(fmin+1.0*df, 0.0, 1.0, 1.0);

    // fire
    } else if (strcmp(name,"fire")==0) {
        tf->SetColorSpaceToRGB();
        tf->AddRGBPoint(fmin+0.0*df, 0.0, 0.0, 0.0);
        tf->AddRGBPoint(fmin+0.4*df, 0.9, 0.0, 0.0);
        tf->AddRGBPoint(fmin+0.8*df, 0.9, 0.9, 0.0);
        tf->AddRGBPoint(fmin+1.0*df, 1.0, 1.0, 1.0);

    // grayscale
    } else if (strcmp(name,"grayscale")==0) {
        tf->SetColorSpaceToRGB();
        tf->AddRGBPoint(fmin+0.0*df, 0.0, 0.0, 0.0);
        tf->AddRGBPoint(fmin+1.0*df, 1.0, 1.0, 1.0);

    // warm
    } else {
        tf->SetColorSpaceToDiverging();
        tf->AddRGBPoint(fmin+0.0*df, 0.230, 0.299, 0.754);
        tf->AddRGBPoint(fmin+1.0*df, 0.706, 0.016, 0.150);
    }

    // lookup table
    if (ncolors < 2) {
        ncolors = 2;
    }
    lt->SetNumberOfTableValues(ncolors);
    lt->SetTableRange(fmin, fmax);
    lt->Build();
    double rgb[3];
    double dx = (fmax - fmin) / static_cast<double>(ncolors - 1);
    for (int i=0; i<ncolors; ++i) {
        double x = fmin + static_cast<double>(i) * dx;
        tf->GetColor(x, rgb);
        lt->SetTableValue(i, rgb[0], rgb[1], rgb[2]);
        //printf("x=%g, r=%g, g=%g, b=%g\n", x, rgb[0], rgb[1], rgb[2]);
    }

    /*
    The TableRange is the input range. No matter how many values are in
    the table, you can map them to any input range.

    The documentation for SetNumberOfColors says "Use
    SetNumberOfTableValues instead", so you should just do that :)

    The ValueRange I guess what you'd want to use for greyscale colors,
    but generally it is the V in the HSV color scheme.

    Say you have this:
        TableRange = [9, 10]
        ValueRange = [125, 255]
        NumberOfColors = 2

    This would take the input value 9.0 and map it to the color with index
    0 in the table, which is (125, something, something) . It will take
    the number 10.0 and map it to the color with index 1 in the table,
    which is (255, something, something) (the somethings depend on your
    choice of SetHueRange and SetSaturationRange.
    */

    // clean up
    tf->Delete();
}

}; // namespace GoslVTK

#endif // GOSLVTK_CLRMAP_H
