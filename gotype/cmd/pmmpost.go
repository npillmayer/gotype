package cmd

/*
----------------------------------------------------------------------

BSD License
Copyright (c) 2017, Norbert Pillmayer <norbert@pillmayer.com>

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

Command line interface for the Poor Man's MetaPost language and
graphical system.

This module creates a possibly interactive CLI, configures the input (file),
trace output, a PMMPost interpreter instance, and necessary backend graphics
modules. Then sets up either an interactive shell or a file input instance.

Usage: pmmpost [-d <output-dir>] [-f png|svg] [-x] [-m] inputfile.pmp

	--outdir    -d   output directory, default = .
	--format    -f   output format,    default = png
	--debug     -x   debug mode,       default = off
	--vi        -m   vi editing mode,  default = off

*/

import (
	"fmt"
	"io"
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/backend/gfx"
	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/pmmpost"
	"github.com/spf13/cobra"
)

// global settings
const pmmpWelcomeMessage = "Welcome to Poor Man's MetaPost [V0.1 experimental]"

// pmmpostCmd represents the pmmp command
var pmmpostCmd = &cobra.Command{
	Use:   "pmmpost",
	Short: "A poor man's MetaPost implementation",
	Long: `pmmpost is a drawing language und engine reminiscend of John Hobby's
	MetaPost sytem. Users supply an input program, either as a text file
	or on the command prompt. Output may be generated as PDF, SVG or PNG.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runPMMPostCmd,
}

// COBRA init method. Defines command and flags.
func init() {
	rootCmd.AddCommand(pmmpostCmd)
	pmmpostCmd.Flags().BoolP("vi", "m", false, "Set vi editing mode")
	pmmpostCmd.Flags().StringP("outdir", "d", ".", "Output directory")
	pmmpostCmd.Flags().StringP("format", "f", "png", "Output format")
	pmmpostCmd.Flags().BoolP("debug", "x", false, "Debug mode")
}

// A type to instantiate a REPL interpreter.
type PMMPostREPL struct {
	BaseREPL
	pmmpintp *pmmpost.PMMPostInterpreter
}

// Bridge to the PMMPost interpreter.
func (pmmprepl *PMMPostREPL) InterpretCommand(input string) {
	errors := pmmprepl.pmmpintp.ParseStatements([]byte(input))
	if errors != nil {
		for _, err := range errors {
			io.WriteString(pmmprepl.readline.Stderr(), err.Error())
		}
	}
}

// The PMMPost command (pmmp).
func runPMMPostCmd(cmd *cobra.Command, args []string) {
	fmt.Println(pmmpWelcomeMessage)
	welcomeMessage = pmmpWelcomeMessage
	var inputfilename string
	if len(args) > 0 {
		T.Infof("input file is %s", args[0])
		inputfilename = args[0]
	}
	defer tracing.Tracefile.Close()
	startPMMPostInput(inputfilename)
}

// Start PMMPost input. if a filename is given, opens the file and reads from
// there. Otherwise starts an interactive shell (REPL).
func startPMMPostInput(inputfilename string) {
	if inputfilename == "" {
		config.IsInteractive = true
		repl := NewPMMPostREPL() // go into interactive mode
		repl.toolname = "pmmpost"
		log.SetOutput(repl.readline.Stderr())
		defer repl.readline.Close()
		repl.doLoop()
	} else {
		//input, err := antlr.NewFileStream(inputfilename) // TODO refactor to get rid of ANTLR
		_, err := antlr.NewFileStream(inputfilename) // TODO refactor to get rid of ANTLR
		if err != nil {
			T.Errorf("cannot open input file")
		} else {
			config.IsInteractive = false
			defer func() {
				if r := recover(); r != nil {
					T.Errorf("error executing PMMPost statement!")
				}
			}()
			//interpreter := pmmpost.NewPMMPostInterpreter()
			//interpreter.SetOutputRoutine(png.NewPNGOutputRoutine()) // will produce PNG format
			//interpreter.ParseStatements(input)
		}
	}
}

// Set up a new REPL entity. It contains a readline-instance (for putting
// out a prompt and read in a line) and a PMMetaPost parser. The REPL will
// then forward PMMetaPost statements to the parser.
func NewPMMPostREPL() *PMMPostREPL {
	rl := NewReadline("pmmpost")
	repl := &PMMPostREPL{}
	repl.readline = rl
	repl.interpreter = repl // we are our own bridge to the interpreter
	repl.pmmpintp = pmmpost.NewPMMPostInterpreter(true, func(pic *gfx.Picture) {
		//
		io.WriteString(repl.readline.Stderr(), fmt.Sprintf("SHIPPING PICTURE '%s'\n", pic.Name))
	})
	//repl.pmmpintp.SetOutputRoutine(png.NewPNGOutputRoutine()) // will produce PNG format
	repl.toolname = "pmmpost"
	return repl
}
