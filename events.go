package h

const (
	// EventAbort fires when an operation is aborted.
	EventAbort = "abort"
	// EventAfterPrint fires after the document is printed.
	EventAfterPrint = "afterprint"
	// EventAnimationEnd fires when a CSS animation ends.
	EventAnimationEnd = "animationend"
	// EventAnimationIteration fires when a CSS animation repeats.
	EventAnimationIteration = "animationiteration"
	// EventAnimationStart fires when a CSS animation starts.
	EventAnimationStart = "animationstart"
	// EventAppInstalled fires when an app is installed.
	EventAppInstalled = "appinstalled"
	// EventAudioProcess fires when an audio buffer is processed.
	EventAudioProcess = "audioprocess"
	// EventAudioEnd fires when audio playback ends.
	EventAudioEnd = "audioend"
	// EventAudioStart fires when audio playback starts.
	EventAudioStart = "audiostart"
	// EventBeforePrint fires before the document is printed.
	EventBeforePrint = "beforeprint"
	// EventBeforeUnload fires before the document is unloaded.
	EventBeforeUnload = "beforeunload"
	// EventBeginEvent fires when a SMIL animation begins.
	EventBeginEvent = "beginEvent"
	// EventBlur fires when an element loses focus.
	EventBlur = "blur"
	// EventBoundary fires when a speech recognition boundary is reached.
	EventBoundary = "boundary"
	// EventCached fires when a cached resource is available.
	EventCached = "cached"
	// EventCanPlay fires when media can start playing.
	EventCanPlay = "canplay"
	// EventCanPlayThrough fires when media can play to the end without buffering.
	EventCanPlayThrough = "canplaythrough"
	// EventChange fires when an input value changes.
	EventChange = "change"
	// EventChargingChange fires when the battery charging state changes.
	EventChargingChange = "chargingchange"
	// EventChargingTimeChange fires when the battery charging time changes.
	EventChargingTimeChange = "chargingtimechange"
	// EventChecking fires when checking for cached resources.
	EventChecking = "checking"
	// EventClick fires when an element is clicked.
	EventClick = "click"
	// EventClose fires when a connection is closed.
	EventClose = "close"
	// EventComplete fires when an operation completes.
	EventComplete = "complete"
	// EventCompositionEnd fires when composition of text is complete.
	EventCompositionEnd = "compositionend"
	// EventCompositionStart fires when composition of text starts.
	EventCompositionStart = "compositionstart"
	// EventCompositionUpdate fires during composition of text.
	EventCompositionUpdate = "compositionupdate"
	// EventContextMenu fires when a context menu is triggered.
	EventContextMenu = "contextmenu"
	// EventCopy fires when text is copied to the clipboard.
	EventCopy = "copy"
	// EventCut fires when text is cut from the clipboard.
	EventCut = "cut"
	// EventDblClick fires when an element is double-clicked.
	EventDblClick = "dblclick"
	// EventDeviceChange fires when a device changes.
	EventDeviceChange = "devicechange"
	// EventDeviceLight fires when ambient light levels change.
	EventDeviceLight = "devicelight"
	// EventDeviceMotion fires when device motion data is available.
	EventDeviceMotion = "devicemotion"
	// EventDeviceOrientation fires when device orientation changes.
	EventDeviceOrientation = "deviceorientation"
	// EventDeviceProximity fires when device proximity is detected.
	EventDeviceProximity = "deviceproximity"
	// EventDischargingTimeChange fires when the battery discharging time changes.
	EventDischargingTimeChange = "dischargingtimechange"
	// EventDomActivate fires when an element is activated.
	EventDomActivate = "DOMActivate"
	// EventDomAttributeNameChanged fires when an attribute name changes.
	EventDomAttributeNameChanged = "DOMAttributeNameChanged"
	// EventDomAttrModified fires when an attribute is modified.
	EventDomAttrModified = "DOMAttrModified"
	// EventDomCharacterDataModified fires when character data is modified.
	EventDomCharacterDataModified = "DOMCharacterDataModified"
	// EventDomContentLoaded fires when the DOM content is loaded.
	EventDomContentLoaded = "DOMContentLoaded"
	// EventDomElementNameChanged fires when an element name changes.
	EventDomElementNameChanged = "DOMElementNameChanged"
	// EventDomFocusIn fires when an element gains focus.
	EventDomFocusIn = "DOMFocusIn"
	// EventDomFocusOut fires when an element loses focus.
	EventDomFocusOut = "DOMFocusOut"
	// EventDomNodeInserted fires when a node is inserted.
	EventDomNodeInserted = "DOMNodeInserted"
	// EventDomNodeInsertedIntoDocument fires when a node is inserted into the document.
	EventDomNodeInsertedIntoDocument = "DOMNodeInsertedIntoDocument"
	// EventDomNodeRemoved fires when a node is removed.
	EventDomNodeRemoved = "DOMNodeRemoved"
	// EventDomNodeRemovedFromDocument fires when a node is removed from the document.
	EventDomNodeRemovedFromDocument = "DOMNodeRemovedFromDocument"
	// EventDomSubtreeModified fires when the subtree is modified.
	EventDomSubtreeModified = "DOMSubtreeModified"
	// EventDownloading fires when a download is in progress.
	EventDownloading = "downloading"
	// EventDrag fires when an element is being dragged.
	EventDrag = "drag"
	// EventDragEnd fires when a drag operation ends.
	EventDragEnd = "dragend"
	// EventDragEnter fires when a dragged element enters a drop target.
	EventDragEnter = "dragenter"
	// EventDragLeave fires when a dragged element leaves a drop target.
	EventDragLeave = "dragleave"
	// EventDragOver fires when a dragged element is over a drop target.
	EventDragOver = "dragover"
	// EventDragStart fires when a drag operation starts.
	EventDragStart = "dragstart"
	// EventDrop fires when an element is dropped.
	EventDrop = "drop"
	// EventDurationChange fires when the duration of media changes.
	EventDurationChange = "durationchange"
	// EventEmptied fires when media becomes empty.
	EventEmptied = "emptied"
	// EventEnd fires when an operation ends.
	EventEnd = "end"
	// EventEnded fires when media playback ends.
	EventEnded = "ended"
	// EventEndEvent fires when a SMIL animation ends.
	EventEndEvent = "endEvent"
	// EventError fires when an error occurs.
	EventError = "error"
	// EventFocus fires when an element gains focus.
	EventFocus = "focus"
	// EventFocusIn fires when an element gains focus.
	EventFocusIn = "focusin"
	// EventFocusOut fires when an element loses focus.
	EventFocusOut = "focusout"
	// EventFullscreenChange fires when fullscreen mode changes.
	EventFullscreenChange = "fullscreenchange"
	// EventFullscreenError fires when fullscreen mode cannot be entered.
	EventFullscreenError = "fullscreenerror"
	// EventGamepadConnected fires when a gamepad is connected.
	EventGamepadConnected = "gamepadconnected"
	// EventGamepadDisconnected fires when a gamepad is disconnected.
	EventGamepadDisconnected = "gamepaddisconnected"
	// EventGotPointerCapture fires when pointer capture is obtained.
	EventGotPointerCapture = "gotpointercapture"
	// EventHashChange fires when the URL hash changes.
	EventHashChange = "hashchange"
	// EventLostPointerCapture fires when pointer capture is lost.
	EventLostPointerCapture = "lostpointercapture"
	// EventInput fires when an input value changes.
	EventInput = "input"
	// EventInvalid fires when an input element is invalid.
	EventInvalid = "invalid"
	// EventKeyDown fires when a key is pressed down.
	EventKeyDown = "keydown"
	// EventKeyPress fires when a key is pressed.
	EventKeyPress = "keypress"
	// EventKeyUp fires when a key is released.
	EventKeyUp = "keyup"
	// EventLanguageChange fires when the language changes.
	EventLanguageChange = "languagechange"
	// EventLevelChange fires when the battery level changes.
	EventLevelChange = "levelchange"
	// EventLoad fires when a resource loads.
	EventLoad = "load"
	// EventLoadedData fires when media data is loaded.
	EventLoadedData = "loadeddata"
	// EventLoadedMetadata fires when media metadata is loaded.
	EventLoadedMetadata = "loadedmetadata"
	// EventLoadEnd fires when a resource finishes loading.
	EventLoadEnd = "loadend"
	// EventLoadStart fires when a resource starts loading.
	EventLoadStart = "loadstart"
	// EventMark fires when a text mark is reached in media.
	EventMark = "mark"
	// EventMessage fires when a message is received.
	EventMessage = "message"
	// EventMessageError fires when a message error occurs.
	EventMessageError = "messageerror"
	// EventMouseDown fires when a mouse button is pressed.
	EventMouseDown = "mousedown"
	// EventMouseEnter fires when the mouse enters an element.
	EventMouseEnter = "mouseenter"
	// EventMouseLeave fires when the mouse leaves an element.
	EventMouseLeave = "mouseleave"
	// EventMouseMove fires when the mouse moves.
	EventMouseMove = "mousemove"
	// EventMouseOut fires when the mouse leaves an element.
	EventMouseOut = "mouseout"
	// EventMouseOver fires when the mouse enters an element.
	EventMouseOver = "mouseover"
	// EventMouseUp fires when a mouse button is released.
	EventMouseUp = "mouseup"
	// EventNoMatch fires when a speech recognition result doesn't match.
	EventNoMatch = "nomatch"
	// EventNotificationClick fires when a notification is clicked.
	EventNotificationClick = "notificationclick"
	// EventNoUpdate fires when no update is available.
	EventNoUpdate = "noupdate"
	// EventObsolete fires when an obsolete API is used.
	EventObsolete = "obsolete"
	// EventOffline fires when the network goes offline.
	EventOffline = "offline"
	// EventOnline fires when the network comes online.
	EventOnline = "online"
	// EventOpen fires when a connection is opened.
	EventOpen = "open"
	// EventOrientationChange fires when the device orientation changes.
	EventOrientationChange = "orientationchange"
	// EventPageHide fires when a page is being hidden.
	EventPageHide = "pagehide"
	// EventPageShow fires when a page is being shown.
	EventPageShow = "pageshow"
	// EventPaste fires when text is pasted from the clipboard.
	EventPaste = "paste"
	// EventPause fires when media playback pauses.
	EventPause = "pause"
	// EventPointerCancel fires when a pointer operation is cancelled.
	EventPointerCancel = "pointercancel"
	// EventPointerDown fires when a pointer is pressed.
	EventPointerDown = "pointerdown"
	// EventPointerEnter fires when a pointer enters an element.
	EventPointerEnter = "pointerenter"
	// EventPointerLeave fires when a pointer leaves an element.
	EventPointerLeave = "pointerleave"
	// EventPointerLockChange fires when pointer lock state changes.
	EventPointerLockChange = "pointerlockchange"
	// EventPointerLockError fires when pointer lock fails.
	EventPointerLockError = "pointerlockerror"
	// EventPointerMove fires when a pointer moves.
	EventPointerMove = "pointermove"
	// EventPointerOut fires when a pointer leaves an element.
	EventPointerOut = "pointerout"
	// EventPointerOver fires when a pointer enters an element.
	EventPointerOver = "pointerover"
	// EventPointerUp fires when a pointer is released.
	EventPointerUp = "pointerup"
	// EventPlay fires when media playback starts.
	EventPlay = "play"
	// EventPlaying fires when media is actively playing.
	EventPlaying = "playing"
	// EventPopState fires when the history state changes.
	EventPopState = "popstate"
	// EventProgress fires when a resource is loading.
	EventProgress = "progress"
	// EventPush fires when a push notification is received.
	EventPush = "push"
	// EventPushSubscriptionChange fires when a push subscription changes.
	EventPushSubscriptionChange = "pushsubscriptionchange"
	// EventRateChange fires when the playback rate changes.
	EventRateChange = "ratechange"
	// EventReadyStateChange fires when the ready state changes.
	EventReadyStateChange = "readystatechange"
	// EventRepeatEvent fires when a SMIL animation repeats.
	EventRepeatEvent = "repeatEvent"
	// EventReset fires when a form is reset.
	EventReset = "reset"
	// EventResize fires when the viewport is resized.
	EventResize = "resize"
	// EventResourceTimingBufferFull fires when the resource timing buffer is full.
	EventResourceTimingBufferFull = "resourcetimingbufferfull"
	// EventResult fires when a speech recognition result is available.
	EventResult = "result"
	// EventResume fires when media playback resumes.
	EventResume = "resume"
	// EventScroll fires when an element is scrolled.
	EventScroll = "scroll"
	// EventSeeked fires when a seek operation completes.
	EventSeeked = "seeked"
	// EventSeeking fires when a seek operation starts.
	EventSeeking = "seeking"
	// EventSelect fires when text is selected.
	EventSelect = "select"
	// EventSelectStart fires when a text selection starts.
	EventSelectStart = "selectstart"
	// EventSelectionChange fires when the selection changes.
	EventSelectionChange = "selectionchange"
	// EventShow fires when a context menu or dialog is shown.
	EventShow = "show"
	// EventSlotChange fires when a slot's content changes.
	EventSlotChange = "slotchange"
	// EventSoundEnd fires when sound ends.
	EventSoundEnd = "soundend"
	// EventSoundStart fires when sound starts.
	EventSoundStart = "soundstart"
	// EventSpeechEnd fires when speech recognition ends.
	EventSpeechEnd = "speechend"
	// EventSpeechStart fires when speech recognition starts.
	EventSpeechStart = "speechstart"
	// EventStalled fires when media loading stalls.
	EventStalled = "stalled"
	// EventStart fires when an operation starts.
	EventStart = "start"
	// EventStorage fires when storage is modified.
	EventStorage = "storage"
	// EventSubmit fires when a form is submitted.
	EventSubmit = "submit"
	// EventSuccess fires when an operation succeeds.
	EventSuccess = "success"
	// EventSuspend fires when media loading is suspended.
	EventSuspend = "suspend"
	// EventSvgAbort fires when SVG loading is aborted.
	EventSvgAbort = "SVGAbort"
	// EventSvgError fires when SVG loading fails.
	EventSvgError = "SVGError"
	// EventSvgLoad fires when SVG loads.
	EventSvgLoad = "SVGLoad"
	// EventSvgResize fires when SVG is resized.
	EventSvgResize = "SVGResize"
	// EventSvgScroll fires when SVG is scrolled.
	EventSvgScroll = "SVGScroll"
	// EventSvgUnload fires when SVG is unloaded.
	EventSvgUnload = "SVGUnload"
	// EventSvgZoom fires when SVG is zoomed.
	EventSvgZoom = "SVGZoom"
	// EventTimeout fires when a timeout occurs.
	EventTimeout = "timeout"
	// EventTimeUpdate fires when the current playback time updates.
	EventTimeUpdate = "timeupdate"
	// EventTouchCancel fires when a touch is cancelled.
	EventTouchCancel = "touchcancel"
	// EventTouchEnd fires when a touch ends.
	EventTouchEnd = "touchend"
	// EventTouchMove fires when a touch moves.
	EventTouchMove = "touchmove"
	// EventTouchStart fires when a touch starts.
	EventTouchStart = "touchstart"
	// EventTransitionEnd fires when a CSS transition ends.
	EventTransitionEnd = "transitionend"
	// EventUnload fires when the document is unloaded.
	EventUnload = "unload"
	// EventUpdateReady fires when an update is ready.
	EventUpdateReady = "updateready"
	// EventUserProximity fires when user proximity is detected.
	EventUserProximity = "userproximity"
	// EventVoicesChanged fires when available voices change.
	EventVoicesChanged = "voiceschanged"
	// EventVisibilityChange fires when document visibility changes.
	EventVisibilityChange = "visibilitychange"
	// EventVolumeChange fires when the volume changes.
	EventVolumeChange = "volumechange"
	// EventWaiting fires when media is waiting for data.
	EventWaiting = "waiting"
	// EventWheel fires when the mouse wheel is rotated.
	EventWheel = "wheel"
)
