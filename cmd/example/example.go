package main

import (
	"github.com/alexflint/go-arg"
	"github.com/kr/pretty"
)

// an example program for testing

func main() {
	var args struct {
		Path    string   `arg:"positional,required"`
		IDs     []string `arg:"--ids"`
		Verbose bool     `arg:"-v"`
	}
	arg.MustParse(&args)
	pretty.Println(args)
}
