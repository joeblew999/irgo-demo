// Package mobile provides gomobile bindings for iOS and Android.
// This package re-exports the irgo mobile bridge functions and initializes
// your app's routes when the mobile app starts.
package mobile

import (
	"fmt"
	"net/http"

	"myapp/app"
	irgomobile "github.com/stukennedy/irgo/mobile"
)

var appRouter http.Handler

// Response is a gomobile-compatible response type.
// This is defined here so gomobile can export it.
type Response struct {
	Status  int
	Headers string
	Body    []byte
}

// BodyString returns the body as a string.
func (r *Response) BodyString() string {
	return string(r.Body)
}

// Initialize sets up the mobile bridge and app routes.
// Called automatically by native code at app startup.
func Initialize() {
	// Set up app routes using shared router setup
	r := app.NewRouter()
	appRouter = r.Handler()

	// Initialize the irgo bridge with our handler
	irgomobile.SetHandler(appRouter)
	irgomobile.Initialize()

	fmt.Println("myapp mobile initialized")
}

// HandleRequest processes an HTTP request from the WebView.
// This is called by native code (Swift/Kotlin) for each request.
func HandleRequest(method, url, headers string, body []byte) *Response {
	coreResp := irgomobile.HandleRequest(method, url, headers, body)
	return &Response{
		Status:  coreResp.Status,
		Headers: coreResp.Headers,
		Body:    coreResp.Body,
	}
}

// HandleRequestSimple processes a simple GET request.
func HandleRequestSimple(method, url string) *Response {
	coreResp := irgomobile.HandleRequestSimple(method, url)
	return &Response{
		Status:  coreResp.Status,
		Headers: coreResp.Headers,
		Body:    coreResp.Body,
	}
}

// RenderInitialPage returns the initial HTML for the WebView.
func RenderInitialPage() string {
	return irgomobile.RenderInitialPage()
}

// IsReady returns true if the bridge is initialized.
func IsReady() bool {
	return irgomobile.IsReady()
}

// Shutdown cleans up the bridge.
func Shutdown() {
	irgomobile.Shutdown()
}

// WebSocketCallback is implemented by Swift/Kotlin to receive WebSocket messages.
// Mirrors irgo/mobile.WebSocketCallback so gomobile can export it to Android/iOS.
type WebSocketCallback interface {
	// OnMessage is called when Go has a message to send to the WebView.
	// data is a JSON-encoded websocket.Envelope.
	OnMessage(sessionID string, data string)

	// OnClose is called when a WebSocket session is closed.
	OnClose(sessionID string, code int, reason string)

	// OnError is called when an error occurs.
	OnError(sessionID string, errorMsg string)
}

// SetWebSocketCallback registers the native callback handler for WebSocket messages.
// Called from Swift/Kotlin during initialization.
func SetWebSocketCallback(cb WebSocketCallback) {
	irgomobile.SetWebSocketCallback(cb)
}

// WebSocketConnect creates a new WebSocket session and returns its session ID.
// Called from the WebView when a WebSocket connection is requested.
func WebSocketConnect(url string) (string, error) {
	return irgomobile.WebSocketConnect(url)
}

// WebSocketSend sends a message through a virtual WebSocket session.
func WebSocketSend(sessionID string, data string) (string, error) {
	return irgomobile.WebSocketSend(sessionID, data)
}

// WebSocketClose closes a virtual WebSocket session.
func WebSocketClose(sessionID string) error {
	return irgomobile.WebSocketClose(sessionID)
}

// --- Full irgo mobile API re-exports --------------------------------------
// The canonical native shells (scaffolded by `irgo new mobile` into
// ios/Example and android/Example) expect the full bridge surface: streaming,
// lifecycle, and native-capability calls. Re-export everything so a
// CLI-scaffolded example compiles against this AAR unchanged.

// StreamCallback is implemented by Swift/Kotlin to receive a response
// progressively (SSE on mobile). Mirrors irgo/mobile.StreamCallback so
// gomobile exports it to Android/iOS.
type StreamCallback interface {
	// OnResponse delivers the status code and JSON-encoded headers.
	OnResponse(status int, headersJSON string)

	// OnChunk delivers a piece of the response body.
	OnChunk(chunk []byte)

	// OnComplete signals the end of the response. errorMessage is "" on
	// success.
	OnComplete(errorMessage string)
}

// StreamHandle lets native code cancel an in-flight streaming request.
type StreamHandle struct {
	handle *irgomobile.StreamHandle
}

// Cancel aborts the request. Safe to call multiple times and after
// completion.
func (h *StreamHandle) Cancel() {
	if h.handle != nil {
		h.handle.Cancel()
	}
}

// HandleRequestStream processes an HTTP request, streaming the response
// through the callback. Returns immediately; callbacks arrive from a
// background goroutine. Use for Datastar/SSE requests.
func HandleRequestStream(method, url, headers string, body []byte, callback StreamCallback) *StreamHandle {
	return &StreamHandle{handle: irgomobile.HandleRequestStream(method, url, headers, body, callback)}
}

// NativeInvoker is implemented by Swift/Kotlin to receive native-capability
// calls made from Go code. Mirrors irgo/mobile.NativeInvoker so gomobile
// exports it to Android/iOS.
type NativeInvoker interface {
	Invoke(callID, method, paramsJSON string)
}

// SetNativeInvoker registers the platform dispatcher for native-capability
// calls. Called by the native shell during initialization.
func SetNativeInvoker(inv NativeInvoker) {
	irgomobile.SetNativeInvoker(inv)
}

// NativeResult delivers the result of a native-capability call back to Go.
// Called by the native shell when a plugin completes (or fails) a call.
func NativeResult(callID string, ok bool, payloadJSON string) {
	irgomobile.NativeResult(callID, ok, payloadJSON)
}

// SetStateDir sets the directory the bridge uses for persistent state
// (cookie jar etc.). Called by native code during initialization.
func SetStateDir(dir string) {
	irgomobile.SetStateDir(dir)
}

// ClearCookies clears the bridge's cookie jar.
func ClearCookies() {
	irgomobile.ClearCookies()
}

// OnBackground is called by native code when the app enters the background
// (iOS: sceneDidEnterBackground, Android: onPause/onStop).
func OnBackground() {
	irgomobile.OnBackground()
}

// OnForeground is called by native code when the app returns to the
// foreground (iOS: sceneWillEnterForeground, Android: onResume).
func OnForeground() {
	irgomobile.OnForeground()
}
