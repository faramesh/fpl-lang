# FPL Conformance Suite

This directory contains language fixtures for parser/compiler conformance.

## Layout

- `valid/`: FPL programs that should parse successfully.
- `invalid/`: FPL programs that should fail parse/validation.

The initial fixture set is intentionally small and grows over time as the language evolves.

## Running conformance

From repository root:

```bash
bash scripts/run-conformance.sh
```

Expected behavior:

- Fixtures in `valid/` must parse successfully.
- Fixtures in `invalid/` must fail parsing.
