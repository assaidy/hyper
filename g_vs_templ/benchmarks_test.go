package g_vs_templ

import (
	"bytes"
	"testing"

	"github.com/assaidy/g"
)

// ============================================================================
// BENCHMARK 1: Simple Element
// Single div with text content
// ============================================================================

func BenchmarkSimpleElement_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		SimpleElementTempl().Render(ctx, &buf)
	}
}

func BenchmarkSimpleElement_G(b *testing.B) {
	page := gg.Div("Hello World")
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 2: Deep Nesting
// Multiple levels of nested elements
// ============================================================================

func BenchmarkDeepNesting_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		DeepNestingTempl().Render(ctx, &buf)
	}
}

func BenchmarkDeepNesting_G(b *testing.B) {
	page := gg.Div(
		gg.Div(
			gg.Div(
				gg.Div(
					gg.Div(
						gg.P("Deep content"),
					),
				),
			),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 3: Many Attributes
// Element with many attributes
// ============================================================================

func BenchmarkManyAttributes_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		ManyAttributesTempl().Render(ctx, &buf)
	}
}

func BenchmarkManyAttributes_G(b *testing.B) {
	page := gg.Div(gg.KV{
		"id":         "main",
		"class":      "container wrapper",
		"data-role":  "content",
		"data-value": "12345",
		"aria-label": "Main content",
		"hidden":     true,
		"disabled":   false,
	})
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 4: Large Text Content
// Element with large text content
// ============================================================================

func BenchmarkLargeText_Templ(b *testing.B) {
	ctx := b.Context()
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		LargeTextTempl(text).Render(ctx, &buf)
	}
}

func BenchmarkLargeText_G(b *testing.B) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	page := gg.P(text)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 5: List Iteration (10 items)
// Rendering a list with 10 items
// ============================================================================

func BenchmarkList10_Templ(b *testing.B) {
	ctx := b.Context()
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		ListTempl(items).Render(ctx, &buf)
	}
}

func BenchmarkList10_G(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	page := gg.Ul(
		gg.MapSlice(items, func(s string) gg.Node {
			return gg.Li(s)
		}),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 6: List Iteration (100 items)
// Rendering a list with 100 items
// ============================================================================

func BenchmarkList100_Templ(b *testing.B) {
	ctx := b.Context()
	items := make([]string, 100)
	for i := range items {
		items[i] = "item"
	}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		ListTempl(items).Render(ctx, &buf)
	}
}

func BenchmarkList100_G(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = "item"
	}
	page := gg.Ul(
		gg.MapSlice(items, func(s string) gg.Node {
			return gg.Li(s)
		}),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 7: Complex Conditionals
// Multiple conditional branches
// ============================================================================

func BenchmarkConditionals_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		ConditionalsTempl(true, false, true).Render(ctx, &buf)
	}
}

func BenchmarkConditionals_G(b *testing.B) {
	page := gg.Div(
		gg.If(true, gg.Span("First")),
		gg.If(false, gg.Span("Second")),
		gg.If(true, gg.Span("Third")),
		gg.IfElse(true, gg.Strong("True"), gg.Em("False")),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 8: Mixed Content
// Various element types and content
// ============================================================================

func BenchmarkMixedContent_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		MixedContentTempl().Render(ctx, &buf)
	}
}

func BenchmarkMixedContent_G(b *testing.B) {
	page := gg.Div(
		gg.H1("Title"),
		gg.P("Paragraph with ", gg.Strong("bold"), " and ", gg.Em("italic"), " text."),
		gg.Ul(
			gg.Li("Item 1"),
			gg.Li(gg.A(gg.KV{"href": "#"}, "Link")),
		),
		gg.Div(gg.KV{"class": "footer"},
			gg.Small("Copyright 2024"),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 9: Void Elements
// Multiple void elements (img, br, hr, input)
// ============================================================================

func BenchmarkVoidElements_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		VoidElementsTempl().Render(ctx, &buf)
	}
}

func BenchmarkVoidElements_G(b *testing.B) {
	page := gg.Div(
		gg.Img(gg.KV{"src": "image.jpg", "alt": "Image"}),
		gg.Br(),
		gg.Hr(),
		gg.Input(gg.KV{"type": "text", "value": "input"}),
		gg.Meta(gg.KV{"charset": "UTF-8"}),
		gg.Link(gg.KV{"rel": "stylesheet", "href": "style.css"}),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 10: HTML Escaping
// Content that requires HTML escaping
// ============================================================================

func BenchmarkHTMLEscaping_Templ(b *testing.B) {
	ctx := b.Context()
	content := "<script>alert('xss')</script> & more <b>bold</b>"
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		HTMLEscapingTempl(content).Render(ctx, &buf)
	}
}

func BenchmarkHTMLEscaping_G(b *testing.B) {
	content := "<script>alert('xss')</script> & more <b>bold</b>"
	page := gg.Div(content)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 11: Table (10x3)
// Complex table structure
// ============================================================================

func BenchmarkTable_Templ(b *testing.B) {
	ctx := b.Context()
	rows := 10
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		TableTempl(rows).Render(ctx, &buf)
	}
}

func BenchmarkTable_G(b *testing.B) {
	rows := 10
	page := gg.Table(
		gg.Thead(
			gg.Tr(
				gg.Th("Name"),
				gg.Th("Value"),
				gg.Th("Action"),
			),
		),
		gg.Tbody(
			gg.Repeat(rows, func() gg.Node {
				return gg.Tr(
					gg.Td("Cell 1"),
					gg.Td("Cell 2"),
					gg.Td(gg.Button("Click")),
				)
			}),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 12: Form Elements
// Complete form with various input types
// ============================================================================

func BenchmarkForm_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		FormTempl().Render(ctx, &buf)
	}
}

func BenchmarkForm_G(b *testing.B) {
	page := gg.Form(gg.KV{"action": "/submit", "method": "POST"},
		gg.Fieldset(
			gg.Legend("User Form"),
			gg.Label(gg.KV{"for": "name"}, "Name:"),
			gg.Input(gg.KV{"type": "text", "id": "name", "name": "name"}),
			gg.Br(),
			gg.Label(gg.KV{"for": "email"}, "Email:"),
			gg.Input(gg.KV{"type": "email", "id": "email", "name": "email"}),
			gg.Br(),
			gg.Button(gg.KV{"type": "submit"}, "Submit"),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 13: Real World Page
// Combination of all features
// ============================================================================

func BenchmarkRealWorld_Templ(b *testing.B) {
	ctx := b.Context()
	users := getBenchmarkData()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		Page(users).Render(ctx, &buf)
	}
}

func BenchmarkRealWorld_G(b *testing.B) {
	users := getBenchmarkData()
	page := gg.Empty(
		gg.DoctypeHTML(),
		gg.Html(
			gg.Head(
				gg.Meta(gg.KV{"charset": "UTF-8"}),
				gg.Title(gg.RawHTML("User Dashboard")),
				gg.Link(gg.KV{"rel": "stylesheet", "href": "/style.css"}),
			),
			gg.Body(
				gg.Header(
					gg.H1("User Dashboard"),
					gg.Nav(
						gg.A(gg.KV{"href": "/"}, "Home"),
						gg.A(gg.KV{"href": "/users"}, "Users"),
						gg.A(gg.KV{"href": "/settings"}, "Settings"),
					),
				),
				gg.Main(
					gg.H2("Users"),
					gg.If(len(users) > 0,
						gg.Table(
							gg.Thead(
								gg.Tr(
									gg.Th("Name"),
									gg.Th("Role"),
								),
							),
							gg.Tbody(
								gg.MapSlice(users, func(u User) gg.Node {
									return gg.Tr(
										gg.Td(u.Name),
										gg.Td(gg.IfElse(u.Admin, gg.Strong("Admin"), gg.Span("User"))),
									)
								}),
							),
						),
					),
					gg.If(len(users) == 0, gg.P("No users found.")),
				),
				gg.Footer(
					gg.P("© 2024 Company"),
				),
			),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 14: Empty Page
// Minimal page structure
// ============================================================================

func BenchmarkEmptyPage_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		EmptyPageTempl().Render(ctx, &buf)
	}
}

func BenchmarkEmptyPage_G(b *testing.B) {
	page := gg.Html(gg.Body())
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 15: RawHTML vs String
// Using RawHTML (no escaping) vs regular string (with escaping)
// ============================================================================

func BenchmarkRawHTML_Templ(b *testing.B) {
	ctx := b.Context()
	html := "<div><span>Content</span></div>"
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		RawHTMLTempl(html).Render(ctx, &buf)
	}
}

func BenchmarkRawHTML_G(b *testing.B) {
	html := "<div><span>Content</span></div>"
	page := gg.Div(gg.RawHTML(html))
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

func BenchmarkRegularString_G(b *testing.B) {
	text := "<div><span>Content</span></div>"
	page := gg.Div(text)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 16: SVG Content
// SVG graphics rendering
// ============================================================================

func BenchmarkSVG_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		SVGTempl().Render(ctx, &buf)
	}
}

func BenchmarkSVG_G(b *testing.B) {
	// NOTE: In real example all the SVG tag is copied and put inside RawHTML.
	page := gg.Svg(gg.KV{"width": "100", "height": "100"},
		gg.RawHTML(`<circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />`),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		gg.Render(&buf, page)
	}
}
