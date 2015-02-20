// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GOSLVTK_MYSTRING_H
#define GOSLVTK_MYSTRING_H

// STL
#include <string>
#include <cstdarg> // for va_list, va_start, va_end
#include <cstdio>  // for vsnprintf
#include <stdio.h> // for printf

namespace GoslVTK {

class String : public std::string
{
public:
    // Constructors
    String () {}
    String (std::string const & Other) : std::string(Other) {}
    String (char        const * Other) : std::string(Other) {}

    // Methods
    int          Printf   (String const & Fmt, ...);                                                        ///< Print with format
    int          PrintfV  (String const & Fmt, va_list ArgList) { return _set_msg (Fmt.c_str(), ArgList); } ///< Print with format and ArgList
    char const * CStr     ()                              const { return this->c_str(); }                   ///< Get C-string
    void         TextFmt  (char const * NF);                                                                ///< Convert NF (ex: "%10g") to text Format (ex: "%10s")
    void         Split    (String & Left, String & Right, char const * Separator=" ");                      ///< Split string into left and right parts separated by Separator
    bool         HasWord  (String const & Word) { return (find(Word)!=npos); }                              ///< Check if string has a word Word
    void         GetFNKey (String & FNKey);                                                                 ///< Return string without ending ".something"

    // For compatibility with wxWidgets
    String       & ToStdString()       { return (*this); }
    String const & ToStdString() const { return (*this); }

private:
    int _set_msg (char const * Fmt, va_list ArgList);
};

String GetFilename(char const * FileName, char const * Suffix, char const * Ext) {
    String fn;
    String aux(FileName);
    aux.GetFNKey(fn);
    fn += String(Suffix) + String(Ext);
    return fn;
}


/////////////////////////////////////////////////////////////////////////////////////////// Implementation /////


int String::Printf(String const & Fmt, ...) {
    int len;
    va_list       arg_list;
    va_start     (arg_list, Fmt);
    len=_set_msg (Fmt.c_str(), arg_list);
    va_end       (arg_list);
    return len;
}

void String::TextFmt(char const * NF) {
    // number format for text
    this->clear();
    this->append(NF);
    size_t pos;
    pos=this->find("g"); while (pos!=String::npos) { this->replace(pos,1,"s"); pos=this->find("g",pos+1); }
    pos=this->find("f"); while (pos!=String::npos) { this->replace(pos,1,"s"); pos=this->find("f",pos+1); }
    pos=this->find("e"); while (pos!=String::npos) { this->replace(pos,1,"s"); pos=this->find("e",pos+1); }
}

void String::Split(String & Left, String & Right, char const * Separator) {
    size_t pos = find(Separator);
    Left  = substr (0,pos);
    if (pos==npos) Right = "";
    else           Right = substr (pos+String(Separator).size());
}

void String::GetFNKey(String & FNKey) {
    size_t pos = rfind(".");
    FNKey = substr (0,pos);
}

int String::_set_msg (char const * Fmt, va_list ArgList) {
    const int size = 4048; // TODO: remove this limitation by using a loop and reallocating space
    char      buffer[size];
    int       len = std::vsnprintf(buffer, size, Fmt, ArgList);
    this->clear();
    if (len<0) this->append("String::_set_msg: INTERNAL ERROR: std::vsnprintf FAILED");
    else
    {
        buffer[len]='\0';
        if (len>size) this->append("String::_set_msg: INTERNAL ERROR: std::vsnprintf MESSAGE TRUNCATED: ");
        this->append(buffer);
    }
    return len;
}

}; // namespace GoslVTK

#endif // GOSLVTK_MYSTRING_H
