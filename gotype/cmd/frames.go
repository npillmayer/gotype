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

Command line interface for the definition for typesetting frames.

This module creates a possibly interactive CLI, configures the input (file),
trace output, an interpreter instance, and necessary backend modules.
Then sets up either an interactive shell or a file input instance.

Usage: frames [-d <output-dir>] [-x] [-m] inputfile.gy

	--outdir    -d   output directory, default = .
	--debug     -x   debug mode,       default = off
	--vi        -m   vi editing mode,  default = off

*/

import (
	"fmt"
	"log"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/syntax/gallery"
	"github.com/spf13/cobra"
)

const fWelcomeMessage = "Welcome to Gallery frames [V0.1 experimental]"
const ftoolname = "frames"

// framesCmd represents the frames command
var framesCmd = &cobra.Command{
	Use:   ftoolname,
	Short: "Define frames for Gallery",
	Long: `Gallery is a language for the definition of frames for 
	typesetting. The syntax is reminiscent of Donald Knuth's MetaFont system.
	Users supply an input program, either as a text file or on the command
	prompt. Output is generated as a JSON file.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runFramesCmd,
}

/* COBRA init method. Defines command and flags.
 */
func init() {
	rootCmd.AddCommand(framesCmd)
	framesCmd.Flags().BoolP("vi", "m", false, "Set vi editing mode")
	framesCmd.Flags().StringP("outdir", "d", ".", "Output directory")
	framesCmd.Flags().BoolP("debug", "x", false, "Debug mode")
}

// A type to instantiate a REPL interpreter.
type FramesREPL struct {
	BaseREPL
	gyintp *gallery.GalleryInterpreter
}

/* Bridge to the Gallery interpreter.
 */
func (frepl *FramesREPL) InterpretCommand(input string) {
	frepl.gyintp.ParseStatements(antlr.NewInputStream(input))
}

/* The FramesCmd command (frames).
 */
func runFramesCmd(cmd *cobra.Command, args []string) {
	fmt.Println(fWelcomeMessage)
	welcomeMessage = fWelcomeMessage
	var inputfilename string
	if len(args) > 0 {
		T.Infof("input file is %s", args[0])
		inputfilename = args[0]
	}
	//tracing.Tracefile = tracing.ConfigTracing(inputfilename)
	defer tracing.Tracefile.Close()
	startFramesInput(inputfilename)
}

/* Start Gallery input. If a filename is given, opens the file and reads from
 * there. Otherwise starts an interactive shell (REPL).
 */
func startFramesInput(inputfilename string) {
	if inputfilename == "" {
		config.IsInteractive = true
		repl := NewFramesREPL() // go into interactive mode
		log.SetOutput(repl.readline.Stderr())
		defer repl.readline.Close()
		repl.doLoop()
	} else {
		input, err := antlr.NewFileStream(inputfilename) // TODO refactor to get rid of ANTLR
		if err != nil {
			T.Errorf("cannot open input file")
		} else {
			config.IsInteractive = false
			defer func() {
				if r := recover(); r != nil {
					T.Errorf("error executing Gallery statement!")
				}
			}()
			interpreter := gallery.NewGalleryInterpreter(true)
			interpreter.ParseStatements(input)
		}
	}
}

/* Set up a new REPL entity. It contains a readline-instance (for putting
 * out a prompt and read in a line) and a PMMetaPost parser. The REPL will
 * then forward PMMetaPost statements to the parser.
 */
func NewFramesREPL() *FramesREPL {
	rl := NewReadline(ftoolname)
	repl := &FramesREPL{}
	repl.readline = rl
	repl.interpreter = repl // we are our own bridge to the interpreter
	repl.gyintp = gallery.NewGalleryInterpreter(true)
	repl.toolname = ftoolname
	return repl
}
