/*
Package pmmpost implements an interpreter for "Poor Man's MetaPost".

This is the implementation of an interpreter for "Poor Man's MetaPost",
my variant of the MetaPost graphical language. There is an accompanying
ANTLR grammar file, which describes the features and limitations of PMMPost.
I will sometimes refer to MetaFont, the original language underlying
MetaPost, as the grammar definitions are taken from Don Knuth's grammar
description in "The METAFONTBook".

The intent of PMMPost is to produce scalable drawings.
To accomplish this, users write a PMMPost program.
A typical PMMPost looks like this (statements for drawing a triangle):

	beginfig("triangle", 1cm, 1cm);
	   pair A, B, C;
	   A:=(0,0); B:=(1cm,0); C:=(0.5cm,1cm);
	   draw A--B--C--cycle;
	endfig;

PMMPost will, depending on the backend configuration, output a picture
for each "beginfig() ... endfig"-statement in the program.
For more information, please refer to
https://www.tug.org/docs/metapost/mpman.pdf.

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
package pmmpost

//go:generate antlr -Dlanguage=Go -o grammar -lib ../corelang -package grammar -Werror PMMPost.g4
//go:generate sh tagdoc.sh

import "github.com/npillmayer/schuko/tracing"

// We will trace to the InterpreterTracer
var T tracing.Trace
