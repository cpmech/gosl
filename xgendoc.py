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
    ("chk", "Checking numerical calculations and testing"),
    ("io",  "Input and output of files and printing and parsing strings"),
    ("utl", "Utilities such as IntRange and LinSpace and allocators"),
    ("mpi", "A lightweight wrapper to MPI"),
    ("la",  "Linear algebra routines"),
    ("plt", "Plotting routines (wrapping matplotlib)"),
    ("fdm", "A simple finite differences solver"),
    ("num", "A few numerical methods"),
    ("fun", "Functions such as y=f(t,x)"),
    ("ode", "Solvers for ordinary differential equations"),
    ("gm",  "Geometry routines"),
    ("tsr", "Tensor algebra and calculus"),
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

    # fix links
    Cmd("sed -i -e 's@/src/target@https://github.com/cpmech/gosl/blob/master/"+pkg[0]+"@g' "+fn+"")

Cmd('echo "</dl>\n</div><!-- manual-nav -->" >> '+idxfn)
Cmd('echo "'+footer()+'" >> '+idxfn)
