package g_vs_templ

import (
	"bytes"
	"testing"

	"github.com/assaidy/g"
)

func getBenchmarkData() []User {
	return []User{
		{Name: "Alice", Admin: true},
		{Name: "Bob", Admin: false},
		{Name: "Charlie", Admin: false},
		{Name: "David", Admin: true},
		{Name: "Eve", Admin: false},
	}
}

// BenchmarkTempl benchmarks the rendering performance using templ library
func BenchmarkTempl(b *testing.B) {
	users := getBenchmarkData()
	ctx := b.Context()
	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer
		Page(users).Render(ctx, &buf)
	}
}

// BenchmarkG benchmarks the rendering performance using gg library
func BenchmarkG(b *testing.B) {
	users := getBenchmarkData()
	page := gg.Empty(
		gg.DoctypeHTML(),
		gg.Html(
			gg.Head(
				gg.Title(gg.RawHTML("Benchmark Page")),
			),
			gg.Body(
				gg.H1(gg.RawHTML("Users List")),
				gg.Ul(
					gg.MapSlice(users, func(u User) gg.Node {
						return gg.Li(
							u.Name,
							gg.If(u.Admin, gg.Span(" (Admin)")),
						)
					}),
				),
				gg.If(len(users) == 0, gg.P("No users found.")),
				gg.Div(
					gg.Repeat(5, func() gg.Node {
						return gg.Hr()
					}),
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
