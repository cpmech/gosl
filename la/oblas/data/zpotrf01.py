import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [+4 + 0j, 0 + 1j, -3 + 1j, 0 + 2j],
    [+0 - 1j, 3 + 0j, +1 + 0j, 2 + 0j],
    [-3 - 1j, 1 + 0j, +4 + 0j, 1 - 1j],
    [+0 - 2j, 2 + 0j, +1 + 1j, 4 + 0j],
],dtype=complex)

print a - a.getH()

res = la.cholesky(a, lower=False)
mprintC('aUp', res)

res = la.cholesky(a, lower=True)
mprintC('aLo', res)
