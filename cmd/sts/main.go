package main

import (
	"flag"
	"log"
	"os"

	"github.com/ekhabarov/sts"
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
	helperpkg = flag.String("hp", "", "Package with helper functions")
	debug     = flag.Bool("debug", false, "Debug")

	version = "0.0.1-alpha-dev"
)

func main() {
	flag.Parse()

	name, content, err := sts.Run(
		*src, *dst,
		*st, *dt,
		*out, *helperpkg, version, *debug)
	must(err)

	file, err := os.Create(name)
	must(err)

	_, err = file.Write(content)
	must(err)
}

func must(err error) {
	if err != nil {
		if *debug {
			log.Fatalf("%+v", err)
		} else {
			log.Fatalf("%v", err)
		}
	}
}
