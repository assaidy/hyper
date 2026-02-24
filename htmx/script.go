package hx

import (
	_ "embed"
	h "github.com/assaidy/hyper"
)

//go:embed htmx@2.0.8.js
var script string

// Script returns the htmx JavaScript as a <script> element.
// Use this to include htmx in your pages.
func Script() h.Node {
	return h.Script(h.RawText(script))
}
