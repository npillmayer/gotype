/*
Package config is the central package for application configuration.

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

Tracing

During configuration all global tracers are set up. Currently tracing
with the Go standard logger is supported, as well as logrus. It is
possible for clients to use a different tracer, e.g. one relying on
https://github.com/golang/glog. To use a different logging implementation,
clients will have to implement an adapter to tracing.Trace (please
refer to the documentation for package tracing as well as to implementations
of adapters for Go log and for logrus).
*/
package config

import (
	"fmt"
	"os"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
	"github.com/npillmayer/gotype/core/config/tracing/logrusadapter"
	"github.com/spf13/viper"
)

// Are we running in interactive mode?
var IsInteractive bool = true

// Initialize is the top level function for setting up the
// application configuration.
//
// It will call InitDefaults(), InitConfigPath() and InitTracing().
func Initialize() {
	InitDefaults()
	InitConfigPath()
	InitTracing(getAdapterFromConfiguration())
}

// InitConfigPath is usually called by Initialize().
func InitConfigPath() {
	viper.SetConfigName("gotype")        // name of config file (without extension)
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.AddConfigPath("$GOPATH/etc/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.gotype") // call multiple times to add many search paths
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

// IsSet is a predicate wether a configuration flag is set to true.
func IsSet(key string) bool {
	return viper.IsSet(key)
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
	"go":     gologadapter.GetAdapter(),
	"logrus": logrusadapter.GetAdapter(),
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
	adapterPackage := viper.GetString("tracing")
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
	if !viper.GetBool("tracingonline") {
		if inputfilename != "" {
			//tracefiledir := viper.GetString("outputdir")
			//file, err := os.OpenFile("_gotype.log", os.O_CREATE|os.O_WRONLY, 0666)
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
	tracing.InterpreterTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracinginterpreter")))
	tracing.CommandTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracingcommands")))
	tracing.EquationsTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracingequations")))
	tracing.SyntaxTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracingsyntax")))
	tracing.GraphicsTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracinggraphics")))
	tracing.ScriptingTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracingscripting")))
	tracing.CoreTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracingcore")))
	tracing.EngineTracer.SetTraceLevel(tracing.TraceLevelFromString(viper.GetString("tracingengine")))
}
