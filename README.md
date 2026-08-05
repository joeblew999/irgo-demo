# irgo-demo

Test repo for [irgo](https://irgo.dev) — a hypermedia framework for building
native iOS, Android, desktop, and web apps from a single Go codebase
(Go + [templ](https://templ.guide/) + [Datastar](https://datastar.fyi)).

[![CI](https://github.com/joeblew999/irgo-demo/actions/workflows/build.yml/badge.svg)](https://github.com/joeblew999/irgo-demo/actions/workflows/build.yml)

## Stack

- **Go 1.26.5** + **node LTS** — managed via [mise](https://mise.jdx.dev) (OS-agnostic toolchain)
- `myapp/` — **generated output**, not source. It is produced wholesale by
  `irgo new myapp`. Never hand-edit it; change the CLI templates and regenerate.

## How the CLI and this repo stay lined up

One value, in `mise.toml`, is the single source of truth:

```toml
[env]
IRGO_REPLACE = "github.com/joeblew999/irgo v0.4.0-androidapi21.24"
```

Both sides read it:

- `irgo new myapp` writes it into `myapp/go.mod` as the `replace` directive
- `mise run env:tools` builds that exact CLI version

So the CLI that generated the app and the CLI you run against it cannot drift.
It lives in `mise.toml` rather than `myapp/go.mod` because `myapp/` is
regenerated wholesale — a pin kept inside it would be destroyed by the very
command that rebuilds it, and it was previously the only record of which CLI
built the project.

To move to a new fork build: bump `IRGO_REPLACE`, then `mise run env:tools`.

## Prerequisites

- [mise](https://mise.jdx.dev/getting-started.html) — installs the pinned
  toolchain (`go`, `node`, `bun`) and the matching irgo CLI
- **Xcode** (App Store) for iOS. Nothing else is required up front.

Everything else installs itself on first use — see below.

## First-time setup (new dev machine)

```bash
mise install          # pinned toolchain (go, node, bun)
mise run env:tools    # the irgo CLI at IRGO_REPLACE + its tools (templ, air, gomobile)
mise run env:setup    # OPTIONAL: OS packages (entr for hot reload, mingw-w64 for Windows cross-compile)
```

Both are idempotent — safe to re-run any time.

**Nothing else to remember, and no ordering to get right.** The CLI provisions
what a command needs, when that command needs it:

- `irgo build android` / `irgo run android` install JDK 17, SDK, NDK (and the
  emulator + AVD for `run`) under `~/.irgo` and `ANDROID_HOME` — no
  brew/apt/winget. Clean down with `irgo uninstall-tools android --remove-jdk`;
  verify with `irgo doctor android`.
- `irgo build ios|android` runs `templ` itself (`_templ.go` is gitignored but
  the mobile package imports it) and installs `gomobile`/`gobind` if missing.
- `irgo build|run ios|android` scaffold `ios/Example` / `android/Example` from
  the CLI's embedded templates (missing-only, idempotent). They are not checked in.
- iOS targets refuse to run off macOS with a clear error instead of failing
  later on a missing `xcodebuild`; `irgo build all` just skips iOS off-macOS.

## Quickstart

```bash
cd myapp
irgo serve       # web server, no file watching
irgo dev         # hot reload (needs entr)
```

Then open http://localhost:8080.

## Commands

The CLI is the reference — `irgo help` and `irgo help <command>` are the source
of truth, including the `IRGO_REPLACE` / `IRGO_PATH` contract in `irgo help new`.
Run these from `myapp/`:

| Command | Description |
|---|---|
| `irgo serve` / `irgo dev` | Web server; `dev` adds hot reload |
| `irgo build ios` | iOS framework (`Irgo.xcframework`) — macOS only |
| `irgo build ios --sim` | Runnable iOS Simulator `.app` (scaffolds + drives xcodebuild) |
| `irgo run ios` / `--dev` | Launch on the Simulator; `--dev` hot-reloads |
| `irgo build android` | Android framework (`irgo.aar`) — self-provisions the toolchain |
| `irgo run android` / `--no-window` | Launch on the emulator; `--no-window` is headless |
| `irgo build all` | Both mobile frameworks (skips iOS off-macOS) |
| `irgo build desktop` | Desktop app for the current OS |
| `irgo doctor android` | Verify the Android toolchain |
| `irgo uninstall-tools android --remove-jdk` | Remove everything the CLI installed |
| `irgo new myapp` | Regenerate the app from the CLI templates |

### What mise is still for

mise's remaining job is lining this repo up with the CLI, plus the OS-level
packages the CLI does not own yet:

| Task | Description |
|---|---|
| `mise run env:tools` | Install the irgo CLI pinned by `IRGO_REPLACE` + its Go tools |
| `mise run env:setup` | OS packages: entr, mingw-w64 (brew/apt/pacman) |
| `mise run env:uninstall` | The exact inverse of `env:setup` |
| `mise run build:assets` | `bun install` + templ + Tailwind CSS (has a Windows git-bash workaround) |
| `mise run build:desktop:{windows,linux,all}` | Cross-compile logic (mingw-w64, GTK/WebKit gating) |

The remaining `build:*` / `run:*` tasks are thin passthroughs to the equivalent
`irgo` command and exist for discoverability; prefer calling `irgo` directly.

## Platform targets

| Platform | How it runs                                   | Command            |
|----------|-----------------------------------------------|--------------------|
| Web      | Standard HTTP server                          | `irgo serve`       |
| Desktop  | Real localhost server + native webview (CGO)  | `irgo build desktop` |
| iOS      | gomobile, in-process Go (virtual HTTP)        | `irgo build ios`   |
| Android  | gomobile, in-process Go (virtual HTTP)        | `irgo build android` |

**Where each target actually runs.** Android no longer needs an opt-in local
setup step — the CLI provisions it on first build. See
[docs/irgo-integration.md](docs/irgo-integration.md) for the workaround
inventory and upstream tracking (stukennedy/irgo#9, PR #10).

| Target | Locally on macOS | CI (per push) |
|--------|:---:|---|
| Web | ✅ `irgo serve` | `check` — tests + web binary |
| macOS desktop | ✅ `irgo build desktop` | `macos-desktop` |
| Windows desktop | opt-in cross-compile (mingw-w64) | `windows-desktop` — native build **+ smoke run** |
| Linux desktop | — (needs GTK/WebKit) | `linux-desktop` (ubuntu-22.04) |
| iOS simulator | ✅ `irgo run ios` | `ios` — framework + simulator app |
| iOS device | opt-in (needs signing) | `ios-device` — only when `IOS_TEAM_ID` set |
| Android | ✅ `irgo run android` (self-provisions) | `android` — AAR; `android-emulator` — headless run on ARM64 |

## Notes

- `irgo dev` (hot reload) shells out to `entr`, which is source-only C with no
  binary releases — so it can't be installed via mise. `irgo serve` is the
  dependency-free alternative.
- `gomobile`, `templ`, `air`, and `irgo` are Go tools installed via
  `go install` (mise can't manage arbitrary Go binaries) — `mise run env:tools`
  installs them, into `$(go env GOPATH)/bin`. Make sure that's on your `PATH`
  so `irgo` is callable directly; CI does this in `.github/actions/setup`.
- Linux desktop builds need GTK3 + WebKit2GTK and can't be cross-compiled from
  macOS — build on a Linux machine/CI or use a Docker image.
- CI (`.github/workflows/build.yml`) builds **every target** on each push —
  see the matrix above. Tag a release
  (`git tag v0.1.0 && git push origin v0.1.0`) to trigger
  `.github/workflows/release.yml`, which attaches all artifacts to a GitHub
  Release (Windows is now built **natively + smoke-tested** before it ships).
  The iOS device build runs only when signing secrets (`IOS_TEAM_ID`) are
  configured. Note that `_templ.go` and `output.css` are gitignored but
  embedded into builds, so CI regenerates them first.
- Keep the templ **generator** and templ **library** versions in sync
  (`go.mod`); otherwise generated code may reference APIs the library lacks.
- App Store / Play Store distribution (signing, metadata, review) is outside
  irgo's scope — see the [irgo deployment docs](https://irgo.dev/docs/deployment).
