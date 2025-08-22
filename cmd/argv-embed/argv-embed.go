package main

// Embed an argument spec in a binary. The binary is modified in-place.

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/alexflint/go-arg"
	"github.com/goccy/go-yaml"
	"go.wit.com/apps/argv"
)

// name of the section we create in the binary
const section = ".argv"

func Main() error {
	var args struct {
		Binary string `arg:"positional,required" help:"path to an executable"`
		Spec   string `arg:"positional,required" help:"path to file describing the command line arguments"`
	}
	arg.MustParse(&args)

	// read input
	buf, err := os.ReadFile(args.Spec)
	if err != nil {
		return fmt.Errorf("error opening spec: %w", err)
	}

	// parse the json
	var spec argv.Command
	err = yaml.Unmarshal(buf, &spec)
	if err != nil {
		return fmt.Errorf("error parsing spec: %w", err)
	}

	// TODO: validate

	// open a temporary file
	f, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// write the base64-encoded json
	//w := base64.NewEncoder(base64.StdEncoding, f)
	enc := json.NewEncoder(f)
	err = enc.Encode(spec)
	if err != nil {
		return fmt.Errorf("error encoding command to json: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("error closing temporary file: %w", err)
	}

	objcopy, err := exec.LookPath("objcopy")
	if err != nil {
		return fmt.Errorf("unable to find objcopy on the system, check that you have it installed: %w", err)
	}

	// run objcopy to embed the spec into the binary
	cmd := exec.Command(
		objcopy,
		"--add-section",
		section+"="+f.Name(),
		"--set-section-flags",
		section+"=noload,readonly",
		args.Binary,
		args.Binary)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("")
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
