/*
Package gconf initializes a global application configuration.

Configuration

All configuration is started explicitely with a call to Initialize().
There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

Tracing

During configuration all global tracers are set up.
To use a concrete logging implementation,
clients will have to use/implement an adapter to tracing.Trace (please
refer to the documentation for package tracing as well as to implementations
of adapters, e.g. for Go log and for logrus).

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
package gconf

import (
	"os"

	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
)

// The global configuration is initially set up to be a no-op.
var globalConf config.Configuration = noconfig{}

type noconfig struct{}

func (nc noconfig) InitDefaults()               {}
func (nc noconfig) IsSet(key string) bool       { return false }
func (nc noconfig) GetString(key string) string { return "" }
func (nc noconfig) GetInt(key string) int       { return 0 }
func (nc noconfig) GetBool(key string) bool     { return false }
func (nc noconfig) IsInteractive() bool         { return false }

// --- Wire Golabl Configuration ----------------------------------------------------

// Initialize is the top level function for setting up the
// application configuration.
// It will call InitDefaults() on the Configuration passed as an argument, and
// make the Configuration available globally.
// Functions in this package will serve as a facade to the Configuration.
func Initialize(conf config.Configuration) {
	globalConf = conf
	globalConf.InitDefaults()
	InitTracing(config.GetAdapterFromConfiguration(conf))
}

// InitTracing sets up all the global module tracers, reading trace levels
// and tracing destinations from the application configuration.
//
// InitTracing is usually not called directly, but rather called by Initialize().
func InitTracing(adapter tracing.Adapter) {
	gtrace.CreateTracers(adapter)
	ConfigureTracing("")
}

// ConfigureTracing sets up the global tracers using default configuration values.
//
// It is exported as it may be useful in testing scenarios.
func ConfigureTracing(inputfilename string) {
	SetDefaultTracingLevels() // set default trace levels from configuration
	if GetBool("tracingonline") {
		if inputfilename != "" {
			file, err := os.Create("__gotype.log")
			if err != nil {
				gtrace.CommandTracer.Errorf("cannot open tracefile, tracing to stderr")
			} else {
				tracing.Tracefile = file
			}
		}
		if tracing.Tracefile != nil {
			gtrace.CommandTracer.SetOutput(tracing.Tracefile)
			gtrace.EquationsTracer.SetOutput(tracing.Tracefile)
			gtrace.SyntaxTracer.SetOutput(tracing.Tracefile)
			gtrace.InterpreterTracer.SetOutput(tracing.Tracefile)
			gtrace.GraphicsTracer.SetOutput(tracing.Tracefile)
			gtrace.ScriptingTracer.SetOutput(tracing.Tracefile)
			gtrace.CoreTracer.SetOutput(tracing.Tracefile)
			gtrace.EngineTracer.SetOutput(tracing.Tracefile)
		}
	}
	gtrace.InterpreterTracer.P("level", gtrace.InterpreterTracer.GetTraceLevel()).Infof("Interpreter-Trace is alive")
	gtrace.CommandTracer.P("level", gtrace.CommandTracer.GetTraceLevel()).Infof("Command-Trace is alive")
	gtrace.EquationsTracer.P("level", gtrace.EquationsTracer.GetTraceLevel()).Infof("Equations-Trace is alive")
	gtrace.SyntaxTracer.P("level", gtrace.SyntaxTracer.GetTraceLevel()).Infof("Syntax-Trace is alive")
	gtrace.GraphicsTracer.P("level", gtrace.GraphicsTracer.GetTraceLevel()).Infof("Graphics-Trace is alive")
	gtrace.ScriptingTracer.P("level", gtrace.ScriptingTracer.GetTraceLevel()).Infof("Scripting-Trace is alive")
	gtrace.CoreTracer.P("level", gtrace.CoreTracer.GetTraceLevel()).Infof("Core-Trace is alive")
	gtrace.EngineTracer.P("level", gtrace.CoreTracer.GetTraceLevel()).Infof("Engine-Trace is alive")
}

// SetDefaultTracingLevels sets all global tracers to their default trace levels,
// read from the application configuration.
func SetDefaultTracingLevels() {
	gtrace.InterpreterTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracinginterpreter")))
	gtrace.CommandTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingcommands")))
	gtrace.EquationsTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingequations")))
	gtrace.SyntaxTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingsyntax")))
	gtrace.GraphicsTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracinggraphics")))
	gtrace.ScriptingTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingscripting")))
	gtrace.CoreTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingcore")))
	gtrace.EngineTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingengine")))
}

// --- Gobal Configuration Facade ---------------------------------------------------

// IsSet is a predicate wether a global configuration property is set.
func IsSet(key string) bool {
	return globalConf.IsSet(key)
}

// GetString returns a global configuration property as a string.
func GetString(key string) string {
	return globalConf.GetString(key)
}

// GetInt returns a global configuration property as an integer.
func GetInt(key string) int {
	return globalConf.GetInt(key)
}

// GetBool returns a global configuration property as a boolean value.
func GetBool(key string) bool {
	return globalConf.GetBool(key)
}

// IsInteractive is a predicate: are we running in interactive mode?
func IsInteractive() bool {
	return globalConf.IsInteractive()
}
