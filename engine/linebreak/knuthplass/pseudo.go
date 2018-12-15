package knuthplass

import "fmt"

/* Data structure for the testing of line breaking algorithms.
 */

// ---------------------------------------------------------------------------

type pseudoBead struct {
	text string // implements Bead
}

func (pb *pseudoBead) String() string {
	var s string = fmt.Sprintf("<bead \"%s\" (%s)>", pb.text, beadtype(pb.BeadType()))
	return s
}

func (pb *pseudoBead) BeadType() int8 {
	if pb.text == " " || pb.text == "\n" || pb.text == "@" {
		return GlueType
	} else if pb.text == "" {
		return KernType
	} else {
		return BoxType
	}
}

func (pb *pseudoBead) Text() string {
	return pb.text
}

func (pb *pseudoBead) Width() WSS { // width, w-shrink, w+stretch
	w := int64(len(pb.text))
	if pb.text == " " || pb.text == "\n" {
		return WSS{w, w, w + 2}
	} else if pb.text == "@" {
		return WSS{0, 0, 1000}
	} else {
		return WSS{w, w, w}
	}
}

func (pb *pseudoBead) Dimens() (w int64, h int64, d int64) {
	return int64(len(pb.text)), 1, 0
}

func newPseudoBead(text string) *pseudoBead {
	return &pseudoBead{text: text}
}

var _ Bead = &pseudoBead{}

// ---------------------------------------------------------------------------

type pseudoBeading struct {
	text string
}

func (pbg *pseudoBeading) GetCursor(cursor BeadingCursor) BeadingCursor {
	bc := &pseudoBeadingCursor{beading: pbg}
	if c, ok := cursor.(*pseudoBeadingCursor); ok {
		if c != nil {
			bc.pos = c.pos
			bc.end = c.end
		}
	}
	return bc
}

type pseudoBeadingCursor struct { // provisional pseudo-implementation, implements Beading
	beading *pseudoBeading
	pos     int
	end     int
}

func (c *pseudoBeadingCursor) String() string {
	return fmt.Sprintf("<pos %d-%d>", c.pos, c.end)
}

func (c *pseudoBeadingCursor) ID() int64 {
	return int64(c.end)
}

func (c *pseudoBeadingCursor) GetBead() Bead {
	//fmt.Printf("text = %s\n", pbg.text[pbg.end:pbg.pos+20])
	t := c.beading.text[c.pos:c.end]
	return newPseudoBead(t)
}

func (c *pseudoBeadingCursor) Advance() bool {
	if c.end >= len(c.beading.text) {
		return false
	}
	c.pos = c.end
	e := c.end
	if c.beading.text[c.pos] == ' ' {
		for e < len(c.beading.text) && c.beading.text[e] == ' ' { // TODO make this: looking for penalty
			e++ // TODO make this unicode-compliant
		}
	} else {
		for e < len(c.beading.text) && c.beading.text[e] != ' ' { // TODO make this: looking for penalty
			e++ // TODO make this unicode-compliant
		}
	}
	c.end = e
	return true
}

func newPseudoBeading(text string) *pseudoBeading {
	return &pseudoBeading{text: text + "@"}
}

var _ Beading = &pseudoBeading{}
var _ BeadingCursor = &pseudoBeadingCursor{}

// ---------------------------------------------------------------------------
