package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"runtime"
	"text/template"
	"time"

	"os"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/npillmayer/gotype/gtcore/unicode/ucd"
)

/*
----------------------------------------------------------------------
Generator for Unicode UAX#14 code-point classes.
UAX#14 is about line break/wrap. For more information see
http://unicode.org/reports/tr14/

Classes are generated from a UAX#14 companion file: "LineBreak.txt".
This is the definite source for UAX#14 code-point classes. The
generator looks for it in a directory "$GOPATH/etc/".
----------------------------------------------------------------------
*/

var logger = log.New(os.Stderr, "UAX#14 generator: ", log.LstdFlags)

// flag: verbose output ?
var verbose bool

var uax14classnames = []string{"AI", "AL", "B2", "BA", "BB", "BK", "CB", "CJ",
	"CL", "CM", "CP", "CR", "EB", "EM", "EX", "GL", "H2", "H3", "HL", "HY",
	"ID", "IN", "IS", "JL", "JT", "JV", "LF", "NL", "NS", "NU", "OP", "PO",
	"PR", "QU", "RI", "SA", "SG", "SP", "SY", "WJ", "XX", "ZW", "ZWJ"}

/*
//go:generate stringer -type=UAX14Class
const (
	// sequence matters!
	AIClass UAX14Class = iota
	ALClass
	BAClass
	BBClass
	BKClass
	CJClass
	CLClass
	CMClass
	CPClass
	CRClass
	EBClass
	EMClass
	EXClass
	GLClass
	H2Class
	H3Class
	HLClass
	HYClass
	IDClass
	ISClass
	JLClass
	JTClass
	JVClass
	LFClass
	NLClass
	NSClass
	NUClass
	OPClass
	POClass
	PRClass
	QUClass
	RIClass
	SAClass
	SPClass
	SYClass
	WJClass
	XXClass
	ZWClass
	ZWJClass
	sot       // pseudo class
	eot       // pseudo class
	optSpaces // pseudo class
)

var uax14ClassFromString = map[string]UAX14Class{
	"CM": CMClass, "BA": BAClass, "LF": LFClass, "BK": BKClass, "CR": CRClass,
	"NL": NLClass, "SP": SPClass, "EX": EXClass, "QU": QUClass, "AL": ALClass,
	"PR": PRClass, "PO": POClass, "OP": OPClass, "CP": CPClass, "IS": ISClass,
	"HY": HYClass, "SY": SYClass, "NU": NUClass, "AI": AIClass, "BB": BBClass,
	"GL": GLClass, "SA": SAClass, "JT": JTClass, "JV": JVClass, "JL": JLClass,
	"NS": NSClass, "ZW": ZWClass, "ZWJ": ZWJClass, "WJ": WJClass, "CL": CLClass,
	"ID": IDClass, "CJ": CJClass, "H2": H2Class, "H3": H3Class, "EB": EBClass,
	"EM": EMClass, "HL": HLClass, "RI": RIClass, "XX": XXClass,
}
*/

// Load the Unicode UAX#14 definition file: LineBreak.txt
func loadUnicodeLineBreakFile() (map[string][]rune, error) {
	if verbose {
		logger.Printf("reading LineBreak.txt")
	}
	defer timeTrack(time.Now(), "loading LineBreak.txt")
	gopath := os.Getenv("GOPATH")
	f, err := os.Open(gopath + "/etc/LineBreak.txt")
	if err != nil {
		fmt.Printf("ERROR loading " + gopath + "/etc/LineBreak.txt\n")
		return nil, err
	}
	defer f.Close()
	p := ucd.NewUCDParser(f)
	lbcs := make(map[string]*arraylist.List, len(uax14classnames))
	for p.Next() {
		from, to := p.Range(0)
		brclzstr := p.String(1)
		list := lbcs[brclzstr]
		if list == nil {
			list = arraylist.New()
		}
		for r := from; r <= to; r++ {
			list.Add(r)
		}
		lbcs[brclzstr] = list
	}
	err = p.Err()
	if err != nil {
		log.Fatal(err)
	}
	runeranges := make(map[string][]rune)
	for k, v := range lbcs {
		runelist := make([]rune, lbcs[k].Size())
		it := v.Iterator()
		i := 0
		for it.Next() {
			runelist[i] = it.Value().(rune)
			i++
		}
		runeranges[k] = runelist
	}
	return runeranges, err
}

// --- Templates --------------------------------------------------------

var header string = `package uax14

// This file has been generated -- you probably should not edit it.
// 
// BSD License, Copyright (c) 2017-18, Norbert Pillmayer

import (
    "strconv"
    "unicode"

    "golang.org/x/text/unicode/rangetable"
)
`

var templateClassType string = `
// Type for UAX#14 code-point classes.
// Must be convertable to int.
type UAX14Class int

// Will be initialized in SetupUAX14Classes()
var rangeFromUAX14Class []*unicode.RangeTable
`

var templateRangeTableVars string = `
// Range tables for UAX#14 code-point classes.
// Will be initialized with setupUAX14Classes().
// Clients can check with unicode.Is(..., rune){{$i:=0}}
var {{range .}}{{$i = inc $i}}{{.}}, {{if modten $i}}
    {{end}}{{end}}unused *unicode.RangeTable
`

var templateClassConsts string = `
// These are all the UAX#14 breaking classes.
const ( {{$i:=0}}
{{range  .}}    {{.}}Class UAX14Class = {{$i}}{{$i = inc $i}}
{{end}})
`

//{{range  $k,$v := .}}    {{$k}}Class UAX14Class = {{$v}}

var templateClassStringer string = `
const _UAX14Class_name = "{{range $c,$name := .}}{{$name}}Class{{end}}"

var _UAX14Class_index = [...]uint16{0, 7, 14, 21, 28, 35, 42, 49, 56, 63, 70,
    77, 84, 91, 98, 105, 112, 119, 126, 133, 140, 147, 154, 161, 168, 175, 182,
    189, 196, 203, 210, 217, 224, 231, 238, 245, 252, 259, 266, 273, 280, 287,
    294, 302}

func (c UAX14Class) String() string {
    if c == sot {
        return "sot"
    } else if c == eot {
        return "eot"
    } else if c < 0 || c >= UAX14Class(len(_UAX14Class_index)-1) {
        return "UAX14Class(" + strconv.FormatInt(int64(c), 10) + ")"
    }
    return _UAX14Class_name[_UAX14Class_index[c]:_UAX14Class_index[c+1]]
}
`

var templateRangeForClass string = `{{$i:=0}}{{range .}}{{if notfirst $i}}, {{if modeight $i}}
    {{end}}{{end}}{{$i = inc $i}}{{printf "%+q" .}}{{end}}`

// Helper functions for templates
var funcMap template.FuncMap = template.FuncMap{
	"modten": func(i int) bool {
		return i%10 == 0
	},
	"modeight": func(i int) bool {
		return (i+2)%8 == 0
	},
	"inc": func(i int) int {
		return i + 1
	},
	"notfirst": func(i int) bool {
		return i > 0
	},
}

func makeTemplate(name string, templString string) *template.Template {
	if verbose {
		logger.Printf("creating %s", name)
	}
	t := template.Must(template.New(name).Funcs(funcMap).Parse(templString))
	return t
}

// --- Main -------------------------------------------------------------

func generateRanges(w *bufio.Writer, codePointLists map[string][]rune) {
	defer timeTrack(time.Now(), "generate range tables")
	w.WriteString("\nfunc setupUAX14Classes() {\n")
	w.WriteString("    rangeFromUAX14Class = make([]*unicode.RangeTable, int(ZWJClass)+1)\n")
	t := makeTemplate("UAX#14 range", templateRangeForClass)
	for key, codepoints := range codePointLists {
		w.WriteString(fmt.Sprintf("\n    // Range for UAX#14 class %s\n", key))
		w.WriteString(fmt.Sprintf("    %s = rangetable.New(", key))
		checkFatal(t.Execute(w, codepoints))
		w.WriteString(")\n")
		w.WriteString(fmt.Sprintf("    rangeFromUAX14Class[int(%sClass)] = %s\n", key, key))
	}
	w.WriteString("}\n")
}

func main() {
	doVerbose := flag.Bool("v", false, "verbose output mode")
	flag.Parse()
	verbose = *doVerbose
	codePointLists, err := loadUnicodeLineBreakFile()
	checkFatal(err)
	if verbose {
		logger.Printf("loaded %d UAX#14 breaking classes\n", len(codePointLists))
	}
	f, ioerr := os.Create("uax14classes.go")
	checkFatal(ioerr)
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString(header)
	w.WriteString(templateClassType)
	t := makeTemplate("UAX#14 classes", templateClassConsts)
	checkFatal(t.Execute(w, uax14classnames))
	t = makeTemplate("UAX#14 range tables", templateRangeTableVars)
	checkFatal(t.Execute(w, uax14classnames))
	t = makeTemplate("UAX#14 classes stringer", templateClassStringer)
	checkFatal(t.Execute(w, uax14classnames))
	generateRanges(w, codePointLists)
	w.Flush()
}

// --- Util -------------------------------------------------------------

// Little helper for testing
func timeTrack(start time.Time, name string) {
	if verbose {
		elapsed := time.Since(start)
		logger.Printf("timing: %s took %s\n", name, elapsed)
	}
}

func checkFatal(err error) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		logger.Fatalln(":", file, ":", line, "-", err)
	}
}
