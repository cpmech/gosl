# Gosl. gm. Geometry algorithms and structures

This package provides some functions to help with the solution of geometry problems. It also
includes some routines loosely related with geometry.

**go doc**

```
package gm // import "gosl/gm"

Package gm (geometry and meshes) implements auxiliary functions for handling
geometry and mesh structures

VARIABLES

var FactoryNurbs = facNurbsT{}
    FactoryNurbs generates NURBS'

var FactoryTfinite = facTfinite{}
    FactoryTfinite generates transfinite mappings

var (
	XDELZERO = 1e-10 // minimum distance between coordinates; i.e. xmax[i]-xmin[i] mininum
)
    consts


FUNCTIONS

func DistPointBoxN(p *PointN, b *BoxN) (dist float64)
    DistPointBoxN returns the distance of a point to the box

        NOTE: If point p lies outside box b, the distance to the nearest point on b is returned.
              If p is inside b or on its surface, zero is returned.

func DistPointLine(p, a, b *Point, tol float64, verbose bool) float64
    DistPointLine computes the distance from p to line passing through a -> b

func DistPointPoint(a, b *Point) float64
    DistPointPoint computes the unsigned distance from a to b

func DistPointPointN(p *PointN, q *PointN) (dist float64)
    DistPointPointN returns the distance between two PointN

func IsPointIn(p *Point, cMin, cMax []float64, tol float64) bool
    IsPointIn returns whether p is inside box with cMin and cMax

func IsPointInLine(p, a, b *Point, zero, told, tolin float64) bool
    IsPointInLine returns whether p is inside line passing through a and b

func MetisAdjacency(edges [][2]int, shares map[int][]int) (xadj, adjncy []int32)
    MetisAdjacency returns adjacency list as a compressed storage format for
    METIS shares is the map returned by MetisShares

func MetisPartition(edges [][2]int, npart int, recursive bool) (shares map[int][]int, objval int32, parts []int32)
    MetisPartition performs graph partitioning using METIS This function
    computes the shares, adjacency list, and partition by calling the other 3
    functions

func MetisPartitionLowLevel(npart, nvert int, xadj, adjncy []int32, recursive bool) (objval int32, parts []int32)
    MetisPartitionLowLevel performs graph partitioning using METIS

func MetisShares(edges [][2]int) (shares map[int][]int)
    MetisShares returns a map of shares owned by vertices i.e. each vertex is
    shared by a number of edges, so, we return the map vertexId => [edgeIds...]
    attached to this vertexId

    Example:

             0        1
         0-------1--------2
         |       |        |
        2|      3|       4|
         |       |        |
         3-------4--------5
             5       6

    Input. edges = {0,1}, {1,2}, {0,3}, {1,4}, {2,5}, {3,4}, {4,5}

    Output. shares = 0:(0,2), 1:(0,1,3), 2:(1,4)

        3:(2,5), 4:(3,5,6), 5:(4,6)

    (notation. vertexID:(firstEdgeID, secondEdgeID))

    NOTE: (1) the pairs or triples will have sorted edgeIDs

        (2) len(shares) = number_of_vertices

func PointsLims(pp []*Point) (cmin, cmax []float64)
    PointsLims returns the limits of a set of points

func VecDot(u, v []float64) float64
    VecDot returns the dot product between two vectors

func VecNew(m float64, u []float64) []float64
    VecNew returns a new vector scaled by m

func VecNewAdd(α float64, u []float64, β float64, v []float64) []float64
    VecNewAdd returns a new vector by adding two other vectors

        w := α*u + β*v

func VecNorm(u []float64) float64
    VecNorm returns the length (Euclidean norm) of a vector


TYPES

type BezierQuad struct {

	// input
	Q [][]float64 // control points; can be set outside

	// auxiliary
	P []float64 // a point on curve
}
    BezierQuad implements a quadratic Bezier curve

        C(t) = (1-t)² Q0  +  2 t (1-t) Q1  +  t² Q2
             = t² A  +  2 t B  +  Q0
        A = Q2 - 2 Q1 + Q0
        B = Q1 - Q0

func (o *BezierQuad) DistPoint(X []float64) float64
    DistPoint returns the distance from a point to this Bezier curve It finds
    the closest projection which is stored in P

func (o *BezierQuad) GetControlCoords() (X, Y, Z []float64)
    GetControlCoords returns the coordinates of control points as 1D arrays
    (e.g. for plotting)

func (o *BezierQuad) GetPoints(T []float64) (X, Y, Z []float64)
    GetPoints returns points along the curve for given parameter values

func (o *BezierQuad) Point(C []float64, t float64)
    Point returns the x-y-z coordinates of a point on Bezier curve

type Bin struct {
	Index   int         // index of bin
	Entries []*BinEntry // entries
}
    Bin defines one bin in Bins (holds entries for search)

func (o Bin) String() string
    String returns the string representation of one Bin

type BinEntry struct {
	ID    int         // object Id
	X     []float64   // entry coordinate (read only)
	Extra interface{} // any entity attached to this entry
}
    BinEntry holds data of an entry to bin

type Bins struct {
	Ndim int       // space dimension
	Xmin []float64 // [ndim] left/lower-most point
	Xmax []float64 // [ndim] right/upper-most point
	Xdel []float64 // [ndim] the lengths along each direction (whole box)
	Size []float64 // size of bins
	Ndiv []int     // [ndim] number of divisions along each direction
	All  []*Bin    // [nbins] all bins (there will be an extra "ghost" bin along each dimension)
	// Has unexported fields.
}
    Bins implements a set of bins holding entries and is used to fast search
    entries by given coordinates.

func (o *Bins) Append(x []float64, id int, extra interface{})
    Append adds a new entry {x, id, something} into the Bins structure

func (o Bins) CalcIndex(x []float64) int
    CalcIndex calculates the bin index where the point x is returns -1 if
    out-of-range

func (o *Bins) Clear()
    Clear clears all bins

func (o Bins) FindAlongSegment(xi, xf []float64, tol float64) []int
    FindAlongSegment gets the ids of entries that lie close to a segment

        Note: the initial (xi) and final (xf) points on segment define a bounding box to filter points

func (o Bins) FindBinByIndex(idx int) *Bin
    FindBinByIndex finds or allocate new bin corresponding to index idx

func (o Bins) FindClosest(x []float64) (idClosest int, sqDistMin float64)
    FindClosest returns the id of the entry whose coordinates are closest to x

         idClosest -- the id of the closest entity. return -1 if out-of-range or not found
         sqDistMin -- the minimum distance (squared) between x and the closest entity in the same bin

        NOTE: FindClosest does search the whole area.
              It only locates neighbours in the same bin where the given x is located.
              So, if there area no entries in the bin containing x, no entry will be found.

func (o *Bins) FindClosestAndAppend(nextID *int, x []float64, extra interface{}, radTol float64, diff func(idOld int, xNew []float64) bool) (id int, existent bool)
    FindClosestAndAppend finds closest point and, if not found, append to bins
    with a new Id

        Input:
          nextId -- is the Id of the next point. Will be incremented if x is a new point to be added.
          x      -- is the point to be added
          extra  -- extra information attached to point
          radTol -- is the tolerance for the radial distance (i.e. NOT squared) to decide
                    whether a new point will be appended or not.
          diff   -- [optional] a function for further check that the new and an eventual existent
                    points are really different, even after finding that they coincide (within tol)
        Output:
          id       -- the id attached to x
          existent -- flag telling if x was found, based on given tolerance

func (o *Bins) GetLimits(idxBin int) (xmin, xmax []float64)
    GetLimits returns limigs of a specific bin

func (o *Bins) Init(xmin, xmax []float64, ndiv []int)
    Init initialise Bins structure

        xmin -- [ndim] min/initial coordinates of the whole space (box/cube)
        xmax -- [ndim] max/final coordinates of the whole space (box/cube)
        ndiv -- [ndim] number of divisions for xmax-xmin

func (o *Bins) Nactive() (nactive int)
    Nactive returns the number of active bins; i.e. non-nil bins

func (o *Bins) Nentries() (nentries int)
    Nentries returns the total number of entries (in active bins)

func (o Bins) String() string
    String returns the string representation of a set of Bins

func (o *Bins) Summary() (l string)
    Summary returns the summary of this Bins' data

type BoxN struct {
	// essential
	Lo *PointN // lower point
	Hi *PointN // higher point

	// auxiliary
	ID int // an auxiliary identification number
}
    BoxN implements a box int he N-dimensional space

func NewBoxN(L ...float64) *BoxN
    NewBoxN creates a new box with given limiting coordinates

        L -- limits [4] or [6]: xmin,xmax, ymin,ymax, {zmin,zmax} optional

func (o *BoxN) GetMid() (mid []float64)
    GetMid gets the centre coordinates of box

func (o *BoxN) GetSize() (delta []float64)
    GetSize returns the size of box (deltas)

func (o BoxN) IsInside(p *PointN) bool
    IsInside tells whether a PointN is inside box or not

type Bspline struct {

	// essential
	T []float64 // array of knots: e.g: T = [t_0, t_1, ..., t_{m-1}]

	// optional
	Q [][]float64 // control points (has to call SetControl to change this)

	// Has unexported fields.
}
    Bspline holds B-spline data

        Reference:
         [1] Piegl L and Tiller W (1995) The NURBS book, Springer, 646p

func NewBspline(T []float64, p int) (o *Bspline)
    NewBspline returns a new B-spline

func (o *Bspline) CalcBasis(t float64)
    CalcBasis computes all non-zero basis functions N[i] @ t Note: use GetBasis
    to get a particular basis function value

func (o *Bspline) CalcBasisAndDerivs(t float64)
    CalcBasisAndDerivs computes all non-zero basis functions N[i] and
    corresponding first order derivatives of basis functions w.r.t t => dR[i]dt
    @ t Note: use GetBasis to get a particular basis function value

        use GetDeriv to get a particular derivative

func (o *Bspline) Elements() (spans [][]int)
    Elements returns the indices of nonzero spans

func (o *Bspline) GetBasis(i int) float64
    GetBasis returns the basis function N[i] just computed by CalcBasis or
    CalcBasisAndDerivs

func (o *Bspline) GetDeriv(i int) float64
    GetDeriv returns the derivative dN[i]dt just computed by CalcBasisAndDerivs

func (o *Bspline) NumBasis() int
    NumBasis returns the number (n) of basis functions == number of control
    points

func (o *Bspline) Point(t float64, option int) (C []float64)
    Point returns the x-y-z coordinates of a point on B-spline option = 0 : use
    CalcBasis

        1 : use RecursiveBasis

func (o *Bspline) RecursiveBasis(t float64, i int) float64
    RecursiveBasis computes one particular basis function N[i] recursively (not
    efficient)

func (o *Bspline) SetControl(Q [][]float64)
    SetControl sets B-spline control points

func (o *Bspline) SetOrder(p int)
    SetOrder sets B-spline order (p)

type Entity interface {
}
    Entity defines the geometric (or not) entity/element to be stored in the
    Octree

type Grid struct {
	// Has unexported fields.
}
    Grid implements (2D/3D) rectangular or curvilinear grid. It also holds
    metrics data related to curvilinear coordinates.

        Notation:
           m,n,p -- indices used for grid points
           i,j,k -- indices used for dimension (x,y,z)
           Ex: the covariant vector @ (m,n,p) is: o.mtr[p][n][m].GovG0
               the i component of this vector is: o.mtr[p][n][m].GovG0[i]

        NOTE: (1) the deep3 structure mtr holds data with the outer index corresponding to z
                  i.e. o.mtr[idxZ][idxY][idxX]
              (2) the reference coordinates of generated rectangular grids are assumed to be
                  -1 ≤ u ≤ +1

func (o *Grid) Boundary(tag int) []int
    Boundary returns a list of indices of nodes on edge or face of boundary

        NOTE: will return empty list if tag is not available
              see EdgeGivenTag() and FaceGivenTag() functions

func (o *Grid) ContraBasis(m, n, p, k int) la.Vector
    ContraBasis returns the [k] contravariant basis g^{k} = d{u[k]}/d{x} [@
    point m,n,p]

func (o *Grid) ContraMatrix(m, n, p int) *la.Matrix
    ContraMatrix returns contravariant metrics g^ij = g^i ⋅ g^j [@ point m,n,p]

func (o *Grid) CovarBasis(m, n, p, k int) la.Vector
    CovarBasis returns the [k] covariant basis g_{k} = d{x}/d{u[k]} [@ point
    m,n,p]

func (o *Grid) CovarMatrix(m, n, p int) *la.Matrix
    CovarMatrix returns the covariant metrics g_ij = g_i ⋅ g_j [@ point m,n,p]

func (o *Grid) DetCovarMatrix(m, n, p int) float64
    DetCovarMatrix returns the determinant of covariant g matrix =
    det(CovariantMatrix) [@ point m,n,p]

func (o *Grid) Edge(iEdge int) []int
    Edge returns the ids of points on edges: [edge0, edge1, edge2, edge3]

               3
         +-----------+    Considering the x-y axes below, the order of indices follows:
         |           |
         |           |       y         0          1          2          3
        0|           |1      ↑    {xmin_edge, xmax_edge, ymin_edge, ymax_edge}
         |           |       │
         |           |       +——→ x
         +-----------+
               2

func (o *Grid) EdgeGivenTag(tag int) []int
    EdgeGivenTag returns a list of nodes marked with given tag

                21
           +-----------+     Considering the x-y axes below, the order of tags follows:
           |           |
           |           |        y      0   1   2   3
         10|           |11      ↑    {10, 11, 20, 21}
           |           |        │
           |           |        +——→ x
           +-----------+
                20

        NOTE: will return empty list if tag is not available

func (o *Grid) Face(iFace int) []int
    Face returns the ids of points on faces: [face0, face1, face2, face3, face4,
    face5]

                  +----------------+   Considering the x-y-z axes below,
                ,'|              ,'|   the order of indices follows:
              ,'  |  ___       ,'  |
            ,'    |,'5,' [0] ,'    |      z         0: xmin_face
          ,'      |~~~     ,'  ,   |      ↑         1: xmax_face
        +'===============+'  ,'|   |      │         2: ymin_face
        |   ,'|   |      |   |3|   |      +——→y     3: ymax_face
        |   |2|   |      |   |,'   |    ,'          4: zmin_face
        |   |,'   +- - - | +- - - -+   x            5: zmax_face
        |   '   ,'       |       ,'
        |     ,' [1]  ___|     ,'
        |   ,'      ,'4,'|   ,'
        | ,'        ~~~  | ,'
        +----------------+'

func (o *Grid) FaceGivenTag(tag int) []int
    FaceGivenTag returns a list of nodes marked with given tag

                    +----------------+   Considering the x-y-z axes below,
                  ,'|              ,'|   the order of tags follows:
                ,'  |  ___       ,'  |
              ,'    |,301' 100 ,'    |      z         100: xmin_face
            ,'      |~~~     ,'  ,'  |      ↑         101: xmax_face
          +'===============+'  ,' |  |      │         200: ymin_face
          |   ,'|   |      |   |201  |      +——→y     201: ymax_face
          |  |200   |      |   |,'   |    ,'          300: zmin_face
          |  | ,'   +- - - | +- - - -+   x            301: zmax_face
          |  ,'   ,'       |       ,'
          |     ,'101   ___|     ,'
          |   ,'      ,300'|   ,'
          | ,'        ~~~  | ,'
          +----------------+'

        NOTE: will return empty list if tag is not available

func (o *Grid) GammaS(m, n, p, k, i, j int) float64
    GammaS returns the [k][i][j] Christoffel coefficients of second kind [@
    point m,n,p]

func (o *Grid) IndexItoMNP(I int) (m, n, p int)
    IndexItoMNP converts node index I into triplet indices (m,n,p)

        2D:   I = m + n⋅n0
              m = I % n0
              n = I / n0

        3D:   I = m + n⋅n0 + p⋅n0⋅n1
              p = I / (n0⋅n1)
              t = I % (n0⋅n1)  (projection @ z=0)
              m = t % n0
              n = t / n0

func (o *Grid) IndexMNPtoI(m, n, p int) (I int)
    IndexMNPtoI converts node triplet indices (m,n,p) into node index I

        2D:   I = m + n⋅n0
              m = I % n0
              n = I / n0

        3D:   I = m + n⋅n0 + p⋅n0⋅n1
              p = I / (n0⋅n1)
              t = I % (n0⋅n1)  (projection @ z=0)
              m = t % n0
              n = t / n0

func (o *Grid) Lcoeff(m, n, p, k int) float64
    Lcoeff returns the [k] L-coefficients = sum(Γ_ij^k ⋅ g^ij) [@ point m,n,p]

func (o *Grid) MapMeshgrid2d(v la.Vector) (V [][]float64)
    MapMeshgrid2d maps vector V into 2D meshgrid using node indices conversion
    IndexMNPtoI()

        vv[ny][nx] -- mapped values: vv[n][m] ⇐ V[I] (see also Meshgrid2d)

func (o *Grid) MapMeshgrid3d(v la.Vector) (V [][][]float64)
    MapMeshgrid3d maps vector V into 3D meshgrid using node indices conversion
    IndexMNPtoI()

        V[nz][ny][nx] -- mapped values: V[p][n][m] ⇐ v[I] (see also Meshgrid2d)

func (o *Grid) Meshgrid2d() (X, Y [][]float64)
    Meshgrid2d extracts 2D meshgrid

        X -- x0[ny][nx]
        Y -- x1[ny][nx]

func (o *Grid) Meshgrid3d() (X, Y, Z [][][]float64)
    Meshgrid3d extracts 3D meshgrid

        X -- x0[nz][ny][nx]
        Y -- x1[nz][ny][nx]
        Z -- x2[nz][ny][nx]

func (o *Grid) Ndim() int
    Ndim returns the number of dimensions (2D or 3D)

func (o *Grid) Node(I int) (x la.Vector)
    Node returns the physical coordinates of node I. See IndexItoMNP(I) ⇒
    (m,n,p).

        x -- slice to position vector @ [p][n][m] [may be used to change values]

func (o *Grid) Npts(idim int) int
    Npts returns number of points along idim dimension

func (o *Grid) RectGenUniform(xmin, xmax []float64, npts []int)
    RectGenUniform generates uniform coordinates of a rectangular grid

        xmin -- min x-y-z values, len==ndim: [xmin, ymin, zmin]
        xmax -- max x-y-z values, len==ndim: [xmax, ymax, zmax]
        npts -- number of points along each direction len==ndim: [n0, n1, n2] (must be greater than 2)

           -1 ≤ u ≤ +1
           x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
           u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
           dx/du = (xmax - xmin) / 2

func (o *Grid) RectSet2d(X, Y []float64)
    RectSet2d sets rectangular grid with given coordinates

        -1 ≤ u ≤ +1
        x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
        u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
        dx/du = (xmax - xmin) / 2

func (o *Grid) RectSet2dU(xmin, xmax, R, S []float64)
    RectSet2dU sets rectangular grid with given reference coordinates and limits

        Input:
          xmin -- min x-y values: [xmin, ymin]
          xmax -- max x-y values: [xmax, ymax]
          R -- reference coordinates along x:  -1 ≤ r ≤ +1
          S -- reference coordinates along y:  -1 ≤ s ≤ +1

           -1 ≤ u ≤ +1
           x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
           u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
           dx/du = (xmax - xmin) / 2

func (o *Grid) RectSet3d(X, Y, Z []float64)
    RectSet3d sets rectangular grid with given coordinates

        -1 ≤ u ≤ +1
        x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
        u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
        dx/du = (xmax - xmin) / 2

func (o *Grid) RectSet3dU(xmin, xmax, R, S, T []float64)
    RectSet3dU sets rectangular grid with given reference coordinates and limits

        Input:
          xmin -- min x-y-z values: [xmin, ymin, zmin]
          xmax -- max x-y-z values: [xmax, ymax, zmax]
          R -- reference coordinates along x:  -1 ≤ r ≤ +1
          S -- reference coordinates along y:  -1 ≤ s ≤ +1
          T -- reference coordinates along z:  -1 ≤ t ≤ +1

           -1 ≤ u ≤ +1
           x(u) = xmin + (xmax - xmin) ⋅ (1 + u) / 2
           u(x) = -1 + 2⋅(x - xmin) / (xmax - xmin)
           dx/du = (xmax - xmin) / 2

func (o *Grid) SetNurbsSolid(nrb *Nurbs, R, S, T []float64)
    SetNurbsSolid sets grid with NURBS solid

        nrb -- NURBS solid
        R   -- [n0] reference coordinates along r-direction  -1 ≤ r ≤ +1
        S   -- [n1] reference coordinates along s-direction  -1 ≤ s ≤ +1
        T   -- [n2] reference coordinates along s-direction  -1 ≤ s ≤ +1

func (o *Grid) SetNurbsSurf2d(nrb *Nurbs, R, S []float64)
    SetNurbsSurf2d sets grid with NURBS surface in 2D (flat surface)

        nrb -- NURBS surface in 2D (flat)
        R   -- [n0] reference coordinates along r-direction  -1 ≤ r ≤ +1
        S   -- [n1] reference coordinates along s-direction  -1 ≤ s ≤ +1

func (o *Grid) SetTransfinite2d(trf *Transfinite, R, S []float64)
    SetTransfinite2d sets grid from 2D transfinite mapping

        trf -- 2D transfinite structure
        R   -- [n0] reference coordinates along r-direction  -1 ≤ r ≤ +1
        S   -- [n1] reference coordinates along s-direction  -1 ≤ s ≤ +1

func (o *Grid) SetTransfinite3d(trf *Transfinite, R, S, T []float64)
    SetTransfinite3d sets grid from 3D transfinite mapping

        trf -- 2D transfinite structure
        R   -- [n0] reference coordinates along r-direction  -1 ≤ r ≤ +1
        S   -- [n1] reference coordinates along s-direction  -1 ≤ s ≤ +1
        T   -- [n2] reference coordinates along s-direction  -1 ≤ t ≤ +1

func (o *Grid) Size() int
    Size returns total number of points

func (o *Grid) U(m, n, p int) la.Vector
    U returns the reference coordinates at point m,n,p

func (o *Grid) Umax(idim int) float64
    Umax returns the maximum reference coordinate at dimension idim

func (o *Grid) Umin(idim int) float64
    Umin returns the minimum reference coordinate at dimension idim

func (o *Grid) UnitNormal(N la.Vector, tag, I int)
    UnitNormal computes the unit normal vector at an edge or face defined by
    "tag" and at a node specified by index "I".

        Input:
          tag -- tag of edge or face; see EdgeGivenTag() or FaceGivenTag()
          I   -- node index; see IndexMNPtoI() [must be a node on boundary; otherwise, panic]
        Output:
          N -- unit normal vector

        NOTE: this function does not check whether point I is on the selected edge or not

func (o *Grid) X(m, n, p int) la.Vector
    X returns the physical coordinates at point m,n,p

func (o *Grid) Xlen(idim int) float64
    Xlen returns the lengths along each direction (whole box) == Xmax(idim) -
    Xmin(idim)

func (o *Grid) Xmax(idim int) float64
    Xmax returns the maximum physical coordinate at dimension idim

func (o *Grid) Xmin(idim int) float64
    Xmin returns the minimum physical coordinate at dimension idim

type Metrics struct {
	U           la.Vector     // reference coordinates {r,s,t}
	X           la.Vector     // physical coordinates {x,y,z}
	CovG0       la.Vector     // covariant basis g_0 = d{x}/dr
	CovG1       la.Vector     // covariant basis g_1 = d{x}/ds
	CovG2       la.Vector     // covariant basis g_2 = d{x}/dt
	CntG0       la.Vector     // contravariant basis g_0 = dr/d{x} (gradients)
	CntG1       la.Vector     // contravariant basis g_1 = ds/d{x} (gradients)
	CntG2       la.Vector     // contravariant basis g_2 = dt/d{x} (gradients)
	CovGmat     *la.Matrix    // covariant metrics g_ij = g_i ⋅ g_j
	CntGmat     *la.Matrix    // contravariant metrics g^ij = g^i ⋅ g^j
	DetCovGmat  float64       // determinant of covariant g matrix = det(CovGmat)
	Homogeneous bool          // homogeneous grid => nil second order derivatives and Christoffel symbols
	GammaS      [][][]float64 // [k][i][j] Christoffel coefficients of second kind (non-homogeneous)
	L           []float64     // [3] L-coefficients = sum(Γ_ij^k ⋅ g^ij) (non-homogeneous)
}
    Metrics holds data related to a position in a space represented by
    curvilinear coordinates

func NewMetrics2d(u, x, dxdr, dxds, ddxdrr, ddxdss, ddxdrs la.Vector) (o *Metrics)
    NewMetrics2d allocate new 2D metrics structure

        NOTE: the second order derivatives (from ddxdrr) may be nil => homogeneous grid

func NewMetrics3d(u, x, dxdr, dxds, dxdt, ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst la.Vector) (o *Metrics)
    NewMetrics3d allocate new 3D metrics structure

        NOTE: the second order derivatives (from ddxdrr) may be nil => homogeneous grid

type Nurbs struct {
	Q [][][][]float64 // Qw: weighted control points and weights [n[0]][n[1]][n[2]][4] (Piegl p120)

	// Has unexported fields.
}
    Nurbs holds NURBS data

        NOTE: (1) Control points must be set after a call to Init
              (2) Either SetControl must be called or the Q array must be directly specified

        Reference:
         [1] Piegl L and Tiller W (1995) The NURBS book, Springer, 646p

func NewNurbs(gnd int, ords []int, knots [][]float64) (o *Nurbs)
    NewNurbs returns a new Nurbs object

func (o *Nurbs) CalcBasis(u []float64)
    CalcBasis computes all non-zero basis functions R[i][j][k] @ u. Note: use
    GetBasisI or GetBasisL to get a particular basis function value

func (o *Nurbs) CalcBasisAndDerivs(u []float64)
    CalcBasisAndDerivs computes all non-zero basis functions R[i][j][k] and
    corresponding first order derivatives of basis functions w.r.t u =>
    dRdu[i][j][k] @ u Note: use GetBasisI or GetBasisL to get a particular basis
    function value

        use GetDerivI or GetDerivL to get a particular derivative

func (o *Nurbs) CloneCtrlsAlongCurve(iAlong, jAt int) (Qnew [][][][]float64)
    CloneCtrlsAlongCurve returns a copy of control points @ 2D boundary

func (o *Nurbs) CloneCtrlsAlongSurface(iAlong, jAlong, kat int) (Qnew [][][][]float64)
    CloneCtrlsAlongSurface returns a copy of control points @ 3D boundary

func (o *Nurbs) ElemBryLocalInds() (I [][]int)
    ElemBryLocalInds returns the local (element) indices of control points @
    boundaries (if element would have all surfaces @ boundaries)

func (o *Nurbs) Elements() (spans [][]int)
    Elements returns the indices of nonzero spans

func (o *Nurbs) ExtractSurfaces() (surfs []*Nurbs)
    ExtractSurfaces returns a new NURBS representing a boundary of this NURBS

func (o *Nurbs) GetBasisI(I []int) float64
    GetBasisI returns the basis function R[i][j][k] just computed by CalcBasis
    or CalcBasisAndDerivs Note: I = [i,j,k]

func (o *Nurbs) GetBasisL(l int) float64
    GetBasisL returns the basis function R[i][j][k] just computed by CalcBasis
    or CalcBasisAndDerivs Note: l = i + j * n0 + k * n0 * n1

func (o *Nurbs) GetDerivI(dRdu []float64, I []int)
    GetDerivI returns the derivative of basis function dR[i][j][k]du just
    computed by CalcBasisAndDerivs Note: I = [i,j,k]

func (o *Nurbs) GetDerivL(dRdu []float64, l int)
    GetDerivL returns the derivative of basis function dR[i][j][k]du just
    computed by CalcBasisAndDerivs Note: l = i + j * n0 + k * n0 * n1

func (o *Nurbs) GetElemNumBasis() (npts int)
    GetElemNumBasis returns the number of control points == basis functions
    needed for one element

        npts := Π_i (p[i] + 1)

func (o *Nurbs) GetLimitsQ() (xmin, xmax []float64)
    GetLimitsQ computes the limits of all coordinates of control points in NURBS

func (o *Nurbs) GetQ(i, j, k int) (x []float64)
    GetQ gets a control point x[i,j,k] (size==4)

func (o *Nurbs) GetQl(l int) (x []float64)
    GetQl gets a control point x[l] (size==4)

func (o *Nurbs) GetU(dir int) []float64
    GetU returns the knots along direction dir

func (o *Nurbs) Gnd() int
    Gnd returns the geometry dimension

func (o *Nurbs) IndBasis(span []int) (L []int)
    IndBasis returns the indices of basis functions == local indices of control
    points

func (o *Nurbs) IndsAlongCurve(iAlong, iSpan0, jAt int) (L []int)
    IndsAlongCurve returns the control points indices along curve

func (o *Nurbs) IndsAlongSurface(iAlong, jAlong, iSpan0, jSpan0, kat int) (L []int)
    IndsAlongSurface return the control points indices along surface

func (o *Nurbs) Krefine(X [][]float64) (O *Nurbs)
    Krefine returns a new Nurbs with knots refined Note: X[gnd][numNewKnots]

func (o *Nurbs) KrefineN(ndiv int, hughesEtAlPaper bool) *Nurbs
    KrefineN return a new Nurbs with each span divided into ndiv parts = [2, 3,
    ...]

func (o *Nurbs) NonZeroSpans(dir int) [][]int
    NonZeroSpans returns the 'elements' along direction dir

func (o *Nurbs) NumBasis(dir int) int
    NumBasis returns the number of basis (controls) along direction dir

func (o *Nurbs) Ord(dir int) int
    Ord returns the order along direction dir

func (o *Nurbs) Point(C, u []float64, ndim int)
    Point returns the x-y-z coordinates of a point on curve/surface/solid

        Input:
          u    -- [gnd] knot values
          ndim -- the dimension of the point. E.g. allows drawing curves in 3D
        Output:
          C -- [ndim] point coordinates
        NOTE: Algorithm A4.1 (p124) of [1]

func (o *Nurbs) PointAndDerivs(x, dxdr, dxds, dxdt,
	ddxdrr, ddxdss, ddxdtt, ddxdrs, ddxdrt, ddxdst, u la.Vector, ndim int)
    PointAndDerivs computes position and the first and second order derivatives
    Using Algorithms A3.2(p93), A3.6(p111), A4.2(p127), and A4.4(p137)

        Input:
          u    -- knot values {r,s,t} [gnd]
          ndim -- the dimension of the point. E.g. allows drawing curves in 3D [ndim=3 if gnd=3]
        Output:
          x      -- position {x,y,z} (the same as the C varible in [1])
          dxdr   -- ∂{x}/∂r
          dxds   -- ∂{x}/∂s    [may be nil] (solid and surfaces)
          dxdt   -- ∂{x}/∂t    [may be nil] (solid)
          ddxdrr -- ∂²{x}/∂r²  [optional]
          ddxdss -- ∂²{x}/∂s²  [optional] [may be nil] (solid and surfaces)
          ddxdtt -- ∂²{x}/∂t²  [optional] [may be nil] (solid)
          ddxdrs -- ∂²{x}/∂r∂s [optional] [may be nil] (solid and surfaces)
          ddxdrt -- ∂²{x}/∂r∂t [optional] [may be nil] (solid)
          ddxdst -- ∂²{x}/∂s∂t [optional] [may be nil] (solid)
        NOTE: the second order derivatives will be ignored if ddxdrr == nil; otherwise, all 2nd derivs will be computed

func (o *Nurbs) PointAndFirstDerivs(dCdu *la.Matrix, C, u []float64, ndim int)
    PointAndFirstDerivs returns the point and first order derivatives with
    respect to the knot values u of the x-y-z coordinates of a point on
    curve/surface/solid

        Input:
          u    -- [gnd] knot values
          ndim -- the dimension of the point. E.g. allows drawing curves in 3D
        Output:
          dCdu -- [ndim][gnd] derivatives dC_i/du_j
          C    -- [ndim] point coordinates

func (o *Nurbs) RecursiveBasis(u []float64, l int) (res float64)
    RecursiveBasis implements basis functions by means of summing all terms in
    Bernstein polynomial using recursive Cox-DeBoor formula (very not efficient)
    Note: l = i + j * n0 + k * n0 * n1

func (o *Nurbs) SetControl(verts [][]float64, ctrls []int)
    SetControl sets control points from list of global vertices

func (o *Nurbs) SetQ(i, j, k int, x []float64)
    SetQ sets a control point x[i,j,k] (size==4)

func (o *Nurbs) SetQl(l int, x []float64)
    SetQl sets a control point x[l] (size==4)

func (o *Nurbs) U(dir, idx int) float64
    U returns the value of a knot along direction dir

func (o *Nurbs) Udelta(dir int) (Δu float64)
    Udelta returns the difference between max knot value and min knot value.
    Returns umax - umin along direnction "dir"

func (o *Nurbs) UfromR(dir int, r float64) (u float64)
    UfromR returns knot value from reference value -1 ≤ r ≤ +1

        Input:
          dir -- direction (dimension)
          r   -- "reference" (e.g. spectral) coordinates -1 ≤ r ≤ +1
        Output:
          u -- knot value

type NurbsExchangeData struct {
	ID    int         `json:"i"` // id of Nurbs
	Gnd   int         `json:"g"` // 1: curve, 2:surface, 3:volume (geometry dimension)
	Ords  []int       `json:"o"` // order along each x-y-z direction [gnd]
	Knots [][]float64 `json:"k"` // knots along each x-y-z direction [gnd][m]
	Ctrls []int       `json:"c"` // global ids of control points
}
    NurbsExchangeData holds all data required to exchange NURBS; e.g. read/save
    files

type NurbsExchangeDataSet []*NurbsExchangeData
    NurbsExchangeDataSet defines a set of nurbs exchange data

type NurbsPatch struct {

	// input/output data
	ControlPoints []*PointExchangeData `json:"points"` // input/output control points
	ExchangeData  NurbsExchangeDataSet `json:"patch"`  // input/output nurbs data

	// Nurbs structures
	Entities []*Nurbs `json:"-"` // pointers to NURBS

	// auxiliary
	Bins Bins `json:"-"` // auxiliary structure to locate points
}
    NurbsPatch defines a patch of many NURBS'

func NewNurbsPatch(binsNdiv int, tolerance float64, entities ...*Nurbs) (o *NurbsPatch)
    NewNurbsPatch returns new patch of NURBS

        tolerance -- tolerance to assume that two control points are the same

func NewNurbsPatchFromFile(filename string, binsNdiv int, tolerance float64) (o *NurbsPatch)
    NewNurbsPatchFromFile allocates a NurbsPatch with data from file

func (o NurbsPatch) LimitsAndNdim() (xmin, xmax []float64, ndim int)
    LimitsAndNdim computes the limits of the patch and max dimension by looping
    over all Entities

func (o *NurbsPatch) ResetFromEntities(binsNdiv int, tolerance float64)
    ResetFromEntities will reset all exchange data with information from
    Entities slice

func (o *NurbsPatch) ResetFromExchangeData(binsNdiv int, tolerance float64)
    ResetFromExchangeData will reset all Entities with information from
    ExchangeData (and ControlPoints)

func (o NurbsPatch) Write(dirout, fnkey string)
    Write writes ExchangeData to json file

type Octree struct {

	// constants
	DIM  uint32 // dimension
	PMAX uint32 // roughly how many levels fit in 32 bits
	QO   uint32 // 4 for quadtree, 8 for octree
	QL   uint32 // offset constant to leftmost daughter

	// Has unexported fields.
}
    Octree implements a Quad-Tree or an Oct-Tree to assist in fast-searching
    elements (entities) in the 2D or 3D space

func NewOctree(L ...float64) (o *Octree)
    NewOctree creates a new Octree

        L -- limits [4] or [6]: xmin,xmax, ymin,ymax, {zmin,zmax} optional

type Point struct {
	X, Y, Z float64
}
    Point holds the Cartesian coordinates of a point in 3D space

func (o *Point) NewCopy() *Point
    NewCopy creates a new copy of Point

func (o *Point) NewDisp(dx, dy, dz float64) *Point
    NewDisp creates a new copy of Point displaced by dx, dy, dz

func (o *Point) String() string
    String outputs Point

type PointExchangeData struct {
	ID  int       `json:"i"` // id
	Tag int       `json:"t"` // tag
	X   []float64 `json:"x"` // coordinates (size==4)
}
    PointExchangeData holds data for exchanging control points

type PointN struct {

	// esssential
	X []float64 // coordinates

	// optional
	ID int // some identification number
}
    PointN implements a point in N-dim space

func NewPointN(X ...float64) (o *PointN)
    NewPointN creats a new PointN with given coordinates; can be any number

func NewPointNdim(ndim uint32) (o *PointN)
    NewPointNdim creates a new PointN with given dimension (ndim)

func (o PointN) AlmostTheSameX(p *PointN, tol float64) bool
    AlmostTheSameX returns true if the X slices of two points have almost the
    same values, for given tolerance

func (o PointN) ExactlyTheSameX(p *PointN) bool
    ExactlyTheSameX returns true if the X slices of two points have exactly the
    same values

func (o PointN) GetCloneX() (p *PointN)
    GetCloneX returns a new point with X cloned, but not the other data

type Segment struct {
	A, B *Point
}
    Segment represents a directed segment from A to B

func NewSegment(a, b *Point) *Segment
    NewSegment creates a new segment from a to b

func (o *Segment) Len() float64
    Len computes the length of Segment == Euclidean norm

func (o *Segment) New(m float64) *Segment
    New creates a new Segment scaled by m and starting from A

func (o *Segment) String() string
    String outputs Segment

func (o *Segment) Vector(m float64) []float64
    Vector returns the vector representing Segment from A to B (scaled by m)

type Transfinite struct {
	// Has unexported fields.
}
    Transfinite maps a reference square [-1,+1]×[-1,+1] into a curve-bounded
    quadrilateral

                                                   B[3](r(x,y)) _,'\
                    B[3](r)                                  _,'    \ B[1](s(x,y))
                   ┌───────┐                              _,'        \
                   │       │                             \            \
            B[0](s)│       │B[1](s)     ⇒                 \         _,'
         s         │       │                  B[0](s(x,y)) \     _,'
         │         └───────┘               y                \ _,'  B[2](r(x,y))
         └──r       B[2](r)                │                 '
                                           └──x

                                          +----------------+
                                        ,'|              ,'|
             t or z                   ,'  |  ___       ,'  |     B[0](s,t)
                ↑                   ,'    |,'5,'  [0],'    |     B[1](s,t)
                |                 ,'      |~~~     ,'      |     B[2](r,t)
                |               +'===============+'  ,'|   |     B[3](r,t)
                |               |   ,'|   |      |   |3|   |     B[4](r,s)
                |     s or y    |   |2|   |      |   |,'   |     B[5](r,s)
                +-------->      |   |,'   +- - - | +- - - -+
              ,'                |       ,'       |       ,'
            ,'                  |     ,' [1]  ___|     ,'
        r or x                  |   ,'      ,'4,'|   ,'
                                | ,'        ~~~  | ,'
                                +----------------+'

        NOTE: the "reference" coordinates {r,s,t} ϵ [-1,+1]×[-1,+1]×[-1,+1] are also symbolised as "u"

func NewTransfinite2d(B, Bd, Bdd []fun.Vs) (o *Transfinite)
    NewTransfinite2d allocates a new structure

        Input:
          B   -- [4] boundary functions
          Bd  -- [4] 1st derivative of boundary functions
          Bdd -- [4 or nil] 2nd derivative of boundary functions [may be nil]

func NewTransfinite3d(B []fun.Vss, Bd []fun.Vvss, Bdd []fun.Vvvss) (o *Transfinite)
    NewTransfinite3d allocates a new structure

        Input:
          B   -- [6] boundary functions
          Bd  -- [6] 1st derivative of boundary functions
          Bdd -- [6 or nil] 2nd derivative of boundary functions [may be nil]

func (o *Transfinite) Point(x, u la.Vector)
    Point computes "real" position x(r,s,t)

        Input:
          u -- the "reference" coordinates {r,s,t} ϵ [-1,+1]×[-1,+1]×[-1,+1]
        Output:
          x -- the "real" coordinates {x,y,z}

func (o *Transfinite) PointAndDerivs(x, dxDr, dxDs, dxDt,
	ddxDrr, ddxDss, ddxDtt, ddxDrs, ddxDrt, ddxDst, u la.Vector)
    PointAndDerivs computes position and the first and second order derivatives

        Input:
          u -- reference coordinates {r,s,t}
        Output:
          x      -- position {x,y,z}
          dxDr   -- ∂{x}/∂r
          dxDs   -- ∂{x}/∂s
          dxDt   -- ∂{x}/∂t
          ddxDrr -- ∂²{x}/∂r²  [may be nil]
          ddxDss -- ∂²{x}/∂s²  [may be nil]
          ddxDtt -- ∂²{x}/∂t²  [may be nil]
          ddxDrs -- ∂²{x}/∂r∂s [may be nil]
          ddxDrt -- ∂²{x}/∂r∂t [may be nil]
          ddxDst -- ∂²{x}/∂s∂t [may be nil]

```
