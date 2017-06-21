import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [1, 2, 0, 1, -1],
    [2, 3,-1, 1, +1],
    [1, 2, 0, 4, -1],
    [4, 0, 3, 1, +1],
], dtype=float)

b = np.matrix([
    [1, 0, 0],
    [0, 0, 3],
    [0, 0, 1],
    [1, 0, 1],
    [0, 2, 0],
], dtype=float)

c = np.matrix([
    [+0.50, 0, +0.25],
    [+0.25, 0, -0.25],
    [-0.25, 0, +0.00],
    [-0.25, 0, +0.00],
], dtype=float)

print 0.5*np.dot(a, b) + 2*c
