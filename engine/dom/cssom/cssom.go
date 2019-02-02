package cssom

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/dom"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"github.com/npillmayer/gotype/engine/tree"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

/*
TODO
+ implement shortcut-properties: "border", "background" etc.
+ style from header: link and sytle tags
+ Locally scoped style (inline or <style> in body)
+ Specifity
+ be independent from goquery
+ be independent form douceur.css
- Fix API:
  + create StyledTree using factory
  + define extension points for future properties: group "X"; registerCompound(key, func(...))
- write lots of tests
- make concurrent
+ Document code
- Create Diagram -> Wiki
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
	rules             map[*html.Node]rulesTree // CSS rules from stylesheets
	defaultProperties *style.PropertyMap       // "user agent" style properties
	compoundSplitters []CompoundPropertiesSplitter
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

// NewCSSOM creates an empty CSSOM.
// Clients are allowed to supply a map of additional/custom CSS property values.
// These may override values of the default ("user-agent") style sheet,
// or introduce completely new styling properties.
func NewCSSOM(additionalProperties []style.KeyValue) CSSOM {
	cssom := CSSOM{}
	cssom.rules = make(map[*html.Node]rulesTree)
	cssom.defaultProperties = style.InitializeDefaultPropertyValues(additionalProperties)
	cssom.compoundSplitters = make([]CompoundPropertiesSplitter, 1)
	cssom.compoundSplitters[0] = style.SplitCompoundProperty
	return cssom
}

// AddStylesForScope includes a stylesheet to a CSSOM and sets the scope for
// the stylesheet. If a stylesheet for the scope already exists, the
// styles are merged. css may be nil. If scope is nil then scope is the
// body (i.e., top-level content element) of a future document.
//
// source hints to where the stylesheet comes from.
// Its value will affect the calculation of specifity for rules of this
// stylesheet.
func (cssom CSSOM) AddStylesForScope(scope *html.Node, css StyleSheet, source PropertySource) error {
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

// bodyElement is a symbolic node to denote the body element of a future
// HTML document. AddStylesFor(...) with nil as a scope will replace it
// with this marker for scoping the complete document body.
var bodyElement *html.Node = &html.Node{}

// RegisterCompoundSplitter allows clients to handle additional compound
// properties. See type CompoundPropertiesSplitter.
func (cssom CSSOM) RegisterCompoundSplitter(splitter CompoundPropertiesSplitter) {
	if splitter != nil {
		cssom.compoundSplitters = append(cssom.compoundSplitters, splitter)
	}
}

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

// Rules matches have to be sorted by specifity. We'll sort the highest
// specifity up and won't overwrite earlier matches with later matches.

// sorter
type byHighestSpecifity []specifity

// make specifities sortable by highest sp.spec
func (sp byHighestSpecifity) Len() int           { return len(sp) }
func (sp byHighestSpecifity) Swap(i, j int)      { sp[i], sp[j] = sp[j], sp[i] }
func (sp byHighestSpecifity) Less(i, j int) bool { return sp[i].spec > sp[j].spec }

// This is a small helper to print out a table with rules matches for node.
func (m *matchesList) String() string {
	s := fmt.Sprintf("match of %d rules:\n", len(m.matchingRules))
	s += "Src +-- Spec. --+------------- Key --------------+------- Value ---------------\n"
	for _, sp := range m.propertiesTable {
		s += fmt.Sprintf("%3d | %9d | %30s | %s\n", sp.source, sp.spec, sp.propertyKey, sp.propertyValue)
	}
	return s
}

// FilterMatchesFor(node) iterates through all the rules relevant at this
// point and looks for rules matching the current HTML node.
// The heavy lifting is done by cascadia. We have to 'compile' all rules
// and will cache compiled rules.
//
// Will return a slice of CSS rules matched for node.
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
				T().Errorf("CSS selector seems not to work: %s", selectorString)
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
	T().Debugf("matching rules for %s", nodePath(node))
	return list
}

// SortProperties takes a slice of CSS rules (matched for an HTML node) and
// extracts all the properties set within the rules. These properties are
// then split into atomic properties, if they are compound properties
// (e.g.,
//     "margin" ⟹ "margin-top", "margin-right", ...
// Finally all property entries are sorted by specifity of the enclosing rule.
func (matches *matchesList) SortProperties(splitters []CompoundPropertiesSplitter) {
	var proptable []specifity
	for rno, rule := range matches.matchingRules {
		for _, propertyKey := range rule.Properties() {
			value := style.Property(rule.Value(propertyKey))
			props, err := splitCompoundProperty(splitters, propertyKey, value)
			if err != nil {
				sp := specifity{Author, rule, propertyKey, value, rule.IsImportant(propertyKey), 0}
				sp.calcSpecifity(rno)
				proptable = append(proptable, sp)
			} else {
				T().Debugf("%s is a compound style", propertyKey)
				for _, kv := range props {
					key := kv.Key
					val := kv.Value
					sp := specifity{Author, rule, key, val, rule.IsImportant(propertyKey), 0}
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

type specifity struct {
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
func (sp *specifity) calcSpecifity(no int) {
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
			//T.Infof("parent is %s", parent)
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

// StyledNode is the node type for the styled tree to build. The CSSOM
// styler builds a tree of styled nodes, more or less corresponding to the
// stylable DOM nodes. We de-couple the Styler from the type of the resulting
// styled nodes by using a builder (see StyledTreeBuilder) and an interface type.
// type StyledNode interface {
// 	ComputedStyles() *style.PropertyMap   // get the computed styles of this styled node
// 	SetComputedStyles(*style.PropertyMap) // set the computed styles of this styled node
// }

// StyledTreeBuilder is a builder to create tree of conrete implementations
// of interface StyleNode (see type StyledNode).
//
// ATTENTION: Tree construction may be performed concurrently, so all methods
// (especially LinkNodeToParent) must be thread-safe!
// type StyledTreeBuilder interface {
// 	MakeNodeFor(*html.Node) StyledNode                 // create a new styled node
// 	LinkNodeToParent(sn StyledNode, parent StyledNode) // attach a node to the tree
// 	WalkUpwards(StyledNode) StyledNode                 // walk to parent of node
// }

//func setupStyledNodeTree(domRoot *html.Node, defaults *style.PropertyMap, builder StyledTreeBuilder) StyledNode {
func setupStyledNodeTree(domRoot *html.Node, defaults *style.PropertyMap,
	creator style.Creator) *tree.Node {
	//
	//viewport := builder.MakeNodeFor(domRoot)
	viewportNode := creator.StyleForHtmlNode(domRoot)
	creator.SetComputedStyles(viewportNode, defaults)
	return viewportNode
}

//func findAncestorWithPropertyGroup(sn StyledNode, group string, builder StyledTreeBuilder) (StyledNode, *style.PropertyGroup) {
func findAncestorWithPropertyGroup(sn *tree.Node, group string,
	creator style.Creator) (*tree.Node, *style.PropertyGroup) {
	//
	var pg *style.PropertyGroup
	it := sn // start search at styled node itself, then proceed upwards
	for it != nil && pg == nil {
		if it == nil {
			panic("this should not happen: user agent defaults not set?")
		}
		styles := creator.ToStyler(sn).ComputedStyles()
		if styles != nil {
			pg = styles.Group(group)
		}
		it = it.Parent()
	}
	return it, pg
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
	if cssom.rules == nil {
		T().Infof("Styling HTML tree without having any CSS rules")
	}
	styledRootNode := setupStyledNodeTree(dom, cssom.defaultProperties, creator)
	walker := tree.NewWalker(styledRootNode) // create a concurrent tree walker
	doStyle := func(node *tree.Node, parent *tree.Node, pos int) (*tree.Node, error) {
		// provide closure with style creator
		return createStyledChildren(node, creator)
	}
	future := walker.TopDown(doStyle).Promise() // build the style tree
	if _, err := future(); err != nil {
		T().Errorf("Error while creating styled tree: %v", err)
		return nil, err
	}
	// runner := newStylingRunner(creator, cssom.compoundS()plitters)
	// err := cssom.attachStylesheets(dom, runner)
	// if err != nil {
	// 	return nil, err
	// }
	// runner.doStyle(runner.startNode, styledRootNode)
	return styledRootNode, nil
}

// Pre-condition: sn has been styled and points to an HTML node.
// Now iterate through the HTML children and create styled nodes for each.
func createStyledChildren(parent *tree.Node, creator style.Creator) (*tree.Node, error) {
	domnode := dom.NewRONode(parent, creator.ToStyler) // interface RODomNode
	h := domnode.HtmlNode()
	if h.Type == html.ElementNode {
		ch := h.FirstChild
		for ch != nil {
			if ch.DataAtom == atom.Style { // <style> element
				T().Debugf("<style> node has Data = %v", h.Data)
				if ch.FirstChild.Type == html.TextNode {
					T().Debugf("first child is text: %v", ch.FirstChild.Data)
				}
				// TODO attach styles from <style>
			} else if isStylable(ch.DataAtom) {
				sn := creator.StyleForHtmlNode(ch)
				parent.AddChild(sn) // sn will be sent to next pipeline stage
				// TODO attach style attributes
			}
			ch = ch.NextSibling
		}
	} else if h.Type == html.TextNode {
		// do not send text node to next pipeline stage => will not be styled
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

// Scopes for all ruletrees have to be HTML element nodes.
// This is enforced during construction of the rules trees (adding to CSSDOM).
func (cssom CSSOM) attachStylesheets(dom *html.Node) error {
	for scope, rulestree := range cssom.rules {
		var stylingRootElement *html.Node
		if scope == bodyElement { // scope is body element, i.e. whole document
			stylingRootElement := findBodyElement(dom)
			if stylingRootElement == nil { // no body element found
				T().Infof("Misconstructed DOM: cannot find <body>. Proceeding.")
				stylingRootElement = dom // root of fragment, try our best
			}
			// runner.activeStylers[stylingRootElement] = rulestree
			// runner.startNode = stylingRootElement
		} else {
			stylingRootElement = findThisNode(dom, scope)
			if stylingRootElement == nil { // scope not not found in DOM
				// do not try to proceed: client meant something else...
				return errors.New(fmt.Sprintf("Scope '%s' not found in DOM", scope.Data))
			}
			if stylingRootElement.Type != html.ElementNode {
				return errors.New(fmt.Sprintf("Scope '%s' is not of element type", shorten(scope.Data)))
			}
			//runner.inactiveStylers[stylingRootElement] = rulestree
		}
	}
	// if runner.startNode == nil {
	// 	runner.startNode = dom
	// }
	return nil
}

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

func calcComputedStylesForNode(node *html.Node, parent *tree.Node) *style.PropertyMap {
	//
	// var matchingRules *matchesList
	// for _, rulesTree := range runner.activeStylers {
	// 	matches := rulesTree.FilterMatchesFor(node)
	// 	matchingRules = matchingRules.mergeMatchesWith(matches)
	// }
	// if matchingRules != nil {
	// 	matchingRules.SortProperties(runner.splitters)
	// 	pmap := matchingRules.createStyleGroups(parent, runner.builder)
	// 	return pmap
	// }
	return nil
}

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
