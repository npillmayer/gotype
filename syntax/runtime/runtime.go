/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer

All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. Neither the name of Norbert Pillmayer or the names of its contributors
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

 * Package for implementing an interpreter runtime, consisting of
 * scopes, memory frames and symbols (variable declarations and references).
 * For a thorough discussion of an interpreter's runtime environment, refer to
 * "Language Implementation Patterns" by Terence Parr.

*/

package runtime

import (
	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

// Configurable trace
var T tracing.Trace = tracing.SyntaxTracer

// A type implementing a runtime environment for an interpreter
type Runtime struct {
	ScopeTree     *ScopeTree        // collect scopes
	MemFrameStack *MemoryFrameStack // runtime stack of memory frames
	ExprStack     *ExprStack        // evaluate arithmetic expressions
	PathBuilder   *PathStack        // construct paths
}

/* Construct a new runtime environment, initialized. Accepts a symbol creator for
 * variable declarations to be used within this runtime environment.
 */
func NewRuntimeEnvironment(withDeclarations func(string) Symbol) *Runtime {
	rt := &Runtime{}
	rt.ScopeTree = new(ScopeTree)                            // scopes for groups and functions
	rt.ScopeTree.PushNewScope("globals", withDeclarations)   // push global scope first
	rt.MemFrameStack = new(MemoryFrameStack)                 // initialize memory frame stack
	mf := rt.MemFrameStack.PushNewMemoryFrame("global", nil) // global memory
	mf.Scope = rt.ScopeTree.Globals()                        // connect the global frame with the global scope
	rt.MemFrameStack.Globals().SymbolTable = NewSymbolTable(withDeclarations)
	rt.ExprStack = NewExprStack()
	rt.PathBuilder = NewPathStack()
	return rt
}

