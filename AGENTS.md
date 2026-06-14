# AGENTS.md — md-viewer-webview

## Build

```bash
./build.sh
```

`build.sh` runs: `go mod tidy` → `swift build -c release` → `go build` → writes `md-viewer.app/Contents/Info.plist` from `Info.plist.template` (with version placeholders filled in) → creates the `.app` bundle.

Note: `build.sh` prepends `/usr/local/go/bin` to `PATH`; if your Go is elsewhere, adjust `GOPATH` or ensure `go` is on PATH.

### App version (About md-viewer)

| What | Where |
|------|--------|
| **CFBundleShortVersionString** (e.g. `1.2.0`) | Set in **`build.sh`** as `MARKETING_VERSION` (default at top of script). One-off override: `MARKETING_VERSION=1.0.0 ./build.sh` |
| **CFBundleVersion** (build number) | **Auto-increments** after a **successful** compile (`swift` + `go build`); persisted in **`.build_number`** in the repo root (commit it if the team should share one sequence). |
| **Copyright, document types, bundle id** | Edit **`Info.plist.template`** (still uses `__MARKETING_VERSION__` / `__BUILD_NUMBER__` placeholders for the two keys above). |

Direct `./md-viewer` (binary not inside `.app`) may not load a plist; the About panel then shows a `0.0.0-dev` fallback. Prefer **`open md-viewer.app`** after `./build.sh` for correct version display.

After Swift changes: re-run `./build.sh` (rebuilds dylib).
After Go changes: `go build -o md-viewer` (faster iteration, no need to re-run build.sh).

## Architecture

```
Go main.go → CGO → libMarkdownEngine.dylib (Swift) → apple/swift-markdown
                  ↓
             webview_go (WKWebView)
```

**Key files:**

| File | Purpose |
|------|---------|
| `main.go` | App entry, NSMenu, settings, file watching, CSS/themes |
| `core/renderer.go` | CGO bridge to Swift dylib |
| `menu.m`, `export.m`, `dragdrop.m` | ObjC bridge for native macOS features |
| `Sources/MarkdownEngine/Engine.swift` | swift-markdown AST → HTML |
| `Package.swift` | Swift package deps (apple/swift-markdown) |

## Tech Stack

- Go 1.26+ with **webview_go** for WKWebView
- Swift 5.7+ with **apple/swift-markdown** (NOT goldmark)
- CGO + ObjC bridging

## Development

- Run directly: `./md-viewer` or `./md-viewer somefile.md`
- Run app bundle: `open md-viewer.app`
- Settings stored at `~/.md-viewer/config.json`

## Run / Lint / Test

No test suite or linter configured. `go build` is the only verification step.