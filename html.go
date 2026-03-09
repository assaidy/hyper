package h

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strings"
	"sync"
)

// Text represents a plain text node that renders HTML-escaped content.
// Unlike HTML elements, Text nodes are not wrapped in tags and are rendered
// as literal text content with HTML entities automatically escaped.
type Text string

func (me Text) Render(w io.Writer) error {
	_, err := io.WriteString(w, html.EscapeString(string(me)))
	return err
}

// RawText represents a text node that renders its content exactly as provided,
// without any HTML escaping.
type RawText string

func (me RawText) Render(w io.Writer) error {
	_, err := io.WriteString(w, string(me))
	return err
}

// KV represents a key-value map for HTML attributes.
//
// The value type must be either string or bool:
//   - string: Attribute will have the format key="value" (HTML-escaped)
//   - bool: If true, attribute appears as key (valueless). If false, attribute is omitted.
//   - any other type triggers an error during rendering.
//
// Example:
//
//	KV{"class": "container", "hidden": true, "disabled": false}
//	// Renders: class="container" hidden
type KV map[string]any

type attribute struct {
	key   string
	value any
}

const (
	// AttrAccept sets the accepted file types for <input type="file">.
	AttrAccept = "accept"
	// AttrAcceptCharset sets the character encodings accepted by the server.
	AttrAcceptCharset = "accept-charset"
	// AttrAccessKey gives keyboard shortcut access to an element.
	AttrAccessKey = "accesskey"
	// AttrAction specifies where to send the form data.
	AttrAction = "action"
	// AttrAlign specifies the alignment of an element.
	AttrAlign = "align"
	// AttrAllow specifies permissions for an iframe.
	AttrAllow = "allow"
	// AttrAlpha sets the alpha transparency level of an element.
	AttrAlpha = "alpha"
	// AttrAlt provides alternative text for an image.
	AttrAlt = "alt"
	// AttrAs specifies the relation between the linked resource and the document.
	AttrAs = "as"
	// AttrAsync indicates that the script should execute asynchronously.
	AttrAsync = "async"
	// AttrAutocapitalize controls whether text input is automatically capitalized.
	AttrAutocapitalize = "autocapitalize"
	// AttrAutocomplete specifies whether an input field should have autocomplete enabled.
	AttrAutocomplete = "autocomplete"
	// AttrAutofocus specifies that an element should automatically get focus on page load.
	AttrAutofocus = "autofocus"
	// AttrAutoplay specifies that the audio/video should automatically start playing.
	AttrAutoplay = "autoplay"
	// AttrBackground specifies the background image URL.
	AttrBackground = "background"
	// AttrBGColor specifies the background color of an element.
	AttrBGColor = "bgcolor"
	// AttrBorder specifies the border width around an element.
	AttrBorder = "border"
	// AttrCapture specifies which camera/mic to use for media capture.
	AttrCapture = "capture"
	// AttrCharset specifies the character encoding of the document.
	AttrCharset = "charset"
	// AttrChecked specifies whether an input checkbox or radio is checked.
	AttrChecked = "checked"
	// AttrCite specifies the source of a quotation.
	AttrCite = "cite"
	// AttrClass specifies one or more class names for an element.
	AttrClass = "class"
	// AttrColor specifies the text color of an element.
	AttrColor = "color"
	// AttrColorSpace specifies the color space for an image.
	AttrColorSpace = "colorspace"
	// AttrCols specifies the number of columns in a textarea.
	AttrCols = "cols"
	// AttrColSpan specifies the number of columns a table cell should span.
	AttrColSpan = "colspan"
	// AttrContent provides metadata about the element.
	AttrContent = "content"
	// AttrContentEditable specifies whether the element is editable.
	AttrContentEditable = "contenteditable"
	// AttrControls shows the audio/video controls.
	AttrControls = "controls"
	// AttrCoords specifies the coordinates of an area in an image map.
	AttrCoords = "coords"
	// AttrCrossOrigin specifies how the element handles cross-origin requests.
	AttrCrossOrigin = "crossorigin"
	// AttrCsp specifies the Content Security Policy for an element.
	AttrCsp = "csp"
	// AttrData specifies the URL of the data for an object element.
	AttrData = "data"
	// AttrDateTime specifies the date and time for an element.
	AttrDateTime = "datetime"
	// AttrDecoding specifies how to decode an image.
	AttrDecoding = "decoding"
	// AttrDefault specifies that a track should be enabled by default.
	AttrDefault = "default"
	// AttrDefer indicates that the script should be executed after the document is parsed.
	AttrDefer = "defer"
	// AttrDir specifies the text direction of an element.
	AttrDir = "dir"
	// AttrDirName specifies the name of the form field used for sending the directionality of the element.
	AttrDirName = "dirname"
	// AttrDisabled specifies that an element should be disabled.
	AttrDisabled = "disabled"
	// AttrDownload specifies that the target should be downloaded when clicked.
	AttrDownload = "download"
	// AttrDraggable specifies whether an element is draggable.
	AttrDraggable = "draggable"
	// AttrEncType specifies how form data should be encoded before sending to a server.
	AttrEncType = "enctype"
	// AttrEnterKeyHint specifies what action label to show on the enter key.
	AttrEnterKeyHint = "enterkeyhint"
	// AttrElementTiming specifies that an element should be observed for performance.
	AttrElementTiming = "elementtiming"
	// AttrForm specifies the id of a form element that the element belongs to.
	AttrForm = "form"
	// AttrFormAction specifies where to send the form data.
	AttrFormAction = "formaction"
	// AttrFormEncType specifies how form data should be encoded.
	AttrFormEncType = "formenctype"
	// AttrFormMethod specifies the HTTP method for form submission.
	AttrFormMethod = "formmethod"
	// AttrFormNoValidate specifies that the form should not be validated.
	AttrFormNoValidate = "formnovalidate"
	// AttrFormTarget specifies where to display the response after form submission.
	AttrFormTarget = "formtarget"
	// AttrFetchPriority indicates the priority of fetching an external resource.
	AttrFetchPriority = "fetchpriority"
	// AttrHeaders specifies the header cells that a table cell relates to.
	AttrHeaders = "headers"
	// AttrHeight specifies the height of an element.
	AttrHeight = "height"
	// AttrHidden specifies that an element is not yet or is no longer relevant.
	AttrHidden = "hidden"
	// AttrHigh specifies the lower bound of a range.
	AttrHigh = "high"
	// AttrHref specifies the URL of a link.
	AttrHref = "href"
	// AttrHrefLang specifies the language of the linked resource.
	AttrHrefLang = "hreflang"
	// AttrHttpEquiv provides an HTTP header for the information in the content attribute.
	AttrHttpEquiv = "http-equiv"
	// AttrID specifies a unique id for an element.
	AttrID = "id"
	// AttrIntegrity specifies a hash of the resource to verify its integrity.
	AttrIntegrity = "integrity"
	// AttrInputMode provides a hint to browsers about the type of data the user should enter.
	AttrInputMode = "inputmode"
	// AttrIsMap specifies that an image is part of a server-side image map.
	AttrIsMap = "ismap"
	// AttrItemProp specifies the property of an item.
	AttrItemProp = "itemprop"
	// AttrKind specifies the kind of text track.
	AttrKind = "kind"
	// AttrLabel specifies the label of an option or track.
	AttrLabel = "label"
	// AttrLang specifies the language of the element.
	AttrLang = "lang"
	// AttrLanguage specifies the scripting language of an element.
	AttrLanguage = "language"
	// AttrLoading specifies whether to load an image lazily.
	AttrLoading = "loading"
	// AttrList refers to a datalist containing predefined options.
	AttrList = "list"
	// AttrLoop specifies whether to loop an audio/video.
	AttrLoop = "loop"
	// AttrLow specifies the upper bound of a range.
	AttrLow = "low"
	// AttrMax specifies the maximum value.
	AttrMax = "max"
	// AttrMaxLength specifies the maximum number of characters allowed.
	AttrMaxLength = "maxlength"
	// AttrMinLength specifies the minimum number of characters required.
	AttrMinLength = "minlength"
	// AttrMedia specifies the media type or device the resource applies to.
	AttrMedia = "media"
	// AttrMethod specifies the HTTP method for form submission.
	AttrMethod = "method"
	// AttrMin specifies the minimum value.
	AttrMin = "min"
	// AttrMultiple specifies that a user can enter more than one value.
	AttrMultiple = "multiple"
	// AttrMuted specifies that the audio should be muted.
	AttrMuted = "muted"
	// AttrName specifies the name of an element.
	AttrName = "name"
	// AttrNoValidate specifies that the form should not be validated.
	AttrNoValidate = "novalidate"
	// AttrOnAbort specifies the event handler for the abort event.
	AttrOnAbort = "onAbort"
	// AttrOnActivate specifies the event handler for the activate event.
	AttrOnActivate = "onActivate"
	// AttrOnAfterPrint specifies the event handler for the afterprint event.
	AttrOnAfterPrint = "onAfterPrint"
	// AttrOnAfterUpdate specifies the event handler for the afterupdate event.
	AttrOnAfterUpdate = "onAfterUpdate"
	// AttrOnBeforeActivate specifies the event handler for the beforeactivate event.
	AttrOnBeforeActivate = "onBeforeActivate"
	// AttrOnBeforeCopy specifies the event handler for the beforecopy event.
	AttrOnBeforeCopy = "onBeforeCopy"
	// AttrOnBeforeCut specifies the event handler for the beforecut event.
	AttrOnBeforeCut = "onBeforeCut"
	// AttrOnBeforeDeactivate specifies the event handler for the beforedeactivate event.
	AttrOnBeforeDeactivate = "onBeforeDeactivate"
	// AttrOnBeforeEditFocus specifies the event handler for the beforeeditfocus event.
	AttrOnBeforeEditFocus = "onBeforeEditFocus"
	// AttrOnBeforePaste specifies the event handler for the beforepaste event.
	AttrOnBeforePaste = "onBeforePaste"
	// AttrOnBeforePrint specifies the event handler for the beforeprint event.
	AttrOnBeforePrint = "onBeforePrint"
	// AttrOnBeforeUnload specifies the event handler for the beforeunload event.
	AttrOnBeforeUnload = "onBeforeUnload"
	// AttrOnBeforeUpdate specifies the event handler for the beforeupdate event.
	AttrOnBeforeUpdate = "onBeforeUpdate"
	// AttrOnBegin specifies the event handler for the begin event.
	AttrOnBegin = "onBegin"
	// AttrOnBlur specifies the event handler for the blur event.
	AttrOnBlur = "onBlur"
	// AttrOnBounce specifies the event handler for the bounce event.
	AttrOnBounce = "onBounce"
	// AttrOnCellChange specifies the event handler for the cellchange event.
	AttrOnCellChange = "onCellChange"
	// AttrOnChange specifies the event handler for the change event.
	AttrOnChange = "onChange"
	// AttrOnClick specifies the event handler for the click event.
	AttrOnClick = "onClick"
	// AttrOnContextMenu specifies the event handler for the contextmenu event.
	AttrOnContextMenu = "onContextMenu"
	// AttrOnControlSelect specifies the event handler for the controlselect event.
	AttrOnControlSelect = "onControlSelect"
	// AttrOnCopy specifies the event handler for the copy event.
	AttrOnCopy = "onCopy"
	// AttrOnCut specifies the event handler for the cut event.
	AttrOnCut = "onCut"
	// AttrOnDataAvailable specifies the event handler for the dataavailable event.
	AttrOnDataAvailable = "onDataAvailable"
	// AttrOnDataSetChanged specifies the event handler for the datasetchanged event.
	AttrOnDataSetChanged = "onDataSetChanged"
	// AttrOnDataSetComplete specifies the event handler for the datasetcomplete event.
	AttrOnDataSetComplete = "onDataSetComplete"
	// AttrOnDblClick specifies the event handler for the dblclick event.
	AttrOnDblClick = "onDblClick"
	// AttrOnDeactivate specifies the event handler for the deactivate event.
	AttrOnDeactivate = "onDeactivate"
	// AttrOnDrag specifies the event handler for the drag event.
	AttrOnDrag = "onDrag"
	// AttrOnDragEnd specifies the event handler for the dragend event.
	AttrOnDragEnd = "onDragEnd"
	// AttrOnDragLeave specifies the event handler for the dragleave event.
	AttrOnDragLeave = "onDragLeave"
	// AttrOnDragEnter specifies the event handler for the dragenter event.
	AttrOnDragEnter = "onDragEnter"
	// AttrOnDragOver specifies the event handler for the dragover event.
	AttrOnDragOver = "onDragOver"
	// AttrOnDragDrop specifies the event handler for the dragdrop event.
	AttrOnDragDrop = "onDragDrop"
	// AttrOnDragStart specifies the event handler for the dragstart event.
	AttrOnDragStart = "onDragStart"
	// AttrOnDrop specifies the event handler for the drop event.
	AttrOnDrop = "onDrop"
	// AttrOnEnd specifies the event handler for the end event.
	AttrOnEnd = "onEnd"
	// AttrOnError specifies the event handler for the error event.
	AttrOnError = "onError"
	// AttrOnErrorUpdate specifies the event handler for the errorupdate event.
	AttrOnErrorUpdate = "onErrorUpdate"
	// AttrOnFilterChange specifies the event handler for the filterchange event.
	AttrOnFilterChange = "onFilterChange"
	// AttrOnFinish specifies the event handler for the finish event.
	AttrOnFinish = "onFinish"
	// AttrOnFocus specifies the event handler for the focus event.
	AttrOnFocus = "onFocus"
	// AttrOnFocusIn specifies the event handler for the focusin event.
	AttrOnFocusIn = "onFocusIn"
	// AttrOnFocusOut specifies the event handler for the focusout event.
	AttrOnFocusOut = "onFocusOut"
	// AttrOnHashChange specifies the event handler for the hashchange event.
	AttrOnHashChange = "onHashChange"
	// AttrOnHelp specifies the event handler for the help event.
	AttrOnHelp = "onHelp"
	// AttrOnInput specifies the event handler for the input event.
	AttrOnInput = "onInput"
	// AttrOnKeyDown specifies the event handler for the keydown event.
	AttrOnKeyDown = "onKeyDown"
	// AttrOnKeyPress specifies the event handler for the keypress event.
	AttrOnKeyPress = "onKeyPress"
	// AttrOnKeyUp specifies the event handler for the keyup event.
	AttrOnKeyUp = "onKeyUp"
	// AttrOnLayoutComplete specifies the event handler for the layoutcomplete event.
	AttrOnLayoutComplete = "onLayoutComplete"
	// AttrOnLoad specifies the event handler for the load event.
	AttrOnLoad = "onLoad"
	// AttrOnLoseCapture specifies the event handler for the losecapture event.
	AttrOnLoseCapture = "onLoseCapture"
	// AttrOnMediaComplete specifies the event handler for the mediacomplete event.
	AttrOnMediaComplete = "onMediaComplete"
	// AttrOnMediaError specifies the event handler for the mediaerror event.
	AttrOnMediaError = "onMediaError"
	// AttrOnMessage specifies the event handler for the message event.
	AttrOnMessage = "onMessage"
	// AttrOnMouseDown specifies the event handler for the mousedown event.
	AttrOnMouseDown = "onMouseDown"
	// AttrOnMouseEnter specifies the event handler for the mouseenter event.
	AttrOnMouseEnter = "onMouseEnter"
	// AttrOnMouseLeave specifies the event handler for the mouseleave event.
	AttrOnMouseLeave = "onMouseLeave"
	// AttrOnMouseMove specifies the event handler for the mousemove event.
	AttrOnMouseMove = "onMouseMove"
	// AttrOnMouseOut specifies the event handler for the mouseout event.
	AttrOnMouseOut = "onMouseOut"
	// AttrOnMouseOver specifies the event handler for the mouseover event.
	AttrOnMouseOver = "onMouseOver"
	// AttrOnMouseUp specifies the event handler for the mouseup event.
	AttrOnMouseUp = "onMouseUp"
	// AttrOnMouseWheel specifies the event handler for the mousewheel event.
	AttrOnMouseWheel = "onMouseWheel"
	// AttrOnMove specifies the event handler for the move event.
	AttrOnMove = "onMove"
	// AttrOnMoveEnd specifies the event handler for the moveend event.
	AttrOnMoveEnd = "onMoveEnd"
	// AttrOnMoveStart specifies the event handler for the movestart event.
	AttrOnMoveStart = "onMoveStart"
	// AttrOnOffline specifies the event handler for the offline event.
	AttrOnOffline = "onOffline"
	// AttrOnOnline specifies the event handler for the online event.
	AttrOnOnline = "onOnline"
	// AttrOnOutOfSync specifies the event handler for the outofsync event.
	AttrOnOutOfSync = "onOutOfSync"
	// AttrOnPaste specifies the event handler for the paste event.
	AttrOnPaste = "onPaste"
	// AttrOnPause specifies the event handler for the pause event.
	AttrOnPause = "onPause"
	// AttrOnPopState specifies the event handler for the popstate event.
	AttrOnPopState = "onPopState"
	// AttrOnProgress specifies the event handler for the progress event.
	AttrOnProgress = "onProgress"
	// AttrOnPropertyChange specifies the event handler for the propertychange event.
	AttrOnPropertyChange = "onPropertyChange"
	// AttrOnReadyStateChange specifies the event handler for the readystatechange event.
	AttrOnReadyStateChange = "onReadyStateChange"
	// AttrOnRedo specifies the event handler for the redo event.
	AttrOnRedo = "onRedo"
	// AttrOnRepeat specifies the event handler for the repeat event.
	AttrOnRepeat = "onRepeat"
	// AttrOnReset specifies the event handler for the reset event.
	AttrOnReset = "onReset"
	// AttrOnResize specifies the event handler for the resize event.
	AttrOnResize = "onResize"
	// AttrOnResizeEnd specifies the event handler for the resizeend event.
	AttrOnResizeEnd = "onResizeEnd"
	// AttrOnResizeStart specifies the event handler for the resizestart event.
	AttrOnResizeStart = "onResizeStart"
	// AttrOnResume specifies the event handler for the resume event.
	AttrOnResume = "onResume"
	// AttrOnReverse specifies the event handler for the reverse event.
	AttrOnReverse = "onReverse"
	// AttrOnRowsEnter specifies the event handler for the rowsenter event.
	AttrOnRowsEnter = "onRowsEnter"
	// AttrOnRowExit specifies the event handler for the rowexit event.
	AttrOnRowExit = "onRowExit"
	// AttrOnRowDelete specifies the event handler for the rowdelete event.
	AttrOnRowDelete = "onRowDelete"
	// AttrOnRowInserted specifies the event handler for the rowinserted event.
	AttrOnRowInserted = "onRowInserted"
	// AttrOnScroll specifies the event handler for the scroll event.
	AttrOnScroll = "onScroll"
	// AttrOnSeek specifies the event handler for the seek event.
	AttrOnSeek = "onSeek"
	// AttrOnSelect specifies the event handler for the select event.
	AttrOnSelect = "onSelect"
	// AttrOnSelectionChange specifies the event handler for the selectionchange event.
	AttrOnSelectionChange = "onSelectionChange"
	// AttrOnSelectStart specifies the event handler for the selectstart event.
	AttrOnSelectStart = "onSelectStart"
	// AttrOnStart specifies the event handler for the start event.
	AttrOnStart = "onStart"
	// AttrOnStop specifies the event handler for the stop event.
	AttrOnStop = "onStop"
	// AttrOnStorage specifies the event handler for the storage event.
	AttrOnStorage = "onStorage"
	// AttrOnSyncRestored specifies the event handler for the syncrestored event.
	AttrOnSyncRestored = "onSyncRestored"
	// AttrOnSubmit specifies the event handler for the submit event.
	AttrOnSubmit = "onSubmit"
	// AttrOnTimeError specifies the event handler for the timeerror event.
	AttrOnTimeError = "onTimeError"
	// AttrOnTrackChange specifies the event handler for the trackchange event.
	AttrOnTrackChange = "onTrackChange"
	// AttrOnUndo specifies the event handler for the undo event.
	AttrOnUndo = "onUndo"
	// AttrOnUnload specifies the event handler for the unload event.
	AttrOnUnload = "onUnload"
	// AttrOnUrlFlip specifies the event handler for the urlflip event.
	AttrOnUrlFlip = "onUrlFlip"
	// AttrOpen specifies whether the element is visible (for details, dialog, etc.).
	AttrOpen = "open"
	// AttrOptimum specifies the optimal value in a range.
	AttrOptimum = "optimum"
	// AttrPattern specifies a regular expression for input validation.
	AttrPattern = "pattern"
	// AttrPing specifies a list of URLs to notify when a link is clicked.
	AttrPing = "ping"
	// AttrPlaceholder provides a hint to the user about what to enter.
	AttrPlaceholder = "placeholder"
	// AttrPlaysInline specifies that the video should play inline.
	AttrPlaysInline = "playsinline"
	// AttrPoster specifies the preview image for a video.
	AttrPoster = "poster"
	// AttrPopoverTargetAction specifies the action to perform with a popover element.
	AttrPopoverTargetAction = "popovertargetaction"
	// AttrPreload specifies how to preload an audio/video.
	AttrPreload = "preload"
	// AttrReadOnly specifies that an input field is read-only.
	AttrReadOnly = "readonly"
	// AttrReferrerPolicy specifies the referrer policy for the resource.
	AttrReferrerPolicy = "referrerpolicy"
	// AttrRel specifies the relationship between the current document and the linked resource.
	AttrRel = "rel"
	// AttrRequired specifies that an input field must be filled out.
	AttrRequired = "required"
	// AttrReversed specifies that the list order should be reversed.
	AttrReversed = "reversed"
	// AttrRole specifies the role of an element for accessibility.
	AttrRole = "role"
	// AttrRows specifies the number of rows in a textarea.
	AttrRows = "rows"
	// AttrRowSpan specifies the number of rows a table cell should span.
	AttrRowSpan = "rowspan"
	// AttrSandbox enables extra restrictions for an iframe.
	AttrSandbox = "sandbox"
	// AttrScope specifies the header cells that a th element applies to.
	AttrScope = "scope"
	// AttrSelected specifies that an option should be pre-selected.
	AttrSelected = "selected"
	// AttrShape specifies the shape of an area in an image map.
	AttrShape = "shape"
	// AttrSize specifies the size of an input field or select element.
	AttrSize = "size"
	// AttrSizes specifies the sizes of an image for different layouts.
	AttrSizes = "sizes"
	// AttrSlot assigns a slot to an element in a shadow DOM.
	AttrSlot = "slot"
	// AttrSpan specifies the number of columns in a colgroup.
	AttrSpan = "span"
	// AttrSpellCheck specifies whether to enable spell checking.
	AttrSpellCheck = "spellcheck"
	// AttrSrc specifies the URL of an image, audio, video, or iframe.
	AttrSrc = "src"
	// AttrSrcDoc specifies the inline HTML for an iframe.
	AttrSrcDoc = "srcdoc"
	// AttrSrcLang specifies the language of the track text.
	AttrSrcLang = "srclang"
	// AttrSrcSet specifies multiple image sources for responsive images.
	AttrSrcSet = "srcset"
	// AttrStart specifies the starting number of an ordered list.
	AttrStart = "start"
	// AttrStep specifies the interval between legal numbers in an input.
	AttrStep = "step"
	// AttrStyle specifies inline CSS styles.
	AttrStyle = "style"
	// AttrSummary provides a summary for a table.
	AttrSummary = "summary"
	// AttrTabIndex specifies the tab order of an element.
	AttrTabIndex = "tabindex"
	// AttrTarget specifies where to open a link or form response.
	AttrTarget = "target"
	// AttrTitle provides advisory information about an element.
	AttrTitle = "title"
	// AttrTranslate specifies whether to translate an element.
	AttrTranslate = "translate"
	// AttrType specifies the type of an input element.
	AttrType = "type"
	// AttrUseMap specifies that an image is a client-side image map.
	AttrUseMap = "usemap"
	// AttrValue specifies the value of an input element.
	AttrValue = "value"
	// AttrWidth specifies the width of an element.
	AttrWidth = "width"
	// AttrWrap specifies how text should wrap in a textarea.
	AttrWrap = "wrap"
)

// Type* constants are valid values for the type attribute on various elements.
const (
	// TypeText creates a single-line text input field.
	TypeText = "text"
	// TypePassword creates a password input field.
	TypePassword = "password"
	// TypeCheckbox creates a checkbox input field.
	TypeCheckbox = "checkbox"
	// TypeRadio creates a radio button input field.
	TypeRadio = "radio"
	// TypeSubmit creates a submit button.
	TypeSubmit = "submit"
	// TypeReset creates a reset button.
	TypeReset = "reset"
	// TypeButton creates a generic button.
	TypeButton = "button"
	// TypeFile creates a file upload input field.
	TypeFile = "file"
	// TypeHidden creates a hidden input field.
	TypeHidden = "hidden"
	// TypeImage creates an image submit button.
	TypeImage = "image"
	// TypeColor creates a color picker input field.
	TypeColor = "color"
	// TypeDate creates a date picker input field.
	TypeDate = "date"
	// TypeDateTime creates a date and time picker input field.
	TypeDateTime = "datetime"
	// TypeDateTimeLocal creates a local date and time picker input field.
	TypeDateTimeLocal = "datetime-local"
	// TypeEmail creates an email input field.
	TypeEmail = "email"
	// TypeMonth creates a month picker input field.
	TypeMonth = "month"
	// TypeNumber creates a number input field.
	TypeNumber = "number"
	// TypeRange creates a range slider input field.
	TypeRange = "range"
	// TypeSearch creates a search input field.
	TypeSearch = "search"
	// TypeTel creates a telephone number input field.
	TypeTel = "tel"
	// TypeTime creates a time picker input field.
	TypeTime = "time"
	// TypeUrl creates a URL input field.
	TypeUrl = "url"
	// TypeWeek creates a week picker input field.
	TypeWeek = "week"
)

// Rel* constants are valid values for the rel attribute.
const (
	// RelAlternate indicates an alternate version of the current document.
	RelAlternate = "alternate"
	// RelAuthor indicates the author of the current document.
	RelAuthor = "author"
	// RelBookmark indicates a bookmark for the current document.
	RelBookmark = "bookmark"
	// RelCanonical indicates the canonical URL of the current document.
	RelCanonical = "canonical"
	// RelCompressionDictionary indicates a compression dictionary resource.
	RelCompressionDictionary = "compression-dictionary"
	// RelDnsPrefetch indicates to pre-resolve the DNS of the linked resource.
	RelDnsPrefetch = "dns-prefetch"
	// RelExternal indicates the linked resource is not part of the current site.
	RelExternal = "external"
	// RelHelp indicates a link to a help resource.
	RelHelp = "help"
	// RelIcon indicates the favicon of the current document.
	RelIcon = "icon"
	// RelLicense indicates the copyright license for the current document.
	RelLicense = "license"
	// RelManifest indicates a Web App Manifest.
	RelManifest = "manifest"
	// RelModulePreload indicates to preload a JavaScript module.
	RelModulePreload = "modulepreload"
	// RelNext indicates the next document in a sequence.
	RelNext = "next"
	// RelNoFollow indicates the link is not endorsed by the author.
	RelNoFollow = "nofollow"
	// RelNoOpener prevents the opened page from accessing the source page.
	RelNoOpener = "noopener"
	// RelNoReferrer indicates not to send a referrer header.
	RelNoReferrer = "noreferrer"
	// RelOpener allows the opened page to access the source page.
	RelOpener = "opener"
	// RelPingback indicates the URL of a pingback server.
	RelPingback = "pingback"
	// RelPreconnect indicates to pre-connect to the linked resource.
	RelPreconnect = "preconnect"
	// RelPrefetch indicates to prefetch the linked resource.
	RelPrefetch = "prefetch"
	// RelPreload indicates to preload the linked resource.
	RelPreload = "preload"
	// RelPrerender indicates to prerender the linked resource.
	RelPrerender = "prerender"
	// RelPrev indicates the previous document in a sequence.
	RelPrev = "prev"
	// RelSearch indicates a search resource for the current document.
	RelSearch = "search"
	// RelStylesheet indicates a stylesheet for the current document.
	RelStylesheet = "stylesheet"
	// RelTag indicates a tag relevant to the current document.
	RelTag = "tag"
)

// Target* constants are valid values for the target attribute.
const (
	// TargetBlank opens the link in a new tab or window.
	TargetBlank = "_blank"
	// TargetSelf opens the link in the same frame as clicked.
	TargetSelf = "_self"
	// TargetParent opens the link in the parent frame.
	TargetParent = "_parent"
	// TargetTop opens the link in the topmost frame.
	TargetTop = "_top"
	// TargetFrameName opens the link in a named frame.
	TargetFrameName = "framename"
)

// Method* constants are valid values for the method attribute on <form>.
const (
	// MethodGet sends form data as URL parameters.
	MethodGet = "get"
	// MethodPost sends form data in the request body.
	MethodPost = "post"
)

// Enctype* constants are valid values for the enctype attribute on <form>.
const (
	// EnctypeUrlEncoded encodes form data as URL-encoded string.
	EnctypeUrlEncoded = "application/x-www-form-urlencoded"
	// EnctypeMultipartForm encodes form data as multipart/form-data.
	EnctypeMultipartForm = "multipart/form-data"
	// EnctypePlainText encodes form data as plain text.
	EnctypePlainText = "text/plain"
)

// CrossOrigin* constants are valid values for the crossorigin attribute.
const (
	// CrossOriginAnonymous allows anonymous cross-origin requests.
	CrossOriginAnonymous = "anonymous"
	// CrossOriginUseCredentials requires credentials for cross-origin requests.
	CrossOriginUseCredentials = "use-credentials"
)

// Dir* constants are valid values for the dir attribute.
const (
	// DirLtr sets text direction to left-to-right.
	DirLtr = "ltr"
	// DirRtl sets text direction to right-to-left.
	DirRtl = "rtl"
	// DirAuto sets text direction to automatically detected.
	DirAuto = "auto"
)

// Preload* constants are valid values for the preload attribute on <audio> and <video>.
const (
	// PreloadNone indicates not to preload the media.
	PreloadNone = "none"
	// PreloadMetadata indicates to preload only metadata.
	PreloadMetadata = "metadata"
	// PreloadAuto indicates to preload the entire media.
	PreloadAuto = "auto"
)

// Loading* constants are valid values for the loading attribute on <img> and <iframe>.
const (
	// LoadingLazy defers loading until the element is near the viewport.
	LoadingLazy = "lazy"
	// LoadingEager loads the element immediately.
	LoadingEager = "eager"
)

// Decoding* constants are valid values for the decoding attribute on <img>.
const (
	// DecodingAsync decodes the image asynchronously.
	DecodingAsync = "async"
	// DecodingSync decodes the image synchronously.
	DecodingSync = "sync"
	// DecodingAuto lets the browser decide.
	DecodingAuto = "auto"
)

// ReferrerPolicy* constants are valid values for the referrerpolicy attribute.
const (
	// ReferrerPolicyNoReferrer does not send a referrer header.
	ReferrerPolicyNoReferrer = "no-referrer"
	// ReferrerPolicyNoReferrerWhenDowngrade sends referrer only for same-origin or secure-to-insecure.
	ReferrerPolicyNoReferrerWhenDowngrade = "no-referrer-when-downgrade"
	// ReferrerPolicyOrigin sends only the origin as referrer.
	ReferrerPolicyOrigin = "origin"
	// ReferrerPolicyOriginWhenCrossOrigin sends full URL for same-origin, origin for cross-origin.
	ReferrerPolicyOriginWhenCrossOrigin = "origin-when-cross-origin"
	// ReferrerPolicySameOrigin sends referrer only for same-origin.
	ReferrerPolicySameOrigin = "same-origin"
	// ReferrerPolicyStrictOriginWhenCrossOrigin sends origin for cross-origin when downgrade.
	ReferrerPolicyStrictOriginWhenCrossOrigin = "strict-origin-when-cross-origin"
	// ReferrerPolicyUnsafeUrl sends the full URL in all requests.
	ReferrerPolicyUnsafeUrl = "unsafe-url"
)

// FetchPriority* constants are valid values for the fetchpriority attribute.
const (
	// FetchPriorityLow indicates low fetch priority.
	FetchPriorityLow = "low"
	// FetchPriorityHigh indicates high fetch priority.
	FetchPriorityHigh = "high"
	// FetchPriorityAuto lets the browser decide the priority.
	FetchPriorityAuto = "auto"
)

// EnterKeyHint* constants are valid values for the enterkeyhint attribute.
const (
	// EnterKeyHintEnter indicates the enter key should insert a newline.
	EnterKeyHintEnter = "enter"
	// EnterKeyHintDone indicates the enter key should indicate "done".
	EnterKeyHintDone = "done"
	// EnterKeyHintGo indicates the enter key should indicate "go".
	EnterKeyHintGo = "go"
	// EnterKeyHintNext indicates the enter key should indicate "next".
	EnterKeyHintNext = "next"
	// EnterKeyHintPrevious indicates the enter key should indicate "previous".
	EnterKeyHintPrevious = "previous"
	// EnterKeyHintSearch indicates the enter key should indicate "search".
	EnterKeyHintSearch = "search"
	// EnterKeyHintSend indicates the enter key should indicate "send".
	EnterKeyHintSend = "send"
)

// Wrap* constants are valid values for the wrap attribute on <textarea>.
const (
	// WrapHard specifies hard wrapping with preserved line breaks.
	WrapHard = "hard"
	// WrapSoft specifies soft wrapping without line breaks in the submitted value.
	WrapSoft = "soft"
	// WrapOff disables wrapping.
	WrapOff = "off"
)

// Shape* constants are valid values for the shape attribute on <area>.
const (
	// ShapeDefault defines the entire region.
	ShapeDefault = "default"
	// ShapeCircle defines a circular region.
	ShapeCircle = "circle"
	// ShapePoly defines a polygonal region.
	ShapePoly = "poly"
	// ShapeRect defines a rectangular region.
	ShapeRect = "rect"
)

// ContentEditable* constants are valid values for the contenteditable attribute.
const (
	// ContentEditableTrue makes the element editable.
	ContentEditableTrue = "true"
	// ContentEditableFalse makes the element non-editable.
	ContentEditableFalse = "false"
	// ContentEditablePlaintextOnly allows only plain text editing.
	ContentEditablePlaintextOnly = "plaintext-only"
)

// InputMode* constants are valid values for the inputmode attribute.
const (
	// InputModeNone indicates no virtual keyboard.
	InputModeNone = "none"
	// InputModeText indicates a text input mode.
	InputModeText = "text"
	// InputModeDecimal indicates a decimal number input mode.
	InputModeDecimal = "decimal"
	// InputModeNumeric indicates a numeric input mode.
	InputModeNumeric = "numeric"
	// InputModeTel indicates a telephone number input mode.
	InputModeTel = "tel"
	// InputModeSearch indicates a search input mode.
	InputModeSearch = "search"
	// InputModeEmail indicates an email input mode.
	InputModeEmail = "email"
	// InputModeUrl indicates a URL input mode.
	InputModeUrl = "url"
)

// PopoverTargetAction* constants are valid values for the popovertargetaction attribute.
const (
	// PopoverTargetActionHide hides the popover.
	PopoverTargetActionHide = "hide"
	// PopoverTargetActionShow shows the popover.
	PopoverTargetActionShow = "show"
	// PopoverTargetActionToggle toggles the popover visibility.
	PopoverTargetActionToggle = "toggle"
)

// Autocomplete* constants are valid values for the autocomplete attribute.
const (
	// AutocompleteOff disables autocomplete.
	AutocompleteOff = "off"
	// AutocompleteOn enables autocomplete.
	AutocompleteOn = "on"
	// AutocompleteName specifies the full name.
	AutocompleteName = "name"
	// AutocompleteHonorificPrefix specifies a honorific prefix (e.g., Mr, Mrs).
	AutocompleteHonorificPrefix = "honorific-prefix"
	// AutocompleteGivenName specifies the given (first) name.
	AutocompleteGivenName = "given-name"
	// AutocompleteAdditionalName specifies an additional name.
	AutocompleteAdditionalName = "additional-name"
	// AutocompleteFamilyName specifies the family (last) name.
	AutocompleteFamilyName = "family-name"
	// AutocompleteHonorificSuffix specifies a honorific suffix (e.g., Jr, III).
	AutocompleteHonorificSuffix = "honorific-suffix"
	// AutocompleteNickname specifies a nickname.
	AutocompleteNickname = "nickname"
	// AutocompleteEmail specifies an email address.
	AutocompleteEmail = "email"
	// AutocompleteUsername specifies a username.
	AutocompleteUsername = "username"
	// AutocompleteNewPassword specifies a new password (for signup).
	AutocompleteNewPassword = "new-password"
	// AutocompleteCurrentPassword specifies the current password (for login).
	AutocompleteCurrentPassword = "current-password"
	// AutocompleteOneTimeCode specifies a one-time code for authentication.
	AutocompleteOneTimeCode = "one-time-code"
	// AutocompleteOrganizationTitle specifies a job title or organizational title.
	AutocompleteOrganizationTitle = "organization-title"
	// AutocompleteOrganization specifies an organization name.
	AutocompleteOrganization = "organization"
	// AutocompleteStreetAddress specifies a street address.
	AutocompleteStreetAddress = "street-address"
	// AutocompleteAddressLine1 specifies the first line of an address.
	AutocompleteAddressLine1 = "address-line1"
	// AutocompleteAddressLine2 specifies the second line of an address.
	AutocompleteAddressLine2 = "address-line2"
	// AutocompleteAddressLine3 specifies the third line of an address.
	AutocompleteAddressLine3 = "address-line3"
	// AutocompleteAddressLevel1 specifies the first address level (e.g., country).
	AutocompleteAddressLevel1 = "address-level1"
	// AutocompleteAddressLevel2 specifies the second address level (e.g., state/province).
	AutocompleteAddressLevel2 = "address-level2"
	// AutocompleteAddressLevel3 specifies the third address level (e.g., city).
	AutocompleteAddressLevel3 = "address-level3"
	// AutocompleteAddressLevel4 specifies the most granular address level (e.g., neighborhood).
	AutocompleteAddressLevel4 = "address-level4"
	// AutocompleteCountry specifies the country code.
	AutocompleteCountry = "country"
	// AutocompleteCountryName specifies the country name.
	AutocompleteCountryName = "country-name"
	// AutocompleteCCName specifies the name on the credit card.
	AutocompleteCCName = "cc-name"
	// AutocompleteCCGivenName specifies the given name on the credit card.
	AutocompleteCCGivenName = "cc-given-name"
	// AutocompleteCCAdditionalName specifies the additional name on the credit card.
	AutocompleteCCAdditionalName = "cc-additional-name"
	// AutocompleteCCFamilyName specifies the family name on the credit card.
	AutocompleteCCFamilyName = "cc-family-name"
	// AutocompleteCCNumber specifies the credit card number.
	AutocompleteCCNumber = "cc-number"
	// AutocompleteCCExp specifies the credit card expiration date.
	AutocompleteCCExp = "cc-exp"
	// AutocompleteCCExpMonth specifies the credit card expiration month.
	AutocompleteCCExpMonth = "cc-exp-month"
	// AutocompleteCCExpYear specifies the credit card expiration year.
	AutocompleteCCExpYear = "cc-exp-year"
	// AutocompleteCCCsc specifies the credit card security code.
	AutocompleteCCCsc = "cc-csc"
	// AutocompleteCCType specifies the credit card type (e.g., Visa, Mastercard).
	AutocompleteCCType = "cc-type"
	// AutocompleteTransactionCurrency specifies the transaction currency.
	AutocompleteTransactionCurrency = "transaction-currency"
	// AutocompleteTransactionAmount specifies the transaction amount.
	AutocompleteTransactionAmount = "transaction-amount"
	// AutocompleteLanguage specifies a language tag.
	AutocompleteLanguage = "language"
	// AutocompleteBday specifies a birth date.
	AutocompleteBday = "bday"
	// AutocompleteBdayDay specifies the day of birth.
	AutocompleteBdayDay = "bday-day"
	// AutocompleteBdayMonth specifies the month of birth.
	AutocompleteBdayMonth = "bday-month"
	// AutocompleteBdayYear specifies the year of birth.
	AutocompleteBdayYear = "bday-year"
	// AutocompleteSex specifies a gender identity.
	AutocompleteSex = "sex"
	// AutocompleteTelCountryCode specifies the country code component of a telephone number.
	AutocompleteTelCountryCode = "tel-country-code"
	// AutocompleteTelNational specifies the telephone number without country code.
	AutocompleteTelNational = "tel-national"
	// AutocompleteTelAreaCode specifies the area code component of a telephone number.
	AutocompleteTelAreaCode = "tel-area-code"
	// AutocompleteTelLocal specifies the local telephone number.
	AutocompleteTelLocal = "tel-local"
	// AutocompleteTelExtension specifies a telephone extension code.
	AutocompleteTelExtension = "tel-extension"
	// AutocompleteImpp specifies an instant messaging protocol URL.
	AutocompleteImpp = "impp"
	// AutocompleteUrl specifies a URL.
	AutocompleteUrl = "url"
	// AutocompletePhoto specifies a photo URL.
	AutocompletePhoto = "photo"
)

// Sandbox* constants are valid values for the sandbox attribute on <iframe>.
const (
	// SandboxAllowForms allows form submission in the iframe.
	SandboxAllowForms = "allow-forms"
	// SandboxAllowModals allows modal dialogs in the iframe.
	SandboxAllowModals = "allow-modals"
	// SandboxAllowOrientationLock allows screen orientation lock in the iframe.
	SandboxAllowOrientationLock = "allow-orientation-lock"
	// SandboxAllowPointerLock allows pointer lock in the iframe.
	SandboxAllowPointerLock = "allow-pointer-lock"
	// SandboxAllowPopups allows popups in the iframe.
	SandboxAllowPopups = "allow-popups"
	// SandboxAllowPopupsToEscapeSandbox allows popups to escape sandbox restrictions.
	SandboxAllowPopupsToEscapeSandbox = "allow-popups-to-escape-sandbox"
	// SandboxAllowPresentation allows presentation mode in the iframe.
	SandboxAllowPresentation = "allow-presentation"
	// SandboxAllowSameOrigin allows the iframe to access same-origin content.
	SandboxAllowSameOrigin = "allow-same-origin"
	// SandboxAllowScripts allows JavaScript execution in the iframe.
	SandboxAllowScripts = "allow-scripts"
	// SandboxAllowTopNavigation allows top-level navigation in the iframe.
	SandboxAllowTopNavigation = "allow-top-navigation"
	// SandboxAllowTopNavigationByUserActivation allows top-level navigation only by user gesture.
	SandboxAllowTopNavigationByUserActivation = "allow-top-navigation-by-user-activation"
)

// Element represents an HTML element with its attributes and children.
type Element struct {
	Tag        string      // HTML tag name
	IsVoid     bool        // Whether the tag is self-closing (e.g., <br>, <img>)
	Attributes []attribute // HTML attributes as key-value pairs
	Children   []HyperNode // Child nodes
}

// Render generates the HTML for the element and its children to the provided writer.
func (me Element) Render(w io.Writer) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	if err := me.renderElement(buf); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

// bufferPool is a sync pool for reusing byte buffers during HTML rendering.
// This reduces allocations when rendering many elements by recycling buffers
// with a pre-allocated capacity of 1KB.
var bufferPool = sync.Pool{
	New: func() any {
		var buf bytes.Buffer
		buf.Grow(1024)
		return &buf
	},
}

// renderElement renders the element to the provided buffer.
func (me Element) renderElement(buf *bytes.Buffer) error {
	if me.Tag == "" {
		return me.renderChildren(buf)
	}

	buf.WriteByte('<')
	buf.WriteString(me.Tag)
	if err := me.renderAttrs(buf); err != nil {
		return err
	}
	buf.WriteByte('>')

	if me.IsVoid {
		return nil
	}

	if err := me.renderChildren(buf); err != nil {
		return err
	}

	buf.WriteString("</")
	buf.WriteString(me.Tag)
	buf.WriteByte('>')
	return nil
}

// renderChildren renders all child nodes to the provided buffer.
func (me Element) renderChildren(buf *bytes.Buffer) error {
	for _, child := range me.Children {
		switch c := child.(type) {
		// I'm tring to pass the concrete type [bytes.Buffer] as possible.
		// That's why I'm not using Render(buf).
		case Element:
			if err := c.renderElement(buf); err != nil {
				return err
			}
		case Text:
			buf.WriteString(html.EscapeString(string(c)))
		case RawText:
			buf.WriteString(string(c))

		default:
			if err := c.Render(buf); err != nil {
				return err
			}
		}
	}

	return nil
}

func (me Element) renderAttrs(buf *bytes.Buffer) error {
	for _, attr := range me.Attributes {
		k := strings.TrimSpace(attr.key)
		if k == "" {
			return fmt.Errorf("empty/whitespace attribute key not allowed.")
		}
		if attr.value == nil {
			return fmt.Errorf("attribute '%s' has nil value", k)
		}

		switch v := attr.value.(type) {
		case string:
			buf.WriteByte(' ')
			buf.WriteString(html.EscapeString(k))
			buf.WriteString(`="`)
			buf.WriteString(strings.ReplaceAll(v, `"`, "&quot;"))
			buf.WriteByte('"')
		case bool:
			if v {
				buf.WriteByte(' ')
				buf.WriteString(html.EscapeString(k))
			}
		default:
			return fmt.Errorf("attribute value must be string or bool, got %T for key '%s'", v, k)
		}
	}
	return nil
}

// newElem creates an HTML element with the given tag name and arguments.
// Arguments can be KV for attributes, Node for children, or other types to convert to text.
func newElem(tag string, args ...any) Element {
	e := Element{Tag: tag}
	for _, arg := range args {
		// NOTE: nil arguments are not checked/filtered - fmt.Sprint() will render them as "Nil",
		// which is intentional for better debugging (makes it obvious when nil values are passed).
		switch value := arg.(type) {
		case KV:
			e.fillAttrsWithKV(value)
		case HyperNode:
			e.Children = append(e.Children, value)

		// Explicit string and fmt.Stringer cases for performance:
		// fmt.Sprint() would handle these, but with overhead from type inspection and buffer allocation.
		case string:
			e.Children = append(e.Children, Text(value))
		case fmt.Stringer:
			e.Children = append(e.Children, Text(value.String()))
		default:
			e.Children = append(e.Children, Text(fmt.Sprint(value)))
		}
	}
	return e
}

// newVoidElem creates a self-closing (void) HTML element with the given tag name and attributes.
func newVoidElem(tag string, attrs ...KV) Element {
	e := Element{Tag: tag, IsVoid: true}
	for _, kv := range attrs {
		e.fillAttrsWithKV(kv)
	}
	return e
}

// fillAttrsWithKV appends key-value attributes to the attributes slice and returns the updated slice.
func (me *Element) fillAttrsWithKV(kv KV) {
	// Two use cases:
	//   - First occurrence (nil attrs): Allocate a slice with exact length from the KV map.
	//     This is the common case where users pass a single KV{}.
	//   - Subsequent calls: Use memory-efficient growth when extending the slice.
	//     This handles the uncommon case where multiple KV{} are passed.
	if me.Attributes == nil {
		me.Attributes = make([]attribute, 0, len(kv))
	} else {
		if len(kv) > cap(me.Attributes)-len(me.Attributes) {
			required := len(me.Attributes) + len(kv)
			newCap := required * 2
			for newCap > 32 {
				newCap = required + 8
			}
			newSlice := make([]attribute, len(me.Attributes), newCap)
			copy(newSlice, me.Attributes)
			me.Attributes = newSlice
		}
	}

	for k, v := range kv {
		me.Attributes = append(me.Attributes, attribute{key: k, value: v})
	}
}

// EMPTY creates an empty element (no tag).
func EMPTY(args ...any) HyperNode {
	return newElem("", args...)
}

// DOCTYPE creates the <!DOCTYPE html> element.
//
// https://developer.mozilla.org/en-US/docs/Glossary/Doctype
func DOCTYPE() HyperNode {
	return newVoidElem("!DOCTYPE html")
}

// HTML creates the root element of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
func HTML(args ...any) HyperNode {
	return newElem("html", args...)
}

// HEAD contains machine-readable information about the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
func HEAD(args ...any) HyperNode {
	return newElem("head", args...)
}

// TITLE defines the document's title that is shown in a browser's title bar or a page's tab.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
func TITLE(args ...any) HyperNode {
	return newElem("title", args...)
}

// LINK specifies relationships between the current document and an external resource.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
func LINK(attrs ...KV) HyperNode {
	return newVoidElem("link", attrs...)
}

// META represents metadata that cannot be represented by other HTML meta-related elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
func META(attrs ...KV) HyperNode {
	return newVoidElem("meta", attrs...)
}

// STYLE contains style information for a document or part of a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
func STYLE(args ...any) HyperNode {
	return newElem("style", args...)
}

// BODY represents the content of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
func BODY(args ...any) HyperNode {
	return newElem("body", args...)
}

// H1 creates a level 1 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h1
func H1(args ...any) HyperNode {
	return newElem("h1", args...)
}

// H2 creates a level 2 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h2
func H2(args ...any) HyperNode {
	return newElem("h2", args...)
}

// H3 creates a level 3 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h3
func H3(args ...any) HyperNode {
	return newElem("h3", args...)
}

// H4 creates a level 4 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h4
func H4(args ...any) HyperNode {
	return newElem("h4", args...)
}

// H5 creates a level 5 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h5
func H5(args ...any) HyperNode {
	return newElem("h5", args...)
}

// H6 creates a level 6 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h6
func H6(args ...any) HyperNode {
	return newElem("h6", args...)
}

// HEADER creates a header element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
func HEADER(args ...any) HyperNode {
	return newElem("header", args...)
}

// FOOTER creates a footer element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
func FOOTER(args ...any) HyperNode {
	return newElem("footer", args...)
}

// NAV creates a navigation element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/nav
func NAV(args ...any) HyperNode {
	return newElem("nav", args...)
}

// MAIN creates a main content element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/main
func MAIN(args ...any) HyperNode {
	return newElem("main", args...)
}

// SECTION creates a section element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/section
func SECTION(args ...any) HyperNode {
	return newElem("section", args...)
}

// ARTICLE creates an article element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
func ARTICLE(args ...any) HyperNode {
	return newElem("article", args...)
}

// ASIDE creates an aside element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
func ASIDE(args ...any) HyperNode {
	return newElem("aside", args...)
}

// HR represents a thematic break between paragraph-level elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hr
func HR(attrs ...KV) HyperNode {
	return newVoidElem("hr", attrs...)
}

// PRE represents preformatted text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/pre
func PRE(args ...any) HyperNode {
	return newElem("pre", args...)
}

// BLOCKQUOTE represents a section quoted from another source.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote
func BLOCKQUOTE(args ...any) HyperNode {
	return newElem("blockquote", args...)
}

// OL represents an ordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ol
func OL(args ...any) HyperNode {
	return newElem("ol", args...)
}

// UL represents an unordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ul
func UL(args ...any) HyperNode {
	return newElem("ul", args...)
}

// LI represents a list item.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/li
func LI(args ...any) HyperNode {
	return newElem("li", args...)
}

// A creates hyperlinks to other web pages, files, locations within the same page, or anything else a URL can address.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a
func A(args ...any) HyperNode {
	return newElem("a", args...)
}

// EM marks text with emphasis.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/em
func EM(args ...any) HyperNode {
	return newElem("em", args...)
}

// STRONG indicates strong importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/strong
func STRONG(args ...any) HyperNode {
	return newElem("strong", args...)
}

// CODE displays its contents styled as computer code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/code
func CODE(args ...any) HyperNode {
	return newElem("code", args...)
}

// VAR represents a variable in a mathematical expression or programming context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/var
func VAR(args ...any) HyperNode {
	return newElem("var", args...)
}

// SAMP represents sample output from a computer program.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/samp
func SAMP(args ...any) HyperNode {
	return newElem("samp", args...)
}

// KBD represents text that the user should enter.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/kbd
func KBD(args ...any) HyperNode {
	return newElem("kbd", args...)
}

// SUB specifies inline text displayed as subscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sub
func SUB(args ...any) HyperNode {
	return newElem("sub", args...)
}

// SUP specifies inline text displayed as superscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sup
func SUP(args ...any) HyperNode {
	return newElem("sup", args...)
}

// I represents text in an alternate voice or mood.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/i
func I(args ...any) HyperNode {
	return newElem("i", args...)
}

// B draws attention to text without conveying importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/b
func B(args ...any) HyperNode {
	return newElem("b", args...)
}

// U represents text with an unarticulated annotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/u
func U(args ...any) HyperNode {
	return newElem("u", args...)
}

// MARK highlights text for reference.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/mark
func MARK(args ...any) HyperNode {
	return newElem("mark", args...)
}

// BDI isolates text for bidirectional text formatting.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdi
func BDI(args ...any) HyperNode {
	return newElem("bdi", args...)
}

// BDO overrides the current text direction.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdo
func BDO(args ...any) HyperNode {
	return newElem("bdo", args...)
}

// BR produces a line break in text (carriage-return). It is useful for writing a poem or an address, where the division of lines is significant.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/br
func BR(attrs ...KV) HyperNode {
	return newVoidElem("br", attrs...)
}

// WBR represents a word break opportunity.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/wbr
func WBR(attrs ...KV) HyperNode {
	return newVoidElem("wbr", attrs...)
}

// IMG embeds an image into the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img
func IMG(attrs ...KV) HyperNode {
	return newVoidElem("img", attrs...)
}

// IFRAME represents a nested browsing context, embedding another HTML page into the current one.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe
func IFRAME(args ...any) HyperNode {
	return newElem("iframe", args...)
}

// EMBED embeds external content at the specified point in the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/embed
func EMBED(attrs ...KV) HyperNode {
	return newVoidElem("embed", attrs...)
}

// OBJECT represents an external resource, which can be treated as an image, a nested browsing context, or a resource to be handled by a plugin.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/object
func OBJECT(args ...any) HyperNode {
	return newElem("object", args...)
}

// PICTURE defines multiple sources for an img element to offer alternative versions of an image for different display/device scenarios.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture
func PICTURE(args ...any) HyperNode {
	return newElem("picture", args...)
}

// SOURCE specifies multiple media resources for the picture, the audio element, or the video element. It is a void element, meaning that it has no content and does not have a closing tag. It is commonly used to offer the same media content in multiple file formats in order to provide compatibility with a broad range of browsers given their differing support for image file formats and media file formats.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/source
func SOURCE(attrs ...KV) HyperNode {
	return newVoidElem("source", attrs...)
}

// TRACK is used as a child of the media elements, audio and video. It lets you specify timed text tracks (or time-based data), for example to automatically handle subtitles. The tracks are formatted in WebVTT format (.vtt files)—Web Video Text Tracks.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/track
func TRACK(attrs ...KV) HyperNode {
	return newVoidElem("track", attrs...)
}

// VIDEO embeds a media player which supports video playback into the document. You can also use &lt;video&gt; for audio content, but the audio element may provide a more appropriate user experience.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video
func VIDEO(args ...any) HyperNode {
	return newElem("video", args...)
}

// AUDIO is used to embed sound content in documents. It may contain one or more audio sources, represented using the src attribute or the source element: the browser will choose the most suitable one. It can also be the destination for streamed media, using a MediaStream.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/audio
func AUDIO(args ...any) HyperNode {
	return newElem("audio", args...)
}

// CANVAS is a container element to use with either the canvas scripting API or the WebGL API to draw graphics and animations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/canvas
func CANVAS(args ...any) HyperNode {
	return newElem("canvas", args...)
}

// MAP is used with &lt;area&gt; elements to define an image map (a clickable link area).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/map
func MAP(args ...any) HyperNode {
	return newElem("map", args...)
}

// AREA defines an area inside an image map that has predefined clickable areas. An image map allows geometric areas on an image to be associated with hyperlink.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/area
func AREA(attrs ...KV) HyperNode {
	return newVoidElem("area", attrs...)
}

// SVG is a container defining a new coordinate system and viewport. It is used as the outermost element of SVG documents, but it can also be used to embed an SVG fragment inside an SVG or HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/svg
func SVG(args ...any) HyperNode {
	return newElem("svg", args...)
}

// MATH is the top-level element in MathML. Every valid MathML instance must be wrapped in it. In addition, you must not nest a second &lt;math&gt; element in another, but you can have an arbitrary number of other child elements in it.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/math
func MATH(args ...any) HyperNode {
	return newElem("math", args...)
}

// SCRIPT is used to embed executable code or data; this is typically used to embed or refer to JavaScript code. The &lt;script&gt; element can also be used with other languages, such as WebGL's GLSL shader programming language and JSON.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script
func SCRIPT(args ...any) HyperNode {
	return newElem("script", args...)
}

// NOSCRIPT defines a section of HTML to be inserted if a script type on the page is unsupported or if scripting is currently turned off in the browser.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/noscript
func NOSCRIPT(args ...any) HyperNode {
	return newElem("noscript", args...)
}

// DEL represents a range of text that has been deleted from a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/del
func DEL(args ...any) HyperNode {
	return newElem("del", args...)
}

// INS represents a range of text that has been added to a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ins
func INS(args ...any) HyperNode {
	return newElem("ins", args...)
}

// TABLE represents tabular data—that is, information presented in a two-dimensional table comprised of rows and columns of cells containing data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/table
func TABLE(args ...any) HyperNode {
	return newElem("table", args...)
}

// CAPTION specifies the caption (or title) of a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/caption
func CAPTION(args ...any) HyperNode {
	return newElem("caption", args...)
}

// COLGROUP defines a group of columns within a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/colgroup
func COLGROUP(args ...any) HyperNode {
	return newElem("colgroup", args...)
}

// COL defines one or more columns in a column group represented by its implicit or explicit parent &lt;colgroup&gt; element. The &lt;col&gt; element is only valid as a child of a &lt;colgroup&gt; element that has no span attribute defined.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/col
func COL(attrs ...KV) HyperNode {
	return newVoidElem("col", attrs...)
}

// THEAD groups the header content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/thead
func THEAD(args ...any) HyperNode {
	return newElem("thead", args...)
}

// TBODY groups the body content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tbody
func TBODY(args ...any) HyperNode {
	return newElem("tbody", args...)
}

// TFOOT groups the footer content in a table with information about the table's columns. This is usually a summary of the columns, e.g., a sum of the given numbers in a column.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tfoot
func TFOOT(args ...any) HyperNode {
	return newElem("tfoot", args...)
}

// TR defines a row of cells in a table. The row's cells can then be established using a mix of &lt;td&gt; (data cell) and &lt;th&gt; (header cell) elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tr
func TR(args ...any) HyperNode {
	return newElem("tr", args...)
}

// TH is a child of the &lt;tr&gt; element, it defines a cell as the header of a group of table cells. The nature of this group can be explicitly defined by the scope and headers attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/th
func TH(args ...any) HyperNode {
	return newElem("th", args...)
}

// TD is a child of the &lt;tr&gt; element, it defines a cell of a table that contains data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/td
func TD(args ...any) HyperNode {
	return newElem("td", args...)
}

// FORM represents a document section containing interactive controls for submitting information.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form
func FORM(args ...any) HyperNode {
	return newElem("form", args...)
}

// FIELDSET is used to group several controls as well as labels (&lt;label&gt;) within a web form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fieldset
func FIELDSET(args ...any) HyperNode {
	return newElem("fieldset", args...)
}

// LEGEND represents a caption for the content of its parent &lt;fieldset&gt;.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/legend
func LEGEND(args ...any) HyperNode {
	return newElem("legend", args...)
}

// LABEL represents a caption for an item in a user interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label
func LABEL(args ...any) HyperNode {
	return newElem("label", args...)
}

// INPUT is used to create interactive controls for web-based forms to accept data from the user; a wide variety of types of input data and control widgets are available, depending on the device and user agent. The &lt;input&gt; element is one of the most powerful and complex in all of HTML due to the sheer number of combinations of input types and attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input
func INPUT(attrs ...KV) HyperNode {
	return newVoidElem("input", attrs...)
}

// BUTTON is an interactive element activated by a user with a mouse, keyboard, finger, voice command, or other assistive technology. Once activated, it performs an action, such as submitting a form or opening a dialog.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/button
func BUTTON(args ...any) HyperNode {
	return newElem("button", args...)
}

// SELECT represents a control that provides a menu of options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/select
func SELECT(args ...any) HyperNode {
	return newElem("select", args...)
}

// DATALIST contains a set of &lt;option&gt; elements that represent the permissible or recommended options available to choose from within other controls.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/datalist
func DATALIST(args ...any) HyperNode {
	return newElem("datalist", args...)
}

// OPTGROUP creates a grouping of options within a &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/optgroup
func OPTGROUP(args ...any) HyperNode {
	return newElem("optgroup", args...)
}

// OPTION is used to define an item contained in a &lt;select&gt;, an &lt;optgroup&gt;, or a &lt;datalist&gt; element. As such, &lt;option&gt; can represent menu items in popups and other lists of items in an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/option
func OPTION(args ...any) HyperNode {
	return newElem("option", args...)
}

// TEXTAREA represents a multi-line plain-text editing control, useful when you want to allow users to enter a sizeable amount of free-form text, for example, a comment on a review or feedback form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/textarea
func TEXTAREA(args ...any) HyperNode {
	return newElem("textarea", args...)
}

// OUTPUT is a container element into which a site or app can inject the results of a calculation or the outcome of a user action.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/output
func OUTPUT(args ...any) HyperNode {
	return newElem("output", args...)
}

// PROGRESS displays an indicator showing the completion progress of a task, typically displayed as a progress bar.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/progress
func PROGRESS(args ...any) HyperNode {
	return newElem("progress", args...)
}

// METER represents either a scalar value within a known range or a fractional value.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meter
func METER(args ...any) HyperNode {
	return newElem("meter", args...)
}

// DETAILS creates a disclosure widget in which information is visible only when the widget is toggled into an "open" state. A summary or label must be provided using the &lt;summary&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details
func DETAILS(args ...any) HyperNode {
	return newElem("details", args...)
}

// SUMMARY specifies a summary, caption, or legend for a details element's disclosure box. Clicking the &lt;summary&gt; element toggles the state of the parent &lt;details&gt; element open and closed.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/summary
func SUMMARY(args ...any) HyperNode {
	return newElem("summary", args...)
}

// DIALOG represents a dialog box or other interactive component, such as a dismissible alert, inspector, or subwindow.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog
func DIALOG(args ...any) HyperNode {
	return newElem("dialog", args...)
}

// SLOT acts as a placeholder inside a web component that you can fill with your own markup, which lets you create separate DOM trees and present them together.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/slot
func SLOT(args ...any) HyperNode {
	return newElem("slot", args...)
}

// TEMPLATE holds HTML that is not to be rendered immediately when a page is loaded but may be instantiated subsequently during runtime using JavaScript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template
func TEMPLATE(args ...any) HyperNode {
	return newElem("template", args...)
}

// FENCEDFRAME represents a nested browsing context, like &lt;iframe&gt; but with more native privacy features built in.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fencedframe
func FENCEDFRAME(args ...any) HyperNode {
	return newElem("fencedframe", args...)
}

// SELECTEDCONTENT displays the content of the currently selected &lt;option&gt; inside a closed &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/selectedcontent
func SELECTEDCONTENT(args ...any) HyperNode {
	return newElem("selectedcontent", args...)
}

// BASE specifies the base URL and default browsing context for relative URLs.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
func BASE(attrs ...KV) HyperNode {
	return newVoidElem("base", attrs...)
}

// HGROUP groups a set of h1–h6 elements when they represent a multi-level heading.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hgroup
func HGROUP(args ...any) HyperNode {
	return newElem("hgroup", args...)
}

// ADDRESS indicates contact information for a person or organization.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
func ADDRESS(args ...any) HyperNode {
	return newElem("address", args...)
}

// SEARCH represents a search or filtering interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/search
func SEARCH(args ...any) HyperNode {
	return newElem("search", args...)
}

// DIV is the generic container for flow content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/div
func DIV(args ...any) HyperNode {
	return newElem("div", args...)
}

// SPAN is the generic inline container for phrasing content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/span
func SPAN(args ...any) HyperNode {
	return newElem("span", args...)
}

// P creates a paragraph element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/p
func P(args ...any) HyperNode {
	return newElem("p", args...)
}

// DL represents a description list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dl
func DL(args ...any) HyperNode {
	return newElem("dl", args...)
}

// DT specifies a term in a description or definition list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dt
func DT(args ...any) HyperNode {
	return newElem("dt", args...)
}

// DD provides the description, definition, or value for the preceding term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dd
func DD(args ...any) HyperNode {
	return newElem("dd", args...)
}

// FIGURE represents self-contained content with an optional caption.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figure
func FIGURE(args ...any) HyperNode {
	return newElem("figure", args...)
}

// FIGCAPTION represents a caption or legend for the contents of its parent figure element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figcaption
func FIGCAPTION(args ...any) HyperNode {
	return newElem("figcaption", args...)
}

// MENU represents a set of commands or options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/menu
func MENU(args ...any) HyperNode {
	return newElem("menu", args...)
}

// SMALL represents side-comments and small print.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/small
func SMALL(args ...any) HyperNode {
	return newElem("small", args...)
}

// S renders text with a strikethrough.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/s
func S(args ...any) HyperNode {
	return newElem("s", args...)
}

// CITE marks the title of a creative work.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/cite
func CITE(args ...any) HyperNode {
	return newElem("cite", args...)
}

// Q indicates a short inline quotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/q
func Q(args ...any) HyperNode {
	return newElem("q", args...)
}

// DFN indicates the defining instance of a term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dfn
func DFN(args ...any) HyperNode {
	return newElem("dfn", args...)
}

// ABBR represents an abbreviation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/abbr
func ABBR(args ...any) HyperNode {
	return newElem("abbr", args...)
}

// RUBY represents ruby annotations for East Asian typography.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ruby
func RUBY(args ...any) HyperNode {
	return newElem("ruby", args...)
}

// RT specifies the ruby text for ruby annotations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rt
func RT(args ...any) HyperNode {
	return newElem("rt", args...)
}

// RP provides parentheses for browsers that don't support ruby text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rp
func RP(args ...any) HyperNode {
	return newElem("rp", args...)
}

// DATA links content with a machine-readable translation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/data
func DATA(args ...any) HyperNode {
	return newElem("data", args...)
}

// TIME represents a specific period in time.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/time
func TIME(args ...any) HyperNode {
	return newElem("time", args...)
}
