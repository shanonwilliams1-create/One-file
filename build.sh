#!/usr/bin/env bash
# Cross-platform build script for OmniConfig
# Builds for desktop (Linux, macOS, Windows) and mobile (iOS, Android)
set -euo pipefail

PROJECT="github.com/omnicofig/cli"
BINARY="omniconfig"
VERSION="${VERSION:-0.1.0-dev}"
BUILD_DIR="${BUILD_DIR:-build}"
LDFLAGS="-X 'github.com/omnicofig/cli/cmd.version=${VERSION}' -s -w"

# Desktop platform matrix: GOOS/GOARCH
DESKTOP_PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Mobile platform matrix
MOBILE_PLATFORMS=(
    "ios/arm64"
    "android/arm64"
    "android/amd64"
)

ALL_PLATFORMS=("${DESKTOP_PLATFORMS[@]}" "${MOBILE_PLATFORMS[@]}")

mkdir -p "${BUILD_DIR}"

echo "==> Building OmniConfig v${VERSION}"
echo ""

for PLATFORM in "${ALL_PLATFORMS[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "${PLATFORM}"
    OUTPUT="${BUILD_DIR}/${BINARY}-${GOOS}-${GOARCH}"

    # Add .exe for Windows
    if [ "${GOOS}" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi

    echo "  Building for ${GOOS}/${GOARCH}..."
    export GOOS="${GOOS}"
    export GOARCH="${GOARCH}"
    export CGO_ENABLED=0

    if [ "${GOOS}" = "ios" ]; then
        # iOS needs special handling - builds as a library for embedding
        echo "    (iOS: building static library for app embedding)"
        CGO_ENABLED=0 GOOS=ios GOARCH=arm64 \
            go build -tags ios -ldflags="${LDFLAGS}" -o "${OUTPUT}" "${PROJECT}" || \
            echo "    ⚠ iOS build requires Xcode toolchain — skipping"
    elif [ "${GOOS}" = "android" ]; then
        # Android builds as a shared library via gomobile or standalone
        echo "    (Android: building native binary for Termux/embedded)"
        CGO_ENABLED=0 GOOS=android GOARCH="${GOARCH}" \
            go build -ldflags="${LDFLAGS}" -o "${OUTPUT}" "${PROJECT}" || \
            echo "    ⚠ Android build may need NDK — see docs for gomobile setup"
    else
        CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
            go build -ldflags="${LDFLAGS}" -o "${OUTPUT}" "${PROJECT}"
    fi

    if [ -f "${OUTPUT}" ]; then
        echo "    ✓ ${OUTPUT}"
    fi
done

echo ""
echo "==> Build complete!"
if command -v sha256sum &>/dev/null; then
    echo "==> Generating checksums..."
    (cd "${BUILD_DIR}" && sha256sum "${BINARY}"-* 2>/dev/null >> checksums.txt || true)
fi
ls -lh "${BUILD_DIR}/" 2>/dev/null || true