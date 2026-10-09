package hyper

import (
	"bytes"
	"io"
)

// HyperNode represents any renderable HTML element or text content.
//
// The HyperNode interface is the core abstraction that allows both HTML elements
// and text content to be treated uniformly when building and rendering HTML
// trees. All elements created by the factory functions (DIV, P, SVG, etc.)
// implement this interface.
//
// Example:
//
//	node := DIV("Hello")
//	err := node.RenderNode(os.Stdout)
type HyperNode interface {
	RenderNode(io.Writer) error
}

// RenderNode writes the HTML representation of a HyperNode to the provided io.Writer.
//
// The writer comes first, before the node: you already have the
// destination in hand (stdout, a file, an http.ResponseWriter) by the
// time you build the node.
//
// Example:
//
//	err := RenderNode(os.Stdout, DIV("Hello")) // Outputs: <div>Hello</div>
func RenderNode(w io.Writer, node HyperNode) error {
	if node == nil {
		panic("nil node passed to RenderNode")
	}
	return node.RenderNode(w)
}

// RenderNodeThen renders a HyperNode and passes the resulting bytes to the
// provided callback function. This is useful for capturing rendered output for further
// processing without writing directly to an io.Writer.
//
// Example:
//
//	err := RenderNodeThen(node, func(data []byte) error {
//		return saveToCache(data)
//	})
func RenderNodeThen(node HyperNode, then func(data []byte) error) error {
	if node == nil {
		panic("nil node passed to RenderNodeThen")
	}
	var buffer bytes.Buffer
	if err := node.RenderNode(&buffer); err != nil {
		return err
	}
	return then(buffer.Bytes())
}
