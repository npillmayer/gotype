/*
Package douceuradapter is a concrete implementation of interface style.StyleSheet.

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
*/
package douceuradapter

import (
	"github.com/aymerick/douceur/css"
	"github.com/npillmayer/gotype/engine/dom/style"
)

// CssStyle is an adapter for interface style.StyleSheet.
// For an explanation of the motivation behind this design, please refer
// to documentation for interface style.StyleSheet.
type CssStyles struct {
	css css.Stylesheet
}

// Wrap a douceur.css.Stylesheet into CssStyles.
// The stylesheet is now managed by the wrapper.
func Wrap(css *css.Stylesheet) *CssStyles {
	sheet := &CssStyles{*css}
	return sheet
}

// Does this stylesheet contain any rules?
//
// Interface style.StyleSheet
func (sheet *CssStyles) Empty() bool {
	return len(sheet.css.Rules) == 0
}

// Append rules from another stylesheet
//
// Interface style.StyleSheet
func (sheet *CssStyles) AppendRules(other style.StyleSheet) {
	othercss := other.(*CssStyles)
	for _, r := range othercss.css.Rules { // append every rule from other
		sheet.css.Rules = append(sheet.css.Rules, r)
	}
}

// All the rules of a stylesheet
//
// Interface style.StyleSheet
func (sheet *CssStyles) Rules() []style.Rule {
	rules := make([]style.Rule, len(sheet.css.Rules))
	for i := range sheet.css.Rules {
		r := sheet.css.Rules[i]
		rules[i] = Rule(*r)
	}
	return rules
}

var _ style.StyleSheet = &CssStyles{}

// Rule is an adapter for interface style.Rule.
type Rule css.Rule

// The prelude / selectors of the rule
func (r Rule) Selector() string {
	return r.Prelude
}

// Property keys, e.g. "margin-top"
func (r Rule) Properties() []string {
	decl := r.Declarations
	props := make([]string, 0, len(decl))
	for _, d := range decl {
		props = append(props, d.Property)
	}
	return props
}

// Property value for key, e.g. "15px"
func (r Rule) Value(key string) string {
	decl := r.Declarations
	for _, d := range decl {
		if d.Property == key {
			return d.Value
		}
	}
	return ""
}

// Is property key marked as important?
func (r Rule) IsImportant(key string) bool {
	decl := r.Declarations
	for _, d := range decl {
		if d.Property == key {
			return d.Important
		}
	}
	return false
}

var _ style.Rule = &Rule{}
