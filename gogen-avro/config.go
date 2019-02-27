package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	nsNone  = "none"
	nsShort = "short"
	nsFull  = "full"

	defaultPackageName     = "avro"
	defaultContainers      = false
	defaultShortUnions     = false
	defaultNamespacedNames = nsNone
)

type config struct {
	packageName     string
	containers      bool
	shortUnions     bool
	namespacedNames string
	targetDir       string
	files           []string
}

// parseCmdLine takes care of building the flagset and checking if the
// number of arguments matches the required ones
func parseCmdLine() config {
	cfg := config{}

	flag.StringVar(&cfg.packageName, "package", defaultPackageName, "Name of generated package.")
	flag.BoolVar(&cfg.containers, "containers", defaultContainers, "Whether to generate container writer methods.")
	flag.BoolVar(&cfg.shortUnions, "short-unions", defaultShortUnions, "Whether to use shorter names for Union types.")
	flag.StringVar(&cfg.namespacedNames, "namespaced-names", defaultNamespacedNames, "Whether to generate namespaced names for types. Default is \"none\"; \"short\" uses the last part of the namespace (last word after a separator); \"full\" uses all namespace string.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <target directory> <schema files>\n\nWhere 'flags' are:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
	}

	cfg.namespacedNames = strings.ToLower(cfg.namespacedNames)
	switch cfg.namespacedNames {
	case nsNone, nsShort, nsFull:
	default:
		fmt.Fprintf(os.Stderr, "namespaced-names: invalid value '%s'\n\n", cfg.namespacedNames)
		flag.Usage()
	}

	cfg.targetDir = flag.Arg(0)
	cfg.files = flag.Args()[1:]
	return cfg
}
