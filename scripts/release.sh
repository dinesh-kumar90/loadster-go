#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
OUT_DIR="dist"
BINARY="loadster"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 v1.0.0"
  exit 1
fi

mkdir -p "$OUT_DIR"
rm -f "$OUT_DIR"/*

platforms=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
  "windows amd64"
  "windows arm64"
)

echo "Building release binaries for $VERSION"

for platform in "${platforms[@]}"; do
  read -r goos goarch <<< "$platform"

  ext=""
  if [[ "$goos" == "windows" ]]; then
    ext=".exe"
  fi

  output="$OUT_DIR/${BINARY}-${VERSION}-${goos}-${goarch}${ext}"

  echo "- $goos/$goarch -> $output"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$output" .
done

(
  cd "$OUT_DIR"
  shasum -a 256 * > checksums.txt
)

echo "Done. Artifacts created in $OUT_DIR/"
ls -1 "$OUT_DIR"
