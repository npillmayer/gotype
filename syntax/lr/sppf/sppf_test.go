package sppf

import (
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"

	"github.com/npillmayer/gotype/syntax/lr"
)

// S' ⟶ S
// S  ⟶ A
// A  ⟶ a
func TestSPPFInsert(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := lr.NewGrammarBuilder("G")
	b.LHS("S").N("A").End()
	r2 := b.LHS("A").T("a", scanner.Ident).End()
	g, err := b.Grammar()
	if err != nil {
		t.Error(err)
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	g.Dump()
	f := NewForest()
	a := f.AddTerminal(r2.RHS()[0], 0)
	A := r2.LHS
	R := f.AddReduction(A, 2, []*SymbolNode{a})
	t.Logf("node A=%v for rule %v", R, g.Rule(2))
	if R == nil {
		t.Errorf("Expected symbol node A=%v for rule %v", R, g.Rule(2))
	}
}

// S' ⟶ S
// S  ⟶ A | B
// A  ⟶ a
// B  ⟶ a
func TestSPPFAmbiguous(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := lr.NewGrammarBuilder("G")
	b.LHS("S").N("A").End()
	b.LHS("S").N("B").End()
	b.LHS("A").T("a", scanner.Ident).End()
	b.LHS("B").T("a", scanner.Ident).End()
	g, err := b.Grammar()
	if err != nil {
		t.Error(err)
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	g.Dump()
	//f := NewForest()
}
