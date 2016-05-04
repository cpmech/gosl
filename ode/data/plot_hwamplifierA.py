# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from pylab import subplot, plot, show
from gosl  import Read, Gll

go = Read('/tmp/gosl/hwamplifierA.res')
fo = Read('data/radau5_hwamplifier.dat') # HW Fortran code results

for i in range(8):
    subplot(8,1,i+1)
    fo_l = 'hw'
    go_l = 'go'
    plot(fo['x'], fo['y%d'%i], '.', label="HW Code (A)")
    plot(go['x'], go['y%d'%i], '+', label="GoSL", ms=8)
    Gll('x','y%d'%i)

show()
