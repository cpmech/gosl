#!/usr/bin/python

# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import subprocess

def Cmd(command, verbose=False, debug=False):
    if debug:
        print '=================================================='
        print cmd
        print '=================================================='
    spr = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out = spr.stdout.read()
    err = spr.stderr.read().strip()
    if verbose:
        print out
        print err
    return out, err

pkgs = [
    ("fdm", "A simple finite differences solver"),
    ("fun", "Functions such as y=f(t,x)"),
    ("gm",  "Geometry routines"),
    ("la",  "Linear algebra routines"),
    ("mpi", "A lightweight wrapper to MPI"),
    ("num", "A few numerical methods"),
    ("ode", "Solvers for ordinary differential equations"),
    ("plt", "Plotting routines (wrapping matplotlib)"),
    ("tsr", "Tensor algebra and calculus"),
    ("utl", "Utilities"),
    ("vtk", "3D visualisation with VTK"),
]

odir  = 'doc/'
idxfn = odir+'index.html'
licen = open('LICENSE', 'r').read()

def header(title):
    return """<html>
<head>
<meta http-equiv=\\"Content-Type\\" content=\\"text/html; charset=utf-8\\">
<title>%s</title>
<link type=\\"text/css\\" rel=\\"stylesheet\\" href=\\"static/style.css\\">
<script type=\\"text/javascript\\" src=\\"static/godocs.js\\"></script>
<style type=\\"text/css\\"></style>
</head>
<body>
<div id=\\"page\\">""" % title

def footer():
    return """</div><!-- page -->
<div id=\\"footer\\">
<br /><br />
<hr>
<pre class=\\"copyright\\">
%s</pre><!-- copyright -->
</div><!-- footer -->
</body>
</html>""" % licen

def pkgheader(pkg):
    return header('Gosl &ndash; package '+pkg[0]) + '<h1>Gosl &ndash; <b>%s</b> &ndash; %s</h1>' % (pkg[0], pkg[1])

def pkgitem(pkg):
    return '<dd><a href=\\"xx%s.html\\"><b>%s</b>: %s</a></dd>' % (pkg[0], pkg[0], pkg[1])

Cmd('echo "'+header('Gosl &ndash; Documentation')+'" > '+idxfn)
Cmd('echo "<h1>Gosl &ndash; Documentation</h1>" >> '+idxfn)
Cmd('echo "<h2 id=\\"pkg-index\\">Index</h2>\n<div id=\\"manual-nav\\">\n<dl>" >> '+idxfn)

for pkg in pkgs:
    fn = odir+'xx'+pkg[0]+'.html'
    Cmd('echo "'+pkgheader(pkg)+'" > '+fn)
    Cmd('godoc -html github.com/cpmech/gosl/'+pkg[0]+' >> '+fn)
    Cmd('echo "'+footer()+'" >> '+fn)
    Cmd('echo "'+pkgitem(pkg)+'" >> '+idxfn)

Cmd('echo "</dl>\n</div><!-- manual-nav -->" >> '+idxfn)
Cmd('echo "'+footer()+'" >> '+idxfn)
