package h_vs_templ

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/assaidy/hyper/v2"
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

func BenchmarkSimpleElement_Hyper(b *testing.B) {
	page := h.DIV("Hello World")
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

func BenchmarkDeepNesting_Hyper(b *testing.B) {
	page := h.DIV(
		h.DIV(
			h.DIV(
				h.DIV(
					h.DIV(
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

func BenchmarkManyAttributes_Hyper(b *testing.B) {
	page := h.DIV(h.KV{
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

func BenchmarkLargeText_Hyper(b *testing.B) {
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

func BenchmarkList10_Hyper(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	page := h.UL(
		h.Range(items, func(s string) h.HyperNode {
			return h.LI(s)
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

func BenchmarkList100_Hyper(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = "item"
	}
	page := h.UL(
		h.Range(items, func(s string) h.HyperNode {
			return h.LI(s)
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

func BenchmarkConditionals_Hyper(b *testing.B) {
	page := h.DIV(
		h.If(true, h.SPAN("First")),
		h.If(false, h.SPAN("Second")),
		h.If(true, h.SPAN("Third")),
		h.IfElse(true, h.STRONG("True"), h.EM("False")),
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

func BenchmarkMixedContent_Hyper(b *testing.B) {
	page := h.DIV(
		h.H1("Title"),
		h.P("Paragraph with ", h.STRONG("bold"), " and ", h.EM("italic"), " text."),
		h.UL(
			h.LI("Item 1"),
			h.LI(h.A(h.KV{"href": "#"}, "Link")),
		),
		h.DIV(h.KV{"class": "footer"},
			h.SMALL("Copyright 2024"),
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

func BenchmarkVoidElements_Hyper(b *testing.B) {
	page := h.DIV(
		h.IMG(h.KV{"src": "image.jpg", "alt": "Image"}),
		h.BR(),
		h.HR(),
		h.INPUT(h.KV{"type": "text", "value": "input"}),
		h.META(h.KV{"charset": "UTF-8"}),
		h.LINK(h.KV{"rel": "stylesheet", "href": "style.css"}),
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

func BenchmarkHTMLEscaping_Hyper(b *testing.B) {
	content := "<script>alert('xss')</script> & more <b>bold</b>"
	page := h.DIV(content)
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

func BenchmarkTable_Hyper(b *testing.B) {
	rows := 10
	page := h.TABLE(
		h.THEAD(
			h.TR(
				h.TH("Name"),
				h.TH("Value"),
				h.TH("Action"),
			),
		),
		h.TBODY(
			h.Repeat(rows, func() h.HyperNode {
				return h.TR(
					h.TD("Cell 1"),
					h.TD("Cell 2"),
					h.TD(h.BUTTON("Click")),
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

func BenchmarkForm_Hyper(b *testing.B) {
	page := h.FORM(h.KV{"action": "/submit", "method": "POST"},
		h.FIELDSET(
			h.LEGEND("User Form"),
			h.LABEL(h.KV{"for": "name"}, "Name:"),
			h.INPUT(h.KV{"type": "text", "id": "name", "name": "name"}),
			h.BR(),
			h.LABEL(h.KV{"for": "email"}, "Email:"),
			h.INPUT(h.KV{"type": "email", "id": "email", "name": "email"}),
			h.BR(),
			h.BUTTON(h.KV{"type": "submit"}, "Submit"),
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

func BenchmarkRealWorld_Hyper(b *testing.B) {
	users := getBenchmarkData()
	page := h.EMPTY(
		h.DOCTYPE(),
		h.HTML(
			h.HEAD(
				h.META(h.KV{"charset": "UTF-8"}),
				h.TITLE(h.RawText("User Dashboard")),
				h.LINK(h.KV{"rel": "stylesheet", "href": "/style.css"}),
			),
			h.BODY(
				h.HEADER(
					h.H1("User Dashboard"),
					h.NAV(
						h.A(h.KV{"href": "/"}, "Home"),
						h.A(h.KV{"href": "/users"}, "Users"),
						h.A(h.KV{"href": "/settings"}, "Settings"),
					),
				),
				h.MAIN(
					h.H2("Users"),
					h.If(len(users) > 0,
						h.TABLE(
							h.THEAD(
								h.TR(
									h.TH("Name"),
									h.TH("Role"),
								),
							),
							h.TBODY(
								h.Range(users, func(u User) h.HyperNode {
									return h.TR(
										h.TD(u.Name),
										h.TD(h.IfElse(u.Admin, h.STRONG("Admin"), h.SPAN("User"))),
									)
								}),
							),
						),
					),
					h.If(len(users) == 0, h.P("No users found.")),
				),
				h.FOOTER(
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

func BenchmarkEmptyPage_Hyper(b *testing.B) {
	page := h.HTML(h.BODY())
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 15: RawText vs String
// Using RawText (no escaping) vs regular string (with escaping)
// ============================================================================

func BenchmarkRawText_Templ(b *testing.B) {
	ctx := b.Context()
	html := "<div><span>Content</span></div>"
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		RawTextTempl(html).Render(ctx, &buf)
	}
}

func BenchmarkRawText_Hyper(b *testing.B) {
	html := "<div><span>Content</span></div>"
	page := h.DIV(h.RawText(html))
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

func BenchmarkRegularString_Hyper(b *testing.B) {
	text := "<div><span>Content</span></div>"
	page := h.DIV(text)
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

func BenchmarkSVG_Hyper(b *testing.B) {
	// NOTE: In real example all the SVG tag is copied and put inside RawText.
	page := h.SVG(h.KV{"width": "100", "height": "100"},
		h.RawText(`<circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />`),
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
func buildRealWorldPage(users []User) h.HyperNode {
	return h.EMPTY(
		h.DOCTYPE(),
		h.HTML(
			h.HEAD(
				h.META(h.KV{"charset": "UTF-8"}),
				h.META(h.KV{"name": "viewport", "content": "width=device-width, initial-scale=1.0"}),
				h.TITLE(h.RawText("Dashboard - User Management")),
				h.LINK(h.KV{"rel": "stylesheet", "href": "/css/main.css"}),
				h.LINK(h.KV{"rel": "icon", "href": "/favicon.ico"}),
			),
			h.BODY(
				h.HEADER(
					h.KV{"class": "site-header"},
					h.NAV(
						h.KV{"class": "main-nav"},
						h.A(h.KV{"href": "/", "class": "nav-link"}, "Home"),
						h.A(h.KV{"href": "/users", "class": "nav-link active"}, "Users"),
						h.A(h.KV{"href": "/settings", "class": "nav-link"}, "Settings"),
						h.A(h.KV{"href": "/logout", "class": "nav-link"}, "Logout"),
					),
				),
				h.MAIN(
					h.KV{"class": "main-content"},
					h.H1("User Management Dashboard"),
					h.P("Welcome to the admin dashboard. Manage users and permissions below."),
					h.If(len(users) > 0,
						h.SECTION(
							h.KV{"class": "users-section"},
							h.H2("Active Users"),
							h.TABLE(
								h.KV{"class": "users-table"},
								h.THEAD(
									h.TR(
										h.TH("ID"),
										h.TH("Name"),
										h.TH("Role"),
										h.TH("Status"),
										h.TH("Actions"),
									),
								),
								h.TBODY(
									h.Range(users, func(u User) h.HyperNode {
										return h.TR(
											h.TD(h.STRONG("#")),
											h.TD(u.Name),
											h.TD(h.IfElse(u.Admin,
												h.SPAN(h.KV{"class": "badge admin"}, "Administrator"),
												h.SPAN(h.KV{"class": "badge user"}, "User"),
											)),
											h.TD(h.SPAN(h.KV{"class": "status active"}, "Active")),
											h.TD(
												h.BUTTON(h.KV{"class": "btn-edit"}, "Edit"),
												h.BUTTON(h.KV{"class": "btn-delete"}, "Delete"),
											),
										)
									}),
								),
							),
						),
					),
					h.If(len(users) == 0,
						h.DIV(
							h.KV{"class": "empty-state"},
							h.P("No users found. Add your first user to get started."),
						),
					),
					h.SECTION(
						h.KV{"class": "quick-stats"},
						h.H3("Quick Stats"),
						h.DIV(
							h.KV{"class": "stats-grid"},
							h.DIV(
								h.KV{"class": "stat-card"},
								h.STRONG(len(users)),
								h.SPAN("Total Users"),
							),
							h.DIV(
								h.KV{"class": "stat-card"},
								h.STRONG(h.IfElse(len(users) > 0, len(users), 0)),
								h.SPAN("Active Now"),
							),
						),
					),
				),
				h.FOOTER(
					h.KV{"class": "site-footer"},
					h.P("2025 Company Inc. All rights reserved."),
				),
			),
		),
	)
}

// BenchmarkSequential_RealWorld_Hyper measures single-threaded (sequential) performance
// This is what you'd see in the standard benchmarks - operations per second on one CPU
func BenchmarkSequential_RealWorld_Hyper(b *testing.B) {
	users := getBenchmarkData()
	page := buildRealWorldPage(users)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// BenchmarkConcurrent10_RealWorld_Hyper simulates 10 concurrent requests
// Typical for a small web application under light load
func BenchmarkConcurrent10_RealWorld_Hyper(b *testing.B) {
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

// BenchmarkConcurrent100_RealWorld_Hyper simulates 100 concurrent requests
// Typical for a medium-traffic web application
func BenchmarkConcurrent100_RealWorld_Hyper(b *testing.B) {
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

// BenchmarkConcurrent1000_RealWorld_Hyper simulates 1000 concurrent requests
// High-load scenario - stress test for the library
func BenchmarkConcurrent1000_RealWorld_Hyper(b *testing.B) {
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

// BenchmarkConcurrentRealistic_RealWorld_Hyper uses GOMAXPROCS goroutines
// This represents the most realistic server scenario where concurrency
// matches available CPU cores (what real servers typically use)
func BenchmarkConcurrentRealistic_RealWorld_Hyper(b *testing.B) {
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

// ============================================================================
// BENCHMARK: Go html/template (pre-parsed)
// ============================================================================

func BenchmarkSimpleElement_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		simpleElementTmpl.Execute(&buf, nil)
	}
}

func BenchmarkDeepNesting_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		deepNestingTmpl.Execute(&buf, nil)
	}
}

func BenchmarkManyAttributes_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		manyAttributesTmpl.Execute(&buf, nil)
	}
}

func BenchmarkLargeText_HtmlTemplate(b *testing.B) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		largeTextTmpl.Execute(&buf, text)
	}
}

func BenchmarkList10_HtmlTemplate(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		listTempl10Tmpl.Execute(&buf, items)
	}
}

func BenchmarkList100_HtmlTemplate(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = "item"
	}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		listTempl100Tmpl.Execute(&buf, items)
	}
}

func BenchmarkConditionals_HtmlTemplate(b *testing.B) {
	data := ConditionalsData{First: true, Second: false, Third: true, True: true}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		conditionalsTmpl.Execute(&buf, data)
	}
}

func BenchmarkMixedContent_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		mixedContentTmpl.Execute(&buf, nil)
	}
}

func BenchmarkVoidElements_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		voidElementsTmpl.Execute(&buf, nil)
	}
}

func BenchmarkHTMLEscaping_HtmlTemplate(b *testing.B) {
	content := "<script>alert('xss')</script> & more <b>bold</b>"
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		htmlEscapingTmpl.Execute(&buf, content)
	}
}

func BenchmarkTable_HtmlTemplate(b *testing.B) {
	rows := 10
	rowData := make([]struct{}, rows)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		tableTmpl.Execute(&buf, rowData)
	}
}

func BenchmarkForm_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		formTempl.Execute(&buf, nil)
	}
}

func BenchmarkEmptyPage_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		emptyPageTmpl.Execute(&buf, nil)
	}
}

func BenchmarkRawText_HtmlTemplate(b *testing.B) {
	html := template.HTML("<div><span>Content</span></div>")
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		rawTextTmpl.Execute(&buf, html)
	}
}

func BenchmarkSVG_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		svgTempl.Execute(&buf, nil)
	}
}

func BenchmarkRealWorld_HtmlTemplate(b *testing.B) {
	users := getBenchmarkData()
	data := RealWorldData{Users: users}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		realWorldTempl.Execute(&buf, data)
	}
}

func BenchmarkSequential_RealWorld_HtmlTemplate(b *testing.B) {
	users := getBenchmarkData()
	data := RealWorldData{Users: users}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		realWorldTempl.Execute(&buf, data)
	}
}

func BenchmarkConcurrent10_RealWorld_HtmlTemplate(b *testing.B) {
	users := getBenchmarkData()
	data := RealWorldData{Users: users}
	b.SetParallelism(10)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			realWorldTempl.Execute(&buf, data)
		}
	})
}

func BenchmarkConcurrent100_RealWorld_HtmlTemplate(b *testing.B) {
	users := getBenchmarkData()
	data := RealWorldData{Users: users}
	b.SetParallelism(100)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			realWorldTempl.Execute(&buf, data)
		}
	})
}

func BenchmarkConcurrent1000_RealWorld_HtmlTemplate(b *testing.B) {
	users := getBenchmarkData()
	data := RealWorldData{Users: users}
	b.SetParallelism(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			realWorldTempl.Execute(&buf, data)
		}
	})
}

func BenchmarkConcurrentRealistic_RealWorld_HtmlTemplate(b *testing.B) {
	users := getBenchmarkData()
	data := RealWorldData{Users: users}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			realWorldTempl.Execute(&buf, data)
		}
	})
}
