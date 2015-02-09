// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_FATAL_H
#define GOSLVTK_FATAL_H

// Std lib
#include <iostream> // for cout
#include <cstdarg>  // for va_list, va_start, va_end
#include <exception>

// Auxiliary
#include "mystring.h"

namespace GoslVTK {

// color codes for the terminal
#define TERM_RST "[0m"
#define TERM_RED "[31m"

// Catch structure
#define GOSLVTK_CATCH catch (GoslVTK::Fatal * e) { e->Cout(); delete e; } \
                      catch (char const * m)     { printf("%sFatal: %s%s\n",TERM_RED,m,TERM_RST); } \
                      catch (std::exception & e) { printf("%sFatal: %s%s\n",TERM_RED,e.what(),TERM_RST); } \
                      catch (...)                { printf("%sFatal: Some exception (...) occurred%s\n",TERM_RED,TERM_RST); }

class Fatal
{
public:
    // Constructor
    Fatal (String const & Fmt, ...) {
        va_list       arg_list;
        va_start     (arg_list, Fmt);
        _msg.PrintfV (Fmt, arg_list);
        va_end       (arg_list);
    }

    // Methods
    void   Cout () const { printf("%sFatal: %s%s\n", TERM_RED, _msg.CStr(), TERM_RST); }
    String Msg  () const { return _msg; }

private:
    String _msg;
};

}; // namespace GoslVTK

#endif // GOSLVTK_FATAL_H
