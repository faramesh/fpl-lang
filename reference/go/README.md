# FPL Reference Parser (Go)

This directory contains an **experimental** reference implementation scaffold, inspired by language repos such as HashiCorp's HCL where syntax tooling and conformance are first-class.

## Scope

Current scope:

- Parse `agent` blocks with `default` and `rules` sections
- Parse rule clauses: `when`, `notify:`, `reason:`, `reeval:`
- Build AST nodes for boolean/comparison conditions (`and`, `or`, `not`, `matches`, `in`, comparison ops)
- Provide a tiny formatter command (`fplfmt`) for canonical output

This is not yet a full implementation of `grammar/fpl.ebnf`.

## Run

```bash
cd reference/go
go test ./...
go run ./cmd/fplparse ../../conformance/valid/basic-agent.fpl
go run ./cmd/fplfmt ../../conformance/valid/basic-agent.fpl
```

## Next milestones

- Full lexical token set (durations, currency, patterns)
- Rule and expression AST
- Parse diagnostics with stable error codes
- Conformance suite execution against valid/invalid fixtures
