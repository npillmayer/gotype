package style

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"github.com/npillmayer/gotype/core/config/tracing"
	"golang.org/x/net/html"
)

/*
TODO
- implement shortcut-properties: "border", "background" etc.
+ style from header: link and sytle tags
- Locally scoped style (inline or <style> in body)
- Specifity
- be independent from goquery
+ be independent form douceur.css
- Fix API
+ Document code
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

// CSSOM is the "CSS Object Model", similar to the DOM for HTML.
// Our CSSOM consists of a set of stylesheets, each relevant for a sub-tree
// of the DOM. This sub-tree is called the "scope" of the stylesheet.
// Sub-trees are identified through the top node.
//
// Stylesheets are wrapped into an internal rules tree.
type CSSOM struct {
	rules             map[*html.Node]rulesTree
	defaultProperties PropertyMap
}

// RulesTree holds the styling rules of a stylesheet.
//
// Status: Currently this is not really a tree.
// Optimize some day (see
// https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/).
type rulesTree struct {
	stylesheet StyleSheet                   // stylesheet holding the rules
	selectors  map[string]cascadia.Selector // cache of compiled selectors
	source     PropertySource               // where do these rules come from?
}

// NewCSSOM creates an empty CSSOM.
// Clients need to supply a map of default property values.
// If pmap is nil, NewCSSOM will create a standard map.
func NewCSSOM(pmap PropertyMap) CSSOM {
	if pmap == nil {
		pmap = InitializeDefaultPropertyValues()
	}
	cssom := CSSOM{}
	cssom.rules = make(map[*html.Node]rulesTree)
	cssom.defaultProperties = pmap
	return cssom
}

// Properties may be defined at different places in HTML: as a sytlesheet
// reference link, within a <script> element in the HTML file, or in an
// attribute value.
//
// PropertySource influences the specifity of rules.
type PropertySource uint8

// Values for property sources, used when adding style sheets.
const (
	Global    PropertySource = iota + 1 // "browser" globals
	Author                              // CSS author (stylesheet link)
	Script                              // <script> element
	Attribute                           // in an element's attribute(s)
)

// AddStylesFor includes a stylesheet to a CSSOM and sets the scope for
// the stylesheet. If a stylesheet for the scope already exists, the
// styles are merged. css may be nil. If scope is nil then scope is the
// body (i.e., top-level content element) of a future document.
//
// source hints to where the stylesheet comes from.
func (cssom CSSOM) AddStylesFor(scope *html.Node, css StyleSheet, source PropertySource) error {
	rules, exists := cssom.rules[scope]
	if exists && css != nil {
		rules.stylesheet.AppendRules(css)
	} else {
		if scope == nil {
			scope = bodyElement
		} else {
			if scope.Type != html.ElementNode {
				return errors.New("Can style element nodes only")
			}
		}
		cssom.rules[scope] = rulesTree{css, nil, source}
	}
	return nil
}

// bodyElement is a symbolic node to denote the body element of a future
// HTML document.
var bodyElement *html.Node = &html.Node{}

// Empty is a predicate wether a stylesheet is empty, i.e. does not contain
// any rules.
func (rt rulesTree) Empty() bool {
	return rt.stylesheet == nil || rt.stylesheet.Empty()
}

// Internal helper for applying rules to an HTML node.
// In a first step it holds all the rules matching for node.
// In a second step it collects all the properties set in those rules,
// then orderes them by specifity.
type matchesList struct {
	matchingRules   []Rule
	propertiesTable []specifity
}

// Rules matches are collected from more than one stylesheet. Matching
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

func (rt *rulesTree) FilterMatchesFor(node *html.Node) *matchesList {
	list := &matchesList{}
	for _, rule := range rt.stylesheet.Rules() {
		selectorString := rule.Selector()
		var sel cascadia.Selector
		found := false
		if sel, found = rt.selectors[selectorString]; !found {
			var err error
			sel, err = cascadia.Compile(selectorString)
			if err != nil {
				tracing.EngineTracer.Errorf("CSS selector seems not to work: %s", selectorString)
				break
			}
			if rt.selectors == nil {
				rt.selectors = make(map[string]cascadia.Selector)
			}
			rt.selectors[selectorString] = sel
		}
		if sel.Match(node) {
			list.matchingRules = append(list.matchingRules, rule)
		}
	}
	tracing.EngineTracer.Debugf("matching rules for %s", nodePath(node))
	return list
}

// small helper to debug-print out a node
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

func (matches *matchesList) SortProperties() {
	var proptable []specifity
	for _, rule := range matches.matchingRules {
		for _, propertyKey := range rule.Properties() {
			value := Property(rule.Value(propertyKey))
			sp := specifity{Author, rule, propertyKey, value, rule.IsImportant(propertyKey), 0}
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
}

func (matches *matchesList) createStyleGroups(parent *StyledNode) map[string]*propertyGroup {
	m := make(map[string]*propertyGroup)
	done := make(map[string]bool, len(matches.propertiesTable))
	for _, pspec := range matches.propertiesTable { // for every specifity entry
		groupname, found := groupNameFromPropertyKey[pspec.propertyKey]
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
		if group, exists := m[groupname]; exists {
			group.Set(pspec.propertyKey, pspec.propertyValue)
		} else {
			_, pg := parent.findAncestorWithPropertyGroup(groupname) // must succeed
			if pg == nil {
				panic(fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", groupname))
			}
			group, isNew := pg.SpawnOn(pspec.propertyKey, pspec.propertyValue, true)
			if isNew { // a new property group has been created
				m[groupname] = group // put it into the group map
			}
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
	source        PropertySource
	rule          Rule
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
	//selectorstring := sp.rule.Selector()
	sp.spec = uint16((sp.source - 1) * 100)
	selspec := uint16(1) // TODO
	sp.spec += uint16(selspec)
	idcnt := uint16(1) // TODO
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

// --- Styled Node Tree -------------------------------------------------

type StyledNodeTree struct {
	root          *StyledNode
	defaultStyles map[string]*propertyGroup
}

func setupStyledNodeTree(domRoot *html.Node) *StyledNodeTree {
	defaultStyles := InitializeDefaultPropertyValues()
	viewport := &StyledNode{}
	viewport.node = domRoot
	viewport.computedStyles = defaultStyles
	tree := &StyledNodeTree{viewport, defaultStyles}
	return tree
}

// StyledNodes are the building blocks of the styled tree.
type StyledNode struct {
	node           *html.Node
	computedStyles map[string]*propertyGroup
	parent         *StyledNode
	children       []*StyledNode
}

func (sn *StyledNode) findAncestorWithPropertyGroup(group string) (*StyledNode, *propertyGroup) {
	var pg *propertyGroup
	it := sn
	for it != nil && pg == nil {
		pg = it.computedStyles[group]
		it = it.parent
	}
	return it, pg
}

func (sn *StyledNode) getProperty(group string, key string) Property {
	_, pg := sn.findAncestorWithPropertyGroup(group) // must succeed
	if pg == nil {
		panic(fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", group))
	}
	return pg.Cascade(key).Get(key) // must succeed
}

// For an explanation what's going on here, refer to
// https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/
// and
// https://limpet.net/mbrubeck/2014/08/23/toy-layout-engine-4-style.html
func (cssom CSSOM) Style(dom *html.Node) (*StyledNodeTree, error) {
	if dom == nil {
		return nil, errors.New("Nothing to style: empty document")
	}
	domDoc := goquery.NewDocumentFromNode(dom)
	styledTree := setupStyledNodeTree(dom)
	runner, err := cssom.attachStylesheets(domDoc)
	if err != nil {
		return nil, err
	}
	runner.doStyle(runner.startNode, styledTree.root)
	return styledTree, nil
}

// Scopes for all ruletrees have to be HTML element nodes.
// This is enforced during construction of the rules trees (adding to CSSDOM).
func (cssom CSSOM) attachStylesheets(domDoc *goquery.Document) (*stylingRunner, error) {
	if len(domDoc.Selection.Nodes) == 0 {
		return nil, errors.New("Nothing to style: empty document")
	}
	runner := newStylingRunner()
	for scope, rulestree := range cssom.rules {
		var stylingRootElement *goquery.Selection
		if scope == bodyElement { // scope is body element, i.e. whole document
			stylingRootElement := domDoc.Find("body")
			if stylingRootElement.Length() == 0 { // no body element found
				tracing.EngineTracer.Infof("Misconstructed DOM: cannot find <body>. Proceeding.")
				stylingRootElement = domDoc.Selection // root of fragment, try our best
			}
			runner.activeStylers[stylingRootElement.Nodes[0]] = rulestree
			runner.startNode = stylingRootElement.Nodes[0]
		} else {
			stylingRootElement = domDoc.FindNodes(scope)
			if stylingRootElement.Length() == 0 { // scope not not found in DOM
				// do not try to proceed: client meant something else...
				return nil, errors.New(fmt.Sprintf("Scope '%s' not found in DOM", scope.Data))
			}
			if stylingRootElement.Nodes[0].Type != html.ElementNode {
				return nil, errors.New(fmt.Sprintf("Scope '%s' is not of element type", shorten(scope.Data)))
			}
			runner.inactiveStylers[stylingRootElement.Nodes[0]] = rulestree
		}
	}
	if runner.startNode == nil {
		runner.startNode = domDoc.Selection.Nodes[0]
	}
	return runner, nil
}

type stylingRunner struct {
	activeStylers   map[*html.Node]rulesTree
	inactiveStylers map[*html.Node]rulesTree
	startNode       *html.Node
}

func newStylingRunner() *stylingRunner {
	runner := &stylingRunner{}
	runner.activeStylers = make(map[*html.Node]rulesTree)
	runner.inactiveStylers = make(map[*html.Node]rulesTree)
	return runner
}

func (runner *stylingRunner) activateStylesheetsFor(node *html.Node) {
	for scope, r := range runner.inactiveStylers {
		if scope == node {
			runner.activeStylers[scope] = r
			delete(runner.inactiveStylers, scope)
		}
	}
}

func (runner *stylingRunner) doStyle(node *html.Node, parent *StyledNode) {
	runner.activateStylesheetsFor(node)
	if createsStyledNode(node.Type) {
		sn := &StyledNode{}
		sn.node = node
		sn.parent = parent
		parent.children = append(parent.children, sn)
		var matchingRules *matchesList
		for _, rulesTree := range runner.activeStylers {
			matches := rulesTree.FilterMatchesFor(node)
			matchingRules = matchingRules.mergeMatchesWith(matches)
		}
		if matchingRules != nil {
			matchingRules.SortProperties()
			groups := matchingRules.createStyleGroups(parent)
			sn.computedStyles = groups // may be nil
		}
		parent = sn // continue with newly created styled node
	}
	node = node.FirstChild // now recurse into children
	for node != nil {
		if node.Type == html.ElementNode {
			runner.doStyle(node, parent)
		}
		node = node.NextSibling
	}
}

// Which HTML node type needs a corresponding styled node?
func createsStyledNode(nodeType html.NodeType) bool {
	if nodeType == html.ElementNode || nodeType == html.TextNode {
		return true
	}
	return false
}

func shorten(s string) string {
	if len(s) > 10 {
		return s[:10] + "..."
	}
	return s
}

// --- old versions -----------------------------------------------------

/*
func ConstructStyledNodeTree(dom *html.Node, rules rulesTree) (*StyledNodeTree, error) {
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
*/

/*
func (rulestree rulesTree) doStyle(scope *html.Node, domDoc *goquery.Document, styledTree *StyledNodeTree) (*StyledNodeTree, error) {
	if scope == nil || len(domDoc.Selection.Nodes) == 0 {
		return nil, errors.New("Nothing to style: empty scope or document")
	}
	if rulestree.Empty() {
		return styledTree, nil // we're done
	}
	var stylingRootElement *goquery.Selection
	if scope == bodyElement { // scope is body element, i.e. whole document
		stylingRootElement := domDoc.Find("body")
		if stylingRootElement.Length() == 0 {
			tracing.EngineTracer.Infof("Misconstructed DOM: cannot find <body>. Proceeding.")
			stylingRootElement = domDoc.Selection // root of fragment, try our best
		}
	} else {
		stylingRootElement = domDoc.FindNodes(scope)
		if stylingRootElement.Length() == 0 {
			// do not try to proceed: client meant something else...
			return nil, errors.New(fmt.Sprintf("Scope '%s' not found in DOM", scope.Data))
		}
	}
	//parentNode := stylingRootElement.Nodes[0].Parent
	//TODO
	//recurseStyling(parent, stylingRootElement.Nodes[0], rulestree)
	return styledTree, nil
}
*/

// node is a (possibly indirect) child of the node corresponding to parent.
/*
func recurseStyling(parent *StyledNode, node *html.Node, rules rulesTree) {
	if parent == nil || node == nil {
		return // recursion termination
	}
	sn := &StyledNode{}
	sn.node = node
	sn.parent = sn
	parent.children = append(parent.children, sn)
	matchingRules := rules.FilterMatchesFor(node)
	if matchingRules != nil {
		matchingRules.SortProperties()
		groups := matchingRules.createStyleGroups(parent)
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
*/

// NewRulesTree wraps a CSS stylesheet. css may be nil to represent an empty
// stylesheet.
/*
func newRulesTree(css *css.Stylesheet) *rulesTree {
	return &rulesTree{css, nil}
}
*/
