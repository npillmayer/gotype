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

Utilities for (interactive) command line interfaces.
*/

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/mitchellh/colorstring"
)

// global settings, to be set by the active command
var welcomeMessage = "Welcome to Gallery [V0.1 experimental]"
var stdprompt = "[green]%s> "
var editmode string = "emacs"

// We support some interactive sub-commands (not part of the language grammar).
func displayCommands(out io.Writer) {
	io.WriteString(out, welcomeMessage)
	io.WriteString(out, "\n\nThe following commands are available:\n\n")
	io.WriteString(out, "help:               print this message\n")
	io.WriteString(out, "bye:                quit\n")
	io.WriteString(out, "mode [mode]:        display or set current editing mode\n")
	io.WriteString(out, "setprompt [prompt]: set current editing mode [to default],\n")
	io.WriteString(out, "                    supports color strings, e.g. '[blue]myprompt#'\n\n")
}

// Completer-tree for interactive frames sub-commands
var replCompleter = readline.NewPrefixCompleter(
	readline.PcItem("help"),
	readline.PcItem("bye"),
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("setprompt"),
)

// A base type to instantiate a REPL interpreter.
type BaseREPL struct {
	readline    *readline.Instance
	interpreter REPLCommandInterpreter
	toolname    string
	helper      func(io.Writer)
}

// All interpreter sub-commands implement this.
type REPLCommandInterpreter interface {
	InterpretCommand(string)
}

/* Create a readline instance.
 */
func NewReadline(toolname string) *readline.Instance {
	histfile := fmt.Sprintf("/tmp/%s-repl-history.tmp", toolname)
	prompt := fmt.Sprintf(stdprompt, toolname)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:              colorstring.Color(prompt),
		HistoryFile:         histfile,
		AutoComplete:        replCompleter,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterReplInput,
	})
	if err != nil {
		panic(err)
	}
	return rl
}

/* Enter a REPL and execute each command.
 * Commands are either tool-commands (setprompt, help, etc.)
 * or Gallery statements.
 */
func (repl *BaseREPL) doLoop() {
	for {
		line, err := repl.readline.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		words := strings.Fields(line)
		command := "<no command>"
		if len(words) > 0 {
			command = words[0]
		}
		if doExit := repl.executeCommand(command, words, line); doExit {
			break
		}
	}
}

/* Central dispatcher function to execute internal commands and PMMetaPost
 * statements. It receives the command (i.e. the first word of the line),
 * a list of words (args) including the command, and the complete line of text.
 * If it returns true, the REPL should terminate.
 */
func (repl *BaseREPL) executeCommand(cmd string, args []string, line string) bool {
	switch {
	case cmd == "help":
		displayCommands(repl.readline.Stderr())
		if repl.helper != nil {
			repl.helper(repl.readline.Stderr())
		}
	case cmd == "bye":
		println("> goodbye!")
		return true
	case cmd == "mode":
		if len(args) > 1 {
			switch args[1] {
			case "vi":
				repl.readline.SetVimMode(true)
				editmode = "vi"
				return false
			case "emacs":
				repl.readline.SetVimMode(false)
				editmode = "emacs"
				return false
			}
		}
		io.WriteString(repl.readline.Stderr(),
			fmt.Sprintf("> current input mode: %s\n", editmode))
	case cmd == "setprompt":
		var prmpt string
		if len(line) <= 10 {
			prmpt = fmt.Sprintf(stdprompt, repl.toolname)
		} else {
			prmpt = line[10:] + " "
		}
		repl.readline.SetPrompt(colorstring.Color(prmpt))
	case cmd == "":
	default:
		T.Debugf("call interpreter on: '%s'", line)
		repl.callInterpreter(line)
	}
	return false // do not exit
}

/* Call the Gallery interpreter, sending a statement.
 */
func (repl *BaseREPL) callInterpreter(line string) {
	defer func() {
		if r := recover(); r != nil {
			io.WriteString(repl.readline.Stderr(), "> error executing statement!\n")
			io.WriteString(repl.readline.Stderr(), fmt.Sprintf("> %v\n", r)) // TODO: get ERROR and print
		}
	}()
	repl.interpreter.InterpretCommand(line)
}

/* Input filter for REPL. Blocks ctrl-z.
 */
func filterReplInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
