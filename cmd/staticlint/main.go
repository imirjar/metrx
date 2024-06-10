package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// реализация будет ниже
	return nil, nil
}

func main() {
	multichecker.Main(
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
	)
}
