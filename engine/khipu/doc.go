/*
Package khipu is about encoding text into typesetting items.

"Khipu were recording devices fashioned
from strings historically used by a number of cultures in the region of
Andean South America.
Khipu is the word for "knot" in Cusco Quechua.
A khipu usually consisted of cotton or camelid fiber strings. The Inca
people used them for collecting data and keeping records, monitoring tax
obligations, properly collecting census records, calendrical information,
and for military organization. The cords stored numeric and other values
encoded as knots, often in a base ten positional system. A khipu could
have only a few or thousands of cords."
––Excerpt from a Wikipedia article about khipus

The Khipukamayuqs (Quechua for “knot-makers”) were the scribes of those
times, tasked with encoding tax figures and other administrative
information in knots.
We will use this analogy to call typesetting items "khipus" or "knots",
and objects which produce khipus will be "Khipukamayuqs".
Knots implement items for typesetting paragraphs. We will use a
box-and-glue model, the various knot types more or less implementing
the corresponding node types from the TeX typesetting system.

A Khipukamayuqs is part of a typsetting pipeline and will transform
text into khipus.
Khipus are the input for linebreakers. The overall process of creating
them and the interaction with line breaking looks like this:

Create Khipus from Text

(1) Normalize Unicode text

	https://godoc.org/golang.org/x/text/unicode/norm

(2) Find natural text wrap opportunities
(words in many scripts, syllables/character in East Asia, etc.)

	https://godoc.org/github.com/npillmayer/gotype/core/uax

(3) Bi-directional text

	https://godoc.org/golang.org/x/text/unicode/bidi
	https://www.w3.org/International/articles/inline-bidi-markup/

(4) Hyphenation:
Lliang patterns + language-specific code

	https://godoc.org/github.com/npillmayer/gotype/core/hyphenation

(5) Translate feasible breakpoints to penalties, glue and discretionaries

    https://wiki.apache.org/xmlgraphics-fop/KnuthsModel

(6) Shape text -> Glyphs
+ alternative glyphs (end-of-line condensed in Arabic, etc.)

	http://behdad.org/text/

At this point, text has been fully converted to khipus.

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package khipu

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
)

// CT traces to the core-tracer.
func CT() tracing.Trace {
	return gtrace.CoreTracer
}
