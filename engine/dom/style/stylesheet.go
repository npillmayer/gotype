package style

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

// In order to de-couple implementations of CSS-stylesheets from the
// construction of the styled node tree, we introduce an interface
// for CSS stylesheets. Clients for the styling engine will have to
// provide a concrete implementation of this interface (e.g., see
// package douceuradapter).
//
// Having this interface imposes a performance hit. However, this
// implementation of CSS-styling will never trade modularity and
// clarity for performance. Clients in need for a production grade
// browser engine (where performance is key) should opt for headless
// versions of the main browser projects.
//
// See interface Rule.
type StyleSheet interface {
	AppendRules(StyleSheet) // append rules from another stylesheet
	Empty() bool            // does this stylesheet contain any rules?
	Rules() []Rule          // all the rules of a stylesheet
}

// Rules belong to a StyleSheet.
//
// See interface StyleSheet.
type Rule interface {
	Selector() string        // the prelude / selectors of the rule
	Properties() []string    // property keys, e.g. "margin-top"
	Value(string) string     // property value for key, e.g. "15px"
	IsImportant(string) bool // is property key marked as important?
}
