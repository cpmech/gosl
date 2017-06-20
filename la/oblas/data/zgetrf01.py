import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.array([[1+1j, 2, 0, 1-1j],
              [2+1j, 3,-1, 1-1j],
              [1+1j, 2, 0, 4-1j],
              [4+1j, 0, 3, 1-1j]], dtype=complex)

lu, piv = la.lu_factor(a)

print lu
print piv

mprintC('lu', lu)
mprintC('ai', la.inv(a))
