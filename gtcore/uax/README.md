# Breaking Unicode Text into Segments

Text processing applications need to segment text into pieces. Segments may be

* words
* sentences
* paragraphs

and so on. For western languages this is often not too hard of a problem, but it may become an involved endeavor if you consider Arabic or Asian languages. From a typographic viewpoint some of these languages present serious challenges for correct segmenting. It becomes even more involved if emojis are considered.

Text Segmentation in Golang
---------------------------

There exist a number of Unicode standards describing best practices for text segmentation. Unfortunately, implementations in Go seem to be sparse. Marcel van Lohuizen from the Go Core Team seems be working on text segmenting, but with low priority. In the long run, it will be best to wait for the standard library to include functions for text segmentation. However, for now I will implement my own.

Implementation
--------------

For some information take a look at the
[GoDoc page for uax](https://godoc.org/github.com/npillmayer/gotype/gtcore/uax).
More information to come, as soon as I collect more real-life experience
with breaking real-life multi-lingual text with these algorithms.
