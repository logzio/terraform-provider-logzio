#!/bin/bash
set -euo pipefail

echo "🔍 Verifying that all test functions are included in test groups..."

# Step 1: Get all test function names from *_test.go files
grep -h -o '^func Test[^(]*' logzio/*_test.go | awk '{print $2}' | sort -u > /tmp/defined_tests.txt

# Step 2: Safely read all group files line by line, even if missing final newline
tmp_group_file="/tmp/grouped_tests.txt"
> "$tmp_group_file"

for f in .github/test-groups/group_*.txt; do
  # Append every line, even the last one if no trailing newline
  while IFS= read -r line || [[ -n "$line" ]]; do
    # Strip leading/trailing whitespace and ignore empty lines
    trimmed=$(echo "$line" | xargs)
    if [[ -n "$trimmed" ]]; then
      echo "$trimmed" >> "$tmp_group_file"
    fi
  done < "$f"
done

# Sort and deduplicate
sort -u "$tmp_group_file" > /tmp/grouped_tests_sorted.txt

# Step 3: Compare
missing=$(comm -23 /tmp/defined_tests.txt /tmp/grouped_tests_sorted.txt || true)
extra=$(comm -13 /tmp/defined_tests.txt /tmp/grouped_tests_sorted.txt || true)

if [[ -n "$missing" ]]; then
  echo "❌ ERROR: The following test functions are defined but NOT included in any group:"
  echo "$missing"
  echo
fi

if [[ -n "$extra" ]]; then
  echo "⚠️ WARNING: The following test functions are in a group but NOT defined in code:"
  echo "$extra"
  echo
fi

if [[ -n "$missing" ]]; then
  exit 1
else
  echo "✅ All defined test functions are included in the test groups."
fi
