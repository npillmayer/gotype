package sppf

import (
	"fmt"
	"testing"
	"text/scanner"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"

	"github.com/npillmayer/gotype/syntax/lr"
)

func TestSigma(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := lr.NewGrammarBuilder("G")
	b.LHS("S").T("<", '<').N("A").N("Z").T(">", '>').End()
	g, _ := b.Grammar()
	s1 := makeSym(g.SymbolByName("A")).spanning(1, 8)
	s2 := makeSym(g.SymbolByName("Z")).spanning(8, 9)
	rhs := []*SymbolNode{s1, s2}
	t.Logf("rhs=%v", rhs)
	sigma := rhsSignature(rhs, 0)
	if sigma != 105544 {
		t.Errorf("sigma expected to be 105544, is %d", sigma)
	}
}

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

// S' ⟶ S
// S  ⟶ A
// A  ⟶ a
func TestTraverse(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	b := lr.NewGrammarBuilder("G")
	r1 := b.LHS("S").N("A").End()
	r2 := b.LHS("A").T("a", scanner.Ident).End()
	G, err := b.Grammar()
	if err != nil {
		t.Error(err)
	}
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	G.Dump()
	f := NewForest()
	a := f.AddTerminal(r2.RHS()[0], 0)
	A := f.AddReduction(r2.LHS, 2, []*SymbolNode{a})
	S := f.AddReduction(r1.LHS, 1, []*SymbolNode{A})
	f.AddReduction(G.SymbolByName("S'"), 0, []*SymbolNode{S})
	if f.Root() == nil {
		t.Errorf("Expected root node S', is nil")
	}
	l := makeListener(G, t)
	c := f.SetCursor(nil, nil)
	c.TopDown(l, LtoR, Continue)
	if !l.(*L).isBack {
		t.Errorf("Exit(S') has not been called")
	}
	if l.(*L).a.Name != "a" {
		t.Errorf("Terminal(a) has not been called")
	}
}

// ---------------------------------------------------------------------------

func makeListener(G *lr.Grammar, t *testing.T) Listener {
	return &L{G: G, t: t}
}

type L struct {
	G      *lr.Grammar
	t      *testing.T
	isBack bool
	a      *lr.Symbol
}

func (l *L) EnterRule(sym *lr.Symbol, rhs []*RuleNode, span lr.Span, level int) bool {
	if sym.IsTerminal() {
		return false
	}
	l.t.Logf("+ enter %v", sym)
	return true
}
func (l *L) ExitRule(sym *lr.Symbol, rhs []*RuleNode, span lr.Span, level int) interface{} {
	if sym.Name == "S'" {
		l.isBack = true
	}
	l.t.Logf("- exit %v", sym)
	return nil
}

func (l *L) Terminal(tokval int, token interface{}, span lr.Span, level int) interface{} {
	tok := l.G.Terminal(tokval)
	l.a = tok
	l.t.Logf("  terminal=%s", tok.Name)
	return tok
}

func (l *L) Conflict(sym *lr.Symbol, rule int, span lr.Span, level int) (int, error) {
	l.t.Error("did not expect conflict")
	return 0, fmt.Errorf("Conflict at symbol %v", sym)
}
