// Package hxws provides constants for the htmx WebSocket extension attributes and events.
//
// NOTE: This package only provides attribute and event name constants.
// To use the ws extension, you must include both htmx and the ws extension
// JavaScript libraries in your HTML and enable it with hx-ext="ws".
// See https://htmx.org/extensions/ws/ for details.
package hxws

// Attr* constants are htmx ws extension attribute names.
const (
	AttrWsConnect = "ws-connect"
	AttrWsSend    = "ws-send"
)

// Event* constants are htmx ws extension event names.
// Use these constants with hx.AttrOn for type-safe inline event handlers.
const (
	EventWsConnecting    = "htmx:ws-connecting"
	EventWsOpen          = "htmx:ws-open"
	EventWsClose         = "htmx:ws-close"
	EventWsError         = "htmx:ws-error"
	EventWsBeforeMessage = "htmx:ws-before-message"
	EventWsAfterMessage  = "htmx:ws-after-message"
	EventWsConfigSend    = "htmx:ws-config-send"
	EventWsBeforeSend    = "htmx:ws-before-send"
	EventWsAfterSend     = "htmx:ws-after-send"
)
