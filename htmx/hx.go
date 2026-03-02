// Package hx provides constants for htmx attributes and events.
//
// NOTE: This package only provides attribute and event name constants.
// To use htmx in your application, you must include the htmx JavaScript library
// in your HTML. See https://htmx.org for installation instructions.
package hx

const (
	// AttrHxGet issues a GET to the specified URL.
	AttrHxGet = "hx-get"
	// AttrHxPost issues a POST to the specified URL.
	AttrHxPost = "hx-post"
	// AttrHxPushUrl pushes a URL into the browser location bar to create history.
	AttrHxPushUrl = "hx-push-url"
	// AttrHxSelect selects content to swap in from a response.
	AttrHxSelect = "hx-select"
	// AttrHxSelectOob selects content to swap in from a response, somewhere other than the target (out of band).
	AttrHxSelectOob = "hx-select-oob"
	// AttrHxSwap controls how content will swap in (outerHTML, beforeend, afterend, ...).
	AttrHxSwap = "hx-swap"
	// AttrHxSwapOob marks element to swap in from a response (out of band).
	AttrHxSwapOob = "hx-swap-oob"
	// AttrHxTarget specifies the target element to be swapped.
	AttrHxTarget = "hx-target"
	// AttrHxTrigger specifies the event that triggers the request.
	AttrHxTrigger = "hx-trigger"
	// AttrHxVals adds values to submit with the request (JSON format).
	AttrHxVals = "hx-vals"
	// AttrHxBoost adds progressive enhancement for links and forms.
	AttrHxBoost = "hx-boost"
	// AttrHxConfirm shows a confirm() dialog before issuing a request.
	AttrHxConfirm = "hx-confirm"
	// AttrHxDelete issues a DELETE to the specified URL.
	AttrHxDelete = "hx-delete"
	// AttrHxDisable disables htmx processing for the given node and any children nodes.
	AttrHxDisable = "hx-disable"
	// AttrHxDisabledElt adds the disabled attribute to the specified elements while a request is in flight.
	AttrHxDisabledElt = "hx-disabled-elt"
	// AttrHxDisinherit control and disable automatic attribute inheritance for child nodes.
	AttrHxDisinherit = "hx-disinherit"
	// AttrHxEncoding changes the request encoding type.
	AttrHxEncoding = "hx-encoding"
	// AttrHxExt extensions to use for this element.
	AttrHxExt = "hx-ext"
	// AttrHxHeaders adds to the headers that will be submitted with the request.
	AttrHxHeaders = "hx-headers"
	// AttrHxHistory prevent sensitive data being saved to the history cache.
	AttrHxHistory = "hx-history"
	// AttrHxHistoryElt the element to snapshot and restore during history navigation.
	AttrHxHistoryElt = "hx-history-elt"
	// AttrHxInclude include additional data in requests.
	AttrHxInclude = "hx-include"
	// AttrHxIndicator the element to put the htmx-request class on during the request.
	AttrHxIndicator = "hx-indicator"
	// AttrHxInherit control and enable automatic attribute inheritance for child nodes if it has been disabled by default.
	AttrHxInherit = "hx-inherit"
	// AttrHxParams filters the parameters that will be submitted with a request.
	AttrHxParams = "hx-params"
	// AttrHxPatch issues a PATCH to the specified URL.
	AttrHxPatch = "hx-patch"
	// AttrHxPreserve specifies elements to keep unchanged between requests.
	AttrHxPreserve = "hx-preserve"
	// AttrHxPrompt shows a prompt() before submitting a request.
	AttrHxPrompt = "hx-prompt"
	// AttrHxPut issues a PUT to the specified URL.
	AttrHxPut = "hx-put"
	// AttrHxReplaceUrl replace the URL in the browser location bar.
	AttrHxReplaceUrl = "hx-replace-url"
	// AttrHxRequest configures various aspects of the request.
	AttrHxRequest = "hx-request"
	// AttrHxSync control how requests made by different elements are synchronized.
	AttrHxSync = "hx-sync"
	// AttrHxValidate force elements to validate themselves before a request.
	AttrHxValidate = "hx-validate"
	// AttrHxVars adds values dynamically to the parameters to submit with the request (deprecated, use hx-vals).
	AttrHxVars = "hx-vars"
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
	// SwapInnerHtml swaps the inner HTML of the target element.
	SwapInnerHtml = "innerHTML"
	// SwapOuterHtml replaces the entire target element with the returned content.
	SwapOuterHtml = "outerHTML"
	// SwapBeforeBegin inserts content before the target element.
	SwapBeforeBegin = "beforebegin"
	// SwapAfterBegin inserts content at the beginning of the target element.
	SwapAfterBegin = "afterbegin"
	// SwapBeforeEnd inserts content at the end of the target element.
	SwapBeforeEnd = "beforeend"
	// SwapAfterEnd inserts content after the target element.
	SwapAfterEnd = "afterend"
	// SwapDelete deletes the target element.
	SwapDelete = "delete"
	// SwapNone does not swap content.
	SwapNone = "none"
)

// SwapOob generates an hx-swap-oob attribute value with a specific target.
// The swap parameter is the swap strategy (e.g., SwapBeforeEnd, SwapInnerHtml).
// The target parameter is a CSS selector (e.g., "#myId", "#table tbody").
//
// Example:
//
//	SwapOob(SwapBeforeEnd, "#notifications") // returns "beforeend:#notifications"
//	SwapOob(SwapBeforeEnd, "#table tbody")   // returns "beforeend:#table tbody"
func SwapOob(swap, target string) string {
	return swap + ":" + target
}

const (
	// EventHtmxAbort send this event to an element to abort a request.
	EventHtmxAbort = "htmx:abort"
	// EventHtmxAfterOnLoad triggered after an AJAX request has completed processing a successful response.
	EventHtmxAfterOnLoad = "htmx:after-on-load"
	// EventHtmxAfterProcessNode triggered after htmx has initialized a node.
	EventHtmxAfterProcessNode = "htmx:after-process-node"
	// EventHtmxAfterRequest triggered after an AJAX request has completed.
	EventHtmxAfterRequest = "htmx:after-request"
	// EventHtmxAfterSettle triggered after the DOM has settled.
	EventHtmxAfterSettle = "htmx:after-settle"
	// EventHtmxAfterSwap triggered after new content has been swapped in.
	EventHtmxAfterSwap = "htmx:after-swap"
	// EventHtmxBeforeCleanupElement triggered before htmx disables an element or removes it from the DOM.
	EventHtmxBeforeCleanupElement = "htmx:before-cleanup-element"
	// EventHtmxBeforeOnLoad triggered before any response processing occurs.
	EventHtmxBeforeOnLoad = "htmx:before-on-load"
	// EventHtmxBeforeProcessNode triggered before htmx initializes a node.
	EventHtmxBeforeProcessNode = "htmx:before-process-node"
	// EventHtmxBeforeRequest triggered before an AJAX request is made.
	EventHtmxBeforeRequest = "htmx:before-request"
	// EventHtmxBeforeSend triggered just before an ajax request is sent.
	EventHtmxBeforeSend = "htmx:before-send"
	// EventHtmxBeforeSettle triggered before the DOM settles.
	EventHtmxBeforeSettle = "htmx:before-settle"
	// EventHtmxBeforeSwap triggered before a swap is done, allows you to configure the swap.
	EventHtmxBeforeSwap = "htmx:before-swap"
	// EventHtmxBeforeTransition triggered before the View Transition wrapped swap occurs.
	EventHtmxBeforeTransition = "htmx:before-transition"
	// EventHtmxBeforeHistorySave triggered before content is saved to the history cache.
	EventHtmxBeforeHistorySave = "htmx:before-history-save"
	// EventHtmxConfigRequest triggered before the request, allows you to customize parameters, headers.
	EventHtmxConfigRequest = "htmx:config-request"
	// EventHtmxConfirm triggered after a trigger occurs on an element, allows you to cancel (or delay) issuing the AJAX request.
	EventHtmxConfirm = "htmx:confirm"
	// EventHtmxHistoryCacheError triggered on an error during cache writing.
	EventHtmxHistoryCacheError = "htmx:history-cache-error"
	// EventHtmxHistoryCacheHit triggered on a cache hit in the history subsystem.
	EventHtmxHistoryCacheHit = "htmx:history-cache-hit"
	// EventHtmxHistoryCacheMiss triggered on a cache miss in the history subsystem.
	EventHtmxHistoryCacheMiss = "htmx:history-cache-miss"
	// EventHtmxHistoryCacheMissError triggered on an unsuccessful remote retrieval.
	EventHtmxHistoryCacheMissError = "htmx:history-cache-miss-error"
	// EventHtmxHistoryCacheMissLoad triggered on a successful remote retrieval.
	EventHtmxHistoryCacheMissLoad = "htmx:history-cache-miss-load"
	// EventHtmxHistoryRestore triggered when htmx handles a history restoration action.
	EventHtmxHistoryRestore = "htmx:history-restore"
	// EventHtmxLoad triggered when new content is added to the DOM.
	EventHtmxLoad = "htmx:load"
	// EventHtmxNoSseSourceError triggered when an element refers to a SSE event in its trigger, but no parent SSE source has been defined.
	EventHtmxNoSseSourceError = "htmx:no-sse-source-error"
	// EventHtmxOnLoadError triggered when an exception occurs during the onLoad handling in htmx.
	EventHtmxOnLoadError = "htmx:on-load-error"
	// EventHtmxOobAfterSwap triggered after an out of band element has been swapped in.
	EventHtmxOobAfterSwap = "htmx:oob-after-swap"
	// EventHtmxOobBeforeSwap triggered before an out of band element swap is done, allows you to configure the swap.
	EventHtmxOobBeforeSwap = "htmx:oob-before-swap"
	// EventHtmxOobErrorNoTarget triggered when an out of band element does not have a matching ID in the current DOM.
	EventHtmxOobErrorNoTarget = "htmx:oob-error-no-target"
	// EventHtmxPrompt triggered after a prompt is shown.
	EventHtmxPrompt = "htmx:prompt"
	// EventHtmxPushedIntoHistory triggered after a url is pushed into history.
	EventHtmxPushedIntoHistory = "htmx:pushed-into-history"
	// EventHtmxReplacedInHistory triggered after a url is replaced in history.
	EventHtmxReplacedInHistory = "htmx:replaced-in-history"
	// EventHtmxResponseError triggered when an HTTP response error (non-200 or 300 response code) occurs.
	EventHtmxResponseError = "htmx:response-error"
	// EventHtmxSendAbort triggered when a request is aborted.
	EventHtmxSendAbort = "htmx:send-abort"
	// EventHtmxSendError triggered when a network error prevents an HTTP request from happening.
	EventHtmxSendError = "htmx:send-error"
	// EventHtmxSseError triggered when an error occurs with a SSE source.
	EventHtmxSseError = "htmx:sse-error"
	// EventHtmxSseOpen triggered when a SSE source is opened.
	EventHtmxSseOpen = "htmx:sse-open"
	// EventHtmxSseMessage triggered when a message is received from a SSE source.
	EventHtmxSseMessage = "htmx:sse-message"
	// EventHtmxSwapError triggered when an error occurs during the swap phase.
	EventHtmxSwapError = "htmx:swap-error"
	// EventHtmxTargetError triggered when an invalid target is specified.
	EventHtmxTargetError = "htmx:target-error"
	// EventHtmxTimeout triggered when a request timeout occurs.
	EventHtmxTimeout = "htmx:timeout"
	// EventHtmxValidationValidate triggered before an element is validated.
	EventHtmxValidationValidate = "htmx:validation:validate"
	// EventHtmxValidationFailed triggered when an element fails validation.
	EventHtmxValidationFailed = "htmx:validation:failed"
	// EventHtmxValidationHalted triggered when a request is halted due to validation errors.
	EventHtmxValidationHalted = "htmx:validation:halted"
	// EventHtmxXhrAbort triggered when an ajax request aborts.
	EventHtmxXhrAbort = "htmx:xhr:abort"
	// EventHtmxXhrLoadend triggered when an ajax request ends.
	EventHtmxXhrLoadend = "htmx:xhr:loadend"
	// EventHtmxXhrLoadstart triggered when an ajax request starts.
	EventHtmxXhrLoadstart = "htmx:xhr:loadstart"
	// EventHtmxXhrProgress triggered periodically during an ajax request that supports progress events.
	EventHtmxXhrProgress = "htmx:xhr:progress"
)
