#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "$0")/.." && pwd)"
parser_cmd=(go run ./cmd/fplparse)

pass_count=0
fail_count=0

run_valid() {
  local file="$1"
  if (cd "$root_dir/reference/go" && "${parser_cmd[@]}" "$file" >/dev/null); then
    echo "PASS valid: $file"
    pass_count=$((pass_count + 1))
  else
    echo "FAIL valid: $file"
    fail_count=$((fail_count + 1))
  fi
}

run_invalid() {
  local file="$1"
  if (cd "$root_dir/reference/go" && "${parser_cmd[@]}" "$file" >/dev/null 2>&1); then
    echo "FAIL invalid (unexpected success): $file"
    fail_count=$((fail_count + 1))
  else
    echo "PASS invalid: $file"
    pass_count=$((pass_count + 1))
  fi
}

while IFS= read -r -d '' file; do
  run_valid "$file"
done < <(find "$root_dir/conformance/valid" -type f -name "*.fpl" -print0 | sort -z)

while IFS= read -r -d '' file; do
  run_invalid "$file"
done < <(find "$root_dir/conformance/invalid" -type f -name "*.fpl" -print0 | sort -z)

if [[ "$fail_count" -gt 0 ]]; then
  echo "conformance suite failed: $fail_count failures, $pass_count passes"
  exit 1
fi

echo "conformance suite passed: $pass_count checks"
