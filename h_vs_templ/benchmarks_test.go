package h_vs_templ

import (
	"bytes"
	"testing"

	"github.com/assaidy/h"
)

// getBenchmarkData returns sample user data for benchmarks
func getBenchmarkData() []User {
	return []User{
		{Name: "Alice", Admin: true},
		{Name: "Bob", Admin: false},
		{Name: "Charlie", Admin: false},
		{Name: "Diana", Admin: true},
		{Name: "Eve", Admin: false},
	}
}

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

func BenchmarkSimpleElement_H(b *testing.B) {
	page := h.Div("Hello World")
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkDeepNesting_H(b *testing.B) {
	page := h.Div(
		h.Div(
			h.Div(
				h.Div(
					h.Div(
						h.P("Deep content"),
					),
				),
			),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkManyAttributes_H(b *testing.B) {
	page := h.Div(h.KV{
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
		h.Render(&buf, page)
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

func BenchmarkLargeText_H(b *testing.B) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	page := h.P(text)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkList10_H(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	page := h.Ul(
		h.MapSlice(items, func(s string) h.Node {
			return h.Li(s)
		}),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkList100_H(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = "item"
	}
	page := h.Ul(
		h.MapSlice(items, func(s string) h.Node {
			return h.Li(s)
		}),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkConditionals_H(b *testing.B) {
	page := h.Div(
		h.If(true, h.Span("First")),
		h.If(false, h.Span("Second")),
		h.If(true, h.Span("Third")),
		h.IfElse(true, h.Strong("True"), h.Em("False")),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkMixedContent_H(b *testing.B) {
	page := h.Div(
		h.H1("Title"),
		h.P("Paragraph with ", h.Strong("bold"), " and ", h.Em("italic"), " text."),
		h.Ul(
			h.Li("Item 1"),
			h.Li(h.A(h.KV{"href": "#"}, "Link")),
		),
		h.Div(h.KV{"class": "footer"},
			h.Small("Copyright 2024"),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkVoidElements_H(b *testing.B) {
	page := h.Div(
		h.Img(h.KV{"src": "image.jpg", "alt": "Image"}),
		h.Br(),
		h.Hr(),
		h.Input(h.KV{"type": "text", "value": "input"}),
		h.Meta(h.KV{"charset": "UTF-8"}),
		h.Link(h.KV{"rel": "stylesheet", "href": "style.css"}),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkHTMLEscaping_H(b *testing.B) {
	content := "<script>alert('xss')</script> & more <b>bold</b>"
	page := h.Div(content)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkTable_H(b *testing.B) {
	rows := 10
	page := h.Table(
		h.Thead(
			h.Tr(
				h.Th("Name"),
				h.Th("Value"),
				h.Th("Action"),
			),
		),
		h.Tbody(
			h.Repeat(rows, func() h.Node {
				return h.Tr(
					h.Td("Cell 1"),
					h.Td("Cell 2"),
					h.Td(h.Button("Click")),
				)
			}),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkForm_H(b *testing.B) {
	page := h.Form(h.KV{"action": "/submit", "method": "POST"},
		h.Fieldset(
			h.Legend("User Form"),
			h.Label(h.KV{"for": "name"}, "Name:"),
			h.Input(h.KV{"type": "text", "id": "name", "name": "name"}),
			h.Br(),
			h.Label(h.KV{"for": "email"}, "Email:"),
			h.Input(h.KV{"type": "email", "id": "email", "name": "email"}),
			h.Br(),
			h.Button(h.KV{"type": "submit"}, "Submit"),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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
		RealWorldTempl(users).Render(ctx, &buf)
	}
}

func BenchmarkRealWorld_H(b *testing.B) {
	users := getBenchmarkData()
	page := h.Empty(
		h.DoctypeHTML(),
		h.Html(
			h.Head(
				h.Meta(h.KV{"charset": "UTF-8"}),
				h.Title(h.RawHTML("User Dashboard")),
				h.Link(h.KV{"rel": "stylesheet", "href": "/style.css"}),
			),
			h.Body(
				h.Header(
					h.H1("User Dashboard"),
					h.Nav(
						h.A(h.KV{"href": "/"}, "Home"),
						h.A(h.KV{"href": "/users"}, "Users"),
						h.A(h.KV{"href": "/settings"}, "Settings"),
					),
				),
				h.Main(
					h.H2("Users"),
					h.If(len(users) > 0,
						h.Table(
							h.Thead(
								h.Tr(
									h.Th("Name"),
									h.Th("Role"),
								),
							),
							h.Tbody(
								h.MapSlice(users, func(u User) h.Node {
									return h.Tr(
										h.Td(u.Name),
										h.Td(h.IfElse(u.Admin, h.Strong("Admin"), h.Span("User"))),
									)
								}),
							),
						),
					),
					h.If(len(users) == 0, h.P("No users found.")),
				),
				h.Footer(
					h.P("© 2024 Company"),
				),
			),
		),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkEmptyPage_H(b *testing.B) {
	page := h.Html(h.Body())
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkRawHTML_H(b *testing.B) {
	html := "<div><span>Content</span></div>"
	page := h.Div(h.RawHTML(html))
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

func BenchmarkRegularString_H(b *testing.B) {
	text := "<div><span>Content</span></div>"
	page := h.Div(text)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
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

func BenchmarkSVG_H(b *testing.B) {
	// NOTE: In real example all the SVG tag is copied and put inside RawHTML.
	page := h.Svg(h.KV{"width": "100", "height": "100"},
		h.RawHTML(`<circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />`),
	)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// ============================================================================
// CONCURRENT BENCHMARKS: Real Server Load Simulation
// These benchmarks simulate real-world server scenarios where multiple requests
// are processed concurrently. This gives you the actual requests/second your
// library can handle under load.
// ============================================================================

// buildRealWorldPage creates a realistic full page using many library utilities
func buildRealWorldPage(users []User) h.Node {
	return h.Empty(
		h.DoctypeHTML(),
		h.Html(
			h.Head(
				h.Meta(h.KV{"charset": "UTF-8"}),
				h.Meta(h.KV{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}),
				h.Title(h.RawHTML("Dashboard - User Management")),
				h.Link(h.KV{"rel": "stylesheet", "href": "/css/main.css"}),
				h.Link(h.KV{"rel": "icon", "href": "/favicon.ico"}),
			),
			h.Body(
				h.Header(
					h.KV{"class": "site-header"},
					h.Nav(
						h.KV{"class": "main-nav"},
						h.A(h.KV{"href": "/", "class": "nav-link"}, "Home"),
						h.A(h.KV{"href": "/users", "class": "nav-link active"}, "Users"),
						h.A(h.KV{"href": "/settings", "class": "nav-link"}, "Settings"),
						h.A(h.KV{"href": "/logout", "class": "nav-link"}, "Logout"),
					),
				),
				h.Main(
					h.KV{"class": "main-content"},
					h.H1("User Management Dashboard"),
					h.P("Welcome to the admin dashboard. Manage users and permissions below."),
					h.If(len(users) > 0,
						h.Section(
							h.KV{"class": "users-section"},
							h.H2("Active Users"),
							h.Table(
								h.KV{"class": "users-table"},
								h.Thead(
									h.Tr(
										h.Th("ID"),
										h.Th("Name"),
										h.Th("Role"),
										h.Th("Status"),
										h.Th("Actions"),
									),
								),
								h.Tbody(
									h.MapSlice(users, func(u User) h.Node {
										return h.Tr(
											h.Td(h.Strong("#")),
											h.Td(u.Name),
											h.Td(h.IfElse(u.Admin,
												h.Span(h.KV{"class": "badge admin"}, "Administrator"),
												h.Span(h.KV{"class": "badge user"}, "User"),
											)),
											h.Td(h.Span(h.KV{"class": "status active"}, "Active")),
											h.Td(
												h.Button(h.KV{"class": "btn-edit"}, "Edit"),
												h.Button(h.KV{"class": "btn-delete"}, "Delete"),
											),
										)
									}),
								),
							),
						),
					),
					h.If(len(users) == 0,
						h.Div(
							h.KV{"class": "empty-state"},
							h.P("No users found. Add your first user to get started."),
						),
					),
					h.Section(
						h.KV{"class": "quick-stats"},
						h.H3("Quick Stats"),
						h.Div(
							h.KV{"class": "stats-grid"},
							h.Div(
								h.KV{"class": "stat-card"},
								h.Strong(len(users)),
								h.Span("Total Users"),
							),
							h.Div(
								h.KV{"class": "stat-card"},
								h.Strong(h.IfElse(len(users) > 0, len(users), 0)),
								h.Span("Active Now"),
							),
						),
					),
				),
				h.Footer(
					h.KV{"class": "site-footer"},
					h.P("2025 Company Inc. All rights reserved."),
				),
			),
		),
	)
}

// BenchmarkSequential_RealWorld_H measures single-threaded (sequential) performance
// This is what you'd see in the standard benchmarks - operations per second on one CPU
func BenchmarkSequential_RealWorld_H(b *testing.B) {
	users := getBenchmarkData()
	page := buildRealWorldPage(users)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// BenchmarkConcurrent10_RealWorld_H simulates 10 concurrent requests
// Typical for a small web application under light load
func BenchmarkConcurrent10_RealWorld_H(b *testing.B) {
	users := getBenchmarkData()
	page := buildRealWorldPage(users)
	b.SetParallelism(10)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			h.Render(&buf, page)
		}
	})
}

// BenchmarkConcurrent100_RealWorld_H simulates 100 concurrent requests
// Typical for a medium-traffic web application
func BenchmarkConcurrent100_RealWorld_H(b *testing.B) {
	users := getBenchmarkData()
	page := buildRealWorldPage(users)
	b.SetParallelism(100)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			h.Render(&buf, page)
		}
	})
}

// BenchmarkConcurrent1000_RealWorld_H simulates 1000 concurrent requests
// High-load scenario - stress test for the library
func BenchmarkConcurrent1000_RealWorld_H(b *testing.B) {
	users := getBenchmarkData()
	page := buildRealWorldPage(users)
	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			h.Render(&buf, page)
		}
	})
}

// BenchmarkConcurrentRealistic_RealWorld_H uses GOMAXPROCS goroutines
// This represents the most realistic server scenario where concurrency
// matches available CPU cores (what real servers typically use)
func BenchmarkConcurrentRealistic_RealWorld_H(b *testing.B) {
	users := getBenchmarkData()
	page := buildRealWorldPage(users)
	// No SetParallelism - uses default GOMAXPROCS
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			h.Render(&buf, page)
		}
	})
}

// ============================================================================
// CONCURRENT BENCHMARKS: Templ Comparison
// Matching benchmarks for templ to compare with H library
// ============================================================================

// BenchmarkSequential_RealWorld_Templ measures templ single-threaded performance
func BenchmarkSequential_RealWorld_Templ(b *testing.B) {
	users := getBenchmarkData()
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		RealWorldTempl(users).Render(ctx, &buf)
	}
}

// BenchmarkConcurrent10_RealWorld_Templ simulates 10 concurrent requests with templ
func BenchmarkConcurrent10_RealWorld_Templ(b *testing.B) {
	users := getBenchmarkData()
	ctx := b.Context()
	b.SetParallelism(10)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			RealWorldTempl(users).Render(ctx, &buf)
		}
	})
}

// BenchmarkConcurrent100_RealWorld_Templ simulates 100 concurrent requests with templ
func BenchmarkConcurrent100_RealWorld_Templ(b *testing.B) {
	users := getBenchmarkData()
	ctx := b.Context()
	b.SetParallelism(100)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			RealWorldTempl(users).Render(ctx, &buf)
		}
	})
}

// BenchmarkConcurrent1000_RealWorld_Templ simulates 1000 concurrent requests with templ
func BenchmarkConcurrent1000_RealWorld_Templ(b *testing.B) {
	users := getBenchmarkData()
	ctx := b.Context()
	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			RealWorldTempl(users).Render(ctx, &buf)
		}
	})
}

// BenchmarkConcurrentRealistic_RealWorld_Templ uses GOMAXPROCS goroutines with templ
func BenchmarkConcurrentRealistic_RealWorld_Templ(b *testing.B) {
	users := getBenchmarkData()
	ctx := b.Context()
	// No SetParallelism - uses default GOMAXPROCS
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			RealWorldTempl(users).Render(ctx, &buf)
		}
	})
}
