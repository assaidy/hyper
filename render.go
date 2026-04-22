package hyper

import (
	"bytes"
	"io"
)

// Render writes the HTML representation of a Node to the provided io.Writer.
//
// This is a convenience function that makes it suitable for writing directly to
// files, HTTP responses, or other output streams.
//
// Example:
//
//	err := Render(os.Stdout, DIV()("Hello")) // Outputs: <div>Hello</div>
func Render(w io.Writer, node HyperNode) error {
	// using Group() because is nil-safe
	return Group(node).Render(w)
}

// RenderThen renders a HyperNode and passes the resulting bytes to the
// provided callback function. This is useful for capturing rendered output for further
// processing without writing directly to an io.Writer.
//
// Example:
//
//	err := RenderThen(node, func(data []byte) error {
//		return saveToDatabase(data)
//	})
func RenderThen(node HyperNode, then func(data []byte) error) error {
	var buffer bytes.Buffer
	if err := Render(&buffer, node); err != nil {
		return err
	}
	return then(buffer.Bytes())
}

// HyperNode represents any renderable HTML element or text content.
//
// The HyperNode interface is the core abstraction that allows both HTML elements
// and text content to be treated uniformly when building and rendering HTML
// trees. All elements created by the factory functions (DIV(), P(), SVG(), etc.)
// implement this interface.
//
// Example:
//
//	var node HyperNode = DIV()("Hello")
//	err := node.Render(os.Stdout)
type HyperNode interface {
	Render(io.Writer) error
}
