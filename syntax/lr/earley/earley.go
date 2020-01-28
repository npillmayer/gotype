/*
Package earley will some day implement an Earley-Parser.

There already exists a solution in

	https://github.com/jakub-m/gearley

which is based on the nice Earley-introduction from

	http://loup-vaillant.fr/tutorials/earley-parsing/

(which boasts an implementation in Lua and OCaml) */
package earley

import (
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
)

// T traces to the global syntax tracer.
func T() tracing.Trace {
	return gtrace.SyntaxTracer
}
