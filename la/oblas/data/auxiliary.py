# Copyright 2016 The Gosl Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

def vprint(name, v):
    l = '%s := []float64{' % name
    for i in range(len(v)):
        if i > 0: l += ','
        l += '%23.15e' % v[i]
    l += '}'
    print l

def mprint(name, m, fmt='%23.15e'):
    l = '%s := [][]float64{\n' % name
    for i in range(len(m)):
        l += '    {'
        for j in range(len(m[i])):
            if j > 0: l += ','
            l += fmt % m[i][j]
        l += '},\n'
    l += '}'
    print l
