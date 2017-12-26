package arithmetic

/* This is mainly a proxy class for a polyclip package from
 * https://github.com/akavel/polyclip-go.
 *
 */

import (
	"bytes"
	"fmt"

	pc "github.com/akavel/polyclip-go"
	"github.com/npillmayer/gotype/gtcore/config/tracing"
	dec "github.com/shopspring/decimal"
)

var L tracing.Trace = tracing.SyntaxTracer

type Path interface {
	fmt.Stringer
	AddPoint(pr Pair)
	AddSubpath(p Path)
	GetPoint(i int) Pair
	DeletePoint(i int)
	Subpath(from, to int)
	Length() int
	IsCycle() bool
	Cycle()
	Copy() Path
	Reverse()
}

type Polygon struct {
	contours pc.Polygon
	isCycle  bool
}

func NullPath() Path {
	pg := &Polygon{}
	return pg
}

func NewPath(prs []Pair) Path {
	pg := NullPath()
	for _, pr := range prs {
		pg.AddPoint(pr)
	}
	return pg
}

func (pg *Polygon) getContour() *pc.Contour {
	if len(pg.contours) == 0 {
		pg.contours.Add(pc.Contour{})
	}
	return &pg.contours[0]
}

/* Interface Path.
 */
func (pg *Polygon) AddPoint(pr Pair) {
	L.Debugf("add point to path = %s", pr)
	pg.getContour().Add(Pr2Pt(pr))
}

/* Interface Path.
 */
func (pg *Polygon) AddSubpath(p Path) {
	L.Debugf("add subpath = %v", p)
	for i := 0; i < p.Length(); i++ {
		pg.AddPoint(p.GetPoint(i))
	}
}

/* Interface Path.
 */
func (pg *Polygon) GetPoint(i int) Pair {
	//L.Errorf("not yet implemented: Polygon.GetPoint()")
	pt := (*pg.getContour())[i]
	return Pt2Pr(pt)
}

/* Interface Path.
 */
func (pg *Polygon) DeletePoint(i int) {
	c := *pg.getContour()
	c = append(c[:i], c[i+1:]...)
	pg.contours[0] = c
}

/* Interface Path.
 */
func (pg *Polygon) Subpath(from, to int) {
	var contour pc.Contour = *pg.getContour()
	pg.contours[0] = contour[from : to+1]
	//L.Debugf("subpath = %v", pg.getContour())
}

/* Interface Path.
 */
func (pg *Polygon) Length() int {
	return pg.contours.NumVertices()
}

/* Interface Path.
 */
func (pg *Polygon) Cycle() {
	pg.isCycle = true
}

/* Interface Path.
 */
func (pg *Polygon) IsCycle() bool {
	return pg.isCycle
}

/* Interface Path.
 */
func (pg *Polygon) Copy() Path {
	cs := pg.contours.Clone()
	p := &Polygon{contours: cs}
	return p
}

/* Interface Path.
 */
func (pg *Polygon) Reverse() {
	c := *pg.getContour()
	for i := 0; i < len(c)/2; i++ { // swap points in place
		j := len(c) - i - 1
		c[i], c[j] = c[j], c[i]
	}
}

// check assignability
var _ Path = &Polygon{}

// ---------------------------------------------------------------------------

/* Pretty Stringer.
 */
func (pg *Polygon) String() string {
	var s bytes.Buffer
	start := true
	for _, pt := range *pg.getContour() {
		if !start {
			s.WriteString("--")
		}
		s.WriteString(fmt.Sprintf("(%.2f,%.2f)", pt.X, pt.Y))
		start = false
	}
	if pg.IsCycle() {
		s.WriteString("--cycle")
	}
	return s.String()
}

// === Utilities =============================================================

func Pr2Pt(pr Pair) pc.Point {
	px, _ := pr.XPart().Float64()
	py, _ := pr.YPart().Float64()
	return pc.Point{
		X: px,
		Y: py,
	}
}

func Pt2Pr(pt pc.Point) Pair {
	px := dec.NewFromFloat(pt.X)
	py := dec.NewFromFloat(pt.Y)
	return MakePair(px, py)
}
