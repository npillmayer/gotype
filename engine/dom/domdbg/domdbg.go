package domdbg

import (
	"html/template"
	"io"

	"github.com/npillmayer/gotype/engine/dom"
)

type GraphParams struct {
	Fontname string
}

func ToGraphViz(dom dom.RODomNode, w io.Writer) {
	tmpl, err := template.New("dom").Parse(graphHeadTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, GraphParams{"Helvetica"})
	if err != nil {
		panic(err)
	}
}

func domNode(node dom.RODomNode, w io.Writer) {
	tmpl, err := template.New("domnode").Parse(domNodeTmpl)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, node)
	if err != nil {
		panic(err)
	}
}

// --- Templates --------------------------------------------------------

const graphHeadTmpl = `digraph g {                                                                                                             
  graph [labelloc="t" label="" splines=true overlap=false rankdir = "TB"];
  graph [{{ .Fontname }} = "helvetica" fontsize=14] ;
   node [fontname = "{{ .Fontname }}" fontsize=14] ;
   edge [fontname = "{{ .Fontname }}" fontsize=14] ;
`

const domNodeTmpl = `{{ if .IsText }}
{{ .String }}	[ shape=box style=filled fillcolor=white ] ;
{{ else }}
{{ .String }}	[ shape=ellipse style=filled fillcolor=white ] ;
{{ end }}
`

const styleGroupTmpl = `{{ .Name }} [ style="filled" penwidth=1 fillcolor="white" shape="Mrecord"
    label=<<table border="0" cellborder="0" cellpadding="2" bgcolor="white">
      <tr><td bgcolor="azure4" align="center" colspan="2"><font color="white">{{ .Name }}</font></td></tr>
      {{ range .Properties }}
      <tr><td align="right">{{ .Key }}:</td><td>{{ .Value }}</td></tr>
      {{ else }}
      <tr><td colspan="2">no styles</td></tr>
      {{ end }}
    </table>> ] ;
`

const domEdgeTmpl = `{{ .N1.String -> .N2.String }} [dir=none, weight=1] ;
`

// ----------------------------------------------------------------------
