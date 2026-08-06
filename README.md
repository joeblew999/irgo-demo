# irgo-demo

Test repo for [irgo](https://irgo.dev) — a hypermedia framework for building
native iOS, Android, desktop, and web apps from a single Go codebase
(Go + [templ](https://templ.guide/) + [Datastar](https://datastar.fyi)).

[![CI](https://github.com/joeblew999/irgo-demo/actions/workflows/build.yml/badge.svg)](https://github.com/joeblew999/irgo-demo/actions/workflows/build.yml)

## Stack

- **Go 1.26.5** + **node LTS** — managed via [mise](https://mise.jdx.dev) (OS-agnostic toolchain)
- **This repo is the generated app.** The root *is* `irgo new` output — never
  hand-edit it; change the CLI templates and regenerate with
  `go tool irgo`. CI enforces that: the `regen` job fails if regenerating
  changes a tracked file.
- `.github/workflows/build.yml` and `release.yml` are **scaffolded by `irgo ci`
  and committed unmodified** — they are what every irgo project gets, so this
  repo proves they work on every push. Demo-only checks live in
  `maintainer.yml`.

## How the CLI and this repo stay lined up

`go.mod` is the only pin, and it is the one every Go project already has:

```
require github.com/stukennedy/irgo v0.4.0
tool    github.com/stukennedy/irgo/cmd/irgo
replace github.com/stukennedy/irgo => github.com/joeblew999/irgo v0.4.0-androidapi21.58
```

`go tool irgo` builds exactly that version — through the `replace`, so this repo
tracks a fork without anything being installed, put on `PATH`, or kept in step
by hand. Move the pin with `go mod edit -replace`, and CI, your shell and the
generated app all follow together.

There is no mise, no bootstrap script and no setup step.

## Prerequisites

- **Go** (the version in `go.mod`; Go downloads a matching toolchain itself)
- **bun** or **npm**, for the Tailwind build
- **Xcode** (App Store) only if you want iOS

Everything else installs itself on first use — see below.

## Getting started

```bash
go tool irgo doctor    # what can this machine build?
go tool irgo dev       # hot-reload server on :8080
```

Because this repo tracks a fork whose module path the public proxy cannot
serve, set `GOPRIVATE` once:

```bash
export GOPRIVATE='github.com/joeblew999/*'
```

A project on a published release needs neither that nor the `replace`.

**Nothing else to remember, and no ordering to get right.** The CLI provisions
what a command needs, when that command needs it:

- `go tool irgo build android` installs JDK 17, SDK, NDK (and the emulator +
  AVD for `run`) under `~/.irgo` and `ANDROID_HOME` — no brew/apt/winget,
  and no Android Studio. Clean down with
  `go tool irgo uninstall-tools android --remove-jdk`.
- Builds regenerate `_templ.go` and the stylesheet themselves — both are
  gitignored but embedded.
- `build`/`run` scaffold `ios/Example` and `android/Example` from the CLI's
  templates. They are not checked in.
- iOS targets refuse to run off macOS with a clear error; `build all` skips
  what the host cannot do.

## Commands

`go tool irgo help` is the reference. `irgo doctor` reports what this host can
build and what needs a manual step.

| Command | Description |
|---|---|
| `go tool irgo doctor [--fix]` | What this host can build; `--fix` repairs the Xcode setup |
| `go tool irgo dev` / `serve` | Web server; `dev` adds hot reload |
| `go tool irgo build ios --sim` | Runnable iOS Simulator app |
| `go tool irgo run ios --device` | Build, install and launch on a USB iPhone |
| `go tool irgo build android` | Android AAR — self-provisions the toolchain |
| `go tool irgo build desktop all` | Every desktop app this host supports |
| `go tool irgo upgrade` | Refresh framework files, leaving your code alone |
| `go tool irgo ci` | Scaffold the GitHub Actions workflows |
| `go tool irgo clean [--all]` | Remove generated output |
| `go tool irgo uninstall [ios\|android\|desktop]` | Remove the installed app |
| `go tool irgo package setup --check` | What each store needs, and how to supply it |

### This repo is generated output

The root **is** `irgo new` output. Never hand-edit it — change the CLI
templates and regenerate. CI enforces that: the `regen` job fails if
regenerating changes a tracked file.

`.github/workflows/build.yml` and `release.yml` are scaffolded by `irgo ci` and
committed **unmodified**, so every push exercises what a developer actually
gets. Demo-only checks live in `maintainer.yml`.

## Platform targets

| Platform | How it runs                                   | Command            |
|----------|-----------------------------------------------|--------------------|
| Web      | Standard HTTP server                          | `go tool irgo serve` |
| Desktop  | Real localhost server + native webview (CGO)  | `irgo build desktop` |
| iOS      | gomobile, in-process Go (virtual HTTP)        | `irgo build ios`   |
| Android  | gomobile, in-process Go (virtual HTTP)        | `irgo build android` |

**Where each target actually runs.** Android no longer needs an opt-in local
setup step — the CLI provisions it on first build. `irgo help` is the reference
for everything the CLI does.

| Target | Locally on macOS | CI (per push) |
|--------|:---:|---|
| Web | ✅ `irgo serve` | `check` — tests + web binary |
| macOS desktop | ✅ `irgo build desktop` | `macos-desktop` |
| Windows desktop | opt-in cross-compile (mingw-w64) | `windows-desktop` — native build **+ smoke run** |
| Linux desktop | — (needs GTK/WebKit) | `linux-desktop` (ubuntu-22.04) |
| iOS simulator | ✅ `irgo run ios` | `ios` — framework + simulator app |
| iOS device | opt-in (needs signing) | `ios-device` — only when `IOS_TEAM_ID` set |
| Android | ✅ `irgo run android` (self-provisions) | `android` — AAR only (no emulator: slow and fragile in CI) |

## Notes

- `irgo dev` (hot reload) shells out to `entr`, which is source-only C with no
  binary releases — so it can't be installed via mise. `irgo serve` is the
  dependency-free alternative.
- `gomobile`, `templ`, `air`, and `irgo` are Go tools installed via
  `go install` (mise can't manage arbitrary Go binaries) — `go tool irgo`
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
