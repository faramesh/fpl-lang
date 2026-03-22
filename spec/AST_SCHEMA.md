# FPL Reference AST Schema (Draft)

Version: 0.1

This document defines the JSON AST shape emitted by the reference parser in reference/go.

## Document

```json
{
  "agents": [Agent]
}
```

## Agent

```json
{
  "name": "string",
  "default": "string (optional)",
  "rules": [Rule] (optional)
}
```

## Rule

```json
{
  "effect": "permit|deny|deny!|defer",
  "tool": "string",
  "condition": Expr (optional),
  "notify": "string" (optional),
  "reason": "string" (optional),
  "reeval": "boolean" (optional)
}
```

## Expr

```json
{
  "kind": "ident|literal|list|unary|binary",
  "op": "string (optional)",
  "value": "string (optional)",
  "left": Expr (optional),
  "right": Expr (optional)
}
```

Rules:
- binary expressions use left and right
- unary expressions use right
- ident/literal/list use value

## Example

```json
{
  "agents": [
    {
      "name": "payment-bot",
      "default": "deny",
      "rules": [
        {
          "effect": "permit",
          "tool": "read_customer",
          "condition": {
            "kind": "binary",
            "op": "and",
            "left": {
              "kind": "binary",
              "op": "<=",
              "left": {"kind": "ident", "value": "amount"},
              "right": {"kind": "literal", "value": "500"}
            },
            "right": {
              "kind": "unary",
              "op": "not",
              "right": {
                "kind": "binary",
                "op": "matches",
                "left": {"kind": "ident", "value": "cmd"},
                "right": {"kind": "literal", "value": "\"rm -rf\""}
              }
            }
          },
          "notify": "finance",
          "reason": "human review",
          "reeval": true
        }
      ]
    }
  ]
}
```

Compatibility note:
- This schema is for the reference parser only and may evolve.
- The production parser/compiler behavior remains in faramesh-core.
