import scipy.special as sp
import numpy as np

def d2r(deg): return deg * np.pi / 180.0
def r2d(rad): return rad * 180.0 / np.pi

# Reference:
# [1] Abramowitz M, Stegun IA (1972) Handbook of Mathematical Functions with Formulas, Graphs,
#     and Mathematical Tables. U.S. Department of Commerce, NIST

# generate table 17.6 page 616 of [1]
def gentable(phis, alps):

    # header
    l = '%7s' % 'alp|phi'
    for phi in phis: l += '%12.0f' % phi
    l += '\n'

    # table
    for i, alp in enumerate(alps):
        k = np.sin(d2r(alp))
        m = k**2.0
        l += '%7.0f' % alp
        for j, phi in enumerate(phis):
            r = sp.ellipeinc(d2r(phi), m)
            l += '%12.8f' % r
        l += '\n'
        if (i+1) % 5 == 0:
            l += '\n'
    return l

# generate data for comparison
def gendata(phis, alps):
    l = '%23s%23s%23s\n' % ('phi', 'k', 'E') # phi [rad]
    for i, alp in enumerate(alps):
        k = np.sin(d2r(alp))
        m = k**2.0
        for j, phi in enumerate(phis):
            p = d2r(phi)
            r = sp.ellipeinc(p, m)
            l += '%23.15e%23.15e%23.15e\n' % (p, k, r)
    return l

# write file
def savefile(l, fn):
    f = open(fn,'w')
    f.write(l)
    f.close()
    print 'file <%s> written' % fn

# Table 17.5a in [1]
phis = np.linspace(0,30,7)
alps = np.linspace(0,90,46)
l = gentable(phis, alps)
l += '\n'
l += gentable(phis, np.linspace(5,85,9))
savefile(l, '/tmp/as-17-elliptic-integrals-table17.6a.txt')

# Table 17.5b in [1]
phis = np.linspace(35,60,6)
alps = np.linspace(0,90,46)
l = gentable(phis, alps)
l += '\n'
l += gentable(phis, np.linspace(5,85,9))
savefile(l, '/tmp/as-17-elliptic-integrals-table17.6b.txt')

# Table 17.5c in [1]
phis = np.linspace(65,90,6)
alps = np.linspace(0,90,46)
l = gentable(phis, alps)
l += '\n'
l += gentable(phis, np.linspace(5,85,9))
savefile(l, '/tmp/as-17-elliptic-integrals-table17.6c.txt')

# data for comparison
phis = [0, 30, 45, 60, 90]
alps = [5, 20, 50, 60, 90]
l = gendata(phis, alps)
savefile(l, '/tmp/as-17-elliptic-integrals-table17.6-small.cmp')

# data for comparison
phis = np.linspace(0,90,19)
alps = np.sort(np.hstack((np.linspace(0,90,46), np.linspace(5,85,9))))
l = gendata(phis, alps)
savefile(l, '/tmp/as-17-elliptic-integrals-table17.6-big.cmp')
