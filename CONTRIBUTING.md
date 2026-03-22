# Contributing to FPL

Thanks for helping improve the Faramesh Policy Language.

## Repository goals

This repo is the language contract for FPL:

- grammar and specification
- examples and conformance fixtures
- editor integration assets
- reference parser/tooling scaffolds

Runtime governance behavior stays in `faramesh-core`.

## Development setup

```bash
make test
```

## Pull request checklist

- Keep grammar and spec changes in sync
- Add or update conformance fixtures in `conformance/`
- Include tests for reference tooling changes (`reference/go`)
- Update docs for user-visible behavior changes

## Change policy

If a grammar or semantic change is introduced, include:

- rationale
- migration notes
- compatibility impact
