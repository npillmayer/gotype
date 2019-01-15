
Typesetting Engine
--------------------------------------------------

If you think about it: a typesetter using the HTML/CSS box model is
effectively a browser with output type PDF.
Browsers are large and complex pieces of code, a fact that implies that
we should seek out where to reduce complexity.

A good explanation of styling may be found in a superb
[blog entry](https://hacks.mozilla.org/2017/08/inside-a-super-fast-css-engine-quantum-css-aka-stylo/)
by Lin Clark.
about the Firefox styling engine.
 
There is not very much open source Go code around for supporting us
in implementing a styling and layout engine, except the great work of
[Cascadia](https://godoc.org/github.com/andybalholm/cascadia).
Therefore we will have to compromise
on many feature in order to complete this in a realistic time frame.
For a reminder of why that is, refer to
[this video](https://www.youtube.com/watch?v=S68fcV09nGQ).

From a high level view the creation of documents from the internal
representation looks something like this:

<div style="width:480px;padding:5px;padding-bottom:10px">
<img alt="gotype document creation flow overview" src="http://pillmayer.com/TySE-Notes/img/TySE-Concurrency-Overview.svg" width="480px">
</div>

