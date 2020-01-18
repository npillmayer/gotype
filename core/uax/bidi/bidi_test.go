package bidi

import (
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
)

func TestParserGenerate(t *testing.T) {
	gtrace.SyntaxTracer = gologadapter.New()
	gtrace.SyntaxTracer.SetTraceLevel(tracing.LevelDebug)
	bidi := NewGrammar()
	bidi.Parse("a b")
}
