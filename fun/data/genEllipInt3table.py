import mpmath as mpm
import numpy as np

def d2r(deg): return deg * np.pi / 180.0
def r2d(rad): return rad * 180.0 / np.pi

# Reference:
# [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
#     and Mathematical Tables. U.S. Department of Commerce, NIST

# generate table 17.9 page 625 of [1]
def gentable(ns, phis, alps):

    # header
    l = '%4s%8s' % ('n', 'alp|phi')
    for phi in phis: l += '%9.0f' % phi
    l += '\n'

    # table
    for n in ns:
        for i, alp in enumerate(alps):
            k = np.sin(d2r(alp))
            m = k**2.0
            l += '%4.1f%8.0f' % (n, alp)
            for j, phi in enumerate(phis):
                p = d2r(phi)
                if np.abs(k-1)<1e-15 and np.abs(p-np.pi/2.0)<1e-15:
                    l += '%9s' % 'inf'
                else:
                    r = mpm.ellippi(n, p, m)
                    l += '%9.5f' % r
            l += '\n'
        l += '\n'
    return l

# generate data for comparison
def gendata(n, phis, alps):
    l = '%23s%23s%23s%23s\n' % ('n', 'phi', 'k', 'PI') # phi [rad]
    for n in ns:
        for i, alp in enumerate(alps):
            k = np.sin(d2r(alp))
            m = k**2.0
            for j, phi in enumerate(phis):
                p = d2r(phi)
                if np.abs(k-1)<1e-15 and np.abs(p-np.pi/2.0)<1e-15:
                    l += '%23.15e%23.15e%23.15e%23s\n' % (n, p, k, 'inf')
                else:
                    r = mpm.ellippi(n, p, m)
                    l += '%23.15e%23.15e%23.15e%23.15e\n' % (n, p, k, r)
    return l

# write file
def savefile(l, fn):
    f = open(fn,'w')
    f.write(l)
    f.close()
    print 'file <%s> written' % fn

# Table 17.5a in [1]
ns = np.linspace(0, 1, 11)
phis = np.array([0,15,30,45,60,75,90],dtype=float)
alps = phis
l = gentable(ns, phis, alps)
savefile(l, '/tmp/as-17-elliptic-integrals-table17.9.txt')

# data for comparison
l = gendata(ns, phis, alps)
savefile(l, '/tmp/as-17-elliptic-integrals-table17.9-big.cmp')

# data for comparison
ns = np.linspace(0, 1, 5)
phis = np.array([0,30,60,90],dtype=float)
alps = phis
l = gendata(ns, phis, alps)
savefile(l, '/tmp/as-17-elliptic-integrals-table17.9-small.cmp')
