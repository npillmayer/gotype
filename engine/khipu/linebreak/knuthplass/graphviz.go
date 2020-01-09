package knuthplass

import (
	"io"
	"text/template"
)

// Parameters for GraphViz drawing.
type graphParamsType struct {
	Fontname string
	//NodeTmpl       *template.Template
}

func toGraphViz(kp *linebreaker, w io.Writer) {
	tmpl, _ := template.New("graph").Parse(graphHeader)
	gparams := graphParamsType{Fontname: "Helvetica"}
	_ = tmpl.Execute(w, gparams)
	w.Write([]byte("}\n"))
}

const graphHeader = `digraph g {                                                                                                             
  graph [labelloc="t" label="" splines=true overlap=false];
  graph [{{ .Fontname }} = "helvetica" fontsize=12] ;
   node [fontname = "{{ .Fontname }}" fontsize=12] ;
   edge [fontname = "{{ .Fontname }}" fontsize=12] ;
`
