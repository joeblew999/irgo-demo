package webtest_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/stukennedy/irgo/pkg/browsertest"
)

// The browser build must translate too, and it is the one target where nothing
// else could tell us.
//
// Every other target negotiates from Accept-Language, which arrives on the
// request. In the browser the app runs inside a service worker answering
// same-origin fetches the browser generated itself, and those carry no language
// preference — the preference is on navigator, not on the request. So
// pkg/i18n.FromBrowser reads navigator.languages there instead, and until this
// test that code had never executed.
//
// Served from the built output rather than from a Go handler on purpose. A
// handler would exercise the server path and prove nothing about wasm; this
// registers the real service worker, instantiates the real binary, and reads
// what the page actually shows.
func TestBrowserBuildSpeaksTheVisitorsLanguage(t *testing.T) {
	const dir = "../build/web"
	if _, err := os.Stat(dir + "/app.wasm"); err != nil {
		t.Skip("no browser build — run: irgo app build web")
	}

	for _, tc := range []struct{ locale, want string }{
		{"de-DE", "Servergesteuerte Hypermedia für Go"},
		{"fr-FR", "Hypermédia piloté par le serveur pour Go"},
		{"ja-JP", "Go のためのサーバー駆動ハイパーメディア"},
		{"en-US", "Server-driven hypermedia for Go"},
		// No catalog for this one. It must fall back to the source language,
		// not to whichever catalog happens to be first.
		{"pt-BR", "Server-driven hypermedia for Go"},
	} {
		t.Run(tc.locale, func(t *testing.T) {
			p := browsertest.OpenAs(t, http.FileServer(http.Dir(dir)), tc.locale)

			// The first load registers the worker and reloads once, so the
			// tagline is not there yet. Waiting on the text rather than a
			// fixed sleep: worker install plus wasm instantiation is slow and
			// varies, and a sleep long enough to be reliable is long enough to
			// make five of these tedious.
			p.MustEventually(
				`document.querySelector('.tagline')?.textContent.includes(`+jsString(tc.want)+`)`,
				"the service worker serves the page from the wasm binary, and it "+
					"must render in the browser's language — read from navigator, "+
					"because a worker's own fetches carry no Accept-Language")
		})
	}
}

// jsString quotes for embedding in the page expression.
func jsString(s string) string {
	out := `"`
	for _, r := range s {
		switch r {
		case '"', '\\':
			out += `\` + string(r)
		default:
			out += string(r)
		}
	}
	return out + `"`
}
