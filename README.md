# irgo-demo

Test repo for [irgo](https://irgo.dev) — a hypermedia framework for building
native iOS, Android, desktop, and web apps from a single Go codebase
(Go + [templ](https://templ.guide/) + [Datastar](https://datastar.fyi)).

[![CI](https://github.com/joeblew999/irgo-demo/actions/workflows/build.yml/badge.svg)](https://github.com/joeblew999/irgo-demo/actions/workflows/build.yml)

## Stack

- **Go** (version from `go.mod`; Go fetches a matching toolchain itself) + **bun**/**npm** for Tailwind
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
  binary releases, so `go tool irgo dev` installs it through your OS package
  manager. `go tool irgo serve` is the dependency-free alternative.
- `gomobile`, `templ` and `air` are installed on demand by the commands that
  need them; `irgo` itself is never installed, since `go tool irgo` builds the
  version `go.mod` requires.
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

---

## Router & Handlers

### Standard Handlers (Full Page Loads)

Standard handlers return `(string, error)`. The string is HTML.

```go
import (
    "github.com/stukennedy/irgo/pkg/router"
    "github.com/stukennedy/irgo/pkg/render"
)

// Full page load
r.GET("/", func(ctx *router.Context) (string, error) {
    return renderer.Render(templates.HomePage())
})
```

### Datastar SSE Handlers

Datastar handlers return `error` only and use `ctx.SSE()` for responses.

```go
// Datastar SSE endpoint
r.DSGet("/greeting", func(ctx *router.Context) error {
    var signals struct {
        Name string `json:"name"`
    }
    ctx.ReadSignals(&signals)

    sse := ctx.SSE()
    return sse.PatchTempl(templates.Greeting(signals.Name))
})

r.DSPost("/todos", createTodo)
r.DSPut("/todos/{id}", updateTodo)
r.DSPatch("/todos/{id}", toggleTodo)
r.DSDelete("/todos/{id}", deleteTodo)
```

### Context Methods

**Input:**
- `ctx.Param("id")` - URL path parameter
- `ctx.Query("q")` - Query string parameter
- `ctx.FormValue("name")` - Form field value
- `ctx.Header("X-Custom")` - Request header
- `ctx.ReadSignals(&signals)` - Parse Datastar signals from request

**Datastar Detection:**
- `ctx.IsDatastar()` - true if Accept: text/event-stream

**SSE Output (for Datastar handlers):**
```go
sse := ctx.SSE()
sse.PatchTempl(templates.Component())      // Patch templ component
sse.PatchHTML(`<div id="x">HTML</div>`)    // Patch raw HTML
sse.PatchSignals(map[string]any{...})      // Update client signals
sse.Remove("#element-id")                   // Remove element
sse.Redirect("/new-url")                    // Navigate browser
```

**Standard Output (for full page handlers):**
- Return HTML string from handler
- `ctx.Redirect("/path")` - HTTP redirect
- `ctx.NotFound("message")` - 404 response
- `ctx.BadRequest("message")` - 400 response
- `ctx.NoContent()` - 204 response

## Templ Templates

Templ is a type-safe HTML templating language that compiles to Go.

### Basic Syntax

```go
// templates/components.templ
package templates

// Component with parameters
templ UserCard(name string, email string) {
    <div class="card">
        <h2>{ name }</h2>
        <p>{ email }</p>
    </div>
}

// Component with children
templ Card(title string) {
    <div class="card">
        <h3>{ title }</h3>
        { children... }
    </div>
}

// Usage
templ ProfilePage() {
    @Card("Profile") {
        <p>Content goes here</p>
    }
}

// Conditionals
templ Status(active bool) {
    if active {
        <span class="text-green-500">Active</span>
    } else {
        <span class="text-red-500">Inactive</span>
    }
}

// Loops
templ UserList(users []User) {
    <ul>
        for _, user := range users {
            <li>{ user.Name }</li>
        }
    </ul>
}

// Conditional attributes
templ Checkbox(checked bool) {
    <input type="checkbox" checked?={ checked }/>
}

// Dynamic classes
templ Item(done bool) {
    <span class={ "item", templ.KV("line-through", done) }>Item</span>
}

// Safe URLs
templ Link(url string) {
    <a href={ templ.SafeURL(url) }>Link</a>
}

// Raw HTML (use sparingly)
templ RawContent(html string) {
    @templ.Raw(html)
}
```

### Rendering in Handlers

```go
renderer := render.NewTemplRenderer()

// Standard handler
func handler(ctx *router.Context) (string, error) {
    return renderer.Render(templates.MyComponent(data))
}

// Datastar handler
func sseHandler(ctx *router.Context) error {
    sse := ctx.SSE()
    return sse.PatchTempl(templates.MyComponent(data))
}
```

## Datastar Patterns

This project uses **Datastar** from `https://data-star.dev/`. Key concepts:
- **Signals**: Reactive client-side state
- **SSE**: Server responses as event streams
- **`data-*` attributes**: Declarative behavior

### Signals (Client-Side State)

```go
// Initialize signals
templ Counter() {
    <div data-signals="{count: 0}">
        <span data-text="$count">0</span>
        <button data-on:click="$count++">+</button>
    </div>
}

// Two-way binding
templ SearchForm() {
    <div data-signals="{query: ''}">
        <input type="text" data-bind:query placeholder="Search..."/>
        <span data-text="$query.length + ' characters'"></span>
    </div>
}
```

### Server Requests

```go
// GET request
templ LoadButton() {
    <button data-on:click="@get('/data')">Load</button>
    <div id="result"></div>
}

// POST request
templ TodoForm() {
    <div data-signals="{title: ''}">
        <input type="text" data-bind:title placeholder="New todo"/>
        <button data-on:click="@post('/todos')">Add</button>
    </div>
    <ul id="todo-list"></ul>
}

// DELETE request
templ DeleteButton(id string) {
    <button data-on:click={ fmt.Sprintf("@delete('/todos/%s')", id) }>
        Delete
    </button>
}
```

### Event Modifiers

```go
// Debounce input (wait 300ms after typing stops)
templ SearchInput() {
    <input
        type="text"
        data-bind:query
        data-on:input__debounce.300ms="@get('/search')"
        placeholder="Search..."
    />
}

// Prevent default form submission
templ Form() {
    <form data-on:submit__prevent="@post('/submit')">
        <input type="text" data-bind:name/>
        <button type="submit">Submit</button>
    </form>
}

// Trigger once (lazy loading)
templ LazyLoad() {
    <div data-on:intersect__once="@get('/lazy-content')">
        Loading...
    </div>
}
```

### Conditional Display

```go
// Show/hide based on signal
templ Modal() {
    <div data-signals="{showModal: false}">
        <button data-on:click="$showModal = true">Open</button>
        <div data-show="$showModal" class="modal">
            <p>Modal content</p>
            <button data-on:click="$showModal = false">Close</button>
        </div>
    </div>
}

// Dynamic classes
templ TabButton(name string) {
    <button
        data-class:active="$activeTab === 'name'"
        data-on:click="$activeTab = 'name'"
    >
        { name }
    </button>
}
```

### Loading Indicators

```go
templ LoadButton() {
    <div data-signals="{loading: false}">
        <button
            data-on:click="@get('/slow-endpoint')"
            data-indicator:loading
            data-attr:disabled="$loading"
        >
            <span data-show="!$loading">Load Data</span>
            <span data-show="$loading">Loading...</span>
        </button>
    </div>
}
```

## Native Capabilities

Call platform features (haptics, clipboard, share, storage, notifications)
with one API on every platform.

From templates / Datastar expressions (promise-based `irgo.native`):

```go
templ ShareButton(text string) {
    <button data-on:click={ fmt.Sprintf("irgo.native('share.text', {text: %q})", text) }>
        Share
    </button>
}
```

From Go handlers:

```go
import "github.com/stukennedy/irgo/pkg/native"

native.Call(ctx.Context(), "haptics.impact", native.Params{"style": "light"})
```

Built-in methods: `device.info`, `haptics.impact/notification/selection`,
`clipboard.read/write`, `share.text`, `browser.open`,
`storage.get/set/remove`, `notifications.requestPermission/show`,
`toast.show` (Android). Unsupported methods return
`native.ErrNotSupported` — degrade gracefully. Register Go fallbacks with
`native.Register(method, handler)` so web/desktop work too. Custom native
features: implement `IrgoPlugin` in `ios/.../IrgoPlugins.swift` or
`android/.../IrgoPlugins.kt` and register it.

## Sessions & Cookies

Standard `http.SetCookie` session auth works on all platforms — the mobile
bridge keeps a persistent cookie jar (sessions survive app restarts). Use
`mobile.ClearCookies()` for logout.

## Script Order

The layout must load the framework JS bridge **before** Datastar (both are
served automatically — no files to copy):

```html
<script src="/_irgo/bridge.js"></script>
<script src="/static/js/datastar.js"></script>
```

## Streaming / Real-Time

SSE streams progressively on every platform, including mobile. Long-lived
handlers must watch for disconnect:

```go
r.DSGet("/live", func(ctx *router.Context) error {
    sse := ctx.SSE()
    for {
        select {
        case <-sse.Context().Done():
            return nil // client went away
        case update := <-updates:
            sse.PatchTempl(templates.LiveRow(update))
        }
    }
})
```

## Build Tags

The framework uses Go build tags to separate platform code:

```go
//go:build !desktop    // Mobile/web builds (main.go)
//go:build desktop     // Desktop builds only (main_desktop.go)
```

- `go build .` → uses `main.go` (mobile/web)
- `go build -tags desktop .` → uses `main_desktop.go`
- `go tool irgo run desktop` → automatically adds `-tags desktop`

## Common Handler Patterns

### CRUD Operations

```go
func Mount(r *router.Router) {
    // Full page - list
    r.GET("/", func(ctx *router.Context) (string, error) {
        items := db.GetItems()
        return renderer.Render(templates.ItemsPage(items))
    })

    // SSE - create
    r.DSPost("/items", func(ctx *router.Context) error {
        var signals struct {
            Name string `json:"name"`
        }
        ctx.ReadSignals(&signals)

        if signals.Name == "" {
            return ctx.SSE().PatchTempl(templates.Error("Name required"))
        }

        item := db.CreateItem(signals.Name)
        sse := ctx.SSE()
        sse.PatchTempl(templates.ItemRow(item))
        sse.PatchSignals(map[string]any{"name": ""}) // Clear input
        return nil
    })

    // SSE - update
    r.DSPatch("/items/{id}", func(ctx *router.Context) error {
        id := ctx.Param("id")
        item := db.ToggleItem(id)
        return ctx.SSE().PatchTempl(templates.ItemRow(item))
    })

    // SSE - delete
    r.DSDelete("/items/{id}", func(ctx *router.Context) error {
        id := ctx.Param("id")
        db.DeleteItem(id)
        return ctx.SSE().Remove("#item-" + id)
    })
}
```

### Validation Errors

```go
r.DSPost("/register", func(ctx *router.Context) error {
    var signals struct {
        Email string `json:"email"`
    }
    ctx.ReadSignals(&signals)

    if !isValidEmail(signals.Email) {
        return ctx.SSE().PatchTempl(templates.FieldError("email", "Invalid email"))
    }

    // Success - redirect to dashboard
    return ctx.SSE().Redirect("/dashboard")
})
```

## Datastar Attribute Reference

| Attribute | Description | Example |
|-----------|-------------|---------|
| `data-signals` | Initialize signals | `data-signals="{count: 0}"` |
| `data-bind:X` | Two-way binding | `data-bind:name` |
| `data-text` | Text content | `data-text="$count"` |
| `data-show` | Show/hide | `data-show="$visible"` |
| `data-class:X` | Conditional class | `data-class:active="$isActive"` |
| `data-attr:X` | Dynamic attribute | `data-attr:disabled="$loading"` |
| `data-on:event` | Event handler | `data-on:click="@get('/data')"` |
| `data-indicator:X` | Loading indicator | `data-indicator:loading` |

### HTTP Actions

| Expression | Description |
|------------|-------------|
| `@get('/url')` | GET request |
| `@post('/url')` | POST request |
| `@put('/url')` | PUT request |
| `@patch('/url')` | PATCH request |
| `@delete('/url')` | DELETE request |

### Event Modifiers

| Modifier | Description |
|----------|-------------|
| `__prevent` | Prevent default |
| `__stop` | Stop propagation |
| `__once` | Trigger once |
| `__debounce.Xms` | Debounce (e.g., `__debounce.300ms`) |
| `__throttle.Xms` | Throttle (e.g., `__throttle.100ms`) |

## Tips

1. **Always read files before editing** - understand existing code first
2. **You do not need to run templ by hand** — every build and `go tool irgo dev` regenerates `_templ.go` and the stylesheet
3. **Use `go tool irgo dev`** during development for hot reload
4. **Return HTML fragments via SSE**, not JSON - this is hypermedia-driven
5. **Elements need IDs** for Datastar to patch them
6. **Use signals for client state** - avoid unnecessary server roundtrips
7. **Prefer small, focused components** that can be reused and patched independently
8. **Test in desktop mode** with `go tool irgo run desktop --dev` for browser devtools
