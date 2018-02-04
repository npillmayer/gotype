/*
Package lr implements a simple SLR(1) parser generator.
It is mainly intended for Markdown parsing, but may be of use for
other purposes, too.

Building a Grammar

Grammars are specified using a grammar builder object. Clients add
rules, consisting of non-terminal symbols and terminals. Terminals
carry a token value of type int. Grammars may contain epsilon-productions.

Example:

    b := lr.NewGrammarBuilder("G")
    b.LHS("S").N("A").T("a", 1).EOF()  // S  ->  A a EOF
    b.LHS("A").N("B").N("D").End()     // A  ->  B D
    b.LHS("B").T("b", 2).End()         // B  ->  b
    b.LHS("B").Epsilon()               // B  ->
    b.LHS("D").T("d", 3).End()         // D  ->  d
    b.LHS("D").Epsilon()               // D  ->

This results in the following trivial grammar:

   b.Grammar().Dump()

   0: [S] ::= [A a #eof]
   1: [A] ::= [B D]
   2: [B] ::= [b]
   3: [B] ::= []
   4: [D] ::= [d]
   5: [D] ::= []

Static Grammar Analysis

After the grammar is complete, it has to be analysed. For this end, the
grammar is subjected to a GrammarAnalysis object, which computes FIRST and
FOLLOW sets for the grammar and determines all epsilon-derivable rules.

Although FIRST and FOLLOW-sets are mainly intended to be used for internal
purposes of constructing the parser tables, methods for getting FIRST(A)
and FOLLOW(A) of non-terminals are defined to be public.

    ga := lr.NewGrammarAnalysis(g)  // analyser for grammar above
    ga.Grammar().EachNonTerminal(
        func(name string, A Symbol) interface{} {             // ad-hoc mapper function
            fmt.Printf("FIRST(%s) = %v", name, ga.First(A))   // get FIRST-set for A
            return nil
        })

    // Output:
    FIRST(S) = [1001 1002 1003]     // terminal token values as int, 1001 = 'a'
    FIRST(A) = [999 1002 1003]      // 999 = epsilon
    FIRST(B) = [999 1002]           // 1002 = 'b'
    FIRST(D) = [999 1003]           // 1003 = 'd'

Parser Construction

Using grammar analysis as input, a bottom-up parser can be constructed.
First a characteristic finite state machine (CFSM) is built from the
grammar, which will then be transformed into a GOTO table (LR(0)-table)
for a SLR parser. The CFSM will not be thrown away, but is made available
to the client.  This is intended
for debugging purposes, but may be useful for error recovery, too.
It can be exported to Graphviz's Dot-format.

Example:

    lrgen := NewLRTableGenerator(ga)  // ga is a GrammarAnalysis, see above
    lrgen.CreateTables()              // construct LR parser tables


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
*/
package lr

import (
	"github.com/npillmayer/gotype/gtcore/config/tracing"
)

// Configurable trace
var T tracing.Trace = tracing.SyntaxTracer
