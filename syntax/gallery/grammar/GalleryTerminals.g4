/*
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

This is a sub-grammar defining the lexical tokens of the 'Galery' DSL.

*/

lexer grammar PMMPostTerminals;

ASSIGN   : ':=' ;
EQUALS   : '=' ;
COLON    : ':' ;
SEMIC    : ';' ;
COMMA    : ',' ;
LPAREN   : '(' ;
RPAREN   : ')' ;
LBRACKET : '[' ;
RBRACKET : ']' ; 

PATHJOIN : ('--' '-'? | '..' '.'? | '&' ) ;

PLUS  : '+' ;
MINUS : '-' ;
TIMES : '*' ;
OVER  : '/' ;

PARALLEL    : '||' ;
PERPENDIC   : '|-' ;
CONGRUENT   : '~' ;

WHATEVER
    : '?'
    | 'whatever'
    ;

UNIT : ( 'bp' | 'mm' | 'cm' | 'in' | 'pt' | 'pc' ) ;

TYPE : 'numeric' | 'pair' | 'path' | 'framedbox' | 'framedpath' | 'transform' ;

INTERNAL : 'width' | 'height' ;

PAIRPART   : 'xpart' | 'ypart' ;
EDGECONSTR : 'top' | 'lft' | 'rt' | 'bot' ;
EDGE       : 'edge' ;
FRAME      : 'frame' ;
BOX        : 'box' ;

MATHFUNC : 'length' | 'floor' | 'ceil' | 'sqrt' ;
SUBPATH  : 'subpath' ;
REVERSE  : 'reverse' ;
WITH     : 'with' ;
POINT    : 'point' ;
OF       : 'of' ;

TRANSFORM : 'scaled' | 'shifted' | 'rotated' | 'transformed' ;

CYCLE    : 'cycle' ;

PATHCLIPOP : 'union' | 'intersection' | 'difference' ;

BEGINGROUP : 'begingroup' ;
ENDGROUP   : 'endgroup' ;

PROOF      : 'proof' ;
SAVE       : 'save' ;
SHOW       : 'show' ;

TAG        : [a-zA-Z]+ ;
MIXEDTAG   : '.'? [a-zA-Z] [.a-zA-Z0-9]* ;

DECIMALTOKEN : '0'..'9'+ ( DOT '0'..'9'+ )? '%'? ;

DOT      : '.' ;

LABEL : '"' [a-zA-Z][a-zA-Z0-9 \-_]+ '"' ;

fragment
LINECOMMENT
    : '//' .*? '\r'? '\n'      // line comment
    ;

WS
    : ( [ \r\n\t]+ | LINECOMMENT )        -> skip
    ;
