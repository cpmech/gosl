// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_LINALG_H
#define GOSLVTK_LINALG_H

#include <cmath>

namespace GoslVTK {

// returns u dot v
double V3_dot(const double u[3], const double v[3]) {
    return u[0]*v[0] + u[1]*v[1] + u[2]*v[2];
}

// returns norm(u)
double V3_norm(const double u[3]) {
    return std::sqrt(u[0]*u[0] + u[1]*u[1] + u[2]*u[2]);
}

// w := u cross v
void V3_cross(double w[3], const double u[3], const double v[3]) {
    w[0] = u[1] * v[2] - u[2] * v[1];
    w[1] = u[2] * v[0] - u[0] * v[2];
    w[2] = u[0] * v[1] - u[1] * v[0];
}

// w := u + v
void V3_add(double w[3], const double u[3], const double v[3]) {
    w[0] = u[0] + v[0];
    w[1] = u[1] + v[1];
    w[2] = u[2] + v[2];
}

// w := u - v
void V3_sub(double w[3], const double u[3], const double v[3]) {
    w[0] = u[0] - v[0];
    w[1] = u[1] - v[1];
    w[2] = u[2] - v[2];
}

// w := a * u + b * v (combine)
void V3_comb(double w[3], double a, const double u[3], double b, const double v[3]) {
    w[0] = a * u[0] + b * v[0];
    w[1] = a * u[1] + b * v[1];
    w[2] = a * u[2] + b * v[2];
}

}; // namespace GoslVTK

#endif // GOSLVTK_LINALG_H
