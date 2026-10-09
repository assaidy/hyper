package hyper

import (
	"fmt"
	"html"
	"io"
	"reflect"
	"strings"
)

// Attribute represents an HTML attribute that can be rendered.
// Implementations include PairAttribute (key="value") and BooleanAttribute (present/absent).
//
// NOTE: Passing a nil attribute to an element constructor panics.
// Use [IfAttr] for conditional attributes.
type Attribute interface {
	RenderAttribute(w io.Writer) error
}

// PairAttribute represents an HTML attribute with a key and value (key="value").
type PairAttribute struct {
	Key   string
	Value string
}

func (me PairAttribute) RenderAttribute(w io.Writer) error {
	k := strings.TrimSpace(me.Key)
	if k == "" {
		return fmt.Errorf("empty/whitespace attribute key not allowed.")
	}

	io.WriteString(w, " ")
	io.WriteString(w, html.EscapeString(k))
	io.WriteString(w, `="`)
	io.WriteString(w, strings.ReplaceAll(me.Value, `"`, "&quot;"))
	io.WriteString(w, `"`)

	return nil
}

// BooleanAttribute represents an HTML boolean attribute that is either present or absent.
// When IsActive is true, the attribute is rendered; otherwise it is omitted.
type BooleanAttribute struct {
	Key      string
	IsActive bool
}

func (me BooleanAttribute) RenderAttribute(w io.Writer) error {
	k := strings.TrimSpace(me.Key)
	if k == "" {
		return fmt.Errorf("empty/whitespace attribute key not allowed.")
	}

	if me.IsActive {
		io.WriteString(w, " ")
		io.WriteString(w, html.EscapeString(k))
	}

	return nil
}

// Attr creates an attribute from a key and value.
// If value is a string, it creates a PairAttribute (key="value").
// If value is a bool, it creates a BooleanAttribute (present when true, absent when false).
//
// NOTE: Passing a nil attribute to an element constructor panics.
// Use [IfAttr] for conditional attributes.
//
// Examples:
//
//	Attr("class", "container")  // -> PairAttribute{Key: "class", Value: "container"}
//	Attr("hidden", true)        // -> BooleanAttribute{Key: "hidden", IsActive: true}
//	Attr("disabled", false)     // -> BooleanAttribute{Key: "disabled", IsActive: false}
func Attr[V ~string | ~bool](key string, value V) Attribute {
	return attrReflect(key, value)
}

func attrReflect(key string, value any) Attribute {
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.String:
		return PairAttribute{Key: key, Value: reflectValue.String()}
	case reflect.Bool:
		return BooleanAttribute{Key: key, IsActive: reflectValue.Bool()}
	default:
		panic("unexpected value type for attribute")
	}
}

// MakePairAttributeConstructor creates a function that produces a [PairAttribute] with the
// given key. The returned function accepts a string value and returns a
// [PairAttribute]. This is useful for defining custom HTML attributes that take
// a string value.
//
// Pair the result with [IfAttr] to conditionally omit the attribute:
// when no branch matches and no else branch was set, nothing is rendered.
func MakePairAttributeConstructor(key string) func(value string) PairAttribute {
	return func(value string) PairAttribute {
		return PairAttribute{Key: key, Value: value}
	}
}

// MakeBooleanAttributeConstructor creates a function that produces a [BooleanAttribute]
// with the given key. The returned function accepts a bool value and returns
// a [BooleanAttribute]. This is useful for defining custom boolean HTML attributes
// (such as "disabled", "checked", etc.).
//
// Pair the result with [IfAttr] to conditionally omit the attribute:
// when no branch matches and no else branch was set, nothing is rendered.
func MakeBooleanAttributeConstructor(key string) func(isActive bool) BooleanAttribute {
	return func(isActive bool) BooleanAttribute {
		return BooleanAttribute{Key: key, IsActive: isActive}
	}
}

var (
	// AttrAccept sets the accepted file types for <input type="file">.
	AttrAccept = MakePairAttributeConstructor("accept")
	// AttrAcceptCharset sets the character encodings accepted by the server.
	AttrAcceptCharset = MakePairAttributeConstructor("accept-charset")
	// AttrAccessKey gives keyboard shortcut access to an element.
	AttrAccessKey = MakePairAttributeConstructor("accesskey")
	// AttrAction specifies where to send the form data.
	AttrAction = MakePairAttributeConstructor("action")
	// AttrAlign specifies the alignment of an element.
	AttrAlign = MakePairAttributeConstructor("align")
	// AttrAllow specifies permissions for an iframe.
	AttrAllow = MakePairAttributeConstructor("allow")
	// AttrAlpha sets the alpha transparency level of an element.
	AttrAlpha = MakePairAttributeConstructor("alpha")
	// AttrAlt provides alternative text for an image.
	AttrAlt = MakePairAttributeConstructor("alt")
	// AttrAs specifies the relation between the linked resource and the document.
	AttrAs = MakePairAttributeConstructor("as")
	// AttrAsync indicates that the script should execute asynchronously.
	AttrAsync = MakeBooleanAttributeConstructor("async")
	// AttrAutocapitalize controls whether text input is automatically capitalized.
	AttrAutocapitalize = MakePairAttributeConstructor("autocapitalize")
	// AttrAutocomplete specifies whether an input field should have autocomplete enabled.
	AttrAutocomplete = MakePairAttributeConstructor("autocomplete")
	// AttrAutofocus specifies that an element should automatically get focus on page load.
	AttrAutofocus = MakeBooleanAttributeConstructor("autofocus")
	// AttrAutoplay specifies that the audio/video should automatically start playing.
	AttrAutoplay = MakeBooleanAttributeConstructor("autoplay")
	// AttrBackground specifies the background image URL.
	AttrBackground = MakePairAttributeConstructor("background")
	// AttrBgColor specifies the background color of an element.
	AttrBgColor = MakePairAttributeConstructor("bgcolor")
	// AttrBorder specifies the border width around an element.
	AttrBorder = MakePairAttributeConstructor("border")
	// AttrCapture specifies which camera/mic to use for media capture.
	AttrCapture = MakePairAttributeConstructor("capture")
	// AttrCharset specifies the character encoding of the document.
	AttrCharset = MakePairAttributeConstructor("charset")
	// AttrChecked specifies whether an input checkbox or radio is checked.
	AttrChecked = MakeBooleanAttributeConstructor("checked")
	// AttrCite specifies the source of a quotation.
	AttrCite = MakePairAttributeConstructor("cite")
	// AttrClass specifies one or more class names for an element.
	AttrClass = MakePairAttributeConstructor("class")
	// AttrColor specifies the text color of an element.
	AttrColor = MakePairAttributeConstructor("color")
	// AttrColorSpace specifies the color space for an image.
	AttrColorSpace = MakePairAttributeConstructor("colorspace")
	// AttrCols specifies the number of columns in a textarea.
	AttrCols = MakePairAttributeConstructor("cols")
	// AttrColSpan specifies the number of columns a table cell should span.
	AttrColSpan = MakePairAttributeConstructor("colspan")
	// AttrContent provides metadata about the element.
	AttrContent = MakePairAttributeConstructor("content")
	// AttrContentEditable specifies whether the element is editable.
	AttrContentEditable = MakePairAttributeConstructor("contenteditable")
	// AttrControls shows the audio/video controls.
	AttrControls = MakeBooleanAttributeConstructor("controls")
	// AttrCoords specifies the coordinates of an area in an image map.
	AttrCoords = MakePairAttributeConstructor("coords")
	// AttrCrossOrigin specifies how the element handles cross-origin requests.
	AttrCrossOrigin = MakePairAttributeConstructor("crossorigin")
	// AttrCsp specifies the Content Security Policy for an element.
	AttrCsp = MakePairAttributeConstructor("csp")
	// AttrData specifies the URL of the data for an object element.
	AttrData = MakePairAttributeConstructor("data")
	// AttrDateTime specifies the date and time for an element.
	AttrDateTime = MakePairAttributeConstructor("datetime")
	// AttrDecoding specifies how to decode an image.
	AttrDecoding = MakePairAttributeConstructor("decoding")
	// AttrDefault specifies that a track should be enabled by default.
	AttrDefault = MakeBooleanAttributeConstructor("default")
	// AttrDefer indicates that the script should be executed after the document is parsed.
	AttrDefer = MakeBooleanAttributeConstructor("defer")
	// AttrDir specifies the text direction of an element.
	AttrDir = MakePairAttributeConstructor("dir")
	// AttrDirName specifies the name of the form field used for sending the directionality of the element.
	AttrDirName = MakePairAttributeConstructor("dirname")
	// AttrDisabled specifies that an element should be disabled.
	AttrDisabled = MakeBooleanAttributeConstructor("disabled")
	// AttrDownload specifies that the target should be downloaded when clicked.
	AttrDownload = MakePairAttributeConstructor("download")
	// AttrDraggable specifies whether an element is draggable.
	AttrDraggable = MakePairAttributeConstructor("draggable")
	// AttrEncType specifies how form data should be encoded before sending to a server.
	AttrEncType = MakePairAttributeConstructor("enctype")
	// AttrEnterKeyHint specifies what action label to show on the enter key.
	AttrEnterKeyHint = MakePairAttributeConstructor("enterkeyhint")
	// AttrElementTiming specifies that an element should be observed for performance.
	AttrElementTiming = MakePairAttributeConstructor("elementtiming")
	// AttrFor links a label to an input by its ID, improving accessibility and usability.
	AttrFor = MakePairAttributeConstructor("for")
	// AttrForm specifies the id of a form element that the element belongs to.
	AttrForm = MakePairAttributeConstructor("form")
	// AttrFormAction specifies where to send the form data.
	AttrFormAction = MakePairAttributeConstructor("formaction")
	// AttrFormEncType specifies how form data should be encoded.
	AttrFormEncType = MakePairAttributeConstructor("formenctype")
	// AttrFormMethod specifies the HTTP method for form submission.
	AttrFormMethod = MakePairAttributeConstructor("formmethod")
	// AttrFormNoValidate specifies that the form should not be validated.
	AttrFormNoValidate = MakeBooleanAttributeConstructor("formnovalidate")
	// AttrFormTarget specifies where to display the response after form submission.
	AttrFormTarget = MakePairAttributeConstructor("formtarget")
	// AttrFetchPriority indicates the priority of fetching an external resource.
	AttrFetchPriority = MakePairAttributeConstructor("fetchpriority")
	// AttrHeaders specifies the header cells that a table cell relates to.
	AttrHeaders = MakePairAttributeConstructor("headers")
	// AttrHeight specifies the height of an element.
	AttrHeight = MakePairAttributeConstructor("height")
	// AttrHidden specifies that an element is not yet or is no longer relevant.
	AttrHidden = MakeBooleanAttributeConstructor("hidden")
	// AttrHigh specifies the lower bound of a range.
	AttrHigh = MakePairAttributeConstructor("high")
	// AttrHref specifies the URL of a link.
	AttrHref = MakePairAttributeConstructor("href")
	// AttrHrefLang specifies the language of the linked resource.
	AttrHrefLang = MakePairAttributeConstructor("hreflang")
	// AttrHttpEquiv provides an HTTP header for the information in the content attribute.
	AttrHttpEquiv = MakePairAttributeConstructor("http-equiv")
	// AttrId specifies a unique id for an element.
	AttrId = MakePairAttributeConstructor("id")
	// AttrIntegrity specifies a hash of the resource to verify its integrity.
	AttrIntegrity = MakePairAttributeConstructor("integrity")
	// AttrInputMode provides a hint to browsers about the type of data the user should enter.
	AttrInputMode = MakePairAttributeConstructor("inputmode")
	// AttrIsMap specifies that an image is part of a server-side image map.
	AttrIsMap = MakeBooleanAttributeConstructor("ismap")
	// AttrItemProp specifies the property of an item.
	AttrItemProp = MakePairAttributeConstructor("itemprop")
	// AttrKind specifies the kind of text track.
	AttrKind = MakePairAttributeConstructor("kind")
	// AttrLabel specifies the label of an option or track.
	AttrLabel = MakePairAttributeConstructor("label")
	// AttrLang specifies the language of the element.
	AttrLang = MakePairAttributeConstructor("lang")
	// AttrLanguage specifies the scripting language of an element.
	AttrLanguage = MakePairAttributeConstructor("language")
	// AttrLoading specifies whether to load an image lazily.
	AttrLoading = MakePairAttributeConstructor("loading")
	// AttrList refers to a datalist containing predefined options.
	AttrList = MakePairAttributeConstructor("list")
	// AttrLoop specifies whether to loop an audio/video.
	AttrLoop = MakeBooleanAttributeConstructor("loop")
	// AttrLow specifies the upper bound of a range.
	AttrLow = MakePairAttributeConstructor("low")
	// AttrMax specifies the maximum value.
	AttrMax = MakePairAttributeConstructor("max")
	// AttrMaxLength specifies the maximum number of characters allowed.
	AttrMaxLength = MakePairAttributeConstructor("maxlength")
	// AttrMinLength specifies the minimum number of characters required.
	AttrMinLength = MakePairAttributeConstructor("minlength")
	// AttrMedia specifies the media type or device the resource applies to.
	AttrMedia = MakePairAttributeConstructor("media")
	// AttrMethod specifies the HTTP method for form submission.
	AttrMethod = MakePairAttributeConstructor("method")
	// AttrMin specifies the minimum value.
	AttrMin = MakePairAttributeConstructor("min")
	// AttrMultiple specifies that a user can enter more than one value.
	AttrMultiple = MakeBooleanAttributeConstructor("multiple")
	// AttrMuted specifies that the audio should be muted.
	AttrMuted = MakeBooleanAttributeConstructor("muted")
	// AttrName specifies the name of an element.
	AttrName = MakePairAttributeConstructor("name")
	// AttrNoValidate specifies that the form should not be validated.
	AttrNoValidate = MakeBooleanAttributeConstructor("novalidate")
	// AttrOnAbort specifies the event handler for the abort event.
	AttrOnAbort = MakePairAttributeConstructor("onAbort")
	// AttrOnActivate specifies the event handler for the activate event.
	AttrOnActivate = MakePairAttributeConstructor("onActivate")
	// AttrOnAfterPrint specifies the event handler for the afterprint event.
	AttrOnAfterPrint = MakePairAttributeConstructor("onAfterPrint")
	// AttrOnAfterUpdate specifies the event handler for the afterupdate event.
	AttrOnAfterUpdate = MakePairAttributeConstructor("onAfterUpdate")
	// AttrOnBeforeActivate specifies the event handler for the beforeactivate event.
	AttrOnBeforeActivate = MakePairAttributeConstructor("onBeforeActivate")
	// AttrOnBeforeCopy specifies the event handler for the beforecopy event.
	AttrOnBeforeCopy = MakePairAttributeConstructor("onBeforeCopy")
	// AttrOnBeforeCut specifies the event handler for the beforecut event.
	AttrOnBeforeCut = MakePairAttributeConstructor("onBeforeCut")
	// AttrOnBeforeDeactivate specifies the event handler for the beforedeactivate event.
	AttrOnBeforeDeactivate = MakePairAttributeConstructor("onBeforeDeactivate")
	// AttrOnBeforeEditFocus specifies the event handler for the beforeeditfocus event.
	AttrOnBeforeEditFocus = MakePairAttributeConstructor("onBeforeEditFocus")
	// AttrOnBeforePaste specifies the event handler for the beforepaste event.
	AttrOnBeforePaste = MakePairAttributeConstructor("onBeforePaste")
	// AttrOnBeforePrint specifies the event handler for the beforeprint event.
	AttrOnBeforePrint = MakePairAttributeConstructor("onBeforePrint")
	// AttrOnBeforeUnload specifies the event handler for the beforeunload event.
	AttrOnBeforeUnload = MakePairAttributeConstructor("onBeforeUnload")
	// AttrOnBeforeUpdate specifies the event handler for the beforeupdate event.
	AttrOnBeforeUpdate = MakePairAttributeConstructor("onBeforeUpdate")
	// AttrOnBegin specifies the event handler for the begin event.
	AttrOnBegin = MakePairAttributeConstructor("onBegin")
	// AttrOnBlur specifies the event handler for the blur event.
	AttrOnBlur = MakePairAttributeConstructor("onBlur")
	// AttrOnBounce specifies the event handler for the bounce event.
	AttrOnBounce = MakePairAttributeConstructor("onBounce")
	// AttrOnCellChange specifies the event handler for the cellchange event.
	AttrOnCellChange = MakePairAttributeConstructor("onCellChange")
	// AttrOnChange specifies the event handler for the change event.
	AttrOnChange = MakePairAttributeConstructor("onChange")
	// AttrOnClick specifies the event handler for the click event.
	AttrOnClick = MakePairAttributeConstructor("onClick")
	// AttrOnContextMenu specifies the event handler for the contextmenu event.
	AttrOnContextMenu = MakePairAttributeConstructor("onContextMenu")
	// AttrOnControlSelect specifies the event handler for the controlselect event.
	AttrOnControlSelect = MakePairAttributeConstructor("onControlSelect")
	// AttrOnCopy specifies the event handler for the copy event.
	AttrOnCopy = MakePairAttributeConstructor("onCopy")
	// AttrOnCut specifies the event handler for the cut event.
	AttrOnCut = MakePairAttributeConstructor("onCut")
	// AttrOnDataAvailable specifies the event handler for the dataavailable event.
	AttrOnDataAvailable = MakePairAttributeConstructor("onDataAvailable")
	// AttrOnDataSetChanged specifies the event handler for the datasetchanged event.
	AttrOnDataSetChanged = MakePairAttributeConstructor("onDataSetChanged")
	// AttrOnDataSetComplete specifies the event handler for the datasetcomplete event.
	AttrOnDataSetComplete = MakePairAttributeConstructor("onDataSetComplete")
	// AttrOnDblClick specifies the event handler for the dblclick event.
	AttrOnDblClick = MakePairAttributeConstructor("onDblClick")
	// AttrOnDeactivate specifies the event handler for the deactivate event.
	AttrOnDeactivate = MakePairAttributeConstructor("onDeactivate")
	// AttrOnDrag specifies the event handler for the drag event.
	AttrOnDrag = MakePairAttributeConstructor("onDrag")
	// AttrOnDragEnd specifies the event handler for the dragend event.
	AttrOnDragEnd = MakePairAttributeConstructor("onDragEnd")
	// AttrOnDragLeave specifies the event handler for the dragleave event.
	AttrOnDragLeave = MakePairAttributeConstructor("onDragLeave")
	// AttrOnDragEnter specifies the event handler for the dragenter event.
	AttrOnDragEnter = MakePairAttributeConstructor("onDragEnter")
	// AttrOnDragOver specifies the event handler for the dragover event.
	AttrOnDragOver = MakePairAttributeConstructor("onDragOver")
	// AttrOnDragDrop specifies the event handler for the dragdrop event.
	AttrOnDragDrop = MakePairAttributeConstructor("onDragDrop")
	// AttrOnDragStart specifies the event handler for the dragstart event.
	AttrOnDragStart = MakePairAttributeConstructor("onDragStart")
	// AttrOnDrop specifies the event handler for the drop event.
	AttrOnDrop = MakePairAttributeConstructor("onDrop")
	// AttrOnEnd specifies the event handler for the end event.
	AttrOnEnd = MakePairAttributeConstructor("onEnd")
	// AttrOnError specifies the event handler for the error event.
	AttrOnError = MakePairAttributeConstructor("onError")
	// AttrOnErrorUpdate specifies the event handler for the errorupdate event.
	AttrOnErrorUpdate = MakePairAttributeConstructor("onErrorUpdate")
	// AttrOnFilterChange specifies the event handler for the filterchange event.
	AttrOnFilterChange = MakePairAttributeConstructor("onFilterChange")
	// AttrOnFinish specifies the event handler for the finish event.
	AttrOnFinish = MakePairAttributeConstructor("onFinish")
	// AttrOnFocus specifies the event handler for the focus event.
	AttrOnFocus = MakePairAttributeConstructor("onFocus")
	// AttrOnFocusIn specifies the event handler for the focusin event.
	AttrOnFocusIn = MakePairAttributeConstructor("onFocusIn")
	// AttrOnFocusOut specifies the event handler for the focusout event.
	AttrOnFocusOut = MakePairAttributeConstructor("onFocusOut")
	// AttrOnHashChange specifies the event handler for the hashchange event.
	AttrOnHashChange = MakePairAttributeConstructor("onHashChange")
	// AttrOnHelp specifies the event handler for the help event.
	AttrOnHelp = MakePairAttributeConstructor("onHelp")
	// AttrOnInput specifies the event handler for the input event.
	AttrOnInput = MakePairAttributeConstructor("onInput")
	// AttrOnKeyDown specifies the event handler for the keydown event.
	AttrOnKeyDown = MakePairAttributeConstructor("onKeyDown")
	// AttrOnKeyPress specifies the event handler for the keypress event.
	AttrOnKeyPress = MakePairAttributeConstructor("onKeyPress")
	// AttrOnKeyUp specifies the event handler for the keyup event.
	AttrOnKeyUp = MakePairAttributeConstructor("onKeyUp")
	// AttrOnLayoutComplete specifies the event handler for the layoutcomplete event.
	AttrOnLayoutComplete = MakePairAttributeConstructor("onLayoutComplete")
	// AttrOnLoad specifies the event handler for the load event.
	AttrOnLoad = MakePairAttributeConstructor("onLoad")
	// AttrOnLoseCapture specifies the event handler for the losecapture event.
	AttrOnLoseCapture = MakePairAttributeConstructor("onLoseCapture")
	// AttrOnMediaComplete specifies the event handler for the mediacomplete event.
	AttrOnMediaComplete = MakePairAttributeConstructor("onMediaComplete")
	// AttrOnMediaError specifies the event handler for the mediaerror event.
	AttrOnMediaError = MakePairAttributeConstructor("onMediaError")
	// AttrOnMessage specifies the event handler for the message event.
	AttrOnMessage = MakePairAttributeConstructor("onMessage")
	// AttrOnMouseDown specifies the event handler for the mousedown event.
	AttrOnMouseDown = MakePairAttributeConstructor("onMouseDown")
	// AttrOnMouseEnter specifies the event handler for the mouseenter event.
	AttrOnMouseEnter = MakePairAttributeConstructor("onMouseEnter")
	// AttrOnMouseLeave specifies the event handler for the mouseleave event.
	AttrOnMouseLeave = MakePairAttributeConstructor("onMouseLeave")
	// AttrOnMouseMove specifies the event handler for the mousemove event.
	AttrOnMouseMove = MakePairAttributeConstructor("onMouseMove")
	// AttrOnMouseOut specifies the event handler for the mouseout event.
	AttrOnMouseOut = MakePairAttributeConstructor("onMouseOut")
	// AttrOnMouseOver specifies the event handler for the mouseover event.
	AttrOnMouseOver = MakePairAttributeConstructor("onMouseOver")
	// AttrOnMouseUp specifies the event handler for the mouseup event.
	AttrOnMouseUp = MakePairAttributeConstructor("onMouseUp")
	// AttrOnMouseWheel specifies the event handler for the mousewheel event.
	AttrOnMouseWheel = MakePairAttributeConstructor("onMouseWheel")
	// AttrOnMove specifies the event handler for the move event.
	AttrOnMove = MakePairAttributeConstructor("onMove")
	// AttrOnMoveEnd specifies the event handler for the moveend event.
	AttrOnMoveEnd = MakePairAttributeConstructor("onMoveEnd")
	// AttrOnMoveStart specifies the event handler for the movestart event.
	AttrOnMoveStart = MakePairAttributeConstructor("onMoveStart")
	// AttrOnOffline specifies the event handler for the offline event.
	AttrOnOffline = MakePairAttributeConstructor("onOffline")
	// AttrOnOnline specifies the event handler for the online event.
	AttrOnOnline = MakePairAttributeConstructor("onOnline")
	// AttrOnOutOfSync specifies the event handler for the outofsync event.
	AttrOnOutOfSync = MakePairAttributeConstructor("onOutOfSync")
	// AttrOnPaste specifies the event handler for the paste event.
	AttrOnPaste = MakePairAttributeConstructor("onPaste")
	// AttrOnPause specifies the event handler for the pause event.
	AttrOnPause = MakePairAttributeConstructor("onPause")
	// AttrOnPopState specifies the event handler for the popstate event.
	AttrOnPopState = MakePairAttributeConstructor("onPopState")
	// AttrOnProgress specifies the event handler for the progress event.
	AttrOnProgress = MakePairAttributeConstructor("onProgress")
	// AttrOnPropertyChange specifies the event handler for the propertychange event.
	AttrOnPropertyChange = MakePairAttributeConstructor("onPropertyChange")
	// AttrOnReadyStateChange specifies the event handler for the readystatechange event.
	AttrOnReadyStateChange = MakePairAttributeConstructor("onReadyStateChange")
	// AttrOnRedo specifies the event handler for the redo event.
	AttrOnRedo = MakePairAttributeConstructor("onRedo")
	// AttrOnRepeat specifies the event handler for the repeat event.
	AttrOnRepeat = MakePairAttributeConstructor("onRepeat")
	// AttrOnReset specifies the event handler for the reset event.
	AttrOnReset = MakePairAttributeConstructor("onReset")
	// AttrOnResize specifies the event handler for the resize event.
	AttrOnResize = MakePairAttributeConstructor("onResize")
	// AttrOnResizeEnd specifies the event handler for the resizeend event.
	AttrOnResizeEnd = MakePairAttributeConstructor("onResizeEnd")
	// AttrOnResizeStart specifies the event handler for the resizestart event.
	AttrOnResizeStart = MakePairAttributeConstructor("onResizeStart")
	// AttrOnResume specifies the event handler for the resume event.
	AttrOnResume = MakePairAttributeConstructor("onResume")
	// AttrOnReverse specifies the event handler for the reverse event.
	AttrOnReverse = MakePairAttributeConstructor("onReverse")
	// AttrOnRowsEnter specifies the event handler for the rowsenter event.
	AttrOnRowsEnter = MakePairAttributeConstructor("onRowsEnter")
	// AttrOnRowExit specifies the event handler for the rowexit event.
	AttrOnRowExit = MakePairAttributeConstructor("onRowExit")
	// AttrOnRowDelete specifies the event handler for the rowdelete event.
	AttrOnRowDelete = MakePairAttributeConstructor("onRowDelete")
	// AttrOnRowInserted specifies the event handler for the rowinserted event.
	AttrOnRowInserted = MakePairAttributeConstructor("onRowInserted")
	// AttrOnScroll specifies the event handler for the scroll event.
	AttrOnScroll = MakePairAttributeConstructor("onScroll")
	// AttrOnSeek specifies the event handler for the seek event.
	AttrOnSeek = MakePairAttributeConstructor("onSeek")
	// AttrOnSelect specifies the event handler for the select event.
	AttrOnSelect = MakePairAttributeConstructor("onSelect")
	// AttrOnSelectionChange specifies the event handler for the selectionchange event.
	AttrOnSelectionChange = MakePairAttributeConstructor("onSelectionChange")
	// AttrOnSelectStart specifies the event handler for the selectstart event.
	AttrOnSelectStart = MakePairAttributeConstructor("onSelectStart")
	// AttrOnStart specifies the event handler for the start event.
	AttrOnStart = MakePairAttributeConstructor("onStart")
	// AttrOnStop specifies the event handler for the stop event.
	AttrOnStop = MakePairAttributeConstructor("onStop")
	// AttrOnStorage specifies the event handler for the storage event.
	AttrOnStorage = MakePairAttributeConstructor("onStorage")
	// AttrOnSyncRestored specifies the event handler for the syncrestored event.
	AttrOnSyncRestored = MakePairAttributeConstructor("onSyncRestored")
	// AttrOnSubmit specifies the event handler for the submit event.
	AttrOnSubmit = MakePairAttributeConstructor("onSubmit")
	// AttrOnTimeError specifies the event handler for the timeerror event.
	AttrOnTimeError = MakePairAttributeConstructor("onTimeError")
	// AttrOnTrackChange specifies the event handler for the trackchange event.
	AttrOnTrackChange = MakePairAttributeConstructor("onTrackChange")
	// AttrOnUndo specifies the event handler for the undo event.
	AttrOnUndo = MakePairAttributeConstructor("onUndo")
	// AttrOnUnload specifies the event handler for the unload event.
	AttrOnUnload = MakePairAttributeConstructor("onUnload")
	// AttrOnUrlFlip specifies the event handler for the urlflip event.
	AttrOnUrlFlip = MakePairAttributeConstructor("onUrlFlip")
	// AttrOpen specifies whether the element is visible (for details, dialog, etc.).
	AttrOpen = MakeBooleanAttributeConstructor("open")
	// AttrOptimum specifies the optimal value in a range.
	AttrOptimum = MakePairAttributeConstructor("optimum")
	// AttrPattern specifies a regular expression for input validation.
	AttrPattern = MakePairAttributeConstructor("pattern")
	// AttrPing specifies a list of URLs to notify when a link is clicked.
	AttrPing = MakePairAttributeConstructor("ping")
	// AttrPlaceholder provides a hint to the user about what to enter.
	AttrPlaceholder = MakePairAttributeConstructor("placeholder")
	// AttrPlaysInline specifies that the video should play inline.
	AttrPlaysInline = MakeBooleanAttributeConstructor("playsinline")
	// AttrPoster specifies the preview image for a video.
	AttrPoster = MakePairAttributeConstructor("poster")
	// AttrPopoverTargetAction specifies the action to perform with a popover element.
	AttrPopoverTargetAction = MakePairAttributeConstructor("popovertargetaction")
	// AttrPreload specifies how to preload an audio/video.
	AttrPreload = MakePairAttributeConstructor("preload")
	// AttrReadOnly specifies that an input field is read-only.
	AttrReadOnly = MakeBooleanAttributeConstructor("readonly")
	// AttrReferrerPolicy specifies the referrer policy for the resource.
	AttrReferrerPolicy = MakePairAttributeConstructor("referrerpolicy")
	// AttrRel specifies the relationship between the current document and the linked resource.
	AttrRel = MakePairAttributeConstructor("rel")
	// AttrRequired specifies that an input field must be filled out.
	AttrRequired = MakeBooleanAttributeConstructor("required")
	// AttrReversed specifies that the list order should be reversed.
	AttrReversed = MakeBooleanAttributeConstructor("reversed")
	// AttrRole specifies the role of an element for accessibility.
	AttrRole = MakePairAttributeConstructor("role")
	// AttrRows specifies the number of rows in a textarea.
	AttrRows = MakePairAttributeConstructor("rows")
	// AttrRowSpan specifies the number of rows a table cell should span.
	AttrRowSpan = MakePairAttributeConstructor("rowspan")
	// AttrSandbox enables extra restrictions for an iframe.
	AttrSandbox = MakePairAttributeConstructor("sandbox")
	// AttrScope specifies the header cells that a th element applies to.
	AttrScope = MakePairAttributeConstructor("scope")
	// AttrSelected specifies that an option should be pre-selected.
	AttrSelected = MakeBooleanAttributeConstructor("selected")
	// AttrShape specifies the shape of an area in an image map.
	AttrShape = MakePairAttributeConstructor("shape")
	// AttrSize specifies the size of an input field or select element.
	AttrSize = MakePairAttributeConstructor("size")
	// AttrSizes specifies the sizes of an image for different layouts.
	AttrSizes = MakePairAttributeConstructor("sizes")
	// AttrSlot assigns a slot to an element in a shadow DOM.
	AttrSlot = MakePairAttributeConstructor("slot")
	// AttrSpan specifies the number of columns in a colgroup.
	AttrSpan = MakePairAttributeConstructor("span")
	// AttrSpellCheck specifies whether to enable spell checking.
	AttrSpellCheck = MakePairAttributeConstructor("spellcheck")
	// AttrSrc specifies the URL of an image, audio, video, or iframe.
	AttrSrc = MakePairAttributeConstructor("src")
	// AttrSrcDoc specifies the inline HTML for an iframe.
	AttrSrcDoc = MakePairAttributeConstructor("srcdoc")
	// AttrSrcLang specifies the language of the track text.
	AttrSrcLang = MakePairAttributeConstructor("srclang")
	// AttrSrcSet specifies multiple image sources for responsive images.
	AttrSrcSet = MakePairAttributeConstructor("srcset")
	// AttrStart specifies the starting number of an ordered list.
	AttrStart = MakePairAttributeConstructor("start")
	// AttrStep specifies the interval between legal numbers in an input.
	AttrStep = MakePairAttributeConstructor("step")
	// AttrStyle specifies inline CSS styles.
	AttrStyle = MakePairAttributeConstructor("style")
	// AttrSummary provides a summary for a table.
	AttrSummary = MakePairAttributeConstructor("summary")
	// AttrTabIndex specifies the tab order of an element.
	AttrTabIndex = MakePairAttributeConstructor("tabindex")
	// AttrTarget specifies where to open a link or form response.
	AttrTarget = MakePairAttributeConstructor("target")
	// AttrTitle provides advisory information about an element.
	AttrTitle = MakePairAttributeConstructor("title")
	// AttrTranslate specifies whether to translate an element.
	AttrTranslate = MakePairAttributeConstructor("translate")
	// AttrType specifies the type of an input element.
	AttrType = MakePairAttributeConstructor("type")
	// AttrUseMap specifies that an image is a client-side image map.
	AttrUseMap = MakePairAttributeConstructor("usemap")
	// AttrValue specifies the value of an input element.
	AttrValue = MakePairAttributeConstructor("value")
	// AttrWidth specifies the width of an element.
	AttrWidth = MakePairAttributeConstructor("width")
	// AttrWrap specifies how text should wrap in a textarea.
	AttrWrap = MakePairAttributeConstructor("wrap")
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
	// MethodDialog indicates that the form is part of a dialog and should use the dialog's method handling.
	MethodDialog = "dialog"
)

// EncType* constants are valid values for the enctype attribute on <form>.
const (
	// EncTypeUrlEncoded encodes form data as URL-encoded string.
	EncTypeUrlEncoded = "application/x-www-form-urlencoded"
	// EncTypeMultipartForm encodes form data as multipart/form-data.
	EncTypeMultipartForm = "multipart/form-data"
	// EncTypePlainText encodes form data as plain text.
	EncTypePlainText = "text/plain"
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
	// ContentEditablePlainTextOnly allows only plain text editing.
	ContentEditablePlainTextOnly = "plaintext-only"
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
	// AutocompleteHonorificPrefix specifies an honorific prefix (e.g., Mr, Mrs).
	AutocompleteHonorificPrefix = "honorific-prefix"
	// AutocompleteGivenName specifies the given (first) name.
	AutocompleteGivenName = "given-name"
	// AutocompleteAdditionalName specifies an additional name.
	AutocompleteAdditionalName = "additional-name"
	// AutocompleteFamilyName specifies the family (last) name.
	AutocompleteFamilyName = "family-name"
	// AutocompleteHonorificSuffix specifies an honorific suffix (e.g., Jr, III).
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
	// AutocompleteCcName specifies the name on the credit card.
	AutocompleteCcName = "cc-name"
	// AutocompleteCcGivenName specifies the given name on the credit card.
	AutocompleteCcGivenName = "cc-given-name"
	// AutocompleteCcAdditionalName specifies the additional name on the credit card.
	AutocompleteCcAdditionalName = "cc-additional-name"
	// AutocompleteCcFamilyName specifies the family name on the credit card.
	AutocompleteCcFamilyName = "cc-family-name"
	// AutocompleteCcNumber specifies the credit card number.
	AutocompleteCcNumber = "cc-number"
	// AutocompleteCcExp specifies the credit card expiration date.
	AutocompleteCcExp = "cc-exp"
	// AutocompleteCcExpMonth specifies the credit card expiration month.
	AutocompleteCcExpMonth = "cc-exp-month"
	// AutocompleteCcExpYear specifies the credit card expiration year.
	AutocompleteCcExpYear = "cc-exp-year"
	// AutocompleteCcCsc specifies the credit card security code.
	AutocompleteCcCsc = "cc-csc"
	// AutocompleteCcType specifies the credit card type (e.g., Visa, Mastercard).
	AutocompleteCcType = "cc-type"
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
