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
	// EventHtmxWsConnecting triggered when a connection to a WebSocket endpoint is being attempted.
	EventHtmxWsConnecting = "htmx:ws-connecting"
	// EventHtmxWsOpen triggered when a connection to a WebSocket endpoint has been established.
	EventHtmxWsOpen = "htmx:ws-open"
	// EventHtmxWsClose triggered when a connection to a WebSocket endpoint has been closed.
	EventHtmxWsClose = "htmx:ws-close"
	// EventHtmxWsError triggered when an error occurs on the WebSocket.
	EventHtmxWsError = "htmx:ws-error"
	// EventHtmxWsBeforeMessage triggered when a message has just been received by a socket, before any processing occurs.
	EventHtmxWsBeforeMessage = "htmx:ws-before-message"
	// EventHtmxWsAfterMessage triggered when a message has been completely processed by htmx and all changes have been settled.
	EventHtmxWsAfterMessage = "htmx:ws-after-message"
	// EventHtmxWsConfigSend triggered when preparing to send a message from a ws-send element.
	EventHtmxWsConfigSend = "htmx:ws-config-send"
	// EventHtmxWsBeforeSend triggered just before sending a message to the WebSocket.
	EventHtmxWsBeforeSend = "htmx:ws-before-send"
	// EventHtmxWsAfterSend triggered just after sending a message to the WebSocket.
	EventHtmxWsAfterSend = "htmx:ws-after-send"
)
