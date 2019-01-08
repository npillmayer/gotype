package cmd

/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer <norbert@pillmayer.com>

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Norbert Pillmayer nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

----------------------------------------------------------------------

Command line interface for HTML/CSS styling.

This module creates a possibly interactive CLI, configures the input (file),
trace output, an interpreter instance, and necessary backend modules.
Then sets up either an interactive shell or a file input instance.

Usage: style [-d <output-dir>] [-x] [-m] inputfile.html

	--outdir    -d   output directory, default = .
	--debug     -x   debug mode,       default = off
	--vi        -m   vi editing mode,  default = off

*/

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/douceuradapter"
	"github.com/npillmayer/gotype/engine/dom/styledtree/builder"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

const sWelcomeMessage = "Welcome to Gallery style [V0.1 experimental]"
const stoolname = "style"

// styleCmd represents the style command
var styleCmd = &cobra.Command{
	Use:   stoolname,
	Short: "Style an HTML document",
	Long:  `Gallery includes an HTML/CSS styling engine.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runStyleCmd,
}

/* COBRA init method. Defines command and flags.
 */
func init() {
	rootCmd.AddCommand(styleCmd)
	styleCmd.Flags().BoolP("vi", "m", false, "Set vi editing mode")
	styleCmd.Flags().StringP("outdir", "d", ".", "Output directory")
	styleCmd.Flags().BoolP("debug", "x", false, "Debug mode")
}

// A type to instantiate a REPL interpreter.
type StyleREPL struct {
	BaseREPL
	//gyintp *gallery.GalleryInterpreter
}

/* Bridge to the Gallery interpreter.
 */
func (frepl *StyleREPL) InterpretCommand(input string) {
	frepl.interpreter.InterpretCommand(input)
}

/* The StyleCmd command (style).
 */
func runStyleCmd(cmd *cobra.Command, args []string) {
	fmt.Println(sWelcomeMessage)
	welcomeMessage = sWelcomeMessage
	var inputfilename string
	if len(args) > 0 {
		T.Infof("input file is %s", args[0])
		inputfilename = args[0]
	}
	//tracing.Tracefile = tracing.ConfigTracing(inputfilename)
	defer tracing.Tracefile.Close()
	startStyleInput(inputfilename)
}

// Start Gallery/Style input. If a filename is given, opens the file and reads
// from there. Otherwise starts an interactive shell (REPL).
func startStyleInput(inputfilename string) {
	if inputfilename == "" {
		config.IsInteractive = true
		repl := NewStyleREPL() // go into interactive mode
		log.SetOutput(repl.readline.Stderr())
		defer repl.readline.Close()
		repl.doLoop()
	} else {
		input, err := antlr.NewFileStream(inputfilename) // TODO refactor to get rid of ANTLR
		if err != nil {
			T.Errorf("cannot open input file")
		} else {
			config.IsInteractive = false
			defer func() {
				if r := recover(); r != nil {
					T.Errorf("error executing Gallery statement!")
				}
			}()
			//interpreter := gallery.NewGalleryInterpreter(true)
			//interpreter.ParseStatements(input)
			T.Debugf("input = %s", input)
			panic("interpreting file not implemented: style")
		}
	}
}

/* Set up a new REPL entity. It contains a readline-instance (for putting
 * out a prompt and read in a line) and a PMMetaPost parser. The REPL will
 * then forward PMMetaPost statements to the parser.
 */
func NewStyleREPL() *StyleREPL {
	rl := NewReadline(stoolname)
	repl := &StyleREPL{}
	repl.readline = rl
	repl.interpreter = repl // we are our own bridge to the interpreter
	repl.interpreter = &styleIntp{}
	repl.toolname = stoolname
	repl.helper = styleCmdHelp
	return repl
}

type styleIntp struct {
	dom *html.Node
	css cssom.StyleSheet
	//rulesTree *style.rulesTree
	query     *goquery.Document
	selector  string
	selection *goquery.Selection
	styleTree cssom.StyledNode
}

func (intp *styleIntp) InterpretCommand(line string) {
	line = strings.TrimSpace(line)
	words := strings.Fields(line)
	command := ""
	if len(words) == 0 {
		return
	}
	command = words[0]
	switch command {
	case "html":
		if len(words) < 2 {
			T.Errorf("need <file> argument for HTML load")
			return
		}
		file, err := os.Open(words[1]) // For read access.
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		doc, err := html.Parse(file)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		T.Infof("OK loading HTML file")
		intp.dom = doc
	case "css":
		if len(words) < 2 {
			T.Errorf("need <file> argument for CSS load")
			return
		}
		bytes, err := ioutil.ReadFile(words[1]) // For read access.
		css, err := parser.Parse(string(bytes))
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		T.Infof("OK loading CSS file")
		intp.css = douceuradapter.Wrap(css)
		//intp.rulesTree = style.NewRulesTree(css)
	case "query":
		if intp.dom == nil {
			T.Errorf("need HTML loaded for query (command 'html')")
			return
		}
		if len(words) < 2 {
			T.Errorf("need <selector> argument for DOM query")
			return
		}
		selectorstring := line[5:]
		if intp.query == nil {
			intp.query = goquery.NewDocumentFromNode(intp.dom)
		}
		intp.selection = intp.query.Find(selectorstring)
		intp.selector = selectorstring
		T.Infof("selection of %d nodes", intp.selection.Length())
	case "text":
		if intp.dom == nil {
			T.Errorf("need HTML loaded for query (command 'html')")
			return
		}
		text := intp.selection.Text()
		T.Infof("Text Content =\n\"%s\"", text)
	case "name":
		if intp.dom == nil {
			T.Errorf("need HTML loaded for query (command 'html')")
			return
		}
		T.Infof("Nodes:")
		intp.selection.Each(elementName)
		/*
			case "match":
				if intp.dom == nil || intp.css == nil {
					T.Errorf("need HTML and CSS loaded for style-rule matching")
					return
				}
				m := intp.rulesTree.FilterMatchesFor(intp.selection.Nodes[0])
				m.SortProperties()
		*/
	case "style":
		if intp.dom == nil || intp.css == nil {
			T.Errorf("need HTML and CSS loaded for styling")
			return
		}
		//tree, err := style.ConstructStyledNodeTree(intp.dom, intp.rulesTree)
		c := cssom.NewCSSOM(nil)
		err := c.AddStylesForScope(nil, intp.css, cssom.Author)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		//intp.styleTree = tree
		T.Debugf("styling DOM")
		//_, err = cssom.Style(intp.dom)
		tree, err := c.Style(intp.dom, builder.Builder{})
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		tracing.With(T).Dump("tree", tree)
		intp.styleTree = tree
	case "tree":
		if intp.styleTree == nil {
			T.Errorf("need to first style the DOM")
			return
		}
	}
}

func styleCmdHelp(out io.Writer) {
	io.WriteString(out, "html <file>:        load HTML from file\n")
	io.WriteString(out, "css <file>:         load CSS from file\n")
	io.WriteString(out, "query <selector>:   select a set of nodes\n")
	io.WriteString(out, "name:               element name(s) of selection\n")
	io.WriteString(out, "text:               textual content of selection\n")
	//io.WriteString(out, "match:              match style-rules for selection\n")
	io.WriteString(out, "style:              style complete DOM\n")
	io.WriteString(out, "tree:               display styled boxes tree\n")
}

func elementName(i int, sel *goquery.Selection) {
	node := sel.Nodes[0]
	if node.Type == html.TextNode {
		T.Infof("- (text)")
	} else if node.Type == html.ElementNode {
		T.Infof("- %s", node.Data)
	} else {
		T.Infof("- (unknown)")
	}
}
