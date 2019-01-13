/*
----------------------------------------------------------------------
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

This is an ANTLR V4 sub-grammar file for an expression centric language 
reminiscent of MetaFont/MetaPost.
The DSL borrows concepts from MetaFont/MetaPost. See the project
documentation for details.

*/

parser grammar CoreLang;

statementlist
    : ( statement )+ EOF
    ;

vardef
    : VARDEF TAG ( COMMA TAG )*
    ;

compound
    : BEGINGROUP statementlist ENDGROUP
    ;

empty
    :
    ;

assignment
    /*: variable ASSIGN expression*/
	: variable ASSIGN path
    ;

constraint
    : equation
    | orientation
    ;

equation
    : expression ( EQUALS expression )+
    ;

orientation
    : tertiary ( (PARALLEL|PERPENDIC|CONGRUENT) tertiary )+
    ;

token
    : ( PLUS | MINUS | TIMES | OVER | ASSIGN | PARALLEL | PERPENDIC | CONGRUENT
        | BEGINGROUP | ENDGROUP | EDGECONSTR | PATHCLIPOP | PATHJOIN
        | EDGE | FRAME | BOX | REVERSE | SUBPATH | PROOF | SAVE | SHOW
        | TRANSFORM | TAG )
    ;

// --- Expressions -----------------------------------------------------------

expression
    : tertiary
    | expression PATHCLIPOP tertiary
    ;

tertiary
    : path                                       # pathtertiary
    | tertiary (PLUS|MINUS) secondary            # term
    | secondary                                  # term
    ;

path
    : secondary ( '..' secondary )+ cycle?
    ; 

cycle
    : pathjoin CYCLE
    ; 

secondary
    : primary                                    # factor
    | secondary (TIMES|OVER) primary             # factor
    | secondary ( TRANSFORM primary )+           # transform
    ;

primary
    : MATHFUNC atom                               # funcatom
    | scalarmulop atom                            # scalaratom
    | numtokenatom LBRACKET tertiary COMMA tertiary RBRACKET  # interpolation 
    | atom LBRACKET tertiary COMMA tertiary RBRACKET          # interpolation
    | atom                                        # simpleatom
    | PAIRPART primary                            # pairpart       // => numeric
    | POINT tertiary OF primary                   # pointof        // => pair
    | REVERSE primary                             # reversepath    // => path
    | SUBPATH tertiary OF primary                 # subpath        // => path
    ;

scalarmulop
    : (PLUS|MINUS)? numtokenatom
    ;

numtokenatom
    : DECIMALTOKEN OVER DECIMALTOKEN
    | DECIMALTOKEN
    ;

atom
    : DECIMALTOKEN UNIT?                             # decimal
    | variable                                       # varatom
    | LPAREN tertiary COMMA tertiary RPAREN          # literalpair
    | LPAREN tertiary RPAREN                         # subexpression
    | BEGINGROUP statementlist tertiary ENDGROUP     # exprgroup
    ;

variable
    : MIXEDTAG ( subscript | anytag )*
    | TAG ( subscript | anytag )*
    | LAMBDAARG
    ;

subscript
    : DECIMALTOKEN
    | LBRACKET tertiary RBRACKET
    ;

anytag
    : TAG
    | MIXEDTAG
    ;

