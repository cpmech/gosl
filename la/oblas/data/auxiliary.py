# Copyright 2016 The Gosl Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import numpy as np

def vprint(name, v, fmt='%+23.15e'):
    l = '%s := []float64{' % name
    for i in range(len(v)):
        if i > 0: l += ','
        l += fmt % v[i]
    l += '}'
    print l

def mprint(name, m, fmt='%+23.15e'):
    l = '%s := [][]float64{\n' % name
    for i in range(len(m)):
        l += '    {'
        for j in range(len(m[i])):
            if j > 0: l += ','
            l += fmt % m[i][j]
        l += '},\n'
    l += '}'
    print l

def vprintC(name, v, fmt='%+23.15e'):
    ff = fmt + ' ' + fmt + 'i'
    l = '%s := []complex128{' % name
    for i in range(len(v)):
        if i > 0: l += ','
        l += ff % (v[i].real, v[i].imag)
    l += '}'
    print l

def mprintC(name, m, fmt='%+23.15e', ztol=1e-16):
    ff = fmt + ' ' + fmt + 'i'
    l = '%s := [][]complex128{\n' % name
    for i in range(len(m)):
        l += '    {'
        for j in range(len(m[i])):
            if j > 0: l += ','
            re, im = m[i][j].real, m[i][j].imag
            if np.abs(re) < ztol: re = 0.0
            if np.abs(im) < ztol: im = 0.0
            l += ff % (re, im)
        l += '},\n'
    l += '}'
    print l
