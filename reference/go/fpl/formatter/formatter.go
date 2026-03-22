package formatter

import (
	"strings"

	"github.com/faramesh/fpl-lang/reference/go/fpl/parser"
)

func FormatDocument(doc *parser.Document) string {
	var b strings.Builder
	for i, ag := range doc.Agents {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("agent ")
		b.WriteString(ag.Name)
		b.WriteString(" {\n")
		if ag.Default != "" {
			b.WriteString("  default ")
			b.WriteString(ag.Default)
			b.WriteString("\n")
		}
		if len(ag.Rules) > 0 {
			b.WriteString("\n")
			b.WriteString("  rules {\n")
			for _, r := range ag.Rules {
				b.WriteString("    ")
				b.WriteString(r.Effect)
				b.WriteString(" ")
				b.WriteString(r.Tool)
				if r.Condition != nil {
					b.WriteString("\n      when ")
					b.WriteString(formatExpr(r.Condition))
				}
				if r.Notify != "" {
					b.WriteString("\n      notify: \"")
					b.WriteString(r.Notify)
					b.WriteString("\"")
				}
				if r.Reason != "" {
					b.WriteString("\n      reason: \"")
					b.WriteString(r.Reason)
					b.WriteString("\"")
				}
				if r.Reeval != nil {
					if *r.Reeval {
						b.WriteString("\n      reeval: true")
					} else {
						b.WriteString("\n      reeval: false")
					}
				}
				b.WriteString("\n")
			}
			b.WriteString("  }\n")
		}
		b.WriteString("}\n")
	}
	return b.String()
}

func formatExpr(expr *parser.Expr) string {
	if expr == nil {
		return ""
	}
	switch expr.Kind {
	case "ident", "literal", "list":
		return expr.Value
	case "unary":
		return expr.Op + " " + formatExpr(expr.Right)
	case "binary":
		left := formatExpr(expr.Left)
		right := formatExpr(expr.Right)
		if expr.Left != nil && expr.Left.Kind == "binary" && precedence(expr.Left.Op) < precedence(expr.Op) {
			left = "(" + left + ")"
		}
		if expr.Right != nil && expr.Right.Kind == "binary" && precedence(expr.Right.Op) < precedence(expr.Op) {
			right = "(" + right + ")"
		}
		return left + " " + expr.Op + " " + right
	default:
		return expr.Value
	}
}

func precedence(op string) int {
	switch op {
	case "or":
		return 1
	case "and":
		return 2
	case ">", ">=", "<", "<=", "==", "!=", "matches", "in":
		return 3
	default:
		return 0
	}
}
