/*
Package style will some day provide functionality for CSS styling.

Status

Nothing to find here yet – development will start soon.

Overview

HTMLbook is the core DOM of our documents.
Background for this decision can be found under
https://www.balisage.net/Proceedings/vol10/print/Kleinfeld01/BalisageVol10-Kleinfeld01.html
and http://radar.oreilly.com/2013/09/html5-is-the-future-of-book-authorship.html
For an in-depth description of HTMLbook please refer to
https://oreillymedia.github.io/HTMLBook/.

We strive to separate content from presentation. In typesetting, this is
probably an impossible claim, but we'll try anyway. Presentation
is governed with CSS (Cascading Style Sheets). CSS uses a box model more
complex than TeX's, which is well described here:

   https://developer.mozilla.org/en-US/docs/Learn/CSS/Introduction_to_CSS/Box_model

If you think about it: a typesetter using the HTML/CSS box model is
effectively a browser with output type PDF.
Browsers are large and complex pieces of code, a fact that implies that
we should seek out where to reduce complexity.

A good explanation of styling may be found in

   https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/

CSSOM is the "CSS Object Model", similar to the DOM for HTML.
There is not very much open source Go code around for supporting us
in implementing a styling engine, except the great work of
https://godoc.org/github.com/andybalholm/cascadia.
Therefore we will have to compromise
on many feature in order to complete this in a realistic time frame.
For a reminder of why that is, refer to
https://www.youtube.com/watch?v=S68fcV09nGQ .

Further to consider:

   https://github.com/aymerick/douceur
   https://golanglibs.com/search?q=css+parser&sort=top

   https://www.mediaevent.de/xhtml/style.html

The styling component is difficult to document/describe without
diagrams. Think about documenting with https://github.com/robertkrimen/godocdown.

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
package style
