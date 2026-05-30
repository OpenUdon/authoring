#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$root"
go test ./...
go vet ./...
GOWORK=off go test ./...
GOWORK=off go vet ./...

forbidden='github.com/OpenUdon/(openudon|ramen|uws|apitools)'
if go list -deps ./... | grep -E "$forbidden"; then
  echo "Authoring import boundary includes a forbidden downstream dependency" >&2
  exit 1
fi

workspace="$(cd "$root/.." && pwd)"
for repo in openudon ramen; do
  if [[ -d "$workspace/$repo" ]]; then
    (cd "$workspace/$repo" && go test ./...)
  fi
done
