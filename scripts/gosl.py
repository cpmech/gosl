# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# for Python 3
from __future__ import print_function

import subprocess
import sys
import os.path
from os import remove
from os.path import basename, exists
#from scipy.special import erfc
from numpy.linalg import norm, eig, solve
from numpy import cosh, sinh, polyfit, NaN
from numpy import pi, sin, cos, tan, arcsin, arccos, arctan, arctan2, log, log10, exp, sqrt
from numpy import array, linspace, insert, repeat, zeros, matrix, ones, eye, arange, diag, dot
from numpy import logical_or, logical_and, delete, hstack, vstack, meshgrid, vectorize, transpose
from pylab import rcParams, gca, gcf, clf, savefig, ScalarFormatter
from pylab import plot, xlabel, ylabel, show, grid, legend, subplot, axis, text, axhline, axvline, title, xticks
from pylab import contour, contourf, colorbar, clabel, xlim, suptitle, loglog, hist, rcdefaults
from pylab import annotate, subplots_adjust, quiver
from pylab import cm as MPLcm
from pylab import close as MPLclose
from matplotlib.transforms      import offset_copy
from matplotlib.patches         import FancyArrowPatch, PathPatch, Polygon
from matplotlib.patches         import Arc    as MPLArc
from matplotlib.patches         import Circle as MPLCircle
from matplotlib.path            import Path   as MPLPath
from matplotlib.lines           import Line2D
from matplotlib.font_manager    import FontProperties
from matplotlib.ticker          import FuncFormatter
from mpl_toolkits.mplot3d       import Axes3D
from mpl_toolkits.mplot3d.art3d import Poly3DCollection
from matplotlib.ticker          import MaxNLocator
#from scipy.interpolate          import UnivariateSpline
from matplotlib.patheffects     import Stroke # for fixing arrow tips
from matplotlib.gridspec        import GridSpec


def SplotGap(w, h):
    subplots_adjust(wspace=w, hspace=h)


def AnnotateXlabels(x, txt, fs=7):
    annotate(txt, xy=(x, -fs-3), xycoords=('data', 'axes points'), va='top', ha='center', size=fs, clip_on=False)


def Remove(filename):
    if os.path.exists(filename):
        os.remove(filename)


def External(command, silent=False):
    spr = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out = spr.stdout.read()
    err = spr.stderr.read().strip()
    if err:
        if not silent:
            print('[1;31m')
            print('out =', out)
            print('err =', err)
            print('Extenal command failed: command == <', command, '>[0m')
        return 1
    return 0


def RunPdfLatex(fnkey):
    if exists('%s.pdf' % fnkey):
        remove('%s.pdf' % fnkey)
    res = subprocess.call('pdflatex -interaction=batchmode -halt-on-error %s.tex' % fnkey, shell=True)
    err = res == 1
    if not exists('%s.pdf' % fnkey):
        err = True
    if err:
        print('[1;31mlatex command failed when processing file <[1;34m%s.tex[0m>[0m' % fnkey)
    else:
        remove('%s.aux' % fnkey)
        remove('%s.log' % fnkey)
        print('file <[1;34m%s.pdf[0m> created' % fnkey)


def SetXlog():
    gca().set_xscale('log')


def SetYlog():
    gca().set_yscale('log')


def SetXhidden():
    gca().get_xaxis().set_visible(False)


def SetYhidden():
    gca().get_yaxis().set_visible(False)


def HideAxes():
    SetXhidden()
    SetYhidden()


def SetXnticks(num):
    if num == 0: gca().get_xaxis().set_ticks([])
    else:        gca().get_xaxis().set_major_locator(MaxNLocator(num))


def SetYnticks(num):
    if num == 0: gca().get_yaxis().set_ticks([])
    else:        gca().get_yaxis().set_major_locator(MaxNLocator(num))


def TgLine(m0, n0, dm, slope, both=True):
    """
    Data for tangent line at (m0,n0)
    ================================
    """
    a, b = m0-dm, m0+dm
    c, d = n0+(a-m0)*slope, n0+(b-m0)*slope
    if both: return [a,  b], [c,  d]
    else:    return [m0, b], [n0, d]


def DrawPath(dat, with_points=True, doshow=True, dolims=True, fc='#edf5ff', ec='black', ptclr='#6b1bfd', lw=2, zorder=None, clip_on=False):
    """
    Draw MPL path
    =============
    """
    if dolims or with_points:
        x = [d[1][0] for d in dat]
        y = [d[1][1] for d in dat]
    if dolims:
        axis('equal')
        cmin = [min(x), min(y)]
        cmax = [max(x), max(y)]
        axis([cmin[0], cmax[0], cmin[1], cmax[1]])
    if with_points:
        zo = None if zorder == None else zorder + 1
        plot(x, y, 'bo', zorder=zo, color=ptclr, clip_on=clip_on)
    commands, vertices = zip(*dat)
    ph = MPLPath(vertices, commands)
    pc = PathPatch(ph, fc=fc, ec=ec, lw=lw, zorder=zorder, clip_on=clip_on)
    gca().add_patch(pc)


def DrawSlope(X, Y, numfmt='%.3f', div=4.0, fsz=8, datline=True, updown=False, label=None, lbHa='left', lbVa='center', *args, **kwargs):
    """
    Find slope and draw icon
    ========================
    """
    m, c   = polyfit(X, Y, 1)
    yf     = lambda x: c + m * x
    x0, x1 = min(X), max(X)
    xx     = array([x0,x1])
    yy     = yf(xx)
    y0, y1 = min(yy), max(yy)
    dx, dy = x1-x0, y1-y0
    xm, ym = (x0+x1)/2.0, (y0+y1)/2.0
    xa, xb = xm-dx/(2.0*div), xm+dx/(2.0*div)
    ya, yb = yf(xa), yf(xb)
    yc     = (ya+yb)/2.0
    if datline: plot(xx, yy, color='gray')
    if not 'color' in kwargs: kwargs['color'] = 'grey'
    if not 'ls'    in kwargs: kwargs['ls'   ] = '-'
    lb = r'$%s$'%(numfmt%m) if label==None else label
    if updown:
        plot([xa,xa,xb], [ya,yb,yb], *args, **kwargs)
        text(xa, yc, lb, fontsize=fsz, color=kwargs['color'], ha=lbHa, va=lbVa)
    else:
        plot([xa,xb,xb], [ya,ya,yb], *args, **kwargs)
        text(xb, yc, lb, fontsize=fsz, color=kwargs['color'], ha=lbHa, va=lbVa)


def ConvIndicator(X, Y, pct=0.1, fs=14, eqaxis=False):
    """
    Convergence indicator icon
    ==========================
    """
    if len(X)<2: raise Exception('at least 2 points are required')
    xx, yy   = log10(X), log10(Y)
    p        = polyfit(xx, yy, 1)
    m        = round(p[0])
    xx0, xx1 = min(xx), max(xx)
    yy0, yy1 = min(yy), max(yy)
    dxx, dyy = xx1-xx0, yy1-yy0
    xxm, yym = (xx0+xx1)/2.0, (yy0+yy1)/2.0
    xxl, xxr = xxm-pct*dxx, xxm+pct*dxx
    shift    = 0.5*pct*dxx*m
    xm,  ym  = 10.0**xxm, 10.0**(yym-shift)
    xl,  xr  = 10.0**xxl, 10.0**xxr
    yl,  yr  = 10.0**(yym+m*(xxl-xxm)-shift),10.0**(yym+m*(xxr-xxm)-shift)
    loglog(X, Y)
    #plot(xm, ym, 'ro')
    #plot(xl, yl, 'go')
    #plot(xr, yr, 'mo')
    points = array([[xl,yl],[xr,yl],[xr,yr]])
    gca().add_patch(Polygon(points, ec='k', fc='None'))
    xxR = xxm+1.2*pct*dxx
    xR  = 10.0**xxR
    text(xR, ym, '%g'%m, ha='left', va='center', fontsize=fs)
    if eqaxis: axis('equal')
    return m

def NlSpace(xmin, xmax, nx, n=2.0, rev=False):
    """
    Non-linear space
    ================
    INPUT:
        xmin : min x-coordinate
        xmax : max x-coordinate
        nx   : number of points
        n    : power of nonlinear term
        rev  : reversed?
    RETURNS:
        a non-linear sequence
    """
    fcn = vectorize(lambda i: xmin + (xmax-xmin)*(float(i)/float(nx-1))**n)
    return fcn(arange(nx))


def SetScientificFmt(axis='y', min_order=-3, max_order=3):
    """
    Scalar format for axes
    ======================
    """
    # applies to current axis
    fmt = ScalarFormatter(useOffset=True)
    fmt.set_powerlimits((min_order,max_order))
    if axis=='y': gca().yaxis.set_major_formatter(fmt)
    else:         gca().xaxis.set_major_formatter(fmt)


def HideFrameLines(spines_to_remove=['top', 'right']):
    """
    Hide lines from frame
    =====================
    """
    for spine in spines_to_remove:
        gca().spines[spine].set_visible(False)


def Plot3dLine(X, Y, Z, first=True, xlbl='X', ylbl='Y', zlbl='Z', zmin=None, zmax=None, splot=111, *args, **kwargs):
    """
    Plot 3D line
    ============
    """
    if first:
        ax = gcf().add_subplot(splot, projection='3d')
        ax.set_xlabel(xlbl)
        ax.set_ylabel(ylbl)
        ax.set_zlabel(zlbl)
    else: ax = gca()
    ax.plot(X,Y,Z, *args, **kwargs)
    if zmin!=None and zmax!=None:
        ax.set_zlim3d(zmin,zmax)
    return ax


def Plot3dPoints(X, Y, Z, xlbl='X', ylbl='Y', zlbl='Z', zmin=None, zmax=None, splot=111,
        preservePrev=False, *args, **kwargs):
    """
    Plot 3D points
    ==============
    """
    if preservePrev: ax = gca()
    else: ax = gcf().add_subplot(splot, projection='3d')
    if xlbl != "": ax.set_xlabel(xlbl)
    if ylbl != "": ax.set_ylabel(ylbl)
    if zlbl != "": ax.set_zlabel(zlbl)
    ax.scatter(X,Y,Z, *args, **kwargs)
    if zmin!=None and zmax!=None:
        ax.set_zlim3d(zmin,zmax)
    return ax


def Wireframe(X, Y, Z, xlbl='X', ylbl='Y', zlbl='Z', zmin=None, zmax=None, cmapidx=0, splot=111,
        rstride=1, cstride=1, preservePrev=False, *args, **kwargs):
    """
    Plot 3D wireframe
    =================
    """
    if preservePrev: ax = gca()
    else: ax = gcf().add_subplot(splot, projection='3d')
    if xlbl != "": ax.set_xlabel(xlbl)
    if ylbl != "": ax.set_ylabel(ylbl)
    if zlbl != "": ax.set_zlabel(zlbl)
    ax.plot_wireframe(X,Y,Z,rstride=rstride,cstride=cstride,cmap=Cmap(cmapidx), *args, **kwargs)
    if zmin!=None and zmax!=None:
        ax.set_zlim3d(zmin,zmax)
    return ax


def Surface(X, Y, Z, xlbl='X', ylbl='Y', zlbl='Z', zmin=None, zmax=None, cmapidx=0, splot=111,
        rstride=1, cstride=1, preservePrev=False, *args, **kwargs):
    """
    Plot 3D surface
    ===============
    """
    if preservePrev: ax = gca()
    else: ax = gcf().add_subplot(splot, projection='3d')
    if xlbl != "": ax.set_xlabel(xlbl)
    if ylbl != "": ax.set_ylabel(ylbl)
    if zlbl != "": ax.set_zlabel(zlbl)
    ax.plot_surface(X,Y,Z,rstride=rstride,cstride=cstride,cmap=Cmap(cmapidx), *args, **kwargs)
    if zmin!=None and zmax!=None:
        ax.set_zlim3d(zmin,zmax)
    return ax


def AddPoly3D(ax, verts, clr='magenta'):
    """
    Plot polygons in 3D
    ===================
      verts = zeros((nnodes, 3))
      ax = gcf().add_subplot(224, projection='3d')
    """
    ax.add_collection3d(Poly3DCollection([zip(verts[:,0], verts[:,1], verts[:,2])], facecolors=[clr]))


def Cmap(idx):
    """
    Get colormap
    ============
    """
    cmaps = [MPLcm.bwr, MPLcm.RdBu, MPLcm.hsv, MPLcm.jet, MPLcm.terrain, MPLcm.pink, MPLcm.Greys]
    return cmaps[idx % len(cmaps)]


def SetFontSize(xylabel_fs=9, leg_fs=8, text_fs=9, xtick_fs=7, ytick_fs=7):
    """
    Set font sizes
    ==============
    """
    params = {
        'axes.labelsize'  : xylabel_fs,
        'legend.fontsize' : leg_fs,
        'font.size'       : text_fs,
        'xtick.labelsize' : xtick_fs,
        'ytick.labelsize' : ytick_fs,
    }
    rcParams.update(params)


def SetForPng(proport=0.75, fig_width_pt=455.24, dpi=150, xylabel_fs=9, leg_fs=8, text_fs=9, xtick_fs=7, ytick_fs=7):
    """
    Set figure proportions
    ======================
    """
    inches_per_pt = 1.0/72.27                   # Convert pt to inch
    fig_width     = fig_width_pt*inches_per_pt  # width in inches
    fig_height    = fig_width*proport           # height in inches
    fig_size      = [fig_width,fig_height]
    params = {
        'axes.labelsize'  : xylabel_fs,
        'font.size'       : text_fs,
        'legend.fontsize' : leg_fs,
        'xtick.labelsize' : xtick_fs,
        'ytick.labelsize' : ytick_fs,
        'figure.figsize'  : fig_size,
        'savefig.dpi'     : dpi,
    }
    #utl.Ff(&bb, "SetFontSize(%s)\n", args)
    #utl.Ff(&bb, "rcParams.update({'savefig.dpi':%d})\n", dpi)
    MPLclose()
    rcdefaults()
    rcParams.update(params)


def SetForEps(proport=0.75, fig_width_pt=455.24, xylabel_fs=9, leg_fs=8, text_fs=9, xtick_fs=7,
        ytick_fs=7, text_usetex=True, mplclose=True):
    """
    Set figure proportions
    ======================
    """
    # fig_width_pt = 455.24411                  # Get this from LaTeX using \showthe\columnwidth
    inches_per_pt = 1.0/72.27                   # Convert pt to inch
    fig_width     = fig_width_pt*inches_per_pt  # width in inches
    fig_height    = fig_width*proport           # height in inches
    fig_size      = [fig_width,fig_height]
    #params = {'mathtext.fontset':'stix', # 'cm', 'stix', 'stixsans', 'custom'
    params = {
        'backend'            : 'ps',
        'axes.labelsize'     : xylabel_fs,
        'font.size'          : text_fs,
        'legend.fontsize'    : leg_fs,
        'xtick.labelsize'    : xtick_fs,
        'ytick.labelsize'    : ytick_fs,
        'text.usetex'        : text_usetex, # very IMPORTANT to avoid Type 3 fonts
        'ps.useafm'          : True, # very IMPORTANT to avoid Type 3 fonts
        'pdf.use14corefonts' : True, # very IMPORTANT to avoid Type 3 fonts
        'figure.figsize'     : fig_size,
    }
    if mplclose: MPLclose()
    rcdefaults()
    rcParams.update(params)


def Save(filename, ea=None, verbose=True):
    """
    Save fig with extra artists
    ===========================
    INPUT:
        ea : extra artists to adjust figure size.
             it can be a list or a matplotlib object
    Note:
        As a workaround, savefig can take bbox_extra_artists keyword (this may
        only be in the svn version though), which is a list artist that needs
        to be accounted for the bounding box calculation. So in your case, the
        below code will work.
        t1 = ax.text(-0.2,0.5,'text',transform=ax.transAxes)
        fig.savefig('test.png', bbox_inches='tight', bbox_extra_artists=[t1])
    """
    if ea==None:
        ea = []
    else:
        if not isinstance(ea, list):
            ea = [ea]
    ea = [x for x in ea if x is not None]
    savefig (filename, bbox_inches='tight', bbox_extra_artists=ea)
    if verbose:
        print('File <[1;34m%s[0m> written'%filename)


def Leg(fsz=8, ncol=1, loc='best', out=False, hlen=3, out_dims=None, frameon=False, *args, **kwargs):
    """
    Legend
    ======
        hlen : handle length
    """
    handles, labels = gca().get_legend_handles_labels()
    if len(handles) < 1 or len(labels) < 1: return None
    if out:
        dims = out_dims if out_dims else (0.,1.02,1.,.102)
        l = legend(bbox_to_anchor=dims, loc=3, ncol=ncol, mode='expand',
                       borderaxespad=0., handlelength=hlen, prop={'size':fsz},
                       columnspacing=1, handletextpad=0.05, *args, **kwargs) # , frameon=frameon
    else:
        l = legend(loc=loc, prop={'size':fsz}, ncol=ncol, handlelength=hlen, *args, **kwargs) # , frameon=frameon
    if not frameon:
        l.get_frame().set_linewidth(0.0)
    return l


def FigGrid(color='grey', zorder=-100):
    """
    Figure Grid
    ===========
    """
    grid (color=color, zorder=zorder)


def Gll(xl, yl, leg=True, grd=True, leg_ncol=1, leg_loc='best', leg_out=False, leg_hlen=3, leg_fsz=8, leg_out_dims=None, leg_frameon=False, hide_trframe=True, *args, **kwargs):
    """
    FigGrid, labels and legend
    ==========================
    """
    xlabel(xl)
    ylabel(yl)
    if hide_trframe: HideFrameLines()
    if grd: FigGrid()
    if leg: return Leg(frameon=leg_frameon, fsz=leg_fsz, ncol=leg_ncol, loc=leg_loc, out=leg_out, hlen=leg_hlen, out_dims=leg_out_dims, *args, **kwargs)


def Cross(x0=0.0, y0=0.0, clr='black', ls='dashed', lw=1, zorder=0):
    """
    Draw cross through zero
    =======================
    """
    axvline(x0, color=clr, linestyle=ls, linewidth=lw, zorder=zorder)
    axhline(y0, color=clr, linestyle=ls, linewidth=lw, zorder=zorder)


def LodeCross(m=1.2, clr='black', ls='dashed', lw=1, zorder=0, **kwargs):
    """
    Draw cross for Lode diagram
    ===========================
    """
    xmin, xmax = gca().get_xlim()
    ymin, ymax = gca().get_ylim()
    if xmin > 0.0: xmin = 0.0
    if xmax < 0.0: xmax = 0.0
    if ymin > 0.0: ymin = 0.0
    if ymax < 0.0: ymax = 0.0
    dx = abs(xmax - xmin)
    dy = abs(ymax - ymin)
    d  = m * max(dx, dy)
    #print('dx=',dx,' dy=',dy,' d=',d)
    #d = sqrt((xmax-xmin)*2.0 + (ymax-ymin)*2.0)
    c30, s30 = sqrt(3.0)/2.0, 0.5
    c60, s60 = 0.5, sqrt(3.0)/2.0
    plot([0,0],     [0,d*c30], color=clr, ls=ls, lw=lw, zorder=zorder, **kwargs)
    plot([0,-d*s30],[0,d*c30], color=clr, ls=ls, lw=lw, zorder=zorder, **kwargs)
    plot([0,-d*s60],[0,d*c60], color=clr, ls=ls, lw=lw, zorder=zorder, **kwargs)


def Text(x, y, txt, x_offset=0, y_offset=0, units='points', va='bottom', ha='left', color='black', fontsize=10):
    """
    Add text
    ========
    """
    trans = offset_copy(gca().transData, fig=gcf(), x=x_offset, y=y_offset, units=units)
    text(x, y, txt, transform=trans, va=va, ha=ha, color=color, fontsize=fontsize)


def TextBox(x, y, txt, fsz=10, ha='left'):
    """
    Add text inside box
    ===================
    """
    text(x, y, txt, bbox={'facecolor':'white'}, fontsize=fsz, ha=ha)


def Circle(xc,yc,R, ec='red', fc='None', lw=1, ls='solid', zorder=None, clip_on=False):
    """
    Draw circle
    ===========
    """
    gca().add_patch(MPLCircle((xc,yc), R, clip_on=clip_on, ls=ls, edgecolor=ec, facecolor=fc, lw=lw, zorder=zorder))


def Arc(xc,yc,R, alp_min=0.0, alp_max=pi, clr='red', lw=1, ls='solid', zorder=None):
    """
    Draw arc
    ========
    """
    gca().add_patch(MPLArc((xc,yc), 2.*R,2.*R, clip_on=False, angle=0, theta1=alp_min*180.0/pi, theta2=alp_max*180.0/pi, ls=ls, color=clr, lw=lw, zorder=zorder))


def Arrow(xi,yi, xf,yf, sc=20, fc='#a2e3a2', ec='black', zo=0, st='simple', label=None, lbHa='left', lbVa='center', lbPos='cen', lbFs=10, lbRot=0, lbDx=0, lbDy=0, *args, **kwargs):
    """
    Draw arrow
    ==========
    default arguments:
        mutation_scale = sc
        fc             = fc
        ec             = ec
        zorder         = zo
        arrowstyle     = st
        lbPos          = {'cen', 'tip', 'tail'}
    styles:
        Curve           -        None
        CurveB          ->       head_length=0.4,head_width=0.2
        BracketB        -[       widthB=1.0,lengthB=0.2,angleB=None
        CurveFilledB    -|>      head_length=0.4,head_width=0.2
        CurveA          <-       head_length=0.4,head_width=0.2
        CurveAB         <->      head_length=0.4,head_width=0.2
        CurveFilledA    <|-      head_length=0.4,head_width=0.2
        CurveFilledAB   <|-|>    head_length=0.4,head_width=0.2
        BracketA        ]-       widthA=1.0,lengthA=0.2,angleA=None
        BracketAB       ]-[      widthA=1.0,lengthA=0.2,angleA=None,widthB=1.0,lengthB=0.2,angleB=None
        Fancy           fancy    head_length=0.4,head_width=0.4,tail_width=0.4
        Simple          simple   head_length=0.5,head_width=0.5,tail_width=0.2
        Wedge           wedge    tail_width=0.3,shrink_factor=0.5
        BarAB           |-|      widthA=1.0,angleA=None,widthB=1.0,angleB=None
    """
    if not 'mutation_scale' in kwargs: kwargs['mutation_scale'] = sc
    if not 'fc'             in kwargs: kwargs['fc']             = fc
    if not 'ec'             in kwargs: kwargs['ec']             = ec
    if not 'zorder'         in kwargs: kwargs['zorder']         = zo
    if not 'arrowstyle'     in kwargs: kwargs['arrowstyle']     = st
    if True: # makes arrow tips sharper
        fa = FancyArrowPatch((xi,yi), (xf,yf), shrinkA=False, shrinkB=False, path_effects=[Stroke(joinstyle='miter')], *args, **kwargs)
    else:
        fa = FancyArrowPatch((xi,yi), (xf,yf), shrinkA=False, shrinkB=False, *args, **kwargs)
    gca().add_patch(fa)
    if label:
        if lbPos == 'cen':
            xm = (xi + xf) / 2.0
            ym = (yi + yf) / 2.0
        if lbPos == 'tip':
            xm, ym = xf,yf
        if lbPos == 'tail':
            xm, ym = xi,yi
        gca().text(xm+lbDx, ym+lbDy, label, ha=lbHa, va=lbVa, fontsize=lbFs, rotation=lbRot)
    return fa


def Vhatch(xi=1.0, xf=2.0, fbot=lambda x: 0.0, ftop=lambda x: x, np=41, *args, **kwargs):
    """
    Draw hatch from xi to xf between two functions fbot(x) and ftop(x)
    ====================================================================
    """
    dx = (xf - xi) / float(np-1)
    for i in range(np):
        x = xi + float(i) * dx
        plot([x,x], [fbot(x),ftop(x)], *args, **kwargs)


def Axes(xi,yi, xf,yf, xlab=r'$x$', ylab=r'$y$', xHa='right', xVa='top', yHa='right', yVa='bottom', sc=10, st='-|>', xDel=0.04, yDel=0.04):
    """
    Axes substitutes axis frame
    ===========================
    """
    axis('off')
    a = Arrow(xi,yi, xf,yi, sc=sc, st=st, ec='black', fc='black', clip_on=0)
    b = Arrow(xi,yi, xi,yf, sc=sc, st=st, ec='black', fc='black', clip_on=0)
    c = text(xf,yi-xDel,xlab,ha=xHa,va=xVa)
    d = text(xi-yDel,yf,ylab,ha=yHa,va=yVa)
    return [a,b,c,d]


def Quad(x0,y0, x1,y1, x2,y2, x3,y3, fc='#e1eeff', ec='black', zorder=0, alpha=1.0, ls='solid', lw=1, clip_on=1):
    """
    Draw quad = tetragon
    ====================
    """
    gca().add_patch(Polygon(array([[x0,y0],[x1,y1],[x2,y2],[x3,y3]]), ec=ec, fc=fc, ls=ls, zorder=zorder, alpha=alpha, lw=lw, clip_on=clip_on))


def Contour(X,Y,Z, label='', levels=None, cmapidx=0, colors=None, fmt='%g', lwd=1, fsz=10,
        inline=0, wire=True, cbar=True, zorder=None, markZero='', clabels=True):
    """
    Plot contour
    ============
    """
    L = None
    if levels != None:
        if not hasattr(levels, "__iter__"): # not a list or array...
            levels = linspace(Z.min(), Z.max(), levels)
    if colors==None:
        c1 = contourf (X,Y,Z, cmap=Cmap(cmapidx), levels=levels, zorder=None)
    else:
        c1 = contourf (X,Y,Z, colors=colors, levels=levels, zorder=None)
    if wire:
        c2 = contour (X,Y,Z, colors=('k'), levels=levels, linewidths=[lwd], zorder=None)
        if clabels:
            clabel (c2, inline=inline, fontsize=fsz)
    if cbar:
        cb = colorbar (c1, format=fmt)
        cb.ax.set_ylabel (label)
    if markZero:
        c3 = contour(X,Y,Z, levels=[0], colors=[markZero], linewidths=[2])
        if clabels:
            clabel(c3, inline=inline, fontsize=fsz)


def GetClr(idx=0, scheme=1): # color
    """
    Get ordered color
    =================
    """
    if scheme==1:
        C = ['blue', 'green', 'magenta', 'orange', 'red', 'cyan', 'black', '#de9700', '#89009d', '#7ad473', '#737ad4', '#d473ce', '#7e6322', '#462222', '#98ac9d', '#37a3e8', 'yellow']
    elif scheme==2:
        C = ['red', 'green', 'blue', 'magenta', 'cyan', 'black', 'orange', '#89009d']
    else:
        C = ['#89009d', '#7ad473', '#737ad4', 'red', '#d473ce', '#de9700', '#7e6322', '#462222', '#98ac9d', '#37a3e8', 'yellow', 'blue', 'green', 'red', 'cyan', 'magenta', 'orange', 'black']
    return C[idx % len(C)]


def GetLightClr(idx=0, scheme=1): # color
    """
    Get ordered light color
    =======================
    """
    #if scheme==1:
    C = ['#64f1c1', '#d2e5ff', '#fff0d2', '#bdb6b9', '#a6c9b7', '#c7c9a6', '#a6a6c9', '#c9a6bf', '#de9700', '#89009d', '#7ad473', '#737ad4', '#d473ce', '#7e6322', '#462222', '#98acdd', '#37a3e8', 'yellow', 'blue', 'green', 'magenta', 'orange', 'red', 'cyan', 'black']
    return C[idx % len(C)]


def GetMrk(idx=0, scheme=1): # marker
    """
    Get marker
    ==========
    """
    if scheme==1:
        M = ['o', '^', '*', 'd', 'v', 's', '<', '>', 'p', 'h', 'D']
    else:
        M = ['+', '1', '2', '3', '4', 'x', '.', '|', '', '_', 'H', '*']
    return M[idx % len(M)]


def GetLst(idx=0): # linestyle
    """
    Get ordered line style
    ======================
    """
    L = ['solid', 'dashed', 'dotted']
    #L = ['solid', 'dashed', 'dash_dot', 'dotted']
    return L[idx % len(L)]


def Read(filename, int_cols=[], make_maps=False):
    """
    Read file with table
    ====================
    dat: dictionary with the following content:
      dat = {'sx':[1,2,3],'ex':[0,1,2]}
      int_cols = ['idx','num']
    """
    filename = os.path.expandvars(filename)
    if not os.path.isfile(filename): raise Exception("[1;31mRead: could not find file <[1;34m%s[0m[1;31m>[0m"%filename)
    file   = open(filename,'r')
    header = file.readline().split()
    if len(header) == 0:
        raise Exception('[1;31mRead: reading header of file <%s> failed: the first line in file must contain the header. ex: time ux uy uz[0m'%filename)
    #while len(header) == 0:
        #header = file.readline().split()
    dat = {}
    im  = ['id','Id','tag','Tag'] # int keys to be mapped
    fm  = ['time','Time']         # float keys to be mapped
    k2m = []                      # all keys to be mapped
    k2m.extend(im)
    k2m.extend(fm)
    int_cols.extend(im)
    for key in header:
        dat[key] = []
        if make_maps:
            if key in k2m: dat['%s2row'%key] = {}
    row = 0
    for lin in file:
        res = lin.split()
        if len(res) == 0: continue
        for i, key in enumerate(header):
            if key in int_cols: dat[key].append(int  (res[i]))
            else:               dat[key].append(float(res[i]))
            if make_maps:
                if key in k2m:
                    if dat[key][row] in dat['%s2row'%key]:
                        if type(dat['%s2row'%key][dat[key][row]]) == int:
                            dat['%s2row'%key][dat[key][row]] = [dat['%s2row'%key][dat[key][row]], row]
                        else:
                            dat['%s2row'%key][dat[key][row]].append(row)
                    else:
                        dat['%s2row'%key][dat[key][row]] = row
        row += 1
    file.close()
    mkeys = ['%s2row'%k for k in k2m]
    for key, val in dat.items(): # convert lists to arrays (floats only)
        if key in int_cols or key in mkeys: continue
        dat[key] = array(val)
    return dat


def ReadT(fnkey, ids, arc_lens=[]):
    """
    Read table with temporal data
    =============================
    """
    r0  = Read('%s_%d.res'%(fnkey,ids[0]))
    str_time = 'Time'
    if 'time' in r0: str_time = 'time'
    if 't'    in r0: str_time = 't'
    nt  = len(r0[str_time])
    np  = len(ids)
    if len(arc_lens)>0: dat = {'arc_len':zeros((np,nt))}
    else:               dat = {}
    for k, v in r0.items(): dat[k] = zeros((np,nt))
    for i, n in enumerate(ids):
        r = Read('%s_%d.res'%(fnkey,n), make_maps=False)
        for k, v in r.items(): dat[k][i,:] = v
        if len(arc_lens)>0:
            dat['arc_len'][i,:] = arc_lens[i]
    return dat


def ReadMany(filenames, int_cols=[]):
    """
    Read many files with tables
    ===========================
    NOTE: filenames must be given in the correct 'time' order
    """
    ddat = Read(filenames[0], make_maps=False)
    for fn in filenames[1:]:
        d = Read(fn, int_cols, make_maps=False)
        for key, val in d.items():
            if key in ddat:
                l = ddat[key].tolist()
                l.extend(val)
                ddat[key] = array(l)
        for key in ddat.keys():
            if not key in d.keys():
                ddat.pop(key)
    return ddat


def RadForm(x, pos=None):
    """
    Radians formatting
    ==================
    """
    n = int((x/(pi/6.0))+pi/12.0)
    if n== 0: return r'$0$'
    if n== 1: return r'$\frac{\pi}{6}$'
    if n== 2: return r'$\frac{\pi}{3}$'
    if n== 3: return r'$\frac{\pi}{2}$'
    if n== 4: return r'$2\frac{\pi}{3}$'
    if n== 6: return r'$\pi$'
    if n== 8: return r'$4\frac{\pi}{3}$'
    if n== 9: return r'$3\frac{\pi}{2}$'
    if n==10: return r'$5\frac{\pi}{3}$'
    if n==12: return r'$2\pi$'
    return r'$%d\frac{\pi}{6}$'%n


def RadDegForm(x, pos=None):
    """
    Radians and degrees formatting
    ==============================
    """
    n = int((x/(pi/6.0))+pi/12.0)
    if n== 0: return r'$0$'
    if n== 1: return r'$\frac{\pi}{6}$''\n$30^\circ$'
    if n== 2: return r'$\frac{\pi}{3}$''\n$60^\circ$'
    if n== 3: return r'$\frac{\pi}{2}$''\n$90^\circ$'
    if n== 4: return r'$2\frac{\pi}{3}$''\n$120^\circ$'
    if n== 6: return r'$\pi$''\n$180^\circ$'
    if n== 8: return r'$4\frac{\pi}{3}$''\n$240^\circ$'
    if n== 9: return r'$3\frac{\pi}{2}$''\n$270^\circ$'
    if n==10: return r'$5\frac{\pi}{3}$''\n$300^\circ$'
    if n==12: return r'$2\pi$''\n$360^\circ$'
    return r'$%d\frac{\pi}{6}$''\n$%g^\circ$'%(n,x*180.0/pi)


def RadFmt(gen_ticks=False, rad_and_deg=True):
    """
    Set radians or degrees formatting
    =================================
    """
    if gen_ticks: xticks (linspace(0,2*pi,13))
    if rad_and_deg: gca().xaxis.set_major_formatter(FuncFormatter(RadDegForm))
    else:           gca().xaxis.set_major_formatter(FuncFormatter(RadForm))


def ColumnNodes(nc=10, o2=True):
    """
    Column nodes
    ============
    nc: number of cells along y
    """
    ny = nc + 1 # number of rows along y
    l  = arange(ny) * 2
    r  = arange(ny) * 2 + 1
    if o2:
        c = 2 * ny + arange(ny)
        m = 3 * ny + arange(nc) * 2
        L = zeros(ny + nc, dtype=int)
        i = arange(ny + nc)
        L[i%2==0] = l
        L[i%2==1] = m
        R = L + 1
        return l, c, r, L, R
    return l, r


def GetITout(all_output_times, time_stations_out, tol=1.0e-8, with_count=False):
    """
    Get indices and output times
    ============================
    INPUT:
        all_output_times  : array with all output times. ex: [0,0.1,0.2,0.22,0.3,0.4]
        time_stations_out : time stations for output: ex: [0,0.2,0.4]  # must be sorted ascending
        tol               : tolerance to compare times
    RETURNS:
        iout : indices of times in all_output_times
        tout : times corresponding to iout
    """
    I, T        = [], []                          # indices and times for output
    lower_index = 0                               # lower index in all_output_times
    len_aotimes = len(all_output_times)           # length of a_o_times
    for t in time_stations_out:                   # for all requested output times
        if t < 0:                                 # final time
            k = len(all_output_times)-1           # last output time
            I.append(k)                           # append index to iout
            T.append(all_output_times[k])         # append last output time
            continue                              # skip search
        for k in range(lower_index, len_aotimes): # search within a_o_times
            if abs(t-all_output_times[k]) < tol:  # found near match
                lower_index += 1                  # update index
                I.append(k)                       # add index to iout
                T.append(all_output_times[k])     # add time to tout
                break                             # stop searching for this 't'
            if all_output_times[k] > t:           # failed to search for 't'
                lower_index = k                   # update idx to start from here on
                break                             # skip this 't' and try the next one
    if with_count: return zip(range(len(I)), I, T)
    else:          return zip(I, T)


def PrintV(name, V, nf='%12.5f', Tol=1.0e-14):
    """
    Pretty-print a vector
    =====================
    """
    print(name)
    lin = ''
    for i in range(len(V)):
        if abs(V[i])<Tol: lin += nf % 0
        else:             lin += nf % V[i]
    lin += '\n'
    print(lin)


def PrintM(name, A, nf='%12.5f', Tol=1.0e-14):
    """
    Pretty-print a matrix
    =====================
    """
    print(name)
    m = A.shape[0] # number of rows
    n = A.shape[1] # number of columns
    lin = ''       # empty string
    for i in range(m):
        for j in range(n):
            if abs(A[i,j])<Tol: lin += nf % 0
            else:               lin += nf % A[i,j]
        lin += '\n'
    print(lin)


def PrintDiff(fmt, a, b, tol):
    """
    Print colored difference between a and b
    ========================================
    """
    if abs(a-b) < tol: print('[1;32m' + fmt % abs(a-b) + '[0m')
    else:              print('[1;31m' + fmt % abs(a-b) + '[0m')


def CompareArrays(a, b, tol=1.0e-12, table=False, namea='a', nameb='b', numformat='%17.10e', spaces=17, stride=0):
    """
    Compare two arrays
    ==================
    """
    if table:
        sfmt = '%%%ds' % spaces
        print('='*(spaces*3+4+5))
        print(sfmt % namea, sfmt % nameb, '[1;37m', sfmt % 'diff', '[0m')
        print('-'*(spaces*3+4+5))
    max_diff = 0.0
    for i in range(len(a)):
        diff = abs(a[i]-b[i])
        if diff > max_diff: max_diff = diff
        if table:
            clr = '[1;31m' if diff > tol else '[1;32m'
            print('%4d'%(i+stride), numformat % a[i], numformat % b[i], clr, numformat % diff, '[0m')
    if table: print('='*(spaces*3+4+5))
    if max_diff < tol: print('max difference = [1;32m%20.15e[0m' % max_diff)
    else:              print('max difference = [1;31m%20.15e[0m' % max_diff)


# testing ############################################################################################
if __name__=='__main__':
    prob = 0

    if prob == 0:
        x = linspace (0, 10, 100)
        y = x**1.5
        plot    (x,y, 'b-', label='sim')
        ver = int(sys.version[0])
        if ver<3:
            Arc (0,0,10)
            Arc (0,0,20, clr='magenta')
        Arrow   (-10,0,max(x),max(y))
        Text    (0,25,r'$\sigma$')
        Text    (0,25,r'$\sigma$',y_offset=-10)
        axvline (0,color='black',zorder=-1)
        axhline (0,color='black',zorder=-1)
        FigGrid ()
        axis    ('equal')
        legend  (loc='upper left')
        xlabel  (r'$x$')
        ylabel  (r'$y$')
        show    ()

    if prob == 1:
        ITout = get_itout([0., 0.1, 0.15, 0.2, 0.23, 0.23, 0.23, 0.3, 0.8, 0.99],
                          [0., 0.1, 0.2, 0.3, 0.4, 0.6, 0.7, 0.8, 0.9, 1.])
        print(ITout)

    if prob == 2:
        ITout = get_itout([0.0, 0.1, 0.2, 0.30000000000000004, 0.4, 0.5, 0.6, 0.7, 0.7999999999999999, 0.8999999999999999, 0.9999999999999999],
                          [0.1, 0.1, 0.2, 0.5, 1.0])
        print(ITout)
