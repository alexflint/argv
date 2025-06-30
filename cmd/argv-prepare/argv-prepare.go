package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"go.wit.com/apps/argv"
)

// Prepare an argument spec for inclusion in a binary

func Main() error {
	var args struct {
		Input  string `arg:"positional,required"`
		Output string `arg:"positional,required"`
	}
	arg.MustParse(&args)

	// read input
	buf, err := os.ReadFile(args.Input)
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	// parse the json
	var cmd argv.Command
	err = json.Unmarshal(buf, &cmd)
	if err != nil {
		return fmt.Errorf("error unmarshaling command: %w", err)
	}

	// TODO: validate

	// marshal to json
	b, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("error marshaling command: %w", err)
	}

	// write it back out
	err = os.WriteFile(args.Output, b, 0777)
	if err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}
