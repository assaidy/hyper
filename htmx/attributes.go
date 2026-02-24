package hx

// Attr* constants are htmx attribute names.
// Using these constants prevents typos and enables LSP autocompletion.
const (
	AttrGet         = "hx-get"
	AttrPost        = "hx-post"
	AttrPushUrl     = "hx-push-url"
	AttrSelect      = "hx-select"
	AttrSelectOob   = "hx-select-oob"
	AttrSwap        = "hx-swap"
	AttrSwapOob     = "hx-swap-oob"
	AttrTarget      = "hx-target"
	AttrTrigger     = "hx-trigger"
	AttrVals        = "hx-vals"
	AttrBoost       = "hx-boost"
	AttrConfirm     = "hx-confirm"
	AttrDelete      = "hx-delete"
	AttrDisable     = "hx-disable"
	AttrDisabledElt = "hx-disabled-elt"
	AttrDisinherit  = "hx-disinherit"
	AttrEncoding    = "hx-encoding"
	AttrExt         = "hx-ext"
	AttrHeaders     = "hx-headers"
	AttrHistory     = "hx-history"
	AttrHistoryElt  = "hx-history-elt"
	AttrInclude     = "hx-include"
	AttrIndicator   = "hx-indicator"
	AttrInherit     = "hx-inherit"
	AttrParams      = "hx-params"
	AttrPatch       = "hx-patch"
	AttrPreserve    = "hx-preserve"
	AttrPrompt      = "hx-prompt"
	AttrPut         = "hx-put"
	AttrReplaceUrl  = "hx-replace-url"
	AttrRequest     = "hx-request"
	AttrSync        = "hx-sync"
	AttrValidate    = "hx-validate"
	AttrVars        = "hx-vars"
)

// AttrOn generates an hx-on attribute for inline event handlers.
// The event parameter should be the event name (e.g., "click", "mouseover", "htmx:before-swap").
//
// Example:
//
//	AttrOn("click")            // returns "hx-on:click"
//	AttrOn("htmx:before-swap") // returns "hx-on:htmx:before-swap"
func AttrOn(event string) string {
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
