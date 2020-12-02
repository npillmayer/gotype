/*
Package terex provides term rewriting expressions as a basis for
rewriting parse-trees and ASTs. It implements types for
a homogenous abstract syntax tree in a Lisp-like fashion.

Parsing generates a parse tree, which is too verbose for further processing.
Instead of long chains of grammar production symbols we usualy prefer a
much more compact AST (abstract syntax tree). One possible variant of
ASTs is a *homogenous* tree, i.e. one where the structure of all nodes
is identical. This makes tree walking easy.

This module provides the core Go data types to create and modify
homogenous trees. Homogenous trees are usually built around some Node type.
However, there is a programming language specialized in homogenous lists and
trees: Lisp (or Clojure, if you prefer). We implement node types which are
reminiscent of Lisp CONS, and call the resulting mini-language TeREx
(Term Rewriting Expressions).

With homogenous tree nodes there is always one caveat: type information of the
implementing programming language is compromised. Therefore, in absence of generics,
the code in this module heavily uses "interface{}" and relies on type switches and
casts. This is sometimes cumbersome to read, but on the other hands brings convenience
for a certain set of operations, including tree walking and tree restructuring.


BSD License

Copyright (c) 2019–21, Norbert Pillmayer

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
package terex

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
)

// T traces to the global syntax tracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}
