import numpy as np
import numpy.linalg as la
from auxiliary import *

A = np.array([
    [1, 0, 0, 0, 2],
    [0, 0, 3, 0, 0],
    [0, 0, 0, 0, 0],
    [0, 2, 0, 0, 0],
], dtype=float)

U, s, Vt = la.svd(A)

mprint('amat', A)
mprint('uCorrect',  U)
vprint('sCorrect',  s)
mprint('vtCorrect', Vt)

m, n = np.shape(A)
ns   = min([m, n])

S = np.zeros((m, n))
for i in range(ns):
    S[i,i] = s[i]
USVt = np.dot(U, np.dot(S, Vt))

print 'A    =\n', A
print 'USVt =\n', USVt
print np.allclose(A, USVt)
