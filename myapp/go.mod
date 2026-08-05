module myapp

go 1.26.5

// Use the joeblew999/irgo fork: it carries the -androidapi fix for gomobile
// bind (unreleased upstream, see stukennedy/irgo#9) plus the gradlew path fix.
// The fork keeps the original module path so this replace works transparently.
require github.com/stukennedy/irgo v0.4.0

replace github.com/stukennedy/irgo => github.com/joeblew999/irgo v0.4.0-androidapi21.1

require github.com/a-h/templ v0.3.977

require (
	github.com/CAFxX/httpcompression v0.0.9 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/go-chi/chi/v5 v5.2.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/starfederation/datastar-go v1.1.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/webview/webview_go v0.0.0-20240831120633-6173450d4dd6 // indirect
)
