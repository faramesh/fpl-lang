# FPL — Faramesh Policy Language

<p align="center">
  <strong>A domain-specific language for AI agent governance.</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-green.svg" alt="MIT License" /></a>
</p>

---

FPL is a purpose-built language for governing AI agent tool calls. It replaces YAML policy files with concise, readable syntax that anyone — engineers, compliance officers, CISOs — can write and review.

FPL compiles to the same internal representation as YAML but provides **agent-native primitives** as first-class language constructs rather than YAML conventions.

## Quick example

```fpl
agent payment-bot {
  default deny
  model "gpt-4o"
  framework "langgraph"

  budget session {
    max $500
    daily $2000
    max_calls 100
    on_exceed deny
  }

  phase intake {
    permit read_customer
    permit get_order
  }

  rules {
    deny! shell/* reason: "never shell"

    defer stripe/refund
      when amount > 500
      notify: "finance"
      reason: "high value refund"

    permit stripe/*
      when amount <= 500
  }

  credential stripe {
    backend vault
    path secret/data/stripe/live
    ttl 15m
  }
}
```

The equivalent YAML is 60+ lines. This FPL is 25 lines with more readable structure.

## Key features

- **Agent-native primitives** — sessions, budgets, phases, delegation, ambient authority, and human approval workflows are language constructs, not conventions.
- **Mandatory deny (`deny!`)** — compile-time enforced. Cannot be overridden by position, child policies, priority, or any subsequent permit rule.
- **Natural language compilation** — `faramesh policy compile "deny all shell commands"` calls an LLM, produces FPL, validates it, and backtests it against real production history before activation.
- **Backtest before activation** — every policy generated from natural language is backtested against real decision-provenance records. You see exactly what the policy would have done to real past decisions before activating it.
- **Four input modes, one engine** — FPL directly, YAML (interchange), natural language (compiled to FPL), code annotations (`@faramesh.tool(defer_above=500)`).
- **GitOps native** — plain text files, version-controlled, validated in CI.

## Comparison

| Feature | FPL | OPA / Rego | Cedar | YAML + expr |
|---------|-----|-----------|-------|-------------|
| Sessions | First-class | No | No | Convention |
| Budget enforcement | First-class | No | No | Convention |
| Workflow phases | First-class | No | No | Convention |
| Delegation chains | First-class | No | No | Convention |
| Ambient authority | First-class | No | No | Convention |
| Human approval | First-class | No | No | Convention |
| Mandatory deny | Compiler-enforced | Runtime only | Runtime only | Documentation |
| NLP compilation | Built-in | No | No | No |
| Backtest | Built-in | Manual | No | Manual |
| Lines for typical policy | ~25 | ~80 | ~50 | ~65 |

## Install

FPL ships with the Faramesh CLI.

```bash
# Homebrew
brew install faramesh/tap/faramesh

# Shell script
curl -fsSL https://raw.githubusercontent.com/faramesh/faramesh-core/main/install.sh | bash

# Go
go install github.com/faramesh/faramesh-core/cmd/faramesh@latest
```

## Usage

```bash
# Validate an FPL file
faramesh policy validate policy.fpl

# Compile natural language to FPL
faramesh policy compile "deny all shell commands, defer refunds over $500 to finance"

# Parse and display structured output
faramesh policy fpl policy.fpl --json

# Run an agent under governance
faramesh run --policy policy.fpl -- python agent.py
```

## Documentation

| Document | Description |
|----------|-------------|
| [Language Reference](docs/LANGUAGE_REFERENCE.md) | Complete language reference — every keyword, block, and syntax construct |
| [Getting Started](docs/GETTING_STARTED.md) | Write your first policy in five minutes |
| [Comparison](docs/COMPARISON.md) | Detailed comparison with OPA, Cedar, and YAML |
| [Specification](spec/SPECIFICATION.md) | Formal language specification with EBNF grammar |
| [Examples](examples/) | Ready-to-use policy files for common agent types |

## File extension

`.fpl`

## Editor support

Planned: VS Code extension, JetBrains plugin, Neovim treesitter grammar. See [editors/](editors/) for details.

## License

[MIT](LICENSE)
