# FPL vs Other Policy Languages

This page compares policy languages for AI governance, agent guardrails,
and execution control in production tool-calling systems.

## Related guides

- [Getting Started](GETTING_STARTED.md)
- [Language Reference](LANGUAGE_REFERENCE.md)
- [Docs Index](README.md)

## FPL vs OPA / Rego

OPA (Open Policy Agent) uses Rego, a general-purpose logic language designed for infrastructure authorization (Kubernetes admission, API gateways, microservices).

| Dimension | OPA / Rego | FPL |
|-----------|-----------|-----|
| **Purpose** | General infrastructure authorization | AI agent governance |
| **Agent primitives** | None — must be modeled as data | First-class: sessions, budgets, phases, delegation, ambient authority |
| **Mandatory deny** | Runtime convention | Compile-time enforced (`deny!`) |
| **Learning curve** | Steep — Datalog-derived logic language | Low — reads like configuration |
| **Lines for typical policy** | 80+ | ~25 |
| **NLP compilation** | No | Built-in |
| **Backtest** | Manual replay scripts | Built-in with `faramesh policy backtest` |
| **Human approval** | External system required | `defer` is a language primitive |
| **Budget tracking** | Not a concept | `budget` block with per-session and daily limits |
| **Credential scoping** | Not a concept | `credential` block with backend integration |

**When to use OPA:** Infrastructure authorization (K8s admission, API gateways). OPA excels when you need a general-purpose policy engine across many systems.

**When to use FPL:** AI agent governance. FPL is purpose-built for the problem of controlling what AI agents can do, with constructs that map directly to agent workflows.

## FPL vs Cedar

Cedar is an authorization language by AWS, designed for application-level access control (who can do what on which resource).

| Dimension | Cedar | FPL |
|-----------|-------|-----|
| **Purpose** | Application access control | AI agent governance |
| **Model** | Principal → Action → Resource | Agent → Tool → Effect |
| **Agent primitives** | None | First-class |
| **Mandatory deny** | `forbid` with runtime enforcement | `deny!` with compile-time enforcement |
| **Workflow phases** | Not a concept | First-class with `phase` blocks |
| **Budget** | Not a concept | First-class |
| **NLP compilation** | No | Built-in |
| **Delegation** | Not a concept | First-class with scope and ceiling |

**When to use Cedar:** Application-level authorization where you need fine-grained access control on resources. Cedar is excellent for "can user X read document Y?" decisions.

**When to use FPL:** AI agent governance where the question is "should this agent be allowed to call this tool with these arguments, given its budget, session state, and delegation chain?"

## FPL vs YAML + expr-lang

Many governance systems use YAML policy files with an expression language (like expr-lang or CEL) for conditions.

| Dimension | YAML + expr | FPL |
|-----------|-------------|-----|
| **Readability** | Verbose, schema-dependent | Concise, self-descriptive |
| **Agent primitives** | Convention over configuration | Language constructs |
| **Mandatory deny** | Documentation convention | Compiler-enforced |
| **Lines for typical policy** | 60+ | ~25 |
| **Validation** | Schema-based (easy to get wrong) | Compiler-based (catches errors early) |
| **NLP compilation** | No | Built-in |
| **Error messages** | Generic YAML errors | Policy-specific errors |
| **Version control** | Works but noisy diffs | Clean diffs |

**When to use YAML:** When you already have a YAML-based toolchain and the overhead of another language is not justified. Faramesh supports YAML policies and always will.

**When to use FPL:** When readability, safety, and agent-native constructs matter. FPL is the recommended format for new Faramesh deployments.

## Side-by-side example

### FPL (25 lines)

```fpl
agent payment-bot {
  default deny

  budget session {
    max $500
    daily $2000
    max_calls 100
    on_exceed deny
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
}
```

### Equivalent YAML (65+ lines)

```yaml
faramesh-version: "1.0"
agent-id: "payment-bot"

tools:
  shell/exec:
    reversibility: irreversible
    blast_radius: high
  stripe/refund:
    reversibility: irreversible
    blast_radius: high
    tags: [financial]

budget:
  max_calls: 100
  session_usd: 500.00
  daily_usd: 2000.00
  on_exceed: deny

rules:
  - id: mandatory-deny-shell
    match:
      tool: "shell/*"
    effect: deny
    mandatory: true
    reason: "never shell"
    reason_code: MANDATORY_DENY

  - id: defer-high-refund
    match:
      tool: "stripe/refund"
      when: "args.amount > 500"
    effect: defer
    reason: "high value refund"
    notification_channel: "finance"

  - id: allow-small-refund
    match:
      tool: "stripe/refund"
      when: "args.amount <= 500"
    effect: permit
    reason: "routine refund"

default_effect: deny
```

### Equivalent Rego (80+ lines)

```rego
package faramesh.payment_bot

import future.keywords.in

default decision = {"effect": "deny", "reason": "no matching rule"}

decision = {"effect": "deny", "reason": "never shell", "mandatory": true} {
    glob.match("shell/*", ["/"], input.tool)
}

decision = {"effect": "defer", "reason": "high value refund", "notify": "finance"} {
    input.tool == "stripe/refund"
    input.args.amount > 500
}

decision = {"effect": "permit"} {
    input.tool == "stripe/refund"
    input.args.amount <= 500
}

# Budget enforcement requires external data and helper rules...
budget_ok {
    data.session.total_calls < 100
    data.session.total_usd < 500
    data.daily.total_usd < 2000
}

final_decision = decision {
    budget_ok
}

final_decision = {"effect": "deny", "reason": "budget exceeded"} {
    not budget_ok
}
```

## Summary

- Choose **FPL** when you need agent-native policy constructs (budgets, phases, delegation, mandatory deny) with readable policy-as-code.
- Choose **OPA/Rego** when you need broad infrastructure authorization outside agent runtime semantics.
- Choose **Cedar** when your primary model is principal/action/resource application authorization.

## Next steps

- [Getting Started](GETTING_STARTED.md)
- [Language Reference](LANGUAGE_REFERENCE.md)
- [Specification](../spec/SPECIFICATION.md)
- [Examples](../examples/)
