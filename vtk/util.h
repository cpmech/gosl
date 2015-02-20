// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_UTIL_H
#define GOSLVTK_UTIL_H

// STL
#include <cmath>
#include <cfloat> // for DBL_EPSILON

namespace GoslVTK {

// Constants
const double ZERO   = sqrt(DBL_EPSILON); ///< Machine epsilon (smaller positive)
const double SQ2    = sqrt(2.0);         ///< \f$ \sqrt{2} \f$
const double SQ3    = sqrt(3.0);         ///< \f$ \sqrt{3} \f$
const double SQ6    = sqrt(6.0);         ///< \f$ \sqrt{6} \f$
const double SQ2BY3 = sqrt(2.0/3.0);     ///< \f$ \sqrt{2/3} \f$
const double PI     = 4.0*atan(1.0);     ///< \f$ \pi \f$

bool IsNan(double Val) {
    return (std::isnan(Val) || ((Val==Val)==false)); // NaN is the only value, for which the expression Val==Val is always false
}

// Minimum between a and b
template <typename Val_T> Val_T Min(Val_T const & a, Val_T const & b) {
    return (a<b ? a : b);
}

// Maximum between a and b
template <typename Val_T> Val_T Max(Val_T const & a, Val_T const & b) {
    return (a>b ? a : b);
}

///< Maximum between a and b and c
template <typename Val_T> Val_T Max(Val_T const & a, Val_T const & b, Val_T const & c) {
    if (a>=b && a>=c) return a;
    if (b>=a && b>=c) return b;
    return c;
}

}; // namespace GoslVTK

#endif // GOSLVTK_UTIL_H
