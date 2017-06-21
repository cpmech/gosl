import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [+3, +0, -3, +0],
    [+0, +3, +1, +2],
    [-3, +1, +4, +1],
    [+0, +2, +1, +3],
],dtype=float)

res = la.cholesky(a, lower=False)
mprint('aUp', res)

res = la.cholesky(a, lower=True)
mprint('aLo', res)
