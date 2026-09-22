package hyper

import (
	"bytes"
	"testing"
)

// stringerType is a test type that implements fmt.Stringer
type stringerType string

func (s stringerType) String() string {
	return string(s)
}

func TestTextNode_Render(t *testing.T) {
	tests := []struct {
		name     string
		text     Text
		expected string
	}{
		{
			name:     "Simple text",
			text:     Text("Hello World"),
			expected: "Hello World",
		},
		{
			name:     "Empty text",
			text:     Text(""),
			expected: "",
		},
		{
			name:     "Text with HTML entities",
			text:     Text("<script>alert('xss')</script>"),
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "Text with quotes",
			text:     Text("Hello \"World\" & 'Universe'"),
			expected: "Hello &#34;World&#34; &amp; &#39;Universe&#39;",
		},
		{
			name:     "Text with leading space",
			text:     Text("  hello"),
			expected: "  hello",
		},
		{
			name:     "Text with trailing space",
			text:     Text("hello  "),
			expected: "hello  ",
		},
		{
			name:     "Text with leading and trailing spaces",
			text:     Text("  hello  "),
			expected: "  hello  ",
		},
		{
			name:     "Text with multiple spaces inside",
			text:     Text("hello    world"),
			expected: "hello    world",
		},
		{
			name:     "Text with newlines and tabs",
			text:     Text("hello\n\tworld"),
			expected: "hello\n\tworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.text.Render(&buf)
			if err != nil {
				t.Errorf("textNode.Render() returned error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("textNode.Render() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestRawText_Render(t *testing.T) {
	tests := []struct {
		name     string
		html     RawText
		expected string
	}{
		{
			name:     "Simple raw HTML",
			html:     RawText("Hello World"),
			expected: "Hello World",
		},
		{
			name:     "Raw HTML with HTML",
			html:     RawText("<script>alert('xss')</script>"),
			expected: "<script>alert('xss')</script>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.html.Render(&buf)
			if err != nil {
				t.Errorf("RawText.Render() returned error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("RawText.Render() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestTextNode_Render_Error(t *testing.T) {
	// textNode.Render() should never return an error based on the implementation
	// This test ensures that behavior
	text := Text("test")
	var buf bytes.Buffer
	err := text.Render(&buf)
	if err != nil {
		t.Errorf("textNode.Render() should not return error, got: %v", err)
	}
}

func TestInsertChildren(t *testing.T) {
	tests := []struct {
		name     string
		initial  []HyperNode
		children []any
		want     []HyperNode
	}{
		{
			name:     "no children",
			initial:  nil,
			children: nil,
			want:     nil,
		},
		{
			name:     "empty children",
			initial:  nil,
			children: []any{},
			want:     nil,
		},
		{
			name:     "single HyperNode",
			initial:  nil,
			children: []any{Text("a")},
			want:     []HyperNode{Text("a")},
		},
		{
			name:     "multiple HyperNodes",
			initial:  nil,
			children: []any{Text("a"), Text("b"), Text("c")},
			want:     []HyperNode{Text("a"), Text("b"), Text("c")},
		},
		{
			name:     "append to existing children",
			initial:  []HyperNode{Text("a")},
			children: []any{Text("b"), Text("c")},
			want:     []HyperNode{Text("a"), Text("b"), Text("c")},
		},
		{
			name:     "string converted to Text",
			initial:  nil,
			children: []any{"hello"},
			want:     []HyperNode{Text("hello")},
		},
		{
			name:     "multiple strings",
			initial:  nil,
			children: []any{"a", "b", "c"},
			want:     []HyperNode{Text("a"), Text("b"), Text("c")},
		},
		{
			name:     "fmt.Stringer converted to Text",
			initial:  nil,
			children: []any{stringerType("foo")},
			want:     []HyperNode{Text("foo")},
		},
		{
			name:     "int converted via fmt.Sprint",
			initial:  nil,
			children: []any{42},
			want:     []HyperNode{Text("42")},
		},
		{
			name:     "bool converted via fmt.Sprint",
			initial:  nil,
			children: []any{true},
			want:     []HyperNode{Text("true")},
		},
		{
			name:     "nil converted via fmt.Sprint",
			initial:  nil,
			children: []any{nil},
			want:     []HyperNode{Text("<nil>")},
		},
		{
			name:     "mixed types",
			initial:  []HyperNode{Text("a")},
			children: []any{Text("b"), "c", stringerType("d"), 42, true},
			want:     []HyperNode{Text("a"), Text("b"), Text("c"), Text("d"), Text("42"), Text("true")},
		},
		{
			name:     "triggers capacity growth",
			initial:  make([]HyperNode, 0, 2),
			children: []any{Text("a"), Text("b"), Text("c")},
			want:     []HyperNode{Text("a"), Text("b"), Text("c")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elem := &Element{Children: tt.initial}
			elem.InsertChildren(tt.children...)
			if len(elem.Children) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(elem.Children), len(tt.want))
			}
			for i := range tt.want {
				if elem.Children[i] != tt.want[i] {
					t.Errorf("Children[%d] = %#v, want %#v", i, elem.Children[i], tt.want[i])
				}
			}
		})
	}
}

func TestElement_Render(t *testing.T) {
	tests := []struct {
		name     string
		element  HyperNode
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple div",
			element:  DIV(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Div with single attribute",
			element:  DIV(AttrClass("container")),
			expected: `<div class="container"></div>`,
			wantErr:  false,
		},
		{
			name: "Div with text child (auto-escaped string)",
			element: func() HyperNode {
				return DIV("Hello World")
			}(),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name: "Div with multiple string children",
			element: func() HyperNode {
				return DIV("Hello", " ", "World")
			}(),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name: "Div with auto-escaped HTML string",
			element: func() HyperNode {
				return DIV("<script>alert('xss')</script>")
			}(),
			expected: "<div>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</div>",
			wantErr:  false,
		},
		{
			name: "Div with RawText (unescaped)",
			element: func() HyperNode {
				return DIV(RawText("<script>alert('xss')</script>"))
			}(),
			expected: "<div><script>alert('xss')</script></div>",
			wantErr:  false,
		},
		{
			name: "Nested elements with strings",
			element: func() HyperNode {
				return DIV(P("Hello"))
			}(),
			expected: "<div><p>Hello</p></div>",
			wantErr:  false,
		},
		{
			name:     "Void element (br)",
			element:  BR(),
			expected: "<br>",
			wantErr:  false,
		},
		{
			name:     "Void element with single attribute (img)",
			element:  IMG(AttrSrc("test.jpg")),
			expected: `<img src="test.jpg">`,
			wantErr:  false,
		},
		{
			name:     "Empty element",
			element:  Group(),
			expected: "",
			wantErr:  false,
		},
		{
			name: "Empty element with string children",
			element: func() HyperNode {
				return Group("Hello")
			}(),
			expected: "Hello",
			wantErr:  false,
		},
		{
			name:     "Boolean attribute true",
			element:  DIV(Attr("hidden", true)),
			expected: `<div hidden></div>`,
			wantErr:  false,
		},
		{
			name:     "Boolean attribute false",
			element:  DIV(Attr("hidden", false)),
			expected: `<div></div>`,
			wantErr:  false,
		},
		{
			name: "Div with integer (auto-converted)",
			element: func() HyperNode {
				return DIV(42)
			}(),
			expected: "<div>42</div>",
			wantErr:  false,
		},
		{
			name: "Div with boolean (auto-converted)",
			element: func() HyperNode {
				return DIV(true)
			}(),
			expected: "<div>true</div>",
			wantErr:  false,
		},
		{
			name: "Div with fmt.Stringer (auto-converted)",
			element: func() HyperNode {
				return DIV(stringerType("hello from stringer"))
			}(),
			expected: "<div>hello from stringer</div>",
			wantErr:  false,
		},
		{
			name: "Div with mixed types",
			element: func() HyperNode {
				return DIV("Count: ", 42, " Active: ", true)
			}(),
			expected: "<div>Count: 42 Active: true</div>",
			wantErr:  false,
		},
		{
			name: "Div with len() result (auto-converted)",
			element: func() HyperNode {
				items := []string{"a", "b", "c"}
				return DIV("Total: ", len(items))
			}(),
			expected: "<div>Total: 3</div>",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.element.Render(&buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Element.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && buf.String() != tt.expected {
				t.Errorf("Element.Render() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestElementAsHyperNode(t *testing.T) {
	tests := []struct {
		name     string
		node     HyperNode
		expected string
	}{
		{
			name:     "DIV without children",
			node:     DIV(),
			expected: "<div></div>",
		},
		{
			name:     "DIV with attributes, no children",
			node:     DIV(AttrClass("container"), AttrId("main")),
			expected: `<div class="container" id="main"></div>`,
		},
		{
			name:     "P without children",
			node:     P(),
			expected: "<p></p>",
		},
		{
			name:     "Element with attributes passed to Render directly",
			node:     H1(AttrId("title")),
			expected: `<h1 id="title"></h1>`,
		},
		{
			name:     "Elements nested as children directly",
			node:     DIV(SPAN(AttrClass("bold")), P()),
			expected: `<div><span class="bold"></span><p></p></div>`,
		},
		{
			name:     "Deeply nested elements",
			node:     DIV(DIV(DIV())),
			expected: "<div><div><div></div></div></div>",
		},
		{
			name:     "Void element used directly",
			node:     BR(),
			expected: "<br>",
		},
		{
			name:     "Mixed attributes and children in one call",
			node:     DIV(AttrClass("container"), "Hello", SPAN("World")),
			expected: `<div class="container">Hello<span>World</span></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Render(&buf, tt.node)
			if err != nil {
				t.Errorf("Render() returned error: %v", err)
				return
			}
			if buf.String() != tt.expected {
				t.Errorf("Render() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestElement_renderAttrs(t *testing.T) {
	tests := []struct {
		name      string
		attrs     []Attribute
		expected  string
		expectErr bool
	}{
		{
			name:      "Single string attribute",
			attrs:     []Attribute{AttrClass("test")},
			expected:  ` class="test"`,
			expectErr: false,
		},
		{
			name:      "Boolean attributes - true",
			attrs:     []Attribute{BooleanAttribute{Key: "hidden", IsActive: true}},
			expected:  ` hidden`,
			expectErr: false,
		},
		{
			name:      "Boolean attributes - false (not rendered)",
			attrs:     []Attribute{BooleanAttribute{Key: "disabled", IsActive: false}},
			expected:  "",
			expectErr: false,
		},
		{
			name:      "Empty key",
			attrs:     []Attribute{PairAttribute{Key: "", Value: "value"}},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Whitespace key",
			attrs:     []Attribute{PairAttribute{Key: "   ", Value: "value"}},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Key with HTML escaping",
			attrs:     []Attribute{PairAttribute{Key: "data-value", Value: "<script>"}},
			expected:  ` data-value="<script>"`,
			expectErr: false,
		},
		{
			name:      "Value with quotes",
			attrs:     []Attribute{PairAttribute{Key: "title", Value: `name is "Ahmad"`}},
			expected:  ` title="name is &quot;Ahmad&quot;"`,
			expectErr: false,
		},
		{
			name:      "Key needing escaping",
			attrs:     []Attribute{PairAttribute{Key: `"> <script>alert(1)</script>`, Value: "value"}},
			expected:  ` &#34;&gt; &lt;script&gt;alert(1)&lt;/script&gt;="value"`,
			expectErr: false,
		},
		{
			name:      "Boolean key needing escaping",
			attrs:     []Attribute{BooleanAttribute{Key: `"> <script>alert(1)</script>`, IsActive: true}},
			expected:  ` &#34;&gt; &lt;script&gt;alert(1)&lt;/script&gt;`,
			expectErr: false,
		},
		{
			name:      "Nil attribute is ignored",
			attrs:     []Attribute{nil},
			expected:  "",
			expectErr: false,
		},
		{
			name:      "Mix of valid and nil attributes",
			attrs:     []Attribute{AttrClass("test"), nil, AttrId("main")},
			expected:  ` class="test" id="main"`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			element := &Element{Attributes: tt.attrs}
			var buf bytes.Buffer
			err := element.renderAttrs(&buf)

			if (err != nil) != tt.expectErr {
				t.Errorf("renderAttrs() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr && buf.String() != tt.expected {
				t.Errorf("renderAttrs() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

func TestNewElement(t *testing.T) {
	t.Run("attributes and children routed", func(t *testing.T) {
		e := NewElement("div", AttrClass("container"), AttrId("main"), "Hello", 42)
		if e.Name != "div" {
			t.Errorf("Name = %q, want %q", e.Name, "div")
		}
		if len(e.Attributes) != 2 {
			t.Fatalf("len(Attributes) = %d, want 2", len(e.Attributes))
		}
		if len(e.Children) != 2 {
			t.Fatalf("len(Children) = %d, want 2", len(e.Children))
		}
		var buf bytes.Buffer
		if err := e.Render(&buf); err != nil {
			t.Fatal(err)
		}
		expected := `<div class="container" id="main">Hello42</div>`
		if buf.String() != expected {
			t.Errorf("Render() = %q, want %q", buf.String(), expected)
		}
	})

	t.Run("HyperNode children pass through", func(t *testing.T) {
		e := NewElement("div", P("inner"), Text("raw"))
		var buf bytes.Buffer
		if err := e.Render(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "<div><p>inner</p>raw</div>"
		if buf.String() != expected {
			t.Errorf("Render() = %q, want %q", buf.String(), expected)
		}
	})

	t.Run("fmt.Stringer converted via String", func(t *testing.T) {
		e := NewElement("span", stringerType("from stringer"))
		var buf bytes.Buffer
		if err := e.Render(&buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "<span>from stringer</span>" {
			t.Errorf("Render() = %q", buf.String())
		}
	})

	t.Run("no args", func(t *testing.T) {
		e := NewElement("div")
		if e.Name != "div" || len(e.Attributes) != 0 || len(e.Children) != 0 {
			t.Errorf("NewElement() = %#v, want empty div", e)
		}
	})
}

func TestNewVoidElement(t *testing.T) {
	t.Run("void without attributes", func(t *testing.T) {
		e := NewVoidElement("br")
		if !e.IsVoid {
			t.Error("IsVoid = false, want true")
		}
		if len(e.Children) != 0 {
			t.Errorf("len(Children) = %d, want 0", len(e.Children))
		}
		var buf bytes.Buffer
		if err := e.Render(&buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "<br>" {
			t.Errorf("Render() = %q, want %q", buf.String(), "<br>")
		}
	})

	t.Run("void with attributes", func(t *testing.T) {
		e := NewVoidElement("img", AttrSrc("a.png"), AttrAlt("pic"))
		var buf bytes.Buffer
		if err := e.Render(&buf); err != nil {
			t.Fatal(err)
		}
		expected := `<img src="a.png" alt="pic">`
		if buf.String() != expected {
			t.Errorf("Render() = %q, want %q", buf.String(), expected)
		}
	})
}

func TestInsertAttributes(t *testing.T) {
	t.Run("no attributes is a no-op", func(t *testing.T) {
		e := &Element{Name: "div"}
		e.InsertAttributes()
		if len(e.Attributes) != 0 {
			t.Errorf("len(Attributes) = %d, want 0", len(e.Attributes))
		}
	})

	t.Run("appends to existing attributes", func(t *testing.T) {
		e := &Element{Name: "div"}
		e.InsertAttributes(AttrClass("a"))
		e.InsertAttributes(AttrId("b"), AttrTitle("c"))
		if len(e.Attributes) != 3 {
			t.Fatalf("len(Attributes) = %d, want 3", len(e.Attributes))
		}
		var buf bytes.Buffer
		if err := e.Render(&buf); err != nil {
			t.Fatal(err)
		}
		expected := `<div class="a" id="b" title="c"></div>`
		if buf.String() != expected {
			t.Errorf("Render() = %q, want %q", buf.String(), expected)
		}
	})

	t.Run("triggers capacity growth", func(t *testing.T) {
		e := &Element{Name: "div", Attributes: make([]Attribute, 0, 1)}
		e.InsertAttributes(AttrClass("a"), AttrId("b"), AttrTitle("c"))
		if len(e.Attributes) != 3 {
			t.Fatalf("len(Attributes) = %d, want 3", len(e.Attributes))
		}
	})
}
