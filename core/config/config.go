/*
Package config is the central package for application configuration.

Configuration

All configuration is started explicitely with a call to Initialize().
There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

Tracing

During configuration all global tracers are set up. Currently tracing
with the Go standard logger is supported, as well as logrus. It is
possible for clients to use a different tracer, e.g. one relying on
https://github.com/golang/glog. To use a different logging implementation,
clients will have to implement an adapter to tracing.Trace (please
refer to the documentation for package tracing as well as to implementations
of adapters for Go log and for logrus).

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
package config

import (
	"os"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
)

// Are we running in interactive mode?
var IsInteractive bool = true

var globalConf Configuration

// Configuration is an interface to be implemented by every configuration
// adapter.
type Configuration interface {
	Init()
	IsSet(key string) bool
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

// Initialize is the top level function for setting up the
// application configuration.
// It will call Init() on the Configuration passed as an argument, and
// make the Configuration available globally.
// Functions in this package will serve as a facade to the Configuration.
func Initialize(conf Configuration) {
	globalConf = conf
	globalConf.Init()
	InitTracing(getAdapterFromConfiguration())
}

// InitTracing sets up all the global module tracers, reading trace levels
// and tracing destinations from the application configuration.
//
// InitTracing is usually called by Initialize().
func InitTracing(adapter tracing.Adapter) {
	tracing.CreateTracers(adapter)
	ConfigureTracing("")
}

var knownTraceAdapters = map[string]tracing.Adapter{
	"go": gologadapter.GetAdapter(),
	//"logrus": logrusadapter.GetAdapter(), // now to be set by AddTraceAdapter()
}

// AddTraceAdapter is an extension point for clients who want to use
// their own tracing adapter implementation.
// key will be used at configuration initialization time to identify
// this adapter, e.g. in configuration files.
//
// Clients will have to call this before any call to Initialize(), otherwise
// the adapter cannot be found.
func AddTraceAdapter(key string, adapter tracing.Adapter) {
	knownTraceAdapters[key] = adapter
}

// Get the concrete tracing implementation adapter from the appcation
// configuration.
// Default is an adapter for the Go standard log package.
func getAdapterFromConfiguration() tracing.Adapter {
	adapterPackage := GetString("tracing")
	adapter := knownTraceAdapters[adapterPackage]
	if adapter == nil {
		adapter = gologadapter.GetAdapter()
	}
	return adapter
}

// ConfigureTracing is usually called by InitTracing().
//
// It is exported as it may be useful in testing scenarios.
func ConfigureTracing(inputfilename string) {
	DefaultTracing() // set default trace levels from configuration
	if GetBool("tracingonline") {
		if inputfilename != "" {
			file, err := os.Create("__gotype.log")
			if err != nil {
				tracing.CommandTracer.Errorf("cannot open tracefile, tracing to stderr")
			}
			tracing.Tracefile = file
		}
		if tracing.Tracefile != nil {
			tracing.CommandTracer.SetOutput(tracing.Tracefile)
			tracing.EquationsTracer.SetOutput(tracing.Tracefile)
			tracing.SyntaxTracer.SetOutput(tracing.Tracefile)
			tracing.InterpreterTracer.SetOutput(tracing.Tracefile)
			tracing.GraphicsTracer.SetOutput(tracing.Tracefile)
			tracing.ScriptingTracer.SetOutput(tracing.Tracefile)
			tracing.CoreTracer.SetOutput(tracing.Tracefile)
			tracing.EngineTracer.SetOutput(tracing.Tracefile)
		}
	}

	tracing.InterpreterTracer.P("level", tracing.InterpreterTracer.GetTraceLevel()).Infof("Interpreter-Trace is alive")
	tracing.CommandTracer.P("level", tracing.CommandTracer.GetTraceLevel()).Infof("Command-Trace is alive")
	tracing.EquationsTracer.P("level", tracing.EquationsTracer.GetTraceLevel()).Infof("Equations-Trace is alive")
	tracing.SyntaxTracer.P("level", tracing.SyntaxTracer.GetTraceLevel()).Infof("Syntax-Trace is alive")
	tracing.GraphicsTracer.P("level", tracing.GraphicsTracer.GetTraceLevel()).Infof("Graphics-Trace is alive")
	tracing.ScriptingTracer.P("level", tracing.ScriptingTracer.GetTraceLevel()).Infof("Scripting-Trace is alive")
	tracing.CoreTracer.P("level", tracing.CoreTracer.GetTraceLevel()).Infof("Core-Trace is alive")
	tracing.EngineTracer.P("level", tracing.CoreTracer.GetTraceLevel()).Infof("Engine-Trace is alive")
}

// DefaultTracing sets all global tracers to their default trace levels,
// read from the application configuration.
func DefaultTracing() {
	tracing.InterpreterTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracinginterpreter")))
	tracing.CommandTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingcommands")))
	tracing.EquationsTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingequations")))
	tracing.SyntaxTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingsyntax")))
	tracing.GraphicsTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracinggraphics")))
	tracing.ScriptingTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingscripting")))
	tracing.CoreTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingcore")))
	tracing.EngineTracer.SetTraceLevel(tracing.TraceLevelFromString(GetString("tracingengine")))
}

// IsSet is a predicate wether a global configuration property is set.
func IsSet(key string) bool {
	return globalConf.IsSet(key)
}

// GetString returns a global configuration property as a string.
func GetString(key string) string {
	return globalConf.GetString(key)
}

// GetString returns a global configuration property as an integer.
func GetInt(key string) int {
	return globalConf.GetInt(key)
}

// GetString returns a global configuration property as a boolean value.
func GetBool(key string) bool {
	return globalConf.GetBool(key)
}
