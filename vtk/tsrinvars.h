// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_TSRINVARS_H
#define GOSLVTK_TSRINVARS_H

#include <cmath>

#include "util.h"

namespace GoslVTK {

// Calc principal values given octahedral invariants. (-pi <= th(rad) <= pi)
// Using cambridge invariants
void pqth2L(double l[3], double p, double q, double th_deg) {
    double th = th_deg * PI / 180.0;
    l[0] = -p + 2.0*q*sin(th-2.0*PI/3.0)/3.0;
    l[1] = -p + 2.0*q*sin(th)           /3.0;
    l[2] = -p + 2.0*q*sin(th+2.0*PI/3.0)/3.0;
}

}; // namespace GoslVTK

#endif // GOSLVTK_TSRINVARS_H
