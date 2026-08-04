# irgo-demo

Test repo for [irgo](https://irgo.dev) — a hypermedia framework for building
native iOS, Android, desktop, and web apps from a single Go codebase
(Go + [templ](https://templ.guide/) + [Datastar](https://datastar.fyi)).

## Stack

- **Go 1.26.5** + **node LTS** — managed via [mise](https://mise.jdx.dev) (OS-agnostic toolchain)
- `myapp/` — the irgo scaffold (`irgo new myapp`)

## Prerequisites

- [mise](https://mise.jdx.dev/getting-started.html) — installs the pinned toolchain
  (`go`, `node`, `bun`) automatically
- **Xcode** (from the App Store, for iOS) / **Android Studio** (for Android)
- The rest is handled by `mise run setup` (see below)

## First-time setup (new dev machine)

```bash
mise install        # install pinned toolchain (go, node, bun)
mise run setup      # install Go tools (irgo, templ, air, gomobile) + OS packages (entr, mingw-w64 via brew)
```

`setup` is idempotent — safe to re-run any time; it detects the platform and
only installs what's missing.

## Quickstart

```bash
mise run serve      # run web server without file watching (no entr needed)
```

Then open http://localhost:8080.

For hot reload (requires `entr` — installed by `mise run setup`):

```bash
mise run dev
```

## Mise tasks

| Task            | Description                                        |
|-----------------|----------------------------------------------------|
| `mise run setup`| One-time: install all dev deps (Go tools + OS pkgs)|
| `mise run dev`  | Dev server with hot reload (needs `entr`)          |
| `mise run serve`| Web server, no file watching (pure mise)           |
| `mise run css`  | Rebuild Tailwind CSS (`input.css` → `output.css`)  |
| `mise run templ`| Generate `_templ.go` from `*.templ`                |
| `mise run test` | `go test ./...`                                    |
| `mise run build:web`     | Build web server binary                    |
| `mise run build:desktop` | Build desktop app for current platform (auto-detects OS) |
| `mise run build:desktop:windows` | Build Windows exe (cross-compiles on macOS) |
| `mise run build:desktop:linux`   | Build Linux app (native Linux only)     |
| `mise run build:desktop:all`     | Build every app the current OS supports |
| `mise run build:ios`     | Build iOS framework (`Irgo.xcframework`, needs Xcode) |
| `mise run build:ios:app` | Build iOS app for the simulator          |
| `mise run run:ios`       | Build + launch iOS app on the Simulator  |

## Platform targets

| Platform | How it runs                                   | Command            |
|----------|-----------------------------------------------|--------------------|
| Web      | Standard HTTP server                          | `mise run serve`   |
| Desktop  | Real localhost server + native webview (CGO)  | `mise run build:desktop` |
| iOS      | gomobile, in-process Go (virtual HTTP)        | `irgo build ios`   |
| Android  | gomobile, in-process Go (virtual HTTP)        | `irgo build android` |

## Notes

- `irgo dev` (hot reload) shells out to `entr`, which is source-only C with no
  binary releases — so it can't be installed via mise. `mise run serve` is the
  dependency-free alternative.
- `gomobile`, `templ`, `air`, and `irgo` are Go tools installed via
  `go install` (mise can't manage arbitrary Go binaries) — `mise run setup`
  installs them.
- Linux desktop builds need GTK3 + WebKit2GTK and can't be cross-compiled from
  macOS — build on a Linux machine/CI or use a Docker image.
- Keep the templ **generator** and templ **library** versions in sync
  (`go.mod`); otherwise generated code may reference APIs the library lacks.
- App Store / Play Store distribution (signing, metadata, review) is outside
  irgo's scope — see the [irgo deployment docs](https://irgo.dev/docs/deployment).
