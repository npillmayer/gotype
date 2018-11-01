package gtcore

/*
Interfaces and methods to create "beadings" of typesetting items.
Beadings are lists of items, such as glyphs, kerns, glue, etc.
Sometimes we will call these kinds of items "beads".

Beadings are the input for linebreakers.
*/

import (
	"github.com/npillmayer/gotype/gtcore/parameters"
)

type BeadingDriver interface {
	Thread(text string, regs *parameters.TypesettingRegisters)
}
