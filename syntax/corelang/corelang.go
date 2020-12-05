/*
Package corelang implements core commands for DSLs dealing with arithmetic
expressions, pairs and paths. It borrows from MetaFont/MetaPost,
as described in the accompanying ANTLR grammar file.

Language Features

This package includes functions for numerous language features common to
MetaFont-/MetaPost-derivated DSLs.

The implementation is tightly coupled to the ANTLR V4 parser generator.
ANTLR is a great tool and I see no use in being independent from it.

Lua Scripting

This package also includes the support for Lua scripting. The DSLs stemming
from this language core are Lua-enabled by default.

For further information please refer to types Scripting and LuaVarRef.


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
package corelang

//go:generate antlr -Dlanguage=Go -o grammar -lib . -package grammar -Werror CoreLang.g4

import "github.com/npillmayer/schuko/tracing"

// Trace to the InterpreterTracer
func T() tracing.Trace {
	return tracing.InterpreterTracer
}

// Trace to the ScriptingTracer
func S() tracing.Trace {
	return tracing.ScriptingTracer
}
