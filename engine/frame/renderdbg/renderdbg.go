package renderdbg

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/frame/layout"
)

// Parameters for GraphViz drawing.
type graphParamsType struct {
	Fontname    string
	StyleGroups []string
	BoxTmpl     *template.Template
	EdgeTmpl    *template.Template
}

// ToGraphViz creates a graphical representation of a render tree.
// It produces a DOT file format suitable as input for Graphviz, given a Writer.
func ToGraphViz(l *layout.Layouter, w io.Writer) {
	header, err := template.New("renderTree").Parse(graphHeadTmpl)
	if err != nil {
		panic(err)
	}
	gparams := graphParamsType{Fontname: "Helvetica"}
	gparams.BoxTmpl, _ = template.New("box").Funcs(
		template.FuncMap{
			"shortstring": shortText,
		}).Parse(boxTmpl)
	gparams.EdgeTmpl = template.Must(template.New("boxedge").Parse(edgeTmpl))
	err = header.Execute(w, gparams)
	if err != nil {
		panic(err)
	}
	dict := make(map[*layout.Container]string, 4096)
	boxes(l.BoxRoot(), w, dict, &gparams)
	w.Write([]byte("}\n"))
}

func boxes(c *layout.Container, w io.Writer, dict map[*layout.Container]string, gparams *graphParamsType) {
	box(c, w, dict, gparams)
	if c.ChildCount() > 0 {
		children := c.Children()
		for _, ch := range children {
			child := ch.Payload.(*layout.Container)
			boxes(child, w, dict, gparams)
			edge(c, child, w, dict, gparams)
		}
	}
}

func box(c *layout.Container, w io.Writer, dict map[*layout.Container]string, gparams *graphParamsType) {
	name := dict[c]
	if name == "" {
		sz := len(dict) + 1
		name = fmt.Sprintf("node%05d", sz)
		dict[c] = name
	}
	if err := gparams.BoxTmpl.Execute(w, &cbox{c, c.DOMNode, name}); err != nil {
		panic(err)
	}
}

// Helper struct
type cbox struct {
	C    *layout.Container
	N    *dom.W3CNode
	Name string
}

func shortText(n *dom.W3CNode) string {
	h := n.HTMLNode()
	s := "\"\\\""
	if len(h.Data) > 10 {
		s += h.Data[:10] + "...\\\"\""
	} else {
		s += h.Data + "\\\"\""
	}
	s = strings.Replace(s, "\n", `\\n`, -1)
	s = strings.Replace(s, "\t", `\\t`, -1)
	s = strings.Replace(s, " ", "\u2423", -1)
	return s
}

type cedge struct {
	N1, N2 cbox
}

func edge(c1 *layout.Container, c2 *layout.Container, w io.Writer, dict map[*layout.Container]string,
	gparams *graphParamsType) {
	//
	//fmt.Printf("dict has %d entries\n", len(dict))
	name1 := dict[c1]
	name2 := dict[c2]
	e := cedge{cbox{c1, c1.DOMNode, name1}, cbox{c2, c2.DOMNode, name2}}
	if err := gparams.EdgeTmpl.Execute(w, e); err != nil {
		panic(err)
	}
}

// --- Templates --------------------------------------------------------

const graphHeadTmpl = `digraph g {                                                                                                             
  graph [labelloc="t" label="" splines=true overlap=false rankdir = "LR"];
  graph [{{ .Fontname }} = "helvetica" fontsize=14] ;
   node [fontname = "{{ .Fontname }}" fontsize=14] ;
   edge [fontname = "{{ .Fontname }}" fontsize=14] ;
`

const boxTmpl = `{{ if eq .N.NodeName "#text" }}
{{ .Name }}	[ label={{ shortstring .N }} shape=box style=filled fillcolor=grey95 fontname="Courier" fontsize=11.0 ] ;
{{ else }}
{{ .Name }}	[ label={{ printf "%q" .N.NodeName }} shape=ellipse style=filled fillcolor=lightblue3 ] ;
{{ end }}
`

//const domEdgeTmpl = `{{ .N1.Name }} -> {{ .N2.Name }} [dir=none weight=1] ;
const edgeTmpl = `{{ .N1.Name }} -> {{ .N2.Name }} [weight=1] ;
`
