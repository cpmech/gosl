import numpy as np
import scipy.linalg as la
from auxiliary import *

a = np.matrix([
    [+0.35, +0.45, -0.14, -0.17],
    [+0.09, +0.07, -0.54, +0.35],
    [-0.44, -0.33, -0.03, +0.17],
    [+0.25, -0.32, -0.13, +0.11],
], dtype=float)

w, vl, vr = la.eig(a, left=True, right=True)

vprintC('w', w)

print
for i in range(4):
    vprintC('vl%d'%i, vl[:,i])

print
for i in range(4):
    vprintC('vr%d'%i, vr[:,i])
