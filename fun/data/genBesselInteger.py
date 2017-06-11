import scipy.special as sp
import numpy as np

# Reference:
# [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
#     and Mathematical Tables. U.S. Department of Commerce, NIST

# generate tables 9.* starting at page 390 of [1]
def gentable(X):

    # header
    l = '%5s  %18s  %13s  %13s  %13s  %13s  %13s\n' % ('x','J0(x)','J1(x)','J2(x)','Y0(x)','Y1(x)','Y2(x)')

    # table
    for i, x in enumerate(X):
        l += '%5.2f  %18.15f  %13.10f  %13.10f  %13.10f  %13.10f  %13.10f\n' % (x, sp.j0(x), sp.j1(x), sp.jv(2,x), sp.y0(x), sp.y1(x), sp.yn(2,x))
        if (i+1) % 5 == 0:
            l += '\n'
    return l

# generate data for comparison
def gendata(X):
    l = '%5s%23s%23s%23s%23s%23s%23s\n' % ('x', 'J0', 'J1', 'J2', 'Y0', 'Y1', 'Y2')
    for i, x in enumerate(X):
        l += '%5.2f%23.15e%23.15e%23.15e%23.15e%23.15e%23.15e\n' % (x, sp.j0(x), sp.j1(x), sp.jv(2,x), sp.y0(x), sp.y1(x), sp.yn(2,x))
    return l

# write file
def savefile(l, fn):
    f = open(fn,'w')
    f.write(l)
    f.close()
    print 'file <%s> written' % fn

# generate tables 9.* starting at page 390 of [1]
l = gentable(np.linspace(0,5.0,51))
savefile(l, '/tmp/as-9-bessel-integer-tables9.txt')

# data for comparison---small
X = np.linspace(0,15,21)
l = gendata(X)
savefile(l, '/tmp/as-9-bessel-integer-sml.cmp')

# data for comparison---big
X = np.linspace(0,15,151)
l = gendata(X)
savefile(l, '/tmp/as-9-bessel-integer-big.cmp')
