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

const (
	// EventAbort send this event to an element to abort a request.
	EventAbort = "htmx:abort"
	// EventAfterOnLoad triggered after an AJAX request has completed processing a successful response.
	EventAfterOnLoad = "htmx:after-on-load"
	// EventAfterProcessNode triggered after htmx has initialized a node.
	EventAfterProcessNode = "htmx:after-process-node"
	// EventAfterRequest triggered after an AJAX request has completed.
	EventAfterRequest = "htmx:after-request"
	// EventAfterSettle triggered after the DOM has settled.
	EventAfterSettle = "htmx:after-settle"
	// EventAfterSwap triggered after new content has been swapped in.
	EventAfterSwap = "htmx:after-swap"
	// EventBeforeCleanupElement triggered before htmx disables an element or removes it from the DOM.
	EventBeforeCleanupElement = "htmx:before-cleanup-element"
	// EventBeforeOnLoad triggered before any response processing occurs.
	EventBeforeOnLoad = "htmx:before-on-load"
	// EventBeforeProcessNode triggered before htmx initializes a node.
	EventBeforeProcessNode = "htmx:before-process-node"
	// EventBeforeRequest triggered before an AJAX request is made.
	EventBeforeRequest = "htmx:before-request"
	// EventBeforeSend triggered just before an ajax request is sent.
	EventBeforeSend = "htmx:before-send"
	// EventBeforeSettle triggered before the DOM settles.
	EventBeforeSettle = "htmx:before-settle"
	// EventBeforeSwap triggered before a swap is done, allows you to configure the swap.
	EventBeforeSwap = "htmx:before-swap"
	// EventBeforeTransition triggered before the View Transition wrapped swap occurs.
	EventBeforeTransition = "htmx:before-transition"
	// EventBeforeHistorySave triggered before content is saved to the history cache.
	EventBeforeHistorySave = "htmx:before-history-save"
	// EventConfigRequest triggered before the request, allows you to customize parameters, headers.
	EventConfigRequest = "htmx:config-request"
	// EventConfirm triggered after a trigger occurs on an element, allows you to cancel (or delay) issuing the AJAX request.
	EventConfirm = "htmx:confirm"
	// EventHistoryCacheError triggered on an error during cache writing.
	EventHistoryCacheError = "htmx:history-cache-error"
	// EventHistoryCacheHit triggered on a cache hit in the history subsystem.
	EventHistoryCacheHit = "htmx:history-cache-hit"
	// EventHistoryCacheMiss triggered on a cache miss in the history subsystem.
	EventHistoryCacheMiss = "htmx:history-cache-miss"
	// EventHistoryCacheMissError triggered on an unsuccessful remote retrieval.
	EventHistoryCacheMissError = "htmx:history-cache-miss-error"
	// EventHistoryCacheMissLoad triggered on a successful remote retrieval.
	EventHistoryCacheMissLoad = "htmx:history-cache-miss-load"
	// EventHistoryRestore triggered when htmx handles a history restoration action.
	EventHistoryRestore = "htmx:history-restore"
	// EventLoad triggered when new content is added to the DOM.
	EventLoad = "htmx:load"
	// EventNoSseSourceError triggered when an element refers to a SSE event in its trigger, but no parent SSE source has been defined.
	EventNoSseSourceError = "htmx:no-sse-source-error"
	// EventOnLoadError triggered when an exception occurs during the onLoad handling in htmx.
	EventOnLoadError = "htmx:on-load-error"
	// EventOobAfterSwap triggered after an out of band element has been swapped in.
	EventOobAfterSwap = "htmx:oob-after-swap"
	// EventOobBeforeSwap triggered before an out of band element swap is done, allows you to configure the swap.
	EventOobBeforeSwap = "htmx:oob-before-swap"
	// EventOobErrorNoTarget triggered when an out of band element does not have a matching ID in the current DOM.
	EventOobErrorNoTarget = "htmx:oob-error-no-target"
	// EventPrompt triggered after a prompt is shown.
	EventPrompt = "htmx:prompt"
	// EventPushedIntoHistory triggered after a url is pushed into history.
	EventPushedIntoHistory = "htmx:pushed-into-history"
	// EventReplacedInHistory triggered after a url is replaced in history.
	EventReplacedInHistory = "htmx:replaced-in-history"
	// EventResponseError triggered when an HTTP response error (non-200 or 300 response code) occurs.
	EventResponseError = "htmx:response-error"
	// EventSendAbort triggered when a request is aborted.
	EventSendAbort = "htmx:send-abort"
	// EventSendError triggered when a network error prevents an HTTP request from happening.
	EventSendError = "htmx:send-error"
	// EventSseError triggered when an error occurs with a SSE source.
	EventSseError = "htmx:sse-error"
	// EventSseOpen triggered when a SSE source is opened.
	EventSseOpen = "htmx:sse-open"
	// EventSseMessage triggered when a message is received from a SSE source.
	EventSseMessage = "htmx:sse-message"
	// EventSwapError triggered when an error occurs during the swap phase.
	EventSwapError = "htmx:swap-error"
	// EventTargetError triggered when an invalid target is specified.
	EventTargetError = "htmx:target-error"
	// EventTimeout triggered when a request timeout occurs.
	EventTimeout = "htmx:timeout"
	// EventValidationValidate triggered before an element is validated.
	EventValidationValidate = "htmx:validation:validate"
	// EventValidationFailed triggered when an element fails validation.
	EventValidationFailed = "htmx:validation:failed"
	// EventValidationHalted triggered when a request is halted due to validation errors.
	EventValidationHalted = "htmx:validation:halted"
	// EventXhrAbort triggered when an ajax request aborts.
	EventXhrAbort = "htmx:xhr:abort"
	// EventXhrLoadend triggered when an ajax request ends.
	EventXhrLoadend = "htmx:xhr:loadend"
	// EventXhrLoadstart triggered when an ajax request starts.
	EventXhrLoadstart = "htmx:xhr:loadstart"
	// EventXhrProgress triggered periodically during an ajax request that supports progress events.
	EventXhrProgress = "htmx:xhr:progress"
)
