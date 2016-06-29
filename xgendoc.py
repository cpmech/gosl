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
    ("chk"    , "check and unit test"),
    ("io"     , "input/output"),
    ("utl"    , "utilities"),
    ("plt"    , "plotting"),
    ("mpi"    , "message passing interface"),
    ("la"     , "linear algebra"),
    ("fdm"    , "finite differences method"),
    ("num"    , "numerical methods"),
    ("fun"    , "scalar functions of one scalar and one vector"),
    ("gm"     , "geometry"),
    ("gm/msh" , "mesh generation"),
    ("gm/tri" , "mesh generation: triangles"),
    ("gm/rw"  , "mesh generation: read/write"),
    ("graph"  , "graph theory"),
    ("ode"    , "ordinary differential equations"),
    ("opt"    , "optimisation"),
    ("rnd"    , "random numbers and probability distributions"),
    ("tsr"    , "tensor algebra and definitions for continuum mechanics"),
    ("vtk"    , "visualisation tool kit"),
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
    fnk = pkg[0].replace("/","-")
    return '<dd><a href=\\"xx%s.html\\"><b>%s</b>: %s</a></dd>' % (fnk, pkg[0], pkg[1])

Cmd('echo "'+header('Gosl &ndash; Documentation')+'" > '+idxfn)
Cmd('echo "<h1>Gosl &ndash; Documentation</h1>" >> '+idxfn)
Cmd('echo "<h2 id=\\"pkg-index\\">Index</h2>\n<div id=\\"manual-nav\\">\n<dl>" >> '+idxfn)

for pkg in pkgs:
    fnk = pkg[0].replace("/","-")
    fn = odir+'xx'+fnk+'.html'
    Cmd('echo "'+pkgheader(pkg)+'" > '+fn)
    Cmd('godoc -html github.com/cpmech/gosl/'+pkg[0]+' >> '+fn)
    Cmd('echo "'+footer()+'" >> '+fn)
    Cmd('echo "'+pkgitem(pkg)+'" >> '+idxfn)

    # fix links
    Cmd("sed -i -e 's@/src/target@https://github.com/cpmech/gosl/blob/master/"+pkg[0]+"@g' "+fn+"")

Cmd('echo "</dl>\n</div><!-- manual-nav -->" >> '+idxfn)
Cmd('echo "'+footer()+'" >> '+idxfn)
