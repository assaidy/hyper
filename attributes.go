package h

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

// Attribute represents an HTML attribute that can be rendered.
// Implementations include PairAttribute (key="value") and BooleanAttribute (present/absent).
type Attribute interface {
	Render(buf *bytes.Buffer) error
}

// PairAttribute represents an HTML attribute with a key and value (key="value").
type PairAttribute struct {
	Key   string
	Value string
}

func (me PairAttribute) Render(buf *bytes.Buffer) error {
	k := strings.TrimSpace(me.Key)
	if k == "" {
		return fmt.Errorf("empty/whitespace attribute key not allowed.")
	}

	buf.WriteByte(' ')
	buf.WriteString(html.EscapeString(k))
	buf.WriteString(`="`)
	buf.WriteString(strings.ReplaceAll(me.Value, `"`, "&quot;"))
	buf.WriteByte('"')

	return nil
}

// BooleanAttribute represents an HTML boolean attribute that is either present or absent.
// When IsActive is true, the attribute is rendered; otherwise it is omitted.
type BooleanAttribute struct {
	Key      string
	IsActive bool
}

func (me BooleanAttribute) Render(buf *bytes.Buffer) error {
	k := strings.TrimSpace(me.Key)
	if k == "" {
		return fmt.Errorf("empty/whitespace attribute key not allowed.")
	}

	if me.IsActive {
		buf.WriteByte(' ')
		buf.WriteString(html.EscapeString(k))
	}

	return nil
}

// Attr creates an attribute from a key and value.
// If value is a string, it creates a PairAttribute (key="value").
// If value is a bool, it creates a BooleanAttribute (present when true, absent when false).
//
// Examples:
//
//	Attr("class", "container")  // -> PairAttribute{Key: "class", Value: "container"}
//	Attr("hidden", true)        // -> BooleanAttribute{Key: "hidden", IsActive: true}
//	Attr("disabled", false)     // -> BooleanAttribute{Key: "disabled", IsActive: false}
func Attr[V ~string | ~bool](key string, value V) Attribute {
	return attr(key, value)
}

func attr(key string, value any) Attribute {
	switch v := value.(type) {
	case string:
		return PairAttribute{Key: key, Value: v}
	case bool:
		return BooleanAttribute{Key: key, IsActive: v}
	default:
		panic("unexpected value type for attribute")
	}
}

func makePairAttribute(key string) func(value string) PairAttribute {
	return func(value string) PairAttribute {
		return PairAttribute{Key: key, Value: value}
	}
}

func makeBooleanAttribute(key string) func(isActive bool) BooleanAttribute {
	return func(isActive bool) BooleanAttribute {
		return BooleanAttribute{Key: key, IsActive: isActive}
	}
}

var (
	// AttrAccept sets the accepted file types for <input type="file">.
	AttrAccept = makePairAttribute("accept")
	// AttrAcceptCharset sets the character encodings accepted by the server.
	AttrAcceptCharset = makePairAttribute("accept-charset")
	// AttrAccessKey gives keyboard shortcut access to an element.
	AttrAccessKey = makePairAttribute("accesskey")
	// AttrAction specifies where to send the form data.
	AttrAction = makePairAttribute("action")
	// AttrAlign specifies the alignment of an element.
	AttrAlign = makePairAttribute("align")
	// AttrAllow specifies permissions for an iframe.
	AttrAllow = makePairAttribute("allow")
	// AttrAlpha sets the alpha transparency level of an element.
	AttrAlpha = makePairAttribute("alpha")
	// AttrAlt provides alternative text for an image.
	AttrAlt = makePairAttribute("alt")
	// AttrAs specifies the relation between the linked resource and the document.
	AttrAs = makePairAttribute("as")
	// AttrAsync indicates that the script should execute asynchronously.
	AttrAsync = makeBooleanAttribute("async")
	// AttrAutocapitalize controls whether text input is automatically capitalized.
	AttrAutocapitalize = makePairAttribute("autocapitalize")
	// AttrAutocomplete specifies whether an input field should have autocomplete enabled.
	AttrAutocomplete = makePairAttribute("autocomplete")
	// AttrAutofocus specifies that an element should automatically get focus on page load.
	AttrAutofocus = makeBooleanAttribute("autofocus")
	// AttrAutoplay specifies that the audio/video should automatically start playing.
	AttrAutoplay = makeBooleanAttribute("autoplay")
	// AttrBackground specifies the background image URL.
	AttrBackground = makePairAttribute("background")
	// AttrBGColor specifies the background color of an element.
	AttrBGColor = makePairAttribute("bgcolor")
	// AttrBorder specifies the border width around an element.
	AttrBorder = makePairAttribute("border")
	// AttrCapture specifies which camera/mic to use for media capture.
	AttrCapture = makePairAttribute("capture")
	// AttrCharset specifies the character encoding of the document.
	AttrCharset = makePairAttribute("charset")
	// AttrChecked specifies whether an input checkbox or radio is checked.
	AttrChecked = makeBooleanAttribute("checked")
	// AttrCite specifies the source of a quotation.
	AttrCite = makePairAttribute("cite")
	// AttrClass specifies one or more class names for an element.
	AttrClass = makePairAttribute("class")
	// AttrColor specifies the text color of an element.
	AttrColor = makePairAttribute("color")
	// AttrColorSpace specifies the color space for an image.
	AttrColorSpace = makePairAttribute("colorspace")
	// AttrCols specifies the number of columns in a textarea.
	AttrCols = makePairAttribute("cols")
	// AttrColSpan specifies the number of columns a table cell should span.
	AttrColSpan = makePairAttribute("colspan")
	// AttrContent provides metadata about the element.
	AttrContent = makePairAttribute("content")
	// AttrContentEditable specifies whether the element is editable.
	AttrContentEditable = makePairAttribute("contenteditable")
	// AttrControls shows the audio/video controls.
	AttrControls = makeBooleanAttribute("controls")
	// AttrCoords specifies the coordinates of an area in an image map.
	AttrCoords = makePairAttribute("coords")
	// AttrCrossOrigin specifies how the element handles cross-origin requests.
	AttrCrossOrigin = makePairAttribute("crossorigin")
	// AttrCsp specifies the Content Security Policy for an element.
	AttrCsp = makePairAttribute("csp")
	// AttrData specifies the URL of the data for an object element.
	AttrData = makePairAttribute("data")
	// AttrDateTime specifies the date and time for an element.
	AttrDateTime = makePairAttribute("datetime")
	// AttrDecoding specifies how to decode an image.
	AttrDecoding = makePairAttribute("decoding")
	// AttrDefault specifies that a track should be enabled by default.
	AttrDefault = makeBooleanAttribute("default")
	// AttrDefer indicates that the script should be executed after the document is parsed.
	AttrDefer = makeBooleanAttribute("defer")
	// AttrDir specifies the text direction of an element.
	AttrDir = makePairAttribute("dir")
	// AttrDirName specifies the name of the form field used for sending the directionality of the element.
	AttrDirName = makePairAttribute("dirname")
	// AttrDisabled specifies that an element should be disabled.
	AttrDisabled = makeBooleanAttribute("disabled")
	// AttrDownload specifies that the target should be downloaded when clicked.
	AttrDownload = makePairAttribute("download")
	// AttrDraggable specifies whether an element is draggable.
	AttrDraggable = makePairAttribute("draggable")
	// AttrEncType specifies how form data should be encoded before sending to a server.
	AttrEncType = makePairAttribute("enctype")
	// AttrEnterKeyHint specifies what action label to show on the enter key.
	AttrEnterKeyHint = makePairAttribute("enterkeyhint")
	// AttrElementTiming specifies that an element should be observed for performance.
	AttrElementTiming = makePairAttribute("elementtiming")
	// AttrForm specifies the id of a form element that the element belongs to.
	AttrForm = makePairAttribute("form")
	// AttrFormAction specifies where to send the form data.
	AttrFormAction = makePairAttribute("formaction")
	// AttrFormEncType specifies how form data should be encoded.
	AttrFormEncType = makePairAttribute("formenctype")
	// AttrFormMethod specifies the HTTP method for form submission.
	AttrFormMethod = makePairAttribute("formmethod")
	// AttrFormNoValidate specifies that the form should not be validated.
	AttrFormNoValidate = makeBooleanAttribute("formnovalidate")
	// AttrFormTarget specifies where to display the response after form submission.
	AttrFormTarget = makePairAttribute("formtarget")
	// AttrFetchPriority indicates the priority of fetching an external resource.
	AttrFetchPriority = makePairAttribute("fetchpriority")
	// AttrHeaders specifies the header cells that a table cell relates to.
	AttrHeaders = makePairAttribute("headers")
	// AttrHeight specifies the height of an element.
	AttrHeight = makePairAttribute("height")
	// AttrHidden specifies that an element is not yet or is no longer relevant.
	AttrHidden = makeBooleanAttribute("hidden")
	// AttrHigh specifies the lower bound of a range.
	AttrHigh = makePairAttribute("high")
	// AttrHref specifies the URL of a link.
	AttrHref = makePairAttribute("href")
	// AttrHrefLang specifies the language of the linked resource.
	AttrHrefLang = makePairAttribute("hreflang")
	// AttrHttpEquiv provides an HTTP header for the information in the content attribute.
	AttrHttpEquiv = makePairAttribute("http-equiv")
	// AttrID specifies a unique id for an element.
	AttrID = makePairAttribute("id")
	// AttrIntegrity specifies a hash of the resource to verify its integrity.
	AttrIntegrity = makePairAttribute("integrity")
	// AttrInputMode provides a hint to browsers about the type of data the user should enter.
	AttrInputMode = makePairAttribute("inputmode")
	// AttrIsMap specifies that an image is part of a server-side image map.
	AttrIsMap = makeBooleanAttribute("ismap")
	// AttrItemProp specifies the property of an item.
	AttrItemProp = makePairAttribute("itemprop")
	// AttrKind specifies the kind of text track.
	AttrKind = makePairAttribute("kind")
	// AttrLabel specifies the label of an option or track.
	AttrLabel = makePairAttribute("label")
	// AttrLang specifies the language of the element.
	AttrLang = makePairAttribute("lang")
	// AttrLanguage specifies the scripting language of an element.
	AttrLanguage = makePairAttribute("language")
	// AttrLoading specifies whether to load an image lazily.
	AttrLoading = makePairAttribute("loading")
	// AttrList refers to a datalist containing predefined options.
	AttrList = makePairAttribute("list")
	// AttrLoop specifies whether to loop an audio/video.
	AttrLoop = makeBooleanAttribute("loop")
	// AttrLow specifies the upper bound of a range.
	AttrLow = makePairAttribute("low")
	// AttrMax specifies the maximum value.
	AttrMax = makePairAttribute("max")
	// AttrMaxLength specifies the maximum number of characters allowed.
	AttrMaxLength = makePairAttribute("maxlength")
	// AttrMinLength specifies the minimum number of characters required.
	AttrMinLength = makePairAttribute("minlength")
	// AttrMedia specifies the media type or device the resource applies to.
	AttrMedia = makePairAttribute("media")
	// AttrMethod specifies the HTTP method for form submission.
	AttrMethod = makePairAttribute("method")
	// AttrMin specifies the minimum value.
	AttrMin = makePairAttribute("min")
	// AttrMultiple specifies that a user can enter more than one value.
	AttrMultiple = makeBooleanAttribute("multiple")
	// AttrMuted specifies that the audio should be muted.
	AttrMuted = makeBooleanAttribute("muted")
	// AttrName specifies the name of an element.
	AttrName = makePairAttribute("name")
	// AttrNoValidate specifies that the form should not be validated.
	AttrNoValidate = makeBooleanAttribute("novalidate")
	// AttrOnAbort specifies the event handler for the abort event.
	AttrOnAbort = makePairAttribute("onAbort")
	// AttrOnActivate specifies the event handler for the activate event.
	AttrOnActivate = makePairAttribute("onActivate")
	// AttrOnAfterPrint specifies the event handler for the afterprint event.
	AttrOnAfterPrint = makePairAttribute("onAfterPrint")
	// AttrOnAfterUpdate specifies the event handler for the afterupdate event.
	AttrOnAfterUpdate = makePairAttribute("onAfterUpdate")
	// AttrOnBeforeActivate specifies the event handler for the beforeactivate event.
	AttrOnBeforeActivate = makePairAttribute("onBeforeActivate")
	// AttrOnBeforeCopy specifies the event handler for the beforecopy event.
	AttrOnBeforeCopy = makePairAttribute("onBeforeCopy")
	// AttrOnBeforeCut specifies the event handler for the beforecut event.
	AttrOnBeforeCut = makePairAttribute("onBeforeCut")
	// AttrOnBeforeDeactivate specifies the event handler for the beforedeactivate event.
	AttrOnBeforeDeactivate = makePairAttribute("onBeforeDeactivate")
	// AttrOnBeforeEditFocus specifies the event handler for the beforeeditfocus event.
	AttrOnBeforeEditFocus = makePairAttribute("onBeforeEditFocus")
	// AttrOnBeforePaste specifies the event handler for the beforepaste event.
	AttrOnBeforePaste = makePairAttribute("onBeforePaste")
	// AttrOnBeforePrint specifies the event handler for the beforeprint event.
	AttrOnBeforePrint = makePairAttribute("onBeforePrint")
	// AttrOnBeforeUnload specifies the event handler for the beforeunload event.
	AttrOnBeforeUnload = makePairAttribute("onBeforeUnload")
	// AttrOnBeforeUpdate specifies the event handler for the beforeupdate event.
	AttrOnBeforeUpdate = makePairAttribute("onBeforeUpdate")
	// AttrOnBegin specifies the event handler for the begin event.
	AttrOnBegin = makePairAttribute("onBegin")
	// AttrOnBlur specifies the event handler for the blur event.
	AttrOnBlur = makePairAttribute("onBlur")
	// AttrOnBounce specifies the event handler for the bounce event.
	AttrOnBounce = makePairAttribute("onBounce")
	// AttrOnCellChange specifies the event handler for the cellchange event.
	AttrOnCellChange = makePairAttribute("onCellChange")
	// AttrOnChange specifies the event handler for the change event.
	AttrOnChange = makePairAttribute("onChange")
	// AttrOnClick specifies the event handler for the click event.
	AttrOnClick = makePairAttribute("onClick")
	// AttrOnContextMenu specifies the event handler for the contextmenu event.
	AttrOnContextMenu = makePairAttribute("onContextMenu")
	// AttrOnControlSelect specifies the event handler for the controlselect event.
	AttrOnControlSelect = makePairAttribute("onControlSelect")
	// AttrOnCopy specifies the event handler for the copy event.
	AttrOnCopy = makePairAttribute("onCopy")
	// AttrOnCut specifies the event handler for the cut event.
	AttrOnCut = makePairAttribute("onCut")
	// AttrOnDataAvailable specifies the event handler for the dataavailable event.
	AttrOnDataAvailable = makePairAttribute("onDataAvailable")
	// AttrOnDataSetChanged specifies the event handler for the datasetchanged event.
	AttrOnDataSetChanged = makePairAttribute("onDataSetChanged")
	// AttrOnDataSetComplete specifies the event handler for the datasetcomplete event.
	AttrOnDataSetComplete = makePairAttribute("onDataSetComplete")
	// AttrOnDblClick specifies the event handler for the dblclick event.
	AttrOnDblClick = makePairAttribute("onDblClick")
	// AttrOnDeactivate specifies the event handler for the deactivate event.
	AttrOnDeactivate = makePairAttribute("onDeactivate")
	// AttrOnDrag specifies the event handler for the drag event.
	AttrOnDrag = makePairAttribute("onDrag")
	// AttrOnDragEnd specifies the event handler for the dragend event.
	AttrOnDragEnd = makePairAttribute("onDragEnd")
	// AttrOnDragLeave specifies the event handler for the dragleave event.
	AttrOnDragLeave = makePairAttribute("onDragLeave")
	// AttrOnDragEnter specifies the event handler for the dragenter event.
	AttrOnDragEnter = makePairAttribute("onDragEnter")
	// AttrOnDragOver specifies the event handler for the dragover event.
	AttrOnDragOver = makePairAttribute("onDragOver")
	// AttrOnDragDrop specifies the event handler for the dragdrop event.
	AttrOnDragDrop = makePairAttribute("onDragDrop")
	// AttrOnDragStart specifies the event handler for the dragstart event.
	AttrOnDragStart = makePairAttribute("onDragStart")
	// AttrOnDrop specifies the event handler for the drop event.
	AttrOnDrop = makePairAttribute("onDrop")
	// AttrOnEnd specifies the event handler for the end event.
	AttrOnEnd = makePairAttribute("onEnd")
	// AttrOnError specifies the event handler for the error event.
	AttrOnError = makePairAttribute("onError")
	// AttrOnErrorUpdate specifies the event handler for the errorupdate event.
	AttrOnErrorUpdate = makePairAttribute("onErrorUpdate")
	// AttrOnFilterChange specifies the event handler for the filterchange event.
	AttrOnFilterChange = makePairAttribute("onFilterChange")
	// AttrOnFinish specifies the event handler for the finish event.
	AttrOnFinish = makePairAttribute("onFinish")
	// AttrOnFocus specifies the event handler for the focus event.
	AttrOnFocus = makePairAttribute("onFocus")
	// AttrOnFocusIn specifies the event handler for the focusin event.
	AttrOnFocusIn = makePairAttribute("onFocusIn")
	// AttrOnFocusOut specifies the event handler for the focusout event.
	AttrOnFocusOut = makePairAttribute("onFocusOut")
	// AttrOnHashChange specifies the event handler for the hashchange event.
	AttrOnHashChange = makePairAttribute("onHashChange")
	// AttrOnHelp specifies the event handler for the help event.
	AttrOnHelp = makePairAttribute("onHelp")
	// AttrOnInput specifies the event handler for the input event.
	AttrOnInput = makePairAttribute("onInput")
	// AttrOnKeyDown specifies the event handler for the keydown event.
	AttrOnKeyDown = makePairAttribute("onKeyDown")
	// AttrOnKeyPress specifies the event handler for the keypress event.
	AttrOnKeyPress = makePairAttribute("onKeyPress")
	// AttrOnKeyUp specifies the event handler for the keyup event.
	AttrOnKeyUp = makePairAttribute("onKeyUp")
	// AttrOnLayoutComplete specifies the event handler for the layoutcomplete event.
	AttrOnLayoutComplete = makePairAttribute("onLayoutComplete")
	// AttrOnLoad specifies the event handler for the load event.
	AttrOnLoad = makePairAttribute("onLoad")
	// AttrOnLoseCapture specifies the event handler for the losecapture event.
	AttrOnLoseCapture = makePairAttribute("onLoseCapture")
	// AttrOnMediaComplete specifies the event handler for the mediacomplete event.
	AttrOnMediaComplete = makePairAttribute("onMediaComplete")
	// AttrOnMediaError specifies the event handler for the mediaerror event.
	AttrOnMediaError = makePairAttribute("onMediaError")
	// AttrOnMessage specifies the event handler for the message event.
	AttrOnMessage = makePairAttribute("onMessage")
	// AttrOnMouseDown specifies the event handler for the mousedown event.
	AttrOnMouseDown = makePairAttribute("onMouseDown")
	// AttrOnMouseEnter specifies the event handler for the mouseenter event.
	AttrOnMouseEnter = makePairAttribute("onMouseEnter")
	// AttrOnMouseLeave specifies the event handler for the mouseleave event.
	AttrOnMouseLeave = makePairAttribute("onMouseLeave")
	// AttrOnMouseMove specifies the event handler for the mousemove event.
	AttrOnMouseMove = makePairAttribute("onMouseMove")
	// AttrOnMouseOut specifies the event handler for the mouseout event.
	AttrOnMouseOut = makePairAttribute("onMouseOut")
	// AttrOnMouseOver specifies the event handler for the mouseover event.
	AttrOnMouseOver = makePairAttribute("onMouseOver")
	// AttrOnMouseUp specifies the event handler for the mouseup event.
	AttrOnMouseUp = makePairAttribute("onMouseUp")
	// AttrOnMouseWheel specifies the event handler for the mousewheel event.
	AttrOnMouseWheel = makePairAttribute("onMouseWheel")
	// AttrOnMove specifies the event handler for the move event.
	AttrOnMove = makePairAttribute("onMove")
	// AttrOnMoveEnd specifies the event handler for the moveend event.
	AttrOnMoveEnd = makePairAttribute("onMoveEnd")
	// AttrOnMoveStart specifies the event handler for the movestart event.
	AttrOnMoveStart = makePairAttribute("onMoveStart")
	// AttrOnOffline specifies the event handler for the offline event.
	AttrOnOffline = makePairAttribute("onOffline")
	// AttrOnOnline specifies the event handler for the online event.
	AttrOnOnline = makePairAttribute("onOnline")
	// AttrOnOutOfSync specifies the event handler for the outofsync event.
	AttrOnOutOfSync = makePairAttribute("onOutOfSync")
	// AttrOnPaste specifies the event handler for the paste event.
	AttrOnPaste = makePairAttribute("onPaste")
	// AttrOnPause specifies the event handler for the pause event.
	AttrOnPause = makePairAttribute("onPause")
	// AttrOnPopState specifies the event handler for the popstate event.
	AttrOnPopState = makePairAttribute("onPopState")
	// AttrOnProgress specifies the event handler for the progress event.
	AttrOnProgress = makePairAttribute("onProgress")
	// AttrOnPropertyChange specifies the event handler for the propertychange event.
	AttrOnPropertyChange = makePairAttribute("onPropertyChange")
	// AttrOnReadyStateChange specifies the event handler for the readystatechange event.
	AttrOnReadyStateChange = makePairAttribute("onReadyStateChange")
	// AttrOnRedo specifies the event handler for the redo event.
	AttrOnRedo = makePairAttribute("onRedo")
	// AttrOnRepeat specifies the event handler for the repeat event.
	AttrOnRepeat = makePairAttribute("onRepeat")
	// AttrOnReset specifies the event handler for the reset event.
	AttrOnReset = makePairAttribute("onReset")
	// AttrOnResize specifies the event handler for the resize event.
	AttrOnResize = makePairAttribute("onResize")
	// AttrOnResizeEnd specifies the event handler for the resizeend event.
	AttrOnResizeEnd = makePairAttribute("onResizeEnd")
	// AttrOnResizeStart specifies the event handler for the resizestart event.
	AttrOnResizeStart = makePairAttribute("onResizeStart")
	// AttrOnResume specifies the event handler for the resume event.
	AttrOnResume = makePairAttribute("onResume")
	// AttrOnReverse specifies the event handler for the reverse event.
	AttrOnReverse = makePairAttribute("onReverse")
	// AttrOnRowsEnter specifies the event handler for the rowsenter event.
	AttrOnRowsEnter = makePairAttribute("onRowsEnter")
	// AttrOnRowExit specifies the event handler for the rowexit event.
	AttrOnRowExit = makePairAttribute("onRowExit")
	// AttrOnRowDelete specifies the event handler for the rowdelete event.
	AttrOnRowDelete = makePairAttribute("onRowDelete")
	// AttrOnRowInserted specifies the event handler for the rowinserted event.
	AttrOnRowInserted = makePairAttribute("onRowInserted")
	// AttrOnScroll specifies the event handler for the scroll event.
	AttrOnScroll = makePairAttribute("onScroll")
	// AttrOnSeek specifies the event handler for the seek event.
	AttrOnSeek = makePairAttribute("onSeek")
	// AttrOnSelect specifies the event handler for the select event.
	AttrOnSelect = makePairAttribute("onSelect")
	// AttrOnSelectionChange specifies the event handler for the selectionchange event.
	AttrOnSelectionChange = makePairAttribute("onSelectionChange")
	// AttrOnSelectStart specifies the event handler for the selectstart event.
	AttrOnSelectStart = makePairAttribute("onSelectStart")
	// AttrOnStart specifies the event handler for the start event.
	AttrOnStart = makePairAttribute("onStart")
	// AttrOnStop specifies the event handler for the stop event.
	AttrOnStop = makePairAttribute("onStop")
	// AttrOnStorage specifies the event handler for the storage event.
	AttrOnStorage = makePairAttribute("onStorage")
	// AttrOnSyncRestored specifies the event handler for the syncrestored event.
	AttrOnSyncRestored = makePairAttribute("onSyncRestored")
	// AttrOnSubmit specifies the event handler for the submit event.
	AttrOnSubmit = makePairAttribute("onSubmit")
	// AttrOnTimeError specifies the event handler for the timeerror event.
	AttrOnTimeError = makePairAttribute("onTimeError")
	// AttrOnTrackChange specifies the event handler for the trackchange event.
	AttrOnTrackChange = makePairAttribute("onTrackChange")
	// AttrOnUndo specifies the event handler for the undo event.
	AttrOnUndo = makePairAttribute("onUndo")
	// AttrOnUnload specifies the event handler for the unload event.
	AttrOnUnload = makePairAttribute("onUnload")
	// AttrOnUrlFlip specifies the event handler for the urlflip event.
	AttrOnUrlFlip = makePairAttribute("onUrlFlip")
	// AttrOpen specifies whether the element is visible (for details, dialog, etc.).
	AttrOpen = makeBooleanAttribute("open")
	// AttrOptimum specifies the optimal value in a range.
	AttrOptimum = makePairAttribute("optimum")
	// AttrPattern specifies a regular expression for input validation.
	AttrPattern = makePairAttribute("pattern")
	// AttrPing specifies a list of URLs to notify when a link is clicked.
	AttrPing = makePairAttribute("ping")
	// AttrPlaceholder provides a hint to the user about what to enter.
	AttrPlaceholder = makePairAttribute("placeholder")
	// AttrPlaysInline specifies that the video should play inline.
	AttrPlaysInline = makeBooleanAttribute("playsinline")
	// AttrPoster specifies the preview image for a video.
	AttrPoster = makePairAttribute("poster")
	// AttrPopoverTargetAction specifies the action to perform with a popover element.
	AttrPopoverTargetAction = makePairAttribute("popovertargetaction")
	// AttrPreload specifies how to preload an audio/video.
	AttrPreload = makePairAttribute("preload")
	// AttrReadOnly specifies that an input field is read-only.
	AttrReadOnly = makeBooleanAttribute("readonly")
	// AttrReferrerPolicy specifies the referrer policy for the resource.
	AttrReferrerPolicy = makePairAttribute("referrerpolicy")
	// AttrRel specifies the relationship between the current document and the linked resource.
	AttrRel = makePairAttribute("rel")
	// AttrRequired specifies that an input field must be filled out.
	AttrRequired = makeBooleanAttribute("required")
	// AttrReversed specifies that the list order should be reversed.
	AttrReversed = makeBooleanAttribute("reversed")
	// AttrRole specifies the role of an element for accessibility.
	AttrRole = makePairAttribute("role")
	// AttrRows specifies the number of rows in a textarea.
	AttrRows = makePairAttribute("rows")
	// AttrRowSpan specifies the number of rows a table cell should span.
	AttrRowSpan = makePairAttribute("rowspan")
	// AttrSandbox enables extra restrictions for an iframe.
	AttrSandbox = makePairAttribute("sandbox")
	// AttrScope specifies the header cells that a th element applies to.
	AttrScope = makePairAttribute("scope")
	// AttrSelected specifies that an option should be pre-selected.
	AttrSelected = makeBooleanAttribute("selected")
	// AttrShape specifies the shape of an area in an image map.
	AttrShape = makePairAttribute("shape")
	// AttrSize specifies the size of an input field or select element.
	AttrSize = makePairAttribute("size")
	// AttrSizes specifies the sizes of an image for different layouts.
	AttrSizes = makePairAttribute("sizes")
	// AttrSlot assigns a slot to an element in a shadow DOM.
	AttrSlot = makePairAttribute("slot")
	// AttrSpan specifies the number of columns in a colgroup.
	AttrSpan = makePairAttribute("span")
	// AttrSpellCheck specifies whether to enable spell checking.
	AttrSpellCheck = makePairAttribute("spellcheck")
	// AttrSrc specifies the URL of an image, audio, video, or iframe.
	AttrSrc = makePairAttribute("src")
	// AttrSrcDoc specifies the inline HTML for an iframe.
	AttrSrcDoc = makePairAttribute("srcdoc")
	// AttrSrcLang specifies the language of the track text.
	AttrSrcLang = makePairAttribute("srclang")
	// AttrSrcSet specifies multiple image sources for responsive images.
	AttrSrcSet = makePairAttribute("srcset")
	// AttrStart specifies the starting number of an ordered list.
	AttrStart = makePairAttribute("start")
	// AttrStep specifies the interval between legal numbers in an input.
	AttrStep = makePairAttribute("step")
	// AttrStyle specifies inline CSS styles.
	AttrStyle = makePairAttribute("style")
	// AttrSummary provides a summary for a table.
	AttrSummary = makePairAttribute("summary")
	// AttrTabIndex specifies the tab order of an element.
	AttrTabIndex = makePairAttribute("tabindex")
	// AttrTarget specifies where to open a link or form response.
	AttrTarget = makePairAttribute("target")
	// AttrTitle provides advisory information about an element.
	AttrTitle = makePairAttribute("title")
	// AttrTranslate specifies whether to translate an element.
	AttrTranslate = makePairAttribute("translate")
	// AttrType specifies the type of an input element.
	AttrType = makePairAttribute("type")
	// AttrUseMap specifies that an image is a client-side image map.
	AttrUseMap = makePairAttribute("usemap")
	// AttrValue specifies the value of an input element.
	AttrValue = makePairAttribute("value")
	// AttrWidth specifies the width of an element.
	AttrWidth = makePairAttribute("width")
	// AttrWrap specifies how text should wrap in a textarea.
	AttrWrap = makePairAttribute("wrap")
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
