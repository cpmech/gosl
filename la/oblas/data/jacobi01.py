import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [2, 0, 0],
    [0, 3, 4],
    [0, 4, 9],
], dtype=float)

w, vl, vr = la.eig(a, left=True, right=True)
print 'w =', w
print 'vl =\n', vl
print 'vr =\n', vr

w, v = la.eigh(a)
print
print 'w =', w
print 'v =\n', v

print
print
print

b = np.matrix([
    [1, 2, 3],
    [2, 3, 2],
    [3, 2, 2],
], dtype=float)

w, vl, vr = la.eig(b, left=True, right=True)
print 'w =', w
print 'vl =\n', vl
print 'vr =\n', vr

w, v = la.eigh(b)
print
print 'w =', w
print 'v =\n', v

print
print
print

c = np.matrix([
    [1, 2, 3, 4, 5],
    [2, 3, 0, 2, 4],
    [3, 0, 2, 1, 3],
    [4, 2, 1, 1, 2],
    [5, 4, 3, 2, 1],
], dtype=float)

for i in range(5):
    for j in range(5):
        if c[i,j] != c[j,i]: raise Exception("non symmetric")

w, vl, vr = la.eig(c, left=True, right=True)
print 'w =', w
print 'vl =\n', vl
print 'vr =\n', vr

w, v = la.eigh(c)
print
print 'w =', w
print 'v =\n', v

vprint('l', w)
mprint('Q', v)
