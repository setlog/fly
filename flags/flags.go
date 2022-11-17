package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const FlagMajor = "major"
const FlagMinor = "minor"
const FlagAuto = "auto"

type Flags struct {
	ScriptPrefix           string
	VersionIncrementMethod VersionIncrementMethod
}

type VersionIncrementMethod int

const IncrementMajor VersionIncrementMethod = 0
const IncrementMinor VersionIncrementMethod = 1
const IncrementAuto VersionIncrementMethod = 2
const IncrementAsk VersionIncrementMethod = 3

func Parse(args []string) *Flags {
	f := Flags{}
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.Usage = func() {}
	var major, minor, auto bool
	fs.BoolVar(&major, FlagMajor, false, "make a major version increment")
	fs.BoolVar(&minor, FlagMinor, false, "make a minor version increment")
	fs.BoolVar(&auto, FlagAuto, false, "automatically decide how to increment version")
	fs.Parse(args)
	if auto && (major || minor) {
		fmt.Println("Explicit --auto disallows usage of --major and --minor.")
		os.Exit(1)
	}
	if major && minor {
		fmt.Println("Cannot use --major and --minor at the same time.")
		os.Exit(1)
	}
	f.VersionIncrementMethod = readVersionIncrementMethod(major, minor, auto)
	if fs.NArg() > 1 {
		fmt.Println("Must provide only script file name suffix as last argument.")
		os.Exit(1)
	}
	if fs.NArg() == 1 {
		f.ScriptPrefix = fs.Arg(0)
		if strings.HasSuffix(f.ScriptPrefix, ".sql") {
			f.ScriptPrefix = f.ScriptPrefix[:len(f.ScriptPrefix)-4]
		}
	} else {
		f.ScriptPrefix = "migration"
	}
	return &f
}

func readVersionIncrementMethod(major, minor, auto bool) VersionIncrementMethod {
	if major {
		return IncrementMajor
	}
	if minor {
		return IncrementMinor
	}
	if auto {
		return IncrementAuto
	}
	return IncrementAsk
}
