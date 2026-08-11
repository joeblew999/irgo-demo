package handlers_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"irgo-demo/app"
)

// The page must come back in the visitor's language, and — the case with no
// symptom — must NOT come back in some other translation when there is no
// catalog for what they asked for.
//
// x/text's matcher never fails. Given de, en, fr and ja and asked for pt-BR,
// it returns whichever catalog is first and reports language.No beside it. So
// `reader, _ := tokibundle.Match(...)` would serve German to a Brazilian
// visitor, and nothing would say so: no error, no log, a page that renders
// perfectly. lang.For goes through i18n.Reader, which reads that confidence.
//
// The pt-BR row is the one that matters here. The rest would pass either way.
func TestHomePageSpeaksTheVisitorsLanguage(t *testing.T) {
	r := app.NewRouter()
	for _, tc := range []struct{ accept, want string }{
		{"de", "Verbindung"},
		{"fr", "Connexion"},
		{"ja", "接続"},
		{"en", "Connection"},
		{"pt-BR", "Connection"}, // no catalog: source language, NOT the first one
		{"de-AT,de;q=0.9", "Verbindung"},
		{"!!!,fr", "Connexion"},
	} {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Language", tc.accept)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if !strings.Contains(w.Body.String(), tc.want) {
			t.Errorf("FAIL %-16s %q not in the page", tc.accept, tc.want)
			continue
		}
		t.Logf("ok   %-16s -> %s", tc.accept, tc.want)
	}
}
