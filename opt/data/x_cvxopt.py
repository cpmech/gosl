# change this in from line # 922 of coneprog.py in /usr/lib/python2.7/dist-packages/cvxopt
#            if iters == 0:
#                print("%2s%16s%16s%10s%10s%10s%16s" %("it","pcost", "dcost", "gap", "pres", "dres", "k/t"))
#            print("%2d%16.8e%16.8e%10.3e%10.3e%10.3e%16.8e" %(iters, pcost, dcost, gap, pres, dres, kappa/tau))
from cvxopt.modeling import op
lp = op()
#lp.fromfile('afiro.mps')
#lp.fromfile('adlittle.mps')
lp.fromfile('share1b.mps')
lp.solve()
res = {}
for i, v in enumerate(lp.variables()):
    #print v.value[0]
    res[v.name] = v.value[0]
