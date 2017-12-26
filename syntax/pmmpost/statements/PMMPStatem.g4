/*
BSD License
Copyright (c) 2017, Norbert Pillmayer <norbert@pillmayer.com>

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

----------------------------------------------------------------------

This is a poor man's implementation of a METAPOST drawing engine.
It is heavily inspired by the METAFONT system by Donald E. Knuth and
METAPOST by John Hobby.

I sketched the grammar for a typesetting project, where polygones are
the target objects, not cubic splines. Therefore the grammar is
vastly simplfied, especially in the area of path expressions.

Other changes to the grammar are rooted in METAFONT's parser implementation:
the original syntax is strongly context-dependent. Given ANTLR V4's excellent
parsing capabilies even that might be handled by a meta-compiler, but
currently I am satisfied with a restricted grammar and won't bother
keeping closer to the original.

Differences to the original language include:

. No variable declaration: variables relate to types by naming convention. TODO
. No cubic splines: we use polygones only.

*/

grammar PMMPStatem;

import PMMPostExpr;

figures
    : figure* EOF
    ;

figure
    : beginfig statementlist endfig
    ;

beginfig
    : 'beginfig' '(' LABEL ',' DECIMALTOKEN UNIT ',' DECIMALTOKEN UNIT ')' SEMIC
    ;

endfig
    : 'endfig' SEMIC
    ;

statementlist
    : ( statement SEMIC )*
    ;

statement
    : compound
    | equation
    | declaration
    | assignment
    | command
    ;

compound
    : BEGINGROUP statementlist ENDGROUP
    ;

equation
    : numtertiary ( EQUALS numtertiary )+     # multiequation
    | pairtertiary (EQUALS pairtertiary )+    # multiequation
    | pathatom EQUALS pathexpression          # pathequation
    ;

declaration
    : mptype tag ( COMMA tag )*
    ;

mptype
    : NUMERIC
    | PAIR
    | PATH
    | PEN
    | COLOR
    ;

assignment
    :  lvalue ASSIGN expression
    ;

lvalue
    : MIXEDTAG ( subscript | anytag )*
    | TAG ( subscript | anytag )*
    | MIXEDPTAG ( subscript | anytag )*
    | PTAG ( subscript | anytag )*
    ;

command
    : saveStmt
    | showvariableCmd
    | drawCmd
    | fillCmd
    | pickupCmd
    ;

saveStmt
    : 'save' tag (COMMA tag)*
    ;

showvariableCmd
    : 'showvariable' tag
    ;

drawCmd
    : 'draw' pathexpression
    ;


fillCmd
    : 'fill' pathexpression
    ;

pickupCmd
    : 'pickup' PEN ( 'scaled' DECIMALTOKEN )? ( 'withcolor' COLOR )?
    ;

