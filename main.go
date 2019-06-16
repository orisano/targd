package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func hasStage(stages []instructions.Stage, name string) (int, bool) {
	if i, ok := instructions.HasStage(stages, name); ok {
		return i, true
	}
	i, err := strconv.Atoi(name)
	if err != nil {
		return -1, false
	}
	if i < 0 || len(stages) <= i {
		return -1, false
	}
	return i, true
}

func main() {
	target := flag.String("target", "", "target stage name (required)")
	p := flag.String("f", "Dockerfile", "Dockerfile path")
	out := flag.String("o", "-", "generated Dockerfile path")
	flag.Parse()

	if *target == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	var r io.Reader
	if *p == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(*p)
		if err != nil {
			log.Fatal("failed to open Dockerfile:", err)
		}
		defer f.Close()
		r = f
	}

	var w io.Writer
	if *out == "-" {
		w = os.Stdout
	} else {
		f, err := os.Create(*out)
		if err != nil {
			log.Fatal("failed to create Dockerfile:", err)
		}
		defer f.Close()
		w = f
	}

	result, err := parser.Parse(r)
	if err != nil {
		log.Fatal("failed to parse Dockerfile:", err)
	}

	stages, _, err := instructions.Parse(result.AST)
	if err != nil {
		log.Fatal("failed to parse instructions:", err)
	}

	byName := map[string]int{}
	for i := range stages {
		if stages[i].Name == "" {
			stages[i].Name = strconv.Itoa(i)
		}
		byName[strings.ToLower(stages[i].Name)] = i
	}

	required := make(map[string]struct{})

	var visitStage func(string)
	visitStage = func(name string) {
		if _, ok := required[strings.ToLower(name)]; ok {
			return
		}
		i, ok := byName[strings.ToLower(name)]
		if !ok {
			log.Fatalf("unknown stage: %v", name)
		}
		required[strings.ToLower(name)] = struct{}{}
		stage := stages[i]
		if _, ok := byName[strings.ToLower(stage.BaseName)]; ok {
			visitStage(stage.BaseName)
		}
		for _, cmd := range stage.Commands {
			copyCmd, ok := cmd.(*instructions.CopyCommand)
			if !ok {
				continue
			}
			if copyCmd.From == "" {
				continue
			}
			if _, ok := byName[strings.ToLower(copyCmd.From)]; ok {
				visitStage(copyCmd.From)
			}
		}
	}
	visitStage(*target)

	for _, stage := range stages {
		name := stage.Name
		if _, ok := required[strings.ToLower(name)]; !ok {
			continue
		}

		fmt.Fprintln(w, stage.SourceCode)
		for _, cmd := range stage.Commands {
			fmt.Fprintln(w, cmd)
		}
		fmt.Fprintln(w)
	}
}
