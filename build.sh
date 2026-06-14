#!/bin/bash

export PATH="/usr/local/go/bin:$PATH"

# md-viewer build & package script
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# --- App version (About panel / CFBundle*) ---
# User-facing version; override when releasing: MARKETING_VERSION=1.0.0 ./build.sh
#MARKETING_VERSION="${MARKETING_VERSION:-0.2.2}"
MARKETING_VERSION="${MARKETING_VERSION:-$(cat VERSION)}"

APP_NAME="md-viewer"
BUNDLE_DIR="md-viewer.app"
CONTENTS_DIR="$BUNDLE_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
FRAMEWORKS_DIR="$CONTENTS_DIR/Frameworks"

echo "🚀 Starting build process..."
echo "📌 CFBundleShortVersionString will be: $MARKETING_VERSION (from build.sh)"

# 0. Tidy Go modules
echo "🧹 Tidying Go modules..."
go mod tidy

# 1. Build Swift MarkdownEngine
echo "📦 Building Swift MarkdownEngine..."
##rm -rf .build/
swift build -c release

# 2. Build Go Application
echo "🐹 Building Go application..."
go build -o $APP_NAME -buildvcs=false

# --- Bump build number only after compile succeeded ---
BUILD_NUMBER_FILE="$SCRIPT_DIR/.build_number"
if [[ -f "$BUILD_NUMBER_FILE" ]]; then
	BUILD_NUMBER="$(tr -d '[:space:]' < "$BUILD_NUMBER_FILE" | head -c 20)"
	[[ -z "$BUILD_NUMBER" ]] && BUILD_NUMBER=0
else
	BUILD_NUMBER=0
fi
if ! [[ "$BUILD_NUMBER" =~ ^[0-9]+$ ]]; then
	echo "⚠️  Invalid .build_number (must be non-negative integer). Resetting to 0."
	BUILD_NUMBER=0
fi
BUILD_NUMBER=$((BUILD_NUMBER + 1))
echo "$BUILD_NUMBER" > "$BUILD_NUMBER_FILE"
echo "📌 CFBundleVersion=$BUILD_NUMBER (auto-increment, saved in .build_number)"

# 3. Prepare .app structure
echo "📂 Preparing .app bundle..."
mkdir -p "$MACOS_DIR"
mkdir -p "$FRAMEWORKS_DIR"
echo "📄 Generating Info.plist..."
while IFS= read -r line || [[ -n "$line" ]]; do
	line="${line//__MARKETING_VERSION__/${MARKETING_VERSION}}"
	line="${line//__BUILD_NUMBER__/${BUILD_NUMBER}}"
	printf '%s\n' "$line"
done < "$SCRIPT_DIR/Info.plist.template" > "$CONTENTS_DIR/Info.plist"

# 4. Copy and fix dynamic library
echo "🔗 Linking Swift library..."
cp .build/release/libMarkdownEngine.dylib "$FRAMEWORKS_DIR/"
# Also keep a copy or symlink in root for direct execution
ln -sf .build/release/libMarkdownEngine.dylib libMarkdownEngine.dylib

# Fix rpath in the executable
echo "🛠 Fixing rpath..."
# Remove old rpaths if they exist to prevent duplication errors
install_name_tool -delete_rpath "@executable_path/../Frameworks/" "$APP_NAME" 2>/dev/null || true
install_name_tool -delete_rpath "./" "$APP_NAME" 2>/dev/null || true
install_name_tool -delete_rpath ".build/release/" "$APP_NAME" 2>/dev/null || true

# Add rpaths
install_name_tool -add_rpath "@executable_path/../Frameworks/" "$APP_NAME"
install_name_tool -add_rpath "./" "$APP_NAME"
install_name_tool -add_rpath ".build/release/" "$APP_NAME"

# 5. Move executable to bundle
echo "🚚 Moving executable to bundle..."
cp $APP_NAME "$MACOS_DIR/"
ls -ltah $APP_NAME "$MACOS_DIR/"

echo "✅ Build complete!"
echo "👉 You can now run the app via: ./$APP_NAME"
echo "👉 Or use the bundle: open $BUNDLE_DIR"
