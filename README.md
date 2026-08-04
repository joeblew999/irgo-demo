# irgo-demo

Test repo for [irgo](https://irgo.dev) — a hypermedia framework for building
native iOS, Android, desktop, and web apps from a single Go codebase
(Go + [templ](https://templ.guide/) + [Datastar](https://datastar.fyi)).

## Stack

- **Go 1.26.5** + **node LTS** — managed via [mise](https://mise.jdx.dev) (OS-agnostic toolchain)
- `myapp/` — the irgo scaffold (`irgo new myapp`)

## Prerequisites

- [mise](https://mise.jdx.dev/getting-started.html) — installs Go + node automatically
- `bun` (or npm) for frontend deps
- For hot-reload dev only: [`entr`](https://github.com/eradman/entr) (`brew install entr`)

## Quickstart

```bash
mise install        # install pinned toolchain (go, node)
mise run serve      # run web server without file watching (no entr needed)
```

Then open http://localhost:8080.

For hot reload (requires `entr`):

```bash
mise run dev
```

## Mise tasks

| Task            | Description                                        |
|-----------------|----------------------------------------------------|
| `mise run dev`  | Dev server with hot reload (needs `entr`)          |
| `mise run serve`| Web server, no file watching (pure mise)           |
| `mise run css`  | Rebuild Tailwind CSS (`input.css` → `output.css`)  |
| `mise run templ`| Generate `_templ.go` from `*.templ`                |
| `mise run test` | `go test ./...`                                    |
| `mise run build:web`     | Build web server binary                            |
| `mise run build:desktop` | Build desktop app for current platform             |

## Platform targets

| Platform | How it runs                                   | Command            |
|----------|-----------------------------------------------|--------------------|
| Web      | Standard HTTP server                          | `mise run serve`   |
| Desktop  | Real localhost server + native webview (CGO)  | `irgo build desktop` |
| iOS      | gomobile, in-process Go (virtual HTTP)        | `irgo build ios`   |
| Android  | gomobile, in-process Go (virtual HTTP)        | `irgo build android` |

## Notes

- `irgo dev` (hot reload) shells out to `entr`, which is source-only C with no
  binary releases — so it can't be installed via mise. `mise run serve` is the
  dependency-free alternative.
- Keep the templ **generator** and templ **library** versions in sync
  (`go.mod`); otherwise generated code may reference APIs the library lacks.
- App Store / Play Store distribution (signing, metadata, review) is outside
  irgo's scope — see the [irgo deployment docs](https://irgo.dev/docs/deployment).
