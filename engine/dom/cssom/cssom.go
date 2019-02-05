package cssom

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/andybalholm/cascadia"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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

// We trace to EngineTracer.
func T() tracing.Trace {
	return tracing.EngineTracer
}

// CSSOM is the "CSS Object Model", similar to the DOM for HTML.
// Our CSSOM consists of a set of stylesheets, each relevant for a sub-tree
// of the DOM. This sub-tree is called the "scope" of the stylesheet.
// Sub-trees are identified through the top node.
//
// Stylesheets are wrapped into an internal rules tree.
type CSSOM struct {
	rulesTree         *rulesTreeType               // style sheets
	defaultProperties *style.PropertyMap           // "user agent" style properties
	compoundSplitters []CompoundPropertiesSplitter // split up compound properties
}

// NewCSSOM creates an empty CSSOM.
// Clients are allowed to supply a map of additional/custom CSS property values.
// These may override values of the default ("user-agent") style sheet,
// or introduce completely new styling properties.
func NewCSSOM(additionalProperties []style.KeyValue) CSSOM {
	cssom := CSSOM{}
	cssom.rulesTree = newRulesTree()
	cssom.defaultProperties = style.InitializeDefaultPropertyValues(additionalProperties)
	cssom.compoundSplitters = make([]CompoundPropertiesSplitter, 1)
	cssom.compoundSplitters[0] = style.SplitCompoundProperty
	return cssom
}

// AddStylesForScope includes a stylesheet to a CSSOM and sets the scope for
// the stylesheet. If a stylesheet for the scope already exists, the
// styles are merged. css may be nil. If scope is nil then scope is the
// root (i.e., top-level content element) of a future document.
//
// The stylse sheet may not be nil.
// source hints to where the stylesheet comes from.
// Its value will affect the calculation of specifity for rules of this
// stylesheet.
func (cssom CSSOM) AddStylesForScope(scope *html.Node, css StyleSheet, source PropertySource) error {
	if scope != nil && scope.Type != html.ElementNode {
		return errors.New("Can style element nodes only")
	}
	if css == nil {
		return errors.New("Style sheet is nil")
	}
	cssom.rulesTree.StoreStylesheetForHtmlNode(scope, css, source)
	return nil
}

// --- A rules tree -----------------------------------------------------

// RulesTree holds the styling rules of a stylesheet.
//
// Status: Currently this is not really a tree.
// Optimize some day (see
// https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/).
type rulesTreeType struct {
	stylesheets *sync.Map                    // of type []stylesheetType
	selectors   map[string]cascadia.Selector // cache of compiled selectors
	source      PropertySource               // where do these rules come from?
}

// ad-hoc container type for stylesheets and their origin.
// To be stored in a map (per HTML node).
type stylesheetType struct {
	stylesheet StyleSheet
	source     PropertySource
}

func newRulesTree() *rulesTreeType {
	rt := &rulesTreeType{}
	rt.stylesheets = &sync.Map{}
	rt.selectors = make(map[string]cascadia.Selector)
	return rt
}

// StylesheetsForHtmlNode retrieves all style sheets registered for
// an html node. If h is nil it is interpreted as the root scope.
func (rt rulesTreeType) StylesheetsForHtmlNode(h *html.Node) []stylesheetType {
	if h == nil {
		h = rootElement
	}
	sheets, found := rt.stylesheets.Load(h)
	if !found {
		return nil
	}
	return sheets.([]stylesheetType)
}

// StoreStylesheetForHtmlNode registers a style sheet for
// an html node. If h is nil it is interpreted as the root scope.
func (rt rulesTreeType) StoreStylesheetForHtmlNode(h *html.Node, sheet StyleSheet,
	source PropertySource) {
	//
	if h == nil {
		h = rootElement
	}
	sheets := rt.StylesheetsForHtmlNode(h)
	if sheets == nil {
		T().Debugf("Adding first style sheet for HTML node %v", h)
		rt.stylesheets.Store(h, []stylesheetType{stylesheetType{sheet, source}})
	} else {
		sheets = append(sheets, stylesheetType{sheet, source})
		rt.stylesheets.Store(h, sheets)
	}
}

// Empty is a predicate wether a rulestree is empty, i.e. does not contain
// any rules.
func (rt *rulesTreeType) Empty() bool {
	if rt == nil {
		return true
	}
	csscnt := 0
	rt.stylesheets.Range(func(interface{}, interface{}) bool {
		csscnt++
		return true
	})
	T().Debugf("Counting style sheet entries in rules tree: %d", csscnt)
	return csscnt == 0
}

// Compound properties are properties which abbreviate the
// setting of more fine grained propertes. An example is
//    padding: 10px 20px
// which sets the following detail properties:
//    padding-top:    10px
//    padding-right:  20px
//    padding-bottom: 10px
//    padding-left:   20px
// Standard CSS compound properties are known by default, but clients are
// allowed to extend the set of compound properties.
type CompoundPropertiesSplitter func(string, style.Property) ([]style.KeyValue, error)

// RegisterCompoundSplitter allows clients to handle additional compound
// properties. See type CompoundPropertiesSplitter.
func (cssom CSSOM) RegisterCompoundSplitter(splitter CompoundPropertiesSplitter) {
	if splitter != nil {
		cssom.compoundSplitters = append(cssom.compoundSplitters, splitter)
	}
}

// --- Style Rule Matching ----------------------------------------------

// Properties may be defined at different places in HTML: as a sytlesheet
// reference link, within a <script> element in the HTML file, or in an
// attribute value.
//
// PropertySource affects the specifity of rules.
type PropertySource uint8

// Values for property sources, used when adding style sheets.
const (
	Global    PropertySource = iota + 1 // "browser" globals
	Author                              // CSS author (stylesheet link)
	Script                              // <script> element
	Attribute                           // in an element's attribute(s)
)

// rootElement is a symbolic node to denote the body element of a future
// HTML document. AddStylesFor(...) with nil as a scope will replace it
// with this marker for scoping the complete document body.
var rootElement *html.Node = &html.Node{Data: "root"}

// Internal helper for applying rules to an HTML node.
// In a first step it holds all the rules matching for an HTML node.
// In a second step it collects all the properties set in those rules,
// then orderes them by specifity.
type matchesList struct {
	matchingRules   []Rule
	propertiesTable []propertyPlusSpecifityType
}

// Rule-matchings are collected from more than one stylesheet. Matching
// rules from these stylesheets will be merged to one list.
func (mlist *matchesList) mergeMatchesWith(m *matchesList) *matchesList {
	if mlist == nil {
		return m
	}
	if m != nil {
		for _, r := range m.matchingRules {
			mlist.matchingRules = append(mlist.matchingRules, r)
		}
	}
	return mlist
}

// Rule-matchings have to be sorted by specifity. We'll sort the highest
// specifity up and won't overwrite earlier matches with later matches.

// sorter
type byHighestSpecifity []propertyPlusSpecifityType

// make specifities sortable by highest sp.spec
func (sp byHighestSpecifity) Len() int           { return len(sp) }
func (sp byHighestSpecifity) Swap(i, j int)      { sp[i], sp[j] = sp[j], sp[i] }
func (sp byHighestSpecifity) Less(i, j int) bool { return sp[i].spec > sp[j].spec }

// This is a small helper to print out a table with rule-matches for a node.
func (m *matchesList) String() string {
	s := fmt.Sprintf("match of %d rules:\n", len(m.matchingRules))
	s += "Src +-- Spec. --+------------- Key --------------+------- Value ---------------\n"
	for _, sp := range m.propertiesTable {
		s += fmt.Sprintf("%3d | %9d | %30s | %s\n", sp.source, sp.spec, sp.propertyKey, sp.propertyValue)
	}
	return s
}

// FilterMatchesFor(node) iterates through all the rules relevant at this
// point and looks for rules matching the current HTML node h.
// The heavy lifting is done by cascadia. We have to 'compile' all rules
// and will cache compiled rules.
//
// Will return a slice of CSS rules matched for h.
func (rt *rulesTreeType) FilterMatchesFor(h *html.Node) *matchesList {
	//list := &matchesList{}
	matchingRules := make([]Rule, 0, 3)
	sheets := rt.StylesheetsForHtmlNode(rootElement)
	for _, s := range sheets {
		for _, rule := range s.stylesheet.Rules() {
			if rt.matchRuleForHtmlNode(h, rule) {
				matchingRules = append(matchingRules, rule)
			}
		}
	}
	sheets = rt.StylesheetsForHtmlNode(h)
	for _, s := range sheets {
		for _, rule := range s.stylesheet.Rules() {
			if rt.matchRuleForHtmlNode(h, rule) {
				matchingRules = append(matchingRules, rule)
			}
		}
	}
	return &matchesList{matchingRules, nil}
}

func (rt *rulesTreeType) matchRuleForHtmlNode(h *html.Node, rule Rule) bool {
	selectorString := rule.Selector()
	if selectorString == "" { // style-attribute local for this HTML node
		//matchingRules = append(matchingRules, rule)
		return true
	} else { // try to match selector for this rule against HTML node
		var sel cascadia.Selector
		found := false
		if sel, found = rt.selectors[selectorString]; !found {
			var err error
			sel, err = cascadia.Compile(selectorString)
			if err != nil {
				T().Errorf("CSS selector seems not to work: %s", selectorString)
				return false
			}
			rt.selectors[selectorString] = sel
		}
		if sel.Match(h) {
			//list.matchingRules = append(list.matchingRules, rule)
			return true
		}
	}
	return false
}

// SortProperties takes a slice of CSS rules (matched for an HTML node) and
// extracts all the properties set within the rules. These properties are
// then split into atomic properties, if they are compound properties
// (e.g.,
//     "margin" ⟹ "margin-top", "margin-right", ...
// Finally all property entries are sorted by specifity of the enclosing rule.
func (matches *matchesList) SortProperties(splitters []CompoundPropertiesSplitter) {
	var proptable []propertyPlusSpecifityType
	for rno, rule := range matches.matchingRules {
		for _, propertyKey := range rule.Properties() {
			value := style.Property(rule.Value(propertyKey))
			props, err := splitCompoundProperty(splitters, propertyKey, value)
			if err != nil {
				sp := propertyPlusSpecifityType{Author, rule, propertyKey, value, rule.IsImportant(propertyKey), 0}
				sp.calcSpecifity(rno)
				proptable = append(proptable, sp)
			} else {
				T().Debugf("%s is a compound style", propertyKey)
				for _, kv := range props {
					key := kv.Key
					val := kv.Value
					sp := propertyPlusSpecifityType{Author, rule, key, val, rule.IsImportant(propertyKey), 0}
					sp.calcSpecifity(rno)
					proptable = append(proptable, sp)
				}
			}
		}
	}
	if len(proptable) > 0 {
		sort.Sort(byHighestSpecifity(proptable))
		matches.propertiesTable = proptable
	}
	if T().GetTraceLevel() >= tracing.LevelDebug {
		T().Debugf(matches.String())
	}
}

// --- Specifity of rules -----------------------------------------------

type propertyPlusSpecifityType struct {
	source        PropertySource // where the property has been defined
	rule          Rule           // the rule containing the property definition
	propertyKey   string         // CSS property name
	propertyValue style.Property // raw string value
	important     bool           // marked as !IMPORTANT ?
	spec          uint32         // specifity value to calculate; higher is more
}

// CalcSpecifity calculates an appromiation to the true W3C specifity.
// https://www.smashingmagazine.com/2007/07/css-specificity-things-you-should-know/
//
// no is a sequence number for rules, ensuring that later rules override
// previously defined rules / properties.
func (sp *propertyPlusSpecifityType) calcSpecifity(no int) {
	if sp.rule.IsImportant(sp.propertyKey) {
		sp.spec = 99999 // max
		return
	}
	sp.spec = uint32(sp.source-1) * 1000
	selectorstring := sp.rule.Selector()
	// simple "parsing" = rough estimate...
	// alternatively use code from cascadia or from
	// https://godoc.org/github.com/ericchiang/css
	sels := strings.Fields(selectorstring)
	var selcnt uint32
	var idcnt uint32
	var classcnt uint32
	for _, sel := range sels {
		selcnt++
		if strings.ContainsRune(sel, ':') {
			selcnt++ // count double
		}
		if strings.ContainsAny(sel, ".[:") {
			classcnt++
		}
		if strings.HasPrefix(sel, "#") {
			idcnt++
		}
	}
	sp.spec += selcnt*10 + classcnt*100 + idcnt*1000 + uint32(no)
}

// --- Style Property Groups --------------------------------------------

func (matches *matchesList) createStyleGroups(parent *tree.Node,
	creator style.Creator) *style.PropertyMap {
	//
	pmap := style.NewPropertyMap()
	done := make(map[string]bool, len(matches.propertiesTable))
	for _, pspec := range matches.propertiesTable { // for every specifity entry
		if done[pspec.propertyKey] {
			// already present in current properties map
			// this must be from previous set with higher specifity
			// => do nothing
			break
		}
		groupname := style.GroupNameFromPropertyKey(pspec.propertyKey)
		group := pmap.Group(groupname)
		if group != nil {
			group.Set(pspec.propertyKey, pspec.propertyValue)
		} else {
			T().Infof("parent is %s, searching for prop group %s", parent, groupname)
			_, pg := findAncestorWithPropertyGroup(parent, groupname, creator) // must succeed
			if pg == nil {
				panic(fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", groupname))
			}
			group, isNew := pg.ForkOnProperty(pspec.propertyKey, pspec.propertyValue, true)
			if isNew { // a new property group has been created
				pmap = pmap.AddAllFromGroup(group, true) // put it into the group map
			}
		}
		done[pspec.propertyKey] = true // remember we're done with this property
	}
	if pmap.Size() == 0 { // no property groups created, no properties set
		return nil
	}
	// if tracing.EngineTracer.GetTraceLevel() >= tracing.LevelDebug {
	// 	for _, v := range m {
	// 		tracing.EngineTracer.Debugf(v.String())
	// 	}
	// }
	return pmap
}

// --- Styled Node Tree -------------------------------------------------

//func setupStyledNodeTree(domRoot *html.Node, defaults *style.PropertyMap, builder StyledTreeBuilder) StyledNode {
func setupStyledNodeTree(domRoot *html.Node, defaults *style.PropertyMap,
	creator style.Creator) *tree.Node {
	//
	rootNode := creator.StyleForHtmlNode(domRoot)
	creator.SetComputedStyles(rootNode, defaults)
	//T().Debugf("UA node has styles = %s", creator.ToStyler(rootNode).ComputedStyles())
	docNode := creator.StyleForHtmlNode(domRoot)
	rootNode.AddChild(docNode)
	return docNode
}

//func findAncestorWithPropertyGroup(sn StyledNode, group string, builder StyledTreeBuilder) (StyledNode, *style.PropertyGroup) {
func findAncestorWithPropertyGroup(sn *tree.Node, group string,
	creator style.Creator) (*tree.Node, *style.PropertyGroup) {
	//
	var pg *style.PropertyGroup
	if sn == nil {
		T().Errorf("Search for ancestor with property group %s started with nil", group)
		return nil, nil
	}
	it := sn // start search at styled node itself, then proceed upwards
	last := sn
	for it != nil && pg == nil {
		styles := creator.ToStyler(it).ComputedStyles()
		if styles != nil {
			pg = styles.Group(group)
		}
		it = it.Parent()
		if it != nil {
			last = it
		}
	}
	// if it == nil {
	// 	T().Debugf("At root of tree searching for property group %s", group)
	// 	if pg == nil {
	// 		T().Errorf("Property group %s not found", group)
	// 		T().Debugf("Property map of last node %v =\n%s", last, creator.ToStyler(last).ComputedStyles())
	// 	}
	// }
	return last, pg
}

// Style() gets things rolling. It styles a DOM, referred to by the root
// node, and returns a tree of styled nodes.
// For an explanation what's going on here, refer to
// https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/
// and
// https://limpet.net/mbrubeck/2014/08/23/toy-layout-engine-4-style.html
//
// If either dom or factory are nil, no tree is returned (but an error).
func (cssom CSSOM) Style(dom *html.Node, creator style.Creator) (*tree.Node, error) {
	if dom == nil {
		return nil, errors.New("Nothing to style: empty document")
	}
	if creator == nil {
		return nil, errors.New("Cannot style: no builder to create styles nodes")
	}
	if cssom.rulesTree.Empty() {
		T().Infof("Styling HTML tree without having any CSS rules")
	}
	T().Debugf("--- Creating style nodes for HTML nodes ----")
	styledRootNode := setupStyledNodeTree(dom, cssom.defaultProperties, creator)
	walker := tree.NewWalker(styledRootNode) // create a concurrent tree walker
	createNodes := func(node *tree.Node, parent *tree.Node, pos int) (*tree.Node, error) {
		return createStyledChildren(node, cssom.rulesTree, creator) // provide closure with style creator
	}
	future := walker.TopDown(createNodes).Promise() // build the style tree
	if _, err := future(); err != nil {
		T().Errorf("Error while creating styled tree: %v", err)
		return nil, err
	}
	T().Debugf("--- Now styling newly created nodes --------")
	walker = tree.NewWalker(styledRootNode)
	createStyles := func(node *tree.Node, parent *tree.Node, pos int) (*tree.Node, error) {
		return createStylesForNode(node, cssom.rulesTree, creator, cssom.compoundSplitters)
	}
	future = walker.TopDown(createStyles).Promise() // build the style tree
	if _, err := future(); err != nil {
		T().Errorf("Error while creating style properties: %v", err)
		return nil, err
	}
	return styledRootNode, nil
}

// Pre-condition: sn has been styled and points to an HTML node.
// Now iterate through the HTML children and create styled nodes for each.
func createStyledChildren(parent *tree.Node, rulesTree *rulesTreeType,
	creator style.Creator) (*tree.Node, error) {
	//
	domnode := dom.NewRONode(parent, creator.ToStyler) // interface RODomNode
	T().Debugf("Input node = %v, creating styled children", domnode)
	h := domnode.HtmlNode()
	if h.Type == html.ElementNode || h.Type == html.DocumentNode {
		ch := h.FirstChild
		for ch != nil {
			T().Debugf("ch.DataAtom = %v  ==================", ch.DataAtom)
			if ch.DataAtom == atom.Style { // <style> element
				T().Debugf("TODO: <style> node has Data = %v", h.Data)
				if ch.FirstChild.Type == html.TextNode {
					T().Debugf("first child is text: %v", ch.FirstChild.Data)
				}
				// TODO attach styles from <style>
			} else if isStylable(ch.DataAtom) {
				sn := creator.StyleForHtmlNode(ch)
				parent.AddChild(sn) // sn will be sent to next pipeline stage
				if styleAttr := getStyleAttribute(ch); styleAttr != nil {
					// attach local style attributes
					rulesTree.StoreStylesheetForHtmlNode(ch, styleAttr, Attribute)
				}
			} else if ch.DataAtom == atom.Head {
				// search for <style> element inside <head>
				style := findStyleElements(ch)
				if len(style) > 0 {
					for _, s := range style {
						rulesTree.StoreStylesheetForHtmlNode(nil, s, Script)
					}
				}
			}
			ch = ch.NextSibling
		}
	} else if h.Type == html.TextNode {
		// do not send text node to next pipeline stage
		return nil, nil
	}
	return parent, nil
}

func isStylable(a atom.Atom) bool {
	switch a {
	case atom.A, atom.Address, atom.Acronym, atom.Article, atom.Aside,
		atom.B, atom.Blink, atom.Blockquote, atom.Body, atom.Br,
		atom.Button, atom.Label, atom.Canvas, atom.Caption,
		atom.Code, atom.Content, atom.Div, atom.Em, atom.Figcaption,
		atom.Figure, atom.Footer, atom.Form, atom.Frame, atom.Hr,
		atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Html,
		atom.I, atom.Img, atom.Input, atom.Li, atom.Main, atom.Math,
		atom.Menu, atom.Menuitem, atom.Nav, atom.Ol, atom.Option,
		atom.P, atom.Picture, atom.Pre, atom.Poster, atom.Q, atom.S,
		atom.Section, atom.Span, atom.Spacer, atom.Strong, atom.Summary,
		atom.Svg, atom.Sup, atom.Table, atom.Td, atom.Tr, atom.Th,
		atom.Textarea, atom.Tfoot, atom.Title, atom.Ul, atom.Video:
		return true
	}
	return false
}

func findStyleElements(n *html.Node) []StyleSheet {
	ch := n.FirstChild
	for ch != nil {
		if ch.DataAtom == atom.Style {
			// TODO create a style sheet from <style> element
			// how to do this without a concrete CSS parser implementation?
		}
		ch = ch.NextSibling
	}
	return nil
}

func createStylesForNode(node *tree.Node, rulesTree *rulesTreeType, creator style.Creator,
	splitters []CompoundPropertiesSplitter) (*tree.Node, error) {
	//
	styler := creator.ToStyler(node)
	matchlist := rulesTree.FilterMatchesFor(styler.HtmlNode())
	if matchlist != nil && len(matchlist.matchingRules) != 0 {
		matchlist.SortProperties(splitters)
		pmap := matchlist.createStyleGroups(node.Parent(), creator)
		T().Debugf("Setting styles for node %v =\n%s", node, pmap)
		creator.SetComputedStyles(node, pmap)
	} else {
		T().Debugf("Node %v matched no style rules", node)
	}
	return node, nil
}

// Scopes for all ruletrees have to be HTML element nodes.
// This is enforced during construction of the rules trees (adding to CSSDOM).
// func (cssom CSSOM) attachStylesheets(dom *html.Node) error {
// 	for scope, rulestree := range cssom.rules {
// 		var stylingRootElement *html.Node
// 		if scope == rootElement { // scope is body element, i.e. whole document
// 			stylingRootElement := findBodyElement(dom)
// 			if stylingRootElement == nil { // no body element found
// 				T().Infof("Misconstructed DOM: cannot find <body>. Proceeding.")
// 				stylingRootElement = dom // root of fragment, try our best
// 			}
// 			// runner.activeStylers[stylingRootElement] = rulestree
// 			// runner.startNode = stylingRootElement
// 		} else {
// 			stylingRootElement = findThisNode(dom, scope)
// 			if stylingRootElement == nil { // scope not not found in DOM
// 				// do not try to proceed: client meant something else...
// 				return errors.New(fmt.Sprintf("Scope '%s' not found in DOM", scope.Data))
// 			}
// 			if stylingRootElement.Type != html.ElementNode {
// 				return errors.New(fmt.Sprintf("Scope '%s' is not of element type", shorten(scope.Data)))
// 			}
// 			//runner.inactiveStylers[stylingRootElement] = rulestree
// 		}
// 	}
// 	// if runner.startNode == nil {
// 	// 	runner.startNode = dom
// 	// }
// 	return nil
// }

// stylingRunner is a helper type for concurrent execution of node styling.
// type stylingRunner struct {
// 	sync.RWMutex                                 // used to protect non-threadsafe code
// 	activeStylers   map[*html.Node]rulesTree     // rules in scope
// 	inactiveStylers map[*html.Node]rulesTree     // rules out of scope
// 	startNode       *html.Node                   // usually document body
// 	builder         StyledTreeBuilder            // the tree builder to use
// 	splitters       []CompoundPropertiesSplitter // splitters for compound style properties
// 	workers         *workergroup                 // concurrent workers
// }

// Create a new stylingRunner.
// func newStylingRunner(builder StyledTreeBuilder, splitters []CompoundPropertiesSplitter) *stylingRunner {
// 	runner := &stylingRunner{}
// 	runner.activeStylers = make(map[*html.Node]rulesTree)
// 	runner.inactiveStylers = make(map[*html.Node]rulesTree)
// 	runner.builder = builder
// 	runner.splitters = splitters
// 	return runner
// }

// Check for a scope (= DOM node) whether new style sheets have to be
// activated.
//
// Concurrency-safe.
// func (runner *stylingRunner) activateStylesheetsFor(node *html.Node) {
// 	var sty *rulesTree
// 	runner.RLock()
// 	for scope, r := range runner.inactiveStylers {
// 		if scope == node {
// 			sty = &r
// 			break
// 		}
// 	}
// 	runner.RUnlock()
// 	if sty != nil {
// 		runner.Lock()
// 		defer runner.Unlock()
// 		runner.activeStylers[node] = *sty
// 		delete(runner.inactiveStylers, node)
// 	}
// }

// Currently starts 3 workers. TODO: make this configurable.
// func (runner *stylingRunner) doStyle(node *html.Node, parent *tree.Node) {
// 	workload := make(chan workPackage)
// 	errors := make(chan error)
// 	workers := launch(3).workers(styleSingleNode, workload, errors)
// 	runner.workers = workers
// 	initialWorkPackage := workPackage{
// 		runner:       runner,
// 		node:         node,
// 		styledParent: parent,
// 	}
// 	watch(workers, initialWorkPackage)
// 	e := collect(errors) // wait for workers to complete
// 	T.Infof("Errors from styling workers: %v", e)
// 	//
// 	/*  // non-concurrent version
// 	runner.activateStylesheetsFor(node)
// 	if createsStyledNode(node.Type) {
// 		sn := builder.MakeNodeFor(node)
// 		builder.LinkNodeToParent(sn, parent)
// 		var matchingRules *matchesList
// 		for _, rulesTree := range runner.activeStylers {
// 			matches := rulesTree.FilterMatchesFor(node)
// 			matchingRules = matchingRules.mergeMatchesWith(matches)
// 		}
// 		if matchingRules != nil {
// 			matchingRules.SortProperties(runner.splitters)
// 			pmap := matchingRules.createStyleGroups(parent, builder)
// 			sn.SetComputedStyles(pmap)
// 		}
// 		parent = sn // continue with newly created styled node
// 	}
// 	node = node.FirstChild // now recurse into children
// 	for node != nil {
// 		if node.Type == html.ElementNode {
// 			runner.doStyle(node, parent, builder) // TODO make this concurrent
// 		}
// 		node = node.NextSibling
// 	}
// 	*/
// }

// workPackage is a type to distribute styling work amongst worker goroutines.
// It contains necessary information to style a single DOM node.
// type workPackage struct {
// 	runner       *stylingRunner // styling controller
// 	node         *html.Node     // the DOM node to style
// 	styledParent StyledNode     // the active styled node above the DOM node
// 	udata        interface{}    // all-purpose extension data
// }

// styleSingleNode is of type workerTask.
// Workers will perform the styling of a single HTML DOM node. This task
// will create a styled node for the DOM node. The newly created styled node
// will be attached to the existing styled parent node.
//
// The worker can rely on being the only worker tasked with styling this
// DOM node. However, linking the new styled node to its parent must be
// a concurrency-safe operation.
/*
func styleSingleNode() error {
	//func styleSingleNode(wp workPackage) error {
	//T.Debugf("worker is grabbing node %s", wp.node.Data)
	var lasterror error
	builder := wp.runner.builder
	wp.runner.activateStylesheetsFor(wp.node) // is concurrency-safe
	if createsStyledNode(wp.node.Type) {
		sn := builder.MakeNodeFor(wp.node)
		builder.LinkNodeToParent(sn, wp.styledParent) // must be concurrency-safe
		pmap := calcComputedStylesForNode(wp.node, sn, wp.runner)
		if pmap != nil {
			sn.SetComputedStyles(pmap)
		}
		parent := sn // continue with newly created styled node as parent for next node
		createWorkPackagesForChildrenNodes(wp.node, parent, wp.runner)
	}
	//T.Debugf("worpackage for node %s done", wp.node.Data)
	return lasterror
}
*/

// func calcComputedStylesForNode(node *html.Node, parent *tree.Node) *style.PropertyMap {
// 	var matchingRules *matchesList
// 	for _, rulesTree := range runner.activeStylers {
// 		matches := rulesTree.FilterMatchesFor(node)
// 		matchingRules = matchingRules.mergeMatchesWith(matches)
// 	}
// 	if matchingRules != nil {
// 		matchingRules.SortProperties(runner.splitters)
// 		pmap := matchingRules.createStyleGroups(parent, runner.builder)
// 		return pmap
// 	}
// 	return nil
// }

// node is exclusive for a worker; no other goroutine can access its children.
// func createWorkPackagesForChildrenNodes(node *html.Node, parent *tree.Node) {
// 	//
// 	T().Debugf("creating worpackage for child-nodes of  %s", node.Data)
// 	node = node.FirstChild // recurse into children
// 	for node != nil {
// 		if node.Type == html.ElementNode {
// 			wp := workPackage{runner: runner}
// 			wp.node = node
// 			wp.styledParent = parent
// 			order(runner.workers, wp)
// 		}
// 		node = node.NextSibling
// 	}
// }

// --- Helpers ----------------------------------------------------------

var errNoSuchCompoundProperty error = errors.New("No such compound property")

// Try to split up a property (which may or may not be a compound
// property) using a set of splitter functions.
// Return a slice of key-value pairs or nil.
func splitCompoundProperty(splitters []CompoundPropertiesSplitter,
	key string, value style.Property) ([]style.KeyValue, error) {
	for _, splitter := range splitters {
		kv, err := splitter(key, value)
		if err != nil {
			return kv, err
		}
	}
	return nil, errNoSuchCompoundProperty
}

// Which HTML node type needs a corresponding styled node?
func createsStyledNode(nodeType html.NodeType) bool {
	if nodeType == html.ElementNode || nodeType == html.TextNode {
		return true
	}
	return false
}

// Helper to find nodes matching a predicate. Currently works recursive.
// Returns a node or nil.
func findNode(node *html.Node, matcher func(n *html.Node) bool) *html.Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if f := findNode(c, matcher); f != nil {
			return f
		}
	}
	return nil
}

func findThisNode(tree *html.Node, nodeToFind *html.Node) *html.Node {
	return findNode(tree, func(n *html.Node) bool {
		return n == nodeToFind
	})
}

func findBodyElement(tree *html.Node) *html.Node {
	return findNode(tree, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "body"
	})
}

// shorten a string
func shorten(s string) string {
	if len(s) > 10 {
		return s[:10] + "..."
	}
	return s
}

// small helper to debug-print out a node. TODO
func nodePath(node *html.Node) string {
	s := ""
	if node.Type == html.TextNode {
		s += "(text)"
	} else if node.Type == html.ElementNode {
		s += fmt.Sprintf("%s", node.Data)
	} else {
		s += "(unknown)"
	}
	return s
}

// --- Local pseudo rules for style-attributes --------------------------

func getStyleAttribute(h *html.Node) *localPseudoStylesheetType {
	if h != nil && h.Type == html.ElementNode {
		for _, attr := range h.Attr {
			if attr.Key == "style" {
				return &localPseudoStylesheetType{newLocalPseudoRule(attr.Val)}
			}
		}
	}
	return nil
}

type localPseudoStylesheetType struct {
	rule localPseudoRuleType
}

type localPseudoRuleType []style.KeyValue

func newLocalPseudoRule(styleAttr string) localPseudoRuleType {
	styles := strings.Split(styleAttr, ";")
	kv := make(localPseudoRuleType, 0, 3)
	for _, st := range styles {
		st = strings.TrimSpace(st)
		if len(st) > 0 {
			s := strings.Split(st, ":")
			if len(s) < 2 {
				T().Errorf("Skipping ill-formed style rule: %s", st)
			} else {
				k := strings.TrimSpace(s[0])
				v := strings.TrimSpace(s[1])
				kv = append(kv, style.KeyValue{k, style.Property(v)})
			}
		}
	}
	return kv
}

func (pseudorule localPseudoRuleType) Selector() string {
	return ""
}

func (pseudorule localPseudoRuleType) Properties() []string {
	var p []string
	for _, kv := range pseudorule {
		p = append(p, kv.Key)
	}
	return p
}
func (pseudorule localPseudoRuleType) Value(key string) style.Property {
	for _, kv := range pseudorule {
		if key == kv.Key {
			return kv.Value
		}
	}
	return style.NullStyle
}

func (pseudorule localPseudoRuleType) IsImportant(string) bool {
	return false
}

func (pseudosheet *localPseudoStylesheetType) AppendRules(s StyleSheet) {
	for _, r := range s.Rules() {
		for _, k := range r.Properties() {
			pseudosheet.rule = append(pseudosheet.rule, style.KeyValue{k, r.Value(k)})
		}
	}
}

func (pseudosheet *localPseudoStylesheetType) Empty() bool {
	return false
}

func (pseudosheet *localPseudoStylesheetType) Rules() []Rule {
	return []Rule{pseudosheet.rule}
}
