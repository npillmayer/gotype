/*
Package dom will some day provide utilities for HTMLbook DOMs.

Status

Early draft—API may change frequently. Please stay patient.

Overview

HTMLbook is the core DOM of our documents.
Background for this decision can be found under
https://www.balisage.net/Proceedings/vol10/print/Kleinfeld01/BalisageVol10-Kleinfeld01.html
and http://radar.oreilly.com/2013/09/html5-is-the-future-of-book-authorship.html

Excerpt: "In this paper, I argue that HTML5 offers unique advantages to authors and publishers in comparison to both traditional word processing and desktop publishing tools like Microsoft Word and Adobe InDesign, as well as other markup vocabularies like DocBook and AsciiDoc. I also consider the drawbacks currently inherent in the HTML5 standard with respect to representing long-form, structured text content, and the challenges O’Reilly has faced in adopting the standard as the new source format for its toolchain. Finally, I discuss how O’Reilly has surmounted these challenges by developing HTMLBook, a new open, HTML5-based XML standard expressly designed for the authoring and production of both print and digital book content."

For an in-depth description of HTMLbook please refer to
https://oreillymedia.github.io/HTMLBook/.

Tree Implementation

Styling and layout of HTML/CSS involves a lot of operations on different trees.
We implement the various trees on top of a general purpose tree type
(package engine/tree), which offers concurrent operations to manipluate
tree nodes.

In a fully object oriented programming language we would subclass this
tree type for every type of tree in use (styled tree, layout tree,
render tree), but in Go we resort to composition, thus including a
generic tree node in every node (sub-)type. The downside of this approach
is that we will have to provide an adapter for every node sub-type
to return the sub-type from the generic type.

BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.  */
package dom

import "github.com/npillmayer/gotype/core/config/tracing"

// T will return a tracer. We are tracing to the EngineTracer
func T() tracing.Trace {
	return tracing.EngineTracer
}
