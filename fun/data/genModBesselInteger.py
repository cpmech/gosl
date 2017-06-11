import scipy.special as sp
import numpy as np

# Reference:
# [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
#     and Mathematical Tables. U.S. Department of Commerce, NIST

# generate tables 9.* starting at page 416 of [1]
def gentable(X):

    # header
    l = '%5s  %13s  %13s  %13s  %13s  %13s  %13s\n' % ('x','exp(-x)*I0(x)','exp(-x)*I1(x)','x^{-2}*I2(x)','exp(x)*K0(x)','exp(x)*K1(x)','x^2*K2(x)')

    # table
    for i, x in enumerate(X):
        a = np.exp(-x)
        A = 1.0/a
        B = x*x
        b = 1.0/B
        l += '%5.2f  %13.10f  %13.10f  %13.10f  %13.10f  %13.10f  %13.10f\n' % (x, a*sp.i0(x), a*sp.i1(x), b*sp.iv(2,x), A*sp.k0(x), A*sp.k1(x), B*sp.kn(2,x))
        if (i+1) % 5 == 0:
            l += '\n'
    return l

# generate data for comparison
def gendata(X):
    l = '%5s%23s%23s%23s%23s%23s%23s%23s%23s\n' % ('x', 'I0', 'I1', 'I2', 'I3', 'K0', 'K1', 'K2', 'K3')
    for i, x in enumerate(X):
        l += '%5.2f%23.15e%23.15e%23.15e%23.15e%23.15e%23.15e%23.15e%23.15e\n' % (x, sp.i0(x), sp.i1(x), sp.iv(2,x), sp.iv(3,x), sp.k0(x), sp.k1(x), sp.kn(2,x), sp.kn(3,x))
    return l

# generate data for comparison
def gendataNeg(X):
    l = '%6s%23s%23s%23s%23s\n' % ('x', 'I0', 'I1', 'I2', 'I3')
    for i, x in enumerate(X):
        l += '%6.2f%23.15e%23.15e%23.15e%23.15e\n' % (x, sp.i0(x), sp.i1(x), sp.iv(2,x), sp.iv(3,x))
    return l

# write file
def savefile(l, fn):
    f = open(fn,'w')
    f.write(l)
    f.close()
    print 'file <%s> written' % fn

# generate tables 9.* starting at page 416 of [1]
l = gentable(np.linspace(0,5.0,51))
savefile(l, '/tmp/as-9-modbessel-integer-tables9.txt')

# data for comparison---small
X = np.linspace(0,20,41)
l = gendata(X)
savefile(l, '/tmp/as-9-modbessel-integer-sml.cmp')

# data for comparison---big
X = np.linspace(0,20,201)
l = gendata(X)
savefile(l, '/tmp/as-9-modbessel-integer-big.cmp')

# data for comparison---negative
X = np.linspace(-10,0,41)
l = gendataNeg(X)
savefile(l, '/tmp/as-9-modbessel-integer-neg.cmp')
