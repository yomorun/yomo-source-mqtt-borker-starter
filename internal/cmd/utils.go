package cmd

import (
	"log"
	"plugin"
	"strings"

	"github.com/yomorun/yomo-source-mqtt-borker-starter/internal/source"
)

type baseOptions struct {
	Filename string
}

func buildSourceFile(opts *baseOptions, args []string) (string, error) {
	if len(args) >= 1 && strings.HasSuffix(args[0], ".go") {
		// the second arg of `yomo-mqtt build xxx.go` is a .go file
		opts.Filename = args[0]
	}

	// build the file first
	log.Print("Building the Source Function File...")
	soFile, err := source.Build(opts.Filename, true)
	if err != nil {
		log.Print("❌ Build the Source file failure with err: ", err)
		return "", err
	}
	return soFile, nil
}

func buildAndLoadHandler(opts *baseOptions, args []string) (plugin.Symbol, error) {
	// build the file first
	soFile, err := buildSourceFile(opts, args)
	if err != nil {
		return nil, err
	}

	// load handle
	slHandler, err := source.LoadHandler(soFile)
	if err != nil {
		log.Print("❌ Load handle from .so file failure with err: ", err)
		return nil, err
	}
	return slHandler, nil
}
