import numpy as np
import scipy.linalg as la
from auxiliary import *

A = np.array([
    [ 0 ,  3  ,  2  ,  1 ],
    [ 1 , 1j  , 1j  , 1j ],
    [ 2 ,  2  , 2j  , 2j ],
    [ 3 ,  3  ,  3  , 3j ],
], dtype=complex)

U, s, Vt = la.svd(A, lapack_driver='gesvd')

mprintC('amat', A)
mprintC('uCorrect',  U)
vprint('sCorrect',  s)
mprintC('vtCorrect', Vt)

m, n = np.shape(A)
ns   = min([m, n])

S = np.zeros((m, n))
for i in range(ns):
    S[i,i] = s[i]
USVt = np.dot(U, np.dot(S, Vt))

print 'A    =\n', A
print 'USVt =\n', USVt
print np.allclose(A, USVt)
