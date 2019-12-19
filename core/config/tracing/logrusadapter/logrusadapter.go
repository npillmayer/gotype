/*
Package logrusadapter implements tracing with the logrus logger.

Tracing/logging is a cross cutting concern. Relying on a specific package
for such a low-level task will create too tight a coupling—more abstract
classes/packages are infected with log classes/packages.

Sub-packages of tracing implement concrete tracers. Package
logrus uses "github.com/sirupsen/logrus" as the means for tracing.

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package logrusadapter

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/npillmayer/gotype/core/config/tracing"
)

// Tracer is our adapter implementation which implements interface
// tracing.Trace, using a logrus logger.
type Tracer struct {
	log *logrus.Logger
	p   *logrus.Entry
}

// New creates a new Tracer instance based on a logrus logger.
func New() tracing.Trace {
	return &Tracer{logrus.New(), nil}
}

// NewAdapter creates an adapter (i.e., factory for tracing.Trace) to
// be used to initialize (global) tracers.
func GetAdapter() tracing.Adapter {
	return New
}

// Interface tracing.Trace
func (t *Tracer) P(key string, val interface{}) tracing.Trace {
	t.p = t.log.WithField(key, val)
	return t
}

// Interface tracing.Trace
func (t *Tracer) Debugf(s string, args ...interface{}) {
	if t.p != nil {
		t.p.Debugf(s, args...)
		t.p = nil
	} else {
		t.log.Debugf(s, args...)
	}
}

// Interface tracing.Trace
func (t *Tracer) Infof(s string, args ...interface{}) {
	if t.p != nil {
		t.p.Infof(s, args...)
		t.p = nil
	} else {
		t.log.Infof(s, args...)
	}
}

// Interface tracing.Trace
func (t *Tracer) Errorf(s string, args ...interface{}) {
	if t.p != nil {
		t.p.Errorf(s, args...)
		t.p = nil
	} else {
		t.log.Errorf(s, args...)
	}
}

// Interface tracing.Trace
func (t *Tracer) SetTraceLevel(l tracing.TraceLevel) {
	t.log.SetLevel(translateTraceLevel(l))
}

// Interface tracing.Trace
func (t *Tracer) GetTraceLevel() tracing.TraceLevel {
	return translateLogLevel(t.log.Level)
}

// Interface tracing.Trace
func (t *Tracer) SetOutput(writer io.Writer) {
	t.log.Out = writer
	t.log.Formatter = &logrus.TextFormatter{}
}

func translateLogLevel(l logrus.Level) tracing.TraceLevel {
	switch l {
	case logrus.DebugLevel:
		return tracing.LevelDebug
	case logrus.InfoLevel:
		return tracing.LevelInfo
	case logrus.ErrorLevel:
		return tracing.LevelError
	}
	return tracing.LevelDebug
}

func translateTraceLevel(l tracing.TraceLevel) logrus.Level {
	switch l {
	case tracing.LevelDebug:
		return logrus.DebugLevel
	case tracing.LevelInfo:
		return logrus.InfoLevel
	case tracing.LevelError:
		return logrus.ErrorLevel
	}
	return logrus.DebugLevel
}
