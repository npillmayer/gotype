package khipu

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/configtestadapter"
	"github.com/npillmayer/gotype/core/config/gconf"
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/core/parameters"
)

func init() {
	gconf.Initialize(configtestadapter.New())
	gtrace.CoreTracer = gotestingadapter.New()
}

func TestDimen(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	if dimen.BP.String() != "65536sp" {
		t.Error("a big point BP should be 65536 scaled points SP")
	}
}

func TestKhipu(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	khipu := NewKhipu()
	khipu.AppendKnot(NewKnot(KTKern)).AppendKnot(NewKnot(KTGlue))
	khipu.AppendKnot(NewTextBox("Hello"))
	t.Logf("khipu = %s\n", khipu.String())
	if khipu.Length() != 3 {
		t.Errorf("Length of khipu should be 3")
	}
}

func TestBreaking1(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	khipu := KnotEncode(strings.NewReader("Hello World!"), nil, regs)
	if khipu.Length() != 9 {
		t.Errorf("khipu length is %d, should be 9", khipu.Length())
	}
}
