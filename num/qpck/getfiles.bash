#!/bin/bash

FILES="\
dqage.f    \
dqag.f     \
dqagie.f   \
dqagi.f    \
dqagpe.f   \
dqagp.f    \
dqagse.f   \
dqags.f    \
dqawce.f   \
dqawc.f    \
dqawfe.f   \
dqawf.f    \
dqawoe.f   \
dqawo.f    \
dqawse.f   \
dqaws.f    \
dqc25c.f   \
dqc25f.f   \
dqc25s.f   \
dqcheb.f   \
dqelg.f    \
dqk15.f    \
dqk15i.f   \
dqk15w.f   \
dqk21.f    \
dqk31.f    \
dqk41.f    \
dqk51.f    \
dqk61.f    \
dqmomo.f   \
dqng.f     \
dqpsrt.f   \
dqwgtc.f   \
dqwgtf.f   \
dqwgts.f   \
"

# from SciPy
#for f in $FILES; do
    #curl https://raw.githubusercontent.com/scipy/scipy/master/scipy/integrate/quadpack/$f > $f
#done

# extra files
EXTRA="\
d1mach.f \
i1mach.f \
xerror.f \
"

for f in $EXTRA; do
    curl https://raw.githubusercontent.com/scipy/scipy/master/scipy/special/mach/$f > $f
done

# from NetLib
#for f in $FILES; do
    #curl http://www.netlib.org/quadpack/$f > $f
#done
