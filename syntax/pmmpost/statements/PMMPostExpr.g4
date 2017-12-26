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

This is a poor man's implementation of a METAFONT expression grammar.
It is heavily inspired by the METAFONT system by Donald E. Knuth and
METAPOST by John Hobby.

I sketched the grammar for a typesetting project, where polygones are
the target objects, not cubic splines. Therefore the grammar is
vastly simplfied, especially in the area of path expressions.

Other changes to the grammar are rooted in METAFONT's parser implementation:
the original syntax is strongly context-dependent. Given ANTLR V4's excellent
parsing capabilies even that might be handled by a meta-compiler, but
currently I am satisfied with this restricted grammar and won't bother
keeping closer to the original.

Differences to the original language include:

. simpler variable names
  - only 1 subscript allowed (e.g., x[k]r , but not x[k]r[s]l)
  - x.1 == x1 == x[1]  and  x1r == x1.r == x[1]r, but x[1].r is not allowed
  - variables starting with 'z' or uppercase letters are implicitly declared to be pairs
. decimals may not start with a '.' (you have to write '0.')
. no begingroup ... endgroup capsules within expressions
. transforms, but no transform-expressions (yet)
. paths are just polygones, i.e. no control points ('--' path joins only)
. paths variables are to be named 'path.'<tag>
. percentages are numbers: '50%' means 0.5
. line comments start with a '#', not '%'
. intersectiontimes syntax is changed
. some slight changes you probably won't notice...

Examples for valid numeric input are:

   y1r + 2/3a.b - 0.4[y5,y6] + x[k+1]r
   1/2[25%,70%] * 123     # interpolate midway between 25% and 75% of 123
   length (17,48) + xpart z0r + floor 3.14

Examples for valid pair-expression input are:

   Origin scaled 1.4 + 1/3(x1,y1)
   (xpart z1, ypart z0) rotated 45
   1/3[z1,z2] + n[Origin,(14,23)
   intersectionpoint path.p1 with path.p2

Examples for valid path-expression input are:

   z0 -- z1 -- subpath (2,5) of path.p -- cycle
   z0 -- point 2 of path.p -- cycle
   (x1,y1) -- reverse path.p1 shifted (13,234) -- z2r

*/

grammar PMMPostExpr;

import PMMPostTerminals;

expression
    : numtertiary
    | pairtertiary
    ;

numtertiary
    : numsecondary
    | numtertiary (PLUS|MINUS) numsecondary
    ; 

numsecondary
    : numprimary
    | numsecondary (TIMES|OVER) numprimary
    ;

numprimary
    : MATHFUNC numatom                               # funcnumatom
    | scalarmulop numatom                            # scalarnumatom
    | numatom LBRACKET numtertiary COMMA numtertiary RBRACKET       # interpolation
    | numtokenatom LBRACKET numtertiary COMMA numtertiary RBRACKET  # interpolation 
    | numatom                                        # simplenumatom
    | LENGTH pairprimary                             # pointdistance
    | PAIRPART pairprimary                           # pairpart
    ;

scalarmulop
    : (PLUS|MINUS)
    | numtokenatom
    ;

numtokenatom
    : DECIMALTOKEN OVER DECIMALTOKEN
    | DECIMALTOKEN
    ;

numatom
    : INTERNAL                                       # internal
    | WHATEVER                                       # whatever
    | MIXEDTAG ( subscript | anytag )*               # variable
    | TAG ( subscript | anytag )*                    # variable
    | DECIMALTOKEN UNIT?                             # decimal
    | LPAREN numtertiary RPAREN                      # subexpression
    | BEGINGROUP statementlist numtertiary ENDGROUP  # exprgroup
    ;


subscript
    : DECIMALTOKEN                                  // # shortsubscript
    | LBRACKET numtertiary RBRACKET                 // # subscriptexpression
    ;

pairtertiary
//  : INTERSECTIONPOINT pathsecondary WITH pathsecondary
    : pairsecondary
    | pairtertiary (PLUS|MINUS) pairsecondary
    ; 

pairsecondary
    : pairprimary                                    # pairsecond
    | pairsecondary (TIMES|OVER) numprimary          # pairsecond
    | numsecondary TIMES pairprimary                 # pairsecond
    | pairsecondary transformer                      # transform
    ;

transformer
    : SCALED numprimary
    | ROTATED numprimary
    | SHIFTED pairprimary
    ;

pairprimary
    : pairatom                                           # simplepairatom
    | scalarmulop pairatom                               # scalarmuloppair
    | POINT numtertiary OF pathprimary                   # pathpoint
    | numatom LBRACKET pairtertiary COMMA pairtertiary RBRACKET       # pairinterpolation
    | numtokenatom LBRACKET pairtertiary COMMA pairtertiary RBRACKET  # pairinterpolation
    ;

pairatom
    : LPAREN numtertiary COMMA numtertiary RPAREN        # literalpair
    | MIXEDPTAG ( subscript | anytag )*                  # pairvariable
    | PTAG ( subscript | anytag )*                       # pairvariable
    | LPAREN pairtertiary RPAREN                         # subpairexpression
    | BEGINGROUP statementlist pairtertiary ENDGROUP     # pairexprgroup
    ;

pathexpression
    : pathtertiary ( PATHCLIPOP pathtertiary )*
    ;

pathtertiary
    : pathfragm ( PATHJOIN pathfragm )* cycle?
    ;

pathfragm
    : pathsecondary
    | pairtertiary
    ;

cycle
    : PATHJOIN CYCLE
    ;

pathsecondary
    : pathprimary
    | pathsecondary transformer
    ;

pathprimary
    : pathatom                              # atomicpath
    | REVERSE pathprimary                   # reversepath
    | SUBPATH pairtertiary OF pathprimary   # subpath
    ;

pathatom
    : PATHTAG ( subscript | anytag )*       # pathvariable
//  | LPAREN pathexpression RPAREN          # nestedpath       // what for?
    ; 

tag
    : TAG
    | PTAG
    ;

anytag
    : TAG
    | MIXEDTAG
    | PTAG
    | MIXEDPTAG
    ;
