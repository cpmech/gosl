import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [+1-1j, +2, +1, +1, -1, +0+0j],
    [+2+0j, +2, +1, +0, +0, +0+1j],
    [+3+1j, +1, +3, +1, +2, -1+0j],
    [+1+0j, +0, +1, -1, +0, +0+1j],
],dtype=complex)

tmp = np.matrix([
    [+1-1j, -1, +0, -1+0j],
    [-1+0j, -1, +1, +0+0j],
    [-1+0j, +1, +1, +1+0j],
    [-1+0j, +0, +1, -1+1j],
],dtype=complex)

c = np.dot(tmp,tmp.getH())


print '....................... c'
print np.max(np.abs(c.getH() - c))
print c
print la.det(c)
print la.eigvals(c)
print la.cholesky(c)

res = 3.0 * np.dot(a,a.getH()) + 1.0 * c

print
print '....................... 3*a*aT + 1*c'
print res
print la.det(res)
print la.eigvals(res)
print la.cholesky(res)
