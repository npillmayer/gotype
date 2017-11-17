/* Package containing the main gotype command line tool, i.e. an
   interactive shell.
   This shell either starts a builtin command, an external
   command executable, or enters REPL mode.
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/docopt/docopt-go"
)

const toolname = "gtcmd"
const welcomeMessage = "Welcome to gtcmd [V0.1 experimental] interactive mode"

// Usage information. The heavy lifting is done by commands, so the
// top level usage information is thin.
var gtcmdUsage = `
Command line interface for the GOTYPE typesetting engine.

Usage:
  gtcmd (-h | --help)
  gtcmd --version
  gtcmd
  gtcmd <command> [ARG ...]

Examples:
  gtcmd                           # enter interactive shell
  gtcmd help typeset              # print help on command 'typeset'
  gtcmd typeset myinput.gt        # execute command 'typeset'

Options:
  -h, --help      Show this screen
  --version       Show version information

`

/* The main program. It parses the command line and then
   either starts a builtin command, an external
   command executable, or enters REPL mode.
*/
func main() {
	arguments, _ := docopt.Parse(gtcmdUsage, nil, true, "", false)
	fmt.Println(arguments)
	if arguments["<command>"] == nil {
		fmt.Println(welcomeMessage)
		repl() // start an interactive shell
	} else {
		commandLine := []string{
			arguments["<command>"].(string),
		}
		commandLine = append(commandLine, arguments["ARG"].([]string)...)
		executeCommand(false, commandLine) // start in batch-mode
	}
}

// We build a map of builtin commands to be executed from our shell

type builtinCommand func(io.Writer, ...string) int

var builtinCommands = map[string]builtinCommand{
	"help": displayUsageInformation,
}

func executeCommand(interactiveMode bool, args []string) (int, bool) {
	fmt.Println("execute Command (TODO):", args)
	fmt.Println("interactive mode is", interactiveMode)
	return 0, false
}

/* Built in command for displaying usage information.
 */
func displayUsageInformation(w io.Writer, args ...string) int {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, replCompleter.Tree("    "))
	return 0
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var replCompleter = readline.NewPrefixCompleter(
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("login"),
	readline.PcItem("say",
		readline.PcItemDynamic(listFiles("./"),
			readline.PcItem("with",
				readline.PcItem("following"),
				readline.PcItem("items"),
			),
		),
		readline.PcItem("hello"),
		readline.PcItem("bye"),
	),
	readline.PcItem("setprompt"),
	readline.PcItem("setpassword"),
	readline.PcItem("bye"),
	readline.PcItem("help"),
	readline.PcItem("go",
		readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
		readline.PcItem("install",
			readline.PcItem("-v"),
			readline.PcItem("-vv"),
			readline.PcItem("-vvv"),
		),
		readline.PcItem("test"),
	),
	readline.PcItem("sleep"),
)

func filterReplInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// REPL for gotype commands.  This is the interactive shell that
// invokes gotype commands, both internal and plugins.
//
func repl() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m»\033[0m ",
		HistoryFile:         "/tmp/gotype-repl-history.tmp",
		AutoComplete:        replCompleter,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterReplInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
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
		command := "no-op"
		if len(words) > 0 {
			command = words[0]
		}
		fmt.Println("Executing command", command)
		fmt.Println("   Arguments are:", words)
		if _, found := executeCommand(true, words); !found {
			fmt.Println("no builtin command found")
		}
		switch {
		case line == "help":
			displayUsageInformation(l.Stderr())
		case line == "bye":
			println("goodbye!")
			goto exit
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				l.SetVimMode(true)
			case "emacs":
				l.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case line == "mode":
			if l.IsVimMode() {
				println("current input mode: vi")
			} else {
				println("current input mode: emacs")
			}
		case line == "login":
			pswd, err := l.ReadPassword("please enter your password: ")
			if err != nil {
				break
			}
			println("you enter:", strconv.Quote(string(pswd)))
		case line == "setpassword":
			pswd, err := l.ReadPasswordWithConfig(setPasswordCfg)
			if err == nil {
				println("you set:", strconv.Quote(string(pswd)))
			}
		case strings.HasPrefix(line, "setprompt"):
			if len(line) <= 10 {
				log.Println("setprompt <prompt>")
				break
			}
			l.SetPrompt(line[10:])
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				log.Println("say what?")
				break
			}
			go func() {
				for range time.Tick(time.Second) {
					log.Println(line)
				}
			}()
		case line == "sleep":
			log.Println("sleep 4 second")
			time.Sleep(4 * time.Second)
		case line == "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}
