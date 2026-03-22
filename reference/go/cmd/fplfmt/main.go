package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/faramesh/fpl-lang/reference/go/fpl/formatter"
	"github.com/faramesh/fpl-lang/reference/go/fpl/parser"
)

func main() {
	write := flag.Bool("w", false, "write result to source file instead of stdout")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: fplfmt [-w] <file.fpl>")
		os.Exit(2)
	}

	path := flag.Arg(0)
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read file: %v\n", err)
		os.Exit(1)
	}

	doc, err := parser.New(string(src)).ParseDocument()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	formatted := formatter.FormatDocument(doc)
	if *write {
		if err := os.WriteFile(path, []byte(formatted), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write file: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Print(formatted)
}
