#!/bin/bash
set -euo pipefail

# Generates a JSON matrix for GitHub Actions by grouping tests per resource.
# Resource + datasource tests for the same resource are merged into one group.
# Example output: {"include":[{"resource":"endpoint","tests":"TestA|TestB"},…]}

declare -A groups

for f in logzio/*_test.go; do
  tests=$(grep -o '^func Test[^(]*' "$f" | awk '{print $2}' || true)
  [ -z "$tests" ] && continue

  base=$(basename "$f" _test.go)
  resource=${base#resource_}
  resource=${resource#datasource_}

  if [ -n "${groups[$resource]:-}" ]; then
    groups[$resource]+=$'\n'"$tests"
  else
    groups[$resource]="$tests"
  fi
done

# Build JSON matrix
echo -n '{"include":['
first=true
for resource in $(printf '%s\n' "${!groups[@]}" | sort); do
  tests=$(echo "${groups[$resource]}" | tr '\n' '|' | sed 's/|$//')
  $first || echo -n ','
  first=false
  printf '{"resource":"%s","tests":"%s"}' "$resource" "$tests"
done
echo ']}'
