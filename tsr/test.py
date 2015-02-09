from numpy.linalg import eig
from numpy import *

AA = [array([[1,2,0],
[2,-2,0],
[0,0,-2]],dtype=float)
, array([[-100,33,0],
[33,-200,0],
[0,0,150]])
, array([[1,2,4],
[2,-2,3],
[4,3,-2]],dtype=float)
, array([[-100,-10,20],
[-10,-200,15],
[20,15,-300]],dtype=float)
, array([[-100,0,-10],
[0,-200,0],
[-10,0,100]],dtype=float)
, array([[0.13,1.2,0],
[1.2,-20,0],
[0,0,-28]],dtype=float)
, array([[-10,3.3,0],
[3.3,-2,0],
[0,0,1.5]],dtype=float)
, array([[0.1,0.2,0.8],
[0.2,-1.3,0.3],
[0.8,0.3,-0.2]])
, array([[-10,-1,2],
[-1,-20,1],
[2,1,-30]],dtype=float)
, array([[-10,0,-1],
[0,-20,0],
[-1,0,10]],dtype=float)
]

def pp(w):
    print '{%23.15e,%23.15e,%23.15e},' % (w[0], w[1], w[2])

for a in AA:
    w, v = eig(a)
    #pp(w)
    print '{'
    for i in range(3):
        print '    {',
        for j in range(3):
            if j > 0: print ',',
            print '%23.15e' % v[i,j],
        print '},'
    print '},'
