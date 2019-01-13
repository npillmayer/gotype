package testimages

import (
	"testing"

	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/path"
)

var T tracing.Trace = tracing.GraphicsTracer

func TestEnvironment(t *testing.T) {
	tracing.GraphicsTracer = gologadapter.New()
	T.SetTraceLevel(tracing.LevelDebug)
}

func TestEmptyPath1(t *testing.T) {
	pic := gfx.NewPicture("empty", 100, 100, "png")
	pic.Shipout()
}

func TestPath1(t *testing.T) {
	pic := gfx.NewPicture("path1", 100, 100, "png")
	p, controls := path.Nullpath().Knot(path.P(0, 0)).Curve().Knot(path.P(50, 50)).Curve().
		Knot(path.P(100, 65)).End()
	controls = path.FindHobbyControls(p, controls)
	pic.Draw(gfx.NewDrawablePath(p, controls))
	pic.Shipout()
}

func TestPath2(t *testing.T) {
	pic := gfx.NewPicture("path2", 100, 100, "png")
	p, controls := path.Nullpath().Knot(path.P(10, 50)).Curve().Knot(path.P(50, 90)).Curve().
		Knot(path.P(90, 50)).End()
	controls = path.FindHobbyControls(p, controls)
	pic.Draw(gfx.NewDrawablePath(p, controls))
	pic.Shipout()
}

func TestPath3(t *testing.T) {
	pic := gfx.NewPicture("path3", 100, 100, "png")
	p, controls := path.Nullpath().Knot(path.P(10, 50)).Curve().Knot(path.P(50, 90)).Curve().
		Knot(path.P(90, 50)).Curve().Knot(path.P(50, 10)).Curve().Cycle()
	controls = path.FindHobbyControls(p, controls)
	pic.Draw(gfx.NewDrawablePath(p, controls))
	pic.Shipout()
}
