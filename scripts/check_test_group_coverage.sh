#!/bin/bash
set -euo pipefail

echo "🔍 Verifying all test files are included in test groups..."

# 1. List all *_test.go files (including subfolders, relative to root)
all_test_files=$(find logzio -type f -name '*_test.go' | sort)

# 2. Gather all grouped files and prepend 'logzio/' to each
grouped_files=$(cat .github/test-groups/group_*.txt | grep -v '^\s*$' | sed 's/^[[:space:]]*//' | sed 's|^|logzio/|' | sort)

# 3. Diff the two
missing_files=$(comm -23 <(echo "$all_test_files") <(echo "$grouped_files"))

if [[ -n "$missing_files" ]]; then
  echo "❌ ERROR: The following test files are not included in any test group:"
  echo "$missing_files"
  exit 1
else
  echo "✅ All test files are assigned to a test group."
fi
