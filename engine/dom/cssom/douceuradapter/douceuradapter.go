/*
Package douceuradapter is a concrete implementation of interface cssom.StyleSheet.

BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package douceuradapter

import (
	"github.com/aymerick/douceur/css"
	"github.com/aymerick/douceur/parser"
	"github.com/npillmayer/gotype/engine/dom/cssom"
	"github.com/npillmayer/gotype/engine/dom/cssom/style"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// CSSStyles is an adapter for interface cssom.StyleSheet.
// For an explanation of the motivation behind this design, please refer
// to documentation for interface cssom.StyleSheet.
type CSSStyles struct {
	css css.Stylesheet
}

// Wrap a douceur.css.Stylesheet into CssStyles.
// The stylesheet is now managed by the wrapper.
func Wrap(css *css.Stylesheet) *CSSStyles {
	sheet := &CSSStyles{*css}
	return sheet
}

// Empty checks if this stylesheet contains any rules.
//
// Interface cssom.StyleSheet
func (sheet *CSSStyles) Empty() bool {
	return len(sheet.css.Rules) == 0
}

// AppendRules appends rules from another stylesheet.
//
// Interface cssom.StyleSheet
func (sheet *CSSStyles) AppendRules(other cssom.StyleSheet) {
	othercss := other.(*CSSStyles)
	for _, r := range othercss.css.Rules { // append every rule from other
		sheet.css.Rules = append(sheet.css.Rules, r)
	}
}

// Rules returns all the rules of a stylesheet.
//
// Interface style.StyleSheet
func (sheet *CSSStyles) Rules() []cssom.Rule {
	rules := make([]cssom.Rule, len(sheet.css.Rules))
	for i := range sheet.css.Rules {
		r := sheet.css.Rules[i]
		rules[i] = Rule(*r)
	}
	return rules
}

var _ cssom.StyleSheet = &CSSStyles{}

// Rule is an adapter for interface cssom.Rule.
type Rule css.Rule

// Selector returns the prelude / selectors of the rule.
func (r Rule) Selector() string {
	return r.Prelude
}

// Properties returns the property keys of a rule,
// e.g. "margin-top"
func (r Rule) Properties() []string {
	decl := r.Declarations
	props := make([]string, 0, len(decl))
	for _, d := range decl {
		props = append(props, d.Property)
	}
	return props
}

// Value returns the property values for given key with this rule, e.g. "15px"
func (r Rule) Value(key string) style.Property {
	decl := r.Declarations
	for _, d := range decl {
		if d.Property == key {
			return style.Property(d.Value)
		}
	}
	return ""
}

// IsImportant returns true if a style key is marked as important ("!").
func (r Rule) IsImportant(key string) bool {
	decl := r.Declarations
	for _, d := range decl {
		if d.Property == key {
			return d.Important
		}
	}
	return false
}

var _ cssom.Rule = &Rule{}

// ExtractStyleElements visits <head> and <body> elements in an HTML parse
// tree and searches for embedded <style>s. It returns the content of
// style-elements as style sheets.
func ExtractStyleElements(htmldoc *html.Node) []*CSSStyles {
	head := findElement(atom.Head, htmldoc)
	body := findElement(atom.Body, htmldoc)
	css := extractStyles(head)
	css2 := extractStyles(body)
	for _, c := range css2 {
		css = append(css, c)
	}
	return css
}

func extractStyles(h *html.Node) []*CSSStyles {
	var css []*CSSStyles
	ch := h.FirstChild
	for ch != nil {
		if ch.DataAtom == atom.Style {
			c, err := parser.Parse(ch.FirstChild.Data)
			if err != nil {
				break
			}
			css = append(css, Wrap(c))
		}
		ch = ch.NextSibling
	}
	return css
}

func findElement(a atom.Atom, h *html.Node) *html.Node {
	if h == nil {
		return nil
	}
	if h.DataAtom == a {
		return h
	}
	ch := h.FirstChild
	for ch != nil {
		r := findElement(a, ch)
		if r != nil && r.DataAtom == a {
			return r
		}
		ch = ch.NextSibling
	}
	return nil
}
