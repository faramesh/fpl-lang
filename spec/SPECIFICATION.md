# FPL Language Specification

Version 1.0

## 1. Overview

FPL (Faramesh Policy Language) is a domain-specific language for expressing AI agent governance policy. An FPL document defines how one or more agents are permitted, denied, or deferred when they attempt tool calls.

## 2. Lexical structure

### 2.1 Source encoding

FPL source files are UTF-8 encoded. The file extension is `.fpl`.

### 2.2 Line structure

FPL is line-oriented. Each statement occupies one or more lines. Continuation is implicit inside braces.

### 2.3 Comments

Comments begin with `#` and extend to the end of the line.

```fpl
# This is a comment
permit http/get  # This is also a comment
```

### 2.4 Whitespace

Whitespace (spaces, tabs) separates tokens but has no semantic meaning outside of string literals. Newlines are significant only as statement terminators.

### 2.5 Keywords

The following identifiers are reserved:

```
agent      system     budget     phase      rules
delegate   ambient    selector   credential manifest
default    model      framework  version    var
permit     deny       deny!      defer      allow
approve    block      reject     when       not
and        or         in         matches    reason
notify     reeval     scope      ttl        ceiling
inherited  source     cache      on_exceed  on_unavailable
on_timeout backend    path       max        daily
max_calls  session    global     true       false
```

### 2.6 Identifiers

An identifier starts with a letter and may contain letters, digits, underscores, and hyphens.

```
Ident = Letter { Letter | Digit | "_" | "-" }
```

### 2.7 Literals

**String literals** are enclosed in double quotes:
```
"hello world"
```

**Integer literals** are sequences of digits:
```
500
```

**Currency literals** are prefixed with `$`:
```
$500
$2000.50
```

**Duration literals** are integers followed by a unit suffix:
```
15m    # 15 minutes
24h    # 24 hours
7d     # 7 days
```

| Unit | Meaning |
|------|---------|
| `s`  | seconds |
| `m`  | minutes |
| `h`  | hours   |
| `d`  | days    |
| `w`  | weeks   |

**Data size literals** are integers followed by a size unit:
```
10mb
1gb
```

| Unit | Meaning |
|------|---------|
| `kb` | kilobytes |
| `mb` | megabytes |
| `gb` | gigabytes |

## 3. Grammar

See [grammar/fpl.ebnf](../grammar/fpl.ebnf) for the complete EBNF grammar.

### 3.1 Document structure

An FPL document is a sequence of top-level declarations:

```
Document = { AgentBlock | SystemBlock | ManifestStmt }
```

A document must contain at least one `agent` block.

### 3.2 Agent block

The `agent` block is the primary unit of governance. It declares everything about how one agent is governed.

```
AgentBlock = "agent" Ident "{" { AgentBody } "}"
```

An agent block may contain, in any order: `default`, `model`, `framework`, `version`, `var`, `budget`, `phase`, `rules`, `delegate`, `ambient`, `selector`, `credential`.

### 3.3 Rules block

```
RulesBlock = "rules" "{" { Rule } "}"
```

Rules are evaluated top-to-bottom. The first matching rule determines the outcome. If no rule matches, the `default` effect applies.

### 3.4 Rule syntax

```
Rule = Effect ToolPattern [ WhenClause ] [ NotifyClause ] [ ReasonClause ]
```

### 3.5 Tool patterns

Tool patterns use `/` as a namespace separator and `*` as a wildcard:

- `stripe/refund` — exact match
- `stripe/*` — matches all tools in the `stripe` namespace
- `*` — matches all tools

### 3.6 Conditions

Conditions follow the `when` keyword. They are boolean expressions:

```
Condition = OrExpr
OrExpr    = AndExpr { "or" AndExpr }
AndExpr   = NotExpr { "and" NotExpr }
NotExpr   = [ "not" ] Comparison
```

Comparison operators: `>`, `>=`, `<`, `<=`, `==`, `!=`, `matches`, `in`.

Built-in variables available in conditions:

| Variable | Type | Description |
|----------|------|-------------|
| `amount` | number | `args.amount` shorthand |
| `cmd` | string | `args.cmd` shorthand |
| `args.*` | any | Tool call arguments |
| `session.sum(...)` | number | Aggregate function over session history |
| `selectors.<name>.<field>` | any | External selector data |

## 4. Effects

| Effect | Meaning | Override behavior |
|--------|---------|-------------------|
| `permit` | Allow the tool call | Can be overridden by `deny` or `deny!` |
| `deny` | Block the tool call | Can be overridden by later `permit` |
| `deny!` | Mandatory deny | **Cannot be overridden** (compile-time enforced) |
| `defer` | Pause for human approval | Can be overridden by `deny!` |

Aliases: `allow` = `permit`, `approve` = `permit`, `block` = `deny`, `reject` = `deny`.

## 5. Semantics

### 5.1 Evaluation order

1. Kill switch check (runtime, not in policy)
2. Phase visibility filter
3. Rules evaluation (top-to-bottom, first match wins)
4. Default effect (if no rule matched)

### 5.2 The `deny!` invariant

A rule with `deny!` creates a compile-time constraint:

- No subsequent `permit` rule in the same `rules` block can match the same tool pattern
- No child policy in an `extends` chain can override it
- No priority mechanism can override it
- The compiler MUST reject policies that violate this invariant

This is a structural guarantee, not a runtime convention.

### 5.3 Scope inheritance

When agent A delegates to agent B:
- B's effective permissions are the intersection of A's permissions and B's declared scope
- The `ceiling` property on the delegation determines the maximum scope
- `ceiling inherited` means B cannot exceed A's scope

### 5.4 Budget enforcement

Budget checks occur before rule evaluation. If the session budget is exhausted, the tool call is denied regardless of what the rules say.

### 5.5 Phase scoping

Tools declared in a `phase` block are only visible when that phase is active. Tool calls to non-visible tools are denied without rule evaluation.

## 6. Compilation

FPL compiles to an internal intermediate representation (IR) that is identical to the IR produced by YAML policy files. This guarantees behavioral equivalence regardless of input format.

The compilation pipeline:

```
FPL source → Lexer → Parser → AST → IR → Policy Engine
YAML source → YAML parser → IR → Policy Engine
NLP text → LLM → FPL source → (same as above)
```

## 7. Conformance

A conforming FPL implementation MUST:

1. Parse all valid FPL documents as defined by the grammar
2. Reject documents that violate the `deny!` invariant at compile time
3. Evaluate rules in document order (first match wins)
4. Apply the `default` effect when no rule matches
5. Enforce budget limits before rule evaluation
6. Restrict phase-scoped tools to their declared phase

## 8. Changes

| Version | Date | Change |
|---------|------|--------|
| 1.0 | 2025-03 | Initial specification |
