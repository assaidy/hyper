// Package hxws provides constants for the htmx WebSocket extension attributes and events.
//
// NOTE: This package only provides attribute and event name constants.
// To use the ws extension, you must include both htmx and the ws extension
// JavaScript libraries in your HTML and enable it with hx-ext="ws".
// See https://htmx.org/extensions/ws/ for details.
package hxws

const (
	// AttrWsConnect establishes a WebSocket connection to the specified URL.
	AttrWsConnect = "ws-connect"
	// AttrWsSend sends a message to the nearest websocket when the element is triggered.
	AttrWsSend = "ws-send"
)

const (
	// EventWsConnecting triggered when a connection to a WebSocket endpoint is being attempted.
	EventWsConnecting = "htmx:ws-connecting"
	// EventWsOpen triggered when a connection to a WebSocket endpoint has been established.
	EventWsOpen = "htmx:ws-open"
	// EventWsClose triggered when a connection to a WebSocket endpoint has been closed.
	EventWsClose = "htmx:ws-close"
	// EventWsError triggered when an error occurs on the WebSocket.
	EventWsError = "htmx:ws-error"
	// EventWsBeforeMessage triggered when a message has just been received by a socket, before any processing occurs.
	EventWsBeforeMessage = "htmx:ws-before-message"
	// EventWsAfterMessage triggered when a message has been completely processed by htmx and all changes have been settled.
	EventWsAfterMessage = "htmx:ws-after-message"
	// EventWsConfigSend triggered when preparing to send a message from a ws-send element.
	EventWsConfigSend = "htmx:ws-config-send"
	// EventWsBeforeSend triggered just before sending a message to the WebSocket.
	EventWsBeforeSend = "htmx:ws-before-send"
	// EventWsAfterSend triggered just after sending a message to the WebSocket.
	EventWsAfterSend = "htmx:ws-after-send"
)
