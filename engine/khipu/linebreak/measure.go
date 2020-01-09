package linebreak

import (
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/engine/khipu"
)

type GlyphMeasure interface {
	Measure(knot khipu.Knot) dimen.Dimen
}
