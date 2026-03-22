.PHONY: test conformance conformance-run ref-parser

test: conformance ref-parser

conformance:
	bash scripts/check-conformance.sh
	bash scripts/run-conformance.sh

conformance-run:
	bash scripts/run-conformance.sh

ref-parser:
	cd reference/go && go test ./...
