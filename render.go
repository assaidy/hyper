package h

import "io"

// Render writes the HTML representation of a Node to the provided io.Writer.
//
// This is a convenience function that makes it suitable for writing directly to
// files, HTTP responses, or other output streams.
//
// Example:
//
//	err := Render(os.Stdout, Div("Hello"))
//	// Outputs: <div>Hello</div>
func Render(w io.Writer, node Node) error {
	return node.Render(w)
}

// Node represents any renderable HTML element or text content.
//
// The Node interface is the core abstraction that allows both HTML elements
// and text content to be treated uniformly when building and rendering HTML
// trees. All elements created by the factory functions (Div(), P(), Svg(), etc.)
// implement this interface.
//
// Example:
//
//	var node Node = Div("Hello")
//	err := node.Render(os.Stdout)
type Node interface {
	Render(io.Writer) error
}
