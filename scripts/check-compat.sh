#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$root"
go test ./...
go vet ./...
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off go test -race ./...

forbidden='github.com/OpenUdon/(openudon|ramen|uws|apitools)'
if go list -deps ./... | grep -E "$forbidden"; then
  echo "Authoring import boundary includes a forbidden downstream dependency" >&2
  exit 1
fi

workspace="$(cd "$root/.." && pwd)"
check_root="$(mktemp -d "${TMPDIR:-/tmp}/authoring-compat.XXXXXX")"
trap 'rm -rf -- "$check_root"' EXIT

echo "Authoring revision: $(git rev-parse --verify HEAD)"
git status --short

# GOFLAGS also binds subprocess Go builds launched by consumer tests. Only
# temporary modfiles change; consumer manifests and the operator workspace do not.
for repo in openudon kinet; do
  consumer="$workspace/$repo"
  [[ -f "$consumer/go.mod" && -f "$consumer/go.sum" ]] || {
    echo "Required unchanged consumer checkout is missing: $consumer" >&2
    exit 1
  }
  mkdir "$check_root/$repo"
  consumer_mod="$check_root/$repo/consumer.mod"
  cp "$consumer/go.mod" "$consumer_mod"
  cp "$consumer/go.sum" "$check_root/$repo/consumer.sum"
  GOWORK=off go mod edit -modfile="$consumer_mod" \
    -replace="github.com/OpenUdon/authoring=$root"
  (
    cd "$consumer"
    export GOWORK=off GOFLAGS="-modfile=$consumer_mod -mod=readonly"
    echo "$repo revision: $(git rev-parse --verify HEAD)"
    go list -m -json github.com/OpenUdon/authoring
    go test ./...
    go vet ./...
  )
done

# Ramen deliberately retains the removed APIs at its old pin. It is frozen,
# and is qualified separately without the new-source override or parent go.work.
[[ -f "$workspace/ramen/go.mod" ]] || {
  echo "Required frozen Ramen checkout is missing" >&2
  exit 1
}
(
  cd "$workspace/ramen"
  export GOWORK=off GOFLAGS=-mod=readonly
  echo "Frozen Ramen revision: $(git rev-parse --verify HEAD)"
  go list -m -json github.com/OpenUdon/authoring
  go build -o "$check_root/ramen" ./cmd/ramen
)
