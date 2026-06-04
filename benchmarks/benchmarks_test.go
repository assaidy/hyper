package benchmarks

import (
	"bytes"
	"html/template"
	"testing"

	h "github.com/assaidy/hyper/v2"
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
	page := h.DIV()("Hello World")
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
	page := h.DIV()(
		h.DIV()(
			h.DIV()(
				h.DIV()(
					h.DIV()(
						h.P()("Deep content"),
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
	page := h.DIV(
		h.Attr("id", "main"),
		h.Attr("class", "container wrapper"),
		h.Attr("data-role", "content"),
		h.Attr("data-value", "12345"),
		h.Attr("aria-label", "Main content"),
		h.Attr("hidden", true),
		h.Attr("disabled", false),
	)()
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
	page := h.P()(text)
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
	page := h.UL()(
		h.Range(items, func(s string) any {
			return h.LI()(s)
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
	page := h.UL()(
		h.Range(items, func(s string) any {
			return h.LI()(s)
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
	cond := h.If(true, h.SPAN()("First"))
	cond = cond.ElseIf(false, h.SPAN()("Second"))
	page := h.DIV()(cond, h.IfElse(true, h.STRONG()("True"), h.EM()("False")))
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// ============================================================================
// BENCHMARK 8: Mixed Content

func BenchmarkMixedContent_Templ(b *testing.B) {
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		MixedContentTempl().Render(ctx, &buf)
	}
}

func BenchmarkMixedContent_Hyper(b *testing.B) {
	page := h.DIV()(
		h.H1()("Title"),
		h.P()("Paragraph with ", h.STRONG()("bold"), " and ", h.EM()("italic"), " text."),
		h.UL()(
			h.LI()("Item 1"),
			h.LI()(h.A(h.AttrHref("#"))("Link")),
		),
		h.DIV(h.AttrClass("footer"))(
			h.SMALL()("Copyright 2024"),
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
	page := h.DIV()(
		h.IMG(h.AttrSrc("image.jpg"), h.AttrAlt("Image")),
		h.BR(),
		h.HR(),
		h.INPUT(h.AttrType("text"), h.AttrValue("input")),
		h.META(h.AttrCharset("UTF-8")),
		h.LINK(h.AttrRel("stylesheet"), h.AttrHref("style.css")),
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
	page := h.DIV()(content)
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
	page := h.TABLE()(
		h.THEAD()(
			h.TR()(
				h.TH()("Name"),
				h.TH()("Value"),
				h.TH()("Action"),
			),
		),
		h.TBODY()(
			h.Repeat(rows, func() any {
				return h.TR()(
					h.TD()("Cell 1"),
					h.TD()("Cell 2"),
					h.TD()(h.BUTTON()("Click")),
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
	page := h.FORM(h.AttrAction("/submit"), h.AttrMethod("POST"))(
		h.FIELDSET()(
			h.LEGEND()("User Form"),
			h.LABEL(h.Attr("for", "name"))("Name:"),
			h.INPUT(h.AttrType("text"), h.AttrId("name"), h.AttrName("name")),
			h.BR(),
			h.LABEL(h.Attr("for", "email"))("Email:"),
			h.INPUT(h.AttrType("email"), h.AttrId("email"), h.AttrName("email")),
			h.BR(),
			h.BUTTON(h.AttrType("submit"))("Submit"),
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
	page := h.Group(
		h.DOCTYPE(),
		h.HTML()(
			h.HEAD()(
				h.META(h.AttrCharset("UTF-8")),
				h.TITLE()(h.RawText("User Dashboard")),
				h.LINK(h.AttrRel("stylesheet"), h.AttrHref("/style.css")),
			),
			h.BODY()(
				h.HEADER()(
					h.H1()("User Dashboard"),
					h.NAV()(
						h.A(h.AttrHref("/"))("Home"),
						h.A(h.AttrHref("/users"))("Users"),
						h.A(h.AttrHref("/settings"))("Settings"),
					),
				),
				h.MAIN()(
					h.H2()("Users"),
					h.If(len(users) > 0,
						h.TABLE()(
							h.THEAD()(
								h.TR()(
									h.TH()("Name"),
									h.TH()("Role"),
								),
							),
							h.TBODY()(
								h.Range(users, func(u User) any {
									return h.TR()(
										h.TD()(u.Name),
										h.TD()(h.IfElse(u.Admin, h.STRONG()("Admin"), h.SPAN()("User"))),
									)
								}),
							),
						),
					),
					h.If(len(users) == 0, h.P()("No users found.")),
				),
				h.FOOTER()(
					h.P()("© 2024 Company"),
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
	page := h.HTML()(h.BODY()())
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
	page := h.DIV()(h.RawText(html))
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

func BenchmarkRegularString_Hyper(b *testing.B) {
	text := "<div><span>Content</span></div>"
	page := h.DIV()(text)
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
	page := h.SVG(h.AttrWidth("100"), h.AttrHeight("100"))(
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
	return h.Group(
		h.DOCTYPE(),
		h.HTML()(
			h.HEAD()(
				h.META(h.AttrCharset("UTF-8")),
				h.META(h.AttrName("viewport"), h.Attr("content", "width=device-width, initial-scale=1.0")),
				h.TITLE()(h.RawText("Dashboard - User Management")),
				h.LINK(h.AttrRel("stylesheet"), h.AttrHref("/css/main.css")),
				h.LINK(h.AttrRel("icon"), h.AttrHref("/favicon.ico")),
			),
			h.BODY()(
				h.HEADER()(
					h.NAV(h.AttrClass("main-nav"))(
						h.A(h.AttrHref("/"), h.AttrClass("nav-link"))("Home"),
						h.A(h.AttrHref("/users"), h.AttrClass("nav-link active"))("Users"),
						h.A(h.AttrHref("/settings"), h.AttrClass("nav-link"))("Settings"),
						h.A(h.AttrHref("/logout"), h.AttrClass("nav-link"))("Logout"),
					),
				),
				h.MAIN(h.AttrClass("main-content"))(
					h.H1()("User Management Dashboard"),
					h.P()("Welcome to the admin dashboard. Manage users and permissions below."),
					h.If(len(users) > 0,
						h.SECTION(h.AttrClass("users-section"))(
							h.H2()("Active Users"),
							h.TABLE(h.AttrClass("users-table"))(
								h.THEAD()(
									h.TR()(
										h.TH()("ID"),
										h.TH()("Name"),
										h.TH()("Role"),
										h.TH()("Status"),
										h.TH()("Actions"),
									),
								),
								h.TBODY()(
									h.Range(users, func(u User) any {
										return h.TR()(
											h.TD()(h.STRONG()("#")),
											h.TD()(u.Name),
											h.TD()(h.IfElse(u.Admin,
												h.SPAN(h.AttrClass("badge admin"))("Administrator"),
												h.SPAN(h.AttrClass("badge user"))("User"),
											)),
											h.TD()(h.SPAN(h.AttrClass("status active"))("Active")),
											h.TD()(
												h.BUTTON(h.AttrClass("btn-edit"))("Edit"),
												h.BUTTON(h.AttrClass("btn-delete"))("Delete"),
											),
										)
									}),
								),
							),
						),
					),
					h.If(len(users) == 0,
						h.DIV(h.AttrClass("empty-state"))(
							h.P()("No users found. Add your first user to get started."),
						),
					),
					h.SECTION(h.AttrClass("quick-stats"))(
						h.H3()("Quick Stats"),
						h.DIV(h.AttrClass("stats-grid"))(
							h.DIV(h.AttrClass("stat-card"))(
								h.STRONG()(len(users)),
								h.SPAN()("Total Users"),
							),
							h.DIV(h.AttrClass("stat-card"))(
								h.STRONG()(h.IfElse(len(users) > 0, len(users), 0)),
								h.SPAN()("Active Now"),
							),
						),
					),
				),
				h.FOOTER(h.AttrClass("site-footer"))(
					h.P()("2025 Company Inc. All rights reserved."),
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
		simpleElementTempl.Execute(&buf, nil)
	}
}

func BenchmarkDeepNesting_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		deepNestingTempl.Execute(&buf, nil)
	}
}

func BenchmarkManyAttributes_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		manyAttributesTempl.Execute(&buf, nil)
	}
}

func BenchmarkLargeText_HtmlTemplate(b *testing.B) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		largeTextTempl.Execute(&buf, text)
	}
}

func BenchmarkList10_HtmlTemplate(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		listTempl10Templ.Execute(&buf, items)
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
		listTempl100Templ.Execute(&buf, items)
	}
}

func BenchmarkConditionals_HtmlTemplate(b *testing.B) {
	data := ConditionalsData{First: true, Second: false, Third: true, True: true}
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		conditionalsTempl.Execute(&buf, data)
	}
}

func BenchmarkMixedContent_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		mixedContentTempl.Execute(&buf, nil)
	}
}

func BenchmarkVoidElements_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		voidElementsTempl.Execute(&buf, nil)
	}
}

func BenchmarkHTMLEscaping_HtmlTemplate(b *testing.B) {
	content := "<script>alert('xss')</script> & more <b>bold</b>"
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		htmlEscapingTempl.Execute(&buf, content)
	}
}

func BenchmarkTable_HtmlTemplate(b *testing.B) {
	rows := 10
	rowData := make([]struct{}, rows)
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		tableTempl.Execute(&buf, rowData)
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
		emptyPageTempl.Execute(&buf, nil)
	}
}

func BenchmarkRawText_HtmlTemplate(b *testing.B) {
	html := template.HTML("<div><span>Content</span></div>")
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		rawTextTempl.Execute(&buf, html)
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

// ============================================================================
// BUILDER BENCHMARKS: Closure vs Struct+Method
// elementBuilder uses struct+method to avoid MakeChildrenInserter closure heap allocs
// ============================================================================

type elementBuilder struct {
	elem h.Element
}

func (b elementBuilder) With(children ...any) h.Element {
	b.elem.InsertChildren(children...)
	return b.elem
}

func elemDIV(attrs ...h.Attribute) elementBuilder {
	return elementBuilder{elem: h.Element{Tag: "div", Attributes: attrs}}
}

func elemP(attrs ...h.Attribute) elementBuilder {
	return elementBuilder{elem: h.Element{Tag: "p", Attributes: attrs}}
}

func elemNAV(attrs ...h.Attribute) elementBuilder {
	return elementBuilder{elem: h.Element{Tag: "nav", Attributes: attrs}}
}

func elemA(attrs ...h.Attribute) elementBuilder {
	return elementBuilder{elem: h.Element{Tag: "a", Attributes: attrs}}
}

func elemUL(attrs ...h.Attribute) elementBuilder {
	return elementBuilder{elem: h.Element{Tag: "ul", Attributes: attrs}}
}

func elemLI(attrs ...h.Attribute) elementBuilder {
	return elementBuilder{elem: h.Element{Tag: "li", Attributes: attrs}}
}

// ============================================================================
// SMALL TREE: div > p > text
// ============================================================================

func BenchmarkSmallTree_Closure_ConstructOnly(b *testing.B) {
	for b.Loop() {
		h.DIV(h.AttrClass("foo"))(h.P()("hello"))
	}
}

func BenchmarkSmallTree_Closure_ConstructAndRender(b *testing.B) {
	for b.Loop() {
		page := h.DIV(h.AttrClass("foo"))(h.P()("hello"))
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

func BenchmarkSmallTree_Struct_ConstructOnly(b *testing.B) {
	for b.Loop() {
		elemDIV(h.AttrClass("foo")).With(elemP().With("hello"))
	}
}

func BenchmarkSmallTree_Struct_ConstructAndRender(b *testing.B) {
	for b.Loop() {
		page := elemDIV(h.AttrClass("foo")).With(elemP().With("hello"))
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// ============================================================================
// MEDIUM TREE: div > nav > 3 links
// ============================================================================

func BenchmarkMediumTree_Closure_ConstructOnly(b *testing.B) {
	for b.Loop() {
		h.DIV(h.AttrClass("nav"))(
			h.NAV()(
				h.A(h.AttrHref("/"))("Home"),
				h.A(h.AttrHref("/users"))("Users"),
				h.A(h.AttrHref("/about"))("About"),
			),
		)
	}
}

func BenchmarkMediumTree_Closure_ConstructAndRender(b *testing.B) {
	for b.Loop() {
		page := h.DIV(h.AttrClass("nav"))(
			h.NAV()(
				h.A(h.AttrHref("/"))("Home"),
				h.A(h.AttrHref("/users"))("Users"),
				h.A(h.AttrHref("/about"))("About"),
			),
		)
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

func BenchmarkMediumTree_Struct_ConstructOnly(b *testing.B) {
	for b.Loop() {
		elemDIV(h.AttrClass("nav")).With(
			elemNAV().With(
				elemA(h.AttrHref("/")).With("Home"),
				elemA(h.AttrHref("/users")).With("Users"),
				elemA(h.AttrHref("/about")).With("About"),
			),
		)
	}
}

func BenchmarkMediumTree_Struct_ConstructAndRender(b *testing.B) {
	for b.Loop() {
		page := elemDIV(h.AttrClass("nav")).With(
			elemNAV().With(
				elemA(h.AttrHref("/")).With("Home"),
				elemA(h.AttrHref("/users")).With("Users"),
				elemA(h.AttrHref("/about")).With("About"),
			),
		)
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// ============================================================================
// ALLOCATION BREAKDOWN: where are the allocs going?
// ============================================================================

func BenchmarkConstructRealWorld_Hyper(b *testing.B) {
	for b.Loop() {
		buildRealWorldPage(getBenchmarkData())
	}
}

func BenchmarkConstructAndRenderRealWorld_Hyper(b *testing.B) {
	for b.Loop() {
		page := buildRealWorldPage(getBenchmarkData())
		var buf bytes.Buffer
		h.Render(&buf, page)
	}
}

// ============================================================================
// BIG STATIC PAGE: templ vs html/template vs hyper (no cache) vs hyper_cache vs hyper_once
// Compare a fully static marketing page rendered via different caching strategies.
// Hyper variants rebuild the tree each iteration to simulate per-request tree construction.
// ============================================================================

func buildBigStaticPage() h.HyperNode {
	return h.Group(
		h.DOCTYPE(),
		h.HTML()(
			h.HEAD()(
				h.META(h.AttrCharset("UTF-8")),
				h.META(h.AttrName("viewport"), h.Attr("content", "width=device-width, initial-scale=1.0")),
				h.TITLE()(h.RawText("MyCompany - Big Static Page")),
				h.LINK(h.AttrRel("stylesheet"), h.AttrHref("/style.css")),
				h.LINK(h.AttrRel("preconnect"), h.AttrHref("https://fonts.googleapis.com")),
				h.LINK(h.AttrRel("icon"), h.AttrHref("/favicon.ico")),
			),
			h.BODY()(
				h.HEADER(h.AttrClass("site-header"))(
					h.DIV(h.AttrClass("container"))(
						h.DIV(h.AttrClass("logo"))(
							h.A(h.AttrHref("/"), h.AttrClass("logo-link"))(
								h.IMG(h.AttrSrc("/logo.svg"), h.AttrAlt("Logo")),
								h.SPAN()("MyCompany"),
							),
						),
						h.NAV(h.AttrClass("main-nav"))(
							h.UL()(
								h.LI()(h.A(h.AttrHref("/"))("Home")),
								h.LI()(h.A(h.AttrHref("/features"))("Features")),
								h.LI()(h.A(h.AttrHref("/pricing"))("Pricing")),
								h.LI()(h.A(h.AttrHref("/about"))("About")),
								h.LI()(h.A(h.AttrHref("/contact"))("Contact")),
								h.LI()(h.A(h.AttrHref("/blog"))("Blog")),
							),
						),
						h.DIV(h.AttrClass("auth-buttons"))(
							h.A(h.AttrHref("/login"), h.AttrClass("btn btn-outline"))("Log In"),
							h.A(h.AttrHref("/signup"), h.AttrClass("btn btn-primary"))("Sign Up"),
						),
					),
				),
				h.MAIN()(
					h.SECTION(h.AttrClass("hero"))(
						h.DIV(h.AttrClass("hero-content"))(
							h.H1()("Welcome to MyCompany"),
							h.P()("The all-in-one platform for modern teams. Build faster, collaborate smarter, and deliver better results."),
							h.DIV(h.AttrClass("hero-cta"))(
								h.A(h.AttrHref("/signup"), h.AttrClass("btn btn-primary btn-large"))("Get Started Free"),
								h.A(h.AttrHref("/demo"), h.AttrClass("btn btn-outline btn-large"))("Watch Demo"),
							),
							h.DIV(h.AttrClass("hero-stats"))(
								h.DIV(h.AttrClass("stat"))(h.STRONG()("10K+"), h.SPAN()("Active Users")),
								h.DIV(h.AttrClass("stat"))(h.STRONG()("99.9%"), h.SPAN()("Uptime")),
								h.DIV(h.AttrClass("stat"))(h.STRONG()("150+"), h.SPAN()("Countries")),
							),
						),
					),
					h.SECTION(h.AttrClass("features"))(
						h.DIV(h.AttrClass("container"))(
							h.H2()("Why Choose MyCompany"),
							h.DIV(h.AttrClass("feature-grid"))(
								h.DIV(h.AttrClass("feature-card"))(
									h.DIV(h.AttrClass("feature-icon"))(h.RawText("🚀")),
									h.H3()("Lightning Fast"),
									h.P()("Optimized performance with sub-millisecond response times and global CDN distribution for your content."),
								),
								h.DIV(h.AttrClass("feature-card"))(
									h.DIV(h.AttrClass("feature-icon"))(h.RawText("🔒")),
									h.H3()("Enterprise Security"),
									h.P()("Bank-grade encryption, SOC 2 compliance, and advanced threat detection to keep your data safe."),
								),
								h.DIV(h.AttrClass("feature-card"))(
									h.DIV(h.AttrClass("feature-icon"))(h.RawText("🎯")),
									h.H3()("Smart Analytics"),
									h.P()("Real-time insights and AI-powered recommendations to help you make data-driven decisions."),
								),
								h.DIV(h.AttrClass("feature-card"))(
									h.DIV(h.AttrClass("feature-icon"))(h.RawText("🌐")),
									h.H3()("Global Scale"),
									h.P()("Deploy to 30+ regions worldwide with automatic scaling and built-in disaster recovery."),
								),
							),
						),
					),
					h.SECTION(h.AttrClass("pricing"))(
						h.DIV(h.AttrClass("container"))(
							h.H2()("Simple, Transparent Pricing"),
							h.DIV(h.AttrClass("pricing-grid"))(
								h.DIV(h.AttrClass("pricing-card"))(
									h.H3()("Starter"),
									h.P(h.AttrClass("price"))(h.RawText("$9"), h.SPAN()("/month")),
									h.UL()(
										h.LI()("Up to 5 projects"),
										h.LI()("10GB storage"),
										h.LI()("Basic analytics"),
										h.LI()("Email support"),
										h.LI()("1 team member"),
									),
									h.A(h.AttrHref("/signup"), h.AttrClass("btn btn-outline"))("Choose Plan"),
								),
								h.DIV(h.AttrClass("pricing-card featured"))(
									h.DIV(h.AttrClass("badge"))("Popular"),
									h.H3()("Professional"),
									h.P(h.AttrClass("price"))(h.RawText("$29"), h.SPAN()("/month")),
									h.UL()(
										h.LI()("Unlimited projects"),
										h.LI()("100GB storage"),
										h.LI()("Advanced analytics"),
										h.LI()("Priority support"),
										h.LI()("10 team members"),
										h.LI()("Custom integrations"),
									),
									h.A(h.AttrHref("/signup"), h.AttrClass("btn btn-primary"))("Choose Plan"),
								),
								h.DIV(h.AttrClass("pricing-card"))(
									h.H3()("Enterprise"),
									h.P(h.AttrClass("price"))(h.RawText("$99"), h.SPAN()("/month")),
									h.UL()(
										h.LI()("Unlimited projects"),
										h.LI()("1TB storage"),
										h.LI()("Enterprise analytics"),
										h.LI()("24/7 phone support"),
										h.LI()("Unlimited team members"),
										h.LI()("Custom integrations"),
										h.LI()("Dedicated account manager"),
										h.LI()("SLA guarantee"),
									),
									h.A(h.AttrHref("/contact"), h.AttrClass("btn btn-outline"))("Contact Sales"),
								),
							),
						),
					),
				),
				h.FOOTER(h.AttrClass("site-footer"))(
					h.DIV(h.AttrClass("container"))(
						h.DIV(h.AttrClass("footer-grid"))(
							h.DIV(h.AttrClass("footer-col"))(
								h.H4()("Product"),
								h.UL()(
									h.LI()(h.A(h.AttrHref("/features"))("Features")),
									h.LI()(h.A(h.AttrHref("/pricing"))("Pricing")),
									h.LI()(h.A(h.AttrHref("/integrations"))("Integrations")),
									h.LI()(h.A(h.AttrHref("/changelog"))("Changelog")),
								),
							),
							h.DIV(h.AttrClass("footer-col"))(
								h.H4()("Company"),
								h.UL()(
									h.LI()(h.A(h.AttrHref("/about"))("About")),
									h.LI()(h.A(h.AttrHref("/blog"))("Blog")),
									h.LI()(h.A(h.AttrHref("/careers"))("Careers")),
									h.LI()(h.A(h.AttrHref("/press"))("Press")),
								),
							),
							h.DIV(h.AttrClass("footer-col"))(
								h.H4()("Support"),
								h.UL()(
									h.LI()(h.A(h.AttrHref("/docs"))("Documentation")),
									h.LI()(h.A(h.AttrHref("/help"))("Help Center")),
									h.LI()(h.A(h.AttrHref("/status"))("Status")),
									h.LI()(h.A(h.AttrHref("/contact"))("Contact")),
								),
							),
							h.DIV(h.AttrClass("footer-col"))(
								h.H4()("Legal"),
								h.UL()(
									h.LI()(h.A(h.AttrHref("/privacy"))("Privacy Policy")),
									h.LI()(h.A(h.AttrHref("/terms"))("Terms of Service")),
									h.LI()(h.A(h.AttrHref("/cookies"))("Cookie Policy")),
									h.LI()(h.A(h.AttrHref("/gdpr"))("GDPR")),
								),
							),
						),
						h.DIV(h.AttrClass("footer-bottom"))(
							h.P()("© 2025 MyCompany Inc. All rights reserved."),
						),
					),
				),
			),
		),
	)
}

func BenchmarkBigStaticPage_Templ(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		BigStaticPageTempl().Render(b.Context(), &buf)
	}
}

func BenchmarkBigStaticPage_HtmlTemplate(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		bigStaticPageTempl.Execute(&buf, nil)
	}
}

func BenchmarkBigStaticPage_HyperNoCache(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, buildBigStaticPage())
	}
}

func BenchmarkBigStaticPage_HyperOnce(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, h.Once(func() h.HyperNode {
			return buildBigStaticPage()
		}))
	}
}

func BenchmarkBigStaticPage_HyperOnceWithKey(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		h.Render(&buf, h.OnceWithKey("key", func() h.HyperNode {
			return buildBigStaticPage()
		}))
	}
}

func BenchmarkOnceKey_WithKey(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		h.OnceWithKey("static-key", func() h.HyperNode { return h.DIV()("hello") })
	}
}
