// Which language this visitor reads.
//
// One place, because every entry point needs it: the page route, the SSE
// handlers that patch fragments afterwards, and every target irgo builds for.
// A helper per handler would drift, and the way it drifts is that one of them
// stops passing the request and silently serves the source language.
package lang

import (
	"net/http"

	"irgo-demo/tokibundle"

	"github.com/stukennedy/irgo/pkg/i18n"
)

// For picks the localization for this request.
//
// i18n.Reader rather than tokibundle.Match directly. The matcher never fails:
// asked for a language with no catalog it returns the first one that exists
// and reports language.No beside it, so discarding that value would serve a
// French visitor German. Reader checks it and falls back to the source
// language instead.
func For(r *http.Request) tokibundle.Reader {
	return i18n.Reader(tokibundle.Match, tokibundle.Default, i18n.Preferred(r)...)
}
