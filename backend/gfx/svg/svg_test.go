package svg

import (
	"image/color"
	"os"
	"testing"

	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/testadapter"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/path"
)

func TestInit0(t *testing.T) {
	config.Initialize(testadapter.New())
	tracing.GraphicsTracer.SetTraceLevel(tracing.LevelDebug)
}

func TestSimple(t *testing.T) {
	f, err := os.Create("test.svg")
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()
	pic := NewPicture(f, 100, 100)
	p, controls := path.Nullpath().Knot(path.P(10, 50)).Curve().Knot(path.P(50, 90)).Curve().
		Knot(path.P(90, 50)).Curve().Knot(path.P(50, 10)).Curve().Cycle()
	controls = path.FindHobbyControls(p, controls)
	pic.AddContour(gfx.NewDrawablePath(p, controls), 2, color.Black, nil)
	pic.Shipout()
}
