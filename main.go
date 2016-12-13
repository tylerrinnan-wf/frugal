package main

import (
	"fmt"
	"os"
	"plugin"
	"sort"
	"strings"

	"github.com/urfave/cli"

	"github.com/Workiva/frugal/compiler"
	"github.com/Workiva/frugal/compiler/generator"
	"github.com/Workiva/frugal/compiler/globals"
	"github.com/Workiva/frugal/compiler/parser"
)

const defaultTopicDelim = "."

var (
	help               bool
	gen                string
	out                string
	delim              string
	audit              string
	recurse            bool
	verbose            bool
	version            bool
	useVendor          bool
)

func main() {
	app := cli.NewApp()
	app.Name = "frugal"
	app.Usage = "a tool for code generation"
	app.Version = globals.Version
	app.HideVersion = true
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "help, h",
			Usage:       "show help",
			Destination: &help,
		},
		cli.StringFlag{
			Name:        "gen",
			Usage:       genUsage(),
			Destination: &gen,
		},
		cli.StringFlag{
			Name:        "out",
			Usage:       "set the output location for generated files (no gen-* folder will be created)",
			Destination: &out,
		},
		cli.StringFlag{
			Name:        "delim",
			Value:       defaultTopicDelim,
			Usage:       "set the delimiter for pub/sub topic tokens",
			Destination: &delim,
		},
		cli.BoolFlag{
			Name:        "recurse, r",
			Usage:       "generate included files",
			Destination: &recurse,
		},
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "verbose mode",
			Destination: &verbose,
		},
		cli.BoolFlag{
			Name:        "version",
			Usage:       "print the version",
			Destination: &version,
		}, cli.StringFlag{
			Name:        "audit",
			Usage:       "frugal file to run audit against",
			Destination: &audit,
		}, cli.BoolFlag{
			Name:        "use-vendor",
			Usage:       "use specified import references for vendored includes and do not generate code for them (supported by go)",
			Destination: &useVendor,
		},
	}

	app.Action = func(c *cli.Context) error {
		if help {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}

		if version {
			cli.ShowVersion(c)
			os.Exit(0)
		}

		if len(c.Args()) == 0 {
			fmt.Printf("Usage: %s [options] file\n\n", app.Name)
			fmt.Printf("Use %s -help for a list of options\n", app.Name)
			os.Exit(1)
		}

		if gen == "" && audit == "" {
			fmt.Println("No output language specified")
			fmt.Printf("Usage: %s [options] file\n\n", app.Name)
			fmt.Printf("Use %s -help for a list of options\n", app.Name)
			os.Exit(1)
		}

		options := compiler.Options{
			Gen:                gen,
			Out:                out,
			Delim:              delim,
			Recurse:            recurse,
			Verbose:            verbose,
			UseVendor:          useVendor,
		}

		// TODO: This is currently a workaround to https://github.com/golang/go/issues/17928.
		// Remove once that has been released.
		_, opts, err := compiler.CleanGenParam(gen)
		if err != nil {
			fmt.Printf("Failed to generate %s:\n\t%s\n", options.File, err.Error())
			os.Exit(1)
		}
		if err := preloadPlugins(opts); err != nil {
			fmt.Printf("Failed to generate %s:\n\t%s\n", options.File, err.Error())
			os.Exit(1)
		}

		// Handle panics for graceful error messages.
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Failed to generate %s:\n\t%s\n", options.File, r)
				os.Exit(1)
			}
		}()

		auditor := parser.NewAuditor()
		for _, options.File = range c.Args() {
			if audit == "" {
				err = compiler.Compile(options)
			} else {
				err = auditor.Audit(audit, options.File)
			}
			if err != nil {
				fmt.Printf("Failed to generate %s:\n\t%s\n", options.File, err.Error())
				os.Exit(1)
			}
		}

		return nil
	}

	app.Run(os.Args)
}

func genUsage() string {
	usage := "generate code with a registered generator and optional parameters " +
		"(lang[:key1=val1[,key2[,key3=val3]]])\n"
	langKeys := make([]string, 0, len(generator.Languages))
	for lang := range generator.Languages {
		langKeys = append(langKeys, lang)
	}
	sort.Strings(langKeys)
	langPrefix := ""
	for _, lang := range langKeys {
		options := generator.Languages[lang]
		optionsStr := ""
		optionKeys := make([]string, 0, len(options))
		for name := range options {
			optionKeys = append(optionKeys, name)
		}
		sort.Strings(optionKeys)
		for _, name := range optionKeys {
			optionsStr += fmt.Sprintf("\n\t        %s:\t%s", name, options[name])
		}
		usage += fmt.Sprintf("%s\t    %s:%s", langPrefix, lang, optionsStr)
		langPrefix = "\n"
	}
	return usage
}

// preloadPlugins opens the Plugins specified by the "plugins" option, if any.
// This is currently a workaround to https://github.com/golang/go/issues/17928.
// TODO: remove once workaround is no longer needed.
func preloadPlugins(options map[string]string) error {
	namesStr, ok := options["plugins"]
	if !ok {
		return nil
	}
	names := strings.Split(namesStr, ":")
	for _, name := range names {
		if _, err := plugin.Open(name); err != nil {
			return err
		}
	}
	return nil
}
