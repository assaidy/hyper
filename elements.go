package hyper

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"math/bits"
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

// Element represents an HTML element with its attributes and children.
type Element struct {
	Name       string      // HTML tag name
	IsVoid     bool        // Whether the tag is self-closing (e.g., <br>, <img>)
	Attributes []Attribute // HTML attributes as [PairAttribute] or [BooleanAttribute]
	Children   []HyperNode // Child nodes
}

// Render generates the HTML for the element and its children to the provided writer.
func (me Element) Render(w io.Writer) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	if err := me.render(buf); err != nil {
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

// render renders the element to the provided buffer.
func (me Element) render(buf *bytes.Buffer) error {
	if me.Name == "" {
		return me.renderChildren(buf)
	}

	buf.WriteByte('<')
	buf.WriteString(me.Name)
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
	buf.WriteString(me.Name)
	buf.WriteByte('>')
	return nil
}

// renderChildren renders all child nodes to the provided buffer.
func (me Element) renderChildren(buf *bytes.Buffer) error {
	for _, child := range me.Children {
		switch c := child.(type) {
		// I'm trying to pass the concrete type [bytes.Buffer] as possible.
		// That's why I'm not just using Render(buf), as in the default case,
		// which accepts io.Writer.
		case Element:
			if err := c.render(buf); err != nil {
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

// renderChildren renders all attributes to the provided buffer.
func (me Element) renderAttrs(buf *bytes.Buffer) error {
	for _, attr := range me.Attributes {
		if attr != nil {
			if err := attr.Render(buf); err != nil {
				return err
			}
		}
	}

	return nil
}

// InsertChildren adds child nodes to an [Element]. It accepts [HyperNode] values,
// strings (converted to [Text]), and other values (converted to [Text] via fmt.Sprint).
func (me *Element) InsertChildren(children ...any) {
	n := len(children)
	if n == 0 {
		return
	}

	oldLen := len(me.Children)
	newLen := oldLen + n

	if newLen > cap(me.Children) {
		// Capacity grows exponentially.
		newCap := 1 << bits.Len(uint(newLen-1))
		newSlice := make([]HyperNode, oldLen, newCap)
		copy(newSlice, me.Children)
		me.Children = newSlice
	}

	for _, child := range children {
		me.Children = append(me.Children, toHyperNode(child))
	}
}

// InsertAttributes appends attributes to an [Element].
func (me *Element) InsertAttributes(attrs ...Attribute) {
	n := len(attrs)
	if n == 0 {
		return
	}

	oldLen := len(me.Attributes)
	newLen := oldLen + n

	if newLen > cap(me.Attributes) {
		// Capacity grows exponentially.
		newCap := 1 << bits.Len(uint(newLen-1))
		newSlice := make([]Attribute, oldLen, newCap)
		copy(newSlice, me.Attributes)
		me.Attributes = newSlice
	}

	me.Attributes = append(me.Attributes, attrs...)
}

// toHyperNode converts an arbitrary value to a [HyperNode].
// Strings become [Text], fmt.Stringer values are converted via String(),
// and all other values use fmt.Sprint.
func toHyperNode(v any) HyperNode {
	switch value := v.(type) {
	case HyperNode:
		return value

	// Explicit string and fmt.Stringer cases for performance:
	// fmt.Sprint() would handle these, but with overhead from type inspection and buffer allocation.
	case string:
		return Text(value)
	case fmt.Stringer:
		return Text(value.String())

	// Nil arguments are not checked/filtered - fmt.Sprint() renders them as "<nil>",
	// which is intentional for better debugging (makes it obvious when nil values are passed).
	default:
		return Text(fmt.Sprint(value))
	}
}

// NewElement creates an [Element] with the given tag name from mixed arguments.
// [Attribute] values become element attributes; all other values become
// children via [InsertChildren].
func NewElement(name string, args ...any) Element {
	e := Element{Name: name}
	for _, a := range args {
		if attr, ok := a.(Attribute); ok {
			e.InsertAttributes(attr)
		} else {
			e.InsertChildren(a)
		}
	}
	return e
}

// NewVoidElement creates a void (self-closing) [Element] with the given tag
// name and attributes. Void elements cannot have children.
func NewVoidElement(name string, attrs ...Attribute) Element {
	return Element{Name: name, Attributes: attrs, IsVoid: true}
}

// DOCTYPE creates the <!DOCTYPE html> element.
//
// https://developer.mozilla.org/en-US/docs/Glossary/Doctype
func DOCTYPE() HyperNode {
	return NewVoidElement("!DOCTYPE html")
}

// HTML creates the root element of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
func HTML(args ...any) Element {
	return NewElement("html", args...)
}

// HEAD contains machine-readable information about the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
func HEAD(args ...any) Element {
	return NewElement("head", args...)
}

// TITLE defines the document's title that is shown in a browser's title bar or a page's tab.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
func TITLE(args ...any) Element {
	return NewElement("title", args...)
}

// LINK specifies relationships between the current document and an external resource.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
func LINK(attrs ...Attribute) Element {
	return NewVoidElement("link", attrs...)
}

// META represents metadata that cannot be represented by other HTML meta-related elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
func META(attrs ...Attribute) Element {
	return NewVoidElement("meta", attrs...)
}

// STYLE contains style information for a document or part of a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
func STYLE(args ...any) Element {
	return NewElement("style", args...)
}

// BODY represents the content of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
func BODY(args ...any) Element {
	return NewElement("body", args...)
}

// H1 creates a level 1 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h1
func H1(args ...any) Element {
	return NewElement("h1", args...)
}

// H2 creates a level 2 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h2
func H2(args ...any) Element {
	return NewElement("h2", args...)
}

// H3 creates a level 3 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h3
func H3(args ...any) Element {
	return NewElement("h3", args...)
}

// H4 creates a level 4 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h4
func H4(args ...any) Element {
	return NewElement("h4", args...)
}

// H5 creates a level 5 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h5
func H5(args ...any) Element {
	return NewElement("h5", args...)
}

// H6 creates a level 6 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h6
func H6(args ...any) Element {
	return NewElement("h6", args...)
}

// HEADER creates a header element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
func HEADER(args ...any) Element {
	return NewElement("header", args...)
}

// FOOTER creates a footer element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
func FOOTER(args ...any) Element {
	return NewElement("footer", args...)
}

// NAV creates a navigation element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/nav
func NAV(args ...any) Element {
	return NewElement("nav", args...)
}

// MAIN creates a main content element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/main
func MAIN(args ...any) Element {
	return NewElement("main", args...)
}

// SECTION creates a section element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/section
func SECTION(args ...any) Element {
	return NewElement("section", args...)
}

// ARTICLE creates an article element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
func ARTICLE(args ...any) Element {
	return NewElement("article", args...)
}

// ASIDE creates an aside element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
func ASIDE(args ...any) Element {
	return NewElement("aside", args...)
}

// HR represents a thematic break between paragraph-level elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hr
func HR(attrs ...Attribute) Element {
	return NewVoidElement("hr", attrs...)
}

// PRE represents preformatted text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/pre
func PRE(args ...any) Element {
	return NewElement("pre", args...)
}

// BLOCKQUOTE represents a section quoted from another source.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote
func BLOCKQUOTE(args ...any) Element {
	return NewElement("blockquote", args...)
}

// OL represents an ordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ol
func OL(args ...any) Element {
	return NewElement("ol", args...)
}

// UL represents an unordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ul
func UL(args ...any) Element {
	return NewElement("ul", args...)
}

// LI represents a list item.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/li
func LI(args ...any) Element {
	return NewElement("li", args...)
}

// A creates hyperlinks to other web pages, files, locations within the same page, or anything else a URL can address.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a
func A(args ...any) Element {
	return NewElement("a", args...)
}

// EM marks text with emphasis.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/em
func EM(args ...any) Element {
	return NewElement("em", args...)
}

// STRONG indicates strong importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/strong
func STRONG(args ...any) Element {
	return NewElement("strong", args...)
}

// CODE displays its contents styled as computer code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/code
func CODE(args ...any) Element {
	return NewElement("code", args...)
}

// VAR represents a variable in a mathematical expression or programming context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/var
func VAR(args ...any) Element {
	return NewElement("var", args...)
}

// SAMP represents sample output from a computer program.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/samp
func SAMP(args ...any) Element {
	return NewElement("samp", args...)
}

// KBD represents text that the user should enter.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/kbd
func KBD(args ...any) Element {
	return NewElement("kbd", args...)
}

// SUB specifies inline text displayed as subscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sub
func SUB(args ...any) Element {
	return NewElement("sub", args...)
}

// SUP specifies inline text displayed as superscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sup
func SUP(args ...any) Element {
	return NewElement("sup", args...)
}

// I represents text in an alternate voice or mood.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/i
func I(args ...any) Element {
	return NewElement("i", args...)
}

// B draws attention to text without conveying importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/b
func B(args ...any) Element {
	return NewElement("b", args...)
}

// U represents text with an unarticulated annotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/u
func U(args ...any) Element {
	return NewElement("u", args...)
}

// MARK highlights text for reference.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/mark
func MARK(args ...any) Element {
	return NewElement("mark", args...)
}

// BDI isolates text for bidirectional text formatting.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdi
func BDI(args ...any) Element {
	return NewElement("bdi", args...)
}

// BDO overrides the current text direction.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdo
func BDO(args ...any) Element {
	return NewElement("bdo", args...)
}

// BR produces a line break in text (carriage-return). It is useful for writing a poem or an address, where the division of lines is significant.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/br
func BR(attrs ...Attribute) Element {
	return NewVoidElement("br", attrs...)
}

// WBR represents a word break opportunity.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/wbr
func WBR(attrs ...Attribute) Element {
	return NewVoidElement("wbr", attrs...)
}

// IMG embeds an image into the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img
func IMG(attrs ...Attribute) Element {
	return NewVoidElement("img", attrs...)
}

// IFRAME represents a nested browsing context, embedding another HTML page into the current one.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe
func IFRAME(args ...any) Element {
	return NewElement("iframe", args...)
}

// EMBED embeds external content at the specified point in the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/embed
func EMBED(attrs ...Attribute) Element {
	return NewVoidElement("embed", attrs...)
}

// OBJECT represents an external resource, which can be treated as an image, a nested browsing context, or a resource to be handled by a plugin.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/object
func OBJECT(args ...any) Element {
	return NewElement("object", args...)
}

// PICTURE defines multiple sources for an img element to offer alternative versions of an image for different display/device scenarios.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture
func PICTURE(args ...any) Element {
	return NewElement("picture", args...)
}

// SOURCE specifies multiple media resources for the picture, the audio element, or the video element. It is a void element, meaning that it has no content and does not have a closing tag. It is commonly used to offer the same media content in multiple file formats in order to provide compatibility with a broad range of browsers given their differing support for image file formats and media file formats.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/source
func SOURCE(attrs ...Attribute) Element {
	return NewVoidElement("source", attrs...)
}

// TRACK is used as a child of the media elements, audio and video. It lets you specify timed text tracks (or time-based data), for example to automatically handle subtitles. The tracks are formatted in WebVTT format (.vtt files)—Web Video Text Tracks.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/track
func TRACK(attrs ...Attribute) Element {
	return NewVoidElement("track", attrs...)
}

// VIDEO embeds a media player which supports video playback into the document. You can also use `<video>` for audio content, but the audio element may provide a more appropriate user experience.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video
func VIDEO(args ...any) Element {
	return NewElement("video", args...)
}

// AUDIO is used to embed sound content in documents. It may contain one or more audio sources, represented using the src attribute or the source element: the browser will choose the most suitable one. It can also be the destination for streamed media, using a MediaStream.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/audio
func AUDIO(args ...any) Element {
	return NewElement("audio", args...)
}

// CANVAS is a container element to use with either the canvas scripting API or the WebGL API to draw graphics and animations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/canvas
func CANVAS(args ...any) Element {
	return NewElement("canvas", args...)
}

// MAP is used with `<area>` elements to define an image map (a clickable link area).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/map
func MAP(args ...any) Element {
	return NewElement("map", args...)
}

// AREA defines an area inside an image map that has predefined clickable areas. An image map allows geometric areas on an image to be associated with hyperlink.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/area
func AREA(attrs ...Attribute) Element {
	return NewVoidElement("area", attrs...)
}

// SVG is a container defining a new coordinate system and viewport. It is used as the outermost element of SVG documents, but it can also be used to embed an SVG fragment inside an SVG or HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/svg
func SVG(args ...any) Element {
	return NewElement("svg", args...)
}

// MATH is the top-level element in MathML. Every valid MathML instance must be wrapped in it. In addition, you must not nest a second `<math>` element in another, but you can have an arbitrary number of other child elements in it.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/math
func MATH(args ...any) Element {
	return NewElement("math", args...)
}

// SCRIPT is used to embed executable code or data; this is typically used to embed or refer to JavaScript code. The `<script>` element can also be used with other languages, such as WebGL's GLSL shader programming language and JSON.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script
func SCRIPT(args ...any) Element {
	return NewElement("script", args...)
}

// NOSCRIPT defines a section of HTML to be inserted if a script type on the page is unsupported or if scripting is currently turned off in the browser.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/noscript
func NOSCRIPT(args ...any) Element {
	return NewElement("noscript", args...)
}

// DEL represents a range of text that has been deleted from a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/del
func DEL(args ...any) Element {
	return NewElement("del", args...)
}

// INS represents a range of text that has been added to a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ins
func INS(args ...any) Element {
	return NewElement("ins", args...)
}

// TABLE represents tabular data—that is, information presented in a two-dimensional table comprised of rows and columns of cells containing data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/table
func TABLE(args ...any) Element {
	return NewElement("table", args...)
}

// CAPTION specifies the caption (or title) of a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/caption
func CAPTION(args ...any) Element {
	return NewElement("caption", args...)
}

// COLGROUP defines a group of columns within a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/colgroup
func COLGROUP(args ...any) Element {
	return NewElement("colgroup", args...)
}

// COL defines one or more columns in a column group represented by its implicit or explicit parent `<colgroup>` element. The `<col>` element is only valid as a child of a `<colgroup>` element that has no span attribute defined.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/col
func COL(attrs ...Attribute) Element {
	return NewVoidElement("col", attrs...)
}

// THEAD groups the header content in a table with information about the table's columns. This is usually in the form of column headers (`<th>` elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/thead
func THEAD(args ...any) Element {
	return NewElement("thead", args...)
}

// TBODY groups the body content in a table. It typically contains the table's data rows (`<tr>` elements with `<td>` cells).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tbody
func TBODY(args ...any) Element {
	return NewElement("tbody", args...)
}

// TFOOT groups the footer content in a table with information about the table's columns. This is usually a summary of the columns, e.g., a sum of the given numbers in a column.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tfoot
func TFOOT(args ...any) Element {
	return NewElement("tfoot", args...)
}

// TR defines a row of cells in a table. The row's cells can then be established using a mix of `<td>` (data cell) and `<th>` (header cell) elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tr
func TR(args ...any) Element {
	return NewElement("tr", args...)
}

// TH is a child of the `<tr>` element, it defines a cell as the header of a group of table cells. The nature of this group can be explicitly defined by the scope and headers attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/th
func TH(args ...any) Element {
	return NewElement("th", args...)
}

// TD is a child of the `<tr>` element, it defines a cell of a table that contains data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/td
func TD(args ...any) Element {
	return NewElement("td", args...)
}

// FORM represents a document section containing interactive controls for submitting information.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form
func FORM(args ...any) Element {
	return NewElement("form", args...)
}

// FIELDSET is used to group several controls as well as labels (`<label>`) within a web form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fieldset
func FIELDSET(args ...any) Element {
	return NewElement("fieldset", args...)
}

// LEGEND represents a caption for the content of its parent `<fieldset>`.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/legend
func LEGEND(args ...any) Element {
	return NewElement("legend", args...)
}

// LABEL represents a caption for an item in a user interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label
func LABEL(args ...any) Element {
	return NewElement("label", args...)
}

// INPUT is used to create interactive controls for web-based forms to accept data from the user; a wide variety of types of input data and control widgets are available, depending on the device and user agent. The `<input>` element is one of the most powerful and complex in all of HTML due to the sheer number of combinations of input types and attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input
func INPUT(attrs ...Attribute) Element {
	return NewVoidElement("input", attrs...)
}

// BUTTON is an interactive element activated by a user with a mouse, keyboard, finger, voice command, or other assistive technology. Once activated, it performs an action, such as submitting a form or opening a dialog.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/button
func BUTTON(args ...any) Element {
	return NewElement("button", args...)
}

// SELECT represents a control that provides a menu of options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/select
func SELECT(args ...any) Element {
	return NewElement("select", args...)
}

// DATALIST contains a set of `<option>` elements that represent the permissible or recommended options available to choose from within other controls.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/datalist
func DATALIST(args ...any) Element {
	return NewElement("datalist", args...)
}

// OPTGROUP creates a grouping of options within a `<select>` element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/optgroup
func OPTGROUP(args ...any) Element {
	return NewElement("optgroup", args...)
}

// OPTION is used to define an item contained in a `<select>`, an `<optgroup>`, or a `<datalist>` element. As such, `<option>` can represent menu items in popups and other lists of items in an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/option
func OPTION(args ...any) Element {
	return NewElement("option", args...)
}

// TEXTAREA represents a multi-line plain-text editing control, useful when you want to allow users to enter a sizeable amount of free-form text, for example, a comment on a review or feedback form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/textarea
func TEXTAREA(args ...any) Element {
	return NewElement("textarea", args...)
}

// OUTPUT is a container element into which a site or app can inject the results of a calculation or the outcome of a user action.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/output
func OUTPUT(args ...any) Element {
	return NewElement("output", args...)
}

// PROGRESS displays an indicator showing the completion progress of a task, typically displayed as a progress bar.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/progress
func PROGRESS(args ...any) Element {
	return NewElement("progress", args...)
}

// METER represents either a scalar value within a known range or a fractional value.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meter
func METER(args ...any) Element {
	return NewElement("meter", args...)
}

// DETAILS creates a disclosure widget in which information is visible only when the widget is toggled into an "open" state. A summary or label must be provided using the `<summary>` element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details
func DETAILS(args ...any) Element {
	return NewElement("details", args...)
}

// SUMMARY specifies a summary, caption, or legend for a details element's disclosure box. Clicking the `<summary>` element toggles the state of the parent `<details>` element open and closed.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/summary
func SUMMARY(args ...any) Element {
	return NewElement("summary", args...)
}

// DIALOG represents a dialog box or other interactive component, such as a dismissible alert, inspector, or subwindow.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog
func DIALOG(args ...any) Element {
	return NewElement("dialog", args...)
}

// SLOT acts as a placeholder inside a web component that you can fill with your own markup, which lets you create separate DOM trees and present them together.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/slot
func SLOT(args ...any) Element {
	return NewElement("slot", args...)
}

// TEMPLATE holds HTML that is not to be rendered immediately when a page is loaded but may be instantiated subsequently during runtime using JavaScript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template
func TEMPLATE(args ...any) Element {
	return NewElement("template", args...)
}

// FENCEDFRAME represents a nested browsing context, like `<iframe>` but with more native privacy features built in.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fencedframe
func FENCEDFRAME(args ...any) Element {
	return NewElement("fencedframe", args...)
}

// SELECTEDCONTENT displays the content of the currently selected `<option>` inside a closed `<select>` element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/selectedcontent
func SELECTEDCONTENT(args ...any) Element {
	return NewElement("selectedcontent", args...)
}

// BASE specifies the base URL and default browsing context for relative URLs.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
func BASE(attrs ...Attribute) Element {
	return NewVoidElement("base", attrs...)
}

// HGROUP groups a set of h1–h6 elements when they represent a multi-level heading.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hgroup
func HGROUP(args ...any) Element {
	return NewElement("hgroup", args...)
}

// ADDRESS indicates contact information for a person or organization.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
func ADDRESS(args ...any) Element {
	return NewElement("address", args...)
}

// SEARCH represents a search or filtering interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/search
func SEARCH(args ...any) Element {
	return NewElement("search", args...)
}

// DIV is the generic container for flow content.
//
// Examples:
//
//	// Empty div
//	DIV()
//
//	// With text content
//	DIV("Hello World")
//
//	// With attributes and children
//	DIV(AttrClass("container"), AttrId("main"), "Content")
//	DIV(Attr("class", "container"), "Content")
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/div
func DIV(args ...any) Element {
	return NewElement("div", args...)
}

// SPAN is the generic inline container for phrasing content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/span
func SPAN(args ...any) Element {
	return NewElement("span", args...)
}

// P creates a paragraph element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/p
func P(args ...any) Element {
	return NewElement("p", args...)
}

// DL represents a description list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dl
func DL(args ...any) Element {
	return NewElement("dl", args...)
}

// DT specifies a term in a description or definition list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dt
func DT(args ...any) Element {
	return NewElement("dt", args...)
}

// DD provides the description, definition, or value for the preceding term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dd
func DD(args ...any) Element {
	return NewElement("dd", args...)
}

// FIGURE represents self-contained content with an optional caption.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figure
func FIGURE(args ...any) Element {
	return NewElement("figure", args...)
}

// FIGCAPTION represents a caption or legend for the contents of its parent figure element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figcaption
func FIGCAPTION(args ...any) Element {
	return NewElement("figcaption", args...)
}

// MENU represents a set of commands or options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/menu
func MENU(args ...any) Element {
	return NewElement("menu", args...)
}

// SMALL represents side-comments and small print.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/small
func SMALL(args ...any) Element {
	return NewElement("small", args...)
}

// S renders text with a strikethrough.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/s
func S(args ...any) Element {
	return NewElement("s", args...)
}

// CITE marks the title of a creative work.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/cite
func CITE(args ...any) Element {
	return NewElement("cite", args...)
}

// Q indicates a short inline quotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/q
func Q(args ...any) Element {
	return NewElement("q", args...)
}

// DFN indicates the defining instance of a term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dfn
func DFN(args ...any) Element {
	return NewElement("dfn", args...)
}

// ABBR represents an abbreviation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/abbr
func ABBR(args ...any) Element {
	return NewElement("abbr", args...)
}

// RUBY represents ruby annotations for East Asian typography.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ruby
func RUBY(args ...any) Element {
	return NewElement("ruby", args...)
}

// RT specifies the ruby text for ruby annotations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rt
func RT(args ...any) Element {
	return NewElement("rt", args...)
}

// RP provides parentheses for browsers that don't support ruby text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rp
func RP(args ...any) Element {
	return NewElement("rp", args...)
}

// DATA links content with a machine-readable translation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/data
func DATA(args ...any) Element {
	return NewElement("data", args...)
}

// TIME represents a specific period in time.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/time
func TIME(args ...any) Element {
	return NewElement("time", args...)
}
