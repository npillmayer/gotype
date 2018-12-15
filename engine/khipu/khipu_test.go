package khipu

import (
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/dimen"
	"github.com/npillmayer/gotype/core/parameters"
)

func Test0(t *testing.T) {
	t.Log("setting core tracer")
	CT = gologadapter.New()
	CT.SetTraceLevel(tracing.LevelDebug)
	tracing.CoreTracer = CT
	CT.Infof("Core tracer alive")
	CT.SetTraceLevel(tracing.LevelError)
}

func TestDimen(t *testing.T) {
	if dimen.BP.String() != "65536sp" {
		t.Error("a big point BP should be 65536 scaled points SP")
	}
}

func TestKhipu(t *testing.T) {
	khipu := NewKhipu()
	khipu.AppendKnot(NewKnot(KTKern)).AppendKnot(NewKnot(KTGlue))
	khipu.AppendKnot(NewTextBox("Hello"))
	t.Logf("khipu = %s\n", khipu.String())
	if khipu.Length() != 3 {
		t.Errorf("Length of khipu should be 3")
	}
}

func TestBreaking1(t *testing.T) {
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	khipu := KnotEncode(strings.NewReader("Hello world!"), nil, regs)
	if khipu.Length() != 9 {
		t.Errorf("khipu length is %d, should be 9", khipu.Length())
	}
}
