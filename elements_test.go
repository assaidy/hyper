package hyper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text passes through", "hello", "hello"},
		{"escapes angle brackets", "<b>bold</b>", "&lt;b&gt;bold&lt;/b&gt;"},
		{"escapes ampersands", "a & b", "a &amp; b"},
		{"escapes double quotes", `he said "hi"`, "he said &#34;hi&#34;"},
		{"escapes single quotes", "it's", "it&#39;s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, mustRender(t, Text(tt.in)))
		})
	}
}

func TestRawText(t *testing.T) {
	require.Equal(t, "<b>raw</b> & unchanged", mustRender(t, RawText("<b>raw</b> & unchanged")))
}

func TestElementRenderNode(t *testing.T) {
	t.Run("renders tag, attributes in order, and children", func(t *testing.T) {
		e := Element{
			Name: "div",
			Attributes: []Attribute{
				PairAttribute{Key: "class", Value: "box"},
				BooleanAttribute{Key: "hidden", IsActive: true},
			},
			Children: []HyperNode{Text("hi")},
		}
		require.Equal(t, `<div class="box" hidden>hi</div>`, mustRender(t, e))
	})

	t.Run("void elements have no closing tag", func(t *testing.T) {
		e := Element{
			Name:       "img",
			IsVoid:     true,
			Attributes: []Attribute{PairAttribute{Key: "src", Value: "a.png"}},
		}
		require.Equal(t, `<img src="a.png">`, mustRender(t, e))
	})

	t.Run("void elements drop children", func(t *testing.T) {
		e := Element{Name: "br", IsVoid: true, Children: []HyperNode{Text("x")}}
		require.Equal(t, "<br>", mustRender(t, e))
	})

	t.Run("returns an error for an empty tag name", func(t *testing.T) {
		e := Element{Children: []HyperNode{Text("a")}}
		out, err := tryRender(t, e)
		require.EqualError(t, err, "empty element name not allowed.")
		require.Empty(t, out)
	})

	t.Run("propagates empty-name errors from nested elements", func(t *testing.T) {
		e := Element{Name: "div", Children: []HyperNode{Element{}}}
		out, err := tryRender(t, e)
		require.EqualError(t, err, "empty element name not allowed.")
		require.Empty(t, out)
	})

	t.Run("renders an empty tag when there are no children", func(t *testing.T) {
		require.Equal(t, "<div></div>", mustRender(t, Element{Name: "div"}))
	})

	t.Run("skips nil attributes", func(t *testing.T) {
		e := Element{
			Name:   "input",
			IsVoid: true,
			Attributes: []Attribute{
				PairAttribute{Key: "type", Value: "text"},
				nil,
				BooleanAttribute{Key: "disabled", IsActive: true},
			},
		}
		require.Equal(t, `<input type="text" disabled>`, mustRender(t, e))
	})

	t.Run("renders every child kind", func(t *testing.T) {
		counter := &countingNode{}
		e := Element{
			Name: "div",
			Children: []HyperNode{
				Element{Name: "b", Children: []HyperNode{Text("fast")}}, // nested fast path
				Text("<escaped>"),
				RawText("<i>raw</i>"),
				counter, // default path
			},
		}
		require.Equal(t, `<div><b>fast</b>&lt;escaped&gt;<i>raw</i>counted</div>`, mustRender(t, e))
		require.Equal(t, 1, counter.calls)
	})

	t.Run("writes nothing when a child fails", func(t *testing.T) {
		e := Element{Name: "div", Children: []HyperNode{Text("partial"), errNode{errBoom}}}
		out, err := tryRender(t, e)
		require.ErrorIs(t, err, errBoom)
		require.Empty(t, out)
	})

	t.Run("writes nothing when an attribute fails", func(t *testing.T) {
		e := Element{Name: "div", Attributes: []Attribute{PairAttribute{Key: " "}}}
		out, err := tryRender(t, e)
		require.EqualError(t, err, "empty/whitespace attribute key not allowed.")
		require.Empty(t, out)
	})

	t.Run("propagates writer errors", func(t *testing.T) {
		err := RenderNode(failingWriter{errBoom}, Element{Name: "div", Children: []HyperNode{Text("x")}})
		require.ErrorIs(t, err, errBoom)
	})
}

func TestElementInsertAttributes(t *testing.T) {
	t.Run("no-op without arguments", func(t *testing.T) {
		e := Element{Name: "div", Attributes: []Attribute{PairAttribute{Key: "class", Value: "a"}}}
		e.InsertAttributes()
		require.Len(t, e.Attributes, 1)
	})

	t.Run("appends and preserves existing attributes", func(t *testing.T) {
		e := Element{Name: "div", Attributes: []Attribute{PairAttribute{Key: "class", Value: "a"}}}
		e.InsertAttributes(
			PairAttribute{Key: "id", Value: "b"},
			BooleanAttribute{Key: "hidden", IsActive: true},
		)
		require.Equal(t, []Attribute{
			PairAttribute{Key: "class", Value: "a"},
			PairAttribute{Key: "id", Value: "b"},
			BooleanAttribute{Key: "hidden", IsActive: true},
		}, e.Attributes)
	})
}

func TestElementInsertChildren(t *testing.T) {
	t.Run("no-op without arguments", func(t *testing.T) {
		e := Element{Name: "div", Children: []HyperNode{Text("keep")}}
		e.InsertChildren()
		require.Len(t, e.Children, 1)
	})

	t.Run("converts values and appends in order", func(t *testing.T) {
		child := P("node")
		e := Element{Name: "div", Children: []HyperNode{Text("first")}}
		e.InsertChildren(child, "text", 42)

		require.Equal(t, []HyperNode{Text("first"), child, Text("text"), Text("42")}, e.Children)
	})
}

type stringerValue struct{ s string }

func (v stringerValue) String() string { return v.s }

func TestToHyperNode(t *testing.T) {
	t.Run("hyper nodes pass through untouched", func(t *testing.T) {
		counter := &countingNode{}
		require.Same(t, counter, toHyperNode(counter))
	})

	t.Run("strings become text nodes", func(t *testing.T) {
		require.Equal(t, Text("hi"), toHyperNode("hi"))
	})

	t.Run("fmt.Stringer values become their string form", func(t *testing.T) {
		require.Equal(t, Text("labeled!"), toHyperNode(stringerValue{"labeled!"}))
	})

	t.Run("other values become their fmt.Sprint form", func(t *testing.T) {
		require.Equal(t, Text("42"), toHyperNode(42))
		require.Equal(t, Text("true"), toHyperNode(true))
		require.Equal(t, Text("<nil>"), toHyperNode(nil))
	})
}

func TestNewElement(t *testing.T) {
	t.Run("separates attributes from children", func(t *testing.T) {
		e := NewElement("div", Attr("class", "box"), "child", Attr("hidden", true), 42)
		require.Equal(t, Element{
			Name: "div",
			Attributes: []Attribute{
				PairAttribute{Key: "class", Value: "box"},
				BooleanAttribute{Key: "hidden", IsActive: true},
			},
			Children: []HyperNode{Text("child"), Text("42")},
		}, e)
	})

	t.Run("builds an empty element without arguments", func(t *testing.T) {
		require.Equal(t, Element{Name: "div"}, NewElement("div"))
	})

	t.Run("panics on a nil argument", func(t *testing.T) {
		require.PanicsWithValue(t, "nil argument passed to NewElement", func() {
			NewElement("div", nil, "keep")
		})
	})

	t.Run("conditional attributes render only when active", func(t *testing.T) {
		require.Equal(t, "<div>x</div>",
			mustRender(t, NewElement("div", IfAttr(false, AttrHidden(true)), "x")))
		require.Equal(t, "<div hidden>x</div>",
			mustRender(t, NewElement("div", IfAttr(true, AttrHidden(true)), "x")))
	})
}

func TestNewVoidElement(t *testing.T) {
	t.Run("marks the element void and keeps the attributes", func(t *testing.T) {
		require.Equal(t, Element{
			Name:       "img",
			Attributes: []Attribute{PairAttribute{Key: "src", Value: "a.png"}},
			IsVoid:     true,
		}, NewVoidElement("img", PairAttribute{Key: "src", Value: "a.png"}))
	})

	t.Run("skips nil attributes when rendering", func(t *testing.T) {
		require.Equal(t, "<input>",
			mustRender(t, NewVoidElement("input", IfAttr(false, AttrDisabled(true)))))
	})
}

func TestDOCTYPE(t *testing.T) {
	require.Equal(t, "<!DOCTYPE html>", mustRender(t, DOCTYPE()))
}
