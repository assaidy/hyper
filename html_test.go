package h

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

func TestElement_Render(t *testing.T) {
	tests := []struct {
		name     string
		element  Node
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple div",
			element:  Div(),
			expected: "<div></div>",
			wantErr:  false,
		},
		{
			name:     "Div with single attribute",
			element:  Div(KV{"class": "container"}),
			expected: `<div class="container"></div>`,
			wantErr:  false,
		},
		{
			name: "Div with text child (auto-escaped string)",
			element: func() Node {
				return Div("Hello World")
			}(),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name: "Div with multiple string children",
			element: func() Node {
				return Div("Hello", " ", "World")
			}(),
			expected: "<div>Hello World</div>",
			wantErr:  false,
		},
		{
			name: "Div with auto-escaped HTML string",
			element: func() Node {
				return Div("<script>alert('xss')</script>")
			}(),
			expected: "<div>&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;</div>",
			wantErr:  false,
		},
		{
			name: "Div with RawText (unescaped)",
			element: func() Node {
				return Div(RawText("<script>alert('xss')</script>"))
			}(),
			expected: "<div><script>alert('xss')</script></div>",
			wantErr:  false,
		},
		{
			name: "Nested elements with strings",
			element: func() Node {
				return Div(P("Hello"))
			}(),
			expected: "<div><p>Hello</p></div>",
			wantErr:  false,
		},
		{
			name:     "Void element (br)",
			element:  Br(),
			expected: "<br>",
			wantErr:  false,
		},
		{
			name:     "Void element with single attribute (img)",
			element:  Img(KV{"src": "test.jpg"}),
			expected: `<img src="test.jpg">`,
			wantErr:  false,
		},
		{
			name:     "Empty element",
			element:  Empty(),
			expected: "",
			wantErr:  false,
		},
		{
			name: "Empty element with string children",
			element: func() Node {
				return Empty("Hello")
			}(),
			expected: "Hello",
			wantErr:  false,
		},
		{
			name:     "Boolean attribute true",
			element:  Div(KV{"hidden": true}),
			expected: `<div hidden></div>`,
			wantErr:  false,
		},
		{
			name:     "Boolean attribute false",
			element:  Div(KV{"hidden": false}),
			expected: `<div></div>`,
			wantErr:  false,
		},
		{
			name: "Div with integer (auto-converted)",
			element: func() Node {
				return Div(42)
			}(),
			expected: "<div>42</div>",
			wantErr:  false,
		},
		{
			name: "Div with boolean (auto-converted)",
			element: func() Node {
				return Div(true)
			}(),
			expected: "<div>true</div>",
			wantErr:  false,
		},
		{
			name: "Div with fmt.Stringer (auto-converted)",
			element: func() Node {
				return Div(stringerType("hello from stringer"))
			}(),
			expected: "<div>hello from stringer</div>",
			wantErr:  false,
		},
		{
			name: "Div with mixed types",
			element: func() Node {
				return Div("Count: ", 42, " Active: ", true)
			}(),
			expected: "<div>Count: 42 Active: true</div>",
			wantErr:  false,
		},
		{
			name: "Div with len() result (auto-converted)",
			element: func() Node {
				items := []string{"a", "b", "c"}
				return Div("Total: ", len(items))
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

func TestElement_renderAttrs(t *testing.T) {
	tests := []struct {
		name      string
		attrs     KV
		expected  string
		expectErr bool
	}{
		{
			name:      "Single string attribute",
			attrs:     KV{"class": "test"},
			expected:  ` class="test"`,
			expectErr: false,
		},
		{
			name:      "Boolean attributes",
			attrs:     KV{"hidden": true, "disabled": false},
			expected:  ` hidden`,
			expectErr: false,
		},
		{
			name:      "Nil value",
			attrs:     KV{"test": nil},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid value type",
			attrs:     KV{"test": 123},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Empty key",
			attrs:     KV{"": "value"},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Whitespace key",
			attrs:     KV{"   ": "value"},
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Key with HTML escaping",
			attrs:     KV{"data-value": "<script>"},
			expected:  ` data-value="<script>"`,
			expectErr: false,
		},
		{
			name:      "Value with quotes",
			attrs:     KV{"title": `name is "Ahmad"`},
			expected:  ` title="name is &quot;Ahmad&quot;"`,
			expectErr: false,
		},
		{
			name:      "Key needing escaping",
			attrs:     KV{`"> <script>alert(1)</script>`: "value"},
			expected:  ` &#34;&gt; &lt;script&gt;alert(1)&lt;/script&gt;="value"`,
			expectErr: false,
		},
		{
			name:      "Boolean key needing escaping",
			attrs:     KV{`"> <script>alert(1)</script>`: true},
			expected:  ` &#34;&gt; &lt;script&gt;alert(1)&lt;/script&gt;`,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			element := &Element{}
			element.fillAttrsWithKV(tt.attrs)
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
