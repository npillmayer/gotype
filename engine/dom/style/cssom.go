package style

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"github.com/aymerick/douceur/css"
	"github.com/npillmayer/gotype/core/config/tracing"
	"golang.org/x/net/html"
)

/*
TODO
- implement shortcut-properties: "border", "background" etc.
- extract style from header: link and sytle tags
- Locally scoped style (inline or <style> in body
- Specifity
- Fix API
- Document code
- Create Diagram -> Wiki
- API for export to DOT
*/

/* -----------------------------------------------------------------
BSD License

Copyright (c) 2017–18, Norbert Pillmayer

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
----------------------------------------------------------------- */

// RulesTree holds the styling rules of stylesheet.
//
// Status: Currently this is not really a tree.
// Optimize some day (see
// https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/).
type RulesTree struct {
	stylesheet *css.Stylesheet
	selectors  map[string]cascadia.Selector
}

type matchesList struct {
	matchingRules   []*css.Rule
	propertiesTable []specifity
}

type byHighestSpecifity []specifity

// make specifities sortable by highest sp.spec
func (sp byHighestSpecifity) Len() int           { return len(sp) }
func (sp byHighestSpecifity) Swap(i, j int)      { sp[i], sp[j] = sp[j], sp[i] }
func (sp byHighestSpecifity) Less(i, j int) bool { return sp[i].spec > sp[j].spec }

func (m *matchesList) String() string {
	s := fmt.Sprintf("match of %d rules:\n", len(m.matchingRules))
	s += "Src +-- Spec. --+------------- Key --------------+------- Value ---------------\n"
	for _, sp := range m.propertiesTable {
		s += fmt.Sprintf("%3d | %9d | %30s | %s\n", sp.source, sp.spec, sp.propertyKey, sp.propertyValue)
	}
	return s
}

type propertySource uint8

const (
	global propertySource = iota + 1
	author
	attribute
)

func NewRulesTree(css *css.Stylesheet) *RulesTree {
	return &RulesTree{css, nil}
}

func (rt *RulesTree) FilterMatchesFor(node *html.Node) *matchesList {
	list := &matchesList{}
	for _, rule := range rt.stylesheet.Rules {
		pre := rule.Prelude
		var sel cascadia.Selector
		found := false
		if sel, found = rt.selectors[pre]; !found {
			var err error
			sel, err = cascadia.Compile(pre)
			if err != nil {
				tracing.EngineTracer.Errorf("CSS selector seems not to work: %s", sel)
				break
			}
			if rt.selectors == nil {
				rt.selectors = make(map[string]cascadia.Selector)
			}
			rt.selectors[pre] = sel
		}
		if sel.Match(node) {
			list.matchingRules = append(list.matchingRules, rule)
		}
	}
	tracing.EngineTracer.Debugf("matching rules for %s", NodePath(node))
	return list
}

func NodePath(node *html.Node) string {
	s := ""
	if node.Type == html.TextNode {
		s += "(text)"
	} else if node.Type == html.ElementNode {
		s += fmt.Sprintf("- %s", node.Data)
	} else {
		s += "- (unknown)"
	}
	return s
}

func (matches *matchesList) SortProperties() *matchesList {
	var proptable []specifity
	for _, rule := range matches.matchingRules {
		for _, decl := range rule.Declarations {
			key := decl.Property
			value := Property(decl.Value)
			sp := specifity{author, rule, key, value, decl.Important, 0}
			sp.calcSpecifity()
			proptable = append(proptable, sp)
		}
	}
	if len(proptable) > 0 {
		sort.Sort(byHighestSpecifity(proptable))
		matches.propertiesTable = proptable
	}
	if tracing.EngineTracer.GetTraceLevel() >= tracing.LevelDebug {
		tracing.EngineTracer.Debugf(matches.String())
	}
	return matches
}

func (matches *matchesList) createStyleGroups(parent *StyledNode) map[string]*PropertyGroup {
	m := make(map[string]*PropertyGroup)
	done := make(map[string]bool)
	for _, pspec := range matches.propertiesTable { // for every specifity entry
		groupname, found := GroupNameFromPropertyKey[pspec.propertyKey]
		if !found {
			tracing.EngineTracer.Infof("Don't know about CSS property: %s", pspec.propertyKey)
			continue
		}
		if done[pspec.propertyKey] {
			// already present in current properties map
			// this must be from previous set with higher specifity
			// => do nothing
			break
		}
		_, pp := parent.findAncestorWithPropertyGroup(groupname) // must succeed
		if pp == nil {
			panic(fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", groupname))
		}
		group, isNew := pp.SpawnOn(pspec.propertyKey, pspec.propertyValue, true)
		if isNew { // a new property group has been created
			m[groupname] = group // put it into the group map
		}
		done[pspec.propertyKey] = true // remember we're done with this property
	}
	if len(m) == 0 { // no property groups created, not properties set
		return nil
	}
	if tracing.EngineTracer.GetTraceLevel() >= tracing.LevelDebug {
		for _, v := range m {
			tracing.EngineTracer.Debugf(v.String())
		}
	}
	return m
}

// --- Everthing in this section is a hack ! ----------------------------

// https://www.smashingmagazine.com/2007/07/css-specificity-things-you-should-know/

type specifity struct {
	source        propertySource
	rule          *css.Rule
	propertyKey   string
	propertyValue Property
	important     bool
	spec          uint16
}

// This is a grotesque hack :-(
func (sp *specifity) calcSpecifity() {
	if sp.important {
		sp.spec = 10000 // max
		return
	}
	sp.spec = uint16((sp.source - 1) * 100)
	selspec := len(sp.rule.Selectors)
	sp.spec += uint16(selspec)
	idcnt := countIdSelectors(sp.rule.Selectors)
	if idcnt > 0 {
		sp.spec += uint16(10) + idcnt // assumes no |selector| > 10 ...
	}
}

func countIdSelectors(sels []string) (cnt uint16) {
	for _, s := range sels {
		if strings.HasPrefix(s, "#") {
			cnt++
		}
	}
	return
}

func isImportant(rule *css.Rule) bool {
	for _, s := range rule.Selectors {
		if strings.EqualFold(s, "!important") {
			return true
		}
	}
	return false
}

// --- CSS Properties ---------------------------------------------------
//
// https://www.tutorialrepublic.com/css-reference/css3-properties.php
// https://www.mediaevent.de/xhtml/kernattribute.html

type Property string

func (p Property) IsInitial() bool {
	return p == "initial"
}

func (p Property) IsInherited() bool {
	return p == "inherited"
}

func (p Property) IsEmpty() bool {
	return p == ""
}

// --- CSS Property Groups ----------------------------------------------
//
// Caching is currently not implemented.

type PropertyGroup struct {
	name          string
	Parent        *PropertyGroup
	propertiesMap map[string]Property
	//signature       uint32 // signature of IDs and classes, used for caching
}

func NewPropertyGroup(ofType string) *PropertyGroup {
	pp := &PropertyGroup{}
	pp.name = ofType
	return pp
}

func (pp *PropertyGroup) String() string {
	s := "[" + pp.name + "] =\n"
	for k, v := range pp.propertiesMap {
		s += fmt.Sprintf("  %s = %s\n", k, v)
	}
	return s
}

func (pp *PropertyGroup) IsSet(key string) bool {
	if pp.propertiesMap == nil {
		return false
	}
	v, ok := pp.propertiesMap[key]
	return ok && !v.IsEmpty()
}

func (pp *PropertyGroup) Get(key string) Property {
	if pp.propertiesMap == nil {
		return ""
	}
	return pp.propertiesMap[key]
}

func (pp *PropertyGroup) Set(key string, p Property) {
	if pp.propertiesMap == nil {
		pp.propertiesMap = make(map[string]Property)
	}
	pp.propertiesMap[key] = p
}

func (pp *PropertyGroup) SpawnOn(key string, p Property, cascade bool) (*PropertyGroup, bool) {
	if cascade {
		found := pp.Cascade(key)
		if found.Get(key) == p {
			return pp, false
		}
	}
	npp := NewPropertyGroup(pp.name)
	//npp.signature = pp.signature
	npp.Set(key, p)
	return npp, true
}

func (pp *PropertyGroup) Cascade(key string) *PropertyGroup {
	it := pp
	for !it.IsSet(key) { // stopper is default partial
		it = it.Parent
	}
	return it
}

// Signature for being able to cache a fragment.
// Caching is currently not implemented.
//
// Class values have to be sorted.
//
// Returns hash and number of ID+class attributes.
func HashSignatureAttributes(htmlNode *html.Node) (uint32, uint8) {
	var hash uint32 = 0
	var count uint8 = 0
	signature := ""
	for _, a := range htmlNode.Attr {
		if a.Key == "id" {
			signature += a.Key
			count += 1
		}
	}
	for _, a := range htmlNode.Attr {
		if a.Key == "class" {
			signature += a.Key
			count += 1
		}
	}
	hash = xxhash.Checksum32([]byte(signature))
	return hash, count
}

var GroupNameFromPropertyKey = map[string]string{
	"margin-top":                 "Margins", // Margins
	"margin-left":                "Margins",
	"margin-right":               "Margins",
	"margin-bottom":              "Margins",
	"padding-top":                "Padding", // Padding
	"padding-left":               "Padding",
	"padding-right":              "Padding",
	"padding-bottom":             "Padding",
	"border-top-color":           "Border", // Border
	"border-left-color":          "Border",
	"border-right-color":         "Border",
	"border-bottom-color":        "Border",
	"border-top-width":           "Border",
	"border-left-width":          "Border",
	"border-right-width":         "Border",
	"border-bottom-width":        "Border",
	"border-top-style":           "Border",
	"border-left-style":          "Border",
	"border-right-style":         "Border",
	"border-bottom-style":        "Border",
	"border-top-left-radius":     "Border",
	"border-top-right-radius":    "Border",
	"border-bottom-left-radius":  "Border",
	"border-bottom-right-radius": "Border",
	"width":                      "Dimension", // Dimension
	"height":                     "Dimension",
	"min-width":                  "Dimension",
	"min-height":                 "Dimension",
	"max-width":                  "Dimension",
	"max-height":                 "Dimension",
}

func InitializeDefaultPropertyValues() map[string]*PropertyGroup {
	m := make(map[string]*PropertyGroup, 15)
	root := NewPropertyGroup("Root")

	margins := NewPropertyGroup("Margins")
	margins.Set("margin-top", "0")
	margins.Set("margin-left", "0")
	margins.Set("margin-right", "0")
	margins.Set("margin-bottom", "0")
	margins.Parent = root
	m["Margins"] = margins

	padding := NewPropertyGroup("Padding")
	padding.Set("padding-top", "0")
	padding.Set("padding-left", "0")
	padding.Set("padding-right", "0")
	padding.Set("padding-bottom", "0")
	padding.Parent = root
	m["Padding"] = padding

	border := NewPropertyGroup("Border")
	border.Set("border-top-color", "black")
	border.Set("border-left-color", "black")
	border.Set("border-right-color", "black")
	border.Set("border-bottom-color", "black")
	border.Set("border-top-width", "medium")
	border.Set("border-left-width", "medium")
	border.Set("border-right-width", "medium")
	border.Set("border-bottom-width", "medium")
	border.Set("border-top-style", "solid")
	border.Set("border-left-style", "solid")
	border.Set("border-right-style", "solid")
	border.Set("border-bottom-style", "solid")
	border.Set("border-top-left-radius", "0")
	border.Set("border-top-right-radius", "0")
	border.Set("border-bottom-left-radius", "0")
	border.Set("border-bottom-right-radius", "0")
	border.Parent = root
	m["Border"] = border

	dimension := NewPropertyGroup("Dimension")
	dimension.Set("width", "10%")
	dimension.Set("width", "100pt")
	dimension.Set("min-width", "0")
	dimension.Set("min-height", "0")
	dimension.Set("max-width", "10000pt")
	dimension.Set("max-height", "10000pt")
	dimension.Parent = root
	m["Dimension"] = dimension

	/*
	   type ColorModel string

	   type Color struct {
	   	Color   color.Color
	   	Model   ColorModel
	   	Opacity uint8
	   }

	   type DisplayStyle struct {
	   	Display    uint8 // https://www.tutorialrepublic.com/css-reference/css-display-property.php
	   	Position   uint8
	   	Top        dimen.Dimen
	   	Left       dimen.Dimen
	   	Right      dimen.Dimen
	   	Bottom     dimen.Dimen
	   	Float      uint8
	   	ZIndex     int
	   	Overflow   uint8
	   	OverflowX  uint8
	   	OverflowY  uint8
	   	Clip       string // geometric shape
	   	Visibility bool
	   }

	   type Background struct {
	   	Color color.Color
	   	//Position TODO
	   	Image  image.Image
	   	Origin dimen.Point
	   	Size   dimen.Point
	   	Clip   uint8
	   }

	   type Font struct {
	   	Family     string
	   	Style      string
	   	Variant    uint16
	   	Stretch    uint8
	   	Size       dimen.Dimen
	   	SizeAdjust dimen.Dimen
	   }

	   type TextProperties struct {
	   	Direction          uint8
	   	WordSpacing        uint8
	   	LetterSpacing      uint8
	   	VerticalAlignment  uint8
	   	TextAlignment      uint8 // + TextJustify
	   	TextAlignLast      uint8
	   	TextIndentation    dimen.Dimen // first line
	   	TabSize            dimen.Dimen
	   	LineHeight         uint8
	   	TextDecoration     uint8
	   	TextTransformation uint8
	   	WordWrap           uint8
	   	WordBreak          uint8
	   	Whitespace         uint8
	   	TextOverflow       uint8
	   }


	   type GeneratedContent struct {
	   	Content          string
	   	Quotes           string
	   	CounterReset     uint8
	   	CounterIncrement uint8
	   }

	   type Print struct {
	   	PageBreakAfter  uint8
	   	PageBreakBefore uint8
	   	PageBreakInside uint8
	   }

	   type Outline struct {
	   	Color  color.Color
	   	Offset dimen.Dimen
	   	Style  uint8
	   	Width  dimen.Dimen
	   }

	   //list-style-type:
	   //	disc | circle | square | decimal | decimal-leading-zero | lower-roman |
	   //  upper-roman | lower-greek | lower-latin | upper-latin | armenian |
	   //  georgian | lower-alpha | upper-alpha | none | initial | inherit
	   type List struct {
	   	StyleImage    image.Image
	   	StylePosition uint8 // inside, outside
	   	StyleType     uint8
	   }

	*/

	return m
}

// --- Styled Node Tree -------------------------------------------------

// StyledNodes are the building blocks of the styled tree.
type StyledNode struct {
	node           *html.Node
	computedStyles map[string]*PropertyGroup
	parent         *StyledNode
	children       []*StyledNode
}

func (sn *StyledNode) findAncestorWithPropertyGroup(group string) (*StyledNode, *PropertyGroup) {
	var pp *PropertyGroup
	it := sn
	for it != nil && pp == nil {
		pp = it.computedStyles[group]
		it = it.parent
	}
	return it, pp
}

func (sn *StyledNode) getProperty(group string, key string) Property {
	_, pp := sn.findAncestorWithPropertyGroup(group) // must succeed
	if pp == nil {
		panic(fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", group))
	}
	return pp.Cascade(key).Get(key) // must succeed
}

type StyledNodeTree struct {
	root          *StyledNode
	defaultStyles map[string]*PropertyGroup
}

// For an explanation what's going on here, refer to
// https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/
// and
// https://limpet.net/mbrubeck/2014/08/23/toy-layout-engine-4-style.html
func ConstructStyledNodeTree(dom *html.Node, rules *RulesTree) (*StyledNodeTree, error) {
	tree := &StyledNodeTree{}
	tree.defaultStyles = InitializeDefaultPropertyValues()
	querydom := goquery.NewDocumentFromNode(dom)
	body := querydom.Find("body").First()
	if body == nil {
		return nil, errors.New("Misconstructed DOM: cannot find <body>")
	}
	viewport := &StyledNode{}
	viewport.node = body.Nodes[0]
	viewport.computedStyles = tree.defaultStyles
	// now recurse
	// first version is without concurrency
	worklist := body.Children().Nodes
	for _, n := range worklist {
		recurseStyling(viewport, n, rules)
	}
	return tree, nil
}

// node is a (possibly indirect) child of the node corresponding to sn.
func recurseStyling(parent *StyledNode, node *html.Node, rules *RulesTree) {
	if parent == nil || node == nil {
		return // recursion termination
	}
	sn := &StyledNode{}
	sn.node = node
	sn.parent = sn
	parent.children = append(parent.children, sn)
	matchingRules := rules.FilterMatchesFor(node)
	if matchingRules != nil {
		propertiesTable := matchingRules.SortProperties()
		groups := propertiesTable.createStyleGroups(parent)
		sn.computedStyles = groups // may be nil
	}
	node = node.FirstChild // now recurse into children
	for node != nil {
		if node.Type == html.ElementNode {
			recurseStyling(sn, node, rules)
		}
		node = node.NextSibling
	}
}
