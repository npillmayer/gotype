/*
BSD License
Copyright (c) 2017, Norbert Pillmayer

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

----------------------------------------------------------------------

Tracing: TODO clean this up. This is only a first draft to get
things going.

*/

package tracing

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	InitTracingBootstrap()
}

const (
	LevelDebug = log.DebugLevel
	LevelInfo  = log.InfoLevel
	LevelError = log.ErrorLevel
)

type Trace struct {
	*log.Logger
	pkgname string
}

func (t Trace) P(key string, value interface{}) *log.Entry {
	return t.WithField(key, value) //.WithField("~", t.pkgname)
}

var Tracefile *os.File

var EquationsTracer Trace
var InterpreterTracer Trace
var SyntaxTracer Trace
var CommandTracer Trace
var GraphicsTracer Trace
var ScriptingTracer Trace

var packageTracers map[string]Trace = make(map[string]Trace)
var levelMap map[string]log.Level = make(map[string]log.Level)

func InitTracingBootstrap() {
	levelMap["DEBUG"] = log.DebugLevel
	levelMap["INFO"] = log.InfoLevel
	//levelMap["WARN"] = log.WarnLevel
	levelMap["ERROR"] = log.ErrorLevel

	equationsTracer := log.New()
	equationsTracer.SetLevel(log.InfoLevel)
	EquationsTracer = Trace{equationsTracer, "EQ"}
	packageTracers["EQ"] = EquationsTracer

	interpreterTracer := log.New()
	interpreterTracer.SetLevel(log.InfoLevel)
	InterpreterTracer = Trace{interpreterTracer, "IN"}
	packageTracers["IN"] = InterpreterTracer

	syntaxTracer := log.New()
	syntaxTracer.SetLevel(log.InfoLevel)
	SyntaxTracer = Trace{syntaxTracer, "SY"}
	packageTracers["SY"] = SyntaxTracer

	cmdTracer := log.New()
	cmdTracer.SetLevel(log.InfoLevel)
	CommandTracer = Trace{cmdTracer, "CD"}
	packageTracers["CD"] = CommandTracer

	gfxTracer := log.New()
	gfxTracer.SetLevel(log.InfoLevel)
	GraphicsTracer = Trace{gfxTracer, "GF"}
	packageTracers["GF"] = GraphicsTracer

	scriptingTracer := log.New()
	scriptingTracer.SetLevel(log.InfoLevel)
	ScriptingTracer = Trace{scriptingTracer, "LS"}
	packageTracers["LS"] = ScriptingTracer
}

func SetTraceLevel(p string) {
	t := packageTracers[p]
	t.SetLevel(levelMap[p])
}

func ConfigTracing(inputfilename string) *os.File {
	// TODO prepare for errorneous flag values
	InterpreterTracer.SetLevel(levelMap[viper.GetString("tracinginterpreter")])
	CommandTracer.SetLevel(levelMap[viper.GetString("tracingcommands")])
	EquationsTracer.SetLevel(levelMap[viper.GetString("tracingequations")])
	SyntaxTracer.SetLevel(levelMap[viper.GetString("tracingsyntax")])
	GraphicsTracer.SetLevel(levelMap[viper.GetString("tracinggraphics")])
	ScriptingTracer.SetLevel(levelMap[viper.GetString("tracingscripting")])

	if !viper.GetBool("tracingonline") {
		//tracefiledir := viper.GetString("outputdir")
		//file, err := os.OpenFile("_gotype.log", os.O_CREATE|os.O_WRONLY, 0666)
		file, err := os.Create("__gotype.log")
		if err != nil {
			CommandTracer.Error("cannot open tracefile, tracing to stderr")
		} else {
			Tracefile = file
			CommandTracer.Out = file
			CommandTracer.Formatter = &log.TextFormatter{}
			CommandTracer.Println("tracing to file")
			EquationsTracer.Out = file
			EquationsTracer.Formatter = &log.TextFormatter{}
			SyntaxTracer.Out = file
			SyntaxTracer.Formatter = &log.TextFormatter{}
			InterpreterTracer.Out = file
			InterpreterTracer.Formatter = &log.TextFormatter{}
			GraphicsTracer.Out = file
			GraphicsTracer.Formatter = &log.TextFormatter{}
			ScriptingTracer.Out = file
			ScriptingTracer.Formatter = &log.TextFormatter{}
		}
	}

	InterpreterTracer.P("level", InterpreterTracer.Level).Info("Interpreter-Trace is alive")
	CommandTracer.P("level", CommandTracer.Level).Info("Command-Trace is alive")
	EquationsTracer.P("level", EquationsTracer.Level).Info("Equations-Trace is alive")
	SyntaxTracer.P("level", SyntaxTracer.Level).Info("Syntax-Trace is alive")
	GraphicsTracer.P("level", GraphicsTracer.Level).Info("Graphics-Trace is alive")
	ScriptingTracer.P("level", ScriptingTracer.Level).Info("Scripting-Trace is alive")

	return Tracefile
}
