from numpy import *
from numpy.linalg import *

def vprint(name, v):
    l = '%s := []float64{' % name
    for i in range(len(v)):
        if i > 0: l += ','
        l += '%24.17e' % v[i]
    l += '}'
    print l

def mprint(name, m):
    l = '%s := [][]float64{\n' % name
    for i in range(len(m)):
        l += '    {'
        for j in range(len(m[i])):
            if j > 0: l += ','
            l += '%25.17e' % m[i][j]
        l += '},\n'
    l += '}'
    print l

if 1:
    A = array([
        [-5.773502691896260e-01, -5.773502691896260e-01, 1],
        [ 5.773502691896260e-01, -5.773502691896260e-01, 1],
        [-5.773502691896260e-01,  5.773502691896260e-01, 1],
        [ 5.773502691896260e-01,  5.773502691896260e-01, 1],
    ], dtype=float)

if 0:
    A = array([
        [64,   2,   3,  61,  60,   6],
        [ 9,  55,  54,  12,  13,  51],
        [17,  47,  46,  20,  21,  43],
        [40,  26,  27,  37,  36,  30],
        [32,  34,  35,  29,  28,  38],
        [41,  23,  22,  44,  45,  19],
        [49,  15,  14,  52,  53,  11],
        [ 8,  58,  59,   5,   4,  62],
    ], dtype=float)

U, s, Vt = svd(A)

mprint('U',  U)
vprint('s',  s)
mprint('Vt', Vt)

m, n = shape(A)
ns   = min([m, n])

S         = zeros((m, n))
S[:n, :n] = diag(s)
USVt      = dot(U, dot(S, Vt))

print 'A    =\n', A
print 'USVt =\n', USVt
print allclose(A, USVt)
