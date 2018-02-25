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

This is an ANTLR V4 grammar file for a language called 'Gallery'. Gallery
is a DSL for placing frames on pages. It is part of a typesetting project
which uses frames to place text on a page.

The DSL borrows concepts from MetaFont / MetaPost. See the project
documentation for details.

https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/

*/

grammar Gallery;

import GalleryTerminals, CoreLang;

program
    : statementlist EOF
    ;

statement
    : compound
    | declaration
    | vardef
    | assignment
    | constraint
    | command
    | empty
    ;

declaration
    : TYPE TAG ( COMMA TAG )*                        # typedecl
    | LOCAL TYPE? TAG ( COMMA TAG )*                 # localdecl
    | PARAMETER variable ( tertiary TO tertiary )?   # parameterdecl
    ;

command
    : SAVE TAG (COMMA TAG)*          # savecmd
    | SHOW TAG (COMMA TAG)*          # showcmd
    | PROOF LABEL                    # proofcmd
    | LET token EQUALS MATHFUNC      # letcmd
    ;

// --- Expressions -----------------------------------------------------------

pathjoin
    : PATHJOIN
    ;
