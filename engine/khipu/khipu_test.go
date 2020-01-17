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
	//gtrace.CoreTracer = gologadapter.New()
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
	kh := NewKhipu()
	kh.AppendKnot(NewKnot(KTKern)).AppendKnot(NewKnot(KTGlue))
	kh.AppendKnot(NewTextBox("Hello"))
	t.Logf("khipu = %s\n", kh.String())
	if kh.Length() != 3 {
		t.Errorf("Length of khipu should be 3")
	}
}

func TestBreaking1(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := KnotEncode(strings.NewReader("Hello World "), nil, regs)
	if kh.Length() != 10 {
		t.Logf("khipu = %s", kh)
		t.Errorf("khipu length is %d, should be 10", kh.Length())
	}
}

func TestBreaking2(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := KnotEncode(strings.NewReader("The quick !"), nil, regs)
	if kh.Length() != 10 {
		t.Logf("khipu = %s", kh)
		t.Errorf("khipu length is %d, should be 10", kh.Length())
	}
}

func TestText(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelInfo)
	text := "The quick brown fox jumps over the lazy dog!"
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := KnotEncode(strings.NewReader(text), nil, regs)
	out := kh.Text(0, kh.Length())
	if out != text {
		t.Logf("Text: %s", out)
		t.Errorf("output text != input text")
	}
}

func TestExHyphen(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	gtrace.CoreTracer.SetTraceLevel(tracing.LevelDebug)
	text := "lime-tree"
	regs := parameters.NewTypesettingRegisters()
	regs.Push(parameters.P_MINHYPHENLENGTH, 3)
	kh := KnotEncode(strings.NewReader(text), nil, regs)
	out := kh.Text(0, kh.Length())
	if out != text {
		t.Logf("Text: %s", out)
		t.Errorf("output text != input text")
	}
}
