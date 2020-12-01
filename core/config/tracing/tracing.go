/*
Package tracing defines an interface used by all other packages.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing will implement concrete tracers.
Additionally, this package will provide helpers for more descriptive
logging output, to be used with any concrete tracing/logging class.


BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package tracing

import (
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// TraceLevel is a type for leveled tracing.
// All concrete Tracer implementations will support trace-levels.
type TraceLevel uint8

// We support three trace levels.
const (
	LevelError TraceLevel = iota
	LevelInfo
	LevelDebug
)

func (tl TraceLevel) String() string {
	switch tl {
	case LevelDebug:
		return "Debug"
	case LevelInfo:
		return "Info"
	case LevelError:
		return "Error"
	}
	return "<unknown>"
}

// TraceLevelFromString will find a trace level from a string.
// It will recognize "Debug", "Info" and "Error". Default is
// LevelDebug.
func TraceLevelFromString(sl string) TraceLevel {
	switch sl {
	case "Debug":
		return LevelDebug
	case "Info":
		return LevelInfo
	case "Error":
		return LevelError
	}
	return LevelInfo // default
}

// Trace is an interface to be implemented by a concrete tracing adapter.
// For examples refer to the sub-packages of package tracing.
//
// Tracers should support parameter/field tracing given by P(...).
// An example would be
//
//    tracer.P("mycontext", "value").Debugf("message within context")
//
// Tracers should be prepared to trace to console as well as to a file.
// By convention, no newlines at the end of tracing messages will be passed
// by clients.
type Trace interface {
	P(string, interface{}) Trace // field tracing
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	SetTraceLevel(TraceLevel)
	GetTraceLevel() TraceLevel
	SetOutput(io.Writer) // route tracing output to a writer
}

// Tracefile is the gloabl file where tracing goes to.
// If tracing goes to a file (globally), variable Tracefile should
// point to it. It need not be set if tracing goes to console.
var Tracefile *os.File

// Adapter is a factory function to create a virgin Trace instance.
type Adapter func() Trace

// With prepares to dump a data structure to a Trace.
// t may not be nil.
//
// Usage:
//
//     tracing.With(mytracer).Dump(anObject)
//
// Dumping uses https://github.com/davecgh/go-spew .
// Dump() is in level Debug.
func With(t Trace) dumper {
	return dumper{&t}
}

// Helper type for dumping of objects.  Created by calls to With().
type dumper struct {
	tracer *Trace
}

// Dump dumps an object using a tracer, in level Debug.
//
// d may not be nil.
func (d dumper) Dump(name string, obj interface{}) {
	if (*d.tracer).GetTraceLevel() >= LevelDebug {
		str := spew.Sdump(obj)
		(*d.tracer).Debugf(name + " = " + str)
	}
}
