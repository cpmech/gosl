import numpy as np
import scipy.linalg as la

a = np.array([
    [+1, +2, +1, +1, -1, +0],
    [+2, +2, +1, +0, +0, +0],
    [+3, +1, +3, +1, +2, -1],
    [+1, +0, +1, -1, +0, +0],
])

tmp = np.array([
    [+1, -1, +0, -1],
    [-1, -1, +1, +0],
    [-1, +1, +1, +1],
    [-1, +0, +1, -1],
])

c = np.dot(tmp,tmp.transpose())

print '....................... c'
print c
print la.det(c)
print la.eigvals(c)
print la.cholesky(c)

res = 3.0 * np.dot(a,a.transpose()) - 1.0 * c

print
print '....................... 3*a*aT - 1*c'
print res
print la.det(res)
print la.eigvals(res)
print la.cholesky(res)
