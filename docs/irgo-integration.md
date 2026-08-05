# irgo integration & workaround inventory

Single source of truth for **how this repo depends on irgo**: what's forked,
what's worked around, and — critically — **when each workaround can be deleted**.
If you find a hack in this repo that isn't in this table, add it here (and
prefer fixing it upstream over documenting it).

Operating model: **the demo is the canary.** CI builds *and runs* every target
on each push; anything it catches gets fixed upstream (fork → PR), and the demo
only keeps workarounds upstream can't fix (or hasn't shipped yet). Tracked below.

- Upstream: [stukennedy/irgo](https://github.com/stukennedy/irgo)
- Fork: [joeblew999/irgo](https://github.com/joeblew999/irgo) (keeps the upstream
  module path so `go.mod` replace works)
- Upstream issue: [#9 — Android build pipeline](https://github.com/stukennedy/irgo/issues/9)
- Upstream PR: [#10 — Fix Android build/run pipeline end-to-end](https://github.com/stukennedy/irgo/pull/10)
- Upstream PR: [#11 — stale go.work regeneration](https://github.com/stukennedy/irgo/pull/11) (ready for review)
- Upstream PR: [#12 — cross-platform Android toolchain (install/uninstall/doctor, self-provisioning)](https://github.com/stukennedy/irgo/pull/12)

---

## 1. Version pinning & the fork

| Piece | Where | Value |
|---|---|---|
| `go.mod` require | `myapp/go.mod` | `github.com/stukennedy/irgo v0.4.0` |
| Fork pin (single source) | `mise.toml` → `[env] IRGO_REPLACE` | `github.com/joeblew999/irgo v0.4.0-androidapi21.24` |
| `go.mod` replace | `myapp/go.mod` | **generated** from `IRGO_REPLACE` by `irgo new` — never hand-edited |
| CLI install | `mise/tasks/dev.toml` → `env:tools` | local checkout (`$IRGO_FORK_DIR`) or clone of `$IRGO_FORK_TAG` (derived from `IRGO_REPLACE`), then `irgo install-tools`; released `v0.3.1` fallback |

`IRGO_REPLACE` is the one value both sides read — `irgo new` writes it into the
generated `go.mod`, `env:tools` builds that same CLI version — so the CLI and the
generated app cannot drift. It lives in `mise.toml`, not in `myapp/go.mod`,
because `myapp/` is CLI *output*: regenerating it would otherwise destroy the pin
that selects its own generator.

Fork tags (newest last, all additive on top of upstream `main`):

- `v0.4.0-androidapi21` — `-androidapi 21` on the Android bind
- `v0.4.0-androidapi21.1` — `./gradlew` cwd fix in `runAndroid`
- `v0.4.0-androidapi21.2` — template fixes (local `Response` type, gradlew quotes, shipped `gradle-wrapper.jar`)
- `v0.4.0-androidapi21.3` — stale `go.work` validation/regeneration ([PR #11](https://github.com/stukennedy/irgo/pull/11), ready)
- `v0.4.0-androidapi21.4` … `v0.4.0-androidapi21.12` — the Android toolchain feature ([PR #12](https://github.com/stukennedy/irgo/pull/12)): `install-tools android` / `uninstall-tools android` / `doctor android`, self-provisioning (`build`/`run` auto-install any missing toolchain), emulator auto-boot + `--no-window`, `go` GOROOT fallback (Windows git-bash), templ pinned to go.mod
- `v0.4.0-androidapi21.13` … `v0.4.0-androidapi21.23` — CLI-scaffolded native shells (`ios/Example`, `android/Example` from embedded templates), store packaging (`irgo package`), review monitoring (`irgo reviews`), deterministic AVD home, arm64-v8a system image on arm64 hosts + `-accel off` fallback
- `v0.4.0-androidapi21.24` — `IRGO_REPLACE` (generated fork pin, so the app is regenerable from the CLI alone), `irgo build ios --sim` (CLI drives xcodebuild), `irgo new mobile` removed (build/run self-scaffold), `irgo build` runs templ itself, iOS gated to macOS, KVM probed by access not existence

## 2. Workaround inventory

| # | Workaround | Where | Root cause | Delete when |
|---|---|---|---|---|
| 1 | `go.mod` `replace` → fork | `myapp/go.mod` | Android fixes unreleased upstream | PR #10 merges **and** a release includes it → bump `require`, drop `replace` |
| 2 | ~~`irgo build android \|\| true` + manual bind~~ — now plain `irgo build android` (fork CLI pins `-androidapi 21`) | `build:android` task | irgo's `buildAndroid` omitted `-androidapi` (API-16 default, NDK r26/r27 reject) — [#9](https://github.com/stukennedy/irgo/issues/9) | resolved in fork CLI `v0.4.0-androidapi21.3`; becomes upstream when PR #10 lands
| 3 | NDK r26 pin — now in the CLI (`ensureAndroidToolchain` applies `ANDROID_NDK_HOME` to the pinned r26) | `cmd/irgo/android_tools.go` | gomobile's API-16 default: only r26 accepts it; r27+ dropped it | keep r26 as the *chosen* version; the hard constraint dies when `-androidapi` ships upstream (#10) |
| 4 | Local `Response` type re-exported from `myapp/mobile/mobile.go` | `myapp/mobile/mobile.go` (+ fork template) | **golang.org/x/mobile** limitation: gobind silently drops types from dependency modules — not an irgo bug | **PERMANENT** — keep the local type; only the template source moves upstream |
| 5 | irgo CLI from fork — local checkout (`IRGO_FORK_DIR`) or clone of `IRGO_FORK_TAG` (from `IRGO_REPLACE`, `v0.4.0-androidapi21.24`) so CI/fresh machines get the fixed CLI | `env:tools` | released CLI lacks the fixes | released irgo ≥ fixes → plain `go install github.com/stukennedy/irgo/cmd/irgo@<tag>`; fork vars become optional |
| 6 | Demo's `android/Example/gradlew` (no literal quotes in `DEFAULT_JVM_OPTS`) + committed `gradle-wrapper.jar` | `myapp/android/Example/` | predates the fork template fix (PR #10) | re-scaffold `android/Example` from fixed templates once merged |
| 7 | ~~go.work "self-heal" block ×6~~ — **deleted**; the irgo CLI now validates/regenerates its own `go.work` | ~~6 tasks~~ (removed) | irgo wrote `go.work` → temp `x/mobile` clone; macOS cleans the temp dir; irgo never validated | resolved in fork CLI — [PR #11](https://github.com/stukennedy/irgo/pull/11), ready for review, validated in demo CI |
| 8 | Windows direct-toolchain build (absolute `GO`/`TEMPL`/`GCC`/`GXX`, `templ generate` directly) | `build:assets` + `build:desktop:windows` | git-bash on Windows strips the `Path` env var → native→native `exec.LookPath` fails | **PERMANENT** (environment limitation, not irgo) |
| 9 | AVD `<build>` placeholder normalization — **moved into the CLI** (`installEmulator`) | `cmd/irgo/android_tools.go` | avdmanager quirk writing `<build>` placeholders | Resolved — owned by the CLI (permanent tooling quirk) |
| 10 | JDK 17 — **owned by the irgo CLI**: `install-tools android` auto-downloads a managed Temurin 17 into `~/.irgo/jdks` (no brew/apt/winget); `build`/`run` self-provision it | `cmd/irgo/android_tools.go` (fork) | Gradle 8.2/AGP 8.2 won't run on JDK 21+; macOS `/usr/bin/java` is a stub | Not a workaround — a requirement. Owned by the CLI ([PR #12](https://github.com/stukennedy/irgo/pull/12)) |
| 11 | sdkmanager/avdmanager location — **moved into the irgo CLI** (`install-tools android`) | `cmd/irgo/android_tools.go` (fork) | brew-cask shims point at the wrong SDK root → broken AVDs | Not a workaround. Now owned by the CLI |

## 3. Local vs CI boundary (what this machine does vs what CI proves)

**Self-provisioning:** since PR #12, `irgo build android` / `irgo run android`
auto-install any missing JDK/SDK/NDK/emulator/AVD — devs and CI never need a
separate setup step. `irgo uninstall-tools android --remove-jdk` cleans a
machine back to zero; `irgo doctor android` verifies the toolchain.

**Local (this machine — macOS):** web, macOS desktop, iOS simulator, and
Android (opt-in, fully reversible).
**CI-only:** Linux desktop, Windows desktop (native build **+ smoke run**),
Android AAR (round-trip: install → doctor → build → uninstall), the emulator
self-provision run (bare runner → `irgo run android --no-window` → doctor →
smoke → screenshot), iOS device (gated on `IOS_TEAM_ID`).

CI is the contract for "it works": see the matrix in the README. The Windows
job smoke-tests the exe precisely because a cross-compiled exe that was never
run is worse than useless.

## 4. Upstream merge-day checklist (when PR #10 merges)

1. Sync fork: `git fetch upstream && git rebase upstream/main && git push` (drop any now-merged commits).
2. Install the released CLI — change `env:tools` to `go install github.com/stukennedy/irgo/cmd/irgo@<new-tag>` (keep `IRGO_FORK_DIR`/`IRGO_FORK_TAG` as optional overrides, not defaults).
3. `build:android` is already plain `irgo build android` (fork CLI pins `-androidapi 21`) — just confirm it still holds.
4. Revisit the NDK pin: keep r26 as the chosen version, but the hard constraint is gone (see #3).
5. Re-scaffold `android/Example` from the fixed templates (kills workaround #6).
6. `go.mod`: bump `require` to the released version, **delete the `replace`**.
7. Run full CI; verify all green (including the Windows smoke run).
8. Update this doc: resolve rows #1–#3, #5, #6; keep #4, #8, #9 as permanent.
9. PR #11 (go.work) and PR #12 (toolchain commands) are ready/validated — on merge, delete the `go.mod` replace and simplify `env:tools` to `go install github.com/stukennedy/irgo/cmd/irgo@<tag>`; file any remaining upstream issues (e.g. the scaffolded-`output.css` follow-up noted in PR #10).
