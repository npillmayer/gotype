package frame

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/dom/styledtree"
	"github.com/npillmayer/gotype/engine/frame/layout"
	"golang.org/x/net/html"
)

// Parameters for GraphViz drawing.
type graphParamsType struct {
	Fontname string
	NodeTmpl *template.Template
	EdgeTmpl *template.Template
	//PgedgeTmpl *template.Template
	//PgpgTmpl   *template.Template
}

// ToGraphViz outputs a diagram for a render tree. The diagram is in
// GraphViz (DOT) format. Clients have to provide the root node of
// the tree and a Writer.
func ToGraphViz(root *layout.Container, w io.Writer) {
	vizheader, err := template.New("dom").Parse(graphHeadTmpl)
	if err != nil {
		panic(err)
	}
	gparams := graphParamsType{Fontname: "Helvetica"}
	gparams.NodeTmpl = template.Must(template.New("domnode").Parse(domNodeTmpl))
	gparams.EdgeTmpl = template.Must(template.New("domedge").Parse(domEdgeTmpl))
	//gparams.PgedgeTmpl = template.Must(template.New("pgedge").Parse(pgEdgeTmpl))
	//gparams.PgpgTmpl = template.Must(template.New("pgpgedge").Parse(pgpgEdgeTmpl))
	err = vizheader.Execute(w, gparams)
	if err != nil {
		panic(err)
	}
	dict := make(map[*layout.Container]string, 4096)
	containers(root, w, dict, &gparams)
	w.Write([]byte("}\n"))
}

type node struct {
	C       *layout.Container
	ID      string
	Name    string
	IsCData bool
}

func stylenode(c *layout.Container) style.Styler {
	return styledtree.Node(c.StyleNode)
}

func nodename(c *layout.Container, dict map[*layout.Container]string) (string, bool) {
	h := stylenode(c).HtmlNode()
	cnt := len(dict) + 1
	if h == nil {
		return fmt.Sprintf("nil_%d", cnt), false
	}
	if h.Type == html.DocumentNode {
		return "root", false
	}
	if h.Type == html.TextNode {
		return fmt.Sprintf("CDATA_%d", cnt), true
	}
	symbol := "▶︎"
	if c.IsBlock() {
		symbol = "◼︎"
	}
	return fmt.Sprintf("%s%s_%d", h.Data, symbol, cnt), false
}

func containers(c *layout.Container, w io.Writer, dict map[*layout.Container]string,
	gparams *graphParamsType) {
	boxNode(c, w, dict, gparams)
	fmt.Printf("container has %d children\n", c.ChildCount())
	for i := 0; i < c.ChildCount(); i++ {
		child, ok := c.Child(i)
		if ok {
			ch := layout.Node(child)
			if ch == c {
				panic("recursion in render tree")
			}
			containers(ch, w, dict, gparams)
			domEdge(c, ch, w, dict, gparams)
		}
	}
}

func boxNode(c *layout.Container, w io.Writer, dict map[*layout.Container]string, gparams *graphParamsType) {
	id, iscdata := nodename(c, dict)
	if id == "" {
		l := len(dict) + 1
		id = fmt.Sprintf("node%05d", l)
	}
	dict[c] = id
	name := strings.Split(id, "_")[0]
	if err := gparams.NodeTmpl.Execute(w, &node{c, id, name, iscdata}); err != nil {
		panic(err)
	}
}

type edge struct {
	N1, N2 node
}

func domEdge(n1 *layout.Container, n2 *layout.Container, w io.Writer, dict map[*layout.Container]string,
	gparams *graphParamsType) {
	//
	//fmt.Printf("dict has %d entries\n", len(dict))
	id1 := dict[n1]
	id2 := dict[n2]
	e := edge{node{n1, "", id1, false}, node{n2, "", id2, false}}
	if err := gparams.EdgeTmpl.Execute(w, e); err != nil {
		panic(err)
	}
}

type pgedge struct {
	Name      string
	PropGroup *style.PropertyGroup
}

// --- Templates --------------------------------------------------------

const graphHeadTmpl = `digraph g {                                                                                                             
	graph [labelloc="t" label="" splines=true overlap=false rankdir = "LR"];
	graph [{{ .Fontname }} = "helvetica" fontsize=14] ;
	 node [fontname = "{{ .Fontname }}" fontsize=14] ;
	 edge [fontname = "{{ .Fontname }}" fontsize=14] ;
`

const domNodeTmpl = `{{ if .IsCData }}
{{ .ID }}	[ label={{ printf "%q" .Name }} shape=box style=filled fillcolor=white ] ;
{{ else }}
{{ .ID }}	[ label={{ printf "%q" .Name }} shape=ellipse style=filled fillcolor=white ] ;
{{ end }}
`

/*
const styleGroupTmpl = `{{ printf "pg%p" . }} [ style="filled" penwidth=1 fillcolor="white" shape="Mrecord" fontsize=12
	label=<<table border="0" cellborder="0" cellpadding="2" cellspacing="0" bgcolor="white">
<tr><td bgcolor="azure4" align="center" colspan="2"><font color="white">{{ .Name }}</font></td></tr>
{{ range .Properties }}
<tr><td align="right">{{ .Key }}:</td><td>{{ .Value }}</td></tr>
{{ else }}
<tr><td colspan="2">no styles</td></tr>
{{ end }}
	</table>> ] ;
`
*/

//const domEdgeTmpl = `{{ .N1.Name }} -> {{ .N2.Name }} [dir=none weight=1] ;
const domEdgeTmpl = `{{ .N1.Name }} -> {{ .N2.Name }} [weight=1] ;
`

/*
const pgEdgeTmpl = `{{ .Name }} -> {{ printf "pg%p" .PropGroup }} [dir=none weight=1 style="dashed"] ;
`

const pgpgEdgeTmpl = `{{ index . 0 | printf "pg%p"  }} -> {{ index . 1 | printf "pg%p" }} [dir=none weight=1 style="dashed"] ;
`
*/
