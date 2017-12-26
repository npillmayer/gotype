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
3. Neither the name of Tom Everett nor the names of its contributors
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

This part of the grammar parses variable references.

* Variable Names
* ==============
* Variables are complex things in MetaFont/MetaPost. These are legal:

   metafont> showvariable x;
   x=1
   x[]=numeric
   x[][]=numeric
   x[][][]=numeric
   x[][][][]=numeric
   x[]r=numeric
   x[]r[]=numeric
   ...

* Identifier-strings are called "tags". In the example above, 'x' is a tag
* and 'r' is a suffix.
*
* Array variables may be referenced without brackets, if the subscript is just
* a numeric literal, i.e. x[2]r and x2r refer to the same variable. We do
* not rely on the parser to decipher these kinds of variable names for us,
* but rather break up x2r16a => x[2]r[16]a by hand. However, the parser will
* split up array indices in brackets, for the subscript may be a complex expression
* ("x[ypart ((8,5) rotated 20)]" is a valid expression in MetaFont).
* Things are further complicated by the fact that subscripts are allowed to
* be decimals: x[1.2] is valid, and may be typed "x1.2".

   metafont> x[ypart ((8,5) rotated 20)] = 1;
   ## x7.4347=1
*/

grammar PMMPVar;

variable
    : tag (suffix | subscript)* MARKER  // we'll always append a marker at the end
    ;

tag
    : PATHTAG             # pathtag
    | TAG                 # simpletag
    ;

suffix
    : DOT TAG
    | TAG
    ;

subscript
    : DECIMAL
    | '[' DECIMAL ']'
    ;

MARKER : '@' ;

PATHTAG
    : 'path.' ('a'..'z'|'A'..'Z')+
    ;

TAG
    : ('a'..'z' | 'A'..'Z')+
    ;

DECIMAL
    : '0'..'9'+ ( DOT '0'..'9'+ )?
    ;

DOT : '.' ;

WS
    : ( [ \r\n\t]+ )        -> skip
    ;
