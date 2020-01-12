package knuthplass

import (
	"fmt"
	"io"
	"text/template"

	"github.com/npillmayer/gotype/engine/khipu"
)

// Parameters for GraphViz drawing.
type graphParamsType struct {
	Fontname string
	//NodeTmpl       *template.Template
}

func allBreakpointBoxes(kp *linebreaker, kh *khipu.Khipu, boxT *template.Template,
	w io.Writer) map[int]*n {
	//
	breakBoxes := make(map[int]*n)
	for _, fb := range kp.nodes {
		box := makeBox(fb, kh)
		if err := boxT.Execute(w, box); err != nil {
			panic(err)
		}
		breakBoxes[fb.mark.Position()] = box
	}
	return breakBoxes
}

func allEdges(kp *linebreaker, kh *khipu.Khipu, boxes map[int]*n, edgeT *template.Template,
	w io.Writer) {
	//
	for _, edge := range kp.Edges() {
		T().Debugf("output of edge %v", edge)
		e := &e{}
		e.N1 = boxes[edge.from]
		e.N2 = boxes[edge.to]
		e.Cost = edge.cost
		e.Total = edge.total
		e.Line = edge.linecount
		start := 0
		if edge.from >= 0 {
			start = edge.from
		}
		e.Text = kh.Text(start, edge.to)
		if err := edgeT.Execute(w, e); err != nil {
			panic(err)
		}
	}
}

func makeBox(fb *feasibleBreakpoint, kh *khipu.Khipu) *n {
	box := &n{
		khipu: kh,
		Mark:  fb.mark,
	}
	if fb.mark.Position() < 0 {
		box.Text = "root"
		box.Name = "root"
	} else {
		box.Text = fmt.Sprintf("%s\\n%v", getText(box), fb.mark.Knot())
		box.Name = fmt.Sprintf("break%d", fb.mark.Position())
	}
	T().Debugf("Text %v = '%v'", fb.mark.Knot(), box.Text)
	return box
}

func (kp *linebreaker) toGraphViz(cursor *khipu.Cursor, results map[int][]khipu.Mark,
	w io.Writer) {
	//
	tmpl, _ := template.New("graph").Parse(graphHeader)
	gparams := graphParamsType{Fontname: "Helvetica"}
	_ = tmpl.Execute(w, gparams)
	boxT := template.Must(template.New("box").Parse(boxTmpl))
	edgeT := template.Must(template.New("edge").Parse(edgeTmpl))
	// cursor.Next()
	boxes := allBreakpointBoxes(kp, cursor.Khipu(), boxT, w)
	allEdges(kp, cursor.Khipu(), boxes, edgeT, w)
	w.Write([]byte("}\n"))
}

type n struct {
	khipu *khipu.Khipu
	Mark  khipu.Mark
	Name  string
	Text  string
}

type e struct {
	N1, N2      *n
	Cost, Total int32
	Line        int
	Text        string
}

const graphHeader = `digraph g {                                                                                                             
  graph [labelloc="t" label="" splines=true overlap=false];
  graph [{{ .Fontname }} = "helvetica" fontsize=12] ;
   node [fontname = "{{ .Fontname }}" fontsize=12] ;
   edge [fontname = "{{ .Fontname }}" fontsize=9 labelfontsize=9 labeldistance=5.0] ;
`

const boxTmpl = `{{ .Name }}	[ label="{{ .Text }}" shape=box style=filled fillcolor=grey95 fontname="Courier" fontsize=11.0 ] ;
`

//const edgeTmpl = `{{ .N1.Name }} -> {{ .N2.Name }} [weight=1] ;
const edgeTmpl = `{{.N1.Name}} -> {{.N2.Name}} [weight=1 label="{{.Cost}} of\n{{.Total}}\nl={{.Line}}" tooltip="{{ .Text }}"] ;
`

// ----------------------------------------------------------------------

func getText(n *n) string {
	if n.Mark.Position() < 0 {
		return "root"
	}
	cursor := khipu.NewCursor(n.khipu)
	// walk cursor to mark
	for i := 0; i <= n.Mark.Position(); i++ {
		cursor.Next()
	}
	// walk back until a box is found
	root := false
	for cursor.Knot().Type() != khipu.KTTextBox {
		cursor.Prev()
		if !cursor.IsValidPosition() {
			n.Text = "root"
			root = true
		}
	}
	if !root {
		n.Text = cursor.AsTextBox().Text()
	}
	return n.Text
}
