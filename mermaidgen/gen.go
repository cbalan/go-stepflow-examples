//go:build exclude

package main

import (
	"fmt"
	"os"

	"github.com/cbalan/go-stepflow-examples/mermaidgen"
	"github.com/cbalan/go-stepflow-examples/util/mermaid"
)

func main() {
	mermaidDiagram, err := mermaid.FromSteps(mermaidgen.DemoSteps())
	if err != nil {
		panic(err)
	}

	content := fmt.Sprintf("# Demo Steps\nFile generated via `go generate demo.go`\n```mermaid\n%s\n```\n", mermaidDiagram)

	if err := os.WriteFile("demo.md", []byte(content), 0644); err != nil {
		panic(err)
	}
}
