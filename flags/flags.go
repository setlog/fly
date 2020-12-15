package flags

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	ScriptPrefix string
}

func Parse(args []string) *Flags {
	f := Flags{}
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.Usage = func() {}
	fs.Parse(args)
	if fs.NArg() > 1 {
		fmt.Println("Provide at most 1 argument: script file name suffix.")
		os.Exit(1)
	}
	if fs.NArg() == 1 {
		f.ScriptPrefix = fs.Arg(0)
	} else {
		f.ScriptPrefix = "migration"
	}
	return &f
}
