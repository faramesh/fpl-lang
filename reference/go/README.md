# FPL Reference Parser (Go)

This directory contains the Go reference parser and formatter for FPL, aligned with the structured AST used by `faramesh-core`.

## Scope

Current scope:

- Parse the structured FPL document model: imports, runtime/provider/identity/trust blocks, agents, systems, and manifest lines
- Parse agent sub-blocks for budgets, phases, delegates, ambients, selectors, credentials, egress, model policy, session, spawn, completion gates, enforcement, and alerts
- Parse rule clauses including `when`, `notify:`, `reason:`, `host:`, `port:`, `method:`, `path:`, `query:`, `header:`, `headers:`, and legacy `reeval:`
- Provide a formatter command (`fplfmt`) for canonical output

This repository tracks the same structured grammar family as `faramesh-core` and is intended for parser/formatter conformance work.

## Run

```bash
cd reference/go
go test ./...
go run ./cmd/fplparse ../../conformance/valid/basic-agent.fpl
go run ./cmd/fplfmt ../../conformance/valid/basic-agent.fpl
```

## Next milestones

- Keep the reference parser and runtime parser in lockstep as the grammar evolves
- Expand conformance coverage for valid and invalid fixtures
- Tighten canonical formatting for the full document model
