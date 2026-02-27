// Package hx provides constants for htmx attributes and events.
//
// NOTE: This package only provides attribute and event name constants.
// To use htmx in your application, you must include the htmx JavaScript library
// in your HTML. See https://htmx.org for installation instructions.
package hx

// Attr* constants are htmx attribute names.
// Using these constants prevents typos and enables LSP autocompletion.
const (
	AttrHxGet         = "hx-get"
	AttrHxPost        = "hx-post"
	AttrHxPushUrl     = "hx-push-url"
	AttrHxSelect      = "hx-select"
	AttrHxSelectOob   = "hx-select-oob"
	AttrHxSwap        = "hx-swap"
	AttrHxSwapOob     = "hx-swap-oob"
	AttrHxTarget      = "hx-target"
	AttrHxTrigger     = "hx-trigger"
	AttrHxVals        = "hx-vals"
	AttrHxBoost       = "hx-boost"
	AttrHxConfirm     = "hx-confirm"
	AttrHxDelete      = "hx-delete"
	AttrHxDisable     = "hx-disable"
	AttrHxDisabledElt = "hx-disabled-elt"
	AttrHxDisinherit  = "hx-disinherit"
	AttrHxEncoding    = "hx-encoding"
	AttrHxExt         = "hx-ext"
	AttrHxHeaders     = "hx-headers"
	AttrHxHistory     = "hx-history"
	AttrHxHistoryElt  = "hx-history-elt"
	AttrHxInclude     = "hx-include"
	AttrHxIndicator   = "hx-indicator"
	AttrHxInherit     = "hx-inherit"
	AttrHxParams      = "hx-params"
	AttrHxPatch       = "hx-patch"
	AttrHxPreserve    = "hx-preserve"
	AttrHxPrompt      = "hx-prompt"
	AttrHxPut         = "hx-put"
	AttrHxReplaceUrl  = "hx-replace-url"
	AttrHxRequest     = "hx-request"
	AttrHxSync        = "hx-sync"
	AttrHxValidate    = "hx-validate"
	AttrHxVars        = "hx-vars"
)

// AttrHxOn generates an hx-on attribute for inline event handlers.
// The event parameter should be the event name (e.g., "click", "mouseover", "htmx:before-swap").
//
// Example:
//
//	AttrHxOn("click")            // returns "hx-on:click"
//	AttrHxOn("htmx:before-swap") // returns "hx-on:htmx:before-swap"
func AttrHxOn(event string) string {
	return "hx-on:" + event
}

// Swap* constants are valid values for the hx-swap attribute.
const (
	SwapInnerHtml   = "innerHTML"
	SwapOuterHtml   = "outerHTML"
	SwapBeforeBegin = "beforebegin"
	SwapAfterBegin  = "afterbegin"
	SwapBeforeEnd   = "beforeend"
	SwapAfterEnd    = "afterend"
	SwapDelete      = "delete"
	SwapNone        = "none"
)

// Event* constants are htmx event names.
// Use these constants with AttrOn for type-safe inline event handlers.
const (
	EventAbort                 = "htmx:abort"
	EventAfterOnLoad           = "htmx:after-on-load"
	EventAfterProcessNode      = "htmx:after-process-node"
	EventAfterRequest          = "htmx:after-request"
	EventAfterSettle           = "htmx:after-settle"
	EventAfterSwap             = "htmx:after-swap"
	EventBeforeCleanupElement  = "htmx:before-cleanup-element"
	EventBeforeOnLoad          = "htmx:before-on-load"
	EventBeforeProcessNode     = "htmx:before-process-node"
	EventBeforeRequest         = "htmx:before-request"
	EventBeforeSend            = "htmx:before-send"
	EventBeforeSettle          = "htmx:before-settle"
	EventBeforeSwap            = "htmx:before-swap"
	EventConfigRequest         = "htmx:config-request"
	EventConfirm               = "htmx:confirm"
	EventHistoryCacheError     = "htmx:history-cache-error"
	EventHistoryCacheMiss      = "htmx:history-cache-miss"
	EventHistoryCacheMissError = "htmx:history-cache-miss-error"
	EventHistoryRestore        = "htmx:history-restore"
	EventLoad                  = "htmx:load"
	EventNoSseSourceError      = "htmx:no-sse-source-error"
	EventOnLoadError           = "htmx:on-load-error"
	EventOobAfterSwap          = "htmx:oob-after-swap"
	EventOobBeforeSwap         = "htmx:oob-before-swap"
	EventOobErrorNoTarget      = "htmx:oob-error-no-target"
	EventPrompt                = "htmx:prompt"
	EventResponseError         = "htmx:response-error"
	EventSendError             = "htmx:send-error"
	EventSseError              = "htmx:sse-error"
	EventSseOpen               = "htmx:sse-open"
	EventSseMessage            = "htmx:sse-message"
	EventSwapError             = "htmx:swap-error"
	EventTargetError           = "htmx:target-error"
	EventTimeout               = "htmx:timeout"
	EventValidationValidate    = "htmx:validation:validate"
	EventValidationFailed      = "htmx:validation:failed"
	EventValidationHalted      = "htmx:validation:halted"
	EventXhrAbort              = "htmx:xhr:abort"
	EventXhrLoadend            = "htmx:xhr:loadend"
	EventXhrLoadstart          = "htmx:xhr:loadstart"
	EventXhrProgress           = "htmx:xhr:progress"
)
