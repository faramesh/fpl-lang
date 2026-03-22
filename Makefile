.PHONY: test conformance ref-parser

test: conformance ref-parser

conformance:
	bash scripts/check-conformance.sh

ref-parser:
	cd reference/go && go test ./...
