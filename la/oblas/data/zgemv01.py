import numpy as np

a = np.array([[0.1 + 3j, 0.2, 0.3 - 0.3j],
              [1.0 + 2j, 0.2, 0.3 - 0.4j],
              [2.0 + 1j, 0.2, 0.3 - 0.5j],
              [3.0 + 0.1j, 0.2, 0.3 - 0.6j]], dtype=complex)

alp = 0.5+1j
bet = 2.0+1j

x = np.array([20, 10, 30])
y = np.array([3, 1, 2, 4])

res = alp*np.dot(a,x) + bet*y
print res

y = res
res = alp*np.dot(a.T,y) + bet*x
print res
