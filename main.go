package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

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
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	}

	result, err := parser.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	stages, _, err := instructions.Parse(result.AST)
	if err != nil {
		log.Fatal(err)
	}

	required := make(map[string]struct{})

	var visitStage func(string)
	visitStage = func(name string) {
		if _, ok := required[name]; ok {
			return
		}
		i, ok := hasStage(stages, name)
		if !ok {
			log.Fatalf("unknown stage: %v", name)
		}

		required[name] = struct{}{}
		stage := stages[i]
		if _, ok := instructions.HasStage(stages, stage.BaseName); ok {
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
			visitStage(copyCmd.From)
		}
	}
	visitStage(*target)

	for i, stage := range stages {
		name := stage.Name
		if name == "" {
			name = strconv.Itoa(i)
		}
		if _, ok := required[name]; !ok {
			continue
		}

		fmt.Println(stage.SourceCode)
		for _, cmd := range stage.Commands {
			fmt.Println(cmd)
		}
		fmt.Println()
	}
}
