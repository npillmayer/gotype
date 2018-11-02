package gtcore

/*
Interfaces and methods to create "bead chains" of typesetting items.
Bead chains are lists of items, such as glyphs, kerns, glue, etc.
Sometimes we will call these kinds of items "beads".

Bead chains are the input for linebreakers.

Finding possible line-wrap points (Unicode spec):
http://unicode.org/reports/tr14/

JavaScript implementation: https://github.com/foliojs/linebreak (basiert
  auf deprecated pair table)
Python: https://uniseg-python.readthedocs.io/en/latest/ bzw.
https://bitbucket.org/emptypage/uniseg-python/src/e4077d17d026c36999b89c10081a85b219e1fb7b/uniseg/?at=default

Golang package unicode provides all sorts of code-point ranges:
https://golang.org/pkg/unicode/

Moreover: https://godoc.org/golang.org/x/text/unicode

Moreover: https://godoc.org/golang.org/x/text/width


A: Finding feasible break positions

  (1) Mandatory breaks + prohibiting no-break points

  (2) Natural text wrap (words in many scripts, syllables/character in east asia, etc.)
      -> Unicode UA#14 Line Breaking (http://unicode.org/reports/tr14/)
      Algorithm: http://unicode.org/reports/tr14/#Algorithm

  (3) Bidi

  (4) Hyphenation
      Lliang patterns + language-specific code

  (5) Translate feasible breakpoints to penalties, glue and dicretionaries

B: Deciding break positions

  (1) Shape text -> Glyphs
      + alternative glyphs (end-of-line condensed in Arabic, etc.)

  (2) Translate everything to node list

  (3) Apply line-breaking algorithm (simple, K&P, etc.)

C: Justify text

  (1) Set glue


// ----------------------------------------------------------------------

Khipus: https://www.sapiens.org/technology/khipu-incas-knotty-history/

the khipukamayuqs (Quechua for “knot-makers/animators”) encoded administrative
data such as census figures and tax allocation in the twisted strings of these
ancient spreadsheets.

*/

type Khipukamayuq interface {
	//KnotEncode(text string, pipeline *TypesettingPipeline, regs *p.TypesettingRegisters) *Khipu
}
