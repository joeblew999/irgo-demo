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
- The [mise VS Code extension](https://marketplace.visualstudio.com/items?itemName=jdx.mise)
  (`jdx.mise`) — task runner + tool integration in the editor
- **Xcode** (from the App Store, for iOS) / **Android Studio** (for Android)
- The rest is handled by `mise run env:setup` (see below)

## First-time setup (new dev machine)

```bash
mise install        # install pinned toolchain (go, node, bun)
mise run env:setup  # install Go tools (irgo, templ, air, gomobile) + OS packages (entr, mingw-w64 via brew)
```

`env:setup` is idempotent — safe to re-run any time; it detects the platform and
only installs what's missing. Everything it installs can be removed with
`mise run env:uninstall`.

## Quickstart

```bash
mise run serve      # run web server without file watching (no entr needed)
```

Then open http://localhost:8080.

For hot reload (requires `entr` — installed by `mise run env:setup`):

```bash
mise run serve:dev
```

## Mise tasks

Tasks are namespaced: `env:` (toolchain), `serve:` (web server), `build:` (builds),
`run:` (launches), `check:` (validation).

| Task            | Description                                        |
|-----------------|----------------------------------------------------|
| `mise run env:setup`| One-time: install all dev deps (Go tools + OS pkgs)|
| `mise run env:tools`| Install pinned Go tools (irgo/templ/air/gomobile) |
| `mise run env:uninstall`| Remove everything `env:setup` installed (inverse) |
| `mise run env:setup:android`| Install Android SDK + NDK (brew + sdkmanager) |
| `mise run env:uninstall:android`| Remove Android SDK + cmdline-tools |
| `mise run serve`| Web server, no file watching (pure mise)           |
| `mise run serve:dev`| Dev server with hot reload (needs `entr`)        |
| `mise run css`  | Rebuild Tailwind CSS (`input.css` → `output.css`)  |
| `mise run css:watch`| Watch Tailwind CSS sources                    |
| `mise run build:web`     | Build web server binary                    |
| `mise run build:templ`   | Generate `_templ.go` from `*.templ`         |
| `mise run build:assets`  | Generate templ + build CSS (run before any build) |
| `mise run build:desktop` | Build desktop app for current platform (auto-detects OS) |
| `mise run build:desktop:windows` | Build Windows exe (cross-compiles on macOS) |
| `mise run build:desktop:linux`   | Build Linux app (native Linux only)     |
| `mise run build:desktop:all`     | Build every app the current OS supports |
| `mise run build:ios`     | Build iOS framework (`Irgo.xcframework`, needs Xcode) |
| `mise run build:ios:sim` | Build iOS app for the simulator (Debug)  |
| `mise run build:ios:prod`| Build iOS app for device/App Store (Release) |
| `mise run build:android` | Build Android framework (`irgo.aar`, needs SDK + NDK) |
| `mise run build:mobile`  | Build all mobile frameworks (iOS + Android) |
| `mise run run:ios`       | Build + launch iOS app on the Simulator  |
| `mise run run:ios:dev`   | Launch iOS app with hot reload            |
| `mise run run:desktop`   | Build + launch desktop app (native webview) |
| `mise run run:desktop:dev`| Launch desktop app with devtools        |
| `mise run check:test`    | Run Go tests                             |

## Platform targets

| Platform | How it runs                                   | Command            |
|----------|-----------------------------------------------|--------------------|
| Web      | Standard HTTP server                          | `mise run serve`   |
| Desktop  | Real localhost server + native webview (CGO)  | `mise run build:desktop` |
| iOS      | gomobile, in-process Go (virtual HTTP)        | `mise run build:ios` |
| Android  | gomobile, in-process Go (virtual HTTP)        | `mise run build:android` |

## Notes

- `irgo dev` (hot reload) shells out to `entr`, which is source-only C with no
  binary releases — so it can't be installed via mise. `mise run serve` is the
  dependency-free alternative.
- `gomobile`, `templ`, `air`, and `irgo` are Go tools installed via
  `go install` (mise can't manage arbitrary Go binaries) — `mise run env:setup`
  installs them.
- Linux desktop builds need GTK3 + WebKit2GTK and can't be cross-compiled from
  macOS — build on a Linux machine/CI or use a Docker image.
- CI (`.github/workflows/build.yml`) builds **every target** on each push —
  Go tests, web binary, Linux/macOS/Windows desktop, iOS framework + simulator,
  and the Android AAR. Tag a release (`git tag v0.1.0 && git push origin v0.1.0`)
  to trigger `.github/workflows/release.yml`, which attaches all artifacts to a
  GitHub Release. The iOS device build runs only when signing secrets
  (`IOS_TEAM_ID`) are configured. Note that `_templ.go` and `output.css` are
  gitignored but embedded into builds, so CI regenerates them first.
- Keep the templ **generator** and templ **library** versions in sync
  (`go.mod`); otherwise generated code may reference APIs the library lacks.
- App Store / Play Store distribution (signing, metadata, review) is outside
  irgo's scope — see the [irgo deployment docs](https://irgo.dev/docs/deployment).
