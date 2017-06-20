import numpy as np
import numpy.linalg as la
from auxiliary import *

A = np.array([
    [64,   2,   3,  61,  60,   6],
    [ 9,  55,  54,  12,  13,  51],
    [17,  47,  46,  20,  21,  43],
    [40,  26,  27,  37,  36,  30],
    [32,  34,  35,  29,  28,  38],
    [41,  23,  22,  44,  45,  19],
    [49,  15,  14,  52,  53,  11],
    [ 8,  58,  59,   5,   4,  62],
], dtype=float)

U, s, Vt = la.svd(A)

mprint('amat', A, '%3g')
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
