Triangle
A Two-Dimensional Quality Mesh Generator and Delaunay Triangulator.
Version 1.6

Show Me
A Display Program for Meshes and More.
Version 1.6

Copyright 1993, 1995, 1997, 1998, 2002, 2005 Jonathan Richard Shewchuk
2360 Woolsey #H
Berkeley, California  94705-1927
Please send bugs and comments to jrs@cs.berkeley.edu

Created as part of the Quake project (tools for earthquake simulation).
Supported in part by NSF Grant CMS-9318163 and an NSERC 1967 Scholarship.
There is no warranty whatsoever.  Use at your own risk.


Triangle generates exact Delaunay triangulations, constrained Delaunay
triangulations, conforming Delaunay triangulations, Voronoi diagrams, and
high-quality triangular meshes.  The latter can be generated with no small
or large angles, and are thus suitable for finite element analysis.
Show Me graphically displays the contents of the geometric files used by
Triangle.  Show Me can also write images in PostScript form.

Information on the algorithms used by Triangle, including complete
references, can be found in the comments at the beginning of the triangle.c
source file.  Another listing of these references, with PostScript copies
of some of the papers, is available from the Web page

    http://www.cs.cmu.edu/~quake/triangle.research.html

------------------------------------------------------------------------------

These programs may be freely redistributed under the condition that the
copyright notices (including the copy of this notice in the code comments
and the copyright notice printed when the `-h' switch is selected) are
not removed, and no compensation is received.  Private, research, and
institutional use is free.  You may distribute modified versions of this
code UNDER THE CONDITION THAT THIS CODE AND ANY MODIFICATIONS MADE TO IT
IN THE SAME FILE REMAIN UNDER COPYRIGHT OF THE ORIGINAL AUTHOR, BOTH
SOURCE AND OBJECT CODE ARE MADE FREELY AVAILABLE WITHOUT CHARGE, AND
CLEAR NOTICE IS GIVEN OF THE MODIFICATIONS.  Distribution of this code as
part of a commercial system is permissible ONLY BY DIRECT ARRANGEMENT
WITH THE AUTHOR.  (If you are not directly supplying this code to a
customer, and you are instead telling them how they can obtain it for
free, then you are not required to make any arrangement with me.)

------------------------------------------------------------------------------

The files included in this distribution are:

  README           The file you're reading now.
  triangle.c       Complete C source code for Triangle.
  showme.c         Complete C source code for Show Me.
  triangle.h       Include file for calling Triangle from another program.
  tricall.c        Sample program that calls Triangle.
  makefile         Makefile for compiling Triangle and Show Me.
  A.poly           A sample input file.

Each of Triangle and Show Me is a single portable C file.  The easiest way
to compile them is to edit and use the included makefile.  Before
compiling, read the makefile, which describes your options, and edit it
accordingly.  You should specify:

  The source and binary directories.

  The C compiler and level of optimization.

  The "correct" directories for include files (especially X include files),
  if necessary.

  Do you want single precision or double?  (The default is double.)  Do you
  want to leave out some of Triangle's features to reduce the size of the
  executable file?  Investigate the SINGLE, REDUCED, and CDT_ONLY symbols.

  If yours is not a Unix system, define the NO_TIMER symbol to remove the
  Unix-specific timing code.  Also, don't try to compile Show Me; it only
  works with X Windows.

  If you are compiling on an Intel x86 CPU and using gcc w/Linux or
  Microsoft C, be sure to define the LINUX or CPU86 (for Microsoft) symbol
  during compilation so that the exact arithmetic works right.

Once you've done this, type "make" to compile the programs.  Alternatively,
the files are usually easy to compile without a makefile:

  cc -O -o triangle triangle.c -lm
  cc -O -o showme showme.c -lX11

On some systems, the C compiler won't be able to find the X include files
or libraries, and you'll need to specify an include path or library path:

  cc -O -I/usr/local/include -o showme showme.c -L/usr/local/lib -lX11

Some processors, including Intel x86 family and possibly Motorola 68xxx
family chips, are IEEE conformant but have extended length internal
floating-point registers that may defeat Triangle's exact arithmetic
routines by failing to cause enough roundoff error!  Typically, there is a
way to set these internal registers so that they are rounded off to IEEE
single or double precision format.  I believe (but I'm not certain) that
Triangle has the right incantations for x86 chips, if you have gcc running
under Linux (define the LINUX compiler symbol) or Microsoft C (define the
CPU86 compiler symbol).

If you have a different processor or operating system, or if I got the
incantations wrong, you should check your C compiler or system manuals to
find out how to configure these internal registers to the precision you are
using.  Otherwise, the exact arithmetic routines won't be exact at all.
See http://www.cs.cmu.edu/~quake/robust.pc.html for details.  Triangle's
exact arithmetic hasn't a hope of working on machines like the Cray C90 or
Y-MP, which are not IEEE conformant and have inaccurate rounding.

Triangle and Show Me have both text and HTML documentation.  The latter is
illustrated.  Find it on the Web at

  http://www.cs.cmu.edu/~quake/triangle.html
  http://www.cs.cmu.edu/~quake/showme.html

Complete text instructions are printed by invoking each program with the
`-h' switch:

  triangle -h
  showme -h

The instructions are long; you'll probably want to pipe the output to
`more' or `lpr' or redirect it to a file.

Both programs give a short list of command line options if they are invoked
without arguments (that is, just type `triangle' or `showme').

Try out Triangle on the enclosed sample file, A.poly:

  triangle -p A
  showme A.poly &

Triangle will read the Planar Straight Line Graph defined by A.poly, and
write its constrained Delaunay triangulation to A.1.node and A.1.ele.
Show Me will display the figure defined by A.poly.  There are two buttons
marked "ele" in the Show Me window; click on the top one.  This will cause
Show Me to load and display the triangulation.

For contrast, try running

  triangle -pq A

Now, click on the same "ele" button.  A new triangulation will be loaded;
this one having no angles smaller than 20 degrees.

To see a Voronoi diagram, try this:

  cp A.poly A.node
  triangle -v A

Click the "ele" button again.  You will see the Delaunay triangulation of
the points in A.poly, without the segments.  Now click the top "voro" button.
You will see the Voronoi diagram corresponding to that Delaunay triangulation.
Click the "Reset" button to see the full extent of the diagram.

------------------------------------------------------------------------------

If you wish to call Triangle from another program, instructions for doing
so are contained in the file `triangle.h' (but read Triangle's regular
instructions first!).  Also look at `tricall.c', which provides an example
of how to call Triangle.

Type "make trilibrary" to create triangle.o, a callable object file.
Alternatively, the object file is usually easy to compile without a
makefile:

  cc -DTRILIBRARY -O -c triangle.c

Type "make distclean" to remove all the object and executable files created
by make.

------------------------------------------------------------------------------

If you use Triangle, and especially if you use it to accomplish real work,
I would like very much to hear from you.  A short letter or email (to
jrs@cs.berkeley.edu) describing how you use Triangle will mean a lot to me.
The more people I know are using this program, the more easily I can
justify spending time on improvements and on the three-dimensional
successor to Triangle, which in turn will benefit you.  Also, I can put you
on a list to receive email whenever a new version of Triangle is available.

If you use a mesh generated by Triangle or plotted by Show Me in a
publication, please include an acknowledgment as well.  And please spell
Triangle with a capital `T'!  If you want to include a citation, use
`Jonathan Richard Shewchuk, ``Triangle:  Engineering a 2D Quality Mesh
Generator and Delaunay Triangulator,'' in Applied Computational Geometry:
Towards Geometric Engineering (Ming C. Lin and Dinesh Manocha, editors),
volume 1148 of Lecture Notes in Computer Science, pages 203-222,
Springer-Verlag, Berlin, May 1996.  (From the First ACM Workshop on Applied
Computational Geometry.)'


Jonathan Richard Shewchuk
July 27, 2005






/*****************************************************************************/
/*                                                                           */
/*  (triangle.h)                                                             */
/*                                                                           */
/*  Include file for programs that call Triangle.                            */
/*                                                                           */
/*  Accompanies Triangle Version 1.6                                         */
/*  July 28, 2005                                                            */
/*                                                                           */
/*  Copyright 1996, 2005                                                     */
/*  Jonathan Richard Shewchuk                                                */
/*  2360 Woolsey #H                                                          */
/*  Berkeley, California  94705-1927                                         */
/*  jrs@cs.berkeley.edu                                                      */
/*                                                                           */
/*****************************************************************************/

/*****************************************************************************/
/*                                                                           */
/*  How to call Triangle from another program                                */
/*                                                                           */
/*                                                                           */
/*  If you haven't read Triangle's instructions (run "triangle -h" to read   */
/*  them), you won't understand what follows.                                */
/*                                                                           */
/*  Triangle must be compiled into an object file (triangle.o) with the      */
/*  TRILIBRARY symbol defined (generally by using the -DTRILIBRARY compiler  */
/*  switch).  The makefile included with Triangle will do this for you if    */
/*  you run "make trilibrary".  The resulting object file can be called via  */
/*  the procedure triangulate().                                             */
/*                                                                           */
/*  If the size of the object file is important to you, you may wish to      */
/*  generate a reduced version of triangle.o.  The REDUCED symbol gets rid   */
/*  of all features that are primarily of research interest.  Specifically,  */
/*  the -DREDUCED switch eliminates Triangle's -i, -F, -s, and -C switches.  */
/*  The CDT_ONLY symbol gets rid of all meshing algorithms above and beyond  */
/*  constrained Delaunay triangulation.  Specifically, the -DCDT_ONLY switch */
/*  eliminates Triangle's -r, -q, -a, -u, -D, -Y, -S, and -s switches.       */
/*                                                                           */
/*  IMPORTANT:  These definitions (TRILIBRARY, REDUCED, CDT_ONLY) must be    */
/*  made in the makefile or in triangle.c itself.  Putting these definitions */
/*  in this file (triangle.h) will not create the desired effect.            */
/*                                                                           */
/*                                                                           */
/*  The calling convention for triangulate() follows.                        */
/*                                                                           */
/*      void triangulate(triswitches, in, out, vorout)                       */
/*      char *triswitches;                                                   */
/*      struct triangulateio *in;                                            */
/*      struct triangulateio *out;                                           */
/*      struct triangulateio *vorout;                                        */
/*                                                                           */
/*  `triswitches' is a string containing the command line switches you wish  */
/*  to invoke.  No initial dash is required.  Some suggestions:              */
/*                                                                           */
/*  - You'll probably find it convenient to use the `z' switch so that       */
/*    points (and other items) are numbered from zero.  This simplifies      */
/*    indexing, because the first item of any type always starts at index    */
/*    [0] of the corresponding array, whether that item's number is zero or  */
/*    one.                                                                   */
/*  - You'll probably want to use the `Q' (quiet) switch in your final code, */
/*    but you can take advantage of Triangle's printed output (including the */
/*    `V' switch) while debugging.                                           */
/*  - If you are not using the `q', `a', `u', `D', `j', or `s' switches,     */
/*    then the output points will be identical to the input points, except   */
/*    possibly for the boundary markers.  If you don't need the boundary     */
/*    markers, you should use the `N' (no nodes output) switch to save       */
/*    memory.  (If you do need boundary markers, but need to save memory, a  */
/*    good nasty trick is to set out->pointlist equal to in->pointlist       */
/*    before calling triangulate(), so that Triangle overwrites the input    */
/*    points with identical copies.)                                         */
/*  - The `I' (no iteration numbers) and `g' (.off file output) switches     */
/*    have no effect when Triangle is compiled with TRILIBRARY defined.      */
/*                                                                           */
/*  `in', `out', and `vorout' are descriptions of the input, the output,     */
/*  and the Voronoi output.  If the `v' (Voronoi output) switch is not used, */
/*  `vorout' may be NULL.  `in' and `out' may never be NULL.                 */
/*                                                                           */
/*  Certain fields of the input and output structures must be initialized,   */
/*  as described below.                                                      */
/*                                                                           */
/*****************************************************************************/

/*****************************************************************************/
/*                                                                           */
/*  The `triangulateio' structure.                                           */
/*                                                                           */
/*  Used to pass data into and out of the triangulate() procedure.           */
/*                                                                           */
/*                                                                           */
/*  Arrays are used to store points, triangles, markers, and so forth.  In   */
/*  all cases, the first item in any array is stored starting at index [0].  */
/*  However, that item is item number `1' unless the `z' switch is used, in  */
/*  which case it is item number `0'.  Hence, you may find it easier to      */
/*  index points (and triangles in the neighbor list) if you use the `z'     */
/*  switch.  Unless, of course, you're calling Triangle from a Fortran       */
/*  program.                                                                 */
/*                                                                           */
/*  Description of fields (except the `numberof' fields, which are obvious): */
/*                                                                           */
/*  `pointlist':  An array of point coordinates.  The first point's x        */
/*    coordinate is at index [0] and its y coordinate at index [1], followed */
/*    by the coordinates of the remaining points.  Each point occupies two   */
/*    REALs.                                                                 */
/*  `pointattributelist':  An array of point attributes.  Each point's       */
/*    attributes occupy `numberofpointattributes' REALs.                     */
/*  `pointmarkerlist':  An array of point markers; one int per point.        */
/*                                                                           */
/*  `trianglelist':  An array of triangle corners.  The first triangle's     */
/*    first corner is at index [0], followed by its other two corners in     */
/*    counterclockwise order, followed by any other nodes if the triangle    */
/*    represents a nonlinear element.  Each triangle occupies                */
/*    `numberofcorners' ints.                                                */
/*  `triangleattributelist':  An array of triangle attributes.  Each         */
/*    triangle's attributes occupy `numberoftriangleattributes' REALs.       */
/*  `trianglearealist':  An array of triangle area constraints; one REAL per */
/*    triangle.  Input only.                                                 */
/*  `neighborlist':  An array of triangle neighbors; three ints per          */
/*    triangle.  Output only.                                                */
/*                                                                           */
/*  `segmentlist':  An array of segment endpoints.  The first segment's      */
/*    endpoints are at indices [0] and [1], followed by the remaining        */
/*    segments.  Two ints per segment.                                       */
/*  `segmentmarkerlist':  An array of segment markers; one int per segment.  */
/*                                                                           */
/*  `holelist':  An array of holes.  The first hole's x and y coordinates    */
/*    are at indices [0] and [1], followed by the remaining holes.  Two      */
/*    REALs per hole.  Input only, although the pointer is copied to the     */
/*    output structure for your convenience.                                 */
/*                                                                           */
/*  `regionlist':  An array of regional attributes and area constraints.     */
/*    The first constraint's x and y coordinates are at indices [0] and [1], */
/*    followed by the regional attribute at index [2], followed by the       */
/*    maximum area at index [3], followed by the remaining area constraints. */
/*    Four REALs per area constraint.  Note that each regional attribute is  */
/*    used only if you select the `A' switch, and each area constraint is    */
/*    used only if you select the `a' switch (with no number following), but */
/*    omitting one of these switches does not change the memory layout.      */
/*    Input only, although the pointer is copied to the output structure for */
/*    your convenience.                                                      */
/*                                                                           */
/*  `edgelist':  An array of edge endpoints.  The first edge's endpoints are */
/*    at indices [0] and [1], followed by the remaining edges.  Two ints per */
/*    edge.  Output only.                                                    */
/*  `edgemarkerlist':  An array of edge markers; one int per edge.  Output   */
/*    only.                                                                  */
/*  `normlist':  An array of normal vectors, used for infinite rays in       */
/*    Voronoi diagrams.  The first normal vector's x and y magnitudes are    */
/*    at indices [0] and [1], followed by the remaining vectors.  For each   */
/*    finite edge in a Voronoi diagram, the normal vector written is the     */
/*    zero vector.  Two REALs per edge.  Output only.                        */
/*                                                                           */
/*                                                                           */
/*  Any input fields that Triangle will examine must be initialized.         */
/*  Furthermore, for each output array that Triangle will write to, you      */
/*  must either provide space by setting the appropriate pointer to point    */
/*  to the space you want the data written to, or you must initialize the    */
/*  pointer to NULL, which tells Triangle to allocate space for the results. */
/*  The latter option is preferable, because Triangle always knows exactly   */
/*  how much space to allocate.  The former option is provided mainly for    */
/*  people who need to call Triangle from Fortran code, though it also makes */
/*  possible some nasty space-saving tricks, like writing the output to the  */
/*  same arrays as the input.                                                */
/*                                                                           */
/*  Triangle will not free() any input or output arrays, including those it  */
/*  allocates itself; that's up to you.  You should free arrays allocated by */
/*  Triangle by calling the trifree() procedure defined below.  (By default, */
/*  trifree() just calls the standard free() library procedure, but          */
/*  applications that call triangulate() may replace trimalloc() and         */
/*  trifree() in triangle.c to use specialized memory allocators.)           */
/*                                                                           */
/*  Here's a guide to help you decide which fields you must initialize       */
/*  before you call triangulate().                                           */
/*                                                                           */
/*  `in':                                                                    */
/*                                                                           */
/*    - `pointlist' must always point to a list of points; `numberofpoints'  */
/*      and `numberofpointattributes' must be properly set.                  */
/*      `pointmarkerlist' must either be set to NULL (in which case all      */
/*      markers default to zero), or must point to a list of markers.  If    */
/*      `numberofpointattributes' is not zero, `pointattributelist' must     */
/*      point to a list of point attributes.                                 */
/*    - If the `r' switch is used, `trianglelist' must point to a list of    */
/*      triangles, and `numberoftriangles', `numberofcorners', and           */
/*      `numberoftriangleattributes' must be properly set.  If               */
/*      `numberoftriangleattributes' is not zero, `triangleattributelist'    */
/*      must point to a list of triangle attributes.  If the `a' switch is   */
/*      used (with no number following), `trianglearealist' must point to a  */
/*      list of triangle area constraints.  `neighborlist' may be ignored.   */
/*    - If the `p' switch is used, `segmentlist' must point to a list of     */
/*      segments, `numberofsegments' must be properly set, and               */
/*      `segmentmarkerlist' must either be set to NULL (in which case all    */
/*      markers default to zero), or must point to a list of markers.        */
/*    - If the `p' switch is used without the `r' switch, then               */
/*      `numberofholes' and `numberofregions' must be properly set.  If      */
/*      `numberofholes' is not zero, `holelist' must point to a list of      */
/*      holes.  If `numberofregions' is not zero, `regionlist' must point to */
/*      a list of region constraints.                                        */
/*    - If the `p' switch is used, `holelist', `numberofholes',              */
/*      `regionlist', and `numberofregions' is copied to `out'.  (You can    */
/*      nonetheless get away with not initializing them if the `r' switch is */
/*      used.)                                                               */
/*    - `edgelist', `edgemarkerlist', `normlist', and `numberofedges' may be */
/*      ignored.                                                             */
/*                                                                           */
/*  `out':                                                                   */
/*                                                                           */
/*    - `pointlist' must be initialized (NULL or pointing to memory) unless  */
/*      the `N' switch is used.  `pointmarkerlist' must be initialized       */
/*      unless the `N' or `B' switch is used.  If `N' is not used and        */
/*      `in->numberofpointattributes' is not zero, `pointattributelist' must */
/*      be initialized.                                                      */
/*    - `trianglelist' must be initialized unless the `E' switch is used.    */
/*      `neighborlist' must be initialized if the `n' switch is used.  If    */
/*      the `E' switch is not used and (`in->numberofelementattributes' is   */
/*      not zero or the `A' switch is used), `elementattributelist' must be  */
/*      initialized.  `trianglearealist' may be ignored.                     */
/*    - `segmentlist' must be initialized if the `p' or `c' switch is used,  */
/*      and the `P' switch is not used.  `segmentmarkerlist' must also be    */
/*      initialized under these circumstances unless the `B' switch is used. */
/*    - `edgelist' must be initialized if the `e' switch is used.            */
/*      `edgemarkerlist' must be initialized if the `e' switch is used and   */
/*      the `B' switch is not.                                               */
/*    - `holelist', `regionlist', `normlist', and all scalars may be ignored.*/
/*                                                                           */
/*  `vorout' (only needed if `v' switch is used):                            */
/*                                                                           */
/*    - `pointlist' must be initialized.  If `in->numberofpointattributes'   */
/*      is not zero, `pointattributelist' must be initialized.               */
/*      `pointmarkerlist' may be ignored.                                    */
/*    - `edgelist' and `normlist' must both be initialized.                  */
/*      `edgemarkerlist' may be ignored.                                     */
/*    - Everything else may be ignored.                                      */
/*                                                                           */
/*  After a call to triangulate(), the valid fields of `out' and `vorout'    */
/*  will depend, in an obvious way, on the choice of switches used.  Note    */
/*  that when the `p' switch is used, the pointers `holelist' and            */
/*  `regionlist' are copied from `in' to `out', but no new space is          */
/*  allocated; be careful that you don't free() the same array twice.  On    */
/*  the other hand, Triangle will never copy the `pointlist' pointer (or any */
/*  others); new space is allocated for `out->pointlist', or if the `N'      */
/*  switch is used, `out->pointlist' remains uninitialized.                  */
/*                                                                           */
/*  All of the meaningful `numberof' fields will be properly set; for        */
/*  instance, `numberofedges' will represent the number of edges in the      */
/*  triangulation whether or not the edges were written.  If segments are    */
/*  not used, `numberofsegments' will indicate the number of boundary edges. */
/*                                                                           */
/*****************************************************************************/
