package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	"honnef.co/go/tools/staticcheck"

	"github.com/imirjar/metrx/cmd/staticlint/customanalyzer"
)

func main() {
	// Собираем стандартные анализаторы
	var analyzers = []*analysis.Analyzer{
		inspect.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		unusedresult.Analyzer,
	}

	// Добавляем анализаторы staticcheck (SA и другие классы)
	checks := map[string]bool{
		"SA": true,
		"ST": true, // Например, добавляем классы ST
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name[:2]] {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// Добавляем собственный анализатор
	analyzers = append(analyzers, customanalyzer.Analyzer)

	multichecker.Main(analyzers...)
}
