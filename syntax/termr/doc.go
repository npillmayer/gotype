package gospel

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
)

// T traces to the global SyntaxTracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}
