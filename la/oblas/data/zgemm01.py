import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [1, 2, 0+1j, 1, -1],
    [2, 3,-1-1j, 1, +1],
    [1, 2, 0+1j, 4, -1],
    [4, 0, 3-1j, 1, +1],
], dtype=complex)

b = np.matrix([
    [1, 0, 0+1j],
    [0, 0, 3-1j],
    [0, 0, 1+1j],
    [1, 0, 1-1j],
    [0, 2, 0+1j],
], dtype=complex)

c = np.matrix([
    [+0.50, 1j, +0.25],
    [+0.25, 1j, -0.25],
    [-0.25, 1j, +0.00],
    [-0.25, 1j, +0.00],
], dtype=complex)

print (0.5-2j)*np.dot(a, b) + (2.0-4j)*c
