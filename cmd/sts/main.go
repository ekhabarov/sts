package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/powerflyco/sts"
)

var (
	src = flag.String("src", "", "Source: <path/to/*.go file>:<structure>")
	dst = flag.String("dst", "", "Dest: <path/to/*.go file>:<structure>")
	out = flag.String("o", "./transform", `Path to output directory.
Last part of this path will be used as output package name.
`)
	st = flag.String("st", "sts", "Field tag in source structure.")
	dt = flag.String("dt", "",
		"List of comma-separated tag on destination structure.",
	)
	helperpkg = flag.String("hp", "", "Package with helper functions.")
	debug     = flag.Bool("debug", false, "Add debug info into output file.")
	ver       = flag.Bool("version", false, "Print current version.")
	cfgmap    = flag.String("map", "", "Path to YAML file with field map config.")

	version = "0.0.6"
)

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("version: %s\n", version)
		os.Exit(0)
	}

	var m *sts.FieldConfig

	if cfgmap != nil && *cfgmap != "" {
		var err error
		m, err = sts.LoadFieldConfigMap(*cfgmap)
		must(err)
	}

	name, content, err := sts.Run(
		*src, *dst,
		*st, *dt,
		*out, *helperpkg, version, *debug,
		m,
	)
	must(err)

	file, err := os.Create(name)
	must(err)

	_, err = file.Write(content)
	must(err)
}

func must(err error) {
	if err != nil {
		if *debug {
			fmt.Printf("%+v", err)
		} else {
			fmt.Printf("%v", err)
		}
		os.Exit(1)
	}
}
