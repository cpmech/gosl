# change this in from line # 922 of coneprog.py in /usr/lib/python2.7/dist-packages/cvxopt
#            if iters == 0:
#                print("%2s%16s%16s%10s%10s%10s%16s" %("it","pcost", "dcost", "gap", "pres", "dres", "k/t"))
#            print("%2d%16.8e%16.8e%10.3e%10.3e%10.3e%16.8e" %(iters, pcost, dcost, gap, pres, dres, kappa/tau))
from cvxopt import matrix, solvers, printing
A = matrix([ [-1.0, -1.0, 0.0, 1.0], [1.0, -1.0, -1.0, -2.0] ])
b = matrix([ 1.0, -2.0, 0.0, 4.0 ])
c = matrix([ 2.0, 1.0 ])
sol = solvers.lp(c,A,b)
printing.options['dformat']='%.15f'
print sol['x']
