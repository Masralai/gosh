#!/bin/sh
set -eu

repo="Masralai/gosh"
bindir="${GOSH_BIN:-$HOME/.local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
	linux | darwin) ;;
	mingw* | msys* | cygwin*) os=windows ;;
	*)
		echo "gosh: unsupported OS: $os" >&2
		exit 1
		;;
esac

arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch=amd64 ;;
	aarch64 | arm64) arch=arm64 ;;
	*)
		echo "gosh: unsupported architecture: $arch" >&2
		exit 1
		;;
esac

ext=""
if [ "$os" = windows ]; then ext=".exe"; fi
asset="gosh-$os-$arch$ext"
base="https://github.com/$repo/releases/latest/download"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "downloading $asset"
curl -fsSL "$base/$asset" -o "$tmp/$asset"
curl -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt"

want=$(awk -v a="$asset" '$2 == a {print $1}' "$tmp/checksums.txt")
if [ -z "$want" ]; then
	echo "gosh: no checksum for $asset" >&2
	exit 1
fi
if command -v sha256sum >/dev/null 2>&1; then
	got=$(sha256sum "$tmp/$asset" | awk '{print $1}')
else
	got=$(shasum -a 256 "$tmp/$asset" | awk '{print $1}')
fi
if [ "$want" != "$got" ]; then
	echo "gosh: checksum mismatch for $asset" >&2
	exit 1
fi
echo "verified $asset"

mkdir -p "$bindir"
cp "$tmp/$asset" "$bindir/gosh$ext"
chmod 0755 "$bindir/gosh$ext"
echo "installed $bindir/gosh$ext"
case ":$PATH:" in
*":$bindir:"*) ;;
*) echo "add $bindir to PATH to run gosh" ;;
esac
