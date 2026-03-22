# FPL Editor Support

## Current status

FPL editor extensions are planned but not yet released.

## Planned extensions

| Editor | Format | Status |
|--------|--------|--------|
| VS Code | TextMate grammar + Language Server | Planned |
| JetBrains (IntelliJ, WebStorm, GoLand) | TextMate grammar | Planned |
| Neovim | Tree-sitter grammar | Planned |
| Vim | Syntax highlighting | Planned |
| Sublime Text | TextMate grammar | Planned |

## Interim highlighting

FPL syntax is close enough to HCL / Terraform that HCL syntax highlighting produces reasonable results. In VS Code, you can associate `.fpl` files with HCL in your settings:

```json
{
  "files.associations": {
    "*.fpl": "hcl"
  }
}
```

This gives you keyword highlighting and brace matching until the dedicated FPL extension is available.

## Contributing

If you'd like to contribute an editor extension, the [grammar/fpl.ebnf](../grammar/fpl.ebnf) file contains the formal grammar, and the [spec/SPECIFICATION.md](../spec/SPECIFICATION.md) defines the full language semantics.

The key syntax elements to highlight:

- **Keywords:** `agent`, `system`, `budget`, `phase`, `rules`, `delegate`, `ambient`, `selector`, `credential`, `default`, `model`, `framework`, `when`, `not`, `and`, `or`, `in`, `matches`
- **Effects:** `permit`, `deny`, `deny!`, `defer`, `allow`, `approve`, `block`, `reject`
- **Clauses:** `reason:`, `notify:`, `reeval:`
- **Literals:** strings (`"..."`), currencies (`$500`), durations (`24h`), data sizes (`10mb`)
- **Comments:** `# ...`
- **Tool patterns:** identifiers separated by `/` with optional `*` wildcard
