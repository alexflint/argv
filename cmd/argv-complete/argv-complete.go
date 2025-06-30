package main

import (
	"debug/elf"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"go.wit.com/apps/argv"
)

// Environment variables from bash:
//
// COMP_LINE
// The current command line.
//
// COMP_POINT
// The index of the current cursor position relative to the beginning of the current command. If the
// current cursor position is at the end of the current command, the value of this variable is equal
// to ${#COMP_LINE}.
//
// COMP_TYPE
// Set to an integer value corresponding to the type of completion attempted that caused a completion
// function to be called: TAB, for normal completion, ‘?’, for listing completions after successive
// tabs, ‘!’, for listing alternatives on partial word completion, ‘@’, to list completions if the
// word is not unmodified, or ‘%’, for menu completion.
//
// COMP_KEY
// The key (or final key of a key sequence) used to invoke the current completion function.
//
// See https://www.gnu.org/software/bash/manual/html_node/Bash-Variables.html

func Main() error {
	// this program is invoked to provide completions to another program
	// this program itself receives command line arguments and environment variables from bash
	// in order to make this program as understandable as possible, we do not use any command
	// line library in this program, even though we could, because it would confuse the issue
	// of whose arguments are being completed

	if len(os.Args) < 2 {
		return fmt.Errorf("expected first argument to be the program but got argv=%v", os.Args)
	}

	// this is the bash completions interface
	program := os.Args[1]                // program being invoked (first token on the command line)
	line := os.Getenv("COMP_LINE")       // the whole command line so far
	cursorstr := os.Getenv("COMP_POINT") // cursor position in terms of runes
	// we also have COMP_TYPE (integer) and COMP_KEY (integer)

	if line == "" {
		return fmt.Errorf("expected COMP_LINE environment variable but it was missing or empty")
	}
	if cursorstr == "" {
		return fmt.Errorf("expected COMP_POINT environment variable but it was missing or empty")
	}
	cursor, err := strconv.Atoi(cursorstr)
	if err != nil {
		return fmt.Errorf("expected COMP_POINT to be an integer but got %q", cursorstr)
	}
	if cursor < 0 {
		return fmt.Errorf("expected COMP_POINT to be positive but got %q", cursorstr)
	}

	// find the executable
	// note that LookPath detects explicit paths and returns it if it points to a file
	path, err := exec.LookPath(program)
	if err != nil {
		return fmt.Errorf("could not find %q in path: %w", program, err)
	}

	// load as an elf binary (TODO: support other platforms)
	f, err := elf.Open(path)
	if err != nil {
		return fmt.Errorf("error opening %q as ELF binary: %w", path, err)
	}
	defer f.Close()

	// find the section
	s := f.Section(".argv")
	if s == nil {
		return fmt.Errorf("argv section not found in %q", path)
	}

	raw, err := io.ReadAll(s.Open())
	if err != nil {
		return fmt.Errorf("error reading contents of argv section in binary: %w", err)
	}

	// parse the spec
	var spec argv.Command
	err = json.Unmarshal(raw, &spec)
	if err != nil {
		return fmt.Errorf("error unmarshaling completions spec: %w", err)
	}

	// cut the line at the cursor
	rs := []rune(line)
	if cursor > len(rs) {
		return fmt.Errorf("expected COMP_POINT to be an index in COMP_LINE but COMP_POINT=%v and len(COMP_LINE)=%v", cursor, len(rs))
	}
	line = string(rs[:cursor])

	// split the string by spaces (TODO: properly handle quotes per bash)
	parts := strings.Split(line, " ")
	lastpart := parts[len(parts)-1]

	// in the below, each line written to stdout is a completion (fmt.Println goes to stdout, log.Println goes to stderr)

	// complete long options
	if strings.HasPrefix(lastpart, "--") {
		for _, arg := range spec.Arguments {
			if strings.HasPrefix(arg.Long, lastpart) {
				fmt.Println(arg.Long)
			}
		}
		return nil
	}

	// complete short options
	if strings.HasPrefix(lastpart, "-") {
		for _, arg := range spec.Arguments {
			if strings.HasPrefix(arg.Short, lastpart) {
				fmt.Println(arg.Short)
			}
		}
		for _, arg := range spec.Arguments {
			if strings.HasPrefix(arg.Long, lastpart) {
				fmt.Println(arg.Long)
			}
		}
		return nil
	}

	return nil
}

func main() {
	log.SetFlags(0)

	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}
