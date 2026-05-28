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
	Tag        string      // HTML tag name
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

// ChildrenInserter is a function that constructs an [Element] with children.
// It is returned by element functions (DIV, P, BODY, etc.) and must be called
// to produce the final [Element].
//
// Example:
//
//	DIV(AttrClass("container"))("Hello")
type ChildrenInserter func(children ...any) Element

// MakeChildrenInserter wraps an [Element] and returns an [ChildrenInserter] that accepts children.
func MakeChildrenInserter(element Element) ChildrenInserter {
	return func(children ...any) Element {
		element.InsertChildren(children...)
		return element
	}
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
		switch value := child.(type) {
		case HyperNode:
			me.Children = append(me.Children, value)
		// Explicit string and fmt.Stringer cases for performance:
		// fmt.Sprint() would handle these, but with overhead from type inspection and buffer allocation.
		case string:
			me.Children = append(me.Children, Text(value))
		case fmt.Stringer:
			me.Children = append(me.Children, Text(value.String()))
		// Nil arguments are not checked/filtered - fmt.Sprint() will render them as "Nil",
		// which is intentional for better debugging (makes it obvious when nil values are passed).
		default:
			me.Children = append(me.Children, Text(fmt.Sprint(value)))
		}
	}
}

// toHyperNode converts an arbitrary value to a [HyperNode].
// Strings become [Text], fmt.Stringer values are converted via String(),
// and all other values use fmt.Sprint.
func toHyperNode(v any) HyperNode {
	switch value := v.(type) {
	case HyperNode:
		return value
	case string:
		return Text(value)
	case fmt.Stringer:
		return Text(value.String())
	default:
		return Text(fmt.Sprint(value))
	}
}

// DOCTYPE creates the <!DOCTYPE html> element.
//
// https://developer.mozilla.org/en-US/docs/Glossary/Doctype
func DOCTYPE() HyperNode {
	return Element{Tag: "!DOCTYPE html", IsVoid: true}
}

// HTML creates the root element of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
func HTML(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "html", Attributes: attrs})
}

// HEAD contains machine-readable information about the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
func HEAD(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "head", Attributes: attrs})
}

// TITLE defines the document's title that is shown in a browser's title bar or a page's tab.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
func TITLE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "title", Attributes: attrs})
}

// LINK specifies relationships between the current document and an external resource.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
func LINK(attrs ...Attribute) HyperNode {
	return Element{Tag: "link", IsVoid: true, Attributes: attrs}
}

// META represents metadata that cannot be represented by other HTML meta-related elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
func META(attrs ...Attribute) HyperNode {
	return Element{Tag: "meta", IsVoid: true, Attributes: attrs}
}

// STYLE contains style information for a document or part of a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
func STYLE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "style", Attributes: attrs})
}

// BODY represents the content of an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
func BODY(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "body", Attributes: attrs})
}

// H1 creates a level 1 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h1
func H1(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "h1", Attributes: attrs})
}

// H2 creates a level 2 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h2
func H2(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "h2", Attributes: attrs})
}

// H3 creates a level 3 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h3
func H3(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "h3", Attributes: attrs})
}

// H4 creates a level 4 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h4
func H4(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "h4", Attributes: attrs})
}

// H5 creates a level 5 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h5
func H5(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "h5", Attributes: attrs})
}

// H6 creates a level 6 heading element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/h6
func H6(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "h6", Attributes: attrs})
}

// HEADER creates a header element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
func HEADER(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "header", Attributes: attrs})
}

// FOOTER creates a footer element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
func FOOTER(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "footer", Attributes: attrs})
}

// NAV creates a navigation element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/nav
func NAV(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "nav", Attributes: attrs})
}

// MAIN creates a main content element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/main
func MAIN(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "main", Attributes: attrs})
}

// SECTION creates a section element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/section
func SECTION(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "section", Attributes: attrs})
}

// ARTICLE creates an article element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
func ARTICLE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "article", Attributes: attrs})
}

// ASIDE creates an aside element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
func ASIDE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "aside", Attributes: attrs})
}

// HR represents a thematic break between paragraph-level elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hr
func HR(attrs ...Attribute) HyperNode {
	return Element{Tag: "hr", IsVoid: true, Attributes: attrs}
}

// PRE represents preformatted text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/pre
func PRE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "pre", Attributes: attrs})
}

// BLOCKQUOTE represents a section quoted from another source.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote
func BLOCKQUOTE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "blockquote", Attributes: attrs})
}

// OL represents an ordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ol
func OL(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "ol", Attributes: attrs})
}

// UL represents an unordered list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ul
func UL(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "ul", Attributes: attrs})
}

// LI represents a list item.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/li
func LI(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "li", Attributes: attrs})
}

// A creates hyperlinks to other web pages, files, locations within the same page, or anything else a URL can address.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a
func A(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "a", Attributes: attrs})
}

// EM marks text with emphasis.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/em
func EM(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "em", Attributes: attrs})
}

// STRONG indicates strong importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/strong
func STRONG(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "strong", Attributes: attrs})
}

// CODE displays its contents styled as computer code.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/code
func CODE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "code", Attributes: attrs})
}

// VAR represents a variable in a mathematical expression or programming context.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/var
func VAR(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "var", Attributes: attrs})
}

// SAMP represents sample output from a computer program.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/samp
func SAMP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "samp", Attributes: attrs})
}

// KBD represents text that the user should enter.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/kbd
func KBD(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "kbd", Attributes: attrs})
}

// SUB specifies inline text displayed as subscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sub
func SUB(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "sub", Attributes: attrs})
}

// SUP specifies inline text displayed as superscript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sup
func SUP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "sup", Attributes: attrs})
}

// I represents text in an alternate voice or mood.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/i
func I(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "i", Attributes: attrs})
}

// B draws attention to text without conveying importance.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/b
func B(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "b", Attributes: attrs})
}

// U represents text with an unarticulated annotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/u
func U(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "u", Attributes: attrs})
}

// MARK highlights text for reference.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/mark
func MARK(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "mark", Attributes: attrs})
}

// BDI isolates text for bidirectional text formatting.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdi
func BDI(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "bdi", Attributes: attrs})
}

// BDO overrides the current text direction.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdo
func BDO(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "bdo", Attributes: attrs})
}

// BR produces a line break in text (carriage-return). It is useful for writing a poem or an address, where the division of lines is significant.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/br
func BR(attrs ...Attribute) HyperNode {
	return Element{Tag: "br", IsVoid: true, Attributes: attrs}
}

// WBR represents a word break opportunity.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/wbr
func WBR(attrs ...Attribute) HyperNode {
	return Element{Tag: "wbr", IsVoid: true, Attributes: attrs}
}

// IMG embeds an image into the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img
func IMG(attrs ...Attribute) HyperNode {
	return Element{Tag: "img", IsVoid: true, Attributes: attrs}
}

// IFRAME represents a nested browsing context, embedding another HTML page into the current one.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe
func IFRAME(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "iframe", Attributes: attrs})
}

// EMBED embeds external content at the specified point in the document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/embed
func EMBED(attrs ...Attribute) HyperNode {
	return Element{Tag: "embed", IsVoid: true, Attributes: attrs}
}

// OBJECT represents an external resource, which can be treated as an image, a nested browsing context, or a resource to be handled by a plugin.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/object
func OBJECT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "object", Attributes: attrs})
}

// PICTURE defines multiple sources for an img element to offer alternative versions of an image for different display/device scenarios.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture
func PICTURE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "picture", Attributes: attrs})
}

// SOURCE specifies multiple media resources for the picture, the audio element, or the video element. It is a void element, meaning that it has no content and does not have a closing tag. It is commonly used to offer the same media content in multiple file formats in order to provide compatibility with a broad range of browsers given their differing support for image file formats and media file formats.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/source
func SOURCE(attrs ...Attribute) HyperNode {
	return Element{Tag: "source", IsVoid: true, Attributes: attrs}
}

// TRACK is used as a child of the media elements, audio and video. It lets you specify timed text tracks (or time-based data), for example to automatically handle subtitles. The tracks are formatted in WebVTT format (.vtt files)—Web Video Text Tracks.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/track
func TRACK(attrs ...Attribute) HyperNode {
	return Element{Tag: "track", IsVoid: true, Attributes: attrs}
}

// VIDEO embeds a media player which supports video playback into the document. You can also use &lt;video&gt; for audio content, but the audio element may provide a more appropriate user experience.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video
func VIDEO(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "video", Attributes: attrs})
}

// AUDIO is used to embed sound content in documents. It may contain one or more audio sources, represented using the src attribute or the source element: the browser will choose the most suitable one. It can also be the destination for streamed media, using a MediaStream.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/audio
func AUDIO(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "audio", Attributes: attrs})
}

// CANVAS is a container element to use with either the canvas scripting API or the WebGL API to draw graphics and animations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/canvas
func CANVAS(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "canvas", Attributes: attrs})
}

// MAP is used with &lt;area&gt; elements to define an image map (a clickable link area).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/map
func MAP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "map", Attributes: attrs})
}

// AREA defines an area inside an image map that has predefined clickable areas. An image map allows geometric areas on an image to be associated with hyperlink.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/area
func AREA(attrs ...Attribute) HyperNode {
	return Element{Tag: "area", IsVoid: true, Attributes: attrs}
}

// SVG is a container defining a new coordinate system and viewport. It is used as the outermost element of SVG documents, but it can also be used to embed an SVG fragment inside an SVG or HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/svg
func SVG(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "svg", Attributes: attrs})
}

// MATH is the top-level element in MathML. Every valid MathML instance must be wrapped in it. In addition, you must not nest a second &lt;math&gt; element in another, but you can have an arbitrary number of other child elements in it.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/math
func MATH(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "math", Attributes: attrs})
}

// SCRIPT is used to embed executable code or data; this is typically used to embed or refer to JavaScript code. The &lt;script&gt; element can also be used with other languages, such as WebGL's GLSL shader programming language and JSON.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script
func SCRIPT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "script", Attributes: attrs})
}

// NOSCRIPT defines a section of HTML to be inserted if a script type on the page is unsupported or if scripting is currently turned off in the browser.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/noscript
func NOSCRIPT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "noscript", Attributes: attrs})
}

// DEL represents a range of text that has been deleted from a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/del
func DEL(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "del", Attributes: attrs})
}

// INS represents a range of text that has been added to a document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ins
func INS(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "ins", Attributes: attrs})
}

// TABLE represents tabular data—that is, information presented in a two-dimensional table comprised of rows and columns of cells containing data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/table
func TABLE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "table", Attributes: attrs})
}

// CAPTION specifies the caption (or title) of a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/caption
func CAPTION(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "caption", Attributes: attrs})
}

// COLGROUP defines a group of columns within a table.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/colgroup
func COLGROUP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "colgroup", Attributes: attrs})
}

// COL defines one or more columns in a column group represented by its implicit or explicit parent &lt;colgroup&gt; element. The &lt;col&gt; element is only valid as a child of a &lt;colgroup&gt; element that has no span attribute defined.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/col
func COL(attrs ...Attribute) HyperNode {
	return Element{Tag: "col", IsVoid: true, Attributes: attrs}
}

// THEAD groups the header content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/thead
func THEAD(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "thead", Attributes: attrs})
}

// TBODY groups the body content in a table with information about the table's columns. This is usually in the form of column headers (&lt;th&gt; elements).
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tbody
func TBODY(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "tbody", Attributes: attrs})
}

// TFOOT groups the footer content in a table with information about the table's columns. This is usually a summary of the columns, e.g., a sum of the given numbers in a column.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tfoot
func TFOOT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "tfoot", Attributes: attrs})
}

// TR defines a row of cells in a table. The row's cells can then be established using a mix of &lt;td&gt; (data cell) and &lt;th&gt; (header cell) elements.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tr
func TR(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "tr", Attributes: attrs})
}

// TH is a child of the &lt;tr&gt; element, it defines a cell as the header of a group of table cells. The nature of this group can be explicitly defined by the scope and headers attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/th
func TH(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "th", Attributes: attrs})
}

// TD is a child of the &lt;tr&gt; element, it defines a cell of a table that contains data.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/td
func TD(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "td", Attributes: attrs})
}

// FORM represents a document section containing interactive controls for submitting information.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form
func FORM(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "form", Attributes: attrs})
}

// FIELDSET is used to group several controls as well as labels (&lt;label&gt;) within a web form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fieldset
func FIELDSET(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "fieldset", Attributes: attrs})
}

// LEGEND represents a caption for the content of its parent &lt;fieldset&gt;.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/legend
func LEGEND(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "legend", Attributes: attrs})
}

// LABEL represents a caption for an item in a user interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label
func LABEL(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "label", Attributes: attrs})
}

// INPUT is used to create interactive controls for web-based forms to accept data from the user; a wide variety of types of input data and control widgets are available, depending on the device and user agent. The &lt;input&gt; element is one of the most powerful and complex in all of HTML due to the sheer number of combinations of input types and attributes.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input
func INPUT(attrs ...Attribute) HyperNode {
	return Element{Tag: "input", IsVoid: true, Attributes: attrs}
}

// BUTTON is an interactive element activated by a user with a mouse, keyboard, finger, voice command, or other assistive technology. Once activated, it performs an action, such as submitting a form or opening a dialog.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/button
func BUTTON(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "button", Attributes: attrs})
}

// SELECT represents a control that provides a menu of options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/select
func SELECT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "select", Attributes: attrs})
}

// DATALIST contains a set of &lt;option&gt; elements that represent the permissible or recommended options available to choose from within other controls.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/datalist
func DATALIST(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "datalist", Attributes: attrs})
}

// OPTGROUP creates a grouping of options within a &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/optgroup
func OPTGROUP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "optgroup", Attributes: attrs})
}

// OPTION is used to define an item contained in a &lt;select&gt;, an &lt;optgroup&gt;, or a &lt;datalist&gt; element. As such, &lt;option&gt; can represent menu items in popups and other lists of items in an HTML document.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/option
func OPTION(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "option", Attributes: attrs})
}

// TEXTAREA represents a multi-line plain-text editing control, useful when you want to allow users to enter a sizeable amount of free-form text, for example, a comment on a review or feedback form.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/textarea
func TEXTAREA(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "textarea", Attributes: attrs})
}

// OUTPUT is a container element into which a site or app can inject the results of a calculation or the outcome of a user action.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/output
func OUTPUT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "output", Attributes: attrs})
}

// PROGRESS displays an indicator showing the completion progress of a task, typically displayed as a progress bar.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/progress
func PROGRESS(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "progress", Attributes: attrs})
}

// METER represents either a scalar value within a known range or a fractional value.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meter
func METER(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "meter", Attributes: attrs})
}

// DETAILS creates a disclosure widget in which information is visible only when the widget is toggled into an "open" state. A summary or label must be provided using the &lt;summary&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details
func DETAILS(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "details", Attributes: attrs})
}

// SUMMARY specifies a summary, caption, or legend for a details element's disclosure box. Clicking the &lt;summary&gt; element toggles the state of the parent &lt;details&gt; element open and closed.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/summary
func SUMMARY(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "summary", Attributes: attrs})
}

// DIALOG represents a dialog box or other interactive component, such as a dismissible alert, inspector, or subwindow.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog
func DIALOG(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "dialog", Attributes: attrs})
}

// SLOT acts as a placeholder inside a web component that you can fill with your own markup, which lets you create separate DOM trees and present them together.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/slot
func SLOT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "slot", Attributes: attrs})
}

// TEMPLATE holds HTML that is not to be rendered immediately when a page is loaded but may be instantiated subsequently during runtime using JavaScript.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template
func TEMPLATE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "template", Attributes: attrs})
}

// FENCEDFRAME represents a nested browsing context, like &lt;iframe&gt; but with more native privacy features built in.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fencedframe
func FENCEDFRAME(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "fencedframe", Attributes: attrs})
}

// SELECTEDCONTENT displays the content of the currently selected &lt;option&gt; inside a closed &lt;select&gt; element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/selectedcontent
func SELECTEDCONTENT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "selectedcontent", Attributes: attrs})
}

// BASE specifies the base URL and default browsing context for relative URLs.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
func BASE(attrs ...Attribute) HyperNode {
	return Element{Tag: "base", IsVoid: true, Attributes: attrs}
}

// HGROUP groups a set of h1–h6 elements when they represent a multi-level heading.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hgroup
func HGROUP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "hgroup", Attributes: attrs})
}

// ADDRESS indicates contact information for a person or organization.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
func ADDRESS(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "address", Attributes: attrs})
}

// SEARCH represents a search or filtering interface.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/search
func SEARCH(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "search", Attributes: attrs})
}

// DIV is the generic container for flow content.
//
// Examples:
//
//	// Empty div
//	DIV()()
//
//	// With text content
//	DIV()("Hello World")
//
//	// With attributes and children
//	DIV(AttrClass("container"), AttrId("main"))("Content")
//	DIV(Attr("class", "container"))("Content")
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/div
func DIV(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "div", Attributes: attrs})
}

// SPAN is the generic inline container for phrasing content.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/span
func SPAN(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "span", Attributes: attrs})
}

// P creates a paragraph element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/p
func P(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "p", Attributes: attrs})
}

// DL represents a description list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dl
func DL(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "dl", Attributes: attrs})
}

// DT specifies a term in a description or definition list.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dt
func DT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "dt", Attributes: attrs})
}

// DD provides the description, definition, or value for the preceding term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dd
func DD(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "dd", Attributes: attrs})
}

// FIGURE represents self-contained content with an optional caption.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figure
func FIGURE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "figure", Attributes: attrs})
}

// FIGCAPTION represents a caption or legend for the contents of its parent figure element.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figcaption
func FIGCAPTION(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "figcaption", Attributes: attrs})
}

// MENU represents a set of commands or options.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/menu
func MENU(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "menu", Attributes: attrs})
}

// SMALL represents side-comments and small print.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/small
func SMALL(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "small", Attributes: attrs})
}

// S renders text with a strikethrough.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/s
func S(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "s", Attributes: attrs})
}

// CITE marks the title of a creative work.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/cite
func CITE(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "cite", Attributes: attrs})
}

// Q indicates a short inline quotation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/q
func Q(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "q", Attributes: attrs})
}

// DFN indicates the defining instance of a term.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dfn
func DFN(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "dfn", Attributes: attrs})
}

// ABBR represents an abbreviation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/abbr
func ABBR(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "abbr", Attributes: attrs})
}

// RUBY represents ruby annotations for East Asian typography.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ruby
func RUBY(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "ruby", Attributes: attrs})
}

// RT specifies the ruby text for ruby annotations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rt
func RT(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "rt", Attributes: attrs})
}

// RP provides parentheses for browsers that don't support ruby text.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rp
func RP(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "rp", Attributes: attrs})
}

// DATA links content with a machine-readable translation.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/data
func DATA(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "data", Attributes: attrs})
}

// TIME represents a specific period in time.
//
// https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/time
func TIME(attrs ...Attribute) ChildrenInserter {
	return MakeChildrenInserter(Element{Tag: "time", Attributes: attrs})
}
