import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.array([
    [+1-1j, +2, +1, +1, -1, +0+0j],
    [+2+0j, +2, +1, +0, +0, +0+1j],
    [+3+1j, +1, +3, +1, +2, -1+0j],
    [+1+0j, +0, +1, -1, +0, +0+1j],
],dtype=complex)

tmp = np.array([
    [+1, -1, +0, -1],
    [-1, -1, +1, +0],
    [-1, +1, +1, +1],
    [-1, +0, +1, -1],
],dtype=complex)

c = np.dot(tmp,tmp.conjugate())
c[0,0] += 1j
c[3,3] -= 1j

print '....................... c'
print c
print la.det(c)
print la.eigvals(c)
print la.cholesky(c)

res = 3.0 * np.dot(a,a.transpose()) + 1.0 * c

print
print '....................... 3*a*aT + 1*c'
print res
print la.det(res)
print la.eigvals(res)
print la.cholesky(res)
