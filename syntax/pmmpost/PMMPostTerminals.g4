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

This is a sub-grammar defining the lexical tokens of the "Poor Man's
MetaPost" graphics language.

*/

lexer grammar PMMPostTerminals;

import CoreTerminals;

PATHJOIN : ('--' '-'? | '..' '.'? | '&' ) ;
CONTROLS : 'controls' ;
AND      : 'and' ;
CURL     : 'curl' ;

TYPE  : 'numeric' | 'pair' | 'path' | 'transform' | 'framedbox' | 'framedpath' ;

FIGURE    : 'figure' ;
ENDFIG    : 'endfig' ;
PICKUP    : 'pickup' ;
PEN       : 'pencircle' | 'pensquare' ;
FILL      : 'fill' ;
DRAW      : 'draw' ;
WITHCOLOR : 'withcolor' ;

COLOR     : '#' [0-9a-f][0-9a-f][0-9a-f] ( [0-9a-f][0-9a-f][0-9a-f] )? ;
