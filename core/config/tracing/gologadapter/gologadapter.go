/*
Package gologadapter implements tracing with the default Go logger.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing implement concrete tracers. Package
gologadapter uses the default Go logging mechanism.

BSD License

Copyright (c) 2017–18, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package gologadapter

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/npillmayer/gotype/core/config/tracing"
)

// Tracer is our adapter implementation which implements interface
// tracing.Trace, using a Go standard logger.
type Tracer struct {
	log   *log.Logger
	p     string
	level tracing.TraceLevel
}

var logLevelPrefix = []string{"ERROR ", "INFO  ", "DEBUG "}

// New creates a new Tracer instance based on a logrus logger.
func New() tracing.Trace {
	return &Tracer{log.New(os.Stderr, logLevelPrefix[0], log.Ltime),
		"",
		tracing.LevelError,
	}
}

// NewAdapter creates an adapter (i.e., factory for tracing.Trace) to
// be used to initialize (global) tracers.
func GetAdapter() tracing.Adapter {
	return New
}

// Interface Trace
func (t *Tracer) P(key string, val interface{}) tracing.Trace {
	t.p = fmt.Sprintf("[%s=%v] ", key, val)
	return t
}

// Interface Trace
func (t *Tracer) output(l tracing.TraceLevel, s string, args ...interface{}) {
	t.log.SetPrefix(logLevelPrefix[int(l)])
	if t.p != "" {
		t.log.Println(fmt.Sprintf(t.p+s, args...))
		t.p = ""
	} else {
		t.log.Printf(s, args...)
	}
}

// Interface Trace
func (t *Tracer) Debugf(s string, args ...interface{}) {
	if t.level < tracing.LevelDebug {
		return
	}
	t.output(tracing.LevelDebug, s, args...)
}

// Interface Trace
func (t *Tracer) Infof(s string, args ...interface{}) {
	if t.level < tracing.LevelInfo {
		return
	}
	t.output(tracing.LevelInfo, s, args...)
}

// Interface Trace
func (t *Tracer) Errorf(s string, args ...interface{}) {
	if t.level < tracing.LevelError {
		return
	}
	t.output(tracing.LevelError, s, args...)
}

// Interface Trace
func (t *Tracer) SetTraceLevel(l tracing.TraceLevel) {
	t.level = l
	t.log.SetPrefix(logLevelPrefix[int(l)])
}

// Interface Trace
func (t *Tracer) GetTraceLevel() tracing.TraceLevel {
	return t.level
}

// Interface Trace
func (t *Tracer) SetOutput(writer io.Writer) {
	t.log.SetOutput(writer)
}
