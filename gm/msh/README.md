# Gosl. gm/msh. Mesh generation and Delaunay triangulation

## Definitions

![](data/meshelems1.png)

### Quadrilaterals and Hexahedra

```
  2D:
                Vertices                              Edges
   y
   |        3      6      2                             2
   +--x      @-----@-----@                        +-----------+
             |           |                        |           |
             |           |                        |           |
           7 @           @ 5                     3|           |1
             |           |                        |           |
             |           |                        |           |
             @-----@-----@                        +-----------+
            0      4      1                             0

  3D:
                  Vertices                            Faces
    z
    |           4        15        7
   ,+--y         @-------@--------@                 +----------------+
 x'            ,'|              ,'|               ,'|              ,'|
          12 @'  |         14 ,'  |             ,'  |  ___       ,'  |
           ,'    |16        ,@    |19         ,'    |,'5,'  [0],'    |
     5   ,'      @      6 ,'      @         ,'      |~~~     ,'      |
       @'=======@=======@'        |       +'===============+'  ,'|   |
       |      13 |      |         |       |   ,'|   |      |   |3|   |
       |         |      |  11     |       |   |2|   |      |   |,'   |
    17 |       0 @- - - | @- - - -@       |   |,'   +- - - | +- - - -+
       @       ,'       @       ,' 3      |       ,'       |       ,'
       |   8 @'      18 |     ,'          |     ,' [1]  ___|     ,'
       |   ,'           |   ,@ 10         |   ,'      ,'4,'|   ,'
       | ,'             | ,'              | ,'        ~~~  | ,'
       @-------@--------@'                +----------------+'
     1         9         2

  3D:
                  Seams (aka Edges)
    z
    |                    7
   ,+--y         +----------------+
 x'            ,'|              ,'|
           4 ,'  |8          6,'  |
           ,'    |          ,'    |
         ,'      |   5    ,'      |11
       +'===============+'        |
       |         |      |         |
       |         |      |  3      |
       |         .- - - | -  - - -+
      9|       ,'       |       ,'
       |    0,'         |10   ,'
       |   ,'           |   ,' 2
       | ,'             | ,'
       +----------------+'
               1
```

### Triangles and Tetrahedra

```
  2D:
             Nodes                      Edges

   y           2
   |           @                          @
   +--x       / \                        / \
           5 /   \ 4                    /   \
            @     @                  2 /     \ 1
           /       \                  /       \
          /         \                /         \
         @-----@-----@              @-----------@
        0      3      1                   0


  3D:              Nodes                                     Faces
    z
    |                |                                          |
   ,+--y             3                                          +
 x'                 /|`.                                       /|`.
                    ||  `,                                     ||  `,
                   / |    ',                                  / |    ',
                   | |      \                                 | |      \
                  /  |       `.                              /  |       `.
                  |  |         `,                            |  |         `,
                 /   7            9                         /   |           ',
                 |   |             \                        |   |             \
                /    |              `.                     /    |     0        `.
                |    |                ',                   |    |                ',
               8     |                  \                 /  1  |                  \
               |     0 ,,_               `.               |     + ,,_               `.
              |     /     ``'-., 6         `.            |     /     ``'-.,,.         `.
              |    /               `''-.,,_  ',          |    /              ``''-.,,_  ',
             |    /                        ``'2 ,,      |    /      \3\               ``'+ ,,
             |   '                       ,.-``          |  ,'                       ,.-``
            |   4                   _,-'`              |  /                    _,-'`
            ' /                 ,.'`                   ' /      2          ,.'`
           | /             _ 5 `                      | /             _.''`
           '/          ,-'`                           '/          ,-'`
          |/      ,.-``                              |/      ,.-``
          /  _,-``                                   /  _,-``
         1 '`                                       +.'`

    3D:
         Seams (aka edges)

               |
               +
              /|`.
              ||  `,
             / |    ',
             | |      \
            /  |       `.
            |  |         `,5
           /   |3          ',
           |   |             \
          /    |              `.
        4 |    |                ',
         /     |                  \
         |     +.,,_               `.
        |     /     ``'-.,,.         `.
        |    /          2   ``''-.,,_  ',
       |    /                        ``'+ ,,
       |  ,'                       ,.-``
      |  / 0                  _,-'`
      ' /                 ,.'`
     | /             _.''`
     '/          ,-'`   1
    |/      ,.-``
    /  _,-``
   +.'`
```

## API

**go doc**

```
package msh // import "gosl/gm/msh"

Package msh defines mesh data structures and implements interpolation
functions for finite element analyses (FEA)

CONSTANTS

const (
	KindLin    = 0 // "lin" cell kind
	KindTri    = 1 // "tri" cell kind
	KindQua    = 2 // "qua" cell kind
	KindTet    = 3 // "tet" cell kind
	KindHex    = 4 // "hex" cell kind
	KindNumMax = 5 // max number of kinds
)
    cell kinds

const (
	TypeLin2   = 0  // Lin2 cell type index
	TypeLin3   = 1  // Lin3 cell type index
	TypeLin4   = 2  // Lin4 cell type index
	TypeLin5   = 3  // Lin5 cell type index
	TypeTri3   = 4  // Tri3 cell type index
	TypeTri6   = 5  // Tri6 cell type index
	TypeTri10  = 6  // Tri10 cell type index
	TypeTri15  = 7  // Tri15 cell type index
	TypeQua4   = 8  // Qua4 cell type index
	TypeQua8   = 9  // Qua8 cell type index
	TypeQua9   = 10 // Qua9 cell type index
	TypeQua12  = 11 // Qua12 cell type index
	TypeQua16  = 12 // Qua16 cell type index
	TypeQua17  = 13 // Qua17 cell type index
	TypeTet4   = 14 // Tet4 cell type index
	TypeTet10  = 15 // Tet10 cell type index
	TypeHex8   = 16 // Hex8 cell type index
	TypeHex20  = 17 // Hex20 cell type index
	TypeNumMax = 18 // max number of types
)
    cell types


VARIABLES

var (
	// IntPoints holds integration points for all kinds of cells: lin,qua,hex,tri,tet
	// It maps [cellKind] => [options][npts][4] where 4 means r,s,t,w
	IntPoints map[int]map[string][][]float64

	// DefaultIntPoints holds the default integration points for all cell types
	// It maps [cellTypeIndex] => [npts][4] where 4 means r,s,t,w
	// NOTE: the highest number of integration points is selected,
	//       thus the default number may not be optimal.
	DefaultIntPoints [][][]float64
)
var (
	// Functions holds functions to compute shape functions and derivatives [TypeNumMax]
	Functions []ShapeFunction

	// TypeKeyToIndex converts type key (e.g. "lin2") to index (e.g. TypeLin2)
	TypeKeyToIndex map[string]int

	// TypeIndexToKey converts type index (e.g. TypeLin2) to key (e.g. "lin2")
	TypeIndexToKey []string

	// TypeIndexToKind converts type index (e.g. TypeLin2) to cell kind (e.g. KindLin)
	TypeIndexToKind []int

	// NumVerts holds the number of vertices on shape [TypeNumMax]
	NumVerts []int

	// GeomNdim holds the geometry number of space dimensions [TypeNumMax]
	GeomNdim []int

	// EdgeLocalVerts holds the local indices of vertices on edges of shape [TypeNumMax][nedges][nverts]
	EdgeLocalVerts [][][]int

	// FaceLocalVerts holds the local indices of vertices on faces of shape [TypeNumMax][nfaces][nverts]
	FaceLocalVerts [][][]int

	//EdgeLocalVertsD holds the local indices (for drawing) of vertices on edges of shape [TypeNumMax][nedges][nverts]
	EdgeLocalVertsD [][][]int

	// FaceLocalVertsD holds the local indices (for drawing) of vertices on faces of shape [TypeNumMax][nfaces][nverts]
	FaceLocalVertsD [][][][]int

	// NatCoords holds the natural coordinates of vertices on shape [TypeNumMax][nverts][gndim]
	NatCoords [][][]float64
)

FUNCTIONS

func FuncHex20(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncHex20 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of hex20 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

                   4_______15_______7
                 ,'|              ,'|
              12'  |            ,'  |
             ,'    16         ,14   |
           ,'      |        ,'      19
         5'=====13========6'        |
         |         |      |         |
         |         |      |         |
         |         0_____ | _11_____3
        17       ,'       |       ,'
         |     8'        18     ,'
         |   ,'           |   ,10
         | ,'             | ,'
         1_______9________2'

func FuncHex8(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncHex8 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of hex8 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

                  4________________7
                ,'|              ,'|
              ,'  |            ,'  |
            ,'    |          ,'    |
          ,'      |        ,'      |
        5'===============6'        |
        |         |      |         |
        |         |      |         |
        |         0_____ | ________3
        |       ,'       |       ,'
        |     ,'         |     ,'
        |   ,'           |   ,'
        | ,'             | ,'
        1________________2'

func FuncLin2(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncLin2 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of lin2 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        -1     0    +1
         0-----------1-->r

func FuncLin3(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncLin3 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of lin3 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        -1     0    +1
         0-----2-----1-->r

func FuncLin4(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncLin4 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of lin4 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        -1                  +1
         @------@-----@------@  --> r
         0      2     3      1

func FuncLin5(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncLin5 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of lin5 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

         @-----@-----@-----@-----@-> r
         0     3     2     4     1
         |           |           |
        r=-1  -1/2   r=0  1/2   r=+1

func FuncQua12(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncQua12 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of qua12 (serendipity) elements at {r,s,t} natural
    coordinates. The derivatives are calculated only if derivs==true.

         3      10       6        2
           @-----@-------@------@
           |               (1,1)|
           |       s ^          |
         7 @         |          @ 9
           |         |          |
           |         +----> r   |
           |       (0,0)        |
        11 @                    @ 5
           |                    |
           |(-1,-1)             |
           @-----@-------@------@
         0       4       8        1

func FuncQua16(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncQua16 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of qua16 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

         3      10       6        2
           @-----@-------@------@
           |               (1,1)|
           |       s ^          |
         7 @   15@   |    @14   @ 9
           |         |          |
           |         +----> r   |
           |       (0,0)        |
        11 @   12@       @13    @ 5
           |                    |
           |(-1,-1)             |
           @-----@-------@------@
         0       4       8        1

func FuncQua17(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncQua17 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of qua17 (serendipity) elements at {r,s,t} natural
    coordinates. The derivatives are calculated only if derivs==true.

            3      14    10     6     2
              @-----@-----@-----@-----@
              |                  (1,1)|
              |                       |
        	7 @                       @ 13
              |         s ^           |
              |           |           |
           11 @           |16         @ 9
              |           @----> r    |
              |         (0,0)         |
           15 @                       @ 5
              |                       |
              |(-1,-1)                |
              @-----@-----@-----@-----@
            0       4     8    12       1

func FuncQua4(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncQua4 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of qua4 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        3-----------2
        |     s     |
        |     |     |
        |     +--r  |
        |           |
        |           |
        0-----------1

func FuncQua8(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncQua8 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of qua8 (serendipity) elements at {r,s,t} natural
    coordinates. The derivatives are calculated only if derivs==true.

        3-----6-----2
        |     s     |
        |     |     |
        7     +--r  5
        |           |
        |           |
        0-----4-----1

func FuncQua9(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncQua9 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of qua9 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        3-----6-----2
        |     s     |
        |     |     |
        7     8--r  5
        |           |
        |           |
        0-----4-----1

func FuncTet10(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncTet10 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of tet10 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

                      t
                      |
                      3
                     /|`.
                     ||  `,
                    / |    ',
                    | |      \
                   /  |       `.
                   |  |         `,
                  /   7            9
                  |   |             \
                 /    |              `.
                 |    |                ',
                8     |                  \
                |     0 ,,_               `.
               |     /     ``'-., 6         `.
               |    /               `''-.,,_  ',
              |    /                        ``'2 ,,s
              |   '                       ,.-``
             |   4                   _,-'`
             ' /                 ,.'`
            | /             _ 5 `
            '/          ,-'`
           |/      ,.-``
           /  _,-``
          1 '`
         /
        r

func FuncTet4(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncTet4 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of tet4 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

                      t
                      |
                      3
                     /|`.
                     ||  `,
                    / |    ',
                    | |      \
                   /  |       `.
                   |  |         `,
                  /   |           `,
                  |   |             \
                 /    |              `.
                 |    |                ',
                /     |                  \
                |     0.,,_               `.
               |     /     ``'-.,,__        `.
               |    /              ``''-.,,_  ',
              |    /                        `` 2 ,,s
              |  ,'                       ,.-``
             |  ,                    _,-'`
             ' /                 ,.'`
            | /             _.-``
            '/          ,-'`
           |/      ,.-``
           /  _,-``
          1 '`
         /
        r

func FuncTri10(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncTri10 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of tri10 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        s
        |
        2, (0,1)
        | ',
        |   ',
        5     '7
        |       ',
        |         ',
        8      9    '4
        |             ',
        | (0,0)         ', (1,0)
        0-----3-----6-----1 ---- r

func FuncTri15(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncTri15 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of tri15 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

          s
           ^
           |
         2
           @,(0,1)
           | ',
           |   ', 9
        10 @     @,
           |  14   ',   4
         5 @    @     @
           |           ',  8
        11 @  12@   @    '@
           |       13      ',
           |(0,0)            ', (1,0)
           @----@----@----@----@  --> r
         0      6    3    7     1

func FuncTri3(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncTri3 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of tri3 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        s
        |
        2, (0,1)
        | ',
        |   ',
        |     ',
        |       ',
        |         ',
        |           ',
        |             ',
        |               ',
        | (0,0)           ', (1,0)
        0-------------------1 ---- r

func FuncTri6(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    FuncTri6 calculates the shape functions (S) and derivatives of shape
    functions (dSdR) of tri6 elements at {r,s,t} natural coordinates. The
    derivatives are calculated only if derivs==true.

        s
        |
        2, (0,1)
        | ',
        |   ',
        |     ',
        |       ',
        5         '4
        |           ',
        |             ',
        |               ',
        | (0,0)           ', (1,0)
        0---------3---------1 ---- r

func IntPointsFindSet(cellKind int, setName string) (P [][]float64)
    IntPointsFindSet finds set of integration points by cell kind and set name

func QuadPointsGaussLegendre(ndim, npts int) (pts [][]float64)
    QuadPointsGaussLegendre generate quadrature points for Gauss-Legendre
    integration

        npts -- is the total number of points; e.g. 27 for 3D (boxes)

func QuadPointsWilson5(w0input float64, p4stable bool) (pts [][]float64)
    QuadPointsWilson5 generates 5 integration points according to Wilson's
    Appendix G-7 formulae

        w0input  -- if w0input > 0, use this value instead of default w0=8/3 (corner)
        p4stable -- if true, use w0=0.004 and wa=0.999 to mimic 4-point rule

func QuadPointsWilson8(wbinput float64) (pts [][]float64)
    QuadPointsWilson8 generates 8 integration points according to Wilson's
    Appendix G-7 formulae

        wbinput -- if wbinput > 0, use this value instead of default wb=40/49

func QuadPointsWilson9(w0input float64, p8stable bool) (pts [][]float64)
    QuadPointsWilson9 computes the 9-points for hexahedra according to Wilson's
    Appendix G-7 formulae

        w0input  -- if w0input > 0, use this value instead of default w0=16/3 (corner)
        p8stable -- if true, use w0=0.008 and wa=0.999 to mimic 8-point rule


TYPES

type BoundaryData struct {
	LocalID int   // edge local id (edgeId) OR face local id (faceId)
	Cell    *Cell // cell
}
    BoundaryData holds ID of edge or face and pointer to Cell at boundary (edge
    or face)

type BoundaryDataSet []*BoundaryData
    BoundaryDataSet defines a set of BoundaryData

type Cell struct {

	// input
	ID       int    `json:"i"`  // identifier
	Tag      int    `json:"t"`  // tag
	Part     int    `json:"p"`  // partition id
	Disabled bool   `json:"d"`  // cell is disabled
	TypeKey  string `json:"y"`  // geometry type; e.g. "lin2"
	V        []int  `json:"v"`  // vertices
	EdgeTags []int  `json:"et"` // edge tags (2D or 3D)
	FaceTags []int  `json:"ft"` // face tags (3D only)
	NurbsID  int    `json:"b"`  // id of NURBS (or something else) that this cell belongs to
	Span     []int  `json:"s"`  // span in NURBS

	// derived
	TypeIndex int        `json:"-"` // type index of cell. converted from TypeKey
	Gndim     int        `json:"-"` // geometry ndim
	X         *la.Matrix `json:"-"` // all vertex coordinates [nverts][ndim]
}
    Cell holds cell data (in e.g. from msh file)

func (o *Cell) String() string
    String returns a JSON representation of *Cell

type CellSet []*Cell
    CellSet defines a set of cells

type Edge struct {
	Verts VertexSet       // vertices on edge
	Bdata BoundaryDataSet // cells attached to edge, including which local edge id of cell is attached
}
    Edge holds the vertices and cells attached to an edge

type EdgeKey struct {
	A int // id of one vertex on edge
	B int // id of another vertex on edge
	C int // id of a third vertex on edge or the number of mesh vertices if edge has only 2 vertices
}
    EdgeKey holds 3 sorted numbers to identify an edge

type EdgesMap map[EdgeKey]*Edge
    EdgesMap is a map of edges

func (o *EdgesMap) Split() (internal, boundary EdgesMap)
    Split splits map into two sets: internal and boundary edges NOTE: boundary
    edge is determined by checking if edge is shared by only cell only

type Integrator struct {

	// input data
	Ctype  int         // cell type index
	Nverts int         // number of vertices = len(X)
	Ndim   int         // space dimension = len(X[0]) == len(P[0])
	Npts   int         // number of integration points = len(P)
	P      [][]float64 // (Gauss) integration points [npts][ndim]

	// slices related to integration points
	ShapeFcns []la.Vector  // shape functions Sm @ all integ points [npts][nverts]
	RefGrads  []*la.Matrix // reference gradients gm = dSm(r)/dr @ all integ points [npts][nverts][ndim]

	JacobianMat *la.Matrix // jacobian matrix Jr of the mapping reference to general coords [ndim][ndim]
	InvJacobMat *la.Matrix // inverse of jacobian matrix [ndim][ndim]
	DetJacobian float64    // determinat of jacobian matrix
	// Has unexported fields.
}
    Integrator implements methods to perform numerical integration over a
    polyhedron/polygon

func NewIntegrator(ctype int, P [][]float64, pName string) (o *Integrator)
    NewIntegrator returns a new object to integrate over polyhedra/polygons
    (cells)

        ctype -- index of cell type; e.g. TypeQuad4
        P     -- integration points [npoints][ndim]. may be nil => default will be selected
        pName -- use integration points from database instead of P or default ones. may be ""

func (o *Integrator) EvalJacobian(X *la.Matrix, ip int)
    EvalJacobian computes the Jacobian of the mapping from general to reference
    space at integration point with index ip

                                    dx          dSⁿ
         x(r) = Σ Sⁿ(r) xⁿ    ⇒     —— = Σ xⁿ ⊗ ———
                n                   dr   n       dr

         ∂xi              dS
         ——— = Σ X[n,i] * ——[n,j]    ⇒    Jmat = Xᵀ · dSdr
         ∂rj   n          dr

                →     _                           _
               dx    |  ∂x0/∂r0  ∂x0/∂r1  ∂x0/∂r2  |                ∂xi
        Jmat = —— =  |  ∂x1/∂r0  ∂x1/∂r1  ∂x1/∂r2  |     Jmat[ij] = ———
                →    |_ ∂x2/∂r0  ∂x2/∂r1  ∂x2/∂r2 _|                ∂rj
               dr

        Input:
          X  -- coordinates of vertices of cell (polyhedron/polygon) [nverts][ndim]
          ip -- index of integration point

        Computed (stored):
          JacobianMat -- reference Jacobian matrix [ndim][ndim]
          InvJacobMat -- inverse of Jmat [ndim][ndim]
          DetJacobian -- determinat of the reference Jacobian matrix

func (o *Integrator) GetXip(X *la.Matrix) (Xip *la.Matrix)
    GetXip calculates coordinates Xip of integration points from X and P

        Input:
          X -- coordinates of vertices of cell (polyhedron/polygon) [nverts][ndim]
        Output:
          Xip -- general (non-reference) coordinate of integ points [npts][ndim]

func (o *Integrator) IntegrateSv(X *la.Matrix, f fun.Sv) (res float64)
    IntegrateSv integrates scalar function of vector argument over Cell

        Computes:

                ⌠⌠⌠   →       ⌠⌠⌠   → →     →       nip-1   →  →      →
          res = │││ f(x) dΩ = │││ f(x(r))⋅J(r) dΩr ≈  Σ   f(xi(ri))⋅J(ri)⋅wi
                ⌡⌡⌡           ⌡⌡⌡                    i=0
                   Ω             Ωr

        where (J = det(Jmat)):

           x(r) ≈ Σ Sⁿ(r) ⋅ xⁿ     ⇒     x[i] = Σ S[n] * X[n,i]     ⇒     x = Xᵀ ⋅ S
                  n                             n
        Input:
          X  -- coordinates of vertices of cell (polyhedron/polygon) [nverts][ndim]
          f  -- integrand function

func (o *Integrator) ResetP(P [][]float64, pName string)
    ResetP resets integration points

        P     -- integration points [npoints][ndim]. may be nil => default will be selected
        pName -- use integration points from database instead of P or default ones. may be ""

type Mesh struct {

	// input
	Verts VertexSet `json:"verts"` // vertices
	Cells CellSet   `json:"cells"` // cells

	// derived
	Ndim int       // max space dimension among all vertices
	Xmin []float64 // min(x) among all vertices [ndim]
	Xmax []float64 // max(x) among all vertices [ndim]

	// auxiliary
	Tmaps *TagMaps // map of tags
}
    Mesh defines mesh data

func GenQuadRegion(ctype, ndivR, ndivS int, circle bool, f func(i, j, nr, ns int) (x, y float64)) (o *Mesh)
    GenQuadRegion generates 2D region made of quads

        ctype  -- one of Type{Qua4,Qua8,Qua9,Qua12,Qua16,Qua17}
        ndivR  -- number of divisions (cells) along r (e.g. x)
        ndivS  -- number of divisions (cells) along s (e.g. y)
        circle -- connect last row (s=ndivS) with the previous one (s=0)

        f(i,j,nr,ns) -- is a function that computes the (x,y) coordinates of grid nodes
                        were nr=ndivR+1 and nr=ndivS+1

        example (to generate a rectangle):

             f := func(i, j, nr, ns int) (x, y float64) {
             	dx := (xmax - xmin) / float64(nr-1)
             	dy := (ymax - ymin) / float64(ns-1)
             	x = xmin + float64(i)*dx
             	y = ymin + float64(j)*dy
             	return
             }

        The boundaries are tagged as below

            34      3       23                  30
              @-----@------@               +-----------+
              |            |               |           |
              |            |               |           |
            4 @  vertices  @ 2          40 |   edges   | 20
              |            |               |           |
              |            |               |           |
              @-----@------@               +-----------+
            41      1       12                  10

func GenQuadRegionHL(ctype, ndivR, ndivS int, xmin, xmax, ymin, ymax float64) (o *Mesh)
    GenQuadRegionHL generates 2D region made of quads (high-level version of
    GenQuadRegion)

        NOTE: see GenQuadRegion for more details

func GenRing2d(ctype int, ndivR, ndivA int, r, R, alpha float64) (o *Mesh)
    GenRing2d generates mesh of quads representing a 2D ring

        ctype -- one of Type{Qua4,Qua8,Qua9,Qua12,Qua16,Qua17}
        ndivR -- number of divisions along radius
        ndivA -- number of divisions along alpha
        r     -- minimum radius
        R     -- maximum radius
        alpha -- maximum alpha
        NOTE: a circular region is created if maxA=2⋅π

func NewMesh(jsonString string) (o *Mesh)
    NewMesh creates mesh from json string

func Read(fn string) (o *Mesh)
    Read reads mesh and call CheckAndCalcDerivedVars

func (o *Mesh) Boundary(tag int) []int
    Boundary returns a list of indices of nodes on edge (2D) or face (3D) of
    boundary

        NOTE: will return empty list if tag is not available

func (o *Mesh) CheckAndCalcDerivedVars()
    CheckAndCalcDerivedVars checks input data and computes derived quantities
    such as the max space dimension, min(x) and max(x) among all vertices,
    cells' gndim, etc. This function will set o.Ndim, o.Xmin and o.Xmax. This
    function will also generate the maps of tags.

func (o *Mesh) ExtractCellCoords(cellID int) (X *la.Matrix)
    ExtractCellCoords extracts cell coordinates

        X -- matrix with coordinates [nverts][gndim]

func (o *Mesh) ExtractEdges() (edges EdgesMap)
    ExtractEdges find edges in mesh

type MeshIntegrator struct {
	M           *Mesh           // the mesh
	Ngoroutines int             // total number of go routines
	Integrators [][]*Integrator // all integrators [Ngoroutines][TypeNumMax]
}
    MeshIntegrator implements methods to perform numerical integration over a
    mesh

func NewMeshIntegrator(mesh *Mesh, Ngoroutines int) (o *MeshIntegrator)
    NewMeshIntegrator returns a new MeshIntegrator

func (o *MeshIntegrator) IntegrateSv(goroutineID int, f fun.Sv) (res float64)
    IntegrateSv integrates scalar function of vector argument over mesh

                ⌠⌠⌠   →
          res = │││ f(x) dΩ
                ⌡⌡⌡
                   Ω
        Input:
          goroutineId -- go routine id to use when performing optimisation (not to partition mesh)

type ShapeFunction func(S la.Vector, dSdR *la.Matrix, R la.Vector, derivs bool)
    ShapeFunction computes the shape function and derivatives

type TagMaps struct {
	VertexTag2verts map[int]VertexSet       // vertex tag => set of vertices
	CellTag2cells   map[int]CellSet         // cell tag => set of cells
	CellType2cells  map[int]CellSet         // cell type => set of cells
	CellPart2cells  map[int]CellSet         // partition number => set of cells
	EdgeTag2cells   map[int]BoundaryDataSet // edge tag => set of cells {cell,boundaryId}
	EdgeTag2verts   map[int]VertexSet       // edge tag => vertices on tagged edge [unique]
	FaceTag2cells   map[int]BoundaryDataSet // face tag => set of cells {cell,boundaryId}
	FaceTag2verts   map[int]VertexSet       // face tag => vertices on tagged edge [unique]
}
    TagMaps holds data for finding information using on tags

type Vertex struct {

	// input
	ID  int       `json:"i"` // identifier
	Tag int       `json:"t"` // tag
	X   []float64 `json:"x"` // coordinates (size==2 or 3)

	// auxiliary
	Entity interface{} `json:"-"` // any entity attached to this vertex
}
    Vertex holds vertex data (e.g. from msh file)

func (o *Vertex) String() string
    String returns a JSON representation of *Vert

type VertexSet []*Vertex
    VertexSet defines a set of vertices

func (o VertexSet) IDs() (ids []int)
    IDs returns the IDs of vertices in VertexSet

func (o VertexSet) Len() int
    Len returns the length of vertex set

func (o VertexSet) Less(i, j int) bool
    Less compares ides in vertex set

func (o VertexSet) Swap(i, j int)
    Swap swaps two entries in vertex set

```
