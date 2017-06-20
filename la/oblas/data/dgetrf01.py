import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.array([[1, 2, 0, 1],
              [2, 3,-1, 1],
              [1, 2, 0, 4],
              [4, 0, 3, 1]], dtype=float)

lu, piv = la.lu_factor(a)

print lu
print piv

mprint('lu', lu)
mprint('ai', la.inv(a))
