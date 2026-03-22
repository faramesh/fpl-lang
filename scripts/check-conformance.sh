#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "$0")/.." && pwd)"

valid_count=$(find "$root_dir/conformance/valid" -type f -name "*.fpl" | wc -l | tr -d ' ')
invalid_count=$(find "$root_dir/conformance/invalid" -type f -name "*.fpl" | wc -l | tr -d ' ')

if [[ "$valid_count" -lt 1 ]]; then
  echo "expected at least one valid conformance fixture"
  exit 1
fi

if [[ "$invalid_count" -lt 1 ]]; then
  echo "expected at least one invalid conformance fixture"
  exit 1
fi

echo "conformance fixture check passed: $valid_count valid, $invalid_count invalid"
