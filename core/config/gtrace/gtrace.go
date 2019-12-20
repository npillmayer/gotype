/*
Package gtrace lets clients set up a set of global tracers.
It infects clients with the semantics of tracing for the gotype-application.
However, the imposed overhead is very slim, i.e. boils down to having
some global variables declared (but not necessarily used):

	CommandTracer     : Tracing interactive or batch commands from users
	CoreTracer        : Tracing application core
	ScriptingTracer   : Tracing embedded scripting host(s)
	InterpreterTracer : Tracing DSL interpreting
	SyntaxTracer      : Tracing lexing/parsing of DSLs
	GraphicsTracer    : Tracing graphics routines
	EngineTracer      : Tracing web-engine routines
	EquationsTracer   : Tracing arithmetic

Clients are free to selectively use any of these tracers.
They are initially set up to do nothing (no-ops).
Clients will use the configuration package to set it up in a
meaningful way.

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
package gtrace

import (
	"errors"
	"io"

	"github.com/npillmayer/gotype/core/config/tracing"
)

// This is the set of standard module tracers for our application.
//
// All tracers are set up to be no-ops, initially.
// This approach poses a little burden on (selective) clients, but is
// useful for de-coupling the various packages and modules from the
// tracing/logging mechanism.
var (
	EquationsTracer   = NoOpTrace
	InterpreterTracer = NoOpTrace
	SyntaxTracer      = NoOpTrace
	CommandTracer     = NoOpTrace
	GraphicsTracer    = NoOpTrace
	ScriptingTracer   = NoOpTrace
	CoreTracer        = NoOpTrace
	EngineTracer      = NoOpTrace
)

// NoOpTrace is a void Trace. Initially, all tracers will be set up to be no-ops.
// Clients will have to configure concrete tracing backends, usually by calling
// application configuration with a tracing adapter.
var NoOpTrace tracing.Trace = nooptrace{}

type nooptrace struct{}

func (nt nooptrace) Debugf(string, ...interface{})       {}
func (nt nooptrace) Infof(string, ...interface{})        {}
func (nt nooptrace) Errorf(string, ...interface{})       {}
func (nt nooptrace) SetTraceLevel(tracing.TraceLevel)    {}
func (nt nooptrace) GetTraceLevel() tracing.TraceLevel   { return tracing.LevelError }
func (nt nooptrace) SetOutput(io.Writer)                 {}
func (nt nooptrace) P(string, interface{}) tracing.Trace { return nt }

// CreateTracers creates all global tracers, given a function to
// create a concrete Trace instance.
func CreateTracers(newTrace tracing.Adapter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unable to create global tracers")
		}
	}()
	EquationsTracer = newTrace()
	InterpreterTracer = newTrace()
	SyntaxTracer = newTrace()
	CommandTracer = newTrace()
	GraphicsTracer = newTrace()
	ScriptingTracer = newTrace()
	CoreTracer = newTrace()
	EngineTracer = newTrace()
	return
}

// Mute sets all global tracers to LevelError.
func Mute() {
	InterpreterTracer.SetTraceLevel(tracing.LevelError)
	CommandTracer.SetTraceLevel(tracing.LevelError)
	EquationsTracer.SetTraceLevel(tracing.LevelError)
	SyntaxTracer.SetTraceLevel(tracing.LevelError)
	GraphicsTracer.SetTraceLevel(tracing.LevelError)
	ScriptingTracer.SetTraceLevel(tracing.LevelError)
	CoreTracer.SetTraceLevel(tracing.LevelError)
	EngineTracer.SetTraceLevel(tracing.LevelError)
}
