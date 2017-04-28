#!/usr/bin/python

# Copyright 2016 The Gosl Authors. All rights reserved.
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
    ("chk"     , "Check code and unit test tools"),
    ("io"      , "Input/output, read/write files, and print commands"),
    ("utl"     , "Utilities. Lists. Dictionaries. Simple Numerics"),
    ("plt"     , "Plotting and drawing (png and eps)"),
    ("mpi"     , "Message Passing Interface for parallel computing"),
    ("la"      , "Linear Algebra and efficient sparse solvers"),
    ("la/oblas", "Lower level linear algebra using OpenBLAS"),
    ("fdm"     , "Simple finite differences method"),
    ("num"     , "Fundamental Numerical methods"),
    ("fun"     , "Scalar functions of one scalar and one vector"),
    ("gm"      , "Geometry algorithms and structures"),
    ("gm/msh"  , "Mesh structures and interpolation functions for FEA"),
    ("gm/tri"  , "Mesh generation: triangles and Delaunay triangulation"),
    ("gm/rw"   , "Mesh generation: read/write routines"),
    ("graph"   , "Graph theory structures and algorithms"),
    ("ode"     , "Ordinary differential equations"),
    ("opt"     , "Optimisation problem solvers"),
    ("rnd"     , "Random numbers and probability distributions"),
    ("tsr"     , "Tensor algebra and definitions for Continuum Mechanics"),
    ("vtk"     , "3D Visualisation with the VTK tool kit"),
]

odir  = 'doc/'
idxfn = odir+'index.html'
licen = open('LICENSE', 'r').read()

def header(title):
    return """<html>
<head>
<meta http-equiv=\\"Content-Type\\" content=\\"text/html; charset=utf-8\\">
<meta name=\\"viewport\\" content=\\"width=device-width, initial-scale=1\\">
<meta name=\\"theme-color\\" content=\\"#375EAB\\">
<title>%s</title>
<link type=\\"text/css\\" rel=\\"stylesheet\\" href=\\"static/style.css\\">
<script type=\\"text/javascript\\" src=\\"static/godocs.js\\"></script>
<style type=\\"text/css\\"></style>
</head>
<body>
<div id=\\"page\\" class=\\wide\\">
<div class=\\"container\\">
""" % title

def footer():
    return """
<div id=\\"footer\\">
<br /><br />
<hr>
<pre class=\\"copyright\\">
%s</pre><!-- copyright -->
</div><!-- footer -->

</div><!-- container -->
</div><!-- page -->
</body>
</html>""" % licen

def pkgheader(pkg):
    return header('Gosl &ndash; package '+pkg[0]) + '<h1>Gosl &ndash; <b>%s</b> &ndash; %s</h1>' % (pkg[0], pkg[1])

def pkgitem(pkg):
    fnk = pkg[0].replace("/","-")
    return '<dd><a href=\\"xx%s.html\\"><b>%s</b>: %s</a></dd>' % (fnk, pkg[0], pkg[1])

Cmd('echo "'+header('Gosl &ndash; Documentation')+'" > '+idxfn)
Cmd('echo "<h1>Gosl &ndash; Documentation</h1>" >> '+idxfn)
Cmd('echo "<h2 id=\\"pkg-index\\">Index</h2>\n<div id=\\"manual-nav\\">\n<dl>" >> '+idxfn)

for pkg in pkgs:

    fnk = pkg[0].replace("/","-")
    fn = odir+'xx'+fnk+'.html'

    # skip some files in subpackage
    if pkg[0] == "vtk":
        Cmd('mv vtk/autogencgoflags.go /tmp/')

    Cmd('echo "'+pkgheader(pkg)+'" > '+fn)
    Cmd('godoc -html github.com/cpmech/gosl/'+pkg[0]+' >> '+fn)
    Cmd('echo "'+footer()+'" >> '+fn)
    Cmd('echo "'+pkgitem(pkg)+'" >> '+idxfn)

    # copy some files bakc to subpackage
    if pkg[0] == "vtk":
        Cmd('mv /tmp/autogencgoflags.go vtk/')

    # fix links
    Cmd("sed -i -e 's@/src/target@https://github.com/cpmech/gosl/blob/master/"+pkg[0]+"@g' "+fn)
    Cmd("sed -i -e 's@/src/github.com/cpmech/gosl/@https://github.com/cpmech/gosl/blob/master/@g' "+fn)

    # fix links to subdirectories (harder to automate)
    subdirs = []
    if pkg[0] == "fun":
        subdirs = ["figs"]

    if pkg[0] == "gm":
        subdirs = ["data", "msh", "rw", "tri"]

    if pkg[0] == "gm/msh":
        subdirs = ["data"]

    if pkg[0] == "gm/rw":
        subdirs = ["data"]

    if pkg[0] == "graph":
        subdirs = ["data"]

    if pkg[0] == "io":
        subdirs = ["data"]

    if pkg[0] == "la":
        subdirs = ["oblas"]

    if pkg[0] == "num":
        subdirs = ["data"]

    if pkg[0] == "ode":
        subdirs = ["data"]

    if pkg[0] == "opt":
        subdirs = ["data"]

    if pkg[0] == "rnd":
        subdirs = ["data", "dsfmt", "sfmt"]

    if pkg[0] == "utl":
        subdirs = ["data"]

    for subdir in subdirs:
        Cmd("sed -i -e 's@<a href=\""+subdir+"/\">@<a href=\"https://github.com/cpmech/gosl/tree/master/"+pkg[0]+"/"+subdir+"\">@g' "+fn)


    # remove link to parent directory
    Cmd("sed -i -e 's@<td colspan=\"2\"><a href=\"..\">..</a></td>@<td></td>@g' "+fn)

Cmd('echo "</dl>\n</div><!-- manual-nav -->" >> '+idxfn)
Cmd('echo "'+footer()+'" >> '+idxfn)
