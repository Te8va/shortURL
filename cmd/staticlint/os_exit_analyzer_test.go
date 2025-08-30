package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func testOsExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osExitAnalyzer, "./...")
}
