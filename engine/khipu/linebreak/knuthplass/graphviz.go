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

func (kp *linebreaker) toGraphViz(cursor *khipu.Cursor, w io.Writer) {
	tmpl, _ := template.New("graph").Parse(graphHeader)
	gparams := graphParamsType{Fontname: "Helvetica"}
	_ = tmpl.Execute(w, gparams)
	boxT := template.Must(template.New("box").Parse(boxTmpl))
	//edgeT := template.Must(template.New("edge").Parse(edgeTmpl))
	// cursor.Next()
	cursor.Next()
	box := &n{
		khipu: cursor.Khipu(),
		Mark:  cursor.Mark(),
	}
	t := getText(box)
	fmt.Printf("Text = %v", t)
	if err := boxT.Execute(w, box); err != nil {
		panic(err)
	}
	w.Write([]byte("}\n"))
}

type n struct {
	khipu *khipu.Khipu
	Mark  khipu.Mark
	Text  string
}

type e struct {
	N1, N2      *n
	Cost, Total int32
}

const graphHeader = `digraph g {                                                                                                             
  graph [labelloc="t" label="" splines=true overlap=false];
  graph [{{ .Fontname }} = "helvetica" fontsize=12] ;
   node [fontname = "{{ .Fontname }}" fontsize=12] ;
   edge [fontname = "{{ .Fontname }}" fontsize=12] ;
`

const boxTmpl = `break{{ .Mark.Position }}	[ label={{ .Text }} shape=box style=filled fillcolor=grey95 fontname="Courier" fontsize=11.0 ] ;
`

//const domEdgeTmpl = `{{ .N1.Name }} -> {{ .N2.Name }} [dir=none weight=1] ;
const edgeTmpl = `{{ .N1.Mark.Position }} -> {{ .N2.Mark.Position }} [weight=1] ;
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
