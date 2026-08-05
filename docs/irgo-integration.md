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

---

## 1. Version pinning & the fork

| Piece | Where | Value |
|---|---|---|
| `go.mod` require | `myapp/go.mod` | `github.com/stukennedy/irgo v0.4.0` |
| `go.mod` replace | `myapp/go.mod` | `=> github.com/joeblew999/irgo v0.4.0-androidapi21.2` |
| CLI install | `mise/tasks/dev.toml` → `env:tools` | built from `$IRGO_FORK_DIR` (`go install ./cmd/irgo`), released `v0.3.1` fallback |

Fork tags (newest last, all additive on top of upstream `main`):

- `v0.4.0-androidapi21` — `-androidapi 21` on the Android bind
- `v0.4.0-androidapi21.1` — `./gradlew` cwd fix in `runAndroid`
- `v0.4.0-androidapi21.2` — template fixes (local `Response` type, gradlew quotes, shipped `gradle-wrapper.jar`)

## 2. Workaround inventory

| # | Workaround | Where | Root cause | Delete when |
|---|---|---|---|---|
| 1 | `go.mod` `replace` → fork | `myapp/go.mod` | Android fixes unreleased upstream | PR #10 merges **and** a release includes it → bump `require`, drop `replace` |
| 2 | `irgo build android \|\| true` + manual `gomobile bind -target android -androidapi 21` | `build:android` task | irgo's `buildAndroid` omits `-androidapi` (gomobile defaults to API 16, NDK r26/r27 reject) — [#9](https://github.com/stukennedy/irgo/issues/9) | PR #10 lands + CLI refreshed → replace with plain `irgo build android` |
| 3 | NDK r26 pin (`ANDROID_NDK_HOME` forced in task, overriding CI's NDK) | `env:setup:android` + `build:android` | same API-16 default: only r26 still accepts it; r27+ dropped it | same as #2 — keep r26 as the *chosen* version, drop the "must be r26" constraint |
| 4 | Local `Response` type re-exported from `myapp/mobile/mobile.go` | `myapp/mobile/mobile.go` (+ fork template) | **golang.org/x/mobile** limitation: gobind silently drops types from dependency modules — not an irgo bug | **PERMANENT** — keep the local type; only the template source moves upstream |
| 5 | irgo CLI installed from fork (`IRGO_FORK_DIR`) | `env:tools` | released CLI lacks the fixes | released irgo ≥ fixes → plain `go install github.com/stukennedy/irgo/cmd/irgo@<tag>`; `IRGO_FORK_DIR` becomes optional |
| 6 | Demo's `android/Example/gradlew` (no literal quotes in `DEFAULT_JVM_OPTS`) + committed `gradle-wrapper.jar` | `myapp/android/Example/` | predates the fork template fix (PR #10) | re-scaffold `android/Example` from fixed templates once merged |
| 7 | go.work "self-heal" block (identical ~20 lines) | **6 tasks**: `build:ios`, `build:ios:sim`, `run:ios`, `run:ios:dev`, `build:ios:prod`, `build:android` | irgo writes `go.work` → temp `x/mobile` clone; macOS cleans the temp dir; irgo never validates | **Phase 1 candidate**: file upstream issue, fix in fork so irgo handles its own `go.work` → delete all 6 blocks |
| 8 | Windows direct-toolchain build (absolute `GO`/`TEMPL`/`GCC`/`GXX`, `templ generate` directly) | `build:assets` + `build:desktop:windows` | git-bash on Windows strips the `Path` env var → native→native `exec.LookPath` fails | **PERMANENT** (environment limitation, not irgo) |
| 9 | AVD `<build>` placeholder normalization (`sed` on `config.ini`) | `env:setup:emulator` | avdmanager quirk writing `<build>` placeholders | **PERMANENT** (tooling quirk) |
| 10 | JDK 17 detection (~15 lines ×3) | `env:setup:android`, `env:setup:emulator`, `run:android` | Gradle 8.2/AGP 8.2 won't run on JDK 21+; macOS `/usr/bin/java` is a stub | **Not a workaround — a requirement.** Just de-duplicate (Phase 1) |
| 11 | sdkmanager/avdmanager location loop (~10 lines ×2) | `env:setup:android`, `env:setup:emulator` | brew-cask shims point at the wrong SDK root → broken AVDs | **Not a workaround.** De-duplicate (Phase 1) |

## 3. Local vs CI boundary (what this machine does vs what CI proves)

**Local (this machine — macOS):** web, macOS desktop, iOS simulator.
**CI-only:** Linux desktop, Windows desktop (native build **+ smoke run**), Android AAR, iOS device (gated on `IOS_TEAM_ID`).
**Opt-in local:** Android — `mise run env:setup:android` (+ `env:setup:emulator` for the emulator), fully reversible with `mise run env:uninstall:android`.

CI is the contract for "it works": see the matrix in the README. The Windows
job smoke-tests the exe precisely because a cross-compiled exe that was never
run is worse than useless.

## 4. Upstream merge-day checklist (when PR #10 merges)

1. Sync fork: `git fetch upstream && git rebase upstream/main && git push` (drop any now-merged commits).
2. Install the released CLI — change `env:tools` to `go install github.com/stukennedy/irgo/cmd/irgo@<new-tag>` (keep `IRGO_FORK_DIR` as an optional override, not a default).
3. Simplify `build:android`: replace `irgo build android || true` + manual `gomobile bind` with plain `irgo build android`.
4. Revisit the NDK pin: keep r26 as the chosen version, but the hard constraint is gone (see #3).
5. Re-scaffold `android/Example` from the fixed templates (kills workaround #6).
6. `go.mod`: bump `require` to the released version, **delete the `replace`**.
7. Run full CI; verify all green (including the Windows smoke run).
8. Update this doc: resolve rows #1–#3, #5, #6; keep #4, #8, #9 as permanent.
9. File the remaining upstream issues: go.work self-heal (#7) and the scaffolded-`output.css` follow-up noted in PR #10.
