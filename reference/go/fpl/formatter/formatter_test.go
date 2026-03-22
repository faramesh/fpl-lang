package formatter

import (
	"strings"
	"testing"

	"github.com/faramesh/fpl-lang/reference/go/fpl/parser"
)

func TestFormatDocument(t *testing.T) {
	src := `agent payment-bot {
  default deny
  rules {
    permit read_customer when amount > 10 and not cmd matches "rm -rf" reason: "ok"
  }
}`

	doc, err := parser.New(src).ParseDocument()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out := FormatDocument(doc)

	if !strings.Contains(out, "agent payment-bot {") {
		t.Fatalf("missing agent header in output: %q", out)
	}
	if !strings.Contains(out, "when amount > 10 and not cmd matches \"rm -rf\"") {
		t.Fatalf("missing formatted condition in output: %q", out)
	}
	if !strings.Contains(out, "reason: \"ok\"") {
		t.Fatalf("missing reason in output: %q", out)
	}
}

func TestFormatIdempotent(t *testing.T) {
	src := `agent payment-bot {
  default deny
  rules {
    permit read_customer when amount <= 500 and not cmd matches "rm -rf" notify: "finance" reason: "review" reeval: true
  }
}`

	doc1, err := parser.New(src).ParseDocument()
	if err != nil {
		t.Fatalf("parse 1: %v", err)
	}
	out1 := FormatDocument(doc1)

	doc2, err := parser.New(out1).ParseDocument()
	if err != nil {
		t.Fatalf("parse formatted: %v\nformatted:\n%s", err, out1)
	}
	out2 := FormatDocument(doc2)

	if out1 != out2 {
		t.Fatalf("formatter not idempotent\nout1:\n%s\nout2:\n%s", out1, out2)
	}
}
